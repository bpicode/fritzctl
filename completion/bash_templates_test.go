package completion

import (
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

// TestBashTemplate tests the bash template for validity.
func TestBashTemplate(t *testing.T) {
	tpl, err := template.New("completion.bash.template.test").Parse(bashCompletionTemplate)
	assert.NoError(t, err)
	assert.NotNil(t, tpl)
}
