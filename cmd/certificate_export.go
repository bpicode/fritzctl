package cmd

import (
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/bpicode/fritzctl/config"
	"github.com/bpicode/fritzctl/internal/console"
	"github.com/bpicode/fritzctl/internal/stringutils"
	"github.com/spf13/cobra"
)

var certExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export the TLS certificate offered by the FRITZ!Box",
	Long: `Export TLS certificate offered by the FRITZ!Box to stdout, encoded in PEM format.
The output can be redirected to a file or copied from the terminal output.
Use the data beginning with "-----BEGIN CERTIFICATE-----" and ending with "-----END CERTIFICATE-----".`,
	Example: "fritzctl certificate export > /path/to/certificate.pem",
	RunE:    certExport,
}

func init() {
	certCmd.AddCommand(certExportCmd)
}

func certExport(_ *cobra.Command, _ []string) error {
	cfg := mustReadConfig()
	conn := mustConnect(cfg)
	defer conn.Close()
	crt := mustHaveCert(conn)
	logCrt(&printableCert{crt}, os.Stderr)
	writeCrt(crt, os.Stdout)
	return nil
}

func mustReadConfig() *config.Config {
	c, err := cfg(defaultConfigPlaces...)
	assertNoErr(err, "cannot parse configuration")
	return c
}

func mustConnect(cfg *config.Config) *tls.Conn {
	conn, err := tls.Dial("tcp", connectionString(cfg), &tls.Config{InsecureSkipVerify: true})
	assertNoErr(err, "certificate export failed, cannot connect to remote")
	return conn
}

func mustHaveCert(conn *tls.Conn) *x509.Certificate {
	state := conn.ConnectionState()
	assertTrue(len(state.PeerCertificates) > 0, errors.New("certificate export failed, list of peer certificates is empty"))
	crt := state.PeerCertificates[len(state.PeerCertificates)-1]
	return crt
}

func writeCrt(crt *x509.Certificate, w io.Writer) {
	p := pem.Block{Type: "CERTIFICATE", Bytes: crt.Raw}
	pem.Encode(w, &p)
}

func logCrt(crt *printableCert, w io.Writer) {
	fmt.Fprintf(w, "%s %s\n", console.Cyan("Serial Number:      "), crt.serial())
	fmt.Fprintf(w, "%s %s\n", console.Cyan("Issued to:          "), crt.issuedTo())
	fmt.Fprintf(w, "%s %s\n", console.Cyan("Validity:           "), crt.validity())
	fmt.Fprintf(w, "%s %s\n", console.Cyan("SHA-256 fingerprint:"), crt.fingerprintSHA256())
}

func connectionString(cfg *config.Config) string {
	return fmt.Sprintf("%s:%s", cfg.Host, stringutils.DefaultIfEmpty(cfg.Port, "443"))
}

type printableCert struct {
	*x509.Certificate
}

func (c *printableCert) serial() string {
	sn := c.SerialNumber.Bytes
	return hex.EncodeToString(sn())
}

func (c *printableCert) fingerprintSHA256() string {
	sum256 := sha256.Sum256(c.Raw)
	return hex.EncodeToString(sum256[:])
}

func (c *printableCert) issuedTo() string {
	return fmt.Sprintf("CN='%s' O='%s'", c.Subject.CommonName, c.Subject.Organization)
}

func (c *printableCert) validity() string {
	return fmt.Sprintf("from %s to %s", c.NotBefore, c.NotAfter)
}
