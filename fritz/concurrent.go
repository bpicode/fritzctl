package fritz

import "sync/atomic"

// result is a simple model of a concurrent task, having a
// simple payload and an error.
type result struct {
	msg string
	err error
}

type successHandler func(string, string) result

type errorHandler func(string, string, error) result

type workTable map[string]func() (string, error)

// scatterGather forks the workTable into separate goroutines
// with callbacks onSuccess and onError. The results are gathered
// in slice. Neither onSuccess nor onError should panic, otherwise
// scatterGather panics.
func scatterGather(wt workTable, onSuccess successHandler, onError errorHandler) []result {
	amountOfWork := len(wt)
	if amountOfWork == 0 {
		return []result{}
	}
	scatterChannel := make(chan result, amountOfWork)
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

	results := make([]result, 0, amountOfWork)
	for res := range scatterChannel {
		results = append(results, res)
	}
	return results
}

func closeOnDone(ops *uint64, amountOfWork uint64, ch chan result) {
	if atomic.AddUint64(ops, 1) == amountOfWork {
		close(ch)
	}
}
