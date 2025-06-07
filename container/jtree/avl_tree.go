package jtree

import (
	"fmt"
	"github.com/e7coding/coding-common/jutil"

	"github.com/e7coding/coding-common/container/jvar"
	"github.com/e7coding/coding-common/internal/json"
	"github.com/emirpasic/gods/trees/avltree"
)

// AVLTree 是对 gods AVL 树的轻量封装（无锁版本）。
type AVLTree struct {
	comparator func(a, b any) int
	tree       *avltree.Tree
}

// NewAVLTree 创建空 AVL 树。
func NewAVLTree(cmp func(a, b any) int) *AVLTree {
	return &AVLTree{comparator: cmp, tree: avltree.NewWith(cmp)}
}

// Clone 返回当前树的深拷贝。
func (t *AVLTree) Clone() *AVLTree {
	clone := NewAVLTree(t.comparator)
	for _, k := range t.tree.Keys() {
		v, _ := t.tree.Get(k)
		clone.tree.Put(k, v)
	}
	return clone
}

// Set 插入或更新 (key, value)。
func (t *AVLTree) Set(key, value any) {
	if f, ok := value.(func() any); ok {
		value = f()
	}
	if value != nil {
		t.tree.Put(key, value)
	}
}

// Get 返回 key 对应的 value，不存在时返回 nil。
func (t *AVLTree) Get(key any) any {
	v, _ := t.tree.Get(key)
	return v
}

// Remove 删除 key，返回旧值。
func (t *AVLTree) Remove(key any) any {
	v, _ := t.tree.Get(key)
	t.tree.Remove(key)
	return v
}

// Contains 判断 key 是否存在。
func (t *AVLTree) Contains(key any) bool {
	_, ok := t.tree.Get(key)
	return ok
}

// Size 返回节点数量。
func (t *AVLTree) Size() int { return t.tree.Size() }

// Clear 清空树。
func (t *AVLTree) Clear() { t.tree.Clear() }

// Keys 返回所有键（已排序）。
func (t *AVLTree) Keys() []any { return t.tree.Keys() }

// Values 返回所有值（已排序）。
func (t *AVLTree) Values() []any { return t.tree.Values() }

// ForEachAsc 按升序遍历，f 返回 false 则中断。
func (t *AVLTree) ForEachAsc(f func(key, value any) bool) {
	it := t.tree.Iterator()
	for it.Begin(); it.Next(); {
		if !f(it.Key(), it.Value()) {
			break
		}
	}
}

// ForEachDesc 按降序遍历，f 返回 false 则中断。
func (t *AVLTree) ForEachDesc(f func(key, value any) bool) {
	it := t.tree.Iterator()
	for it.End(); it.Prev(); {
		if !f(it.Key(), it.Value()) {
			break
		}
	}
}

// GetOrSet 如果 key 存在则返回旧值，否则设置并返回 value。
func (t *AVLTree) GetOrSet(key, value any) any {
	if v, ok := t.tree.Get(key); ok {
		return v
	}
	t.tree.Put(key, value)
	return value
}

// GetOrSetFunc 如果 key 不存在则执行 f() 得到值并设置。
func (t *AVLTree) GetOrSetFunc(key any, f func() any) any {
	if v, ok := t.tree.Get(key); ok {
		return v
	}
	val := f()
	t.tree.Put(key, val)
	return val
}

// ToMap 按序返回 map。注意 map 本身无序。
func (t *AVLTree) ToMap() map[any]any {
	m := make(map[any]any, t.tree.Size())
	t.ForEachAsc(func(k, v any) bool { m[k] = v; return true })
	return m
}

// String 返回类似 "{k1:v1 k2:v2}"。
func (t *AVLTree) String() string { return fmt.Sprintf("%v", t.ToMap()) }

// MarshalJSON 直接复用内置。
func (t *AVLTree) MarshalJSON() ([]byte, error) { return t.tree.MarshalJSON() }

// UnmarshalJSON 假定 key 为 string。
func (t *AVLTree) UnmarshalJSON(b []byte) error {
	if t.comparator == nil {
		t.comparator = jutil.ComparatorString
		t.tree = avltree.NewWith(t.comparator)
	}
	var data map[string]any
	if err := json.UnmarshalUseNumber(b, &data); err != nil {
		return err
	}
	for k, v := range data {
		t.tree.Put(k, v)
	}
	return nil
}

// GetVarOrSet 和 GetVarOrSetFunc 返回 jvar.Var。
func (t *AVLTree) GetVar(key any) *jvar.Var {
	return jvar.New(t.Get(key))
}
func (t *AVLTree) GetVarOrSet(key, v any) *jvar.Var {
	return jvar.New(t.GetOrSet(key, v))
}
func (t *AVLTree) GetVarFunc(key any, f func() any) *jvar.Var {
	return jvar.New(t.GetOrSetFunc(key, f))
}
