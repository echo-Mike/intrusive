//go:build debug

package dlist

import (
	"fmt"
)

func (d *DList[T]) verifyNotEmpty() {
	if d.first == nil || d.last == nil || d.size == 0 {
		panic(fmt.Sprintf("unexpected empty list: DList %p", d))
	}
}

func (d *DList[T]) verifyElementNotLinked(element *T) {
	hook := d.hookFunc(element)
	if hook.next != nil || hook.prev != nil {
		panic(fmt.Sprintf("already linked element detected: DList %p element: %p", d, element))
	}
}

func (d *DList[T]) verifyIsMemberOfCurrent(element *T) {
	for e := d.first; e != nil; e = d.hookFunc(e).next {
		if e == element {
			return
		}
	}
	panic(fmt.Sprintf("not member of detected: DList %p element: %p", d, element))
}

func (d *DList[T]) verifySize() {
	count := 0
	for e := d.first; e != nil; e = d.hookFunc(e).next {
		count++
		if count > d.size {
			panic(fmt.Sprintf("size of list is greater than expected: DList %p", d))
		}
	}
	if count != d.size {
		panic(fmt.Sprintf("size mismatch: expected %d, got %d: DList %p", d.size, count, d))
	}
}

func (d *DList[T]) verifyNoCycle() {
	visited := make(map[*T]bool)
	for e := d.first; e != nil; e = d.hookFunc(e).next {
		if visited[e] {
			panic(fmt.Sprintf("found a cycle: DList %p", d))
		}
		visited[e] = true
	}

	visited = make(map[*T]bool)
	for e := d.last; e != nil; e = d.hookFunc(e).prev {
		if visited[e] {
			panic(fmt.Sprintf("found a cycle: DList %p", d))
		}
		visited[e] = true
	}
}