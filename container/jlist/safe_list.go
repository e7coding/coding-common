// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with l file,
// You can obtain one at https://github.com/gogf/gf.
//

package jlist

import (
	"bytes"
	"container/list"
	"encoding/json"
	"errors"
	"github.com/e7coding/coding-common/internal/deepcopy"
	"github.com/e7coding/coding-common/jutil/jconv"
	"sync"
)

// SafeList 双向链表，支持并发安全
// 默认线程安全
// 提供基本增删查遍历功能
type SafeList struct {
	mu   sync.RWMutex
	list *list.List
}

// NewSafeList 创建空列表
func NewSafeList() *SafeList {
	return &SafeList{list: list.New()}
}

// NewSafeListFrom 根据切片创建列表
func NewSafeListFrom(vals []interface{}) *SafeList {
	l := NewSafeList()
	sz := len(vals)
	if sz > 0 {
		l.mu.Lock()
		for _, v := range vals {
			l.list.PushBack(v)
		}
		l.mu.Unlock()
	}
	return l
}

// Len 长度
func (l *SafeList) Len() int {
	l.mu.RLock()
	n := l.list.Len()
	l.mu.RUnlock()
	return n
}

// FrontValue 首元素值，空列表返回 nil
func (l *SafeList) FrontValue() interface{} {
	l.mu.RLock()
	e := l.list.Front()
	l.mu.RUnlock()
	if e != nil {
		return e.Value
	}
	return nil
}

// BackValue 尾元素值，空列表返回 nil
func (l *SafeList) BackValue() interface{} {
	l.mu.RLock()
	e := l.list.Back()
	l.mu.RUnlock()
	if e != nil {
		return e.Value
	}
	return nil
}

// PushFront 在头部插入元素
func (l *SafeList) PushFront(v interface{}) *list.Element {
	l.mu.Lock()
	if l.list == nil {
		l.list = list.New()
	}
	e := l.list.PushFront(v)
	l.mu.Unlock()
	return e
}

// PushBack 在尾部插入元素
func (l *SafeList) PushBack(v interface{}) *list.Element {
	l.mu.Lock()
	if l.list == nil {
		l.list = list.New()
	}
	e := l.list.PushBack(v)
	l.mu.Unlock()
	return e
}

// PopFront 弹出头部元素
func (l *SafeList) PopFront() (interface{}, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list.Len() == 0 {
		return nil, errors.New("empty list")
	}
	e := l.list.Front()
	v := l.list.Remove(e)
	return v, nil
}
func (l *SafeList) Remove(e *Element) (value interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
	}
	value = l.list.Remove(e)
	return
}

// PopBack 弹出尾部元素
func (l *SafeList) PopBack() (interface{}, error) {
	l.mu.Lock()
	if l.list.Len() == 0 {
		l.mu.Unlock()
		return nil, errors.New("empty list")
	}
	e := l.list.Back()
	v := l.list.Remove(e)
	l.mu.Unlock()
	return v, nil
}

// ForEach 从头遍历，回调返回 false 停止
func (l *SafeList) ForEach(fn func(v interface{}) bool) {
	l.mu.RLock()
	e := l.list.Front()
	for e != nil {
		v := e.Value
		l.mu.RUnlock()
		if !fn(v) {
			return
		}
		l.mu.RLock()
		e = e.Next()
	}
	l.mu.RUnlock()
}

// ForEachReverse 从尾遍历，回调返回 false 停止
func (l *SafeList) ForEachReverse(fn func(v interface{}) bool) {
	l.mu.RLock()
	e := l.list.Back()
	for e != nil {
		v := e.Value
		l.mu.RUnlock()
		if !fn(v) {
			return
		}
		l.mu.RLock()
		e = e.Prev()
	}
	l.mu.RUnlock()
}

// Join 拼接所有元素为字符串
func (l *SafeList) Join(sep string) string {
	l.mu.RLock()
	buf := &bytes.Buffer{}
	e := l.list.Front()
	first := true
	for e != nil {
		if !first {
			buf.WriteString(sep)
		}
		buf.WriteString(jconv.String(e.Value))
		first = false
		e = e.Next()
	}
	l.mu.RUnlock()
	return buf.String()
}

// Clear 清空列表
func (l *SafeList) Clear() {
	l.mu.Lock()
	l.list = list.New()
	l.mu.Unlock()
}

// Clone 深拷贝列表
func (l *SafeList) Clone() *SafeList {
	l.mu.RLock()
	sz := l.list.Len()
	vals := make([]interface{}, 0, sz)
	e := l.list.Front()
	for e != nil {
		vals = append(vals, deepcopy.Copy(e.Value))
		e = e.Next()
	}
	l.mu.RUnlock()
	return NewSafeListFrom(vals)
}

// MarshalJSON JSON 序列化
func (l *SafeList) MarshalJSON() ([]byte, error) {
	l.mu.RLock()
	vals := make([]interface{}, 0, l.list.Len())
	e := l.list.Front()
	for e != nil {
		vals = append(vals, e.Value)
		e = e.Next()
	}
	l.mu.RUnlock()
	return json.Marshal(vals)
}

// UnmarshalJSON JSON 反序列化
func (l *SafeList) UnmarshalJSON(b []byte) error {
	l.mu.Lock()
	var vals []interface{}
	if err := json.Unmarshal(b, &vals); err != nil {
		l.mu.Unlock()
		return err
	}
	l.list = list.New()
	for _, v := range vals {
		l.list.PushBack(v)
	}
	l.mu.Unlock()
	return nil
}
