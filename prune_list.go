package util

import (
	"time"
)

type PruneList struct {
	hash	map[string]*PruneItem
	wait	chan interface{}
	ni	chan *PruneItem
}

type PruneItem struct {
	at	time.Time
	key	string
	value	interface{}
}

func NewPruneList() *PruneList {
	pl := &PruneList{}

	pl.hash = make(map[string]*PruneItem)
	pl.wait = make(chan interface{})
	pl.ni = make(chan *PruneItem)

	go pl.loop()

	return pl
}

func (pl *PruneList)Schedule(key string, value interface{}, in time.Duration) {
	pi := &PruneItem{}
	pi.key = key
	pi.value = value
	pi.at = time.Now().Add(in)
	pl.ni <- pi
}

func (pl *PruneList)Unschedule(key string) {
	pi := &PruneItem{}
	pi.key = key
	pi.value = nil
	pl.ni <- pi
}

func (pl *PruneList)min() *PruneItem {
	/* FIXME -- use google/btree or emirpasic/gods */
	var mi *PruneItem

	for _, i := range pl.hash {
		if mi == nil || i.at.Before(mi.at) {
			mi = i
		}
	}

	return mi
}

func (pl *PruneList)loop() {
	var next <-chan time.Time
	var t *time.Timer

	for {
		mi := pl.min()
		if mi == nil {
			next = nil /* will block */
			t = nil
		} else {
			dur := mi.at.Sub(time.Now())
			t = time.NewTimer(dur)
			next = t.C
		}

		select {
		case ni := <-pl.ni:
			if t != nil {
				t.Stop()
			}

			if ni.value == nil {
				delete(pl.hash, ni.key)
			} else {
				ci := pl.hash[ni.key]
				if ci != nil {
					ci.at = ni.at
				} else {
					pl.hash[ni.key] = ni
				}
			}
		case <-next:
			delete(pl.hash, mi.key)
			pl.wait <-mi.value
		}
	}
}

func (pl *PruneList)Wait() <-chan interface{} {
	return pl.wait
}
