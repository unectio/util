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
