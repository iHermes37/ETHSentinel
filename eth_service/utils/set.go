package utils

type Set[T comparable] map[T]struct{}

func NewSet[T comparable](items ...T) Set[T] {
	s := make(Set[T])
	for _, item := range items {
		s.Add(item)
	}
	return s
}

func (s Set[T]) Add(item T) {
	s[item] = struct{}{}
}

func (s Set[T]) Remove(item T) {
	delete(s, item)
}

func (s Set[T]) Contains(item T) bool {
	_, ok := s[item]
	return ok
}

func (s Set[T]) Size() int {
	return len(s)
}

// s := NewSet("apple", "banana")
// s.Add("cherry")
// fmt.Println(s.Contains("banana")) // true
// s.Remove("banana")
// fmt.Println(s.Size()) // 2
