package main

import "sync"

// AtomicQueue представляет потокобезопасную очередь для элементов типа T.
type AtomicQueue[T any] struct {
	mu   sync.RWMutex // Мьютекс для синхронизации доступа к очереди.
	list []T          // Срез, содержащий элементы очереди.
}

// Add добавляет элемент value в конец очереди.
//
// @param value элемент типа T, который будет добавлен.
func (queue *AtomicQueue[T]) Add(value T) {
	queue.mu.Lock()
	defer queue.mu.Unlock()
	queue.list = append(queue.list, value)
}

// Get возвращает первый элемент очереди без его удаления.
// Если очередь пуста, возвращает нулевое значение типа T.
//
// @return первый элемент очереди или нулевое значение.
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

// Pop возвращает и удаляет первый элемент очереди.
// Если очередь пуста, возвращает нулевое значение типа T.
//
// @return первый элемент очереди или нулевое значение.
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

// Size возвращает количество элементов в очереди.
//
// @return количество элементов в очереди.
func (queue *AtomicQueue[T]) Size() int {
	queue.mu.RLock()
	defer queue.mu.RUnlock()
	return len(queue.list)
}
