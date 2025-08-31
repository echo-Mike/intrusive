package slist

import "testing"

type fuzzEmbedItem struct {
	Hook[fuzzEmbedItem]
	value  int
	isUsed bool
	id     int
}

func fuzzEmbedHook(self *fuzzEmbedItem) *Hook[fuzzEmbedItem] {
	return &self.Hook
}

func newFuzzList() SList[fuzzEmbedItem] {
	return New(fuzzEmbedHook)
}

func newFuzz(value, id int) fuzzEmbedItem {
	return fuzzEmbedItem{Hook: NewHook[fuzzEmbedItem](), value: value, isUsed: false, id: id}
}

func lessFuzz(lhs, rhs *fuzzEmbedItem) bool {
	return lhs.value < rhs.value
}

const (
	opInsertAfter byte = iota
	opRemoveAfter
	opSpliceAfter
	opPushFront
	opPopFront
	opSpliceFront
	opPushBack
	opSpliceBack
	opClear
	opReverse
	opMerge
	opSort
	opUnique
	opRemoveIf
	opVerifyList
	opCOUNT
)

func elementAt(l SList[fuzzEmbedItem], pos int) (r *fuzzEmbedItem) {
	r = l.Front()
	for i := 0; i < pos && r != nil; i++ {
		r = r.Next()
	}
	return
}

func verifyListConsistency(t *testing.T, l *SList[fuzzEmbedItem]) {
	if l.Empty() {
		if l.Front() != nil || l.Back() != nil || l.Size() != 0 {
			t.Errorf("Empty list inconsistency: front=%v, back=%v, size=%d", l.Front(), l.Back(), l.Size())
		}
		return
	}

	// Verify forward traversal
	count := 0
	current := l.Front()
	var last *fuzzEmbedItem

	for current != nil {
		count++
		last = current
		current = l.hookFunc(current).next
	}

	// Verify size matches traversal count
	if count != l.Size() {
		t.Errorf("Size inconsistency: forward=%d, stored=%d", count, l.Size())
	}

	// Verify back pointer points to last element
	if l.Back() != last {
		t.Errorf("Back pointer inconsistency: expected=%v, actual=%v", last, l.Back())
	}

	// Verify last element's next is nil
	if l.hookFunc(l.Back()).next != nil {
		t.Errorf("Last element should not have a next pointer")
	}
}

func nextState(t *testing.T, items []fuzzEmbedItem, lists []SList[fuzzEmbedItem]) func(op, arg1, arg2, arg3 byte) {
	return func(op, arg1, arg2, arg3 byte) {
		listIdx := int(arg1) % len(lists)
		itemIdx := int(arg2) % len(items)
		positionIdx := int(arg3)

		l := &lists[listIdx]
		i := &items[itemIdx]

		switch op % opCOUNT {
		case opInsertAfter:
			if !i.isUsed {
				if pos := elementAt(*l, positionIdx); pos != nil {
					l.InsertAfter(pos, i)
					i.isUsed = true
				}
			}

		case opRemoveAfter:
			if pos := elementAt(*l, positionIdx); pos != nil {
				if removed := l.RemoveAfter(pos); removed != nil {
					removed.isUsed = false
				}
			}

		case opSpliceAfter:
			if listIdx != int(arg2)%len(lists) {
				l2 := &lists[int(arg2)%len(lists)]
				if pos := elementAt(*l, positionIdx); pos != nil {
					l.SpliceAfter(pos, l2)
				}
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

		case opSpliceFront:
			if listIdx != int(arg2)%len(lists) {
				l2 := &lists[int(arg2)%len(lists)]
				l.SpliceFront(l2)
			}

		case opPushBack:
			if !i.isUsed {
				l.PushBack(i)
				i.isUsed = true
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

func FuzzSListOps(f *testing.F) {
	const numItems = 256
	const numLists = 8

	items := make([]fuzzEmbedItem, numItems)
	for i := range items {
		items[i] = newFuzz(i%32, i)
	}

	lists := make([]SList[fuzzEmbedItem], numLists)
	for i := range lists {
		lists[i] = newFuzzList()
	}

	f.Fuzz(func(t *testing.T, commands []byte) {
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

		for i := 0; i+3 < len(commands); i += 4 {
			next(commands[i], commands[i+1], commands[i+2], commands[i+3])
		}

		for i := range lists {
			verifyListConsistency(t, &lists[i])
		}
	})
}
