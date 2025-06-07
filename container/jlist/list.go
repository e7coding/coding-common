package jlist

import (
	"container/list"
	"encoding/json"

	"github.com/e7coding/coding-common/internal/deepcopy"
)

// List 是一个简单的双向链表包装。
type List struct {
	l *list.List
}
type Element = list.Element

// NewList 创建一个空链表。
func NewList() *List {
	return &List{l: list.New()}
}

// NewListFrom 从切片创建链表。
func NewListFrom(vals []interface{}) *List {
	l := list.New()
	for _, v := range vals {
		l.PushBack(v)
	}
	return &List{l: l}
}

// Remove 从链表中删除指定元素 e，并返回它的 Value
func (l *List) Remove(e *Element) interface{} {
	return l.l.Remove(e)
}

// RemoveAt 删除第 i 个元素并返回它的值
func (l *List) RemoveAt(i int) interface{} {
	var target *Element
	l.ForEach(func(e *Element) bool {
		if i == 0 {
			target = e
			return false
		}
		i--
		return true
	})
	if target != nil {
		return l.Remove(target)
	}
	return nil
}

// PushFront 在表头插入一个元素。
func (l *List) PushFront(v interface{}) *list.Element {
	return l.l.PushFront(v)
}

// PushBack 在表尾插入一个元素。
func (l *List) PushBack(v interface{}) *list.Element {
	return l.l.PushBack(v)
}

// PopFront 删除并返回表头元素，空时返回 nil。
func (l *List) PopFront() interface{} {
	if e := l.l.Front(); e != nil {
		return l.l.Remove(e)
	}
	return nil
}

// PopBack 删除并返回表尾元素，空时返回 nil。
func (l *List) PopBack() interface{} {
	if e := l.l.Back(); e != nil {
		return l.l.Remove(e)
	}
	return nil
}

// Len 返回链表长度。
func (l *List) Len() int {
	return l.l.Len()
}

// Empty 判断是否为空。
func (l *List) Empty() bool {
	return l.l.Len() == 0
}

// FrontAll 以切片形式返回从头到尾的所有值。
func (l *List) FrontAll() []interface{} {
	var res []interface{}
	for e := l.l.Front(); e != nil; e = e.Next() {
		res = append(res, e.Value)
	}
	return res
}

// BackAll 以切片形式返回从尾到头的所有值。
func (l *List) BackAll() []interface{} {
	var res []interface{}
	for e := l.l.Back(); e != nil; e = e.Prev() {
		res = append(res, e.Value)
	}
	return res
}

// ForEach 从头到尾遍历，f 返回 false 时中断。
func (l *List) ForEach(f func(e *list.Element) bool) {
	for e := l.l.Front(); e != nil; e = e.Next() {
		if !f(e) {
			break
		}
	}
}

// ForEachReverse 从尾到头遍历，f 返回 false 时中断。
func (l *List) ForEachReverse(f func(e *list.Element) bool) {
	for e := l.l.Back(); e != nil; e = e.Prev() {
		if !f(e) {
			break
		}
	}
}

// Clear 清空链表。
func (l *List) Clear() {
	l.l.Init()
}

// Clone 返回当前链表值的浅拷贝。
func (l *List) Clone() *List {
	return NewListFrom(l.FrontAll())
}

// MarshalJSON 支持 JSON 序列化。
func (l *List) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.FrontAll())
}

// UnmarshalJSON 支持 JSON 反序列化。
func (l *List) UnmarshalJSON(b []byte) error {
	var a []interface{}
	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}
	l.l = list.New()
	for _, v := range a {
		l.l.PushBack(v)
	}
	return nil
}

// DeepCopy 深拷贝链表及其元素。
func (l *List) DeepCopy() interface{} {
	values := l.FrontAll()
	for i, v := range values {
		values[i] = deepcopy.Copy(v)
	}
	return NewListFrom(values)
}
