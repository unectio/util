package k8s

import (
	"time"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/apimachinery/pkg/fields"
)

type event struct {
	up	bool
	pod	*v1.Pod
}

type EventHandlers struct {
	PodLifecycle func(*Pod, bool, interface{})
	PodError func(error, interface{})
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

func (kc *KubeNsClient)Notify(handlers *EventHandlers, data interface{}) {
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

	watchlist := cache.NewListWatchFromClient(kc.c.Core().RESTClient(), "pods", kc.ns, fields.Everything())
	_, controller := cache.NewInformer(watchlist, &v1.Pod{}, time.Second * 0,
			cache.ResourceEventHandlerFuncs{
				AddFunc:	func (obj interface{}) {
					pod := obj.(*v1.Pod)
					if pod.Status.PodIP != "" {
						events <-&event{true, pod}
					}
				},

				DeleteFunc:	func(obj interface{}) {
					pod := obj.(*v1.Pod)
					if pod.Status.PodIP != "" {
						events <-&event{false, pod}
					}
				},

				UpdateFunc:	func(oobj, nobj interface{}) {
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
