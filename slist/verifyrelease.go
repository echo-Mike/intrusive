//go:build !debug

package slist

func (s *SList[T]) verifyNotEmpty() {
}

func (s *SList[T]) verifyElementNotLinked(element *T) {
}

func (s *SList[T]) verifyIsMemberOfCurrent(element *T) {
}

func (s *SList[T]) verifySize() {
}

func (s *SList[T]) verifyNoCycle() {
}
