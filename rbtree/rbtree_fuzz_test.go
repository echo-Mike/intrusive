package rbtree

import (
	"encoding/binary"
	"testing"
)

type fuzzEmbedItem struct {
	Hook[fuzzEmbedItem]
	value     int
	isUsed    bool
	treeIndex int
	id        int
}

func fuzzEmbedHook(self *fuzzEmbedItem) *Hook[fuzzEmbedItem] {
	return &self.Hook
}

func newFuzzRbTree() *RbTree[fuzzEmbedItem] {
	return NewRbTree(fuzzEmbedHook, lessFuzz)
}

func newFuzz(value, id int) fuzzEmbedItem {
	return fuzzEmbedItem{Hook: NewHook[fuzzEmbedItem](), value: value, isUsed: false, treeIndex: 0, id: id}
}

func lessFuzz(lhs, rhs *fuzzEmbedItem) bool {
	return lhs.value < rhs.value
}

const (
	opInsert byte = iota
	opErase
	opClear
	opFind
	opLowerBound
	opUpperBound
	opMerge
	opIncludes
	opDifference
	opIntersection
	opSymDifference
	opUnion
	opEraseIf
	opVerifyTree
	opSize
	opEmpty
	opFront
	opBack
	opSwap
	opInit
	opCOUNT
)

func verifyTreeConsistency(t *testing.T, tree *RbTree[fuzzEmbedItem], treeIdx int) {
	if tree.Empty() {
		if tree.Front() != nil || tree.Back() != nil || tree.Size() != 0 {
			t.Errorf("Empty tree inconsistency: front=%v, back=%v, size=%d", tree.Front(), tree.Back(), tree.Size())
		}
		return
	}

	// Verify tree properties
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Tree verification failed: %v", r)
		}
	}()
	tree.verify()

	// Verify size matches in-order traversal count
	count := 0
	current := tree.Front()
	prev := current
	for current != nil {
		count++
		if prev != current && tree.Prev(current) != prev {
			t.Errorf("Prev pointer inconsistency at node %v", current)
		}
		prev = current
		current = tree.Next(current)
	}

	if count != tree.Size() {
		t.Errorf("Size inconsistency: in-order=%d, stored=%d", count, tree.Size())
	}

	// Verify first and last pointers
	if tree.min(tree.root) != tree.Front() {
		t.Errorf("First pointer inconsistency")
	}

	if tree.max(tree.root) != tree.Back() {
		t.Errorf("Last pointer inconsistency")
	}

	tree.TraversePreOrder(func(node *fuzzEmbedItem) {
		// Verify all nodes have proper parent pointers
		if tree.left(node) != nil && tree.parent(tree.left(node)) != node {
			t.Errorf("Left child parent pointer mismatch at node %v", node)
		}
		if tree.right(node) != nil && tree.parent(tree.right(node)) != node {
			t.Errorf("Right child parent pointer mismatch at node %v", node)
		}
		// Verify that all nodes are from this tree
		if node.treeIndex != treeIdx {
			t.Errorf("Node %v thinks it's from other tree", node)
		}
	})
}

func referenceIncludes[T any](t1, t2 *RbTree[T]) bool {
	for node := t2.Front(); node != nil; node = t2.Next(node) {
		if t1.Find(node) == nil {
			return false
		}
	}
	return true
}

func referenceDifference[T any](t1, t2 *RbTree[T]) []*T {
	var result []*T
	for node := t1.Front(); node != nil; node = t1.Next(node) {
		if t2.Find(node) == nil {
			result = append(result, node)
		}
	}
	return result
}

func referenceIntersection[T any](t1, t2 *RbTree[T]) []*T {
	var result []*T
	for node := t1.Front(); node != nil; node = t1.Next(node) {
		if t2.Find(node) != nil {
			result = append(result, node)
		}
	}
	return result
}

func referenceSymDifference[T any](t1, t2 *RbTree[T]) []*T {
	diff1 := referenceDifference(t1, t2)
	diff2 := referenceDifference(t2, t1)
	return append(diff1, diff2...)
}

func referenceUnion[T any](t1, t2 *RbTree[T]) []*T {
	var result []*T
	addUnique := func(node *T) {
		for _, n := range result {
			if n == node {
				return
			}
		}
		result = append(result, node)
	}

	for node := t1.Front(); node != nil; node = t1.Next(node) {
		addUnique(node)
	}
	for node := t2.Front(); node != nil; node = t2.Next(node) {
		addUnique(node)
	}
	return result
}

func referenceFind[T any](tree *RbTree[T], item *T) *T {
	for node := tree.Front(); node != nil; node = tree.Next(node) {
		if !tree.lessFunc(node, item) && !tree.lessFunc(item, node) {
			return node
		}
	}
	return nil
}

func referenceLowerBound[T any](tree *RbTree[T], item *T) *T {
	var candidate *T
	for node := tree.Front(); node != nil; node = tree.Next(node) {
		if !tree.lessFunc(node, item) && (candidate == nil || tree.lessFunc(node, candidate)) {
			candidate = node
		}
	}
	return candidate
}

func referenceUpperBound[T any](tree *RbTree[T], item *T) *T {
	var candidate *T
	for node := tree.Front(); node != nil; node = tree.Next(node) {
		if tree.lessFunc(item, node) && (candidate == nil || tree.lessFunc(node, candidate)) {
			candidate = node
		}
	}
	return candidate
}

func compareSlices[T any](a, b []*T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func nextState(t *testing.T, items []fuzzEmbedItem, trees []*RbTree[fuzzEmbedItem]) func(op, arg1, arg2, arg3 byte, arg4 uint32) {
	return func(op, arg1, arg2, arg3 byte, arg4 uint32) {
		treeIdx := int(arg1) % len(trees)
		itemIdx := int(arg4) % len(items)
		tree2Idx := int(arg3) % len(trees)

		tree := trees[treeIdx]
		item := &items[itemIdx]
		tree2 := trees[tree2Idx]

		switch op % opCOUNT {
		case opInsert:
			if !item.isUsed {
				if tree.Insert(item) {
					item.isUsed = true
					item.treeIndex = treeIdx
				} else if !tree.Contains(item) {
					t.Errorf("Failed to insert item %v", item)
				}
			}

		case opErase:
			if item.isUsed && item.treeIndex == treeIdx {
				if tree.Erase(item) {
					item.isUsed = false
					item.treeIndex = 0
				} else {
					t.Errorf("Failed to erase item %v", item)
				}
			}

		case opClear:
			cleared := tree.Clear()
			for _, it := range cleared {
				it.isUsed = false
				it.treeIndex = 0
			}

		case opFind:
			expected := referenceFind(tree, item)
			actual := tree.Find(item)
			if expected != actual {
				t.Errorf("Find mismatch: expected %v, got %v", expected, actual)
			}

		case opLowerBound:
			expected := referenceLowerBound(tree, item)
			actual := tree.LowerBound(item)
			if expected != actual {
				t.Errorf("LowerBound mismatch: expected %v, got %v", expected, actual)
			}

		case opUpperBound:
			expected := referenceUpperBound(tree, item)
			actual := tree.UpperBound(item)
			if expected != actual {
				t.Errorf("UpperBound mismatch: expected %v, got %v", expected, actual)
			}

		case opMerge:
			if tree != tree2 {
				originalSize := tree.Size() + tree2.Size()
				tree.Merge(tree2)
				if tree.Size() != originalSize {
					t.Errorf("Merge size inconsistency: expected %d, got %d", originalSize, tree.Size())
				}
				tree.Traverse(func(node *fuzzEmbedItem) {
					node.treeIndex = treeIdx
				})
			}

		case opIncludes:
			if tree != tree2 {
				expected := referenceIncludes(tree, tree2)
				actual := tree.Includes(tree2)
				if expected != actual {
					t.Errorf("Includes mismatch: expected %v, got %v", expected, actual)
				}
			}

		case opDifference:
			if tree != tree2 {
				expected := referenceDifference(tree, tree2)
				actual := tree.Difference(tree2)
				if !compareSlices(expected, actual) {
					t.Errorf("Difference mismatch")
				}
			}

		case opIntersection:
			if tree != tree2 {
				expected := referenceIntersection(tree, tree2)
				actual := tree.Intersection(tree2)
				if !compareSlices(expected, actual) {
					t.Errorf("Intersection mismatch")
				}
			}

		case opSymDifference:
			if tree != tree2 {
				expected := referenceSymDifference(tree, tree2)
				actual := tree.SymDifference(tree2)
				if !compareSlices(expected, actual) {
					t.Errorf("SymDifference mismatch")
				}
			}

		case opUnion:
			if tree != tree2 {
				expected := referenceUnion(tree, tree2)
				actual := tree.Union(tree2)
				if !compareSlices(expected, actual) {
					t.Errorf("Union mismatch")
				}
			}

		case opEraseIf:
			var predicate func(*fuzzEmbedItem) bool
			switch int(arg2) % 8 {
			case 0:
				predicate = func(e *fuzzEmbedItem) bool { return e.value%2 == 0 }
			case 1:
				predicate = func(e *fuzzEmbedItem) bool { return e.value%2 == 1 }
			case 2:
				predicate = func(e *fuzzEmbedItem) bool { return e.value < 8 }
			case 3:
				predicate = func(e *fuzzEmbedItem) bool { return e.value >= 8 }
			case 4:
				predicate = func(e *fuzzEmbedItem) bool { return e.id%3 == 0 }
			case 5:
				predicate = func(e *fuzzEmbedItem) bool { return e.id%5 == 0 }
			case 6:
				predicate = func(e *fuzzEmbedItem) bool { return e.value == 0 }
			case 7:
				predicate = func(e *fuzzEmbedItem) bool { return e.value > 1000 }
			}

			erased := tree.EraseIf(predicate)
			for _, e := range erased {
				e.isUsed = false
				e.treeIndex = 0
			}

		case opVerifyTree:
			verifyTreeConsistency(t, tree, treeIdx)

		case opSize:
			if tree.Size() < 0 {
				t.Errorf("Negative tree size: %d", tree.Size())
			}

		case opEmpty:
			empty := tree.Empty()
			size := tree.Size()

			// Check consistency between Empty() and Size()
			if empty && size != 0 {
				t.Errorf("Empty() returned true but Size() returned %d", size)
			}
			if !empty && size == 0 {
				t.Errorf("Empty() returned false but Size() returned 0")
			}

			// Check that front and back pointers are nil when empty
			if empty {
				if tree.Front() != nil {
					t.Errorf("Empty tree has non-nil Front(): %v", tree.Front())
				}
				if tree.Back() != nil {
					t.Errorf("Empty tree has non-nil Back(): %v", tree.Back())
				}
				if tree.root != nil {
					t.Errorf("Empty tree has non-nil root: %v", tree.root)
				}
			} else {
				// For non-empty trees, verify front and back are not nil
				if tree.Front() == nil {
					t.Errorf("Non-empty tree has nil Front()")
				}
				if tree.Back() == nil {
					t.Errorf("Non-empty tree has nil Back()")
				}
				if tree.root == nil {
					t.Errorf("Non-empty tree has nil root")
				}
			}

			for range 3 {
				if tree.Empty() != empty {
					t.Errorf("Empty() returned inconsistent results: expected %v, got %v", empty, tree.Empty())
				}
			}

		case opFront:
			if front := tree.Front(); front != nil && tree.Prev(front) != nil {
				t.Errorf("Front element has predecessor: %v", front)
			}

		case opBack:
			if back := tree.Back(); back != nil && tree.Next(back) != nil {
				t.Errorf("Back element has successor: %v", back)
			}

		case opSwap:
			if tree != tree2 {
				size1, size2 := tree.Size(), tree2.Size()
				tree.Swap(tree2)
				if tree.Size() != size2 || tree2.Size() != size1 {
					t.Errorf("Swap size inconsistency")
				}
				tree.Traverse(func(node *fuzzEmbedItem) {
					node.treeIndex = treeIdx
				})
				tree2.Traverse(func(node *fuzzEmbedItem) {
					node.treeIndex = tree2Idx
				})
			}

		case opInit:
			if cleared := tree.Clear(); len(cleared) > 0 {
				for _, it := range cleared {
					it.isUsed = false
					it.treeIndex = 0
				}
			}
			tree.Init()
		}
	}
}

func FuzzRbTreeOps(f *testing.F) {
	const numItems = 2048
	const numTrees = 16

	items := make([]fuzzEmbedItem, numItems)
	for i := range items {
		items[i] = newFuzz(i%64, i)
	}

	trees := make([]*RbTree[fuzzEmbedItem], numTrees)
	for i := range trees {
		trees[i] = newFuzzRbTree()
	}

	f.Fuzz(func(t *testing.T, commands []byte) {
		// Reset all trees and items
		for i := range trees {
			if elements := trees[i].Clear(); len(elements) > 0 {
				for _, e := range elements {
					e.isUsed = false
					e.treeIndex = 0
				}
			}
			trees[i].Init()
		}

		for i := range items {
			items[i].isUsed = false
			items[i].treeIndex = 0
			items[i].Hook.Init()
		}

		next := nextState(t, items, trees)

		for i := 0; i+7 < len(commands); i += 8 {
			indexArg := binary.LittleEndian.Uint32([]byte{commands[i+4], commands[i+5], commands[i+6], commands[i+7]})
			next(commands[i], commands[i+1], commands[i+2], commands[i+3], indexArg)
		}

		// Verify all trees at the end
		for i := range trees {
			verifyTreeConsistency(t, trees[i], i)
		}

		// Additional verification: check that all items are either in exactly one tree or not in any tree
		inTreeCount := 0
		for i := range items {
			count := 0
			for _, tree := range trees {
				if tree.Find(&items[i]) == &items[i] {
					count++
				}
			}
			if count > 1 {
				t.Errorf("Item %v found in multiple trees", items[i])
			}
			if items[i].isUsed && count == 0 {
				t.Errorf("Item %v marked as used but not found in any tree", items[i])
			}
			if !items[i].isUsed && count > 0 {
				t.Errorf("Item %v not marked as used but found in tree", items[i])
			}
			inTreeCount += count
		}

		// Verify total items in trees matches sum of sizes
		totalSize := 0
		for _, tree := range trees {
			totalSize += tree.Size()
		}
		if totalSize != inTreeCount {
			t.Errorf("Total size mismatch: sum of sizes=%d, actual items in trees=%d", totalSize, inTreeCount)
		}
	})
}
