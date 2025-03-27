package concurrency

import (
	"log"
	"sync"
	"time"
)

type WorkerPool struct {
	NumOfWorkers int16
	Tasks        chan Executor
	Quit         chan int
	Errors       chan error
	Results      chan bool
	start        sync.Once
	stop         sync.Once
}

func NewWorkerPool(numOfWorkers int16) *WorkerPool {
	return &WorkerPool{
		NumOfWorkers: numOfWorkers,
		Tasks:        make(chan Executor, 256),
		Quit:         make(chan int),
		Errors:       make(chan error),
		Results:      make(chan bool),
	}
}
func (wp *WorkerPool) run() {
	for i := 0; i <= int(wp.NumOfWorkers); i++ {
		go func() {
			for {
				select {
				case task, ok := <-wp.Tasks:
					if !ok {
						log.Fatal("Fatal. WorkerPool channel closed")
						return
					}
					err := task.Execute()
					wp.Results <- true
					if err != nil {
						task.ErrorHandler(err)
					}

				case <-wp.Quit:
					return
				}

			}
		}()
	}
}
func (wp *WorkerPool) Start() {
	wp.start.Do(func() {
		wp.run()
	})
}
func (wp *WorkerPool) WaitForTimeout(d time.Duration) {
	wp.stop.Do(func() {
		timer := time.NewTimer(d)
		for {
			select {
			case <-wp.Results:
				timer.Stop()
				timer = time.NewTimer(d)
			case <-timer.C:
				close(wp.Quit)
				log.Println("WorkerPool Quit")
				return
			}
		}
	},
	)
}

func (wp *WorkerPool) Stop() {
	wp.stop.Do(func() {
		close(wp.Quit)
	})
}

func (wp *WorkerPool) AddTask(task Executor) {
	select {
	case wp.Tasks <- task:
	case <-wp.Quit:
	}
}
