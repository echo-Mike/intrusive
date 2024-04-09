package slist

type (
	// Hook structure to insert/embed into concrete types
	// of elements of singly-linked list intrusive container
	Hook[T any] struct {
		next *T
	}

	Operations[T any] interface {
		Hook(self *T) *Hook[T]
	}

	// Head structure of singly-linked list intrusive container
	SList[Ops Operations[T], T any] struct {
		ops         Ops
		size        int
		first, last *T
	}
)

// Return next element if this object is part of some list
// or nil if this object is last in some list or not part of any list
func (h Hook[T]) Next() *T {
	return h.next
}

// Initialize hook to empty state.
//
// WARNING: Calling this function on linked Hook will damage SList structure
func (h *Hook[T]) Init() {
	h.next = nil
}

// Create hook in empty state
func NewHook[T any]() Hook[T] {
	return Hook[T]{next: nil}
}

// Create new SList container
func New[Ops Operations[T], T any](ops Ops) SList[Ops, T] {
	return SList[Ops, T]{ops: ops, size: 0, first: nil, last: nil}
}

// Initialize SList to empty state
func (s *SList[Ops, T]) Init() {
	s.first = nil
	s.last = nil
	s.size = 0
}

// This function is not exported as it do not swap ops and exporting it will lead to one of:
//
// 1. Unexpected behavior from users (as ops are ton swapped)
//
// 2. Requirement for Ops type to implement swap or be trivially swappable (via copy)
func (s *SList[Ops, T]) swap(other *SList[Ops, T]) {
	other.first, s.first = s.first, other.first
	other.last, s.last = s.last, other.last
	other.size, s.size = s.size, other.size
}

// Get current length of SList
func (s SList[Ops, T]) Len() int {
	return s.size
}

// Insert new element after specified. Position SHOULD be part of current SList
func (s *SList[Ops, T]) InsertAfter(position, element *T) {
	s.verifyNotEmpty()
	s.verifyElementNotLinked(element)
	s.verifyIsMemberOfCurrent(position)
	defer s.verifyIsMemberOfCurrent(element)
	defer s.verifyNoCycle()
	defer s.verifySize()

	s.ops.Hook(element).next = s.ops.Hook(position).next
	s.ops.Hook(position).next = element
	if s.last == position {
		s.last = element
	}
	s.size++
}

// Unlink and return element after specified. Position SHOULD be part of current SList. Return nil if position is at the end of SList
func (s *SList[Ops, T]) RemoveAfter(position *T) (popped *T) {
	s.verifyNotEmpty()
	s.verifyNoCycle()
	s.verifyIsMemberOfCurrent(position)
	defer s.verifySize()

	if popped = s.ops.Hook(position).next; popped != nil {
		s.ops.Hook(position).next = s.ops.Hook(popped).next
		if s.last == popped {
			s.last = position
		}
		s.size--
		s.ops.Hook(popped).Init()
	}
	return
}

// Move elements from other SList to be included in current SList after position
func (s *SList[Ops, T]) SpliceAfter(position *T, other *SList[Ops, T]) {
	if other.first == nil {
		return
	}
	other.verifyNotEmpty()
	other.verifyNoCycle()
	s.verifyNotEmpty()
	s.verifyNoCycle()
	s.verifyIsMemberOfCurrent(position)
	defer s.verifyNoCycle()
	defer s.verifySize()

	s.ops.Hook(other.last).next = s.ops.Hook(position).next
	s.ops.Hook(position).next = other.first
	if s.last == position {
		s.last = other.last
	}
	s.size += other.size
	other.Init()
}

// Return first element in SList
func (s SList[Ops, T]) Front() *T {
	return s.first
}

// Insert new element into the head of SList
func (s *SList[Ops, T]) PushFront(element *T) {
	s.verifyElementNotLinked(element)
	defer s.verifyIsMemberOfCurrent(element)
	defer s.verifyNoCycle()
	defer s.verifySize()

	if s.first == nil {
		s.last = element
	} else {
		s.ops.Hook(element).next = s.first
	}
	s.first = element
	s.size++
}

// Unlink and return element at the head of a SList
func (s *SList[Ops, T]) PopFront() (popped *T) {
	if s.first == nil {
		return nil
	}
	s.verifyNotEmpty()
	defer func(popped **T) { s.verifyElementNotLinked(*popped) }(&popped)
	defer s.verifyNoCycle()
	defer s.verifySize()

	popped = s.first
	if s.first == s.last {
		s.first = nil
		s.last = nil
	} else {
		s.first = s.ops.Hook(popped).next
	}
	s.size--
	s.ops.Hook(popped).Init()
	return
}

// Move elements from other SList to be included at the head of current SList
func (s *SList[Ops, T]) SpliceFront(other *SList[Ops, T]) {
	other.SpliceBack(s)
	s.swap(other)
}

// Return last element in SList
func (s SList[Ops, T]) Back() *T {
	return s.last
}

// Insert new element into the tail of current SList
func (s *SList[Ops, T]) PushBack(element *T) {
	s.verifyElementNotLinked(element)
	defer s.verifyIsMemberOfCurrent(element)
	defer s.verifyNoCycle()
	defer s.verifySize()

	if s.last == nil {
		s.first = element
	} else {
		s.ops.Hook(s.last).next = element
	}
	s.last = element
	s.size++
	s.ops.Hook(element).Init()
}

// Move elements from other SList to be included at the tail of current SList
func (s *SList[Ops, T]) SpliceBack(other *SList[Ops, T]) {
	if s.first == nil {
		s.swap(other)
		return
	}
	s.SpliceAfter(s.last, other)
}

// Clear SList and return all currently linked elements as slice.
// Use Init() to clear SList without allocations
func (s *SList[Ops, T]) Clear() (elements []*T) {
	s.verifyNoCycle()

	elements = make([]*T, 0, s.size)
	e := s.first
	for e != nil {
		elements = append(elements, e)
		h := s.ops.Hook(e)
		e = h.Next()
		h.Init()
	}
	s.Init()
	return
}

// Reverse current SList in place
func (s *SList[Ops, T]) Reverse() {
	s.verifyNoCycle()
	defer s.verifyNoCycle()
	defer s.verifySize()

	var prev *T = nil
	e := s.first
	for e != nil {
		next := s.ops.Hook(e).next
		s.ops.Hook(e).next = prev
		prev = e
		e = next
	}
	s.first, s.last = s.last, s.first
}

func (s *SList[Ops, T]) median(l *T) (slow *T) {
	slow = l
	fast := s.ops.Hook(l).next

	for fast != nil {
		fast = s.ops.Hook(fast).next
		if fast != nil {
			slow = s.ops.Hook(slow).next
			fast = s.ops.Hook(fast).next
		}
	}
	return
}

// Return element in the center of SList.
//
// If list size is even last element in first half is returned
func (s *SList[Ops, T]) Median() (median *T) {
	s.verifySize()
	s.verifyNoCycle()

	median = s.first
	half := (s.size - s.size%2) / 2
	for i := 0; i < half; i++ {
		median = s.ops.Hook(median).next
	}
	return
}

func (s *SList[Ops, T]) merge(a, b *T, less func(lhs, rhs *T) bool) (first, last *T) {
	if a != nil && (b != nil && less(a, b) || b == nil) {
		first = a
		a = s.ops.Hook(a).next
	} else {
		first = b
		if b != nil {
			b = s.ops.Hook(b).next
		}
	}
	last = first
	for a != nil || b != nil {
		if a != nil && (b != nil && less(a, b) || b == nil) {
			s.ops.Hook(last).next = a
			last = a
			a = s.ops.Hook(a).next
		} else {
			s.ops.Hook(last).next = b
			last = b
			b = s.ops.Hook(b).next
		}
	}
	return
}

// Merge sorted lists into current SList
func (s *SList[Ops, T]) Merge(other *SList[Ops, T], less func(lhs, rhs *T) bool) {
	s.verifyNoCycle()
	other.verifyNoCycle()
	defer s.verifyNoCycle()
	defer s.verifySize()

	s.first, s.last = s.merge(s.first, other.first, less)
	s.size += other.size
	other.Init()
}

// Merge sort based implementation
func (s *SList[Ops, T]) sort(head *T, less func(lhs, rhs *T) bool) (first, last *T) {
	if head == nil || s.ops.Hook(head).next == nil {
		return
	}

	m := s.median(head)
	tail := s.ops.Hook(m).next
	s.ops.Hook(m).next = nil

	head, _ = s.sort(head, less)
	tail, _ = s.sort(tail, less)

	first, last = s.merge(head, tail, less)
	return
}

// Sort current SList in place
func (s *SList[Ops, T]) Sort(less func(lhs, rhs *T) bool) {
	if s.first == nil || s.first == s.last {
		return
	}
	s.verifyNotEmpty()
	s.verifyNoCycle()
	defer s.verifyNoCycle()
	defer s.verifySize()

	// One iteration here just will be a bit faster as Median is based on size
	m := s.Median()
	tail := s.ops.Hook(m).next
	s.ops.Hook(m).next = nil

	head, _ := s.sort(s.first, less)
	tail, _ = s.sort(tail, less)

	s.first, s.last = s.merge(head, tail, less)
}

// Iterate over current SList applying f to each element and their parent
//
// This function iterates over SList safely meaning users can delete cur in f
func (s *SList[Ops, T]) adjacent(f func(prev, cur *T)) {
	s.verifyNotEmpty()

	var p *T = s.first
	e := s.ops.Hook(p).next
	for e != nil {
		f(p, e)
		n := s.ops.Hook(p).next
		if n == e {
			p = e
			e = s.ops.Hook(e).next
		} else {
			e = n
		}
	}
}

// Remove consecutive duplicate elements and return all removed elements as slice
func (s *SList[Ops, T]) Unique(less func(lhs, rhs *T) bool) (elements []*T) {
	elements = make([]*T, 0)
	if s.first == nil {
		return
	}
	s.adjacent(func(prev, cur *T) {
		if !(less(prev, cur) || less(cur, prev)) {
			elements = append(elements, s.RemoveAfter(prev))
		}
	})
	return
}

// Remove elements satisfying predicate and return all removed elements as slice
func (s *SList[Ops, T]) RemoveIf(predicate func(value *T) bool) (elements []*T) {
	elements = make([]*T, 0)
	if s.first == nil {
		return
	}
	front := s.first
	for front != nil && predicate(front) {
		front = s.ops.Hook(front).next
		elements = append(elements, s.PopFront())
	}
	if front != nil {
		s.adjacent(func(prev, cur *T) {
			if predicate(cur) {
				elements = append(elements, s.RemoveAfter(prev))
			}
		})
	}
	return
}
