package configurer

import (
	"fmt"

	"bufio"
	"os"

	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/meta"

	"encoding/json"

	"strconv"

	"github.com/bpicode/fritzctl/stringutils"
	"golang.org/x/crypto/ssh/terminal"
)

// ExtendedConfig contains the fritz core config along with
// other data (like config file location).
type ExtendedConfig struct {
	fritzCfg fritz.Config
	file     string
}

// InteractiveCLI provides funtions to obtain user data from
// stdin and write the result to a file.
type InteractiveCLI interface {
	InitWithDefaultVaules(cfg ExtendedConfig)
	Obtain() ExtendedConfig
	Write() error
}

// CLI creates an InteractiveCLI instance.
func CLI() InteractiveCLI {
	return &interactiveCLI{}
}

// Defaults constructs an ExtendedConfig with default values.
func Defaults() ExtendedConfig {
	return ExtendedConfig{
		file: meta.DefaultConfigFileAbsolute(),
		fritzCfg: fritz.Config{
			Protocol:       "https",
			Host:           "fritz.box",
			Password:       "",
			LoginURL:       "/login_sid.lua",
			Username:       "",
			SkipTLSVerify:  false,
			CerificateFile: "/etc/fritzctl/fritz.pem",
		}}
}

type interactiveCLI struct {
	defaultValues ExtendedConfig
	userValues    ExtendedConfig
}

// InitWithDefaultVaules backs the interactiveCLI with default values that
// will be applied if the user does not want to change the field.
func (iCLI *interactiveCLI) InitWithDefaultVaules(cfg ExtendedConfig) {
	iCLI.defaultValues = cfg
	iCLI.userValues = iCLI.defaultValues
}

// Obtain starts the dialog session, asking for the values to fill
// an ExtendedConfig.
func (iCLI *interactiveCLI) Obtain() ExtendedConfig {
	scanner := bufio.NewScanner(os.Stdin)
	iCLI.userValues.file = next(fmt.Sprintf("Enter config file location [%s]: ",
		iCLI.defaultValues.file), scanner, iCLI.defaultValues.file)
	iCLI.userValues.fritzCfg.Protocol = next(fmt.Sprintf("Enter FRITZ!Box communication protocol [%s]: ",
		iCLI.defaultValues.fritzCfg.Protocol), scanner, iCLI.defaultValues.fritzCfg.Protocol)
	iCLI.userValues.fritzCfg.Host = next(fmt.Sprintf("Enter FRITZ!Box hostname/ip(:port) [%s]: ",
		iCLI.defaultValues.fritzCfg.Host), scanner, iCLI.defaultValues.fritzCfg.Host)
	iCLI.userValues.fritzCfg.LoginURL = next(fmt.Sprintf("Enter FRITZ!Box login path [%s]: ",
		iCLI.defaultValues.fritzCfg.LoginURL), scanner, iCLI.defaultValues.fritzCfg.LoginURL)
	iCLI.userValues.fritzCfg.Username = next(fmt.Sprintf("Enter FRITZ!Box username [%s]: ",
		iCLI.defaultValues.fritzCfg.Username), scanner, iCLI.defaultValues.fritzCfg.Username)
	iCLI.userValues.fritzCfg.Password = nextCredential("Enter FRITZ!Box password: ", iCLI.defaultValues.fritzCfg.Password)
	fmt.Println()
	defaultSkipCert := strconv.FormatBool(iCLI.defaultValues.fritzCfg.SkipTLSVerify)
	doSkipCert := next(fmt.Sprintf("Skip TLS certificate validation [%s]: ", defaultSkipCert), scanner, defaultSkipCert)
	iCLI.userValues.fritzCfg.SkipTLSVerify, _ = strconv.ParseBool(doSkipCert)
	iCLI.userValues.fritzCfg.CerificateFile = next(fmt.Sprintf("Enter path to certificate file [%s]: ",
		iCLI.defaultValues.fritzCfg.CerificateFile), scanner, iCLI.defaultValues.fritzCfg.CerificateFile)
	return iCLI.userValues
}

// Write writes the user data to the configured file.
func (iCLI *interactiveCLI) Write() error {
	f, err := os.OpenFile(iCLI.userValues.file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	return encoder.Encode(&iCLI.userValues.fritzCfg)
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
	return stringutils.DefaultIfEmpty(string(pwBytes), defaultValue)
}
