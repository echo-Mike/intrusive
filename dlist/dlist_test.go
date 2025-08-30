package dlist

import (
	"testing"
)

type fuzzEmbedItem struct {
	Hook[fuzzEmbedItem]
	value  int
	isUsed bool
	id     int // Unique identifier for better debugging
}

func fuzzEmbedHook(self *fuzzEmbedItem) *Hook[fuzzEmbedItem] {
	return &self.Hook
}

func newFuzzList() DList[fuzzEmbedItem] {
	return New(fuzzEmbedHook)
}

func newFuzz(value, id int) fuzzEmbedItem {
	return fuzzEmbedItem{
		Hook:   NewHook[fuzzEmbedItem](),
		value:  value,
		isUsed: false,
		id:     id,
	}
}

func lessFuzz(lhs, rhs *fuzzEmbedItem) bool {
	return lhs.value < rhs.value
}

const (
	opInsert byte = iota
	opErase
	opPushFront
	opPopFront
	opPushBack
	opPopBack
	opSplice
	opSpliceFront
	opSpliceBack
	opClear
	opReverse
	opMerge
	opSort
	opUnique
	opRemoveIf
	opVerifyList
	opInsertNil
	opCOUNT
)

func elementAt(d DList[fuzzEmbedItem], pos int) (r *fuzzEmbedItem) {
	r = d.Front()
	for i := 0; i < pos && r != nil; i++ {
		r = r.Next()
	}
	return
}

func verifyListConsistency(t *testing.T, l *DList[fuzzEmbedItem]) {
	if l.Empty() {
		if l.Front() != nil || l.Back() != nil || l.Size() != 0 {
			t.Errorf("Empty list inconsistency: front=%v, back=%v, size=%d", l.Front(), l.Back(), l.Size())
		}
		return
	}

	// Verify forward traversal
	count := 0
	current := l.Front()
	var prev *fuzzEmbedItem

	for current != nil {
		count++
		hook := l.hookFunc(current)

		// Check prev pointer consistency
		if hook.prev != prev {
			t.Errorf("Prev pointer inconsistency at element %d", current.id)
		}

		prev = current
		current = hook.next
	}

	// Verify backward traversal
	countBack := 0
	current = l.Back()
	var next *fuzzEmbedItem

	for current != nil {
		countBack++
		hook := l.hookFunc(current)

		// Check next pointer consistency
		if hook.next != next {
			t.Errorf("Next pointer inconsistency at element %d", current.id)
		}

		next = current
		current = hook.prev
	}

	// Verify size matches traversal count
	if count != l.Size() || countBack != l.Size() {
		t.Errorf("Size inconsistency: forward=%d, backward=%d, stored=%d", count, countBack, l.Size())
	}

	// Verify front and back pointers
	if l.hookFunc(l.Front()).prev != nil {
		t.Errorf("Front element should not have a previous pointer")
	}

	if l.hookFunc(l.Back()).next != nil {
		t.Errorf("Back element should not have a next pointer")
	}
}

func nextState(t *testing.T, items []fuzzEmbedItem, lists []DList[fuzzEmbedItem]) func(op, arg1, arg2, arg3 byte) {
	return func(op, arg1, arg2, arg3 byte) {
		listIdx := int(arg1) % len(lists)
		itemIdx := int(arg2) % len(items)
		positionIdx := int(arg3)

		l := &lists[listIdx]
		i := &items[itemIdx]

		switch op % opCOUNT {
		case opInsert:
			if !i.isUsed {
				if pos := elementAt(*l, positionIdx); pos != nil {
					l.Insert(pos, i)
					i.isUsed = true
				} else {
					l.Insert(nil, i)
					i.isUsed = true
				}
			}

		case opInsertNil:
			if !i.isUsed {
				l.Insert(nil, i)
				i.isUsed = true
			}

		case opErase:
			if el := elementAt(*l, positionIdx); el != nil {
				l.Erase(el)
				el.isUsed = false
			}

		case opPushFront:
			if !i.isUsed {
				l.PushFront(i)
				i.isUsed = true
			}

		case opPopFront:
			if el := l.PopFront(); el != nil {
				el.isUsed = false
			}

		case opPushBack:
			if !i.isUsed {
				l.PushBack(i)
				i.isUsed = true
			}

		case opPopBack:
			if el := l.PopBack(); el != nil {
				el.isUsed = false
			}

		case opSplice:
			if listIdx != int(arg2)%len(lists) {
				l2 := &lists[int(arg2)%len(lists)]
				if pos := elementAt(*l, positionIdx); pos != nil {
					l.Splice(pos, l2)
				} else {
					l.Splice(nil, l2)
				}
			}

		case opSpliceFront:
			if listIdx != int(arg2)%len(lists) {
				l2 := &lists[int(arg2)%len(lists)]
				l.SpliceFront(l2)
			}

		case opSpliceBack:
			if listIdx != int(arg2)%len(lists) {
				l2 := &lists[int(arg2)%len(lists)]
				l.SpliceBack(l2)
			}

		case opClear:
			if elements := l.Clear(); len(elements) > 0 {
				for _, e := range elements {
					e.isUsed = false
				}
			}

		case opReverse:
			l.Reverse()

		case opMerge:
			if listIdx != int(arg2)%len(lists) {
				l2 := &lists[int(arg2)%len(lists)]
				l.Sort(lessFuzz)
				l2.Sort(lessFuzz)
				l.Merge(l2, lessFuzz)
			}

		case opSort:
			l.Sort(lessFuzz)

		case opUnique:
			if int(arg2)%2 == 0 {
				l.Sort(lessFuzz)
			}
			if elements := l.Unique(lessFuzz); len(elements) > 0 {
				for _, e := range elements {
					e.isUsed = false
				}
			}

		case opRemoveIf:
			var predicate func(*fuzzEmbedItem) bool
			switch int(arg2) % 4 {
			case 0:
				predicate = func(e *fuzzEmbedItem) bool { return e.value%2 == 0 }
			case 1:
				predicate = func(e *fuzzEmbedItem) bool { return e.value%2 == 1 }
			case 2:
				predicate = func(e *fuzzEmbedItem) bool { return e.value < 8 }
			case 3:
				predicate = func(e *fuzzEmbedItem) bool { return e.value >= 8 }
			}

			if elements := l.RemoveIf(predicate); len(elements) > 0 {
				for _, e := range elements {
					e.isUsed = false
				}
			}

		case opVerifyList:
			verifyListConsistency(t, l)
		}
	}
}

func FuzzDListOps(f *testing.F) {
	// Use more items and lists for better coverage
	const numItems = 512
	const numLists = 8

	items := make([]fuzzEmbedItem, numItems)
	for i := range items {
		items[i] = newFuzz(i%32, i) // More varied values
	}

	lists := make([]DList[fuzzEmbedItem], numLists)
	for i := range lists {
		lists[i] = newFuzzList()
	}

	f.Fuzz(func(t *testing.T, commands []byte) {
		// Reset all lists and items before each fuzz run
		for i := range lists {
			if elements := lists[i].Clear(); len(elements) > 0 {
				for _, e := range elements {
					e.isUsed = false
				}
			}
		}

		for i := range items {
			items[i].isUsed = false
			items[i].Hook.Init()
		}

		next := nextState(t, items, lists)

		// Process commands in chunks of 4 bytes
		for i := 0; i+3 < len(commands); i += 4 {
			next(commands[i], commands[i+1], commands[i+2], commands[i+3])
		}

		// Verify all lists at the end of the sequence
		for i := range lists {
			verifyListConsistency(t, &lists[i])
		}
	})
}
