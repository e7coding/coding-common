package jpool

import (
	"context"
	"github.com/e7coding/coding-common/errs/jcode"
	"github.com/e7coding/coding-common/errs/jerr"
	"sync"
	"time"

	"github.com/e7coding/coding-common/container/jlist"
)

// ErrPoolEmpty 当池中无可用对象且 NewFunc 也未设置时返回。
var ErrPoolEmpty = jerr.WithCode(jcode.NewErrCode(jcode.OptErr), "pool is empty")

// Pool 是一个对象可复用的简单池。
type Pool struct {
	idle       *jlist.SafeList // 空闲对象列表
	ttl        time.Duration   // 对象存活时间
	newFunc    func() (interface{}, error)
	expireFunc func(interface{})

	mu           sync.Mutex
	closed       bool
	cancelReaper context.CancelFunc
}

type poolItem struct {
	obj      interface{}
	expireAt int64 // 毫秒级时间戳
}

// New 创建一个对象池：
//
//	ttl == 0 永不过期
//	ttl < 0 用过即过期
//	ttl > 0 存活 ttl 后过期
func New(
	ttl time.Duration,
	newFunc func() (interface{}, error),
	expireFunc ...func(interface{}),
) *Pool {
	p := &Pool{
		idle:    jlist.NewSafeList(),
		ttl:     ttl,
		newFunc: newFunc,
	}
	if len(expireFunc) > 0 {
		p.expireFunc = expireFunc[0]
	}
	// 启动后台清理器
	var ctx context.Context
	ctx, p.cancelReaper = context.WithCancel(context.Background())
	go p.runReaper(ctx)
	return p
}

// Put 放回一个对象到池中，可能会被 later 过期回收。
func (p *Pool) Put(obj interface{}) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return jerr.WithCode(jcode.NewErrCode(jcode.OptErr), "pool closed")
	}
	item := &poolItem{obj, 0}
	if p.ttl > 0 {
		item.expireAt = time.Now().Add(p.ttl).UnixMilli()
	}
	p.idle.PushBack(item)
	return nil
}

// Get 从池中取一个对象：
//  1. 如果有未过期的，返回之；
//  2. 否则尝试 newFunc；
//  3. 还没拿到，则 ErrPoolEmpty。
func (p *Pool) Get() (interface{}, error) {
	// 先从空闲列表中扫一遍
	for {
		r, _ := p.idle.PopFront()
		if r == nil {
			break
		}
		item := r.(*poolItem)
		// never expire
		if p.ttl == 0 || item.expireAt > time.Now().UnixMilli() {
			return item.obj, nil
		}
		// 过期：调用回调
		if p.expireFunc != nil {
			p.expireFunc(item.obj)
		}
	}
	// 池空，尝试 newFunc
	if p.newFunc != nil {
		return p.newFunc()
	}
	return nil, ErrPoolEmpty
}

// Size 返回当前空闲对象数。
func (p *Pool) Size() int {
	return p.idle.Len()
}

// Clear 立即清空所有对象，并调用 ExpireFunc。
func (p *Pool) Clear() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.purgeAll()
}

// Close 关闭池子：
//   - 标记关闭，
//   - 停止后台清理，
//   - 清空并过期所有对象。
func (p *Pool) Close() {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return
	}
	p.closed = true
	p.cancelReaper()
	p.purgeAll()
	p.mu.Unlock()
}

// runReaper 每秒扫描一次，清理过期对象。
func (p *Pool) runReaper(ctx context.Context) {
	if p.ttl <= 0 {
		// 永不过期或用即过期，无需后台
		return
	}
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.purgeExpired()
		}
	}
}

// purgeExpired 移除所有已过期的对象（expireAt <= now），并回调 ExpireFunc。
func (p *Pool) purgeExpired() {
	now := time.Now().UnixMilli()
	for {
		r, _ := p.idle.PopFront()
		if r == nil {
			return
		}
		item := r.(*poolItem)
		if item.expireAt > now && p.ttl > 0 {
			// 还没过期，推回头部
			p.idle.PushFront(item)
			return
		}
		if p.expireFunc != nil {
			p.expireFunc(item.obj)
		}
	}
}

// purgeAll 清空所有对象，并调用 ExpireFunc。
func (p *Pool) purgeAll() {
	for {
		r, _ := p.idle.PopFront()
		if r == nil {
			return
		}
		if p.expireFunc != nil {
			p.expireFunc(r.(*poolItem).obj)
		}
	}
}
