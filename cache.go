package main

import (
	"errors"
	"sync"
)

var NotFoundError = errors.New("run not found")

type RunCache struct {
	store map[string]*Run
	sync.RWMutex
}

func NewCache() *RunCache {
	r := new(RunCache)
	r.store = make(map[string]*Run)
	return r
}

func (r *RunCache) Set(key string, value *Run) {
	r.Lock()
	defer r.Unlock()
	r.store[key] = value
}

func (r *RunCache) Get(key string) (*Run, error) {
	r.RLock()
	defer r.RUnlock()
	value, ok := r.store[key]
	if !ok {
		return nil, NotFoundError
	}
	return value, nil
}
