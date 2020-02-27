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
	"sync"
	"errors"
	"k8s.io/api/extensions/v1beta1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	inmemPodAddr	= "in.mem.pod.addr"
)

type ievent struct {
	pod	*Pod
	up	bool
}

type KubeOSet struct {
	deps	sync.Map
	evs	chan *ievent
	notify	func(*Pod, bool)
}

func (k *KubeOSet)Deps() DepAPI { return k }

func (k *KubeOSet)Notify(h *EventHandlers, data interface{}) {
	k.notify = func(pod *Pod, up bool) {
		h.PodLifecycle(pod, up, data)
	}
}

func KubeInMem() Client {
	ret := &KubeOSet{}
	ret.evs = make(chan *ievent)
	go func() {
		for ev := range ret.evs {
			if ret.notify != nil {
				ret.notify(ev.pod, ev.up)
			}
		}
	}()
	return ret
}

func (k *KubeOSet)Create(spec *v1beta1.Deployment) (*v1beta1.Deployment, error) {
	n := spec.ObjectMeta.Name
	k.deps.Store(n, spec)
	go func() {
		p := &Pod { UID: "0" }
		p.Addr = inmemPodAddr
		p.scanEnv(spec.Spec.Template.Spec.Containers[0].Env)
		k.evs <- &ievent{ pod: p, up: true }
	}()
	return spec, nil
}

func (k *KubeOSet)Delete(name string, ops *meta.DeleteOptions) error {
	k.deps.Delete(name)
	return nil
}

func (k *KubeOSet)Get(name string, ops meta.GetOptions) (*v1beta1.Deployment, error) {
	x, ok := k.deps.Load(name)
	if !ok {
		return nil, errors.New("not found")
	}

	return x.(*v1beta1.Deployment), nil
}

func (k *KubeOSet)Update(spec *v1beta1.Deployment) (*v1beta1.Deployment, error) {
	n := spec.ObjectMeta.Name
	k.deps.Store(n, spec)
	return spec, nil
}
