package slist

import "testing"

/// Possible use cases: 1. embedded hook 2. member hook

type testEmbedItem struct {
	Hook[testEmbedItem]
	value int
}

type embedOps struct{}

func (embedOps) Hook(self *testEmbedItem) *Hook[testEmbedItem] {
	return &self.Hook
}

func newEmbedList() SList[embedOps, testEmbedItem] {
	return New(embedOps{})
}

func newEmbed(value int) testEmbedItem {
	return testEmbedItem{Hook: NewHook[testEmbedItem](), value: value}
}

func newEmbedListGenerate(count int, generator func(position int) int) (l SList[embedOps, testEmbedItem]) {
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

type testMemberItem struct {
	hook  Hook[testMemberItem]
	value int
}

type memberOps struct{}

func (memberOps) Hook(self *testMemberItem) *Hook[testMemberItem] {
	return &self.hook
}

func newMemberList() SList[memberOps, testMemberItem] {
	return New(memberOps{})
}

func newMember(value int) testMemberItem {
	return testMemberItem{hook: NewHook[testMemberItem](), value: value}
}

func newMemberListGenerate(count int, generator func(position int) int) (l SList[memberOps, testMemberItem]) {
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

func iOTA(p int) int { return p }

/// Empty list

func TestEmptyListIsEmptyAfterInit(t *testing.T) {
	e := newEmbedList()
	e.Init()
	if size := e.Len(); size != 0 {
		t.Errorf("init embedded element list has size not equal to than 0: %v", size)
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
		t.Errorf("init member element list has size not equal to than 0: %v", size)
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
		t.Errorf("new embedded element list has size not equal to than 0: %v", size)
	}

	m := newMemberList()
	if size := m.Len(); size != 0 {
		t.Errorf("new member element list has size not equal to than 0: %v", size)
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
		t.Errorf("reversed embedded element list has size not equal to than 0: %v", size)
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
		t.Errorf("reversed member element list has size not equal to than 0: %v", size)
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
		t.Errorf("sorted embedded element list has size not equal to than 0: %v", size)
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
		t.Errorf("sorted member element list has size not equal to than 0: %v", size)
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

func TestEmptyListRemoveIfDoNotReturnAnyElements(t *testing.T) {
	e := newEmbedList()
	if c := e.RemoveIf(func(*testEmbedItem) bool { return true }); len(c) > 0 {
		t.Errorf("new embedded element list removeif returned some elements %v", c)
	}

	m := newMemberList()
	if c := m.RemoveIf(func(*testMemberItem) bool { return true }); len(c) > 0 {
		t.Errorf("new member element list removeif returned some elements %v", c)
	}
}

/// One element in list

func TestOneElementListInitClearsList(t *testing.T) {
	e := newEmbedListGenerate(1, iOTA)
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

	m := newMemberListGenerate(1, iOTA)
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

func TestOneElementListInsertAfterIncreasesSize(t *testing.T) {
	e := newEmbedListGenerate(1, iOTA)
	if el := newEmbed(0); true {
		e.InsertAfter(e.Front(), &el)
		if size := e.Len(); size != 2 {
			t.Errorf("embedded element list after insert has size not equal to %v: %v", 2, size)
		}
	}

	m := newMemberListGenerate(1, iOTA)
	if el := newMember(0); true {
		m.InsertAfter(m.Front(), &el)
		if size := m.Len(); size != 2 {
			t.Errorf("member element list after insert has size not equal to %v: %v", 2, size)
		}
	}
}

func TestOneElementListInsertAfterChangesBackToInserted(t *testing.T) {
	e := newEmbedListGenerate(1, iOTA)
	if el := newEmbed(0); true {
		e.InsertAfter(e.Front(), &el)
		if tel := e.Back(); &el != tel {
			t.Errorf("embedded element list after insert has incorrect back %p %p", &el, tel)
		}
	}

	m := newMemberListGenerate(1, iOTA)
	if el := newMember(0); true {
		m.InsertAfter(m.Front(), &el)
		if tel := m.Back(); &el != tel {
			t.Errorf("member element list after insert has incorrect back %p %p", &el, tel)
		}
	}
}

func TestOneElementListInsertAfterDoNotChangeFront(t *testing.T) {
	e := newEmbedListGenerate(1, iOTA)
	if tiel := newEmbed(0); true {
		if el := e.Front(); true {
			e.InsertAfter(e.Front(), &tiel)
			if tel := e.Front(); tel != el {
				t.Errorf("embedded element list after remove has incorrect front %p %p", el, tel)
			}
		}
	}

	m := newMemberListGenerate(1, iOTA)
	if tiel := newMember(0); true {
		if el := m.Front(); true {
			m.InsertAfter(m.Front(), &tiel)
			if tel := m.Front(); tel != el {
				t.Errorf("member element list after remove has incorrect front %p %p", el, tel)
			}
		}
	}
}

func TestOneElementListRemoveAfterDoNotChangeListSize(t *testing.T) {
	e := newEmbedListGenerate(1, iOTA)
	e.RemoveAfter(e.Front())
	if size := e.Len(); size != 1 {
		t.Errorf("embedded element list after remove has size not equal to %v: %v", 1, size)
	}

	m := newMemberListGenerate(1, iOTA)
	m.RemoveAfter(m.Front())
	if size := m.Len(); size != 1 {
		t.Errorf("member element list after remove has size not equal to %v: %v", 1, size)
	}
}

func TestOneElementListRemoveAfterDoNotChangeBack(t *testing.T) {
	e := newEmbedListGenerate(1, iOTA)
	if el := e.Back(); true {
		e.RemoveAfter(e.Front())
		if tel := e.Back(); tel != el {
			t.Errorf("embedded element list after remove has incorrect back %p %p", el, tel)
		}
	}

	m := newMemberListGenerate(1, iOTA)
	if el := m.Back(); true {
		m.RemoveAfter(m.Front())
		if tel := m.Back(); tel != el {
			t.Errorf("member element list after remove has incorrect back %p %p", el, tel)
		}
	}
}

func TestOneElementListPushFrontIncreasesSize(t *testing.T) {
	e := newEmbedListGenerate(1, iOTA)
	if el := newEmbed(0); true {
		e.PushFront(&el)
		if size := e.Len(); size != 2 {
			t.Errorf("embedded element list after pushfront has size not equal to %v: %v", 2, size)
		}
	}

	m := newMemberListGenerate(1, iOTA)
	if el := newMember(0); true {
		m.PushFront(&el)
		if size := m.Len(); size != 2 {
			t.Errorf("member element list after pushfront has size not equal to %v: %v", 2, size)
		}
	}
}

func TestOneElementListPushFrontChangesFrontToInserted(t *testing.T) {
	e := newEmbedListGenerate(1, iOTA)
	if el := newEmbed(0); true {
		e.PushFront(&el)
		if tel := e.Front(); &el != tel {
			t.Errorf("embedded element list after pushfront has incorrect front %p %p", &el, tel)
		}
	}

	m := newMemberListGenerate(1, iOTA)
	if el := newMember(0); true {
		m.PushFront(&el)
		if tel := m.Front(); &el != tel {
			t.Errorf("member element list after pushfront has incorrect front %p %p", &el, tel)
		}
	}
}

func TestOneElementListPushFrontDoNotChangeBack(t *testing.T) {
	e := newEmbedListGenerate(1, iOTA)
	if tiel := newEmbed(0); true {
		if el := e.Back(); true {
			e.PushFront(&tiel)
			if tel := e.Back(); tel != el {
				t.Errorf("embedded element list after pushfront has incorrect back %p %p", el, tel)
			}
		}
	}

	m := newMemberListGenerate(1, iOTA)
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
	e := newEmbedListGenerate(1, iOTA)
	if f := e.Front(); true {
		if el := e.PopFront(); el != f {
			t.Errorf("embedded element list popfront do not return front %p %p", el, f)
		}
	}

	m := newMemberListGenerate(1, iOTA)
	if f := m.Front(); true {
		if el := m.PopFront(); el != f {
			t.Errorf("member element list popfront do not return front %p %p", el, f)
		}
	}
}

func TestOneElementListPopFrontReturnsBack(t *testing.T) {
	e := newEmbedListGenerate(1, iOTA)
	if b := e.Back(); true {
		if el := e.PopFront(); el != b {
			t.Errorf("embedded element list popfront do not return back %p %p", el, b)
		}
	}

	m := newMemberListGenerate(1, iOTA)
	if b := m.Back(); true {
		if el := m.PopFront(); el != b {
			t.Errorf("member element list popfront do not return back %p %p", el, b)
		}
	}
}

func TestOneElementListPopFrontMakeListEmpty(t *testing.T) {
	e := newEmbedListGenerate(1, iOTA)
	e.PopFront()
	if size := e.Len(); size != 0 {
		t.Errorf("embedded element list after popfront has size not equal to than 0: %v", size)
	}
	if f := e.Front(); f != nil {
		t.Errorf("embedded element list after popfront has front %p", f)
	}
	if b := e.Back(); b != nil {
		t.Errorf("embedded element list after popfront has back %p", b)
	}

	m := newMemberListGenerate(1, iOTA)
	m.PopFront()
	if size := m.Len(); size != 0 {
		t.Errorf("member element list after popfront has size not equal to than 0: %v", size)
	}
	if f := m.Front(); f != nil {
		t.Errorf("member element list after popfront has front %p", f)
	}
	if b := m.Back(); b != nil {
		t.Errorf("member element list after popfront has back %p", b)
	}
}

func TestOneElementListPushBackIncreasesSize(t *testing.T) {
	e := newEmbedListGenerate(1, iOTA)
	if el := newEmbed(0); true {
		e.PushBack(&el)
		if size := e.Len(); size != 2 {
			t.Errorf("embedded element list after pushback has size not equal to %v: %v", 2, size)
		}
	}

	m := newMemberListGenerate(1, iOTA)
	if el := newMember(0); true {
		m.PushBack(&el)
		if size := m.Len(); size != 2 {
			t.Errorf("member element list after pushback has size not equal to %v: %v", 2, size)
		}
	}
}

func TestOneElementListPushBackChangesBackToInserted(t *testing.T) {
	e := newEmbedListGenerate(1, iOTA)
	if el := newEmbed(0); true {
		e.PushBack(&el)
		if tel := e.Back(); &el != tel {
			t.Errorf("embedded element list after pushback has incorrect back %p %p", &el, tel)
		}
	}

	m := newMemberListGenerate(1, iOTA)
	if el := newMember(0); true {
		m.PushBack(&el)
		if tel := m.Back(); &el != tel {
			t.Errorf("member element list after pushback has incorrect back %p %p", &el, tel)
		}
	}
}

func TestOneElementListPushBackDoNotChangeFront(t *testing.T) {
	e := newEmbedListGenerate(1, iOTA)
	if tiel := newEmbed(0); true {
		if el := e.Front(); true {
			e.PushBack(&tiel)
			if tel := e.Front(); tel != el {
				t.Errorf("embedded element list after pushback has incorrect front %p %p", el, tel)
			}
		}
	}

	m := newMemberListGenerate(1, iOTA)
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
	e := newEmbedListGenerate(1, iOTA)
	if l := e.Clear(); len(l) != 1 {
		t.Errorf("embedded element list after clear return incorrect count of elements: %v", len(l))
	}

	m := newMemberListGenerate(1, iOTA)
	if l := m.Clear(); len(l) != 1 {
		t.Errorf("member element list after clear return incorrect count of elements: %v", len(l))
	}
}

func TestOneElementListClearReturnsPreviouslyLinkedElement(t *testing.T) {
	e := newEmbedListGenerate(1, iOTA)
	if el := e.Front(); true {
		if l := e.Clear(); l[0] != el {
			t.Errorf("embedded element list after clear return incorrect element %p %p", el, l[0])
		}
	}

	m := newMemberListGenerate(1, iOTA)
	if el := m.Front(); true {
		if l := m.Clear(); l[0] != el {
			t.Errorf("member element list after clear return incorrect element %p %p", el, l[0])
		}
	}
}

func TestOneElementListClearReturnedElementIsUnlinked(t *testing.T) {
	e := newEmbedListGenerate(1, iOTA)
	if l := e.Clear(); l[0].Next() != nil {
		t.Errorf("embedded element list after clear return linked element %p", l[0].Next())
	}

	m := newMemberListGenerate(1, iOTA)
	if l := m.Clear(); l[0].hook.Next() != nil {
		t.Errorf("member element list after clear return linked element %p", l[0].hook.Next())
	}
}

func TestOneElementListClearMakesListEmpty(t *testing.T) {
	e := newEmbedListGenerate(1, iOTA)
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

	m := newMemberListGenerate(1, iOTA)
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
	e := newEmbedListGenerate(1, iOTA)
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

	m := newMemberListGenerate(1, iOTA)
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
	e := newEmbedListGenerate(1, iOTA)
	if el := e.Front(); true {
		if tel := e.Median(); tel != el {
			t.Errorf("embedded element list median returns incorrect element %p %p", el, tel)
		}
	}

	m := newMemberListGenerate(1, iOTA)
	if el := m.Front(); true {
		if tel := m.Median(); tel != el {
			t.Errorf("member element list median returns incorrect element %p %p", el, tel)
		}
	}
}

func TestOneElementListSortDoseNothing(t *testing.T) {
	e := newEmbedListGenerate(1, iOTA)
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

	m := newMemberListGenerate(1, iOTA)
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
	e := newEmbedListGenerate(1, iOTA)
	if c := e.Unique(lessEmbed); len(c) > 0 {
		t.Errorf("embedded element list unique returned some elements %v", c)
	}

	m := newMemberListGenerate(1, iOTA)
	if c := m.Unique(lessMember); len(c) > 0 {
		t.Errorf("member element list unique returned some elements %v", c)
	}
}

func TestOneElementListRemoveIfDoNotReturnAnyElements(t *testing.T) {
	e := newEmbedListGenerate(1, iOTA)
	if c := e.RemoveIf(func(*testEmbedItem) bool { return false }); len(c) > 0 {
		t.Errorf("embedded element list removeif returned some elements %v", c)
	}

	m := newMemberListGenerate(1, iOTA)
	if c := m.RemoveIf(func(*testMemberItem) bool { return false }); len(c) > 0 {
		t.Errorf("member element list removeif returned some elements %v", c)
	}
}

func TestOneElementListRemoveIfAskedReturnsOneElement(t *testing.T) {
	e := newEmbedListGenerate(1, iOTA)
	if l := e.RemoveIf(func(*testEmbedItem) bool { return true }); len(l) != 1 {
		t.Errorf("embedded element list after removeif return incorrect count of elements: %v", len(l))
	}

	m := newMemberListGenerate(1, iOTA)
	if l := m.RemoveIf(func(*testMemberItem) bool { return true }); len(l) != 1 {
		t.Errorf("member element list after removeif return incorrect count of elements: %v", len(l))
	}
}

func TestOneElementListRemoveIfAskedReturnsPreviouslyLinkedElement(t *testing.T) {
	e := newEmbedListGenerate(1, iOTA)
	if el := e.Front(); true {
		if l := e.RemoveIf(func(*testEmbedItem) bool { return true }); l[0] != el {
			t.Errorf("embedded element list after removeif return incorrect element %p %p", el, l[0])
		}
	}

	m := newMemberListGenerate(1, iOTA)
	if el := m.Front(); true {
		if l := m.RemoveIf(func(*testMemberItem) bool { return true }); l[0] != el {
			t.Errorf("member element list after removeif return incorrect element %p %p", el, l[0])
		}
	}
}

func TestOneElementListRemoveIfAskedReturnedElementIsUnlinked(t *testing.T) {
	e := newEmbedListGenerate(1, iOTA)
	if l := e.RemoveIf(func(*testEmbedItem) bool { return true }); l[0].Next() != nil {
		t.Errorf("embedded element list after removeif return linked element %p", l[0].Next())
	}

	m := newMemberListGenerate(1, iOTA)
	if l := m.RemoveIf(func(*testMemberItem) bool { return true }); l[0].hook.Next() != nil {
		t.Errorf("member element list after removeif return linked element %p", l[0].hook.Next())
	}
}

func TestOneElementListRemoveIfAskedMakesListEmpty(t *testing.T) {
	e := newEmbedListGenerate(1, iOTA)
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

	m := newMemberListGenerate(1, iOTA)
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
}

/// Two element in list

/// Three element in list (lists of all other lengths are same as three element list)
