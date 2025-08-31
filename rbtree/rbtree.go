package rbtree

type color int

const (
	black color = 1
	red   color = 2
)

type (
	// Hook contains tree structure information for a value
	Hook[T any] struct {
		left, parent, right *T
		color               color
	}

	// RbTree implements a red-black tree data structure.
	// This structure have a set semantic - meaning the total order
	// of element as compared by lessFunc should not change while
	// it is inside tree
	RbTree[T any] struct {
		hookFunc          func(*T) *Hook[T]
		lessFunc          func(*T, *T) bool
		size              int
		first, root, last *T
	}
)

// Initialize hook to empty state.
//
// WARNING: Calling this function on linked Hook will damage RbTree structure
func (h *Hook[T]) Init() {
	h.left = nil
	h.right = nil
	h.parent = nil
	h.color = black
}

// NewHook creates a new initialized Hook
func NewHook[T any]() Hook[T] {
	return Hook[T]{left: nil, parent: nil, right: nil, color: black}
}

// NewRbTree creates a new Red-Black Tree
func NewRbTree[T any](hookFunc func(*T) *Hook[T], lessFunc func(*T, *T) bool) *RbTree[T] {
	return &RbTree[T]{
		hookFunc: hookFunc,
		lessFunc: lessFunc,
	}
}

// Next returns the next node in in-order traversal
func (t RbTree[T]) Next(node *T) *T {
	t.verifyIsMemberOfCurrent(node)
	return t.next(node)
}

// Prev returns the previous node in in-order traversal
func (t RbTree[T]) Prev(node *T) *T {
	t.verifyIsMemberOfCurrent(node)
	return t.prev(node)
}

// Init initializes the tree to empty state
func (t *RbTree[T]) Init() {
	t.root = nil
	t.first = nil
	t.last = nil
	t.size = 0
}

func (t RbTree[T]) getHook(node *T) *Hook[T] {
	if node == nil {
		return nil
	}
	return t.hookFunc(node)
}

func (t RbTree[T]) left(node *T) *T {
	return t.getHook(node).left
}

func (t RbTree[T]) setLeft(node *T, left *T) {
	t.getHook(node).left = left
}

func (t RbTree[T]) right(node *T) *T {
	return t.getHook(node).right
}

func (t RbTree[T]) setRight(node *T, right *T) {
	t.getHook(node).right = right
}

func (t RbTree[T]) parent(node *T) *T {
	return t.getHook(node).parent
}

func (t RbTree[T]) setParent(node *T, parent *T) {
	t.getHook(node).parent = parent
}

func (t RbTree[T]) color(node *T) color {
	if node == nil {
		return black
	}
	return t.getHook(node).color
}

func (t RbTree[T]) setColor(node *T, color color) {
	t.getHook(node).color = color
}

func (t RbTree[T]) min(node *T) *T {
	for t.left(node) != nil {
		node = t.left(node)
	}
	return node
}

func (t RbTree[T]) max(node *T) *T {
	for t.right(node) != nil {
		node = t.right(node)
	}
	return node
}

func (t RbTree[T]) next(node *T) *T {
	if t.right(node) != nil {
		return t.min(t.right(node))
	}
	parent := t.parent(node)
	for parent != nil && node == t.right(parent) {
		node = parent
		parent = t.parent(parent)
	}
	return parent
}

func (t RbTree[T]) prev(node *T) *T {
	if t.left(node) != nil {
		return t.max(t.left(node))
	}
	parent := t.parent(node)
	for parent != nil && node == t.left(parent) {
		node = parent
		parent = t.parent(parent)
	}
	return parent
}

func (t *RbTree[T]) rotateLeft(x *T) {
	y := t.right(x)
	t.setRight(x, t.left(y))
	if t.left(y) != nil {
		t.setParent(t.left(y), x)
	}
	t.setParent(y, t.parent(x))
	if t.parent(x) == nil {
		t.root = y
	} else if x == t.left(t.parent(x)) {
		t.setLeft(t.parent(x), y)
	} else {
		t.setRight(t.parent(x), y)
	}
	t.setLeft(y, x)
	t.setParent(x, y)
}

func (t *RbTree[T]) rotateRight(x *T) {
	y := t.left(x)
	t.setLeft(x, t.right(y))
	if t.right(y) != nil {
		t.setParent(t.right(y), x)
	}
	t.setParent(y, t.parent(x))
	if t.parent(x) == nil {
		t.root = y
	} else if x == t.right(t.parent(x)) {
		t.setRight(t.parent(x), y)
	} else {
		t.setLeft(t.parent(x), y)
	}
	t.setRight(y, x)
	t.setParent(x, y)
}

func (t *RbTree[T]) insertFixup(z *T) {
	for t.parent(z) != nil && t.color(t.parent(z)) == red {
		if t.parent(z) == t.left(t.parent(t.parent(z))) {
			y := t.right(t.parent(t.parent(z)))
			if t.color(y) == red {
				t.setColor(t.parent(z), black)
				t.setColor(y, black)
				t.setColor(t.parent(t.parent(z)), red)
				z = t.parent(t.parent(z))
			} else {
				if z == t.right(t.parent(z)) {
					z = t.parent(z)
					t.rotateLeft(z)
				}
				t.setColor(t.parent(z), black)
				t.setColor(t.parent(t.parent(z)), red)
				t.rotateRight(t.parent(t.parent(z)))
			}
		} else {
			y := t.left(t.parent(t.parent(z)))
			if t.color(y) == red {
				t.setColor(t.parent(z), black)
				t.setColor(y, black)
				t.setColor(t.parent(t.parent(z)), red)
				z = t.parent(t.parent(z))
			} else {
				if z == t.left(t.parent(z)) {
					z = t.parent(z)
					t.rotateRight(z)
				}
				t.setColor(t.parent(z), black)
				t.setColor(t.parent(t.parent(z)), red)
				t.rotateLeft(t.parent(t.parent(z)))
			}
		}
	}
	t.setColor(t.root, black)
}

func (t *RbTree[T]) deleteFixup(x *T, parentOfX *T) {
	for x != t.root && t.color(x) == black {
		if x == t.left(parentOfX) {
			w := t.right(parentOfX)
			if t.color(w) == red {
				t.setColor(w, black)
				t.setColor(parentOfX, red)
				t.rotateLeft(parentOfX)
				w = t.right(parentOfX)
			}
			if t.color(t.left(w)) == black && t.color(t.right(w)) == black {
				t.setColor(w, red)
				x = parentOfX
				parentOfX = t.parent(x)
			} else {
				if t.color(t.right(w)) == black {
					t.setColor(t.left(w), black)
					t.setColor(w, red)
					t.rotateRight(w)
					w = t.right(parentOfX)
				}
				t.setColor(w, t.color(parentOfX))
				t.setColor(parentOfX, black)
				t.setColor(t.right(w), black)
				t.rotateLeft(parentOfX)
				x = t.root
				parentOfX = nil
			}
		} else {
			w := t.left(parentOfX)
			if t.color(w) == red {
				t.setColor(w, black)
				t.setColor(parentOfX, red)
				t.rotateRight(parentOfX)
				w = t.left(parentOfX)
			}
			if t.color(t.right(w)) == black && t.color(t.left(w)) == black {
				t.setColor(w, red)
				x = parentOfX
				parentOfX = t.parent(x)
			} else {
				if t.color(t.left(w)) == black {
					t.setColor(t.right(w), black)
					t.setColor(w, red)
					t.rotateLeft(w)
					w = t.left(parentOfX)
				}
				t.setColor(w, t.color(parentOfX))
				t.setColor(parentOfX, black)
				t.setColor(t.left(w), black)
				t.rotateRight(parentOfX)
				x = t.root
				parentOfX = nil
			}
		}
	}
	if x != nil {
		t.setColor(x, black)
	}
}

func (t *RbTree[T]) transplant(u, v *T) {
	if t.parent(u) == nil {
		t.root = v
	} else if u == t.left(t.parent(u)) {
		t.setLeft(t.parent(u), v)
	} else {
		t.setRight(t.parent(u), v)
	}
	if v != nil {
		t.setParent(v, t.parent(u))
	}
}

// Empty returns true if tree is empty
func (t RbTree[T]) Empty() bool {
	return t.size == 0
}

// Size returns the number of elements in the tree
func (t RbTree[T]) Size() int {
	return t.size
}

// Len returns the number of elements in the tree
func (t RbTree[T]) Len() int {
	return t.size
}

// Swap exchanges contents with another tree
func (t *RbTree[T]) Swap(other *RbTree[T]) {
	other.hookFunc, t.hookFunc = t.hookFunc, other.hookFunc
	other.lessFunc, t.lessFunc = t.lessFunc, other.lessFunc
	t.root, other.root = other.root, t.root
	t.first, other.first = other.first, t.first
	t.last, other.last = other.last, t.last
	t.size, other.size = other.size, t.size
}

// Front returns the first (leftmost) node in the tree
func (t RbTree[T]) Front() *T {
	return t.first
}

// Back returns the last (rightmost) node in the tree
func (t RbTree[T]) Back() *T {
	return t.last
}

// Clear removes all nodes from the tree
func (t *RbTree[T]) Clear() []*T {
	nodes := make([]*T, 0, t.size)

	t.TraversePostOrder(func(node *T) {
		nodes = append(nodes, node)
		t.getHook(node).Init()
	})

	t.Init()
	return nodes
}

// Traverse traverses tree in-order
func (t RbTree[T]) Traverse(f func(*T)) {
	var traverse func(*T)
	traverse = func(node *T) {
		if node == nil {
			return
		}
		hook := t.getHook(node)
		traverse(hook.left)
		f(node)
		traverse(hook.right)
	}
	traverse(t.root)
}

// TraversePreOrder traverses tree in pre-order
func (t RbTree[T]) TraversePreOrder(f func(*T)) {
	var traverse func(*T)
	traverse = func(node *T) {
		if node == nil {
			return
		}
		hook := t.getHook(node)
		f(node)
		traverse(hook.left)
		traverse(hook.right)
	}
	traverse(t.root)
}

// TraversePostOrder traverses tree in post-order
func (t RbTree[T]) TraversePostOrder(f func(*T)) {
	var traverse func(*T)
	traverse = func(node *T) {
		if node == nil {
			return
		}
		hook := t.getHook(node)
		traverse(hook.left)
		traverse(hook.right)
		f(node)
	}
	traverse(t.root)
}

// Insert adds a new node to the tree
func (t *RbTree[T]) Insert(item *T) bool {
	t.verifyElementNotLinked(item)
	defer t.verify()

	var y *T
	x := t.root
	for x != nil {
		y = x
		if t.lessFunc(item, x) {
			x = t.left(x)
		} else if t.lessFunc(x, item) {
			x = t.right(x)
		} else {
			return false
		}
	}
	t.setParent(item, y)
	if y == nil {
		t.root = item
		t.first = item
		t.last = item
	} else if t.lessFunc(item, y) {
		t.setLeft(y, item)
		if y == t.first {
			t.first = item
		}
	} else {
		t.setRight(y, item)
		if y == t.last {
			t.last = item
		}
	}
	t.setLeft(item, nil)
	t.setRight(item, nil)
	t.setColor(item, red)
	t.insertFixup(item)
	t.size++
	return true
}

// Erase removes a node from the tree
func (t *RbTree[T]) Erase(item *T) bool {
	if item == nil {
		return false
	}
	t.verifyNotEmpty()
	t.verifyIsMemberOfCurrent(item)
	defer t.verifyElementNotLinked(item)
	defer t.verify()

	originalColor := t.color(item)
	var x *T
	var xParent *T
	if t.left(item) == nil {
		x = t.right(item)
		xParent = t.parent(item)
		t.transplant(item, t.right(item))
	} else if t.right(item) == nil {
		x = t.left(item)
		xParent = t.parent(item)
		t.transplant(item, t.left(item))
	} else {
		y := t.min(t.right(item))
		originalColor = t.color(y)
		x = t.right(y)
		if t.parent(y) == item {
			xParent = y
		} else {
			xParent = t.parent(y)
			t.transplant(y, t.right(y))
			t.setRight(y, t.right(item))
			t.setParent(t.right(y), y)
		}
		t.transplant(item, y)
		t.setLeft(y, t.left(item))
		t.setParent(t.left(y), y)
		t.setColor(y, t.color(item))
	}
	if originalColor == black {
		if x == nil {
			t.deleteFixup(x, xParent)
		} else {
			t.deleteFixup(x, t.parent(x))
		}
	}
	if item == t.first {
		t.first = t.next(item)
	}
	if item == t.last {
		t.last = t.prev(item)
	}
	t.setLeft(item, nil)
	t.setRight(item, nil)
	t.setParent(item, nil)
	t.size--
	return true
}

// Merge combines two trees
func (t *RbTree[T]) Merge(other *RbTree[T]) {
	if other.size == 0 {
		return
	}
	defer t.verify()
	defer other.verify()

	node := other.Front()
	for node != nil {
		next := other.Next(node)

		if t.Find(node) == nil {
			other.Erase(node)
			t.Insert(node)
		}
		node = next
	}
}

// Contains checks if element that compares equal with item exists in tree
func (t RbTree[T]) Contains(item *T) bool {
	return t.Find(item) != nil
}

// Find searches for an element that compares equal with item
func (t RbTree[T]) Find(item *T) *T {
	current := t.root
	for current != nil {
		if t.lessFunc(item, current) {
			current = t.left(current)
		} else if t.lessFunc(current, item) {
			current = t.right(current)
		} else {
			return current
		}
	}
	return nil
}

// LowerBound finds first element not less than item
func (t RbTree[T]) LowerBound(item *T) *T {
	var candidate *T
	current := t.root
	for current != nil {
		if !t.lessFunc(current, item) {
			candidate = current
			current = t.left(current)
		} else {
			current = t.right(current)
		}
	}
	return candidate
}

// UpperBound finds first element greater than item
func (t RbTree[T]) UpperBound(item *T) *T {
	var candidate *T
	current := t.root
	for current != nil {
		if t.lessFunc(item, current) {
			candidate = current
			current = t.left(current)
		} else {
			current = t.right(current)
		}
	}
	return candidate
}

// EraseIf removes nodes matching predicate
func (t *RbTree[T]) EraseIf(predicate func(*T) bool) (erased []*T) {
	defer t.verify()
	defer func() {
		for _, n := range erased {
			t.verifyElementNotLinked(n)
		}
	}()

	erased = make([]*T, 0)
	var toErase []*T

	// First collect nodes to erase
	for node := t.Front(); node != nil; node = t.Next(node) {
		if predicate(node) {
			toErase = append(toErase, node)
		}
	}

	// Erase collected nodes
	for _, node := range toErase {
		if t.Erase(node) {
			erased = append(erased, node)
		}
	}

	return erased
}

// Includes checks if tree contains all elements of another tree
func (t RbTree[T]) Includes(other *RbTree[T]) bool {
	if other.size == 0 {
		return true
	}
	a := t.Front()
	b := other.Front()
	for a != nil && b != nil {
		if t.lessFunc(a, b) {
			a = t.Next(a)
		} else if t.lessFunc(b, a) {
			return false
		} else {
			a = t.Next(a)
			b = other.Next(b)
		}
	}
	return b == nil
}

// Difference returns elements in tree but not in other
func (t RbTree[T]) Difference(other *RbTree[T]) []*T {
	var result []*T
	a := t.Front()
	b := other.Front()
	for a != nil && b != nil {
		if t.lessFunc(a, b) {
			result = append(result, a)
			a = t.Next(a)
		} else if t.lessFunc(b, a) {
			b = other.Next(b)
		} else {
			a = t.Next(a)
			b = other.Next(b)
		}
	}
	for a != nil {
		result = append(result, a)
		a = t.Next(a)
	}
	return result
}

// Intersection returns elements common to both trees
func (t RbTree[T]) Intersection(other *RbTree[T]) []*T {
	var result []*T
	a := t.Front()
	b := other.Front()
	for a != nil && b != nil {
		if t.lessFunc(a, b) {
			a = t.Next(a)
		} else if t.lessFunc(b, a) {
			b = other.Next(b)
		} else {
			result = append(result, a)
			a = t.Next(a)
			b = other.Next(b)
		}
	}
	return result
}

// SymDifference returns elements not common to both trees
func (t RbTree[T]) SymDifference(other *RbTree[T]) []*T {
	var result []*T
	a := t.Front()
	b := other.Front()
	for a != nil && b != nil {
		if t.lessFunc(a, b) {
			result = append(result, a)
			a = t.Next(a)
		} else if t.lessFunc(b, a) {
			result = append(result, b)
			b = other.Next(b)
		} else {
			a = t.Next(a)
			b = other.Next(b)
		}
	}
	for a != nil {
		result = append(result, a)
		a = t.Next(a)
	}
	for b != nil {
		result = append(result, b)
		b = other.Next(b)
	}
	return result
}

// Union returns all elements from both trees
func (t RbTree[T]) Union(other *RbTree[T]) []*T {
	var result []*T
	a := t.Front()
	b := other.Front()
	for a != nil && b != nil {
		if t.lessFunc(a, b) {
			result = append(result, a)
			a = t.Next(a)
		} else if t.lessFunc(b, a) {
			result = append(result, b)
			b = other.Next(b)
		} else {
			result = append(result, a)
			a = t.Next(a)
			b = other.Next(b)
		}
	}
	for a != nil {
		result = append(result, a)
		a = t.Next(a)
	}
	for b != nil {
		result = append(result, b)
		b = other.Next(b)
	}
	return result
}
