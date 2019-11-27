package sync

import (
	"sync"
	"time"
)

func CondWaitTmo(cond *sync.Cond, timeout time.Duration) {
	timer := time.AfterFunc(timeout, func() { cond.Signal() })
	defer timer.Stop()
	cond.Wait()
}

type Till time.Duration

const (
	TillMinute Till	= Till(time.Minute)
	TillHour Till	= Till(time.Hour)
	TillDay Till	= Till(24 * time.Hour)
)

func (t Till)Sleep() {
	var left time.Duration

	now := time.Now()
	left = now.Truncate(time.Duration(t)).Add(time.Duration(t)).Sub(now)
	time.Sleep(left)
}

type WaitGroupErr struct {
	sync.WaitGroup
	err error
}

func (wg *WaitGroupErr)Inc() {
	wg.WaitGroup.Add(1)
}

func (wg *WaitGroupErr)Wait() error {
	wg.WaitGroup.Wait()
	return wg.err
}

func (wg *WaitGroupErr)Done(err error) {
	wg.err = err
	wg.WaitGroup.Done()
}
