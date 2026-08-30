package set

import (
	"iter"
	"slices"
	"sync"
)

type Set[T comparable] struct {
	mu       sync.RWMutex
	elements map[T]struct{}
}

// Helper method to clone maps
func cloneMap[T comparable](m map[T]struct{}) *Set[T] {
	if m == nil {
		return New[T]()
	}
	res := NewWithCapacity[T](len(m))
	for k := range m {
		res.elements[k] = struct{}{}
	}
	return res
}

// Constructors
func New[T comparable]() *Set[T] {
	return &Set[T]{
		elements: make(map[T]struct{}),
	}
}

func NewWithCapacity[T comparable](cap int) *Set[T] {
	return &Set[T]{
		elements: make(map[T]struct{}, cap),
	}
}

// NewFromSlice creates a set from a slice, deduplicating elements
func NewFromSlice[T comparable](items []T) *Set[T] {
	result := NewWithCapacity[T](len(items))
	for _, k := range items {
		result.elements[k] = struct{}{}
	}
	return result
}

// --- Query / Inspection ---

// Verify if an element exists in the Set
func (s *Set[T]) Contains(e T) bool {
	if s == nil || s.elements == nil {
		return false
	}
	_, exists := s.elements[e]
	return exists
}

// Returns the cardinality of the Set
func (s *Set[T]) Size() int {
	if s == nil || s.elements != nil {
		return len(s.elements)
	}
	return 0
}

// Checks if the Set does not contain elements
func (s *Set[T]) IsEmpty() bool {
	return s == nil || len(s.elements) == 0
}

// --- Modification ---

// Inserts a new element into the Set, returns false if the element was already present
func (s *Set[T]) Insert(e T) bool {
	if s == nil || s.elements == nil {
		s.elements = make(map[T]struct{})
	}
	if _, exists := s.elements[e]; !exists {
		s.elements[e] = struct{}{}
		return true
	}
	return false
}

// Removes an element from the set. Returns true if the element was present
func (s *Set[T]) Remove(e T) bool {
	if s == nil || s.elements == nil {
		return false
	}
	if _, exists := s.elements[e]; exists {
		delete(s.elements, e)
		return true
	}
	return false
}

// Delete all the elements
func (s *Set[T]) Clear() {
	if s == nil || s.elements != nil {
		clear(s.elements)
	}
}

// Take removes and returns element if present
func (s *Set[T]) Take(e T) (T, bool) {

	var zero T

	if s == nil || s.elements == nil {
		return zero, false
	}
	if _, exists := s.elements[e]; !exists {
		return zero, false
	}
	delete(s.elements, e)
	return e, true
}

// Retain keeps only the elements for which f returns true. Mutates the Set in-place
func (s *Set[T]) Retain(f func(T) bool) {
	if s == nil || s.elements == nil || f == nil {
		return
	}

	for k := range s.elements {
		if !f(k) {
			delete(s.elements, k)
		}
	}
}

// Set Operations

// Returns a new set with the union of both
func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	if s == nil || s.elements == nil {
		if other == nil || other.elements == nil {
			return New[T]()
		}
		return cloneMap(other.elements)
	}
	if other == nil || other.elements == nil {
		return cloneMap(s.elements)
	}

	union := NewWithCapacity[T](len(s.elements) + len(other.elements))

	for k := range s.elements {
		union.elements[k] = struct{}{}
	}
	for k := range other.elements {
		union.elements[k] = struct{}{}
	}

	return union
}

// Returns a new set with the elements in common
func (s *Set[T]) Intersection(other *Set[T]) *Set[T] {
	if s == nil || s.elements == nil || other == nil || other.elements == nil {
		return New[T]()
	}

	small, large := s.elements, other.elements
	if len(small) > len(large) {
		small, large = large, small
	}

	intersection := NewWithCapacity[T](len(small))
	for k := range small {
		if _, exists := large[k]; exists {
			intersection.elements[k] = struct{}{}
		}
	}
	return intersection
}

// Returns the elements from the original set that are not present in the other one
func (s *Set[T]) Difference(other *Set[T]) *Set[T] {
	if s == nil || s.elements == nil {
		return New[T]()
	}

	if other == nil || other.elements == nil {
		return cloneMap(s.elements)
	}

	result := NewWithCapacity[T](len(s.elements))
	for k := range s.elements {
		if _, exists := other.elements[k]; !exists {
			result.elements[k] = struct{}{}
		}
	}

	return result
}

// Returns the elements that are in each set but they are not in common
func (s *Set[T]) SymmetricDifference(other *Set[T]) *Set[T] {
	if s == nil || s.elements == nil {
		if other == nil || other.elements == nil {
			return New[T]()
		}
		return cloneMap(other.elements)
	}
	if other == nil || other.elements == nil {
		return cloneMap(s.elements)
	}

	result := NewWithCapacity[T](len(s.elements) + len(other.elements))

	for k := range s.elements {
		if _, exists := other.elements[k]; !exists {
			result.elements[k] = struct{}{}
		}
	}

	for k := range other.elements {
		if _, exists := s.elements[k]; !exists {
			result.elements[k] = struct{}{}
		}
	}

	return result
}

// --- Set Comparisons ---
// IsDisjoint reports whether s and other share no elements
func (s *Set[T]) IsDisjoint(other *Set[T]) bool {
	if s == nil || s.elements == nil || other == nil || other.elements == nil {
		return true
	}

	small, large := s.elements, other.elements
	if len(small) > len(large) {
		small, large = large, small
	}

	for k := range small {
		if _, exists := large[k]; exists {
			return false
		}
	}

	return true
}

// IsSubset reports whether all elements of s are in other (s ⊆ other)
func (s *Set[T]) IsSubset(other *Set[T]) bool {
	if s == nil || len(s.elements) == 0 {
		return true
	}
	if other == nil || len(s.elements) > len(other.elements) {
		return false
	}

	for k := range s.elements {
		if _, exists := other.elements[k]; !exists {
			return false
		}
	}
	return true
}

// IsSuperset reports whether s contains all elements of other (s ⊇ other)
func (s *Set[T]) IsSuperset(other *Set[T]) bool {
	if other == nil {
		return true
	}
	return other.IsSubset(s)
}

// Equal reports whether s and other contain exactly the same elements
func (s *Set[T]) Equal(other *Set[T]) bool {
	if s == other {
		return true
	}
	if s == nil || other == nil {
		return false
	}
	if len(s.elements) != len(other.elements) {
		return false
	}

	for k := range s.elements {
		if _, exists := other.elements[k]; !exists {
			return false
		}
	}
	return true
}

// --- Iteration ---

// All returns an iterator over the elements of the set (Go 1.23+ range-over-func)
func (s *Set[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		if s == nil {
			return
		}
		for k := range s.elements {
			if !yield(k) {
				return
			}
		}
	}
}

// Drain removes all elements, returning an iterator over the removed values
func (s *Set[T]) Drain() iter.Seq[T] {
	return func(yield func(T) bool) {
		if s == nil || s.elements == nil {
			return
		}
		for k := range s.elements {
			delete(s.elements, k)
			if !yield(k) {
				return
			}
		}
	}
}

// --- Conversion ---

// ToSlice returns the elements of the set as a slice
func (s *Set[T]) ToSlice() []T {
	if s == nil {
		return []T{}
	}

	return slices.Collect(s.All())
}
