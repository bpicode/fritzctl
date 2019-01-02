package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/bpicode/fritzctl/internal/console"
	"github.com/stretchr/testify/assert"
)

func Test_bow_license_best_match(t *testing.T) {
	ble := bowEngine{}
	ble.initialize()

	for _, l := range licenses {
		bag := ble.bag(strings.NewReader(l.text))
		b, _ := ble.bestMatch(bag)
		assert.Equal(t, b, l)
	}
}

func Test_bow_license_best_match_not_good_enough(t *testing.T) {
	ble := bowEngine{}
	ble.initialize()
	ble.start()
	ble.analyze("this is an unknown license")
	err := ble.stop()
	assert.Error(t, err)
}

func Test_bow_license_matrix(t *testing.T) {
	ble := bowEngine{}
	ble.initialize()
	mat := calcMatrix(ble)
	lName := licShortNames()
	table := table(lName, mat)
	table.Print(os.Stdout)
}

func licShortNames() []string {
	lName := make([]string, len(licenses)+1)
	lName[0] = ""
	for i, l := range licenses {
		lName[i+1] = l.shortName
	}
	return lName
}

func calcMatrix(ble bowEngine) [][]float64 {
	mat := make([][]float64, len(licenses))
	for i := range mat {
		mat[i] = make([]float64, len(licenses))
	}
	for i, l1 := range licenses {
		b1 := ble.bag(strings.NewReader(l1.text))
		for j, l2 := range licenses {
			b2 := ble.bag(strings.NewReader(l2.text))
			mat[i][j] = b1.cos(b2)
		}
	}
	return mat
}

func table(lName []string, mat [][]float64) *console.Table {
	table := console.NewTable(console.Headers(lName...))
	for i, r := range mat {
		row := make([]string, len(r)+1)
		row[0] = licenses[i].shortName
		for j, s := range r {
			row[j+1] = fmt.Sprintf("%f", s)
		}
		table.Append(row)
	}
	return table
}
