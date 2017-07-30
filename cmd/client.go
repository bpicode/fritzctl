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
	assertNoError(err, "unable to create FRITZ!Box client:", err)
	client, err := fritz.NewClient(configFile)
	assertNoError(err, "unable to create FRITZ!Box client:", err)
	err = client.Login()
	assertNoError(err, "unable to login:", err)
	return client
}

func homeAutoClient() fritz.HomeAuto {
	opts := findOptions(config.FindConfigFile)
	h := fritz.NewHomeAuto(opts...)
	err := h.Login()
	assertNoError(err, "unable to login:", err)
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
	assertNoError(err, "cannot apply configuration:", err)
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
		assertNoError(err, "cannot read certificate file:", err)
		opts = append(opts, fritz.Certificate(bs))
	}
	return opts

}

func loginOptions(opts []fritz.Option, login *config.Login) []fritz.Option {
	opts = append(opts, fritz.Credentials(login.Username, login.Password))
	opts = append(opts, fritz.AuthEndpoint(login.LoginURL))
	return opts
}
