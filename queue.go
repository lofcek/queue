package queue

import (
	"errors"
	"sync"
)

// Queue represents bocking queue of string
type Queue struct {
	data       []string
	isFinished bool
	mtx        sync.Mutex
	cv         *sync.Cond
}

// New creates a queue
func New() *Queue {
	q := Queue{}
	q.cv = sync.NewCond(&q.mtx)
	return &q
}

// Push appends new string at the end the queue.
func (q *Queue) Push(s string) {
	q.mtx.Lock()
	defer q.mtx.Unlock()

	if q.isFinished {
		return
	}

	if len(q.data) == 0 {
		q.cv.Signal()
	}

	q.data = append(q.data, s)
}

// Pop returns first string from queue or block.
//
// After call Finish it returns immediately error "ErrDone"
func (q *Queue) Pop() (string, error) {
	q.mtx.Lock()
	defer q.mtx.Unlock()

	for !q.isFinished && len(q.data) == 0 {
		q.cv.Wait()
	}

	if q.isFinished {
		return "", ErrDone
	}

	ret := q.data[0]
	q.data = q.data[1:]
	return ret, nil
}

// Finish ensures that next Pop return error "ErrDone" without blocking.
func (q *Queue) Finish() {
	q.mtx.Lock()
	defer q.mtx.Unlock()

	q.isFinished = true
	q.cv.Broadcast()
}

// ErrDone is a error returned after every Pop after Finish
var ErrDone = errors.New("Queue is finished")
