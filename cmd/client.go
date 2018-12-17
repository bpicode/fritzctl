package cmd

import (
	"io/ioutil"
	"net/url"
	"os/user"

	"github.com/bpicode/fritzctl/config"
	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/logger"
)

var defaultConfigPlaces = []config.Place{
	config.InHomeDir(user.Current, ".fritzctl/config.yml", config.YAML()),
	config.InHomeDir(user.Current, ".fritzctl/config.json", config.JSON()),
	config.InHomeDir(user.Current, ".fritzctl/fritzctl.yml", config.YAML()),
	config.InHomeDir(user.Current, ".fritzctl/fritzctl.json", config.JSON()),
	config.InHomeDir(user.Current, ".fritzctl.yml", config.YAML()),
	config.InHomeDir(user.Current, ".fritzctl.json", config.JSON()),
	config.InDir("", "fritzctl.yml", config.YAML()),
	config.InDir("", "fritzctl.json", config.JSON()),
	config.InDir("", ".fritzctl.yml", config.YAML()),
	config.InDir("", ".fritzctl.json", config.JSON()),
	config.InDir("/etc/fritzctl", "config.yml", config.YAML()),
	config.InDir("/etc/fritzctl", "config.json", config.JSON()),
	config.InDir("/etc/fritzctl", "fritzctl.yml", config.YAML()),
	config.InDir("/etc/fritzctl", "fritzctl.json", config.JSON()),
}

func clientLogin() *fritz.Client {
	conf, err := cfg(defaultConfigPlaces...)
	assertNoErr(err, "cannot parse configuration")
	client := fritz.NewClientFromConfig(conf)
	err = client.Login()
	assertNoErr(err, "login failed")
	return client
}

func homeAutoClient(overrides ...fritz.Option) fritz.HomeAuto {
	opts := optsFromPlaces(defaultConfigPlaces...)
	opts = append(opts, overrides...)
	h := fritz.NewHomeAuto(opts...)
	err := h.Login()
	assertNoErr(err, "login failed")
	return h
}

func optsFromPlaces(places ...config.Place) []fritz.Option {
	opts := make([]fritz.Option, 0)
	cfg, err := cfg(places...)
	if err != nil {
		logger.Warn("Using default configuration because no config file could be inferred:", err)
		return make([]fritz.Option, 0)
	}
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

func cfg(places ...config.Place) (*config.Config, error) {
	p := config.NewParser(places...)
	return p.Parse()
}

func mustList() *fritz.Devicelist {
	c := homeAutoClient()
	devs, err := c.List()
	assertNoErr(err, "cannot obtain device data")
	return devs
}
