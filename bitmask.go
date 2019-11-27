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
