package cmd

import (
	"io/ioutil"
	"net/url"

	"github.com/bpicode/fritzctl/config"
	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/logger"
)

func clientLogin() *fritz.Client {
	configFile, err := config.FindConfigFile()
	assertNoErr(err, "cannot find configuration file")
	client, err := fritz.NewClient(configFile)
	assertNoErr(err, "failed to create FRITZ!Box client")
	err = client.Login()
	assertNoErr(err, "login failed")
	return client
}

func homeAutoClient(overrides ...fritz.Option) fritz.HomeAuto {
	opts := findOptions(config.FindConfigFile)
	opts = append(opts, overrides...)
	h := fritz.NewHomeAuto(opts...)
	err := h.Login()
	assertNoErr(err, "login failed")
	return h
}

type cfgFileFinder func() (string, error)

func findOptions(finder cfgFileFinder) []fritz.Option {
	path, err := finder()
	if err != nil {
		logger.Warn("Using default configuration because no config file could be inferred:", err)
		return make([]fritz.Option, 0)
	}
	return fromFile(path)
}

func fromFile(path string) []fritz.Option {
	opts := make([]fritz.Option, 0)
	cfg, err := config.New(path)
	assertNoErr(err, "cannot apply configuration")
	opts = networkOptions(opts, cfg.Net)
	opts = certificateOptions(opts, cfg.Pki)
	opts = loginOptions(opts, cfg.Login)
	return opts
}

func networkOptions(opts []fritz.Option, net *config.Net) []fritz.Option {
	return append(opts, fritz.URL(&url.URL{Host: net.Host + ":" + net.Port, Scheme: net.Protocol}))
}

func certificateOptions(opts []fritz.Option, pki *config.Pki) []fritz.Option {
	if pki.SkipTLSVerify {
		opts = append(opts, fritz.SkipTLSVerify())
		return opts
	}
	if pki.CertificateFile != "" {
		bs, err := ioutil.ReadFile(pki.CertificateFile)
		assertNoErr(err, "cannot read certificate file")
		opts = append(opts, fritz.Certificate(bs))
	}
	return opts

}

func loginOptions(opts []fritz.Option, login *config.Login) []fritz.Option {
	opts = append(opts, fritz.Credentials(login.Username, login.Password))
	opts = append(opts, fritz.AuthEndpoint(login.LoginURL))
	return opts
}
