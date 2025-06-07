package jqueue

import (
	"math"

	"github.com/e7coding/coding-common/container/jatomic"
	"github.com/e7coding/coding-common/container/jlist"
)

// Queue 是一个支持可选容量限制的并发安全队列。
// 当指定了 limit（>0）时，内部直接使用带缓冲的 channel；
// 否则使用无界链表 + 后台 goroutine 推送到 channel。
type Queue struct {
	limit  int              // 队列容量限制，<=0 表示不限
	list   *jlist.SafeList  // 无界模式下的链表存储
	closed *jatomic.Bool    // 是否已关闭
	events chan struct{}    // 无界模式下的 push 触发信号
	C      chan interface{} // 对外可读的 channel
}

const (
	defaultChanSize = 10000 // 默认的 channel 缓冲区大小
)

// New 创建并返回一个新队列。
// 可选参数 limit 指定容量限制（>0），否则为不限制。
func New(limit ...int) *Queue {
	q := &Queue{closed: jatomic.NewBool()}
	if len(limit) > 0 && limit[0] > 0 {
		q.limit = limit[0]
		q.C = make(chan interface{}, q.limit)
	} else {
		q.list = jlist.NewSafeList()
		q.events = make(chan struct{}, math.MaxInt32)
		q.C = make(chan interface{}, defaultChanSize)
		go q.loopPump()
	}
	return q
}

// Push 将 v 放入队列；在关闭后调用会 panic。
func (q *Queue) Push(v interface{}) {
	if q.limit > 0 {
		q.C <- v
	} else {
		q.list.PushBack(v)
		select {
		case q.events <- struct{}{}:
		default:
		}
	}
}

// Pop 从队列中取出一个元素；关闭后会立即返回 nil。
func (q *Queue) Pop() interface{} {
	return <-q.C
}

// Close 关闭队列，所有阻塞的 Pop 将被唤醒并返回 nil。
func (q *Queue) Close() {
	if !q.closed.CAS(false, true) {
		return
	}
	if q.events != nil {
		close(q.events)
		// 等待后台完成所有数据推送并关闭 C
	} else {
		close(q.C)
	}
}

// Len 返回当前队列长度；无界模式下可能不完全精确。
func (q *Queue) Len() int64 {
	if q.limit > 0 {
		return int64(len(q.C))
	}
	return int64(q.list.Len()) + int64(len(q.C))
}

// loopPump 是无界模式下的后台协程：
// 从链表拉数据，推入 C，直到队列关闭。
func (q *Queue) loopPump() {
	defer close(q.C)
	for {
		// 等待 Push 触发或者队列被关闭
		_, ok := <-q.events
		if !ok || q.closed.Val() {
			return
		}
		// 将链表中的所有待处理项推到 C
		for {
			v, _ := q.list.PopFront()
			if v == nil {
				break
			}
			q.C <- v
		}
		// 如果还有残余事件，丢弃多余的触发信号
		select {
		case <-q.events:
		default:
		}
	}
}
