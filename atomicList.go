package main

import "sync"

type AtomicQueue[T any] struct {
	mu   sync.RWMutex
	list []T
}

func (queue *AtomicQueue[T]) Add(value T) {
	queue.mu.Lock()
	defer queue.mu.Unlock()
	queue.list = append(queue.list, value)
}

func (queue *AtomicQueue[T]) Get() T {
	queue.mu.RLock()
	defer queue.mu.RUnlock()

	if len(queue.list) == 0 {
		var zeroValue T
		return zeroValue
	}

	value := queue.list[0]
	return value
}

func (queue *AtomicQueue[T]) Pop() T {
	queue.mu.Lock()
	defer queue.mu.Unlock()

	if len(queue.list) == 0 {
		var zeroValue T
		return zeroValue
	}

	value := queue.list[0]
	queue.list = queue.list[1:]
	return value
}

func (queue *AtomicQueue[T]) Size() int {
	queue.mu.RLock()
	defer queue.mu.RUnlock()
	return len(queue.list)

}