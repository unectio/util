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
	"sync"
	"time"
	"strconv"
)

const (
	VcDying	= "-1"
)

type VCount struct {
	lock	sync.RWMutex
	counts	map[string]*count
	wake	*sync.Cond
}

type count struct {
	vs	map[string]uint
}

func (c *count)dead() bool {
	return len(c.vs) == 0
}

func newCount() *count {
	c := &count{}
	c.vs = make(map[string]uint)
	return c
}

func MakeVCount() *VCount {
	vc := &VCount{}
	vc.wake = sync.NewCond(vc.lock.RLocker())
	vc.counts = make(map[string]*count)
	return vc
}

func (vc *VCount)Add(id, ver string) {
	vc.lock.Lock()
	defer vc.lock.Unlock()

	cm, ok := vc.counts[id]
	if !ok {
		cm = newCount()
		vc.counts[id] = cm
	}
	cm.vs[ver] = cm.vs[ver] + 1
	vc.wake.Broadcast()
}

func (vc *VCount)Del(id, ver string) {
	vc.lock.Lock()
	defer vc.lock.Unlock()

	cm, ok := vc.counts[id]
	if !ok {
		/* pod that fails to start still emits stop event */
		return
	}

	nv := cm.vs[ver] - 1
	if nv == 0 {
		delete(cm.vs, ver)
		if cm.dead() {
			delete(vc.counts, id)
		}
	} else {
		cm.vs[ver] = nv
	}

	vc.wake.Broadcast()
}

func (vc *VCount)List(id string) []string {
	var ret []string

	vc.lock.RLock()
	defer vc.lock.RUnlock()

	cm, ok := vc.counts[id]
	if ok {
		for v, _ := range cm.vs {
			ret = append(ret, v)
		}
	}

	return ret
}

func equal(vers map[string]uint, ver string) bool {
	_, ok := vers[ver]
	return ok
}

func over(vers map[string]uint, ver string) bool {
	if equal(vers, ver) {
		return true
	}

	veri, _ := strconv.Atoi(ver)
	for v, _ := range vers {
		vi, _ := strconv.Atoi(v)
		if vi > veri {
			return true
		}
	}

	return false
}

/*
 * Returns two bools -- timeout or the id has gone empty. If the timeout
 * is true, the second bool is undefined, so the former must be checked
 * first.
 */
func (vc *VCount)Wait(id string, tmo time.Duration, vermatch string) (bool, bool) {
	timeout := false
	timer := time.AfterFunc(tmo, func() {
			timeout = true
			vc.wake.Broadcast()
		})

	vc.lock.RLock()
	var gone bool
	for {
		if timeout {
			break
		}

		vers, has := vc.counts[id]
		if has {
			if _, ok := vers.vs[VcDying]; ok {
				gone = true
				break
			}

			if vermatch[0] == '=' {
				if equal(vers.vs, vermatch[1:]) {
					break
				}
			} else if vermatch[0] == '+' {
				if over(vers.vs, vermatch[1:]) {
					break
				}
			}
		}

		vc.wake.Wait()
	}
	vc.lock.RUnlock()

	timer.Stop()

	return timeout, gone
}
