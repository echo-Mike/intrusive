package dlist

type (
	// Hook structure to insert/embed into concrete types
	// of elements of doubly-linked list intrusive container
	Hook[T any] struct {
		next, prev *T
	}

	// Head structure of doubly-linked list intrusive container
	DList[T any] struct {
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

// Return previous element if this object is part of some list
// or nil if this object is first in some list or not part of any list
func (h Hook[T]) Prev() *T {
	return h.prev
}

// Initialize hook to empty state.
//
// WARNING: Calling this function on linked Hook will damage DList structure
func (h *Hook[T]) Init() {
	h.next = nil
	h.prev = nil
}

// Create hook in empty state
func NewHook[T any]() Hook[T] {
	return Hook[T]{next: nil, prev: nil}
}

// Create new DList container
func New[T any](hookFunc func(*T) *Hook[T]) DList[T] {
	return DList[T]{hookFunc: hookFunc, size: 0, first: nil, last: nil}
}

// Initialize DList to empty state
func (d *DList[T]) Init() {
	d.first = nil
	d.last = nil
	d.size = 0
}

// Check if DList is empty
func (d DList[T]) Empty() bool {
	return d.size == 0
}

// Get current length of DList
func (d DList[T]) Size() int {
	return d.size
}

// Get current length of DList
func (s DList[T]) Len() int {
	return s.size
}

// Swap content of two DList heads
func (d *DList[T]) Swap(other *DList[T]) {
	other.hookFunc, d.hookFunc = d.hookFunc, other.hookFunc
	other.first, d.first = d.first, other.first
	other.last, d.last = d.last, other.last
	other.size, d.size = d.size, other.size
}

// Return first element in DList
func (d DList[T]) Front() *T {
	return d.first
}

// Return last element in DList
func (d DList[T]) Back() *T {
	return d.last
}

// Insert new element before specified position
func (d *DList[T]) Insert(position, element *T) {
	d.verifyElementNotLinked(element)
	defer d.verifyIsMemberOfCurrent(element)
	defer d.verifyNoCycle()
	defer d.verifySize()

	if position == nil {
		d.PushBack(element)
		return
	}

	d.verifyIsMemberOfCurrent(position)

	hook := d.hookFunc(element)
	posHook := d.hookFunc(position)

	hook.next = position
	hook.prev = posHook.prev

	if posHook.prev != nil {
		d.hookFunc(posHook.prev).next = element
	} else {
		d.first = element
	}

	posHook.prev = element
	d.size++
}

// Remove element from DList
func (d *DList[T]) Erase(element *T) {
	d.verifyNotEmpty()
	d.verifyIsMemberOfCurrent(element)
	defer d.verifyElementNotLinked(element)
	defer d.verifyNoCycle()
	defer d.verifySize()

	hook := d.hookFunc(element)

	if hook.prev != nil {
		d.hookFunc(hook.prev).next = hook.next
	} else {
		d.first = hook.next
	}

	if hook.next != nil {
		d.hookFunc(hook.next).prev = hook.prev
	} else {
		d.last = hook.prev
	}

	hook.Init()
	d.size--
}

// Insert new element at the front of DList
func (d *DList[T]) PushFront(element *T) {
	d.verifyElementNotLinked(element)
	defer d.verifyIsMemberOfCurrent(element)
	defer d.verifyNoCycle()
	defer d.verifySize()

	hook := d.hookFunc(element)

	if d.first == nil {
		d.last = element
		hook.next = nil
		hook.prev = nil
	} else {
		hook.next = d.first
		hook.prev = nil
		d.hookFunc(d.first).prev = element
	}

	d.first = element
	d.size++
}

// Remove and return element from the front of DList
func (d *DList[T]) PopFront() (popped *T) {
	if d.first == nil {
		return nil
	}
	d.verifyNotEmpty()
	defer d.verifyNoCycle()
	defer d.verifySize()

	popped = d.first
	hook := d.hookFunc(popped)

	if hook.next != nil {
		d.hookFunc(hook.next).prev = nil
	} else {
		d.last = nil
	}

	d.first = hook.next
	hook.Init()
	d.size--
	return
}

// Insert new element at the back of DList
func (d *DList[T]) PushBack(element *T) {
	d.verifyElementNotLinked(element)
	defer d.verifyIsMemberOfCurrent(element)
	defer d.verifyNoCycle()
	defer d.verifySize()

	hook := d.hookFunc(element)

	if d.last == nil {
		d.first = element
		hook.next = nil
		hook.prev = nil
	} else {
		hook.next = nil
		hook.prev = d.last
		d.hookFunc(d.last).next = element
	}
	d.last = element

	d.size++
}

// Remove and return element from the back of DList
func (d *DList[T]) PopBack() (popped *T) {
	if d.last == nil {
		return nil
	}
	d.verifyNotEmpty()
	defer d.verifyNoCycle()
	defer d.verifySize()

	popped = d.last
	hook := d.hookFunc(popped)

	if hook.prev != nil {
		d.hookFunc(hook.prev).next = nil
	} else {
		d.first = nil
	}

	d.last = hook.prev
	hook.Init()
	d.size--
	return
}

// Move elements from other DList to be included in current DList at the beginning
func (d *DList[T]) SpliceFront(other *DList[T]) {
	d.Splice(d.first, other)
}

// Move elements from other DList to be included in current DList at the end
func (d *DList[T]) SpliceBack(other *DList[T]) {
	d.Splice(nil, other)
}

// Move elements from other DList to be included in current DList before position
func (d *DList[T]) Splice(position *T, other *DList[T]) {
	if other.first == nil || d == other {
		return
	}
	other.verifyNotEmpty()
	other.verifyNoCycle()
	defer d.verifyNoCycle()
	defer d.verifySize()

	if position == nil {
		// Append to the end
		if d.last != nil {
			d.hookFunc(other.first).prev = d.last
			d.hookFunc(d.last).next = other.first
		} else {
			d.first = other.first
		}
		d.last = other.last
	} else {
		d.verifyIsMemberOfCurrent(position)
		// Insert before position
		posHook := d.hookFunc(position)
		otherFirstHook := d.hookFunc(other.first)
		otherLastHook := d.hookFunc(other.last)

		otherFirstHook.prev = posHook.prev
		otherLastHook.next = position

		if posHook.prev != nil {
			d.hookFunc(posHook.prev).next = other.first
		} else {
			d.first = other.first
		}
		posHook.prev = other.last
	}

	d.size += other.size
	other.Init()
}

// Clear DList and return all currently linked elements as slice.
func (d *DList[T]) Clear() (elements []*T) {
	d.verifyNoCycle()

	elements = make([]*T, 0, d.size)
	e := d.first
	for e != nil {
		elements = append(elements, e)
		h := d.hookFunc(e)
		e = h.Next()
		h.Init()
	}
	d.Init()
	return
}

// Reverse current DList in place
func (d *DList[T]) Reverse() {
	d.verifyNoCycle()
	defer d.verifyNoCycle()
	defer d.verifySize()

	current := d.first
	for current != nil {
		hook := d.hookFunc(current)
		hook.prev, hook.next = hook.next, hook.prev
		current = hook.prev
	}
	d.first, d.last = d.last, d.first
}

// Merge sorted lists into current DList
func (d *DList[T]) Merge(other *DList[T], less func(lhs, rhs *T) bool) {
	if other.first == nil {
		return
	}
	d.verifyNoCycle()
	other.verifyNoCycle()
	defer d.verifyNoCycle()
	defer d.verifySize()

	if d.first == nil {
		d.first = other.first
		d.last = other.last
		d.size = other.size
		other.Init()
		return
	}

	var head, tail *T
	a, b := d.first, other.first

	// Determine the head of the merged list
	if less(a, b) {
		head = a
		tail = a
		a = d.hookFunc(a).next
	} else {
		head = b
		tail = b
		b = d.hookFunc(b).next
	}

	// Merge the two lists
	for a != nil && b != nil {
		if less(a, b) {
			d.hookFunc(tail).next = a
			d.hookFunc(a).prev = tail
			tail = a
			a = d.hookFunc(a).next
		} else {
			d.hookFunc(tail).next = b
			d.hookFunc(b).prev = tail
			tail = b
			b = d.hookFunc(b).next
		}
	}

	// Append the remaining elements
	if a != nil {
		d.hookFunc(tail).next = a
		d.hookFunc(a).prev = tail
		tail = d.last
	} else if b != nil {
		d.hookFunc(tail).next = b
		d.hookFunc(b).prev = tail
		tail = other.last
	}

	d.first = head
	d.last = tail
	d.size += other.size
	other.Init()
}

// Remove elements satisfying predicate and return all removed elements as slice
func (d *DList[T]) RemoveIf(predicate func(value *T) bool) (elements []*T) {
	elements = make([]*T, 0)

	current := d.first
	for current != nil {
		next := d.hookFunc(current).next
		if predicate(current) {
			elements = append(elements, current)
			d.Erase(current)
		}
		current = next
	}
	return
}

// Remove consecutive duplicate elements and return all removed elements as slice
func (d *DList[T]) Unique(less func(lhs, rhs *T) bool) (elements []*T) {
	elements = make([]*T, 0)
	current := d.first
	for current != nil && d.hookFunc(current).next != nil {
		next := d.hookFunc(current).next
		if !less(current, next) && !less(next, current) {
			elements = append(elements, next)
			d.Erase(next)
		} else {
			current = next
		}
	}
	return
}

// Sort current DList in place using merge sort
func (d *DList[T]) Sort(less func(lhs, rhs *T) bool) {
	if d.size <= 1 {
		return
	}
	d.verifyNotEmpty()
	d.verifyNoCycle()
	defer d.verifyNoCycle()
	defer d.verifySize()

	// Split the list into two halves
	mid := d.size / 2
	current := d.first
	for i := 0; i < mid-1; i++ {
		current = d.hookFunc(current).next
	}

	rightFirst := d.hookFunc(current).next
	d.hookFunc(current).next = nil
	d.hookFunc(rightFirst).prev = nil

	left := &DList[T]{
		hookFunc: d.hookFunc,
		first:    d.first,
		last:     current,
		size:     mid,
	}
	right := &DList[T]{
		hookFunc: d.hookFunc,
		first:    rightFirst,
		last:     d.last,
		size:     d.size - mid,
	}

	left.Sort(less)
	right.Sort(less)

	d.Init()
	d.Merge(left, less)
	d.Merge(right, less)
}
