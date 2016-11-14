package concurrent

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestScatterGatherNoWork test the ScatterGather when there is nothing to do.
func TestScatterGatherNoWork(t *testing.T) {

	work := map[string]func() (string, error){}

	ok := func(string, string) Result {
		return Result{Msg: "OK", Err: nil}
	}

	nok := func(string, string, error) Result {
		return Result{Msg: "", Err: errors.New("An error")}
	}

	results := ScatterGather(work, ok, nok)
	assert.NotNil(t, results)
	assert.Empty(t, results)
}

// TestScatterGatherAllOk test the ScatterGather where every
// goroutine succeeds.
func TestScatterGatherAllOk(t *testing.T) {

	work := map[string]func() (string, error){
		"1": func() (string, error) { return "1 says OK", nil },
		"2": func() (string, error) { return "2 says OK", nil },
		"3": func() (string, error) { return "3 says OK", nil },
		"4": func() (string, error) { return "4 says OK", nil },
		"5": func() (string, error) { return "5 says OK", nil },
		"6": func() (string, error) { return "6 says OK", nil },
		"7": func() (string, error) { return "7 says OK", nil },
		"8": func() (string, error) { return "8 says OK", nil },
		"9": func() (string, error) { return "9 says OK", nil },
	}

	ok := func(string, string) Result {
		return Result{Msg: "OK", Err: nil}
	}

	nok := func(string, string, error) Result {
		panic("i should not be called")
	}

	results := ScatterGather(work, ok, nok)
	assert.NotNil(t, results)
	for _, r := range results {
		assert.NoError(t, r.Err)
	}
}

// TestScatterGatherMixedResults test the ScatterGather where some
// goroutine fail and some succeed.
func TestScatterGatherMixedResults(t *testing.T) {

	work := map[string]func() (string, error){
		"1": func() (string, error) { return "1 says OK", nil },
		"2": func() (string, error) { return "", errors.New("2 says not ok") },
		"3": func() (string, error) { return "", errors.New("3 says not ok") },
		"4": func() (string, error) { return "4 says OK", nil },
		"5": func() (string, error) { return "5 says OK", nil },
		"6": func() (string, error) { return "6 says OK", nil },
		"7": func() (string, error) { return "", errors.New("7 says not ok") },
		"8": func() (string, error) { return "", errors.New("8 says not ok") },
		"9": func() (string, error) { return "", errors.New("9 says not ok") },
	}

	ok := func(string, string) Result {
		return Result{Msg: "OK", Err: nil}
	}

	nok := func(key, msg string, err error) Result {
		return Result{Msg: "Propagting", Err: err}
	}

	results := ScatterGather(work, ok, nok)
	assert.NotNil(t, results)
}

// TestScatterGatherAllNotOk test the ScatterGather where all
// goroutine fail.
func TestScatterGatherAllNotOk(t *testing.T) {

	work := map[string]func() (string, error){
		"1": func() (string, error) { return "", errors.New("1 says not ok") },
		"2": func() (string, error) { return "", errors.New("2 says not ok") },
		"3": func() (string, error) { return "", errors.New("3 says not ok") },
		"4": func() (string, error) { return "", errors.New("4 says not ok") },
		"5": func() (string, error) { return "", errors.New("5 says not ok") },
		"6": func() (string, error) { return "", errors.New("6 says not ok") },
		"7": func() (string, error) { return "", errors.New("7 says not ok") },
		"8": func() (string, error) { return "", errors.New("8 says not ok") },
		"9": func() (string, error) { return "", errors.New("9 says not ok") },
	}

	ok := func(string, string) Result {
		return Result{Msg: "OK", Err: nil}
	}

	nok := func(key, msg string, err error) Result {
		return Result{Msg: "Propagting", Err: err}
	}

	results := ScatterGather(work, ok, nok)
	assert.NotNil(t, results)
	for _, r := range results {
		assert.Error(t, r.Err)
	}
}
