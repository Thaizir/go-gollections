package set

import (
	"iter"
	"sync"
)

// ConcurrentSet is a thread-safe wrapper around a set of comparable elements.
// It is safe to share a *ConcurrentSet[T] across goroutines: every method
// takes the internal RWMutex before touching the underlying map.
//
// Unlike Set[T], this type deliberately opts into the cost of locking on
// every call. If you don't need concurrent access, use Set[T] instead.
//
// A ConcurrentSet must not be copied after first use (like sync.Mutex).
type ConcurrentSet[T comparable] struct {
	mu  sync.RWMutex
	set Set[T]
}

// --- Constructors ---

func NewConcurrent[T comparable]() *ConcurrentSet[T] {
	return &ConcurrentSet[T]{set: Set[T]{elements: make(map[T]struct{})}}
}

func NewConcurrentWithCapacity[T comparable](cap int) *ConcurrentSet[T] {
	return &ConcurrentSet[T]{set: Set[T]{elements: make(map[T]struct{}, cap)}}
}

// NewConcurrentFromSlice creates a ConcurrentSet from a slice, deduplicating elements.
func NewConcurrentFromSlice[T comparable](items []T) *ConcurrentSet[T] {
	cs := NewConcurrentWithCapacity[T](len(items))
	for _, item := range items {
		cs.set.elements[item] = struct{}{}
	}
	return cs
}

// FromSet wraps an existing Set[T] to make it concurrency-safe going forward.
// The original Set should not be used directly after this call, since both
// would share the same underlying map without synchronization.
func FromSet[T comparable](s *Set[T]) *ConcurrentSet[T] {
	if s == nil || s.elements == nil {
		return NewConcurrent[T]()
	}
	return &ConcurrentSet[T]{set: Set[T]{elements: s.elements}}
}

// --- Query / Inspection ---

func (cs *ConcurrentSet[T]) Contains(e T) bool {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	return cs.set.Contains(e)
}

func (cs *ConcurrentSet[T]) Size() int {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	return cs.set.Size()
}

func (cs *ConcurrentSet[T]) IsEmpty() bool {
	if cs == nil {
		return true
	}
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	return cs.set.IsEmpty()
}

// --- Modification ---

func (cs *ConcurrentSet[T]) Insert(e T) bool {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	return cs.set.Insert(e)
}

func (cs *ConcurrentSet[T]) Remove(e T) bool {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	return cs.set.Remove(e)
}

func (cs *ConcurrentSet[T]) Clear() {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.set.Clear()
}

func (cs *ConcurrentSet[T]) Take(e T) (T, bool) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	return cs.set.Take(e)
}

func (cs *ConcurrentSet[T]) Retain(f func(T) bool) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.set.Retain(f)
}

// --- Set Operations ---
// These return a plain *Set[T] (a point-in-time snapshot), not another
// ConcurrentSet. Wrap the result with FromSet if you need to keep sharing it.

func (cs *ConcurrentSet[T]) Union(other *ConcurrentSet[T]) *Set[T] {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	other.mu.RLock()
	defer other.mu.RUnlock()
	return cs.set.Union(&other.set)
}

func (cs *ConcurrentSet[T]) Intersection(other *ConcurrentSet[T]) *Set[T] {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	other.mu.RLock()
	defer other.mu.RUnlock()
	return cs.set.Intersection(&other.set)
}

func (cs *ConcurrentSet[T]) Difference(other *ConcurrentSet[T]) *Set[T] {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	other.mu.RLock()
	defer other.mu.RUnlock()
	return cs.set.Difference(&other.set)
}

// --- Set Comparisons ---

func (cs *ConcurrentSet[T]) Equal(other *ConcurrentSet[T]) bool {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	other.mu.RLock()
	defer other.mu.RUnlock()
	return cs.set.Equal(&other.set)
}

// --- Iteration ---

// All returns an iterator over a SNAPSHOT of the elements taken under the
// read lock. This avoids holding the lock for the entire iteration (which
// would deadlock if the consumer called another method on cs mid-loop) at
// the cost of not reflecting concurrent mutations made during iteration.
func (cs *ConcurrentSet[T]) All() iter.Seq[T] {
	cs.mu.RLock()
	snapshot := make([]T, 0, len(cs.set.elements))
	for k := range cs.set.elements {
		snapshot = append(snapshot, k)
	}
	cs.mu.RUnlock()

	return func(yield func(T) bool) {
		for _, k := range snapshot {
			if !yield(k) {
				return
			}
		}
	}
}

// ToSlice returns a snapshot of the elements as a slice.
func (cs *ConcurrentSet[T]) ToSlice() []T {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	return cs.set.ToSlice()
}

// --- Escape hatch for batched operations ---

// WithLock executes f while holding the write lock, giving direct access to
// the underlying *Set[T]. Use this to batch several operations atomically
// (e.g. check-then-insert) without paying per-call lock overhead or risking
// a race between the check and the mutation.
//
// f must not retain or leak the *Set[T] it receives beyond the call.
func (cs *ConcurrentSet[T]) WithLock(f func(s *Set[T])) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	f(&cs.set)
}
