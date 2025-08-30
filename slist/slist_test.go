package slist

import (
	"testing"

	"github.com/echo-Mike/intrusive/internal/pkg/fn"
)

/// Possible use cases: 1. embedded hook 2. member hook

type testEmbedItem struct {
	Hook[testEmbedItem]
	value int
}

func embedHook(self *testEmbedItem) *Hook[testEmbedItem] {
	return &self.Hook
}

func newEmbedList() SList[testEmbedItem] {
	return New(embedHook)
}

func newEmbed(value int) testEmbedItem {
	return testEmbedItem{Hook: NewHook[testEmbedItem](), value: value}
}

func newEmbedListGenerate(count int, generator func(position int) int) (l SList[testEmbedItem]) {
	l = newEmbedList()
	for i := 0; i < count; i++ {
		item := newEmbed(generator(i))
		l.PushBack(&item)
	}
	return
}

func lessEmbed(lhs, rhs *testEmbedItem) bool {
	return lhs.value < rhs.value
}

func nextEmbed(l SList[testEmbedItem]) func() int {
	first := l.Front()
	return func() int {
		defer func() { first = first.Next() }()
		return first.value
	}
}

type testMemberItem struct {
	hook  Hook[testMemberItem]
	value int
}

func memberHook(self *testMemberItem) *Hook[testMemberItem] {
	return &self.hook
}

func newMemberList() SList[testMemberItem] {
	return New(memberHook)
}

func newMember(value int) testMemberItem {
	return testMemberItem{hook: NewHook[testMemberItem](), value: value}
}

func newMemberListGenerate(count int, generator func(position int) int) (l SList[testMemberItem]) {
	l = newMemberList()
	for i := 0; i < count; i++ {
		item := newMember(generator(i))
		l.PushBack(&item)
	}
	return
}

func lessMember(lhs, rhs *testMemberItem) bool {
	return lhs.value < rhs.value
}

func nextMember(l SList[testMemberItem]) func() int {
	first := l.Front()
	return func() int {
		defer func() { first = first.hook.Next() }()
		return first.value
	}
}

func increment(n int) func(int) int {
	return func(p int) int {
		return n + p
	}
}

func decrement(n int) func(int) int {
	return func(p int) int {
		return n - p
	}
}

func static(n int) func(int) int {
	return func(int) int {
		return n
	}
}

func isSorted(g func() int, n int) bool {
	type T struct {
		int
		bool
	}
	return fn.Reduce(g, func(a int, b T) T {
		b.bool = b.bool && b.int <= a
		b.int = a
		return b
	}, T{g(), true}, n-1).bool
}

/// Empty list

func TestEmptyListIsEmptyAfterInit(t *testing.T) {
	e := newEmbedList()
	e.Init()
	if size := e.Len(); size != 0 {
		t.Errorf("init embedded element list has size not equal to 0: %v", size)
	}
	if f := e.Front(); f != nil {
		t.Errorf("init embedded element list has front %p", f)
	}
	if b := e.Back(); b != nil {
		t.Errorf("init embedded element list has back %p", b)
	}

	m := newMemberList()
	m.Init()
	if size := m.Len(); size != 0 {
		t.Errorf("init member element list has size not equal to 0: %v", size)
	}
	if f := m.Front(); f != nil {
		t.Errorf("init member element list has front %p", f)
	}
	if b := m.Back(); b != nil {
		t.Errorf("init member element list has back %p", b)
	}
}

func TestEmptyListIsZeroSize(t *testing.T) {
	e := newEmbedList()
	if size := e.Len(); size != 0 {
		t.Errorf("new embedded element list has size not equal to 0: %v", size)
	}

	m := newMemberList()
	if size := m.Len(); size != 0 {
		t.Errorf("new member element list has size not equal to 0: %v", size)
	}
}

func TestEmptyListDoNotHaveFront(t *testing.T) {
	e := newEmbedList()
	if f := e.Front(); f != nil {
		t.Errorf("new embedded element list has front %p", f)
	}

	m := newMemberList()
	if f := m.Front(); f != nil {
		t.Errorf("new member element list has front %p", f)
	}
}

func TestEmptyListPushFrontAddOneElement(t *testing.T) {
	e := newEmbedList()
	if el := newEmbed(0); true {
		e.PushFront(&el)
		if size := e.Len(); size != 1 {
			t.Errorf("new embedded element list pushfront do not increase size %v", size)
		}
	}

	m := newMemberList()
	if el := newMember(0); true {
		m.PushFront(&el)
		if size := m.Len(); size != 1 {
			t.Errorf("new member element list pushfront do not increase size %v", size)
		}
	}
}

func TestEmptyListPushFrontAddOneElementAtFront(t *testing.T) {
	e := newEmbedList()
	if el := newEmbed(0); true {
		e.PushFront(&el)
		if tel := e.Front(); &el != tel {
			t.Errorf("new embedded element list pushfront front is not correct %p %p", &el, tel)
		}
	}

	m := newMemberList()
	if el := newMember(0); true {
		m.PushFront(&el)
		if tel := m.Front(); &el != tel {
			t.Errorf("new member element list pushfront front is not correct %p %p", &el, tel)
		}
	}
}

func TestEmptyListPushFrontAddOneElementAtFrontAndItIsTheBackElement(t *testing.T) {
	e := newEmbedList()
	if el := newEmbed(0); true {
		e.PushFront(&el)
		if tel := e.Back(); &el != tel {
			t.Errorf("new embedded element list pushfront back is not correct %p %p", &el, tel)
		}
	}

	m := newMemberList()
	if el := newMember(0); true {
		m.PushFront(&el)
		if tel := m.Back(); &el != tel {
			t.Errorf("new member element list pushfront back is not correct %p %p", &el, tel)
		}
	}
}

func TestEmptyListDoNotPopFront(t *testing.T) {
	e := newEmbedList()
	if b := e.PopFront(); b != nil {
		t.Errorf("new embedded element list poped front %p", b)
	}

	m := newMemberList()
	if b := m.PopFront(); b != nil {
		t.Errorf("new member element list poped front %p", b)
	}
}

func TestEmptyListDoNotHaveBack(t *testing.T) {
	e := newEmbedList()
	if b := e.Back(); b != nil {
		t.Errorf("new embedded element list has back %p", b)
	}

	m := newMemberList()
	if b := m.Back(); b != nil {
		t.Errorf("new member element list has back %p", b)
	}
}

func TestEmptyListPushBackAddOneElement(t *testing.T) {
	e := newEmbedList()
	if el := newEmbed(0); true {
		e.PushBack(&el)
		if size := e.Len(); size != 1 {
			t.Errorf("new embedded element list pushback do not increase size %v", size)
		}
	}

	m := newMemberList()
	if el := newMember(0); true {
		m.PushBack(&el)
		if size := m.Len(); size != 1 {
			t.Errorf("new member element list pushback do not increase size %v", size)
		}
	}
}

func TestEmptyListPushBackAddOneElementAtBack(t *testing.T) {
	e := newEmbedList()
	if el := newEmbed(0); true {
		e.PushBack(&el)
		if tel := e.Back(); &el != tel {
			t.Errorf("new embedded element list pushback back is not correct %p %p", &el, tel)
		}
	}

	m := newMemberList()
	if el := newMember(0); true {
		m.PushBack(&el)
		if tel := m.Back(); &el != tel {
			t.Errorf("new member element list pushback back is not correct %p %p", &el, tel)
		}
	}
}

func TestEmptyListPushBackAddOneElementAtBackAndItIsTheFrontElement(t *testing.T) {
	e := newEmbedList()
	if el := newEmbed(0); true {
		e.PushBack(&el)
		if tel := e.Front(); &el != tel {
			t.Errorf("new embedded element list pushback front is not correct %p %p", &el, tel)
		}
	}

	m := newMemberList()
	if el := newMember(0); true {
		m.PushBack(&el)
		if tel := m.Front(); &el != tel {
			t.Errorf("new member element list pushback front is not correct %p %p", &el, tel)
		}
	}
}

func TestEmptyListClearDoNotReturnAnyElements(t *testing.T) {
	e := newEmbedList()
	if c := e.Clear(); len(c) > 0 {
		t.Errorf("new embedded element list clear returned some elements %v", c)
	}

	m := newMemberList()
	if c := m.Clear(); len(c) > 0 {
		t.Errorf("new member element list clear returned some elements %v", c)
	}
}

func TestEmptyListReversDoseNothing(t *testing.T) {
	e := newEmbedList()
	e.Reverse()
	if size := e.Len(); size != 0 {
		t.Errorf("reversed embedded element list has size not equal to 0: %v", size)
	}
	if f := e.Front(); f != nil {
		t.Errorf("reversed embedded element list has front %p", f)
	}
	if b := e.Back(); b != nil {
		t.Errorf("reversed embedded element list has back %p", b)
	}

	m := newMemberList()
	m.Reverse()
	if size := m.Len(); size != 0 {
		t.Errorf("reversed member element list has size not equal to 0: %v", size)
	}
	if f := m.Front(); f != nil {
		t.Errorf("reversed member element list has front %p", f)
	}
	if b := m.Back(); b != nil {
		t.Errorf("reversed member element list has back %p", b)
	}
}

func TestEmptyListMedianIsNotReturned(t *testing.T) {
	e := newEmbedList()
	if b := e.Median(); b != nil {
		t.Errorf("new embedded element list returned some median element %p", b)
	}

	m := newMemberList()
	if b := m.Median(); b != nil {
		t.Errorf("new member element list returned some median element %p", b)
	}
}

func TestEmptyListSortDoseNothing(t *testing.T) {
	e := newEmbedList()
	e.Sort(lessEmbed)
	if size := e.Len(); size != 0 {
		t.Errorf("sorted embedded element list has size not equal to 0: %v", size)
	}
	if f := e.Front(); f != nil {
		t.Errorf("sorted embedded element list has front %p", f)
	}
	if b := e.Back(); b != nil {
		t.Errorf("sorted embedded element list has back %p", b)
	}

	m := newMemberList()
	m.Sort(lessMember)
	if size := m.Len(); size != 0 {
		t.Errorf("sorted member element list has size not equal to 0: %v", size)
	}
	if f := m.Front(); f != nil {
		t.Errorf("sorted member element list has front %p", f)
	}
	if b := m.Back(); b != nil {
		t.Errorf("sorted member element list has back %p", b)
	}
}

func TestEmptyListUniqueDoNotReturnAnyElements(t *testing.T) {
	e := newEmbedList()
	if c := e.Unique(lessEmbed); len(c) > 0 {
		t.Errorf("new embedded element list unique returned some elements %v", c)
	}

	m := newMemberList()
	if c := m.Unique(lessMember); len(c) > 0 {
		t.Errorf("new member element list unique returned some elements %v", c)
	}
}

/// One element in list

func TestOneElementListInitClearsList(t *testing.T) {
	e := newEmbedListGenerate(1, increment(0))
	e.Init()
	if size := e.Len(); size != 0 {
		t.Errorf("embedded element list after init has size not equal to 0: %v", size)
	}
	if f := e.Front(); f != nil {
		t.Errorf("embedded element list after init has front %p", f)
	}
	if b := e.Back(); b != nil {
		t.Errorf("embedded element list after init has back %p", b)
	}

	m := newMemberListGenerate(1, increment(0))
	m.Init()
	if size := m.Len(); size != 0 {
		t.Errorf("member element list after init has size not equal to 0: %v", size)
	}
	if f := m.Front(); f != nil {
		t.Errorf("member element list after init has front %p", f)
	}
	if b := m.Back(); b != nil {
		t.Errorf("member element list after init has back %p", b)
	}
}

func TestOneElementListHsCorrectSize(t *testing.T) {
	e := newEmbedListGenerate(1, increment(0))
	if size := e.Len(); size != 1 {
		t.Errorf("new embedded element list has size not equal to %v: %v", 1, size)
	}

	m := newMemberListGenerate(1, increment(0))
	if size := m.Len(); size != 1 {
		t.Errorf("new member element list has size not equal to %v: %v", 1, size)
	}
}

func TestOneElementListInsertAfterFrontIncreasesSize(t *testing.T) {
	e := newEmbedListGenerate(1, increment(0))
	if el := newEmbed(0); true {
		e.InsertAfter(e.Front(), &el)
		if size := e.Len(); size != 2 {
			t.Errorf("embedded element list after insert has size not equal to %v: %v", 2, size)
		}
	}

	m := newMemberListGenerate(1, increment(0))
	if el := newMember(0); true {
		m.InsertAfter(m.Front(), &el)
		if size := m.Len(); size != 2 {
			t.Errorf("member element list after insert has size not equal to %v: %v", 2, size)
		}
	}
}

func TestOneElementListInsertAfterFrontChangesBackToInserted(t *testing.T) {
	e := newEmbedListGenerate(1, increment(0))
	if el := newEmbed(0); true {
		e.InsertAfter(e.Front(), &el)
		if tel := e.Back(); &el != tel {
			t.Errorf("embedded element list after insert has incorrect back %p %p", &el, tel)
		}
	}

	m := newMemberListGenerate(1, increment(0))
	if el := newMember(0); true {
		m.InsertAfter(m.Front(), &el)
		if tel := m.Back(); &el != tel {
			t.Errorf("member element list after insert has incorrect back %p %p", &el, tel)
		}
	}
}

func TestOneElementListInsertAfterFrontDoNotChangeFront(t *testing.T) {
	e := newEmbedListGenerate(1, increment(0))
	if tiel := newEmbed(0); true {
		if el := e.Front(); true {
			e.InsertAfter(e.Front(), &tiel)
			if tel := e.Front(); tel != el {
				t.Errorf("embedded element list after remove has incorrect front %p %p", el, tel)
			}
		}
	}

	m := newMemberListGenerate(1, increment(0))
	if tiel := newMember(0); true {
		if el := m.Front(); true {
			m.InsertAfter(m.Front(), &tiel)
			if tel := m.Front(); tel != el {
				t.Errorf("member element list after remove has incorrect front %p %p", el, tel)
			}
		}
	}
}

func TestOneElementListRemoveAfterFrontDoNotChangeListSize(t *testing.T) {
	e := newEmbedListGenerate(1, increment(0))
	e.RemoveAfter(e.Front())
	if size := e.Len(); size != 1 {
		t.Errorf("embedded element list after remove has size not equal to %v: %v", 1, size)
	}

	m := newMemberListGenerate(1, increment(0))
	m.RemoveAfter(m.Front())
	if size := m.Len(); size != 1 {
		t.Errorf("member element list after remove has size not equal to %v: %v", 1, size)
	}
}

func TestOneElementListRemoveAfterFrontDoNotChangeBack(t *testing.T) {
	e := newEmbedListGenerate(1, increment(0))
	if el := e.Back(); true {
		e.RemoveAfter(e.Front())
		if tel := e.Back(); tel != el {
			t.Errorf("embedded element list after remove has incorrect back %p %p", el, tel)
		}
	}

	m := newMemberListGenerate(1, increment(0))
	if el := m.Back(); true {
		m.RemoveAfter(m.Front())
		if tel := m.Back(); tel != el {
			t.Errorf("member element list after remove has incorrect back %p %p", el, tel)
		}
	}
}

func TestOneElementListPushFrontIncreasesSize(t *testing.T) {
	e := newEmbedListGenerate(1, increment(0))
	if el := newEmbed(0); true {
		e.PushFront(&el)
		if size := e.Len(); size != 2 {
			t.Errorf("embedded element list after pushfront has size not equal to %v: %v", 2, size)
		}
	}

	m := newMemberListGenerate(1, increment(0))
	if el := newMember(0); true {
		m.PushFront(&el)
		if size := m.Len(); size != 2 {
			t.Errorf("member element list after pushfront has size not equal to %v: %v", 2, size)
		}
	}
}

func TestOneElementListPushFrontChangesFrontToInserted(t *testing.T) {
	e := newEmbedListGenerate(1, increment(0))
	if el := newEmbed(0); true {
		e.PushFront(&el)
		if tel := e.Front(); &el != tel {
			t.Errorf("embedded element list after pushfront has incorrect front %p %p", &el, tel)
		}
	}

	m := newMemberListGenerate(1, increment(0))
	if el := newMember(0); true {
		m.PushFront(&el)
		if tel := m.Front(); &el != tel {
			t.Errorf("member element list after pushfront has incorrect front %p %p", &el, tel)
		}
	}
}

func TestOneElementListPushFrontDoNotChangeBack(t *testing.T) {
	e := newEmbedListGenerate(1, increment(0))
	if tiel := newEmbed(0); true {
		if el := e.Back(); true {
			e.PushFront(&tiel)
			if tel := e.Back(); tel != el {
				t.Errorf("embedded element list after pushfront has incorrect back %p %p", el, tel)
			}
		}
	}

	m := newMemberListGenerate(1, increment(0))
	if tiel := newMember(0); true {
		if el := m.Back(); true {
			m.PushFront(&tiel)
			if tel := m.Back(); tel != el {
				t.Errorf("member element list after pushfront has incorrect back %p %p", el, tel)
			}
		}
	}
}

func TestOneElementListPopFrontReturnsFront(t *testing.T) {
	e := newEmbedListGenerate(1, increment(0))
	if f := e.Front(); true {
		if el := e.PopFront(); el != f {
			t.Errorf("embedded element list popfront do not return front %p %p", el, f)
		}
	}

	m := newMemberListGenerate(1, increment(0))
	if f := m.Front(); true {
		if el := m.PopFront(); el != f {
			t.Errorf("member element list popfront do not return front %p %p", el, f)
		}
	}
}

func TestOneElementListPopFrontReturnsBack(t *testing.T) {
	e := newEmbedListGenerate(1, increment(0))
	if b := e.Back(); true {
		if el := e.PopFront(); el != b {
			t.Errorf("embedded element list popfront do not return back %p %p", el, b)
		}
	}

	m := newMemberListGenerate(1, increment(0))
	if b := m.Back(); true {
		if el := m.PopFront(); el != b {
			t.Errorf("member element list popfront do not return back %p %p", el, b)
		}
	}
}

func TestOneElementListPopFrontMakeListEmpty(t *testing.T) {
	e := newEmbedListGenerate(1, increment(0))
	e.PopFront()
	if size := e.Len(); size != 0 {
		t.Errorf("embedded element list after popfront has size not equal to 0: %v", size)
	}
	if f := e.Front(); f != nil {
		t.Errorf("embedded element list after popfront has front %p", f)
	}
	if b := e.Back(); b != nil {
		t.Errorf("embedded element list after popfront has back %p", b)
	}

	m := newMemberListGenerate(1, increment(0))
	m.PopFront()
	if size := m.Len(); size != 0 {
		t.Errorf("member element list after popfront has size not equal to 0: %v", size)
	}
	if f := m.Front(); f != nil {
		t.Errorf("member element list after popfront has front %p", f)
	}
	if b := m.Back(); b != nil {
		t.Errorf("member element list after popfront has back %p", b)
	}
}

func TestOneElementListPushBackIncreasesSize(t *testing.T) {
	e := newEmbedListGenerate(1, increment(0))
	if el := newEmbed(0); true {
		e.PushBack(&el)
		if size := e.Len(); size != 2 {
			t.Errorf("embedded element list after pushback has size not equal to %v: %v", 2, size)
		}
	}

	m := newMemberListGenerate(1, increment(0))
	if el := newMember(0); true {
		m.PushBack(&el)
		if size := m.Len(); size != 2 {
			t.Errorf("member element list after pushback has size not equal to %v: %v", 2, size)
		}
	}
}

func TestOneElementListPushBackChangesBackToInserted(t *testing.T) {
	e := newEmbedListGenerate(1, increment(0))
	if el := newEmbed(0); true {
		e.PushBack(&el)
		if tel := e.Back(); &el != tel {
			t.Errorf("embedded element list after pushback has incorrect back %p %p", &el, tel)
		}
	}

	m := newMemberListGenerate(1, increment(0))
	if el := newMember(0); true {
		m.PushBack(&el)
		if tel := m.Back(); &el != tel {
			t.Errorf("member element list after pushback has incorrect back %p %p", &el, tel)
		}
	}
}

func TestOneElementListPushBackDoNotChangeFront(t *testing.T) {
	e := newEmbedListGenerate(1, increment(0))
	if tiel := newEmbed(0); true {
		if el := e.Front(); true {
			e.PushBack(&tiel)
			if tel := e.Front(); tel != el {
				t.Errorf("embedded element list after pushback has incorrect front %p %p", el, tel)
			}
		}
	}

	m := newMemberListGenerate(1, increment(0))
	if tiel := newMember(0); true {
		if el := m.Front(); true {
			m.PushBack(&tiel)
			if tel := m.Front(); tel != el {
				t.Errorf("member element list after pushback has incorrect front %p %p", el, tel)
			}
		}
	}
}

func TestOneElementListClearReturnsOneElement(t *testing.T) {
	e := newEmbedListGenerate(1, increment(0))
	if l := e.Clear(); len(l) != 1 {
		t.Errorf("embedded element list after clear return incorrect count of elements: %v", len(l))
	}

	m := newMemberListGenerate(1, increment(0))
	if l := m.Clear(); len(l) != 1 {
		t.Errorf("member element list after clear return incorrect count of elements: %v", len(l))
	}
}

func TestOneElementListClearReturnsPreviouslyLinkedElement(t *testing.T) {
	e := newEmbedListGenerate(1, increment(0))
	if el := e.Front(); true {
		if l := e.Clear(); l[0] != el {
			t.Errorf("embedded element list after clear return incorrect element %p %p", el, l[0])
		}
	}

	m := newMemberListGenerate(1, increment(0))
	if el := m.Front(); true {
		if l := m.Clear(); l[0] != el {
			t.Errorf("member element list after clear return incorrect element %p %p", el, l[0])
		}
	}
}

func TestOneElementListClearReturnedElementIsUnlinked(t *testing.T) {
	e := newEmbedListGenerate(1, increment(0))
	if l := e.Clear(); l[0].Next() != nil {
		t.Errorf("embedded element list after clear return linked element %p", l[0].Next())
	}

	m := newMemberListGenerate(1, increment(0))
	if l := m.Clear(); l[0].hook.Next() != nil {
		t.Errorf("member element list after clear return linked element %p", l[0].hook.Next())
	}
}

func TestOneElementListClearMakesListEmpty(t *testing.T) {
	e := newEmbedListGenerate(1, increment(0))
	e.Clear()
	if size := e.Len(); size != 0 {
		t.Errorf("embedded element list after clear has size not equal to 0: %v", size)
	}
	if f := e.Front(); f != nil {
		t.Errorf("embedded element list after clear has front %p", f)
	}
	if b := e.Back(); b != nil {
		t.Errorf("embedded element list after clear has back %p", b)
	}

	m := newMemberListGenerate(1, increment(0))
	m.Clear()
	if size := m.Len(); size != 0 {
		t.Errorf("member element list after clear has size not equal to 0: %v", size)
	}
	if f := m.Front(); f != nil {
		t.Errorf("member element list after clear has front %p", f)
	}
	if b := m.Back(); b != nil {
		t.Errorf("member element list after clear has back %p", b)
	}
}

func TestOneElementListReverseDoseNothing(t *testing.T) {
	e := newEmbedListGenerate(1, increment(0))
	if el := e.Front(); true {
		e.Reverse()
		if size := e.Len(); size != 1 {
			t.Errorf("embedded element list after reverse has size not equal to %v: %v", 1, size)
		}
		if f := e.Front(); f != el {
			t.Errorf("embedded element list after reverse has incorrect front %p %p", el, f)
		}
		if b := e.Back(); b != el {
			t.Errorf("embedded element list after reverse has incorrect back %p %p", el, b)
		}
	}

	m := newMemberListGenerate(1, increment(0))
	if el := m.Front(); true {
		m.Reverse()
		if size := m.Len(); size != 1 {
			t.Errorf("member element list after reverse has size not equal to %v: %v", 1, size)
		}
		if f := m.Front(); f != el {
			t.Errorf("member element list after reverse has incorrect front %p %p", el, f)
		}
		if b := m.Back(); b != el {
			t.Errorf("member element list after reverse has incorrect back %p %p", el, b)
		}
	}
}

func TestOneElementListMedianReturnsSingleElement(t *testing.T) {
	e := newEmbedListGenerate(1, increment(0))
	if el := e.Front(); true {
		if tel := e.Median(); tel != el {
			t.Errorf("embedded element list median returns incorrect element %p %p", el, tel)
		}
	}

	m := newMemberListGenerate(1, increment(0))
	if el := m.Front(); true {
		if tel := m.Median(); tel != el {
			t.Errorf("member element list median returns incorrect element %p %p", el, tel)
		}
	}
}

func TestOneElementListSortDoseNothing(t *testing.T) {
	e := newEmbedListGenerate(1, increment(0))
	if el := e.Front(); true {
		e.Sort(lessEmbed)
		if size := e.Len(); size != 1 {
			t.Errorf("embedded element list after sort has size not equal to %v: %v", 1, size)
		}
		if f := e.Front(); f != el {
			t.Errorf("embedded element list after sort has incorrect front %p %p", el, f)
		}
		if b := e.Back(); b != el {
			t.Errorf("embedded element list after sort has incorrect back %p %p", el, b)
		}
	}

	m := newMemberListGenerate(1, increment(0))
	if el := m.Front(); true {
		m.Sort(lessMember)
		if size := m.Len(); size != 1 {
			t.Errorf("member element list after sort has size not equal to %v: %v", 1, size)
		}
		if f := m.Front(); f != el {
			t.Errorf("member element list after sort has incorrect front %p %p", el, f)
		}
		if b := m.Back(); b != el {
			t.Errorf("member element list after sort has incorrect back %p %p", el, b)
		}
	}
}

func TestOneElementListUniqueDoNotReturnAnyElements(t *testing.T) {
	e := newEmbedListGenerate(1, increment(0))
	if c := e.Unique(lessEmbed); len(c) > 0 {
		t.Errorf("embedded element list unique returned some elements %v", c)
	}

	m := newMemberListGenerate(1, increment(0))
	if c := m.Unique(lessMember); len(c) > 0 {
		t.Errorf("member element list unique returned some elements %v", c)
	}
}

/// Two element in list

func TestTwoElementListInitClearsList(t *testing.T) {
	e := newEmbedListGenerate(2, increment(0))
	e.Init()
	if size := e.Len(); size != 0 {
		t.Errorf("embedded element list after init has size not equal to %v: %v", 0, size)
	}
	if f := e.Front(); f != nil {
		t.Errorf("embedded element list after init has front %p", f)
	}
	if b := e.Back(); b != nil {
		t.Errorf("embedded element list after init has back %p", b)
	}

	m := newMemberListGenerate(2, increment(0))
	m.Init()
	if size := m.Len(); size != 0 {
		t.Errorf("member element list after init has size not equal to %v: %v", 0, size)
	}
	if f := m.Front(); f != nil {
		t.Errorf("member element list after init has front %p", f)
	}
	if b := m.Back(); b != nil {
		t.Errorf("member element list after init has back %p", b)
	}
}

func TestTwoElementListHasCorrectSize(t *testing.T) {
	e := newEmbedListGenerate(2, increment(0))
	if size := e.Len(); size != 2 {
		t.Errorf("new embedded element list has size not equal to %v: %v", 2, size)
	}

	m := newMemberListGenerate(2, increment(0))
	if size := m.Len(); size != 2 {
		t.Errorf("new member element list has size not equal to %v: %v", 2, size)
	}
}

func TestTwoElementListInsertAfterFrontDoNotChangeBack(t *testing.T) {
	e := newEmbedListGenerate(2, increment(0))
	if el := e.Back(); true {
		ie := newEmbed(0)
		e.InsertAfter(e.Front(), &ie)
		if tel := e.Back(); tel != el {
			t.Errorf("embedded element list middle insert changes back %p %p", el, tel)
		}
	}

	m := newMemberListGenerate(2, increment(0))
	if el := m.Back(); true {
		ie := newMember(0)
		m.InsertAfter(m.Front(), &ie)
		if tel := m.Back(); tel != el {
			t.Errorf("member element list middle insert changes back %p %p", el, tel)
		}
	}
}

func TestTwoElementListInsertAfterBackChangesBack(t *testing.T) {
	e := newEmbedListGenerate(2, increment(0))
	if el := e.Back(); true {
		ie := newEmbed(0)
		e.InsertAfter(el, &ie)
		if tel := e.Back(); tel == el {
			t.Errorf("embedded element list back insert do not change back %p %p", el, tel)
		}
	}

	m := newMemberListGenerate(2, increment(0))
	if el := m.Back(); true {
		ie := newMember(0)
		m.InsertAfter(el, &ie)
		if tel := m.Back(); tel == el {
			t.Errorf("member element list back insert do not change back %p %p", el, tel)
		}
	}
}

func TestTwoElementListRemoveAfterFrontChangesSize(t *testing.T) {
	e := newEmbedListGenerate(2, increment(0))
	if e.RemoveAfter(e.Front()); e.Len() != 1 {
		t.Errorf("embedded element list remove has size not equal to %v: %v", 1, e.Len())
	}

	m := newMemberListGenerate(2, increment(0))
	if m.RemoveAfter(m.Front()); m.Len() != 1 {
		t.Errorf("member element list remove has size not equal to %v: %v", 1, m.Len())
	}
}

func TestTwoElementListRemoveAfterFrontReturnsBack(t *testing.T) {
	e := newEmbedListGenerate(2, increment(0))
	if el := e.Back(); true {
		if tel := e.RemoveAfter(e.Front()); tel != el {
			t.Errorf("embedded element list remove after front do not return back %p %p", el, tel)
		}
	}

	m := newMemberListGenerate(2, increment(0))
	if el := m.Back(); true {
		if tel := m.RemoveAfter(m.Front()); tel != el {
			t.Errorf("member element list remove after front do not return back %p %p", el, tel)
		}
	}
}

func TestTwoElementListRemoveAfterFrontChangesBack(t *testing.T) {
	e := newEmbedListGenerate(2, increment(0))
	if el := e.Front(); true {
		if e.RemoveAfter(el); e.Back() != el {
			t.Errorf("embedded element list remove after front do not change back %p %p", el, e.Back())
		}
	}

	m := newMemberListGenerate(2, increment(0))
	if el := m.Back(); true {
		if m.RemoveAfter(el); m.Back() != el {
			t.Errorf("member element list remove after front do not change back %p %p", el, e.Back())
		}
	}
}

func TestTwoElementListPopFrontDoNotChangeBack(t *testing.T) {
	e := newEmbedListGenerate(2, increment(0))
	if el := e.Back(); true {
		if e.PopFront(); e.Back() != el {
			t.Errorf("embedded element list popfront changed back %p %p", el, e.Back())
		}
	}

	m := newMemberListGenerate(2, increment(0))
	if el := m.Back(); true {
		if m.PopFront(); m.Back() != el {
			t.Errorf("member element list popfront changed back %p %p", el, m.Back())
		}
	}
}

func TestTwoElementListReverseChangesFrontAndBack(t *testing.T) {
	e := newEmbedListGenerate(2, increment(0))
	if f, b := e.Front(), e.Back(); true {
		if e.Reverse(); e.Front() != b || e.Back() != f {
			t.Errorf("embedded element list reverse do not change order of elements %p %p : %p %p", f, b, e.Front(), e.Back())
		}
	}

	m := newMemberListGenerate(2, increment(0))
	if f, b := m.Front(), m.Back(); true {
		if m.Reverse(); m.Front() != b || m.Back() != f {
			t.Errorf("member element list reverse do not change order of elements %p %p : %p %p", f, b, e.Front(), e.Back())
		}
	}
}

func TestTwoElementListMedianReturnsFront(t *testing.T) {
	e := newEmbedListGenerate(2, increment(0))
	if el := e.Front(); true {
		if tel := e.Median(); tel != el {
			t.Errorf("embedded element list median do not return front %p %p", el, tel)
		}
	}

	m := newMemberListGenerate(2, increment(0))
	if el := m.Front(); true {
		if tel := m.Median(); tel != el {
			t.Errorf("member element list median do not return front %p %p", el, tel)
		}
	}
}

func TestTwoElementListSortOrdersElements(t *testing.T) {
	e := newEmbedListGenerate(2, decrement(2))
	if e.Sort(lessEmbed); !isSorted(nextEmbed(e), 2) {
		t.Errorf("embedded element list sort do not order elements")
	}

	m := newMemberListGenerate(2, decrement(2))
	if m.Sort(lessMember); !isSorted(nextMember(m), 2) {
		t.Errorf("member element list sort do not order elements")
	}
}

func TestTwoElementListUniqueRemovesBackOnEquals(t *testing.T) {
	e := newEmbedListGenerate(2, static(1))
	if el := e.Back(); true {
		if tel := e.Unique(lessEmbed)[0]; tel != el {
			t.Errorf("embedded element list unique do have not removed back %p %p", el, tel)
		}
	}

	m := newMemberListGenerate(2, static(1))
	if el := m.Back(); true {
		if tel := m.Unique(lessMember)[0]; tel != el {
			t.Errorf("member element list unique do have not removed back %p %p", el, tel)
		}
	}
}

func TestTwoElementListUniqueDoNotRemoveAnyIfNotEquals(t *testing.T) {
	e := newEmbedListGenerate(2, increment(0))
	if l := e.Unique(lessEmbed); 0 < len(l) {
		t.Errorf("embedded element list unique removed some elements %v", len(l))
	}

	m := newMemberListGenerate(2, increment(0))
	if l := m.Unique(lessMember); 0 < len(l) {
		t.Errorf("member element list unique removed some elements %v", len(l))
	}
}

/// Three element in list (lists of all other lengths are same as three element list)

func TestThreeElementListNextWillIterateOverAllElements(t *testing.T) {
	e := newEmbedListGenerate(3, func(n int) int { return 1 << n })
	if a := fn.Reduce(nextEmbed(e), func(e int, n int) int { return n | e }, 0, e.Len()); a&(a+1) != 0 {
		t.Errorf("embedded element list next function do not iterate over all elements %b", a)
	}

	m := newMemberListGenerate(3, func(n int) int { return 1 << n })
	if a := fn.Reduce(nextMember(m), func(e int, n int) int { return n | e }, 0, m.Len()); a&(a+1) != 0 {
		t.Errorf("member element list next function do not iterate over all elements %b", a)
	}
}

func TestThreeElementListHasCorrectSize(t *testing.T) {
	e := newEmbedListGenerate(3, increment(0))
	if s := e.Len(); s != 3 {
		t.Errorf("embedded element list is not of an expected size %v %v", 3, e.Len())
	}

	m := newMemberListGenerate(3, increment(0))
	if s := m.Len(); s != 3 {
		t.Errorf("member element list is not of an expected size %v %v", 3, m.Len())
	}
}

func TestThreeElementListInsertAfterMiddleDoNotChangeBack(t *testing.T) {
	e := newEmbedListGenerate(3, increment(0))
	if middle, ne, el := e.Front().Next(), newEmbed(50), e.Back(); true {
		e.InsertAfter(middle, &ne)
		if tel := e.Back(); tel != el {
			t.Errorf("embedded element list insert after middle changes back %p %p", tel, el)
		}
	}

	m := newMemberListGenerate(3, increment(0))
	if middle, ne, el := m.Front().hook.Next(), newMember(50), m.Back(); true {
		m.InsertAfter(middle, &ne)
		if tel := m.Back(); tel != el {
			t.Errorf("member element list insert after middle changes back %p %p", tel, el)
		}
	}
}

func TestThreeElementListRemoveAfterFrontChangesSize(t *testing.T) {
	e := newEmbedListGenerate(3, increment(0))
	if sz := e.Len(); true {
		e.RemoveAfter(e.Front())
		if size := e.Len(); size == sz {
			t.Errorf("embedded element list remove after front do not change size %v %v", size, sz)
		}
	}

	m := newMemberListGenerate(3, increment(0))
	if sz := m.Len(); true {
		m.RemoveAfter(m.Front())
		if size := m.Len(); size == sz {
			t.Errorf("member element list remove after front do not change size %v %v", size, sz)
		}
	}
}

func TestThreeElementListRemoveAfterFrontReturnsMiddle(t *testing.T) {
	e := newEmbedListGenerate(3, increment(0))
	if el := e.Front().Next(); true {
		if tel := e.RemoveAfter(e.Front()); tel != el {
			t.Errorf("embedded element list remove after front do not return middle %p %p", tel, el)
		}
	}

	m := newMemberListGenerate(3, increment(0))
	if el := m.Front().hook.Next(); true {
		if tel := m.RemoveAfter(m.Front()); tel != el {
			t.Errorf("member element list remove after front do not return middle %p %p", tel, el)
		}
	}
}

func TestThreeElementListRemoveAfterMiddleChangesBack(t *testing.T) {
	e := newEmbedListGenerate(3, increment(0))
	if el := e.Front().Next(); true {
		e.RemoveAfter(el)
		if tel := e.Back(); tel != el {
			t.Errorf("embedded element list remove after middle do not change back to middle %p %p", tel, el)
		}
	}

	m := newMemberListGenerate(3, increment(0))
	if el := m.Front().hook.Next(); true {
		m.RemoveAfter(el)
		if tel := m.Back(); tel != el {
			t.Errorf("member element list remove after middle do not change back to middle %p %p", tel, el)
		}
	}
}

func TestThreeElementListRemoveAfterMiddleReturnsBack(t *testing.T) {
	e := newEmbedListGenerate(3, increment(0))
	if el, middle := e.Back(), e.Front().Next(); true {
		if tel := e.RemoveAfter(middle); tel != el {
			t.Errorf("embedded element list remove after middle do not return back %p %p", tel, el)
		}
	}

	m := newMemberListGenerate(3, increment(0))
	if el, middle := m.Back(), m.Front().hook.Next(); true {
		if tel := m.RemoveAfter(middle); tel != el {
			t.Errorf("member element list remove after middle do not return back %p %p", tel, el)
		}
	}
}

func TestThreeElementListReverseMiddleStaysTheSame(t *testing.T) {
	e := newEmbedListGenerate(3, increment(0))
	if el := e.Front().Next(); true {
		e.Reverse()
		if tel := e.Front().Next(); tel != el {
			t.Errorf("embedded element list reverse changes middle element %p %p", tel, el)
		}
	}

	m := newMemberListGenerate(3, increment(0))
	if el := m.Front().hook.Next(); true {
		m.Reverse()
		if tel := m.Front().hook.Next(); tel != el {
			t.Errorf("member element list reverse changes middle element %p %p", tel, el)
		}
	}
}

func TestThreeElementListMedianReturnsMiddle(t *testing.T) {
	e := newEmbedListGenerate(3, increment(0))
	if el := e.Front().Next(); true {
		if tel := e.Median(); tel != el {
			t.Errorf("embedded element list median do not return middle %p %p", tel, el)
		}
	}

	m := newMemberListGenerate(3, increment(0))
	if el := m.Front().hook.Next(); true {
		if tel := m.Median(); tel != el {
			t.Errorf("member element list median do not return middle %p %p", tel, el)
		}
	}
}

func TestThreeElementListSortOrdersElements(t *testing.T) {
	e := newEmbedListGenerate(3, decrement(3))
	if e.Sort(lessEmbed); !isSorted(nextEmbed(e), 3) {
		t.Errorf("embedded element list sort do not order elements")
	}

	m := newMemberListGenerate(3, decrement(3))
	if m.Sort(lessMember); !isSorted(nextMember(m), 3) {
		t.Errorf("member element list sort do not order elements")
	}
}

func TestThreeElementListUniqueRemovesMiddleAndBackInOrderIfAsked(t *testing.T) {
	e := newEmbedListGenerate(3, static(1))
	if middle, back := e.Front().Next(), e.Back(); true {
		if tel := e.Unique(lessEmbed); !(len(tel) == 2 && tel[0] == middle && tel[1] == back) {
			t.Errorf("embedded element list unique do not remove exact amount of elements with order specified %p %p %v", middle, back, tel)
		}
	}

	m := newMemberListGenerate(3, static(1))
	if middle, back := m.Front().hook.Next(), m.Back(); true {
		if tel := m.Unique(lessMember); !(len(tel) == 2 && tel[0] == middle && tel[1] == back) {
			t.Errorf("member element list unique do not remove exact amount of elements with order specified %p %p %v", middle, back, tel)
		}
	}
}

/// RemoveIf is a bit special it should be tested nearly the same way for all base cases so we use table tests

func TestListRemoveIfDoNotReturnAnyElements(t *testing.T) {
	tests := map[string]struct {
		count int
	}{
		"empty list":    {0},
		"one element":   {1},
		"two elements":  {2},
		"three element": {3},
	}
	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			e := newEmbedListGenerate(testCase.count, increment(0))
			if c := e.RemoveIf(func(*testEmbedItem) bool { return false }); len(c) > 0 {
				t.Errorf("embedded element list removeif returned some elements %v", c)
			}

			m := newMemberListGenerate(testCase.count, increment(0))
			if c := m.RemoveIf(func(*testMemberItem) bool { return false }); len(c) > 0 {
				t.Errorf("member element list removeif returned some elements %v", c)
			}
		})
	}
}

func TestListRemoveIfAskedReturnsExactCountOfElementsThatShouldBeRemoved(t *testing.T) {
	tests := map[string]struct {
		count int
	}{
		"one element":   {1},
		"two elements":  {2},
		"three element": {3},
	}
	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			e := newEmbedListGenerate(testCase.count, increment(0))
			if l := e.RemoveIf(func(*testEmbedItem) bool { return true }); len(l) != testCase.count {
				t.Errorf("embedded element list after removeif return incorrect count of elements: %v", len(l))
			}

			m := newMemberListGenerate(testCase.count, increment(0))
			if l := m.RemoveIf(func(*testMemberItem) bool { return true }); len(l) != testCase.count {
				t.Errorf("member element list after removeif return incorrect count of elements: %v", len(l))
			}
		})
	}
}

func TestListRemoveIfAskedReturnsPreviouslyLinkedElement(t *testing.T) {
	tests := map[string]struct {
		count int
	}{
		"one element":   {1},
		"two elements":  {2},
		"three element": {3},
	}
	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			e := newEmbedListGenerate(testCase.count, increment(0))
			if el := e.Front(); true {
				if l := e.RemoveIf(func(*testEmbedItem) bool { return true }); l[0] != el {
					t.Errorf("embedded element list after removeif return incorrect element %p %p", el, l[0])
				}
			}

			m := newMemberListGenerate(testCase.count, increment(0))
			if el := m.Front(); true {
				if l := m.RemoveIf(func(*testMemberItem) bool { return true }); l[0] != el {
					t.Errorf("member element list after removeif return incorrect element %p %p", el, l[0])
				}
			}
		})
	}
}

func TestListRemoveIfAskedReturnedElementsAreUnlinked(t *testing.T) {
	tests := map[string]struct {
		count int
	}{
		"one element":   {1},
		"two elements":  {2},
		"three element": {3},
	}
	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			e := newEmbedListGenerate(testCase.count, increment(0))
			if l := e.RemoveIf(func(*testEmbedItem) bool { return true }); fn.AnyOf(l, func(e *testEmbedItem) bool { return e.Next() != nil }) {
				t.Errorf("embedded element list after removeif return linked element %v", fn.Filter(l, func(e *testEmbedItem) bool { return e.Hook.Next() != nil }))
			}

			m := newMemberListGenerate(testCase.count, increment(0))
			if l := m.RemoveIf(func(*testMemberItem) bool { return true }); fn.AnyOf(l, func(e *testMemberItem) bool { return e.hook.Next() != nil }) {
				t.Errorf("member element list after removeif return linked element %v", fn.Filter(l, func(e *testMemberItem) bool { return e.hook.Next() != nil }))
			}
		})
	}
}

func TestListRemoveIfAskedMakesListEmpty(t *testing.T) {
	tests := map[string]struct {
		count int
	}{
		"empty list":    {0},
		"one element":   {1},
		"two elements":  {2},
		"three element": {3},
	}
	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			e := newEmbedListGenerate(testCase.count, increment(0))
			e.RemoveIf(func(*testEmbedItem) bool { return true })
			if size := e.Len(); size != 0 {
				t.Errorf("embedded element list after removeif has size not equal to 0: %v", size)
			}
			if f := e.Front(); f != nil {
				t.Errorf("embedded element list after removeif has front %p", f)
			}
			if b := e.Back(); b != nil {
				t.Errorf("embedded element list after removeif has back %p", b)
			}

			m := newMemberListGenerate(testCase.count, increment(0))
			m.RemoveIf(func(*testMemberItem) bool { return true })
			if size := m.Len(); size != 0 {
				t.Errorf("member element list after removeif has size not equal to 0: %v", size)
			}
			if f := m.Front(); f != nil {
				t.Errorf("member element list after removeif has front %p", f)
			}
			if b := m.Back(); b != nil {
				t.Errorf("member element list after removeif has back %p", b)
			}
		})
	}
}

/// Two lists

func TestTwoListsSpliceAfterCorrectlyOrdersElementsOfLists(t *testing.T) {
	tests := map[string]struct {
		listSizes [2]int
		at        int
		order     []int
	}{
		"110": {[2]int{1, 1}, 0, []int{0, 3}},
		"120": {[2]int{1, 2}, 0, []int{0, 3, 4}},
		"130": {[2]int{1, 3}, 0, []int{0, 3, 4, 5}},
		"210": {[2]int{2, 1}, 0, []int{0, 3, 1}},
		"211": {[2]int{2, 1}, 1, []int{0, 1, 3}},
		"220": {[2]int{2, 2}, 0, []int{0, 3, 4, 1}},
		"221": {[2]int{2, 2}, 1, []int{0, 1, 3, 4}},
		"230": {[2]int{2, 3}, 0, []int{0, 3, 4, 5, 1}},
		"231": {[2]int{2, 3}, 1, []int{0, 1, 3, 4, 5}},
		"310": {[2]int{3, 1}, 0, []int{0, 3, 1, 2}},
		"311": {[2]int{3, 1}, 1, []int{0, 1, 3, 2}},
		"312": {[2]int{3, 1}, 2, []int{0, 1, 2, 3}},
		"320": {[2]int{3, 2}, 0, []int{0, 3, 4, 1, 2}},
		"321": {[2]int{3, 2}, 1, []int{0, 1, 3, 4, 2}},
		"322": {[2]int{3, 2}, 2, []int{0, 1, 2, 3, 4}},
		"330": {[2]int{3, 3}, 0, []int{0, 3, 4, 5, 1, 2}},
		"331": {[2]int{3, 3}, 1, []int{0, 1, 3, 4, 5, 2}},
		"332": {[2]int{3, 3}, 2, []int{0, 1, 2, 3, 4, 5}},
	}
	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			e := newEmbedListGenerate(testCase.listSizes[0], increment(0))
			if at := e.Front(); true {
				for range testCase.at {
					at = at.Next()
				}
				other := newEmbedListGenerate(testCase.listSizes[1], increment(3))
				e.SpliceAfter(at, &other)
				if o := fn.Next(testCase.order); !fn.Reduce(nextEmbed(e), func(e int, a bool) bool {
					return a && e == o()
				}, true, e.Len()) {
					t.Errorf("embedded element list spliceafter do not retains set order of elements %v %v", testCase.order, fn.Apply(nextEmbed(e), fn.I, e.Len()))
				}
			}

			m := newMemberListGenerate(testCase.listSizes[0], increment(0))
			if at := m.Front(); true {
				for range testCase.at {
					at = at.hook.Next()
				}
				other := newMemberListGenerate(testCase.listSizes[1], increment(3))
				m.SpliceAfter(at, &other)
				if o := fn.Next(testCase.order); !fn.Reduce(nextMember(m), func(e int, a bool) bool {
					return a && e == o()
				}, true, m.Len()) {
					t.Errorf("member element list spliceafter do not retains set order of elements %v %v", testCase.order, fn.Apply(nextMember(m), fn.I, m.Len()))
				}
			}
		})
	}
}

func TestTwoListsSpliceAfterResultingListHasSizeOfTwoListsCombined(t *testing.T) {
	tests := map[string]struct {
		listSizes [2]int
		at        int
		size      int
	}{
		"110": {[2]int{1, 1}, 0, 2},
		"120": {[2]int{1, 2}, 0, 3},
		"130": {[2]int{1, 3}, 0, 4},
		"210": {[2]int{2, 1}, 0, 3},
		"211": {[2]int{2, 1}, 1, 3},
		"220": {[2]int{2, 2}, 0, 4},
		"221": {[2]int{2, 2}, 1, 4},
		"230": {[2]int{2, 3}, 0, 5},
		"231": {[2]int{2, 3}, 1, 5},
		"310": {[2]int{3, 1}, 0, 4},
		"311": {[2]int{3, 1}, 1, 4},
		"312": {[2]int{3, 1}, 2, 4},
		"320": {[2]int{3, 2}, 0, 5},
		"321": {[2]int{3, 2}, 1, 5},
		"322": {[2]int{3, 2}, 2, 5},
		"330": {[2]int{3, 3}, 0, 6},
		"331": {[2]int{3, 3}, 1, 6},
		"332": {[2]int{3, 3}, 2, 6},
	}
	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			e := newEmbedListGenerate(testCase.listSizes[0], increment(0))
			if at := e.Front(); true {
				for range testCase.at {
					at = at.Next()
				}
				other := newEmbedListGenerate(testCase.listSizes[1], increment(0))
				e.SpliceAfter(at, &other)
				if size := e.Len(); size != testCase.size {
					t.Errorf("embedded element list spliceafter has invalid size %v %v", testCase.size, size)
				}
			}

			m := newMemberListGenerate(testCase.listSizes[0], increment(0))
			if at := m.Front(); true {
				for range testCase.at {
					at = at.hook.Next()
				}
				other := newMemberListGenerate(testCase.listSizes[1], increment(0))
				m.SpliceAfter(at, &other)
				if size := m.Len(); size != testCase.size {
					t.Errorf("member element list spliceafter has invalid size %v %v", testCase.size, size)
				}
			}
		})
	}
}

func TestTwoListsSpliceAfterOtherListIsEmptyAfterwards(t *testing.T) {
	tests := map[string]struct {
		listSizes [2]int
		at        int
		size      int
	}{
		"110": {[2]int{1, 1}, 0, 2},
		"120": {[2]int{1, 2}, 0, 3},
		"130": {[2]int{1, 3}, 0, 4},
		"210": {[2]int{2, 1}, 0, 3},
		"211": {[2]int{2, 1}, 1, 3},
		"220": {[2]int{2, 2}, 0, 4},
		"221": {[2]int{2, 2}, 1, 4},
		"230": {[2]int{2, 3}, 0, 5},
		"231": {[2]int{2, 3}, 1, 5},
		"310": {[2]int{3, 1}, 0, 4},
		"311": {[2]int{3, 1}, 1, 4},
		"312": {[2]int{3, 1}, 2, 4},
		"320": {[2]int{3, 2}, 0, 5},
		"321": {[2]int{3, 2}, 1, 5},
		"322": {[2]int{3, 2}, 2, 5},
		"330": {[2]int{3, 3}, 0, 6},
		"331": {[2]int{3, 3}, 1, 6},
		"332": {[2]int{3, 3}, 2, 6},
	}
	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			e := newEmbedListGenerate(testCase.listSizes[0], increment(0))
			if at := e.Front(); true {
				for range testCase.at {
					at = at.Next()
				}
				other := newEmbedListGenerate(testCase.listSizes[1], increment(0))
				e.SpliceAfter(at, &other)
				if size := other.Len(); size != 0 {
					t.Errorf("embedded element list after spliceafter has size not equal to 0: %v", size)
				}
				if f := other.Front(); f != nil {
					t.Errorf("embedded element list after spliceafter has front %p", f)
				}
				if b := other.Back(); b != nil {
					t.Errorf("embedded element list after spliceafter has back %p", b)
				}
			}

			m := newMemberListGenerate(testCase.listSizes[0], increment(0))
			if at := m.Front(); true {
				for range testCase.at {
					at = at.hook.Next()
				}
				other := newMemberListGenerate(testCase.listSizes[1], increment(0))
				m.SpliceAfter(at, &other)
				if size := other.Len(); size != 0 {
					t.Errorf("member element list after spliceafter has size not equal to 0: %v", size)
				}
				if f := other.Front(); f != nil {
					t.Errorf("member element list after spliceafter has front %p", f)
				}
				if b := other.Back(); b != nil {
					t.Errorf("member element list after spliceafter has back %p", b)
				}
			}
		})
	}
}

func TestTwoListsSpliceAfterWithEmptyListDoNotChangeAnything(t *testing.T) {
	tests := map[string]struct {
		listSize int
		at       int
		size     int
	}{
		"101": {1, 0, 1},
		"202": {2, 0, 2},
		"212": {2, 1, 2},
		"303": {3, 0, 3},
		"313": {3, 1, 3},
		"323": {3, 2, 3},
	}
	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			e := newEmbedListGenerate(testCase.listSize, increment(0))
			if at := e.Front(); true {
				for range testCase.at {
					at = at.Next()
				}
				other := newEmbedListGenerate(0, increment(0))
				if f, b := e.Front(), e.Back(); true {
					e.SpliceAfter(at, &other)
					if size := e.Len(); size != testCase.size {
						t.Errorf("embedded element list after spliceafter has size not equal to %v %v", testCase.size, size)
					}
					if e.Front() != f {
						t.Errorf("embedded element list after spliceafter has front %p", f)
					}
					if e.Back() != b {
						t.Errorf("embedded element list after spliceafter has back %p", b)
					}
				}
			}

			m := newMemberListGenerate(testCase.listSize, increment(0))
			if at := m.Front(); true {
				for range testCase.at {
					at = at.hook.Next()
				}
				other := newMemberListGenerate(0, increment(0))
				if f, b := m.Front(), m.Back(); true {
					m.SpliceAfter(at, &other)
					if size := m.Len(); size != testCase.size {
						t.Errorf("embedded element list after spliceafter has size not equal to %v %v", testCase.size, size)
					}
					if m.Front() != f {
						t.Errorf("embedded element list after spliceafter has front %p", f)
					}
					if m.Back() != b {
						t.Errorf("embedded element list after spliceafter has back %p", b)
					}
				}
			}
		})
	}
}

func TestTwoListsSpliceFrontOrdersAllElementsOfOneBeforeOther(t *testing.T) {
	tests := map[string]struct {
		listSizes [2]int
		order     []int
	}{
		"110": {[2]int{1, 1}, []int{3, 0}},
		"120": {[2]int{1, 2}, []int{3, 4, 0}},
		"130": {[2]int{1, 3}, []int{3, 4, 5, 0}},
		"210": {[2]int{2, 1}, []int{3, 0, 1}},
		"220": {[2]int{2, 2}, []int{3, 4, 0, 1}},
		"230": {[2]int{2, 3}, []int{3, 4, 5, 0, 1}},
		"310": {[2]int{3, 1}, []int{3, 0, 1, 2}},
		"320": {[2]int{3, 2}, []int{3, 4, 0, 1, 2}},
		"330": {[2]int{3, 3}, []int{3, 4, 5, 0, 1, 2}},
	}
	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			e := newEmbedListGenerate(testCase.listSizes[0], increment(0))
			if other := newEmbedListGenerate(testCase.listSizes[1], increment(3)); true {
				e.SpliceFront(&other)
				if o := fn.Next(testCase.order); !fn.Reduce(nextEmbed(e), func(e int, a bool) bool {
					return a && e == o()
				}, true, e.Len()) {
					t.Errorf("embedded element list splicefront do not retains set order of elements %v %v", testCase.order, fn.Apply(nextEmbed(e), fn.I, e.Len()))
				}
			}

			m := newMemberListGenerate(testCase.listSizes[0], increment(0))
			if other := newMemberListGenerate(testCase.listSizes[1], increment(3)); true {
				m.SpliceFront(&other)
				if o := fn.Next(testCase.order); !fn.Reduce(nextMember(m), func(e int, a bool) bool {
					return a && e == o()
				}, true, m.Len()) {
					t.Errorf("member element list splicefront do not retains set order of elements %v %v", testCase.order, fn.Apply(nextMember(m), fn.I, m.Len()))
				}
			}
		})
	}
}

func TestTwoListsSpliceFrontResultingListHasSizeOfTwoListsCombined(t *testing.T) {
	tests := map[string]struct {
		listSizes [2]int
		size      int
	}{
		"110": {[2]int{1, 1}, 2},
		"120": {[2]int{1, 2}, 3},
		"130": {[2]int{1, 3}, 4},
		"210": {[2]int{2, 1}, 3},
		"220": {[2]int{2, 2}, 4},
		"230": {[2]int{2, 3}, 5},
		"310": {[2]int{3, 1}, 4},
		"320": {[2]int{3, 2}, 5},
		"330": {[2]int{3, 3}, 6},
	}
	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			e := newEmbedListGenerate(testCase.listSizes[0], increment(0))
			if other := newEmbedListGenerate(testCase.listSizes[1], increment(0)); true {
				e.SpliceFront(&other)
				if size := e.Len(); size != testCase.size {
					t.Errorf("embedded element list splicefront has invalid size %v %v", testCase.size, size)
				}
			}

			m := newMemberListGenerate(testCase.listSizes[0], increment(0))
			if other := newMemberListGenerate(testCase.listSizes[1], increment(0)); true {
				m.SpliceFront(&other)
				if size := m.Len(); size != testCase.size {
					t.Errorf("embedded element list splicefront has invalid size %v %v", testCase.size, size)
				}
			}
		})
	}
}

func TestTwoListsSpliceFrontOtherListIsEmptyAfterwards(t *testing.T) {
	tests := map[string]struct {
		listSizes [2]int
		size      int
	}{
		"110": {[2]int{1, 1}, 2},
		"120": {[2]int{1, 2}, 3},
		"130": {[2]int{1, 3}, 4},
		"210": {[2]int{2, 1}, 3},
		"220": {[2]int{2, 2}, 4},
		"230": {[2]int{2, 3}, 5},
		"310": {[2]int{3, 1}, 4},
		"320": {[2]int{3, 2}, 5},
		"330": {[2]int{3, 3}, 6},
	}
	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			e := newEmbedListGenerate(testCase.listSizes[0], increment(0))
			if other := newEmbedListGenerate(testCase.listSizes[1], increment(0)); true {
				e.SpliceFront(&other)
				if size := other.Len(); size != 0 {
					t.Errorf("embedded element list after splicefront has size not equal to 0: %v", size)
				}
				if f := other.Front(); f != nil {
					t.Errorf("embedded element list after splicefront has front %p", f)
				}
				if b := other.Back(); b != nil {
					t.Errorf("embedded element list after splicefront has back %p", b)
				}
			}

			m := newMemberListGenerate(testCase.listSizes[0], increment(0))
			if other := newMemberListGenerate(testCase.listSizes[1], increment(0)); true {
				m.SpliceFront(&other)
				if size := other.Len(); size != 0 {
					t.Errorf("member element list after splicefront has size not equal to 0: %v", size)
				}
				if f := other.Front(); f != nil {
					t.Errorf("member element list after splicefront has front %p", f)
				}
				if b := other.Back(); b != nil {
					t.Errorf("member element list after splicefront has back %p", b)
				}
			}
		})
	}
}

func TestTwoListsSpliceFrontWithEmptyListDoNotChangeAnything(t *testing.T) {
	tests := map[string]struct {
		listSize int
		at       int
		size     int
	}{
		"101": {1, 0, 1},
		"202": {2, 0, 2},
		"212": {2, 1, 2},
		"303": {3, 0, 3},
		"313": {3, 1, 3},
		"323": {3, 2, 3},
	}
	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			e := newEmbedListGenerate(testCase.listSize, increment(0))
			if other := newEmbedListGenerate(0, increment(0)); true {
				if f, b := e.Front(), e.Back(); true {
					e.SpliceFront(&other)
					if size := e.Len(); size != testCase.size {
						t.Errorf("embedded element list after spliceafter has size not equal to %v %v", testCase.size, size)
					}
					if e.Front() != f {
						t.Errorf("embedded element list after spliceafter has front %p", f)
					}
					if e.Back() != b {
						t.Errorf("embedded element list after spliceafter has back %p", b)
					}
				}
			}

			m := newMemberListGenerate(testCase.listSize, increment(0))
			if other := newMemberListGenerate(0, increment(0)); true {
				if f, b := m.Front(), m.Back(); true {
					m.SpliceFront(&other)
					if size := m.Len(); size != testCase.size {
						t.Errorf("embedded element list after spliceafter has size not equal to %v %v", testCase.size, size)
					}
					if m.Front() != f {
						t.Errorf("embedded element list after spliceafter has front %p", f)
					}
					if m.Back() != b {
						t.Errorf("embedded element list after spliceafter has back %p", b)
					}
				}
			}
		})
	}
}

func TestTwoListsSpliceBackOrdersAllElementsOfOneBeforeOther(t *testing.T) {
	tests := map[string]struct {
		listSizes [2]int
		order     []int
	}{
		"110": {[2]int{1, 1}, []int{0, 3}},
		"120": {[2]int{1, 2}, []int{0, 3, 4}},
		"130": {[2]int{1, 3}, []int{0, 3, 4, 5}},
		"210": {[2]int{2, 1}, []int{0, 1, 3}},
		"220": {[2]int{2, 2}, []int{0, 1, 3, 4}},
		"230": {[2]int{2, 3}, []int{0, 1, 3, 4, 5}},
		"310": {[2]int{3, 1}, []int{0, 1, 2, 3}},
		"320": {[2]int{3, 2}, []int{0, 1, 2, 3, 4}},
		"330": {[2]int{3, 3}, []int{0, 1, 2, 3, 4, 5}},
	}
	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			e := newEmbedListGenerate(testCase.listSizes[0], increment(0))
			if other := newEmbedListGenerate(testCase.listSizes[1], increment(3)); true {
				e.SpliceBack(&other)
				if o := fn.Next(testCase.order); !fn.Reduce(nextEmbed(e), func(e int, a bool) bool {
					return a && e == o()
				}, true, e.Len()) {
					t.Errorf("embedded element list spliceback do not retains set order of elements %v %v", testCase.order, fn.Apply(nextEmbed(e), fn.I, e.Len()))
				}
			}

			m := newMemberListGenerate(testCase.listSizes[0], increment(0))
			if other := newMemberListGenerate(testCase.listSizes[1], increment(3)); true {
				m.SpliceBack(&other)
				if o := fn.Next(testCase.order); !fn.Reduce(nextMember(m), func(e int, a bool) bool {
					return a && e == o()
				}, true, m.Len()) {
					t.Errorf("member element list spliceback do not retains set order of elements %v %v", testCase.order, fn.Apply(nextMember(m), fn.I, m.Len()))
				}
			}
		})
	}
}

func TestTwoListsSpliceBackResultingListHasSizeOfTwoListsCombined(t *testing.T) {
	tests := map[string]struct {
		listSizes [2]int
		size      int
	}{
		"110": {[2]int{1, 1}, 2},
		"120": {[2]int{1, 2}, 3},
		"130": {[2]int{1, 3}, 4},
		"210": {[2]int{2, 1}, 3},
		"220": {[2]int{2, 2}, 4},
		"230": {[2]int{2, 3}, 5},
		"310": {[2]int{3, 1}, 4},
		"320": {[2]int{3, 2}, 5},
		"330": {[2]int{3, 3}, 6},
	}
	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			e := newEmbedListGenerate(testCase.listSizes[0], increment(0))
			if other := newEmbedListGenerate(testCase.listSizes[1], increment(0)); true {
				e.SpliceBack(&other)
				if size := e.Len(); size != testCase.size {
					t.Errorf("embedded element list spliceback has invalid size %v %v", testCase.size, size)
				}
			}

			m := newMemberListGenerate(testCase.listSizes[0], increment(0))
			if other := newMemberListGenerate(testCase.listSizes[1], increment(0)); true {
				m.SpliceBack(&other)
				if size := m.Len(); size != testCase.size {
					t.Errorf("embedded element list spliceback has invalid size %v %v", testCase.size, size)
				}
			}
		})
	}
}

func TestTwoListsSpliceBackOtherListIsEmptyAfterwards(t *testing.T) {
	tests := map[string]struct {
		listSizes [2]int
		size      int
	}{
		"110": {[2]int{1, 1}, 2},
		"120": {[2]int{1, 2}, 3},
		"130": {[2]int{1, 3}, 4},
		"210": {[2]int{2, 1}, 3},
		"220": {[2]int{2, 2}, 4},
		"230": {[2]int{2, 3}, 5},
		"310": {[2]int{3, 1}, 4},
		"320": {[2]int{3, 2}, 5},
		"330": {[2]int{3, 3}, 6},
	}
	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			e := newEmbedListGenerate(testCase.listSizes[0], increment(0))
			if other := newEmbedListGenerate(testCase.listSizes[1], increment(0)); true {
				e.SpliceBack(&other)
				if size := other.Len(); size != 0 {
					t.Errorf("embedded element list after spliceback has size not equal to 0: %v", size)
				}
				if f := other.Front(); f != nil {
					t.Errorf("embedded element list after spliceback has front %p", f)
				}
				if b := other.Back(); b != nil {
					t.Errorf("embedded element list after spliceback has back %p", b)
				}
			}

			m := newMemberListGenerate(testCase.listSizes[0], increment(0))
			if other := newMemberListGenerate(testCase.listSizes[1], increment(0)); true {
				m.SpliceBack(&other)
				if size := other.Len(); size != 0 {
					t.Errorf("member element list after spliceback has size not equal to 0: %v", size)
				}
				if f := other.Front(); f != nil {
					t.Errorf("member element list after spliceback has front %p", f)
				}
				if b := other.Back(); b != nil {
					t.Errorf("member element list after spliceback has back %p", b)
				}
			}
		})
	}
}

func TestTwoListsSpliceBackWithEmptyListDoNotChangeAnything(t *testing.T) {
	tests := map[string]struct {
		listSize int
		at       int
		size     int
	}{
		"101": {1, 0, 1},
		"202": {2, 0, 2},
		"212": {2, 1, 2},
		"303": {3, 0, 3},
		"313": {3, 1, 3},
		"323": {3, 2, 3},
	}
	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			e := newEmbedListGenerate(testCase.listSize, increment(0))
			if other := newEmbedListGenerate(0, increment(0)); true {
				if f, b := e.Front(), e.Back(); true {
					e.SpliceBack(&other)
					if size := e.Len(); size != testCase.size {
						t.Errorf("embedded element list after spliceafter has size not equal to %v %v", testCase.size, size)
					}
					if e.Front() != f {
						t.Errorf("embedded element list after spliceafter has front %p", f)
					}
					if e.Back() != b {
						t.Errorf("embedded element list after spliceafter has back %p", b)
					}
				}
			}

			m := newMemberListGenerate(testCase.listSize, increment(0))
			if other := newMemberListGenerate(0, increment(0)); true {
				if f, b := m.Front(), m.Back(); true {
					m.SpliceBack(&other)
					if size := m.Len(); size != testCase.size {
						t.Errorf("embedded element list after spliceafter has size not equal to %v %v", testCase.size, size)
					}
					if m.Front() != f {
						t.Errorf("embedded element list after spliceafter has front %p", f)
					}
					if m.Back() != b {
						t.Errorf("embedded element list after spliceafter has back %p", b)
					}
				}
			}
		})
	}
}

func TestTwoListsMergeOrdersElementsCorrectly(t *testing.T) {
	tests := map[string]struct {
		listSizes [2]int
		order     []int
	}{
		"11": {[2]int{1, 1}, []int{0, 3}},
		"12": {[2]int{1, 2}, []int{0, 3, 4}},
		"13": {[2]int{1, 3}, []int{0, 3, 4, 5}},
		"21": {[2]int{2, 1}, []int{0, 1, 3}},
		"22": {[2]int{2, 2}, []int{0, 1, 3, 4}},
		"23": {[2]int{2, 3}, []int{0, 1, 3, 4, 5}},
		"31": {[2]int{3, 1}, []int{0, 1, 2, 3}},
		"32": {[2]int{3, 2}, []int{0, 1, 2, 3, 4}},
		"33": {[2]int{3, 3}, []int{0, 1, 2, 3, 4, 5}},
	}
	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			e := newEmbedListGenerate(testCase.listSizes[0], increment(0))
			if other := newEmbedListGenerate(testCase.listSizes[1], increment(3)); true {
				e.Merge(&other, lessEmbed)
				e.SpliceBack(&other)
				if o := fn.Next(testCase.order); !fn.Reduce(nextEmbed(e), func(e int, a bool) bool {
					return a && e == o()
				}, true, e.Len()) {
					t.Errorf("embedded element list merge do not put elements in correct order %v %v", testCase.order, fn.Apply(nextEmbed(e), fn.I, e.Len()))
				}
			}

			m := newMemberListGenerate(testCase.listSizes[0], increment(0))
			if other := newMemberListGenerate(testCase.listSizes[1], increment(3)); true {
				m.Merge(&other, lessMember)
				if o := fn.Next(testCase.order); !fn.Reduce(nextMember(m), func(e int, a bool) bool {
					return a && e == o()
				}, true, m.Len()) {
					t.Errorf("member element list merge do not put elements in correct order %v %v", testCase.order, fn.Apply(nextMember(m), fn.I, m.Len()))
				}
			}
		})
	}
}

/// Other tests for 100% code coverage

func TestListMedianOfLengthFive(t *testing.T) {
	e := newEmbedListGenerate(5, decrement(5))
	if e.Sort(lessEmbed); !isSorted(nextEmbed(e), 5) {
		t.Errorf("embedded element list sort do not order elements")
	}

	m := newMemberListGenerate(5, decrement(5))
	if m.Sort(lessMember); !isSorted(nextMember(m), 5) {
		t.Errorf("member element list sort do not order elements")
	}
}

func TestListRemoveIfNotFromFront(t *testing.T) {
	tests := map[string]struct {
		size     int
		removed  []int
		remained []int
	}{
		"2": {2, []int{2}, []int{1}},
		"3": {3, []int{2}, []int{1, 3}},
	}
	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			e := newEmbedListGenerate(testCase.size, increment(1))
			if removed := e.RemoveIf(func(value *testEmbedItem) bool { return value.value%2 == 0 }); true {
				if !fn.AllOf(removed, func(value *testEmbedItem) bool { return value.value%2 == 0 }) {
					t.Errorf("embedded element list removeif failed to remove all even elements %v", removed)
				}
				if !fn.AllOf(fn.Apply(nextEmbed(e), fn.I, e.Len()), func(value int) bool { return value%2 != 0 }) {
					t.Errorf("embedded element list removeif failed to remove all even elements %v", fn.Apply(nextEmbed(e), fn.I, e.Len()))
				}
			}

			m := newMemberListGenerate(testCase.size, increment(0))
			if removed := m.RemoveIf(func(value *testMemberItem) bool { return value.value%2 == 0 }); true {
				if !fn.AllOf(removed, func(value *testMemberItem) bool { return value.value%2 == 0 }) {
					t.Errorf("member element list removeif failed to remove all even elements %v", removed)
				}
				if !fn.AllOf(fn.Apply(nextEmbed(e), fn.I, e.Len()), func(value int) bool { return value%2 != 0 }) {
					t.Errorf("member element list removeif failed to remove all even elements %v", fn.Apply(nextMember(m), fn.I, m.Len()))
				}
			}
		})
	}
}

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
