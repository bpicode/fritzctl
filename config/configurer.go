package config

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path"

	"github.com/bpicode/fritzctl/internal/console"
	"gopkg.in/yaml.v2"
)

// ExtendedConfig contains the fritz core config along with
// other data (like config file location).
type ExtendedConfig struct {
	fritzCfg Config
	file     string
}

// Configurer provides functions to obtain user data from
// stdin and write the result to a file.
type Configurer interface {
	Greet()
	Obtain(r io.Reader) (ExtendedConfig, error)
}

// NewConfigurer creates a Configurer instance.
func NewConfigurer() Configurer {
	return &cliConfigurer{}
}

type cliConfigurer struct {
}

// Greet prints a small greeting.
func (iCLI *cliConfigurer) Greet() {
	fmt.Println("Configure fritzctl: hit [ENTER] to accept the default value, hit [^C] to abort")
}

// Obtain starts the dialog session, asking for the values to fill
// an ExtendedConfig.
func (iCLI *cliConfigurer) Obtain(r io.Reader) (ExtendedConfig, error) {
	f := ""
	c := Config{}
	err := proceedUntilFirstError(
		func() (err error) {
			f, err = iCLI.obtainFileLocation(r)
			return err
		},
		func() (err error) {
			c.Net, err = iCLI.obtainNetConfig(r)
			return err
		},
		func() (err error) {
			c.Login, err = iCLI.obtainLoginConfig(r)
			return err
		},
		func() (err error) {
			c.Pki, err = iCLI.obtainPkiConfig(r)
			return err
		},
	)
	return ExtendedConfig{file: f, fritzCfg: c}, err
}

func proceedUntilFirstError(fs ...func() error) error {
	for _, f := range fs {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}

// Write writes the user data to the configured file.
func (c *ExtendedConfig) Write() error {
	dir, _ := path.Split(c.file)
	os.MkdirAll(dir, 0700)
	f, err := os.OpenFile(c.file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	encoder := yaml.NewEncoder(f)
	defer encoder.Close()
	return encoder.Encode(c.fritzCfg)
}

func (iCLI *cliConfigurer) obtainFileLocation(r io.Reader) (string, error) {
	f := struct{ File string }{}
	s := console.Survey{In: r, Out: os.Stdout}
	l := iCLI.defaultConfigLocation(homeDirOf(user.Current))
	err := s.Ask(
		[]console.Question{console.ForString("file", "Config file location", l)},
		&f)
	return f.File, err
}

func (iCLI *cliConfigurer) defaultConfigLocation(pathFunc func(string) (string, error)) string {
	loc, err := pathFunc(path.Join(".fritzctl", "config.yml"))
	if err != nil {
		return "fritzctl.yml"
	}
	return loc
}

func (iCLI *cliConfigurer) obtainNetConfig(r io.Reader) (*Net, error) {
	netCfg := Net{}
	survey := console.Survey{In: r, Out: os.Stdout}
	err := survey.Ask(
		[]console.Question{
			console.ForString("protocol", "Communication protocol", "https"),
			console.ForString("host", "Hostname/IP", "fritz.box"),
			console.ForString("port", "Port", ""),
		}, &netCfg)
	return &netCfg, err
}

func (iCLI *cliConfigurer) obtainLoginConfig(r io.Reader) (*Login, error) {
	login := Login{}
	survey := console.Survey{In: r, Out: os.Stdout}
	err := survey.Ask(
		[]console.Question{
			console.ForString("loginURL", "Login path", "/login_sid.lua"),
			console.ForString("username", "Username", ""),
			console.ForPassword("password", "Password"),
		}, &login)
	return &login, err
}

func (iCLI *cliConfigurer) obtainPkiConfig(r io.Reader) (*Pki, error) {
	pki := Pki{}
	survey := console.Survey{In: r, Out: os.Stdout}
	err := survey.Ask(
		[]console.Question{
			console.ForBool("skipTlsVerify", "Skip TLS certificate validation", false),
			console.ForString("certificateFile", "Path to PEM-formatted certificate file", ""),
		}, &pki)
	return &pki, err
}
