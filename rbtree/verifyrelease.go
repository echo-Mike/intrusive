//go:build !debug

package rbtree

func (t *RbTree[T]) verifyNotEmpty() {
}

func (t *RbTree[T]) verifyElementNotLinked(element *T) {
}

func (t *RbTree[T]) verifyIsMemberOfCurrent(element *T) {
}

func (t *RbTree[T]) verify() {
}
