package jtree

import (
	"fmt"
	"github.com/e7coding/coding-common/text/jstr"

	"github.com/e7coding/coding-common/jutil/jconv"
	"github.com/emirpasic/gods/trees/btree"
)

// BTree implements a B-tree without internal locking.
type BTree struct {
	order int
	cmp   func(a, b any) int
	tree  *btree.Tree
}

// BEntry is a key/value pair.
type BEntry struct {
	Key, Value any
}

// NewBTree creates a B-tree of order m (>=3) with cmp.
func NewBTree(m int, cmp func(a, b any) int) *BTree {
	return &BTree{order: m, cmp: cmp, tree: btree.NewWith(m, cmp)}
}

// NewBTreeFrom creates a B-tree prefilled from data.
func NewBTreeFrom(m int, cmp func(a, b any) int, data map[any]any) *BTree {
	t := NewBTree(m, cmp)
	for k, v := range data {
		t.Put(k, v)
	}
	return t
}

// Clone makes a shallow copy.
func (t *BTree) Clone() *BTree {
	nt := NewBTree(t.order, t.cmp)
	nt.Range(func(k, v any) bool {
		nt.Put(k, v)
		return true
	})
	return nt
}

// Put inserts or updates.
func (t *BTree) Put(key, val any) {
	t.tree.Put(key, val)
}

// PutIfAbsent only puts when key missing.
func (t *BTree) PutIfAbsent(key, val any) bool {
	if _, ok := t.tree.Get(key); !ok {
		t.tree.Put(key, val)
		return true
	}
	return false
}

// Get returns value or nil.
func (t *BTree) Get(key any) any {
	if v, _ := t.tree.Get(key); v != nil {
		return v
	}
	return nil
}

// GetOK returns (value, found).
func (t *BTree) GetOK(key any) (any, bool) {
	return t.tree.Get(key)
}

// GetOrPut returns existing or puts the new value.
func (t *BTree) GetOrPut(key, val any) any {
	if v, ok := t.tree.Get(key); ok {
		return v
	}
	t.tree.Put(key, val)
	return val
}

// Has reports whether key exists.
func (t *BTree) Has(key any) bool {
	_, ok := t.tree.Get(key)
	return ok
}

// Delete removes key, returns old value.
func (t *BTree) Delete(key any) any {
	if v, ok := t.tree.Get(key); ok {
		t.tree.Remove(key)
		return v
	}
	return nil
}

func foreachGetIndex(key any, keys []any, match bool) (index int, canIterator bool) {
	if match {
		for i, k := range keys {
			if k == key {
				canIterator = true
				index = i
			}
		}
	} else {
		if i, ok := key.(int); ok {
			canIterator = true
			index = i
		}
	}
	return
}

func (t *BTree) ForeachFrom(key any, match bool, f func(key, value any) bool) {
	var keys = t.tree.Keys()
	index, canIterator := foreachGetIndex(key, keys, match)
	if !canIterator {
		return
	}
	for ; index < len(keys); index++ {
		f(keys[index], t.Get(keys[index]))
	}
}
func (t *BTree) Foreach(f func(key, value any) bool) {
	var (
		ok bool
		it = t.tree.Iterator()
	)
	for it.Begin(); it.Next(); {
		index, value := it.Key(), it.Value()
		if ok = f(index, value); !ok {
			break
		}
	}
}
func (t *BTree) ForeachReverse(f func(key, value any) bool) {
	var (
		ok bool
		it = t.tree.Iterator()
	)
	for it.End(); it.Prev(); {
		index, value := it.Key(), it.Value()
		if ok = f(index, value); !ok {
			break
		}
	}
}

// Clear removes all nodes.
func (t *BTree) Clear() {
	t.tree.Clear()
}

// Size returns number of elements.
func (t *BTree) Size() int {
	return t.tree.Size()
}

// Empty reports whether tree is empty.
func (t *BTree) Empty() bool {
	return t.tree.Size() == 0
}

// Keys returns all keys in order.
func (t *BTree) Keys() []any {
	return t.tree.Keys()
}

// Values returns all values in order.
func (t *BTree) Values() []any {
	return t.tree.Values()
}

// Map builds a snapshot map[key]=value.
func (t *BTree) Map() map[any]any {
	m := make(map[any]any, t.Size())
	t.Range(func(k, v any) bool {
		m[k] = v
		return true
	})
	return m
}

// MapStr returns map[string]any.
func (t *BTree) MapStr() map[string]any {
	m := make(map[string]any, t.Size())
	t.Range(func(k, v any) bool {
		m[jconv.String(k)] = v
		return true
	})
	return m
}

// Range calls fn for each kv in ascending order, stop if fn returns false.
func (t *BTree) Range(fn func(key, value any) bool) {
	it := t.tree.Iterator()
	for it.Begin(); it.Next(); {
		if !fn(it.Key(), it.Value()) {
			break
		}
	}
}

// RangeRev is like Range but in descending order.
func (t *BTree) RangeRev(fn func(key, value any) bool) {
	it := t.tree.Iterator()
	for it.End(); it.Prev(); {
		if !fn(it.Key(), it.Value()) {
			break
		}
	}
}

// Print for debugging.
func (t *BTree) Print() {
	fmt.Println(t.String())
}

// String returns a compact dump.
func (t *BTree) String() string {
	return jstr.Replace(t.tree.String(), "BTree\n", "")
}

// MarshalJSON delegates.
func (t *BTree) MarshalJSON() ([]byte, error) {
	return t.tree.MarshalJSON()
}
