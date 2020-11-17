//////////////////////////////////////////////////////////////////////////////
//
// (C) Copyright 2019-2020 by Unectio, Inc.
//
// The information contained herein is confidential, proprietary to Unectio,
// Inc.
//
//////////////////////////////////////////////////////////////////////////////

package workq

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func run(wq *Pool, key, seq string, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		wq.Run(key, func(w interface{}) {
			time.Sleep(500 * time.Millisecond)
			fmt.Printf("%s.%s.%s\n", key, w.(string), seq)
			wg.Done()
		})
	}()
}

func TestWqueue(t *testing.T) {
	var wg sync.WaitGroup

	wq := Make()
	wq.AddWorker("1")
	wq.AddWorker("2")

	run(wq, "a", "1", &wg)
	run(wq, "a", "2", &wg)
	run(wq, "a", "3", &wg)
	run(wq, "a", "4", &wg)
	run(wq, "b", "1", &wg)
	run(wq, "b", "2", &wg)
	run(wq, "c", "1", &wg)
	run(wq, "d", "1", &wg)
	run(wq, "d", "2", &wg)
	run(wq, "d", "3", &wg)
	run(wq, "d", "4", &wg)
	run(wq, "d", "5", &wg)
	run(wq, "d", "6", &wg)

	wg.Wait()
}
