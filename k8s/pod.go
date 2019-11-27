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
