package util

import (
	"sync"
)

type CachedValue struct {
	v	interface{}
	load	func() (interface{}, error)
	lock	sync.RWMutex
}

func MakeCachedValue(load func() (interface{}, error)) *CachedValue {
	return &CachedValue{ load: load }
}

func (cv *CachedValue)Get() (interface{}, error) {
	cv.lock.RLock()
	x := cv.v
	cv.lock.RUnlock()

	if x != nil {
		return x, nil
	}

	cv.lock.Lock()
	defer cv.lock.Unlock()

	x = cv.v
	if x == nil {
		var err error

		x, err = cv.load()
		if err != nil {
			return nil, err
		}

		cv.v = x
	}

	return x, nil
}

func (cv *CachedValue)Invalidate() {
	cv.lock.Lock()
	defer cv.lock.Unlock()
	cv.v = nil
}
