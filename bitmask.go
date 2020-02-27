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

type Bitmask uint64

func (f *Bitmask) flag(nr uint) Bitmask {
	if nr > 63 { panic("Too big bitnr") }
	return 1 << nr
}

func (f *Bitmask) Check(nr uint) bool { return (*f) & f.flag(nr) != 0 }
func (f *Bitmask) Set(nr uint)     { *f |= f.flag(nr) }
func (f *Bitmask) Clear(nr uint)   { *f &= ^f.flag(nr) }
func (f *Bitmask) Toggle(nr uint)  { *f ^= f.flag(nr) }

func NewBitmask(nr ...uint) *Bitmask {
	var ret Bitmask
	for _, n := range nr {
		ret.Set(n)
	}
	return &ret
}

func NewFullBitmask() *Bitmask {
	var ret Bitmask
	ret = ^ret
	return &ret
}

func NewEmptyBitmask() *Bitmask {
	var ret Bitmask
	return &ret
}
