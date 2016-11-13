package fritz

import "sync/atomic"

type result struct {
	msg string
	err error
}

func scatterGather(workTable map[string]func() (string, error), onSuccess func(string, string) result, onError func(string, string, error) result) []result {
	amountOfWork := len(workTable)
	scatterChannel := make(chan result, amountOfWork)
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

	results := make([]result, 0, amountOfWork)
	for res := range scatterChannel {
		results = append(results, res)
	}
	return results
}
