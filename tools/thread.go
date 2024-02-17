package tools

import (
	"sync"
	"time"
)

type Thread struct {
	max   int
	count int
	m     sync.Mutex
}

func NewThread() *Thread {
	return &Thread{
		max:   3,
		count: 0,
		m:     sync.Mutex{},
	}
}

func (t *Thread) SetMax(m int) {
	t.max = m
}

func (t *Thread) Count() int {
	return t.count
}

func (t *Thread) Max() int {
	return t.max
}

func (t *Thread) Add() {
	t.m.Lock()
	defer t.m.Unlock()
	t.count += 1
}

func (t *Thread) Done() {
	t.m.Lock()
	defer t.m.Unlock()
	t.count -= 1
}

func (t *Thread) Wait() {
	for t.count > t.max {
		time.Sleep(time.Millisecond * 100)
	}
}

func (t *Thread) WaitAll() {
	for t.count > 0 {
		time.Sleep(time.Millisecond * 100)
	}
}
