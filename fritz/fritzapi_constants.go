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
	queryParams []string
	paths       []string
}

func (fb *fritzURLBuilderImpl) build() string {
	qs := "?" + strings.Join(fb.queryParams, "&")
	return fb.baseURL + stringutils.DefaultIf(qs, "", "?")
}

func (fb *fritzURLBuilderImpl) query(key, value string) fritzURLBuilder {
	fb.queryParams = append(fb.queryParams, fmt.Sprintf("%s=%s", key, value))
	return fb
}

func (fb *fritzURLBuilderImpl) path(ps ...string) fritzURLBuilder {
	fb.paths = append(fb.paths, ps...)
	return fb
}

func newURLBuilder(cfg *config.Config) fritzURLBuilder {
	return &fritzURLBuilderImpl{baseURL: fmt.Sprintf("%s://%s%s", cfg.Protocol, cfg.Host, stringutils.DefaultIf(":"+cfg.Port, "", ":"))}
}

func ain(a string) string {
	return fmt.Sprintf("ain=%s", a)
}
