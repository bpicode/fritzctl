package fritz

import (
	"fmt"

	"strings"

	"github.com/bpicode/fritzctl/config"
	"github.com/bpicode/fritzctl/stringutils"
)

const (
	homeAutomationURI = "/webservices/homeautoswitch.lua"
	queryURI          = "/query.lua"
)

type fritzURLBuilder interface {
	query(key, value string) fritzURLBuilder
	path(ps ...string) fritzURLBuilder
	build() string
}

type fritzURLBuilderImpl struct {
	protocol    string
	host        string
	port        string
	queryParams map[string]string
	paths       []string
}

func (fb *fritzURLBuilderImpl) build() string {
	path := fmt.Sprintf("%s://%s%s", fb.protocol, fb.host, stringutils.DefaultIf(":"+fb.port, "", ":")) + strings.Replace("/"+strings.Join(fb.paths, "/"), "//", "/", -1)
	qs := "?" + strings.Join(stringutils.Contract(fb.queryParams, func(k, v string) string {
		return k + "=" + v
	}), "&")
	return path + stringutils.DefaultIf(qs, "", "?")
}

func (fb *fritzURLBuilderImpl) query(key, value string) fritzURLBuilder {
	fb.queryParams[key] = value
	return fb
}

func (fb *fritzURLBuilderImpl) path(ps ...string) fritzURLBuilder {
	fb.paths = append(fb.paths, ps...)
	return fb
}

func newURLBuilder(cfg *config.Config) fritzURLBuilder {
	return &fritzURLBuilderImpl{protocol: cfg.Net.Protocol, host: cfg.Net.Host, port: cfg.Net.Port, queryParams: map[string]string{}}
}
