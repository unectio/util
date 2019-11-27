package workq

import (
	"sync"
	"container/list"
)

type Pool struct {
	in	chan *request
	queues	*list.List
	lock	sync.Mutex
	wake	*sync.Cond
}

type Worker struct {
	stop	bool
	p	*Pool
}

type queue struct {
	key	string
	reqs	[]*request
	l	*list.Element
}

type request struct {
	key	string
	fn	func(interface{})
	done	chan bool
}

func Make() *Pool {
	ret := &Pool{}

	ret.in = make(chan *request)
	ret.queues = list.New()
	ret.wake = sync.NewCond(&ret.lock)

	return ret
}

func (pool *Pool)queue(key string) *queue {
	for e := pool.queues.Front(); e != nil; e = e.Next() {
		q := e.Value.(*queue)
		if q.key == key {
			return q
		}
	}

	return nil
}

func (pool *Pool)enqueue(rq *request) {
	pool.lock.Lock()
	defer func() {
		pool.wake.Signal()
		pool.lock.Unlock()
	}()

	q := pool.queue(rq.key)
	if q != nil {
		q.reqs = append(q.reqs, rq)
	} else {
		q = &queue{key: rq.key, reqs: []*request{rq}}
		q.l = pool.queues.PushBack(q)
	}
}

func (w *Worker)next() *request {
	var x *list.Element

	pool := w.p
	pool.lock.Lock()
	for {
		if w.stop {
			pool.lock.Unlock()
			return nil
		}
		x = pool.queues.Front()
		if x != nil {
			break
		}

		pool.wake.Wait()
	}

	q := x.Value.(*queue)
	rq := q.reqs[0]
	q.reqs = q.reqs[1:]
	if len(q.reqs) > 0 {
		pool.queues.MoveToBack(x)
	} else {
		pool.queues.Remove(x)
	}
	pool.lock.Unlock()

	return rq
}

func (w *Worker)work(wi interface{}) {
	for {
		rq := w.next()
		if rq == nil {
			break
		}
		rq.fn(wi)
		close(rq.done)
	}
}

func (w *Worker)Stop() {
	w.stop = true
	w.p.wake.Broadcast()
}

func (pool *Pool)AddWorker(wi interface{}) *Worker {
	w := &Worker{p: pool}
	go w.work(wi)
	return w
}

func (pool *Pool)Run(key string, f func(w interface{})) {
	rq := &request{ key: key, fn: f, done: make(chan bool) }
	pool.enqueue(rq)
	<-rq.done
}
