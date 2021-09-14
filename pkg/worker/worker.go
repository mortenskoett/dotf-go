/* The worker package contains an implementation of a worker that can be used
to run and manage background goroutines that should run with intervals. */
package worker

import (
	"log"
	"time"
)

type Worker struct {
	Interval time.Duration // Interval between Action is executed.
	Action   func()        // Action is a simple function that is called at every interval.
	shutdown chan string   // Channel used to indicate when to shutdown the worker.
}

func NewWorker() *Worker {
	return &Worker{
		shutdown: make(chan string),
	}
}

func NewWorkerParam(interval time.Duration, action func()) *Worker {
	return &Worker{
		Interval: interval,
		Action:   action,
		shutdown: make(chan string),
	}
}

func (w *Worker) Start() {
	go runWorker(w)
}

func runWorker(w *Worker) {
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

func (w *Worker) Stop() {
	log.Println("Worker stopping")
	w.shutdown <- "Shutdown"
}
