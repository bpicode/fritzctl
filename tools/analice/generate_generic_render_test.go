package main

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_logErrs(t *testing.T) {
	cmd := genericRenderCmd{}
	assert.NotPanics(t, func() {
		cmd.logErrs([]error{nil, errors.New("an error"), nil})
	})
}

type faultyScanner struct {
}

func (faultyScanner) scan(opt depScannerOptions, patterns ...string) ([]pkg, error) {
	return nil, errors.New("not working")
}

func Test_genericRenderCmd_run_errors_scan(t *testing.T) {
	cmd := genericRenderCmd{scanner: faultyScanner{}}
	err := cmd.run(&cobra.Command{}, []string{"./..."})
	assert.Error(t, err)
}
