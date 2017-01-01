package configurer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/bpicode/fritzctl/config"
	"github.com/bpicode/fritzctl/meta"
	"github.com/bpicode/fritzctl/stringutils"
)

// ExtendedConfig contains the fritz core config along with
// other data (like config file location).
type ExtendedConfig struct {
	fritzCfg config.Config
	file     string
}

// Configurer provides funtions to obtain user data from
// stdin and write the result to a file.
type Configurer interface {
	Greet()
	ApplyDefaults(cfg ExtendedConfig)
	Obtain() ExtendedConfig
	Write() error
}

// New creates a Configurer instance.
func New() Configurer {
	return &cliConfigurer{}
}

// Defaults constructs an ExtendedConfig with default values.
func Defaults() ExtendedConfig {
	return ExtendedConfig{
		file: meta.DefaultConfigFileAbsolute(),
		fritzCfg: config.Config{
			Net: &config.Net{
				Protocol: "https",
				Host:     "fritz.box",
				Port:     "",
			},
			Login: &config.Login{
				Password: "",
				LoginURL: "/login_sid.lua",
				Username: "",
			},
			Pki: &config.Pki{
				SkipTLSVerify:   false,
				CertificateFile: "/etc/fritzctl/fritz.pem",
			},
		}}
}

type cliConfigurer struct {
	defaultValues ExtendedConfig
	userValues    ExtendedConfig
}

// Greet prints a small greeting.
func (iCLI *cliConfigurer) Greet() {
	fmt.Println("Configure fritzctl: hit [ENTER] to accept the default value, hit [^C] to abort")
}

// ApplyDefaults backs the cliConfigurer with default values that
// will be applied if the user does not want to change the field.
func (iCLI *cliConfigurer) ApplyDefaults(cfg ExtendedConfig) {
	iCLI.defaultValues = cfg
	iCLI.userValues = iCLI.defaultValues
}

// Obtain starts the dialog session, asking for the values to fill
// an ExtendedConfig.
func (iCLI *cliConfigurer) Obtain() ExtendedConfig {
	scanner := bufio.NewScanner(os.Stdin)
	iCLI.userValues.file = next(fmt.Sprintf("Config file location [%s]: ", iCLI.defaultValues.file), scanner, iCLI.defaultValues.file)
	iCLI.obtainNetConfig(scanner)
	iCLI.obtainLoginConfig(scanner)
	iCLI.obtainPkiConfig(scanner)
	return iCLI.userValues
}

// Write writes the user data to the configured file.
func (iCLI *cliConfigurer) Write() error {
	f, err := os.OpenFile(iCLI.userValues.file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	return encoder.Encode(struct {
		*config.Net
		*config.Login
		*config.Pki
	}{iCLI.userValues.fritzCfg.Net, iCLI.userValues.fritzCfg.Login, iCLI.userValues.fritzCfg.Pki})
}

func (iCLI *cliConfigurer) obtainNetConfig(scanner *bufio.Scanner) {
	iCLI.userValues.fritzCfg.Net.Protocol = next(fmt.Sprintf("FRITZ!Box communication protocol [%s]: ",
		iCLI.defaultValues.fritzCfg.Net.Protocol), scanner, iCLI.defaultValues.fritzCfg.Net.Protocol)
	iCLI.userValues.fritzCfg.Net.Host = next(fmt.Sprintf("FRITZ!Box hostname/ip [%s]: ",
		iCLI.defaultValues.fritzCfg.Net.Host), scanner, iCLI.defaultValues.fritzCfg.Net.Host)
	iCLI.userValues.fritzCfg.Net.Port = next(fmt.Sprintf("FRITZ!Box port [%s]: ",
		iCLI.defaultValues.fritzCfg.Net.Port), scanner, iCLI.defaultValues.fritzCfg.Net.Port)
}

func (iCLI *cliConfigurer) obtainLoginConfig(scanner *bufio.Scanner) {
	iCLI.userValues.fritzCfg.Login.LoginURL = next(fmt.Sprintf("FRITZ!Box login path [%s]: ",
		iCLI.defaultValues.fritzCfg.Login.LoginURL), scanner, iCLI.defaultValues.fritzCfg.Login.LoginURL)
	iCLI.userValues.fritzCfg.Login.Username = next(fmt.Sprintf("FRITZ!Box username [%s]: ",
		iCLI.defaultValues.fritzCfg.Login.Username), scanner, iCLI.defaultValues.fritzCfg.Login.Username)
	iCLI.userValues.fritzCfg.Login.Password = nextCredential("FRITZ!Box password: ", iCLI.defaultValues.fritzCfg.Login.Password)
}

func (iCLI *cliConfigurer) obtainPkiConfig(scanner *bufio.Scanner) {
	defaultSkipCert := strconv.FormatBool(iCLI.defaultValues.fritzCfg.Pki.SkipTLSVerify)
	doSkipCert := next(fmt.Sprintf("Skip TLS certificate validation [%s]: ", defaultSkipCert), scanner, defaultSkipCert)
	iCLI.userValues.fritzCfg.Pki.SkipTLSVerify, _ = strconv.ParseBool(doSkipCert)
	iCLI.userValues.fritzCfg.Pki.CertificateFile = next(fmt.Sprintf("Path to certificate file [%s]: ",
		iCLI.defaultValues.fritzCfg.Pki.CertificateFile), scanner, iCLI.defaultValues.fritzCfg.Pki.CertificateFile)
}

func next(prompt string, scanner *bufio.Scanner, defaultValue string) string {
	fmt.Print(prompt)
	scanner.Scan()
	val := scanner.Text()
	return stringutils.DefaultIfEmpty(val, defaultValue)
}

func nextCredential(prompt string, defaultValue string) string {
	fmt.Print(prompt)
	pwBytes, _ := terminal.ReadPassword(0)
	fmt.Println()
	return stringutils.DefaultIfEmpty(string(pwBytes), defaultValue)
}
