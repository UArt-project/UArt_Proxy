// Package workerpool provides a workerpool for concurrent task handling.
package workerpool

// WorkerPool provides functionality of running multiple tasks in separate routines.
type WorkerPool struct {
	workerNumber int
	taskChan     chan func()
	stopChan     chan struct{}
}

// NewPool returns new instance of a workerpool.
func NewPool(workerNumber int) *WorkerPool {
	if workerNumber < 1 {
		workerNumber = 1
	}

	taskChan := make(chan func(), workerNumber)
	stopChan := make(chan struct{})

	return &WorkerPool{
		workerNumber: workerNumber,
		taskChan:     taskChan,
		stopChan:     stopChan,
	}
}

// Start runs workerpool and makes it available for task handling.
func (p *WorkerPool) Start() {
	for i := 0; i < p.workerNumber; i++ {
		go func() {
			for {
				select {
				case <-p.stopChan:
					break
				case task := <-p.taskChan:
					task()
				}
			}
		}()
	}
}

// AddTask throws new tasks into the queue for the workers to complete.
func (p *WorkerPool) AddTask(task func()) {
	p.taskChan <- task
}

// Stop shuts workers down.
func (p *WorkerPool) Stop() {
	close(p.stopChan)
}
