package dlist

import (
	"testing"

	"github.com/echo-Mike/intrusive/internal/pkg/fn"
)

// Test structures and helper functions similar to SList tests

type testEmbedItem struct {
	Hook[testEmbedItem]
	value int
}

func embedHook(self *testEmbedItem) *Hook[testEmbedItem] {
	return &self.Hook
}

func newEmbedList() DList[testEmbedItem] {
	return New(embedHook)
}

func newEmbed(value int) testEmbedItem {
	return testEmbedItem{Hook: NewHook[testEmbedItem](), value: value}
}

func newEmbedListGenerate(count int, generator func(position int) int) (l DList[testEmbedItem]) {
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

func nextEmbed(l DList[testEmbedItem]) func() int {
	current := l.Front()
	return func() int {
		if current == nil {
			return 0
		}
		val := current.value
		current = current.Next()
		return val
	}
}

func prevEmbed(l DList[testEmbedItem]) func() int {
	current := l.Back()
	return func() int {
		if current == nil {
			return 0
		}
		val := current.value
		current = current.Prev()
		return val
	}
}

type testMemberItem struct {
	hook  Hook[testMemberItem]
	value int
}

func memberHook(self *testMemberItem) *Hook[testMemberItem] {
	return &self.hook
}

func newMemberList() DList[testMemberItem] {
	return New(memberHook)
}

func newMember(value int) testMemberItem {
	return testMemberItem{hook: NewHook[testMemberItem](), value: value}
}

func newMemberListGenerate(count int, generator func(position int) int) (l DList[testMemberItem]) {
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

func nextMember(l DList[testMemberItem]) func() int {
	current := l.Front()
	return func() int {
		if current == nil {
			return 0
		}
		val := current.value
		current = current.hook.Next()
		return val
	}
}

func prevMember(l DList[testMemberItem]) func() int {
	current := l.Back()
	return func() int {
		if current == nil {
			return 0
		}
		val := current.value
		current = current.hook.Prev()
		return val
	}
}

// Helper functions
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
	type state struct {
		prev  int
		valid bool
	}
	return fn.Reduce(g, func(a int, s state) state {
		return state{
			prev:  a,
			valid: s.valid && s.prev <= a,
		}
	}, state{g(), true}, n-1).valid
}

// / Empty list tests
func TestDListEmptyListIsEmptyAfterInit(t *testing.T) {
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

func TestDListEmptyListIsZeroSize(t *testing.T) {
	e := newEmbedList()
	if size := e.Len(); size != 0 {
		t.Errorf("new embedded element list has size not equal to 0: %v", size)
	}

	m := newMemberList()
	if size := m.Len(); size != 0 {
		t.Errorf("new member element list has size not equal to 0: %v", size)
	}
}

func TestDListEmptyListDoNotHaveFront(t *testing.T) {
	e := newEmbedList()
	if f := e.Front(); f != nil {
		t.Errorf("new embedded element list has front %p", f)
	}

	m := newMemberList()
	if f := m.Front(); f != nil {
		t.Errorf("new member element list has front %p", f)
	}
}

func TestDListEmptyListDoNotHaveBack(t *testing.T) {
	e := newEmbedList()
	if b := e.Back(); b != nil {
		t.Errorf("new embedded element list has back %p", b)
	}

	m := newMemberList()
	if b := m.Back(); b != nil {
		t.Errorf("new member element list has back %p", b)
	}
}

func TestDListEmptyListPushFrontAddOneElement(t *testing.T) {
	e := newEmbedList()
	el := newEmbed(0)
	e.PushFront(&el)
	if size := e.Len(); size != 1 {
		t.Errorf("new embedded element list pushfront do not increase size %v", size)
	}

	m := newMemberList()
	ml := newMember(0)
	m.PushFront(&ml)
	if size := m.Len(); size != 1 {
		t.Errorf("new member element list pushfront do not increase size %v", size)
	}
}

func TestDListEmptyListPushFrontAddOneElementAtFront(t *testing.T) {
	e := newEmbedList()
	el := newEmbed(0)
	e.PushFront(&el)
	if tel := e.Front(); &el != tel {
		t.Errorf("new embedded element list pushfront front is not correct %p %p", &el, tel)
	}

	m := newMemberList()
	ml := newMember(0)
	m.PushFront(&ml)
	if tel := m.Front(); &ml != tel {
		t.Errorf("new member element list pushfront front is not correct %p %p", &ml, tel)
	}
}

func TestDListEmptyListPushFrontAddOneElementAtFrontAndItIsTheBackElement(t *testing.T) {
	e := newEmbedList()
	el := newEmbed(0)
	e.PushFront(&el)
	if tel := e.Back(); &el != tel {
		t.Errorf("new embedded element list pushfront back is not correct %p %p", &el, tel)
	}

	m := newMemberList()
	ml := newMember(0)
	m.PushFront(&ml)
	if tel := m.Back(); &ml != tel {
		t.Errorf("new member element list pushfront back is not correct %p %p", &ml, tel)
	}
}

func TestDListEmptyListDoNotPopFront(t *testing.T) {
	e := newEmbedList()
	if b := e.PopFront(); b != nil {
		t.Errorf("new embedded element list poped front %p", b)
	}

	m := newMemberList()
	if b := m.PopFront(); b != nil {
		t.Errorf("new member element list poped front %p", b)
	}
}

func TestDListEmptyListPushBackAddOneElement(t *testing.T) {
	e := newEmbedList()
	el := newEmbed(0)
	e.PushBack(&el)
	if size := e.Len(); size != 1 {
		t.Errorf("new embedded element list pushback do not increase size %v", size)
	}

	m := newMemberList()
	ml := newMember(0)
	m.PushBack(&ml)
	if size := m.Len(); size != 1 {
		t.Errorf("new member element list pushback do not increase size %v", size)
	}
}

func TestDListEmptyListPushBackAddOneElementAtBack(t *testing.T) {
	e := newEmbedList()
	el := newEmbed(0)
	e.PushBack(&el)
	if tel := e.Back(); &el != tel {
		t.Errorf("new embedded element list pushback back is not correct %p %p", &el, tel)
	}

	m := newMemberList()
	ml := newMember(0)
	m.PushBack(&ml)
	if tel := m.Back(); &ml != tel {
		t.Errorf("new member element list pushback back is not correct %p %p", &ml, tel)
	}
}

func TestDListEmptyListPushBackAddOneElementAtBackAndItIsTheFrontElement(t *testing.T) {
	e := newEmbedList()
	el := newEmbed(0)
	e.PushBack(&el)
	if tel := e.Front(); &el != tel {
		t.Errorf("new embedded element list pushback front is not correct %p %p", &el, tel)
	}

	m := newMemberList()
	ml := newMember(0)
	m.PushBack(&ml)
	if tel := m.Front(); &ml != tel {
		t.Errorf("new member element list pushback front is not correct %p %p", &ml, tel)
	}
}

func TestDListEmptyListDoNotPopBack(t *testing.T) {
	e := newEmbedList()
	if b := e.PopBack(); b != nil {
		t.Errorf("new embedded element list poped back %p", b)
	}

	m := newMemberList()
	if b := m.PopBack(); b != nil {
		t.Errorf("new member element list poped back %p", b)
	}
}

func TestDListEmptyListClearDoNotReturnAnyElements(t *testing.T) {
	e := newEmbedList()
	if c := e.Clear(); len(c) > 0 {
		t.Errorf("new embedded element list clear returned some elements %v", c)
	}

	m := newMemberList()
	if c := m.Clear(); len(c) > 0 {
		t.Errorf("new member element list clear returned some elements %v", c)
	}
}

func TestDListEmptyListReverseDoesNothing(t *testing.T) {
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

// Continue with one element tests, two elements tests, etc.

func TestDListOneElementListInitClearsList(t *testing.T) {
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

func TestDListOneElementListHasCorrectSize(t *testing.T) {
	e := newEmbedListGenerate(1, increment(0))
	if size := e.Len(); size != 1 {
		t.Errorf("new embedded element list has size not equal to %v: %v", 1, size)
	}

	m := newMemberListGenerate(1, increment(0))
	if size := m.Len(); size != 1 {
		t.Errorf("new member element list has size not equal to %v: %v", 1, size)
	}
}

func TestDListOneElementListInsertBeforeFrontChangesSize(t *testing.T) {
	e := newEmbedListGenerate(1, increment(0))
	el := newEmbed(0)
	e.Insert(e.Front(), &el)
	if size := e.Len(); size != 2 {
		t.Errorf("embedded element list after insert has size not equal to %v: %v", 2, size)
	}

	m := newMemberListGenerate(1, increment(0))
	ml := newMember(0)
	m.Insert(m.Front(), &ml)
	if size := m.Len(); size != 2 {
		t.Errorf("member element list after insert has size not equal to %v: %v", 2, size)
	}
}

func TestDListOneElementListInsertBeforeFrontChangesFront(t *testing.T) {
	e := newEmbedListGenerate(1, increment(0))
	el := newEmbed(0)
	oldFront := e.Front()
	e.Insert(oldFront, &el)
	newFront := e.Front()
	if newFront != &el {
		t.Errorf("embedded element list after insert has incorrect front %p %p", &el, newFront)
	}
	if next := newFront.Next(); next != oldFront {
		t.Errorf("embedded element list after insert has incorrect next pointer %p %p", oldFront, next)
	}
	if prev := oldFront.Prev(); prev != newFront {
		t.Errorf("embedded element list after insert has incorrect prev pointer %p %p", newFront, prev)
	}

	m := newMemberListGenerate(1, increment(0))
	ml := newMember(0)
	oldMFront := m.Front()
	m.Insert(oldMFront, &ml)
	if newMFront := m.Front(); newMFront != &ml {
		t.Errorf("member element list after insert has incorrect front %p %p", &ml, newMFront)
	}
}

func TestDListTwoElementsReverseSwapsFrontAndBack(t *testing.T) {
	e := newEmbedListGenerate(2, increment(0))
	f, b := e.Front(), e.Back()
	e.Reverse()
	if e.Front() != b || e.Back() != f {
		t.Errorf("embedded element list reverse did not swap front and back")
	}
	if e.Front().Next() != f || e.Back().Prev() != b {
		t.Errorf("embedded element list reverse did not update pointers correctly")
	}

	m := newMemberListGenerate(2, increment(0))
	mf, mb := m.Front(), m.Back()
	m.Reverse()
	if m.Front() != mb || m.Back() != mf {
		t.Errorf("member element list reverse did not swap front and back")
	}
}

func TestDListOneElementListPopFrontReturnsFront(t *testing.T) {
	e := newEmbedListGenerate(1, increment(0))
	f := e.Front()
	if el := e.PopFront(); el != f {
		t.Errorf("embedded element list popfront do not return front %p %p", el, f)
	}

	m := newMemberListGenerate(1, increment(0))
	mf := m.Front()
	if el := m.PopFront(); el != mf {
		t.Errorf("member element list popfront do not return front %p %p", el, mf)
	}
}

func TestDListOneElementListPopFrontReturnsBack(t *testing.T) {
	e := newEmbedListGenerate(1, increment(0))
	b := e.Back()
	if el := e.PopFront(); el != b {
		t.Errorf("embedded element list popfront do not return back %p %p", el, b)
	}

	m := newMemberListGenerate(1, increment(0))
	mb := m.Back()
	if el := m.PopFront(); el != mb {
		t.Errorf("member element list popfront do not return back %p %p", el, mb)
	}
}

func TestDListOneElementListPopFrontMakesListEmpty(t *testing.T) {
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

func TestDListOneElementListPopBackReturnsBack(t *testing.T) {
	e := newEmbedListGenerate(1, increment(0))
	b := e.Back()
	if el := e.PopBack(); el != b {
		t.Errorf("embedded element list popback do not return back %p %p", el, b)
	}

	m := newMemberListGenerate(1, increment(0))
	mb := m.Back()
	if el := m.PopBack(); el != mb {
		t.Errorf("member element list popback do not return back %p %p", el, mb)
	}
}

func TestDListOneElementListPopBackReturnsFront(t *testing.T) {
	e := newEmbedListGenerate(1, increment(0))
	f := e.Front()
	if el := e.PopBack(); el != f {
		t.Errorf("embedded element list popback do not return front %p %p", el, f)
	}

	m := newMemberListGenerate(1, increment(0))
	mf := m.Front()
	if el := m.PopBack(); el != mf {
		t.Errorf("member element list popback do not return front %p %p", el, mf)
	}
}

func TestDListOneElementListPopBackMakesListEmpty(t *testing.T) {
	e := newEmbedListGenerate(1, increment(0))
	e.PopBack()
	if size := e.Len(); size != 0 {
		t.Errorf("embedded element list after popback has size not equal to 0: %v", size)
	}
	if f := e.Front(); f != nil {
		t.Errorf("embedded element list after popback has front %p", f)
	}
	if b := e.Back(); b != nil {
		t.Errorf("embedded element list after popback has back %p", b)
	}

	m := newMemberListGenerate(1, increment(0))
	m.PopBack()
	if size := m.Len(); size != 0 {
		t.Errorf("member element list after popback has size not equal to 0: %v", size)
	}
	if f := m.Front(); f != nil {
		t.Errorf("member element list after popback has front %p", f)
	}
	if b := m.Back(); b != nil {
		t.Errorf("member element list after popback has back %p", b)
	}
}

func TestDListOneElementListEraseMakesListEmpty(t *testing.T) {
	e := newEmbedListGenerate(1, increment(0))
	el := e.Front()
	e.Erase(el)
	if size := e.Len(); size != 0 {
		t.Errorf("embedded element list after erase has size not equal to 0: %v", size)
	}
	if f := e.Front(); f != nil {
		t.Errorf("embedded element list after erase has front %p", f)
	}
	if b := e.Back(); b != nil {
		t.Errorf("embedded element list after erase has back %p", b)
	}

	m := newMemberListGenerate(1, increment(0))
	ml := m.Front()
	m.Erase(ml)
	if size := m.Len(); size != 0 {
		t.Errorf("member element list after erase has size not equal to 0: %v", size)
	}
	if f := m.Front(); f != nil {
		t.Errorf("member element list after erase has front %p", f)
	}
	if b := m.Back(); b != nil {
		t.Errorf("member element list after erase has back %p", b)
	}
}

func TestDListTwoElementListHasCorrectPrevPointers(t *testing.T) {
	e := newEmbedListGenerate(2, increment(0))
	f, b := e.Front(), e.Back()
	if f.Next() != b {
		t.Errorf("embedded element list front next is not back")
	}
	if b.Prev() != f {
		t.Errorf("embedded element list back prev is not front")
	}

	m := newMemberListGenerate(2, increment(0))
	mf, mb := m.Front(), m.Back()
	if mf.hook.Next() != mb {
		t.Errorf("member element list front next is not back")
	}
	if mb.hook.Prev() != mf {
		t.Errorf("member element list back prev is not front")
	}
}

func TestDListThreeElementListHasCorrectPointers(t *testing.T) {
	e := newEmbedListGenerate(3, increment(0))
	f, m, b := e.Front(), e.Front().Next(), e.Back()

	if f.Next() != m {
		t.Errorf("embedded element list front next is not middle")
	}
	if m.Next() != b {
		t.Errorf("embedded element list middle next is not back")
	}
	if m.Prev() != f {
		t.Errorf("embedded element list middle prev is not front")
	}
	if b.Prev() != m {
		t.Errorf("embedded element list back prev is not middle")
	}

	// Test member version similarly
	ml := newMemberListGenerate(3, increment(0))
	mf, mm, mb := ml.Front(), ml.Front().hook.Next(), ml.Back()

	if mf.hook.Next() != mm {
		t.Errorf("member element list front next is not middle")
	}
	if mm.hook.Next() != mb {
		t.Errorf("member element list middle next is not back")
	}
	if mm.hook.Prev() != mf {
		t.Errorf("member element list middle prev is not front")
	}
	if mb.hook.Prev() != mm {
		t.Errorf("member element list back prev is not middle")
	}
}

func TestDListInsertBeforeFrontUpdatesPointersCorrectly(t *testing.T) {
	e := newEmbedListGenerate(2, increment(0))
	newEl := newEmbed(5)
	oldFront := e.Front()

	e.Insert(oldFront, &newEl)

	if e.Front() != &newEl {
		t.Errorf("insert before front did not update front pointer")
	}
	if newEl.Next() != oldFront {
		t.Errorf("new element next pointer incorrect")
	}
	if oldFront.Prev() != &newEl {
		t.Errorf("old front prev pointer incorrect")
	}
	if e.Back().Prev() != oldFront {
		t.Errorf("back element prev pointer incorrect")
	}

	// Test member version similarly
	ml := newMemberListGenerate(2, increment(0))
	newMl := newMember(5)
	oldMFront := ml.Front()

	ml.Insert(oldMFront, &newMl)

	if ml.Front() != &newMl {
		t.Errorf("insert before front did not update front pointer")
	}
	if newMl.hook.Next() != oldMFront {
		t.Errorf("new element next pointer incorrect")
	}
	if oldMFront.hook.Prev() != &newMl {
		t.Errorf("old front prev pointer incorrect")
	}
}

func TestDListInsertBeforeMiddleUpdatesPointersCorrectly(t *testing.T) {
	e := newEmbedListGenerate(3, increment(0))
	newEl := newEmbed(5)
	middle := e.Front().Next()

	e.Insert(middle, &newEl)

	if middle.Prev() != &newEl {
		t.Errorf("middle prev pointer not updated correctly")
	}
	if newEl.Next() != middle {
		t.Errorf("new element next pointer incorrect")
	}
	if newEl.Prev() != e.Front() {
		t.Errorf("new element prev pointer incorrect")
	}
	if e.Front().Next() != &newEl {
		t.Errorf("front next pointer not updated correctly")
	}

	// Test member version similarly
	ml := newMemberListGenerate(3, increment(0))
	newMl := newMember(5)
	mMiddle := ml.Front().hook.Next()

	ml.Insert(mMiddle, &newMl)

	if mMiddle.hook.Prev() != &newMl {
		t.Errorf("middle prev pointer not updated correctly")
	}
	if newMl.hook.Next() != mMiddle {
		t.Errorf("new element next pointer incorrect")
	}
	if newMl.hook.Prev() != ml.Front() {
		t.Errorf("new element prev pointer incorrect")
	}
	if ml.Front().hook.Next() != &newMl {
		t.Errorf("front next pointer not updated correctly")
	}
}

func TestDListEraseMiddleUpdatesPointersCorrectly(t *testing.T) {
	e := newEmbedListGenerate(3, increment(0))
	f, m, b := e.Front(), e.Front().Next(), e.Back()

	e.Erase(m)

	if f.Next() != b {
		t.Errorf("front next pointer not updated correctly after erase")
	}
	if b.Prev() != f {
		t.Errorf("back prev pointer not updated correctly after erase")
	}
	if m.Prev() != nil || m.Next() != nil {
		t.Errorf("erased element pointers not cleared")
	}

	// Test member version similarly
	ml := newMemberListGenerate(3, increment(0))
	mf, mm, mb := ml.Front(), ml.Front().hook.Next(), ml.Back()

	ml.Erase(mm)

	if mf.hook.Next() != mb {
		t.Errorf("front next pointer not updated correctly after erase")
	}
	if mb.hook.Prev() != mf {
		t.Errorf("back prev pointer not updated correctly after erase")
	}
	if mm.hook.Prev() != nil || mm.hook.Next() != nil {
		t.Errorf("erased element pointers not cleared")
	}
}

func TestDListReverseUpdatesAllPointersCorrectly(t *testing.T) {
	e := newEmbedListGenerate(3, increment(0))
	f, m, b := e.Front(), e.Front().Next(), e.Back()

	e.Reverse()

	if e.Front() != b || e.Back() != f {
		t.Errorf("reverse did not swap front and back correctly")
	}
	if b.Next() != m || m.Next() != f {
		t.Errorf("reverse did not update next pointers correctly")
	}
	if f.Prev() != m || m.Prev() != b {
		t.Errorf("reverse did not update prev pointers correctly")
	}

	// Test member version similarly
	ml := newMemberListGenerate(3, increment(0))
	mf, mm, mb := ml.Front(), ml.Front().hook.Next(), ml.Back()

	ml.Reverse()

	if ml.Front() != mb || ml.Back() != mf {
		t.Errorf("reverse did not swap front and back correctly")
	}
	if mb.hook.Next() != mm || mm.hook.Next() != mf {
		t.Errorf("reverse did not update next pointers correctly")
	}
	if mf.hook.Prev() != mm || mm.hook.Prev() != mb {
		t.Errorf("reverse did not update prev pointers correctly")
	}
}

func TestDListSpliceFrontMovesElementsCorrectly(t *testing.T) {
	e1 := newEmbedListGenerate(2, increment(0))
	e2 := newEmbedListGenerate(2, increment(10))

	e1.SpliceFront(&e2)

	if e1.Len() != 4 {
		t.Errorf("splice front did not combine sizes correctly")
	}
	if e2.Len() != 0 {
		t.Errorf("spliced list should be empty")
	}

	// Check order is correct
	expected := []int{10, 11, 0, 1}
	current := e1.Front()
	for i, exp := range expected {
		if current.value != exp {
			t.Errorf("position %d: expected %d, got %d", i, exp, current.value)
		}
		current = current.Next()
	}

	// Check prev pointers are correct
	current = e1.Front()
	for current != nil {
		if current.Next() != nil && current.Next().Prev() != current {
			t.Errorf("prev pointer incorrect for element %d", current.value)
		}
		current = current.Next()
	}

	// Test member version similarly
	m1 := newMemberListGenerate(2, increment(0))
	m2 := newMemberListGenerate(2, increment(10))

	m1.SpliceFront(&m2)

	if m1.Len() != 4 {
		t.Errorf("splice front did not combine sizes correctly")
	}
	if m2.Len() != 0 {
		t.Errorf("spliced list should be empty")
	}
}

func TestDListSpliceBackMovesElementsCorrectly(t *testing.T) {
	e1 := newEmbedListGenerate(2, increment(0))
	e2 := newEmbedListGenerate(2, increment(10))

	e1.SpliceBack(&e2)

	if e1.Len() != 4 {
		t.Errorf("splice back did not combine sizes correctly")
	}
	if e2.Len() != 0 {
		t.Errorf("spliced list should be empty")
	}

	// Check order is correct
	expected := []int{0, 1, 10, 11}
	current := e1.Front()
	for i, exp := range expected {
		if current.value != exp {
			t.Errorf("position %d: expected %d, got %d", i, exp, current.value)
		}
		current = current.Next()
	}

	// Check prev pointers are correct
	current = e1.Front()
	for current != nil {
		if current.Next() != nil && current.Next().Prev() != current {
			t.Errorf("prev pointer incorrect for element %d", current.value)
		}
		current = current.Next()
	}

	// Test member version similarly
	m1 := newMemberListGenerate(2, increment(0))
	m2 := newMemberListGenerate(2, increment(10))

	m1.SpliceBack(&m2)

	if m1.Len() != 4 {
		t.Errorf("splice back did not combine sizes correctly")
	}
	if m2.Len() != 0 {
		t.Errorf("spliced list should be empty")
	}
}

func TestDListSortOrdersElementsCorrectly(t *testing.T) {
	e := newEmbedListGenerate(5, decrement(5)) // Creates [4, 3, 2, 1, 0]
	e.Sort(lessEmbed)

	if !isSorted(nextEmbed(e), 5) {
		t.Errorf("embedded element list sort did not order elements correctly")
	}

	// Check prev pointers are correct
	current := e.Front()
	for current != nil {
		if current.Next() != nil && current.Next().Prev() != current {
			t.Errorf("prev pointer incorrect after sort for element %d", current.value)
		}
		current = current.Next()
	}

	// Test member version similarly
	m := newMemberListGenerate(5, decrement(5))
	m.Sort(lessMember)

	if !isSorted(nextMember(m), 5) {
		t.Errorf("member element list sort did not order elements correctly")
	}
}

func TestDListUniqueRemovesDuplicates(t *testing.T) {
	e := newEmbedListGenerate(5, func(n int) int {
		return n / 2 // Creates [0, 0, 1, 1, 2]
	})

	removed := e.Unique(lessEmbed)

	if len(removed) != 2 {
		t.Errorf("unique should have removed 2 elements, got %d", len(removed))
	}
	if e.Len() != 3 {
		t.Errorf("unique should leave 3 elements, got %d", e.Len())
	}

	// Check order is correct
	expected := []int{0, 1, 2}
	current := e.Front()
	for i, exp := range expected {
		if current.value != exp {
			t.Errorf("position %d: expected %d, got %d", i, exp, current.value)
		}
		current = current.Next()
	}

	// Test member version similarly
	m := newMemberListGenerate(5, func(n int) int {
		return n / 2
	})

	mRemoved := m.Unique(lessMember)

	if len(mRemoved) != 2 {
		t.Errorf("unique should have removed 2 elements, got %d", len(mRemoved))
	}
}

func TestDListRemoveIfFiltersElements(t *testing.T) {
	e := newEmbedListGenerate(5, increment(0)) // Creates [0, 1, 2, 3, 4]

	removed := e.RemoveIf(func(value *testEmbedItem) bool {
		return value.value%2 == 0 // Remove even numbers
	})

	if len(removed) != 3 {
		t.Errorf("removeif should have removed 3 elements, got %d", len(removed))
	}
	if e.Len() != 2 {
		t.Errorf("removeif should leave 2 elements, got %d", e.Len())
	}

	// Check remaining elements are correct
	current := e.Front()
	for current != nil {
		if current.value%2 != 1 {
			t.Errorf("element %d should not have been removed", current.value)
		}
		current = current.Next()
	}

	// Test member version similarly
	m := newMemberListGenerate(5, increment(0))

	mRemoved := m.RemoveIf(func(value *testMemberItem) bool {
		return value.value%2 == 0
	})

	if len(mRemoved) != 3 {
		t.Errorf("removeif should have removed 3 elements, got %d", len(mRemoved))
	}
}

func TestDListInsertBeforeNilAppendsToBack(t *testing.T) {
	e := newEmbedListGenerate(2, increment(0))
	newEl := newEmbed(5)
	oldBack := e.Back()

	e.Insert(nil, &newEl)

	if e.Back() != &newEl {
		t.Errorf("insert before nil did not append to back")
	}
	if oldBack.Next() != &newEl {
		t.Errorf("old back next pointer not updated correctly")
	}
	if newEl.Prev() != oldBack {
		t.Errorf("new element prev pointer incorrect")
	}
	if newEl.Next() != nil {
		t.Errorf("new element next pointer should be nil")
	}

	// Test member version
	m := newMemberListGenerate(2, increment(0))
	newMl := newMember(5)

	m.Insert(nil, &newMl)

	if m.Back() != &newMl {
		t.Errorf("insert before nil did not append to back")
	}
}

func TestDListSpliceAtPositionUpdatesPointersCorrectly(t *testing.T) {
	e1 := newEmbedListGenerate(3, increment(0))  // [0, 1, 2]
	e2 := newEmbedListGenerate(2, increment(10)) // [10, 11]
	position := e1.Front().Next()                // Element with value 1

	e1.Splice(position, &e2)

	if e1.Len() != 5 {
		t.Errorf("splice did not combine sizes correctly")
	}
	if e2.Len() != 0 {
		t.Errorf("spliced list should be empty")
	}

	// Check order: [0, 10, 11, 1, 2]
	expected := []int{0, 10, 11, 1, 2}
	current := e1.Front()
	for i, exp := range expected {
		if current.value != exp {
			t.Errorf("position %d: expected %d, got %d", i, exp, current.value)
		}
		current = current.Next()
	}

	// Check prev pointers
	current = e1.Front()
	for current != nil {
		if current.Next() != nil && current.Next().Prev() != current {
			t.Errorf("prev pointer incorrect for element %d", current.value)
		}
		current = current.Next()
	}

	// Test member version
	m1 := newMemberListGenerate(3, increment(0))
	m2 := newMemberListGenerate(2, increment(10))
	mPosition := m1.Front().hook.Next()

	m1.Splice(mPosition, &m2)

	if m1.Len() != 5 {
		t.Errorf("splice did not combine sizes correctly")
	}
}

func TestDListMergeSortedLists(t *testing.T) {
	e1 := newEmbedListGenerate(3, func(n int) int { return n * 2 })   // [0, 2, 4]
	e2 := newEmbedListGenerate(3, func(n int) int { return n*2 + 1 }) // [1, 3, 5]

	e1.Merge(&e2, lessEmbed)

	if e1.Len() != 6 {
		t.Errorf("merge did not combine sizes correctly")
	}
	if e2.Len() != 0 {
		t.Errorf("merged list should be empty")
	}

	// Check order: [0, 1, 2, 3, 4, 5]
	expected := []int{0, 1, 2, 3, 4, 5}
	current := e1.Front()
	for i, exp := range expected {
		if current.value != exp {
			t.Errorf("position %d: expected %d, got %d", i, exp, current.value)
		}
		current = current.Next()
	}

	// Check prev pointers
	current = e1.Front()
	for current != nil {
		if current.Next() != nil && current.Next().Prev() != current {
			t.Errorf("prev pointer incorrect for element %d", current.value)
		}
		current = current.Next()
	}

	// Test member version
	m1 := newMemberListGenerate(3, func(n int) int { return n * 2 })
	m2 := newMemberListGenerate(3, func(n int) int { return n*2 + 1 })

	m1.Merge(&m2, lessMember)

	if m1.Len() != 6 {
		t.Errorf("merge did not combine sizes correctly")
	}
}

func TestDListMergeWithEmptyList(t *testing.T) {
	e1 := newEmbedListGenerate(3, increment(0))
	e2 := newEmbedList()

	e1.Merge(&e2, lessEmbed)

	if e1.Len() != 3 {
		t.Errorf("merge with empty list should not change size")
	}
	if e2.Len() != 0 {
		t.Errorf("empty list should remain empty")
	}

	// Check original order preserved
	expected := []int{0, 1, 2}
	current := e1.Front()
	for i, exp := range expected {
		if current.value != exp {
			t.Errorf("position %d: expected %d, got %d", i, exp, current.value)
		}
		current = current.Next()
	}

	// Test member version
	m1 := newMemberListGenerate(3, increment(0))
	m2 := newMemberList()

	m1.Merge(&m2, lessMember)

	if m1.Len() != 3 {
		t.Errorf("merge with empty list should not change size")
	}
}

func TestDListMergeEmptyWithNonEmpty(t *testing.T) {
	e1 := newEmbedList()
	e2 := newEmbedListGenerate(3, increment(0))

	e1.Merge(&e2, lessEmbed)

	if e1.Len() != 3 {
		t.Errorf("merge should copy elements from non-empty list")
	}
	if e2.Len() != 0 {
		t.Errorf("merged list should be empty")
	}

	// Check order preserved
	expected := []int{0, 1, 2}
	current := e1.Front()
	for i, exp := range expected {
		if current.value != exp {
			t.Errorf("position %d: expected %d, got %d", i, exp, current.value)
		}
		current = current.Next()
	}

	// Test member version
	m1 := newMemberList()
	m2 := newMemberListGenerate(3, increment(0))

	m1.Merge(&m2, lessMember)

	if m1.Len() != 3 {
		t.Errorf("merge should copy elements from non-empty list")
	}
}

func TestDListClearReturnsAllElements(t *testing.T) {
	e := newEmbedListGenerate(4, increment(0))

	elements := e.Clear()

	if len(elements) != 4 {
		t.Errorf("clear should return all elements, got %d", len(elements))
	}
	if e.Len() != 0 {
		t.Errorf("list should be empty after clear")
	}
	if e.Front() != nil {
		t.Errorf("front should be nil after clear")
	}
	if e.Back() != nil {
		t.Errorf("back should be nil after clear")
	}

	// Check all returned elements are unlinked
	for _, el := range elements {
		if el.Next() != nil || el.Prev() != nil {
			t.Errorf("returned element should be unlinked")
		}
	}

	// Test member version
	m := newMemberListGenerate(4, increment(0))

	mElements := m.Clear()

	if len(mElements) != 4 {
		t.Errorf("clear should return all elements, got %d", len(mElements))
	}
}

func TestDListBackwardTraversal(t *testing.T) {
	e := newEmbedListGenerate(4, increment(0)) // [0, 1, 2, 3]

	// Test backward traversal
	expected := []int{3, 2, 1, 0}
	current := e.Back()
	for i, exp := range expected {
		if current.value != exp {
			t.Errorf("backward position %d: expected %d, got %d", i, exp, current.value)
		}
		current = current.Prev()
	}

	if current != nil {
		t.Errorf("should reach nil at end of backward traversal")
	}

	// Test member version
	m := newMemberListGenerate(4, increment(0))

	mExpected := []int{3, 2, 1, 0}
	mCurrent := m.Back()
	for i, exp := range mExpected {
		if mCurrent.value != exp {
			t.Errorf("backward position %d: expected %d, got %d", i, exp, mCurrent.value)
		}
		mCurrent = mCurrent.hook.Prev()
	}
}

func TestDListEraseFrontUpdatesPointers(t *testing.T) {
	e := newEmbedListGenerate(3, increment(0))
	oldFront := e.Front()
	newFront := oldFront.Next()

	e.Erase(oldFront)

	if e.Front() != newFront {
		t.Errorf("front pointer not updated after erase")
	}
	if newFront.Prev() != nil {
		t.Errorf("new front should not have prev pointer")
	}
	if oldFront.Next() != nil || oldFront.Prev() != nil {
		t.Errorf("erased element pointers not cleared")
	}

	// Test member version
	m := newMemberListGenerate(3, increment(0))
	oldMFront := m.Front()
	newMFront := oldMFront.hook.Next()

	m.Erase(oldMFront)

	if m.Front() != newMFront {
		t.Errorf("front pointer not updated after erase")
	}
}

func TestDListEraseBackUpdatesPointers(t *testing.T) {
	e := newEmbedListGenerate(3, increment(0))
	oldBack := e.Back()
	newBack := oldBack.Prev()

	e.Erase(oldBack)

	if e.Back() != newBack {
		t.Errorf("back pointer not updated after erase")
	}
	if newBack.Next() != nil {
		t.Errorf("new back should not have next pointer")
	}
	if oldBack.Next() != nil || oldBack.Prev() != nil {
		t.Errorf("erased element pointers not cleared")
	}

	// Test member version
	m := newMemberListGenerate(3, increment(0))
	oldMBack := m.Back()
	newMBack := oldMBack.hook.Prev()

	m.Erase(oldMBack)

	if m.Back() != newMBack {
		t.Errorf("back pointer not updated after erase")
	}
}

func TestDListSpliceIntoEmptyList(t *testing.T) {
	e1 := newEmbedList()
	e2 := newEmbedListGenerate(2, increment(0))

	e1.SpliceFront(&e2)

	if e1.Len() != 2 {
		t.Errorf("splice into empty list should copy all elements")
	}
	if e2.Len() != 0 {
		t.Errorf("spliced list should be empty")
	}
	if e1.Front().value != 0 || e1.Back().value != 1 {
		t.Errorf("elements not copied correctly")
	}

	// Test member version
	m1 := newMemberList()
	m2 := newMemberListGenerate(2, increment(0))

	m1.SpliceFront(&m2)

	if m1.Len() != 2 {
		t.Errorf("splice into empty list should copy all elements")
	}
}

func TestDListSpliceEmptyListDoesNothing(t *testing.T) {
	e1 := newEmbedListGenerate(2, increment(0))
	e2 := newEmbedList()
	oldFront := e1.Front()
	oldBack := e1.Back()
	oldSize := e1.Len()

	e1.SpliceFront(&e2)

	if e1.Len() != oldSize {
		t.Errorf("splicing empty list should not change size")
	}
	if e1.Front() != oldFront {
		t.Errorf("front should not change when splicing empty list")
	}
	if e1.Back() != oldBack {
		t.Errorf("back should not change when splicing empty list")
	}
	if e2.Len() != 0 {
		t.Errorf("empty list should remain empty")
	}

	// Test member version
	m1 := newMemberListGenerate(2, increment(0))
	m2 := newMemberList()
	oldMSize := m1.Len()

	m1.SpliceFront(&m2)

	if m1.Len() != oldMSize {
		t.Errorf("splicing empty list should not change size")
	}
}

func TestDListSortSingleElement(t *testing.T) {
	e := newEmbedListGenerate(1, increment(0))
	oldElement := e.Front()

	e.Sort(lessEmbed)

	if e.Len() != 1 {
		t.Errorf("sort should not change size of single element list")
	}
	if e.Front() != oldElement {
		t.Errorf("sort should not change single element")
	}
	if e.Back() != oldElement {
		t.Errorf("sort should not change single element")
	}

	// Test member version
	m := newMemberListGenerate(1, increment(0))

	m.Sort(lessMember)

	if m.Len() != 1 {
		t.Errorf("sort should not change size of single element list")
	}
}

func TestDListSortAlreadySorted(t *testing.T) {
	e := newEmbedListGenerate(4, increment(0)) // [0, 1, 2, 3]

	e.Sort(lessEmbed)

	// Should remain sorted
	expected := []int{0, 1, 2, 3}
	current := e.Front()
	for i, exp := range expected {
		if current.value != exp {
			t.Errorf("position %d: expected %d, got %d", i, exp, current.value)
		}
		current = current.Next()
	}

	// Check prev pointers
	current = e.Front()
	for current != nil {
		if current.Next() != nil && current.Next().Prev() != current {
			t.Errorf("prev pointer incorrect for element %d", current.value)
		}
		current = current.Next()
	}

	// Test member version
	m := newMemberListGenerate(4, increment(0))

	m.Sort(lessMember)

	if m.Len() != 4 {
		t.Errorf("sort should not change size")
	}
}

func TestDListReverseSingleElement(t *testing.T) {
	e := newEmbedListGenerate(1, increment(0))
	oldElement := e.Front()

	e.Reverse()

	if e.Len() != 1 {
		t.Errorf("reverse should not change size of single element list")
	}
	if e.Front() != oldElement {
		t.Errorf("reverse should not change single element")
	}
	if e.Back() != oldElement {
		t.Errorf("reverse should not change single element")
	}
	if oldElement.Next() != nil || oldElement.Prev() != nil {
		t.Errorf("single element pointers should remain nil")
	}

	// Test member version
	m := newMemberListGenerate(1, increment(0))

	m.Reverse()

	if m.Len() != 1 {
		t.Errorf("reverse should not change size of single element list")
	}
}

func TestDListRemoveIfAllElements(t *testing.T) {
	e := newEmbedListGenerate(3, increment(0))

	removed := e.RemoveIf(func(value *testEmbedItem) bool {
		return true // Remove all elements
	})

	if len(removed) != 3 {
		t.Errorf("should have removed all elements, got %d", len(removed))
	}
	if e.Len() != 0 {
		t.Errorf("list should be empty after removing all elements")
	}
	if e.Front() != nil {
		t.Errorf("front should be nil after removing all elements")
	}
	if e.Back() != nil {
		t.Errorf("back should be nil after removing all elements")
	}

	// Test member version
	m := newMemberListGenerate(3, increment(0))

	mRemoved := m.RemoveIf(func(value *testMemberItem) bool {
		return true
	})

	if len(mRemoved) != 3 {
		t.Errorf("should have removed all elements, got %d", len(mRemoved))
	}
}

func TestDListRemoveIfNoElements(t *testing.T) {
	e := newEmbedListGenerate(3, increment(0))
	oldSize := e.Len()

	removed := e.RemoveIf(func(value *testEmbedItem) bool {
		return false // Remove no elements
	})

	if len(removed) != 0 {
		t.Errorf("should not have removed any elements, got %d", len(removed))
	}
	if e.Len() != oldSize {
		t.Errorf("list size should not change when no elements removed")
	}

	// Test member version
	m := newMemberListGenerate(3, increment(0))

	mRemoved := m.RemoveIf(func(value *testMemberItem) bool {
		return false
	})

	if len(mRemoved) != 0 {
		t.Errorf("should not have removed any elements, got %d", len(mRemoved))
	}
}

func TestDListUniqueNoDuplicates(t *testing.T) {
	e := newEmbedListGenerate(3, increment(0)) // [0, 1, 2]
	oldSize := e.Len()

	removed := e.Unique(lessEmbed)

	if len(removed) != 0 {
		t.Errorf("should not remove any elements when no duplicates, got %d", len(removed))
	}
	if e.Len() != oldSize {
		t.Errorf("list size should not change when no duplicates")
	}

	// Test member version
	m := newMemberListGenerate(3, increment(0))

	mRemoved := m.Unique(lessMember)

	if len(mRemoved) != 0 {
		t.Errorf("should not remove any elements when no duplicates, got %d", len(mRemoved))
	}
}

func TestDListInitOnNewList(t *testing.T) {
	e := newEmbedList()
	e.Init()
	if !e.Empty() {
		t.Error("Init should make list empty")
	}

	m := newMemberList()
	m.Init()
	if !m.Empty() {
		t.Error("Init should make list empty")
	}
}

func TestDListEmptyOnNonEmptyList(t *testing.T) {
	e := newEmbedListGenerate(1, increment(0))
	if e.Empty() {
		t.Error("Empty should return false for non-empty list")
	}

	m := newMemberListGenerate(1, increment(0))
	if m.Empty() {
		t.Error("Empty should return false for non-empty list")
	}
}

func TestDListSwap(t *testing.T) {
	e1 := newEmbedListGenerate(2, increment(0))
	e2 := newEmbedListGenerate(3, increment(10))

	oldE1Size := e1.Len()
	oldE2Size := e2.Len()

	e1.Swap(&e2)

	if e1.Len() != oldE2Size {
		t.Error("Swap should transfer size from e2 to e1")
	}
	if e2.Len() != oldE1Size {
		t.Error("Swap should transfer size from e1 to e2")
	}

	// Test member version
	m1 := newMemberListGenerate(2, increment(0))
	m2 := newMemberListGenerate(3, increment(10))

	oldM1Size := m1.Len()
	oldM2Size := m2.Len()

	m1.Swap(&m2)

	if m1.Len() != oldM2Size {
		t.Error("Swap should transfer size from m2 to m1")
	}
	if m2.Len() != oldM1Size {
		t.Error("Swap should transfer size from m1 to m2")
	}
}

func TestDListPushFrontNonEmpty(t *testing.T) {
	e := newEmbedListGenerate(1, increment(0)) // List with one element
	el := newEmbed(5)

	e.PushFront(&el)

	if e.Front() != &el {
		t.Error("PushFront should add element to front")
	}
	if e.Front().Next().value != 0 {
		t.Error("Previous front should now be second element")
	}
	if e.Back().value != 0 {
		t.Error("Back should remain the same")
	}

	// Test member version
	m := newMemberListGenerate(1, increment(0))
	ml := newMember(5)

	m.PushFront(&ml)

	if m.Front() != &ml {
		t.Error("PushFront should add element to front")
	}
}

func TestDListPopFrontMultipleElements(t *testing.T) {
	e := newEmbedListGenerate(2, increment(0)) // List with two elements

	popped := e.PopFront()

	if popped.value != 0 {
		t.Error("PopFront should return the front element")
	}
	if e.Front().value != 1 {
		t.Error("Second element should become new front")
	}
	if e.Back().value != 1 {
		t.Error("Back should be updated to the remaining element")
	}
	if e.Len() != 1 {
		t.Error("Size should be decreased by 1")
	}

	// Test member version
	m := newMemberListGenerate(2, increment(0))

	mPopped := m.PopFront()

	if mPopped.value != 0 {
		t.Error("PopFront should return the front element")
	}
}

func TestDListPopBackMultipleElements(t *testing.T) {
	e := newEmbedListGenerate(2, increment(0)) // List with two elements

	popped := e.PopBack()

	if popped.value != 1 {
		t.Error("PopBack should return the back element")
	}
	if e.Front().value != 0 {
		t.Error("Front should remain the same")
	}
	if e.Back().value != 0 {
		t.Error("Front element should become new back")
	}
	if e.Len() != 1 {
		t.Error("Size should be decreased by 1")
	}

	// Test member version
	m := newMemberListGenerate(2, increment(0))

	mPopped := m.PopBack()

	if mPopped.value != 1 {
		t.Error("PopBack should return the back element")
	}
}

func TestDListSizeMethod(t *testing.T) {
	e := newEmbedListGenerate(3, increment(0))

	if e.Size() != 3 {
		t.Error("Size should return the correct number of elements")
	}

	e.PopFront()
	if e.Size() != 2 {
		t.Error("Size should update after operations")
	}

	// Test member version
	m := newMemberListGenerate(3, increment(0))

	if m.Size() != 3 {
		t.Error("Size should return the correct number of elements")
	}
}

func TestDListVerifyNotEmptyPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("verifyNotEmpty should panic on empty list")
		}
	}()

	e := newEmbedList()
	e.verifyNotEmpty()
}

func TestDListVerifyElementNotLinkedPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("verifyElementNotLinked should panic on linked element")
		}
	}()

	e := newEmbedListGenerate(1, increment(0))
	element := e.Front()
	e.verifyElementNotLinked(element)
}

func TestDListVerifyIsMemberOfCurrentPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("verifyIsMemberOfCurrent should panic for non-member element")
		}
	}()

	e := newEmbedList()
	element := newEmbed(0)
	e.verifyIsMemberOfCurrent(&element)
}

func TestDListVerifySizePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("verifySize should panic on size mismatch")
		}
	}()

	e := newEmbedListGenerate(2, increment(0))
	e.size = 3 // Incorrect size
	e.verifySize()
}

func TestDListVerifyNoCyclePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("verifyNoCycle should panic on cycle detection")
		}
	}()

	e := newEmbedListGenerate(2, increment(0))
	// Create a cycle manually
	f, s := e.Front(), e.Back()
	e.hookFunc(f).next = s
	e.hookFunc(s).next = f // Creates cycle
	e.verifyNoCycle()
}
