//go:build !debug

package slist

func (s *SList[Ops, T]) verifyNotEmpty() {
}

func (s *SList[Ops, T]) verifyElementNotLinked(element *T) {
}

func (s *SList[Ops, T]) verifyIsMemberOfCurrent(element *T) {
}

func (s *SList[Ops, T]) verifySize() {
}

func (s *SList[Ops, T]) verifyNoCycle() {
}
