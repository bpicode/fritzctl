package fritz

import (
	"sync"
)

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
	ch := fanOut(wt, onSuccess, onError)
	res := fanIn(ch)
	results := <-res
	return results
}

func fanOut(wt workTable, onSuccess successHandler, onError errorHandler) <-chan result {
	wg := new(sync.WaitGroup)
	wg.Add(len(wt))
	ch := scatter(wt, onSuccess, onError, wg)
	go func() {
		wg.Wait()
		close(ch)
	}()
	return ch
}

func scatter(wt workTable, onSuccess successHandler, onError errorHandler, wg *sync.WaitGroup) chan result {
	ch := make(chan result)
	for key, work := range wt {
		go func(k string, w func() (string, error)) {
			defer wg.Done()
			msg, err := w()
			if err == nil {
				ch <- onSuccess(k, msg)
			} else {
				ch <- onError(k, msg, err)
			}
		}(key, work)
	}
	return ch
}

func fanIn(ch <-chan result) <-chan []result {
	collect := make(chan []result)
	go func() {
		results := make([]result, 0)
		for res := range ch {
			results = append(results, res)
		}
		collect <- results
		close(collect)
	}()
	return collect
}
