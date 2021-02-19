package syncx

import (
	"errors"
	"sync"
)

var ErrCounterEmpty = errors.New("counter empty. someone returned multiple times")

type Counter struct {
	Max  int32
	Curr int32
	mu   sync.Mutex
}

func (ct *Counter) TryBorrow() bool {
	defer ct.mu.Unlock()

	ct.mu.Lock()
	if ct.Curr+1 > ct.Max {
		return false
	} else {
		ct.Curr++
		return true
	}
}

func (ct *Counter) Return() error {
	defer ct.mu.Unlock()

	ct.mu.Lock()
	ct.Curr--
	if ct.Curr < 0 {
		ct.Curr = 0
		return ErrCounterEmpty
	}
	return nil
}
