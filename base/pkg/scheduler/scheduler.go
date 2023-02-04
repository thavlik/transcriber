package scheduler

import (
	"errors"
)

var ErrLocked = errors.New("entity is locked")

type Lock interface {
	Extend() error
	Release() error
}

type Scheduler interface {
	Add(entities ...string) error
	List() ([]string, error)
	//Rand() (id string, err error)
	Remove(entity string) error
	Lock(entity string) (Lock, error)
	Notify() <-chan struct{}
}
