package jtree

import (
	"fmt"
	"github.com/e7coding/coding-common/jutil"

	"github.com/e7coding/coding-common/internal/json"
	"github.com/emirpasic/gods/trees/redblacktree"
)

// RedBlackTree 是对 gods 红黑树的轻量封装（无锁版本）。
type RedBlackTree struct {
	comparator func(a, b any) int
	tree       *redblacktree.Tree
}

// NewRedBlackTree 创建一棵新的空树，comparator 不能为空。
func NewRedBlackTree(cmp func(a, b any) int) *RedBlackTree {
	return &RedBlackTree{
		comparator: cmp,
		tree:       redblacktree.NewWith(cmp),
	}
}
func NewRedBlackTreeFrom(comparator func(v1, v2 any) int, data map[any]any) *RedBlackTree {
	tree := NewRedBlackTree(comparator)
	for k, v := range data {
		tree.Set(k, v)
	}
	return tree
}

// Clone 深拷贝当前树（不含 lock）。
func (t *RedBlackTree) Clone() *RedBlackTree {
	clone := NewRedBlackTree(t.comparator)
	for _, key := range t.tree.Keys() {
		val, _ := t.tree.Get(key)
		clone.tree.Put(key, val)
	}
	return clone
}

// Set 插入或更新一对 key/value。
func (t *RedBlackTree) Set(key, value any) {
	t.tree.Put(key, value)
}

// Get 返回 key 对应的 value；若不存在则返回 nil。
func (t *RedBlackTree) Get(key any) any {
	v, _ := t.tree.Get(key)
	return v
}

// Remove 删除 key 并返回其旧值；不存在则返回 nil。
func (t *RedBlackTree) Remove(key any) any {
	old, _ := t.tree.Get(key)
	t.tree.Remove(key)
	return old
}

// Contains 判断 key 是否存在。
func (t *RedBlackTree) Contains(key any) bool {
	_, ok := t.tree.Get(key)
	return ok
}

// Size 返回节点总数。
func (t *RedBlackTree) Size() int {
	return t.tree.Size()
}

// Clear 清空整棵树。
func (t *RedBlackTree) Clear() {
	t.tree.Clear()
}

// Keys 按照 comparator 排序返回所有 key。
func (t *RedBlackTree) Keys() []any {
	return t.tree.Keys()
}

// Values 按照 comparator 排序返回所有 value。
func (t *RedBlackTree) Values() []any {
	return t.tree.Values()
}

// ForEachAsc 从最小到最大遍历，f 返回 false 则中断。
func (t *RedBlackTree) ForEachAsc(f func(key, value any) bool) {
	it := t.tree.Iterator()
	for it.Begin(); it.Next(); {
		if !f(it.Key(), it.Value()) {
			break
		}
	}
}

// ForEachDesc 从最大到最小遍历，f 返回 false 则中断。
func (t *RedBlackTree) ForEachDesc(f func(key, value any) bool) {
	it := t.tree.Iterator()
	for it.End(); it.Prev(); {
		if !f(it.Key(), it.Value()) {
			break
		}
	}
}

// GetOrSet 如果 key 存在则直接返回旧值，否则设置并返回 value。
func (t *RedBlackTree) GetOrSet(key, value any) any {
	if v, ok := t.tree.Get(key); ok {
		return v
	}
	t.tree.Put(key, value)
	return value
}

// GetOrSetFunc 如果 key 不存在则调用 f() 得到值并插入，然后返回。
func (t *RedBlackTree) GetOrSetFunc(key any, f func() any) any {
	if v, ok := t.tree.Get(key); ok {
		return v
	}
	v := f()
	t.tree.Put(key, v)
	return v
}

// ToMap 按序返回一个 map（注意，map 本身无序）。
func (t *RedBlackTree) ToMap() map[any]any {
	m := make(map[any]any, t.tree.Size())
	t.ForEachAsc(func(k, v any) bool {
		m[k] = v
		return true
	})
	return m
}

// MarshalJSON 直接利用内置实现。
func (t *RedBlackTree) MarshalJSON() ([]byte, error) {
	return t.tree.MarshalJSON()
}

// UnmarshalJSON 从 JSON 恢复（假设 key 可转为 string）。
func (t *RedBlackTree) UnmarshalJSON(b []byte) error {
	// 默认使用字符串比较器
	if t.comparator == nil {
		t.comparator = jutil.ComparatorString
		t.tree = redblacktree.NewWith(t.comparator)
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

// String 返回类似 "{k1:v1 k2:v2}" 的简洁表示。
func (t *RedBlackTree) String() string {
	return fmt.Sprintf("%v", t.ToMap())
}

// Flip 交换 key<->value 并可选传入新的 comparator。
func (t *RedBlackTree) Flip(cmp ...func(a, b any) int) {
	newCmp := t.comparator
	if len(cmp) > 0 {
		newCmp = cmp[0]
	}
	flipped := NewRedBlackTree(newCmp)
	t.ForEachAsc(func(k, v any) bool {
		flipped.tree.Put(v, k)
		return true
	})
	*t = *flipped
}

// Floor 返回小于等于 key 的最大键值对（如果存在）。
func (t *RedBlackTree) Floor(key any) (k, v any, ok bool) {
	if node, found := t.tree.Floor(key); found {
		return node.Key, node.Value, true
	}
	return nil, nil, false
}

// Ceiling 返回大于等于 key 的最小键值对（如果存在）。
func (t *RedBlackTree) Ceiling(key any) (k, v any, ok bool) {
	if node, found := t.tree.Ceiling(key); found {
		return node.Key, node.Value, true
	}
	return nil, nil, false
}
