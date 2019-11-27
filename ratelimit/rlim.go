package ratelimit

import (
	"sync"
	"time"
)

type Filter struct {
	rate	uint
	burst	uint
	base	time.Duration

	stock	uint
	ts	time.Time

	lock	sync.Mutex
}

func (f *Filter)gain(d time.Duration) uint {
	/* One stock point is acquired for time.Second / rate time */
	gain := uint(uint64(d) * uint64(f.rate) / uint64(f.base))
	if gain > f.burst {
		gain = f.burst
	}
	return gain
}

func (f *Filter)FullLocked() bool {
	return f.stock >= f.burst || time.Now().Sub(f.ts) >= f.base
}

func (f *Filter)Full() bool {
	f.lock.Lock()
	defer f.lock.Unlock()
	return f.FullLocked()
}

func (f *Filter)refill() bool {
	now := time.Now()
	d := now.Sub(f.ts)
	if d >= f.base {
		f.stock = f.burst
		f.ts = now
		return true
	}

	gain := f.gain(d)
	if gain == 0 {
		return false
	}

	f.stock = gain
	f.ts = f.ts.Add(f.base * time.Duration(gain) / time.Duration(f.rate))
	return true
}

func (f *Filter)StepLocked() bool {
	if f.rate == 0 {
		return true
	}

	if f.stock == 0 && !f.refill() {
		return false
	}

	f.stock--
	return true
}

func (f *Filter)Step() bool {
	f.lock.Lock()
	ok := f.StepLocked()
	f.lock.Unlock()
	return ok
}

func (f *Filter)UndoLocked() {
	if f.rate != 0 && f.stock < f.burst {
		f.stock++
	}
}

func (f *Filter)Undo() {
	f.lock.Lock()
	f.UndoLocked()
	f.lock.Unlock()
}

func (f *Filter)ResetLocked(rate, burst uint, base time.Duration) {
	if f.rate != rate || f.burst != burst + 1 {
		f.rate = rate
		f.burst = burst + 1
		f.base = base

		f.stock = burst + 1
		f.ts = time.Now()
	}
}

func (f *Filter)Reset(rate, burst uint) {
	f.lock.Lock()
	f.ResetLocked(rate, burst, f.base)
	f.lock.Unlock()
}

func NewFilterWithBase(rate, burst uint, base time.Duration) *Filter {
	f := &Filter{}
	f.ResetLocked(rate, burst, base)
	return f
}

func NewFilter(rate, burst uint) *Filter {
	return NewFilterWithBase(rate, burst, time.Second)
}
