package main

import (
	"fmt"
	"runtime/debug"
	"sync"
	"time"
)

// Worker Pool Stucture
type WorkerPool struct {
	Count  int
	Queue  chan func()
	Waiter *sync.WaitGroup
}

var instance *WorkerPool
var once sync.Once

// NewWorkerPool - creates a new worker pool 
// count - number of workers in the pool.
func NewWorkerPool(count int) *WorkerPool {
	once.Do(func() {
		instance = &WorkerPool{
			Count:  count,
			Queue:  make(chan func(), 1024),
			Waiter: &sync.WaitGroup{},
		}
		instance.start()
	})
	return instance
}

func (pool *WorkerPool) exec(function func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error: recovered from ", string(debug.Stack()))
		}
	}()
	function()
}

func (pool *WorkerPool) worker() {
	for {
		select {
		case function, ok := <-pool.Queue:
			if !ok {
				return
			}
			pool.exec(function)
			pool.Waiter.Done()
		default:
			time.Sleep(time.Millisecond * 1)
		}
	}
}

func (pool *WorkerPool) start() {
	for i := 0; i < pool.Count; i++ {
		go pool.worker()
	}
}

// Enqueue task to worker pool
// Any function can be enqueued
func (pool *WorkerPool) Enqueue(function func()) {
	if function == nil {
		fmt.Println("Error: cannot enqueue 'nil' function for execution")
	}
	pool.Waiter.Add(1)
	pool.Queue <- function
}

func (pool *WorkerPool) Await() {
	pool.Waiter.Wait()
}
