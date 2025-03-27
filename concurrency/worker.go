package concurrency

type Executor struct {
	Execute      func() error
	ErrorHandler func(error)
}

type Worker interface {
	Start()
	AddTask(task Executor)
}
