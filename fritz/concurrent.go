package fritz

import "sync/atomic"

type result struct {
	msg string
	err error
}

func scatterGather(workTable map[string]func() (string, error), onSuccess func(string, string) result, onError func(string, string, error) result) []result {
	amountOfWork := len(workTable)

	scatterChannel := make(chan result, amountOfWork)
	gatherChannel := make(chan result, amountOfWork)
	for key, work := range workTable {
		go func(k string, w func() (string, error)) {
			msg, err := w()
			if err == nil {
				scatterChannel <- onSuccess(k, msg)
			} else {
				scatterChannel <- onError(k, msg, err)
			}
		}(key, work)
	}

	var ops uint64
	go func() {
		for {
			res := <-scatterChannel
			atomic.AddUint64(&ops, 1)
			gatherChannel <- res
			if atomic.LoadUint64(&ops) == uint64(amountOfWork) {
				close(scatterChannel)
				close(gatherChannel)
				return
			}
		}
	}()

	results := make([]result, 0, amountOfWork)
	for res := range gatherChannel {
		results = append(results, res)
	}
	return results
}
