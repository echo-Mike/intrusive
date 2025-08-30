package slist

import "testing"

type fuzzEmbedItem struct {
	Hook[fuzzEmbedItem]
	value  int
	isUsed bool
}

func fuzzEmbedHook(self *fuzzEmbedItem) *Hook[fuzzEmbedItem] {
	return &self.Hook
}

func newFuzzList() SList[fuzzEmbedItem] {
	return New(fuzzEmbedHook)
}

func newFuzz(value int) fuzzEmbedItem {
	return fuzzEmbedItem{Hook: NewHook[fuzzEmbedItem](), value: value, isUsed: false}
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
	opCOUNT
)

func elementAt(l SList[fuzzEmbedItem], pos int) (r *fuzzEmbedItem) {
	r = l.Front()
	for range pos {
		if r != nil {
			r = r.Next()
		}
	}
	return
}

func nextState(items []fuzzEmbedItem, lists []SList[fuzzEmbedItem]) func(op, arg1, arg2, arg3 byte) {
	return func(op, arg1, arg2, arg3 byte) {
		switch op % opCOUNT {
		case opInsertAfter:
			{
				l := &lists[int(arg1)%len(lists)]
				i := &items[int(arg2)%len(items)]
				if el := elementAt(*l, int(arg3)); el != nil && !i.isUsed {
					l.InsertAfter(el, i)
					i.isUsed = true
				}
			}
		case opRemoveAfter:
			{
				l := &lists[int(arg1)%len(lists)]
				if el := elementAt(*l, int(arg2)); el != nil {
					i := l.RemoveAfter(el)
					if i != nil {
						i.isUsed = false
					}
				}
			}
		case opSpliceAfter:
			{
				if int(arg1)%len(lists) != int(arg2)%len(lists) {
					l1 := &lists[int(arg1)%len(lists)]
					l2 := &lists[int(arg2)%len(lists)]
					if el := elementAt(*l1, int(arg3)); el != nil {
						l1.SpliceAfter(el, l2)
					}
				}
			}
		case opPushFront:
			{
				l := &lists[int(arg1)%len(lists)]
				i := &items[int(arg2)%len(items)]
				if !i.isUsed {
					l.PushFront(i)
					i.isUsed = true
				}
			}
		case opPopFront:
			{
				l := &lists[int(arg1)%len(lists)]
				if el := l.PopFront(); el != nil {
					el.isUsed = false
				}
			}
		case opSpliceFront:
			{
				if int(arg1)%len(lists) != int(arg2)%len(lists) {
					l1 := &lists[int(arg1)%len(lists)]
					l2 := &lists[int(arg2)%len(lists)]
					l1.SpliceFront(l2)
				}
			}
		case opPushBack:
			{
				l := &lists[int(arg1)%len(lists)]
				i := &items[int(arg2)%len(items)]
				if !i.isUsed {
					l.PushBack(i)
					i.isUsed = true
				}
			}
		case opSpliceBack:
			{
				if int(arg1)%len(lists) != int(arg2)%len(lists) {
					l1 := &lists[int(arg1)%len(lists)]
					l2 := &lists[int(arg2)%len(lists)]
					l1.SpliceBack(l2)
				}
			}
		case opClear:
			{
				l := &lists[int(arg1)%len(lists)]
				if el := l.Clear(); 0 < len(el) {
					for _, e := range el {
						e.isUsed = false
					}
				}
			}
		case opReverse:
			{
				l := &lists[int(arg1)%len(lists)]
				l.Reverse()
			}
		case opMerge:
			{
				if int(arg1)%len(lists) != int(arg2)%len(lists) {
					l1 := &lists[int(arg1)%len(lists)]
					l2 := &lists[int(arg2)%len(lists)]
					l1.Sort(lessFuzz)
					l2.Sort(lessFuzz)
					l1.Merge(l2, lessFuzz)
				}
			}
		case opSort:
			{
				l := &lists[int(arg1)%len(lists)]
				l.Sort(lessFuzz)
			}
		case opUnique:
			{
				l := &lists[int(arg1)%len(lists)]
				if int(arg2)%2 == 0 {
					l.Sort(lessFuzz)
				}
				if el := l.Unique(lessFuzz); 0 < len(el) {
					for _, e := range el {
						e.isUsed = false
					}
				}
			}
		case opRemoveIf:
			{
				l := &lists[int(arg1)%len(lists)]
				var f func(*fuzzEmbedItem) bool
				if int(arg2)%2 == 0 {
					f = func(e *fuzzEmbedItem) bool { return e.value%2 == 0 }
				} else {
					f = func(e *fuzzEmbedItem) bool { return e.value%2 == 1 }
				}
				if el := l.RemoveIf(f); 0 < len(el) {
					for _, e := range el {
						e.isUsed = false
					}
				}
			}
		}
	}
}

func FuzzSListOps(f *testing.F) {
	items := make([]fuzzEmbedItem, 0, 256)
	for i := range cap(items) {
		items = append(items, newFuzz(i%16))
	}
	lists := make([]SList[fuzzEmbedItem], 0, 5)
	for range cap(lists) {
		lists = append(lists, newFuzzList())
	}
	next := nextState(items, lists)
	f.Fuzz(func(t *testing.T, commands []byte) {
		for i := 0; i < len(commands)/4; i++ {
			next(commands[0+4*i], commands[1+4*i], commands[2+4*i], commands[3+4*i])
		}
	})
}
