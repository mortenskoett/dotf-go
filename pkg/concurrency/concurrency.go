// The concurrency package contains concurrency related implementations
package concurrency

import (
	"log"
	"time"
)

// IntervalWorker is an implementation of a worker that can be used to run and manage background
// goroutines that should run with intervals.
type IntervalWorker struct {
	Interval time.Duration // Interval between Action is executed.
	Action   func()        // Action is a simple function that is called at every interval.
	shutdown chan string   // Channel used to indicate when to shutdown the worker.
}

func NewIntervalWorker() *IntervalWorker {
	return &IntervalWorker{
		shutdown: make(chan string),
	}
}

func NewIntervalWorkerParam(interval time.Duration, action func()) *IntervalWorker {
	return &IntervalWorker{
		Interval: interval,
		Action:   action,
		shutdown: make(chan string),
	}
}

func (w *IntervalWorker) Start() {
	go runWorker(w)
	log.Println("Worker started")
}

func runWorker(w *IntervalWorker) {
	ticker := time.NewTicker(w.Interval)

	for {
		select {
		case <-w.shutdown:
			ticker.Stop()
			close(w.shutdown)
			return
		case <-ticker.C:
			w.Action()
		}
	}
}

func (w *IntervalWorker) Stop() {
	log.Println("Worker stopping")
	w.shutdown <- "Shutdown"
}
