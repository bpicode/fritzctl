package fritz

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/bpicode/fritzctl/config"
	"github.com/stretchr/testify/assert"
)

// TestURLAllFeatures test if urlBuilderImpl correctly returns a URL when all features are used.
func TestURLAllFeatures(t *testing.T) {
	s := newURLBuilder(&config.Config{Net: &config.Net{Protocol: "https", Host: "192.168.127.4", Port: "4443"}}).query("key", "val ue").path("/alpha", "/beta").path("/gamma").query("key2", "value2").build()
	assert.Contains(t, s, "key")
	assert.Contains(t, s, "val+ue")
	assert.Contains(t, s, "key=val+ue")
	assert.Contains(t, s, "key2=value2")
	assert.Contains(t, s, "/alpha/beta/gamma")
	assert.Contains(t, s, "https://192.168.127.4:4443/alpha/beta/gamma?")
	u, err := url.Parse(s)
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Contains(t, u.RawQuery, "key")
	assert.Contains(t, u.RawQuery, "val+ue")
	assert.Contains(t, u.RawQuery, "key=val+ue")
	assert.Contains(t, u.RawQuery, "key2=value2")
}

// TestURLBuilder test if urlBuilderImpl correctly returns URLs.
func TestURLBuilder(t *testing.T) {
	testCases := []fritzURLBuilder{
		newURLBuilder(&config.Config{Net: &config.Net{Protocol: "https", Host: "192.168.127.4"}}),
		newURLBuilder(&config.Config{Net: &config.Net{Protocol: "https", Host: "192.168.127.4"}}).query("key", "value"),
		newURLBuilder(&config.Config{Net: &config.Net{Protocol: "https", Host: "192.168.127.4", Port: "443"}}).query("key", "value"),
		newURLBuilder(&config.Config{Net: &config.Net{Protocol: "https", Host: "192.168.127.4", Port: "443"}}).query("key", "value").path("a"),
	}
	for i, testcase := range testCases {
		t.Run(fmt.Sprintf("Test url builder %d", i), func(t *testing.T) {
			s := testcase.build()
			assert.NotNil(t, s)
			assert.Contains(t, s, "https")
			assert.Contains(t, s, "192.168.127.4")
			u, err := url.Parse(s)
			assert.NoError(t, err)
			assert.NotNil(t, u)
		})
	}
}
