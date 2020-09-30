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
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
)

type event struct {
	up  bool
	pod *v1.Pod
}

type EventHandlers struct {
	PodLifecycle    func(*Pod, bool, interface{})
	PodError        func(error, interface{})
	PodUnshedulable func(string, interface{})
}

func checkSchedulable(pod *v1.Pod) bool {
	for _, cond := range pod.Status.Conditions {
		if cond.Type == v1.PodScheduled && cond.Status == v1.ConditionFalse && cond.Reason == v1.PodReasonUnschedulable {
			return true
		}
	}
	return false
}

func (kc *KubeNsClient) Notify(handlers *EventHandlers, data interface{}) {
	events := make(chan *event)

	go func() {
		for {
			e := <-events
			pod, err := toPod(e.pod)
			if err == nil {
				handlers.PodLifecycle(pod, e.up, data)
			} else if handlers.PodError != nil {
				handlers.PodError(err, data)
			}
		}
	}()

	watchlist := cache.NewListWatchFromClient(kc.c.CoreV1().RESTClient(), "pods", kc.ns, fields.Everything())
	_, controller := cache.NewInformer(watchlist, &v1.Pod{}, time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				pod := obj.(*v1.Pod)
				if pod.Status.PodIP != "" {
					events <- &event{true, pod}
				}
			},

			DeleteFunc: func(obj interface{}) {
				pod := obj.(*v1.Pod)
				if pod.Status.PodIP != "" {
					events <- &event{false, pod}
				}
			},

			UpdateFunc: func(oobj, nobj interface{}) {
				opod := oobj.(*v1.Pod)
				npod := nobj.(*v1.Pod)

				if handlers.PodUnshedulable != nil {
					if checkSchedulable(npod) {
						uid := string(npod.ObjectMeta.UID)
						handlers.PodUnshedulable(uid, data)
					}
				}

				if opod.Status.PodIP == "" && npod.Status.PodIP != "" {
					events <- &event{true, npod}
				} else if opod.Status.PodIP != "" && npod.Status.PodIP == "" {
					events <- &event{false, npod}
				}
			},
		})

	go controller.Run(make(chan struct{}))
}
