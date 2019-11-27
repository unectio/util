package context

type global struct {
	desc		string
}

func (g *global)Scope() string { return "GLOBAL:" + g.desc }

func Global(desc string) LoginInfo {
	return &global{
		desc:	desc,
	}
}
