package domain

import "sync"

type Res struct {
	Errs map[string]string
	Err  error
}

func CloseChannel(wg *sync.WaitGroup, ch chan Res) {
	wg.Wait()
	close(ch)
}
