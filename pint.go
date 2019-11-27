package util

/* Use with lock */

type PInt struct {
	Val	uint
	Peak	uint
}

func (p *PInt)Inc() {
	p.Val++
	if p.Val > p.Peak {
		p.Peak = p.Val
	}
}

func (p *PInt)Dec() {
	p.Val--
}
