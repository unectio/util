package util

import (
	"os"
	"fmt"
	"regexp"
)

type Creds struct {
	Usr	string
	Pwd	string
	Adr	string
	Prt	string
	Dom	string
}

func (c *Creds)String() string {
	return fmt.Sprintf("%s:%s@%s:%s/%s", c.Usr, c.Pwd, c.Adr, c.Prt, c.Dom)
}

func CredsParse(url string) *Creds {
	/* user:pass@host:port[/domain] */
	re := regexp.MustCompile("^([a-zA-Z0-9_@.-]+):(.*)@([a-zA-Z0-9._-]+):([0-9]+)(/([a-zA-Z0-9_]+))?$")
	m := re.FindAllStringSubmatch(url, -1)
	if len(m) == 0 {
		return nil
	}
	if len(m[0]) != 7 {
		return nil
	}

	return &Creds {
		Usr:	checkEnv(m[0][1]),
		Pwd:	checkEnv(m[0][2]),
		Adr:	checkEnv(m[0][3]),
		Prt:	m[0][4],
		Dom:	checkEnv(m[0][6]),
	}
}

func checkEnv(v string) string {
	if v == "-" {
		return ""
	}

	for _, l := range v {
		if ! ((l >= 'A' && l <= 'Z') || l == '_') {
			return v
		}
	}

	/* It's all caps and _-s */
	x := os.Getenv(v)
	if x == "" {
		return v
	}

	return x
}
