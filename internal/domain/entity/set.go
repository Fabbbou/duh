package entity

import "fmt"

// Set is a generic set type
type Set[T comparable] struct {
	items map[T]struct{}
}

// NewSet creates a new empty set
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{items: make(map[T]struct{})}
}

// Add inserts an element
func (s *Set[T]) Add(item T) {
	s.items[item] = struct{}{}
}

// Remove deletes an element
func (s *Set[T]) Remove(item T) {
	delete(s.items, item)
}

// Contains checks if an element exists
func (s *Set[T]) Contains(item T) bool {
	_, exists := s.items[item]
	return exists
}

// Values returns all elements as a slice
func (s *Set[T]) Values() []T {
	vals := make([]T, 0, len(s.items))
	for k := range s.items {
		vals = append(vals, k)
	}
	return vals
}

func main() {
	s := NewSet[string]()
	s.Add("apple")
	s.Add("banana")
	fmt.Println(s.Contains("apple")) // true
	fmt.Println(s.Contains("pear"))  // false
	s.Remove("banana")
	fmt.Println(s.Values()) // [apple]
}
