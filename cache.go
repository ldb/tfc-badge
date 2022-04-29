package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	r.set(key, value)
}

func (r *RunCache) set(key string, value *Run) {
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

func (r *RunCache) List() []*Run {
	r.RLock()
	defer r.RUnlock()
	rr := make([]*Run, 0, len(r.store))
	for _, r := range r.store {
		rr = append(rr, r)
	}
	return rr
}

func (r *RunCache) Dump(w io.Writer) error {
	r.RLock()
	defer r.RUnlock()
	e := json.NewEncoder(w)
	for id, run := range r.store {
		if err := e.Encode(run); err != nil {
			return fmt.Errorf("error encoding run %s: %v", id, err)
		}
	}
	return nil
}

func (r *RunCache) Restore(reader io.Reader) error {
	r.Lock()
	defer r.Unlock()
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		run := new(Run)
		if err := json.Unmarshal(scanner.Bytes(), run); err != nil {
			return fmt.Errorf("error decoding data: %v", err)
		}
		r.set(run.WorkspaceID, run)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error scanning data: %v", err)
	}
	return nil
}
