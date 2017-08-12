package fritz

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestApiUsage tests how the API feels.
func TestApiUsage(t *testing.T) {
	assertions := assert.New(t)
	u, _ := url.Parse("https://localhost:12345")
	h := NewHomeAuto(
		Credentials("", "password"),
		SkipTLSVerify(),
		Certificate([]byte{}),
		URL(u),
		AuthEndpoint("/login_sid.lua"),
	)
	assertions.NotNil(h)
}

// TestNoPanic tests that the API functions should not panic.
func TestNoPanic(t *testing.T) {
	assertions := assert.New(t)
	u, _ := url.Parse("https://localhost:12345")
	h := NewHomeAuto(
		Credentials("", "password"),
		SkipTLSVerify(),
		URL(u),
	)
	assertions.NotPanics(func() {
		h.Login()
		h.List()
		h.Temp(20.0, "dev_name")
		h.On("dev_name")
		h.Off("dev_name")
		h.Toggle("dev_name")
	})

}
