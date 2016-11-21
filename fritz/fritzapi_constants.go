package fritz

import (
	"fmt"

	"strings"

	"github.com/bpicode/fritzctl/config"
	"github.com/bpicode/fritzctl/stringutils"
)

const (
	homeAutomationURI = "/webservices/homeautoswitch.lua"
)

type fritzURLBuilder interface {
	query(key, value string) fritzURLBuilder
	path(ps ...string) fritzURLBuilder
	build() string
}

type fritzURLBuilderImpl struct {
	baseURL     string
	queryParams map[string]string
	paths       []string
}

func (fb *fritzURLBuilderImpl) build() string {
	path := strings.Replace(fb.baseURL+"/"+strings.Join(fb.paths, "/"), "//", "/", -1)
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
	return &fritzURLBuilderImpl{baseURL: fmt.Sprintf("%s://%s%s", cfg.Protocol, cfg.Host, stringutils.DefaultIf(":"+cfg.Port, "", ":")), queryParams: map[string]string{}}
}
