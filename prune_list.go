/////////////////////////////////////////////////////////////////////////////////
//
// Copyright (C) 2019-2020, Unectio Inc, All Right Reserved.
//
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
/////////////////////////////////////////////////////////////////////////////////

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
