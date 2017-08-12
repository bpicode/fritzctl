package concurrent

import "sync/atomic"

// Result is a simple model of a concurrent task, having a
// simple payload and an error.
type Result struct {
	Msg string
	Err error
}

type successHandler func(string, string) Result

type errorHandler func(string, string, error) Result

type workTable map[string]func() (string, error)

// ScatterGather forks the workTable into separate goroutines
// with callbacks onSuccess and onError. The results are gathered
// in slice. Neither onSuccess nor onError should panic, otherwise
// ScatterGather panics.
func ScatterGather(wt workTable, onSuccess successHandler, onError errorHandler) []Result {
	amountOfWork := len(wt)
	if amountOfWork == 0 {
		return []Result{}
	}
	scatterChannel := make(chan Result, amountOfWork)
	var ops uint64
	for key, work := range wt {
		go func(k string, w func() (string, error)) {
			defer closeOnDone(&ops, uint64(amountOfWork), scatterChannel)
			msg, err := w()
			if err == nil {
				scatterChannel <- onSuccess(k, msg)
			} else {
				scatterChannel <- onError(k, msg, err)
			}
		}(key, work)
	}

	results := make([]Result, 0, amountOfWork)
	for res := range scatterChannel {
		results = append(results, res)
	}
	return results
}

func closeOnDone(ops *uint64, amountOfWork uint64, ch chan Result) {
	if atomic.AddUint64(ops, 1) == amountOfWork {
		close(ch)
	}
}
