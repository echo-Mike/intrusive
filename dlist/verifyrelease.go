//go:build !debug

package dlist

func (d *DList[T]) verifyNotEmpty() {
}

func (d *DList[T]) verifyElementNotLinked(element *T) {
}

func (d *DList[T]) verifyIsMemberOfCurrent(element *T) {
}

func (d *DList[T]) verifySize() {
}

func (d *DList[T]) verifyNoCycle() {
}
