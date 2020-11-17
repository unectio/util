package router

import (
	"fmt"
)

func debug_print(f string, args ...interface{}) {
	fmt.Printf("ROUTER:"+f, args...)
}

var debug func(f string, args ...interface{})

func Debug() {
	debug = debug_print
}

func (r *Router) Print() {
	fmt.Print("======\n")
	r.root.print("\t")
}

func (l *layer) print(pfx string) {
	if l.match != nil {
		fmt.Printf("%s=->%v\n", pfx, l.match)
	}
	for mn, n := range l.exact {
		fmt.Printf("%s%s->\n", pfx, mn)
		n.print(pfx + "\t")
	}
	for pn, n := range l.param {
		fmt.Printf("%s{%s}->\n", pfx, pn)
		n.print(pfx + "\t")
	}
}
