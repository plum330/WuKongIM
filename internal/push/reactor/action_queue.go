package reactor

import (
	"sync"

	"github.com/WuKongIM/WuKongIM/internal/reactor"
	clusterReactor "github.com/WuKongIM/WuKongIM/pkg/cluster/reactor"
	"github.com/WuKongIM/WuKongIM/pkg/wklog"
	"go.uber.org/zap"
)

type actionQueue struct {
	ch            chan struct{}
	rl            *clusterReactor.RateLimiter // 限制字节流量速度
	lazyFreeCycle uint64                      // 懒惰释放周期，n表示n次释放一次
	size          uint64
	left          []reactor.PushAction // 左边队列
	right         []reactor.PushAction // 右边队列, 左右的目的是为了重复利用内存
	nodrop        []reactor.PushAction // 不能drop的消息
	mu            sync.Mutex
	leftInWrite   bool   // 写入时是否使用左边队列
	idx           uint64 // 当前写入的位置下标
	oldIdx        uint64
	cycle         uint64
	wklog.Log
}

func newActionQueue(size uint64, ch bool,
	lazyFreeCycle uint64, maxMemorySize uint64) *actionQueue {

	q := &actionQueue{
		rl:            clusterReactor.NewRateLimiter(maxMemorySize),
		size:          size,
		lazyFreeCycle: lazyFreeCycle,
		left:          make([]reactor.PushAction, size),
		right:         make([]reactor.PushAction, size),
		nodrop:        make([]reactor.PushAction, 0),
		Log:           wklog.NewWKLog("push.actionQueue"),
	}
	if ch {
		q.ch = make(chan struct{}, 1)
	}
	return q
}

func (q *actionQueue) add(msg reactor.PushAction) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.idx >= q.size {
		return false
	}
	if !q.tryAdd(msg) {
		return false
	}

	w := q.targetQueue()
	w[q.idx] = msg
	q.idx++
	return true

}

// 必须要添加的消息不接受drop
func (q *actionQueue) mustAdd(msg reactor.PushAction) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.nodrop = append(q.nodrop, msg)
}

func (q *actionQueue) get() []reactor.PushAction {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.cycle++
	sz := q.idx
	q.idx = 0
	t := q.targetQueue()
	q.leftInWrite = !q.leftInWrite
	q.gc()
	q.oldIdx = sz
	if q.rl.Enabled() {
		q.rl.Set(0)
	}
	if len(q.nodrop) == 0 {
		return t[:sz]
	}

	var result []reactor.PushAction
	if len(q.nodrop) > 0 {
		ssm := q.nodrop
		q.nodrop = make([]reactor.PushAction, 0)
		result = append(result, ssm...)
	}
	return append(result, t[:sz]...)
}

func (q *actionQueue) targetQueue() []reactor.PushAction {
	var t []reactor.PushAction
	if q.leftInWrite {
		t = q.left
	} else {
		t = q.right
	}
	return t
}

func (q *actionQueue) tryAdd(msg reactor.PushAction) bool {
	if !q.rl.Enabled() {
		return true
	}
	if q.rl.RateLimited() {
		q.Warn("rate limited dropped", zap.String("actionType", msg.Type.String()))
		return false
	}
	q.rl.Increase(msg.Size())
	return true
}

func (q *actionQueue) gc() {
	if q.lazyFreeCycle > 0 {
		oldq := q.targetQueue()
		if q.lazyFreeCycle == 1 {
			for i := uint64(0); i < q.oldIdx; i++ {
				oldq[i].Messages = nil
			}
		} else if q.cycle%q.lazyFreeCycle == 0 {
			for i := uint64(0); i < q.size; i++ {
				oldq[i].Messages = nil
			}
		}
	}
}
