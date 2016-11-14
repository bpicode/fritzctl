package concurrent

import "sync/atomic"

// Result is a simple model of a concurrent task, having a
// simple payload and an error.
type Result struct {
	Msg string
	Err error
}

// ScatterGather forks the workTable into separate goroutines
// with callbacks onSuccess and onError. The results are gathered
// in slice.
func ScatterGather(workTable map[string]func() (string, error), onSuccess func(string, string) Result, onError func(string, string, error) Result) []Result {
	amountOfWork := len(workTable)
	if amountOfWork == 0 {
		return []Result{}
	}
	scatterChannel := make(chan Result, amountOfWork)
	var ops uint64
	for key, work := range workTable {
		go func(k string, w func() (string, error)) {
			msg, err := w()
			if err == nil {
				scatterChannel <- onSuccess(k, msg)
			} else {
				scatterChannel <- onError(k, msg, err)
			}
			if atomic.AddUint64(&ops, 1) == uint64(amountOfWork) {
				close(scatterChannel)
			}
		}(key, work)
	}

	results := make([]Result, 0, amountOfWork)
	for res := range scatterChannel {
		results = append(results, res)
	}
	return results
}
