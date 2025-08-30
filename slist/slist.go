package slist

type (
	// Hook structure to insert/embed into concrete types
	// of elements of singly-linked list intrusive container
	Hook[T any] struct {
		next *T
	}

	// Head structure of singly-linked list intrusive container
	SList[T any] struct {
		hookFunc    func(*T) *Hook[T]
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
func New[T any](hookFunc func(*T) *Hook[T]) SList[T] {
	return SList[T]{hookFunc: hookFunc, size: 0, first: nil, last: nil}
}

// Initialize SList to empty state
func (s *SList[T]) Init() {
	s.first = nil
	s.last = nil
	s.size = 0
}

// Swap content of two SList heads
func (s *SList[T]) Swap(other *SList[T]) {
	other.hookFunc, s.hookFunc = s.hookFunc, other.hookFunc
	other.first, s.first = s.first, other.first
	other.last, s.last = s.last, other.last
	other.size, s.size = s.size, other.size
}

// Get current length of SList
func (s SList[T]) Len() int {
	return s.size
}

// Insert new element after specified. Position SHOULD be part of current SList
func (s *SList[T]) InsertAfter(position, element *T) {
	s.verifyNotEmpty()
	s.verifyElementNotLinked(element)
	s.verifyIsMemberOfCurrent(position)
	defer s.verifyIsMemberOfCurrent(element)
	defer s.verifyNoCycle()
	defer s.verifySize()

	s.hookFunc(element).next = s.hookFunc(position).next
	s.hookFunc(position).next = element
	if s.last == position {
		s.last = element
	}
	s.size++
}

// Unlink and return element after specified. Position SHOULD be part of current SList. Return nil if position is at the end of SList
func (s *SList[T]) RemoveAfter(position *T) (popped *T) {
	s.verifyNotEmpty()
	s.verifyNoCycle()
	s.verifyIsMemberOfCurrent(position)
	defer s.verifySize()

	if popped = s.hookFunc(position).next; popped != nil {
		s.hookFunc(position).next = s.hookFunc(popped).next
		if s.last == popped {
			s.last = position
		}
		s.size--
		s.hookFunc(popped).Init()
	}
	return
}

// Move elements from other SList to be included in current SList after position
func (s *SList[T]) SpliceAfter(position *T, other *SList[T]) {
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

	s.hookFunc(other.last).next = s.hookFunc(position).next
	s.hookFunc(position).next = other.first
	if s.last == position {
		s.last = other.last
	}
	s.size += other.size
	other.Init()
}

// Return first element in SList
func (s SList[T]) Front() *T {
	return s.first
}

// Insert new element into the head of SList
func (s *SList[T]) PushFront(element *T) {
	s.verifyElementNotLinked(element)
	defer s.verifyIsMemberOfCurrent(element)
	defer s.verifyNoCycle()
	defer s.verifySize()

	if s.first == nil {
		s.last = element
	} else {
		s.hookFunc(element).next = s.first
	}
	s.first = element
	s.size++
}

// Unlink and return element at the head of a SList
func (s *SList[T]) PopFront() (popped *T) {
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
		s.first = s.hookFunc(popped).next
	}
	s.size--
	s.hookFunc(popped).Init()
	return
}

// Move elements from other SList to be included at the head of current SList
func (s *SList[T]) SpliceFront(other *SList[T]) {
	other.SpliceBack(s)
	s.Swap(other)
}

// Return last element in SList
func (s SList[T]) Back() *T {
	return s.last
}

// Insert new element into the tail of current SList
func (s *SList[T]) PushBack(element *T) {
	s.verifyElementNotLinked(element)
	defer s.verifyIsMemberOfCurrent(element)
	defer s.verifyNoCycle()
	defer s.verifySize()

	if s.last == nil {
		s.first = element
	} else {
		s.hookFunc(s.last).next = element
	}
	s.last = element
	s.size++
	s.hookFunc(element).Init()
}

// Move elements from other SList to be included at the tail of current SList
func (s *SList[T]) SpliceBack(other *SList[T]) {
	if s.first == nil {
		s.Swap(other)
		return
	}
	s.SpliceAfter(s.last, other)
}

// Clear SList and return all currently linked elements as slice.
// Use Init() to clear SList without allocations
func (s *SList[T]) Clear() (elements []*T) {
	s.verifyNoCycle()

	elements = make([]*T, 0, s.size)
	e := s.first
	for e != nil {
		elements = append(elements, e)
		h := s.hookFunc(e)
		e = h.Next()
		h.Init()
	}
	s.Init()
	return
}

// Reverse current SList in place
func (s *SList[T]) Reverse() {
	s.verifyNoCycle()
	defer s.verifyNoCycle()
	defer s.verifySize()

	var prev *T = nil
	e := s.first
	for e != nil {
		next := s.hookFunc(e).next
		s.hookFunc(e).next = prev
		prev = e
		e = next
	}
	s.first, s.last = s.last, s.first
}

func (s *SList[T]) median(l *T) (slow *T) {
	slow = l
	fast := s.hookFunc(l).next

	for fast != nil {
		fast = s.hookFunc(fast).next
		if fast != nil {
			slow = s.hookFunc(slow).next
			fast = s.hookFunc(fast).next
		}
	}
	return
}

// Return element in the center of SList.
//
// If list size is even last element in first half is returned
func (s *SList[T]) Median() (median *T) {
	s.verifySize()
	s.verifyNoCycle()

	median = s.first
	half := (s.size + s.size%2) / 2
	for i := 0; i < half-1; i++ {
		median = s.hookFunc(median).next
	}
	return
}

func (s *SList[T]) merge(a, b *T, less func(lhs, rhs *T) bool) (first, last *T) {
	if a != nil && (b != nil && less(a, b) || b == nil) {
		first = a
		a = s.hookFunc(a).next
	} else {
		first = b
		if b != nil {
			b = s.hookFunc(b).next
		}
	}
	last = first
	for a != nil || b != nil {
		if a != nil && (b != nil && less(a, b) || b == nil) {
			s.hookFunc(last).next = a
			last = a
			a = s.hookFunc(a).next
		} else {
			s.hookFunc(last).next = b
			last = b
			b = s.hookFunc(b).next
		}
	}
	return
}

// Merge sorted lists into current SList
func (s *SList[T]) Merge(other *SList[T], less func(lhs, rhs *T) bool) {
	s.verifyNoCycle()
	other.verifyNoCycle()
	defer s.verifyNoCycle()
	defer s.verifySize()

	s.first, s.last = s.merge(s.first, other.first, less)
	s.size += other.size
	other.Init()
}

// Merge sort based implementation
func (s *SList[T]) sort(head *T, less func(lhs, rhs *T) bool) (first, last *T) {
	if head == nil || s.hookFunc(head).next == nil {
		first = head
		last = first
		return
	}

	m := s.median(head)
	tail := s.hookFunc(m).next
	s.hookFunc(m).next = nil

	head, _ = s.sort(head, less)
	tail, _ = s.sort(tail, less)

	first, last = s.merge(head, tail, less)
	return
}

// Sort current SList in place
func (s *SList[T]) Sort(less func(lhs, rhs *T) bool) {
	if s.first == nil || s.first == s.last {
		return
	}
	s.verifyNotEmpty()
	s.verifyNoCycle()
	defer s.verifyNoCycle()
	defer s.verifySize()

	// One iteration here just will be a bit faster as Median is based on size
	m := s.Median()
	tail := s.hookFunc(m).next
	s.hookFunc(m).next = nil

	head, _ := s.sort(s.first, less)
	tail, _ = s.sort(tail, less)

	s.first, s.last = s.merge(head, tail, less)
}

// Iterate over current SList applying f to each element and their parent
//
// This function iterates over SList safely meaning users can delete cur in f
func (s *SList[T]) adjacent(f func(prev, cur *T)) {
	s.verifyNotEmpty()

	var p *T = s.first
	e := s.hookFunc(p).next
	for e != nil {
		f(p, e)
		n := s.hookFunc(p).next
		if n == e {
			p = e
			e = s.hookFunc(e).next
		} else {
			e = n
		}
	}
}

// Remove consecutive duplicate elements and return all removed elements as slice
func (s *SList[T]) Unique(less func(lhs, rhs *T) bool) (elements []*T) {
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
func (s *SList[T]) RemoveIf(predicate func(value *T) bool) (elements []*T) {
	elements = make([]*T, 0)
	if s.first == nil {
		return
	}
	front := s.first
	for front != nil && predicate(front) {
		front = s.hookFunc(front).next
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