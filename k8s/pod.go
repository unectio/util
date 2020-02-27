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

package k8s

import (
	"fmt"
	"github.com/unectio/api"
	"k8s.io/api/core/v1"
)

type Pod struct {
	UID	string
	Addr	string
	Port	string
	Host	string
	Helper	string
	Token	string

	DepDesc
}

func (p *Pod)URL(proto string) string {
	if proto != "" {
		proto += "://"
	}

	return proto + p.Addr + ":" + p.Port
}

const (
	LabelType	= "langlet-type"
	LabelLang	= "langlet-lang"

	LabelFunc	= "function"
	LabelHelper	= "helper"
)

func (p *Pod)scanEnv(env []v1.EnvVar) {
	for _, v := range env {
		switch v.Name {
		case api.EnvPort:
			p.Port = v.Value
		case api.EnvProject:
			p.proj = v.Value
		case api.EnvLang:
			p.lang = v.Value
		case api.EnvClass:
			p.class = v.Value
		case api.EnvToken:
			p.Token = v.Value
		}
	}
}

func toPod(pod *v1.Pod) (*Pod, error) {
	p := &Pod {
		UID:		string(pod.ObjectMeta.UID),
		Addr:		pod.Status.PodIP,
		Host:		pod.Status.HostIP,
	}

	typ := pod.ObjectMeta.Labels[LabelType]
	switch typ {
	case LabelHelper:
		p.Helper = pod.ObjectMeta.Labels[LabelLang]
	case LabelFunc:
		; /* that's OK */
	default:
		return nil, fmt.Errorf("Non-langlet POD %s popped up", p.UID)
	}

	for _, c := range pod.Spec.Containers {
		p.scanEnv(c.Env)
	}

	if p.Helper != "" && p.Helper != p.lang {
		return nil, fmt.Errorf("Lang mismatch: %s != %s", p.Helper, p.Lang)
	}

	return p, nil
}
