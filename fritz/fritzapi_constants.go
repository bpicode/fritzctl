package fritz

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/bpicode/fritzctl/config"
	"github.com/bpicode/fritzctl/internal/stringutils"
)

const (
	homeAutomationURI = "/webservices/homeautoswitch.lua"
	queryURI          = "/query.lua"
	inetStatURI       = "/internet/inetstat_monitor.lua"
	phoneListURI      = "/fon_num/foncalls_list.lua"
	systemStatusURI   = "/cgi-bin/system_status"
)

type fritzURLBuilder interface {
	query(key, value string) fritzURLBuilder
	path(ps ...string) fritzURLBuilder
	build() string
}

type fritzURLBuilderImpl struct {
	target      *target
	queryParams queryParams
	paths       []string
}

type target struct {
	protocol string
	host     string
	port     string
}

func (t *target) String() string {
	return fmt.Sprintf("%s://%s%s", t.protocol, t.host, stringutils.DefaultIf(":"+t.port, "", ":"))
}

type queryParams map[string]string

func (qs queryParams) String() string {
	joined := "?" + strings.Join(stringutils.Contract(qs, func(k, v string) string {
		return k + "=" + v
	}), "&")
	return stringutils.DefaultIf(joined, "", "?")
}

func (fb *fritzURLBuilderImpl) build() string {
	t := fb.target
	path := t.String() + strings.Replace("/"+strings.Join(fb.paths, "/"), "//", "/", -1)
	qs := fb.queryParams.String()
	return path + qs
}

func (fb *fritzURLBuilderImpl) query(key, value string) fritzURLBuilder {
	fb.queryParams[url.QueryEscape(key)] = url.QueryEscape(value)
	return fb
}

func (fb *fritzURLBuilderImpl) path(ps ...string) fritzURLBuilder {
	fb.paths = append(fb.paths, ps...)
	return fb
}

func newURLBuilder(cfg *config.Config) fritzURLBuilder {
	t := &target{protocol: cfg.Net.Protocol, host: cfg.Net.Host, port: cfg.Net.Port}
	return &fritzURLBuilderImpl{target: t, queryParams: queryParams{}}
}
