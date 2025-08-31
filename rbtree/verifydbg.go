//go:build debug

package rbtree

import (
	"fmt"
)

func (t *RbTree[T]) verifyNotEmpty() {
	if t.root == nil || t.first == nil || t.last == nil || t.size == 0 {
		panic(fmt.Sprintf("unexpected empty tree: RbTree %p", t))
	}
}

func (t *RbTree[T]) verifyElementNotLinked(element *T) {
	hook := t.getHook(element)
	if hook.left != nil || hook.right != nil || hook.parent != nil {
		panic(fmt.Sprintf("already linked element detected: RbTree %p element: %p", t, element))
	}
}

func (t *RbTree[T]) verifyIsMemberOfCurrent(element *T) {
	current := t.root
	for current != nil {
		if current == element {
			return
		}
		if t.lessFunc(element, current) {
			current = t.left(current)
		} else {
			current = t.right(current)
		}
	}
	panic(fmt.Sprintf("not member of detected: RbTree %p element: %p", t, element))
}

func (t *RbTree[T]) verifySize() {
	count := 0
	var traverse func(*T)
	traverse = func(node *T) {
		if node == nil {
			return
		}
		count++
		traverse(t.left(node))
		traverse(t.right(node))
	}
	traverse(t.root)

	if count != t.size {
		panic(fmt.Sprintf("size mismatch: expected %d, got %d: RbTree %p", t.size, count, t))
	}
}

func (t *RbTree[T]) verifyBSTProperty(node *T, min, max *T) {
	if node == nil {
		return
	}

	if min != nil && t.lessFunc(node, min) {
		panic(fmt.Sprintf("BST property violation: node %p < min %p: RbTree %p", node, min, t))
	}

	if max != nil && t.lessFunc(max, node) {
		panic(fmt.Sprintf("BST property violation: node %p > max %p: RbTree %p", node, max, t))
	}

	t.verifyBSTProperty(t.left(node), min, node)
	t.verifyBSTProperty(t.right(node), node, max)
}

func (t *RbTree[T]) verifyRedBlackProperties(node *T) int {
	if node == nil {
		return 1
	}

	leftBlackHeight := t.verifyRedBlackProperties(t.left(node))
	rightBlackHeight := t.verifyRedBlackProperties(t.right(node))

	if leftBlackHeight != rightBlackHeight {
		panic(fmt.Sprintf("black height mismatch: left %d, right %d: RbTree %p node: %p",
			leftBlackHeight, rightBlackHeight, t, node))
	}

	if t.color(node) == red {
		if t.color(t.left(node)) == red || t.color(t.right(node)) == red {
			panic(fmt.Sprintf("red node with red child: RbTree %p node: %p", t, node))
		}
		return leftBlackHeight
	}

	return leftBlackHeight + 1
}

func (t *RbTree[T]) verifyParentPointers(node *T) {
	if node == nil {
		return
	}

	if t.left(node) != nil && t.parent(t.left(node)) != node {
		panic(fmt.Sprintf("left child parent pointer mismatch: RbTree %p node: %p", t, node))
	}

	if t.right(node) != nil && t.parent(t.right(node)) != node {
		panic(fmt.Sprintf("right child parent pointer mismatch: RbTree %p node: %p", t, node))
	}

	t.verifyParentPointers(t.left(node))
	t.verifyParentPointers(t.right(node))
}

func (t *RbTree[T]) verifyFirstLast() {
	if t.size == 0 {
		if t.first != nil || t.last != nil {
			panic(fmt.Sprintf("non-nil first/last in empty tree: RbTree %p", t))
		}
		return
	}

	if t.min(t.root) != t.first {
		panic(fmt.Sprintf("first pointer mismatch: expected %p, got %p: RbTree %p",
			t.min(t.root), t.first, t))
	}

	if t.max(t.root) != t.last {
		panic(fmt.Sprintf("last pointer mismatch: expected %p, got %p: RbTree %p",
			t.max(t.root), t.last, t))
	}
}

func (t *RbTree[T]) verifyNoCycle() {
	visited := make(map[*T]bool)
	var traverse func(*T)
	traverse = func(node *T) {
		if node == nil {
			return
		}
		if visited[node] {
			panic(fmt.Sprintf("cycle detected: RbTree %p node: %p", t, node))
		}
		visited[node] = true
		traverse(t.left(node))
		traverse(t.right(node))
	}
	traverse(t.root)
}

func (t *RbTree[T]) verify() {
	t.verifySize()
	t.verifyFirstLast()
	t.verifyNoCycle()
	if t.root != nil {
		t.verifyBSTProperty(t.root, nil, nil)
		t.verifyRedBlackProperties(t.root)
		t.verifyParentPointers(t.root)
		if t.color(t.root) != black {
			panic(fmt.Sprintf("root is not black: RbTree %p", t))
		}
	}
}
