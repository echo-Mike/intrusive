//go:build debug

package slist

import (
	"fmt"
)

func (s *SList[T]) verifyNotEmpty() {
	if s.first == nil || s.last == nil || s.size == 0 {
		panic(fmt.Sprintf("unexpected empty list: SList %p", s))
	}
}

func (s *SList[T]) verifyElementNotLinked(element *T) {
	if s.hookFunc(element).next != nil {
		panic(fmt.Sprintf("already linked element detected: SList %p element: %p", s, element))
	}
}

func (s *SList[T]) verifyIsMemberOfCurrent(element *T) {
	for e := s.first; e != nil; e = s.hookFunc(e).next {
		if e == element {
			return
		}
	}
	panic(fmt.Sprintf("not member of detected: SList %p element: %p", s, element))
}

func (s *SList[T]) verifySize() {
	e := s.first
	for i := 0; i < s.size; i++ {
		if e == nil {
			panic(fmt.Sprintf("size of list is less than expected: SList %p", s))
		}
		e = s.hookFunc(e).next
	}
	if e != nil {
		panic(fmt.Sprintf("size of list is greater than expected: SList %p", s))
	}
}

func (s *SList[T]) verifyNoCycle() {
	walked := make(map[*T]bool)
	for e := s.first; e != nil; e = s.hookFunc(e).next {
		if walked[e] {
			panic(fmt.Sprintf("found a cycle: SList %p", s))
		}
		walked[e] = true
	}
}
