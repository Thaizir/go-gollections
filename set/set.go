package set

type Set[T comparable] struct {
	Elements map[T]struct{}
}

// Constructors
func New[T comparable]() *Set[T] {
	return &Set[T]{
		Elements: make(map[T]struct{}),
	}
}

func NewWithCapacity[T comparable](cap int) *Set[T] {
	return &Set[T]{
		Elements: make(map[T]struct{}, cap),
	}
}

// Basic Operations

// Inserts a new element into the Set, returns false if the element was already present
func (s *Set[T]) Insert(e T) bool {
	if s.Elements == nil {
		s.Elements = make(map[T]struct{})
	}
	if _, exists := s.Elements[e]; !exists {
		s.Elements[e] = struct{}{}
		return true
	}
	return false
}

// Removes an element from the set. Returns true if the element was present
func (s *Set[T]) Remove(e T) bool {
	if s.Elements == nil {
		return false
	}

	if _, exists := s.Elements[e]; !exists {
		return false
	}

	delete(s.Elements, e)
	return true
}

// Verify if and element exists in the Set
func (s *Set[T]) Contains(e T) bool {
	if s.Elements == nil {
		return false
	}

	if _, exists := s.Elements[e]; !exists {
		return false
	}
	return true
}

// Delete all the elements
func (s *Set[T]) Clear() {
	if s.Elements != nil {
		clear(s.Elements)
	}
}

// Returns the cardinality of the Set
func (s *Set[T]) Size() int {
	if s.Elements != nil {
		return len(s.Elements)
	}
	return 0
}

// Checks if the Set does not contain elements
func (s *Set[T]) IsEmpty() bool {
	if s.Elements != nil {
		return true
	}
	return false
}

// Set Algebra

// Returns a new set with the union of both
func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	return other
}

// Returns a new set with the elements in common
func (s *Set[T]) Intersection(other *Set[T]) *Set[T] {
	return other
}

// Returns the elements from the original set that are not present in the other one
func (s *Set[T]) Difference(other *Set[T]) *Set[T] {
	return other
}

// Returns the elements that are in each set but they are not in common
func (s *Set[T]) SymmetricDifference(other *Set[T]) *Set[T] {
	return other
}

// Relations between sets (predicates)

// Filtering and Iteration
