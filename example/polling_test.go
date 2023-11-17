package example

import (
	"context"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const BUFFER_SIZE = 256

type Request string

type HookName string

type Poller struct {
	notNull    [BUFFER_SIZE]bool
	data       [BUFFER_SIZE]*Request
	polling    [BUFFER_SIZE]bool
	createTime [BUFFER_SIZE]int64
	lock       sync.Mutex
	initLock   [BUFFER_SIZE]sync.Mutex
}

func (plr *Poller) TryAddRequest(req *Request) (success bool) {
	for i := range BUFFER_SIZE {
		if plr.initLock[i].TryLock() {
			if plr.notNull[i] {
				plr.initLock[i].Unlock()
				continue
			}
			plr.data[i] = req
			plr.polling[i] = false
			plr.createTime[i] = time.Now().UnixNano()
			plr.notNull[i] = true
			plr.initLock[i].Unlock()
			return true
		}
	}
	return false
}

func (plr *Poller) Poll(ctx context.Context, action func(*Request)) {
	var beforeVisit func()
	hook := ctx.Value(HookName("beforeVisit"))
	if hook != nil {
		if f, isFunc := hook.(func()); isFunc {
			beforeVisit = f
		}
	}
	for {
		// get the least recently-polled Resource
		// and mark it as being polled
		plr.lock.Lock()
		select {
		case <-ctx.Done():
			plr.lock.Unlock()
			return
		default:
		}
		handleIdx := -1
		for i := range BUFFER_SIZE {
			if !plr.notNull[i] || plr.polling[i] {
				continue
			}
			if handleIdx == -1 || plr.createTime[i] < plr.createTime[handleIdx] {
				handleIdx = i
			}
		}
		if handleIdx != -1 {
			plr.polling[handleIdx] = true
			beforeVisit()
		}
		plr.lock.Unlock()
		if handleIdx == -1 {
			continue
		}
		action(plr.data[handleIdx])
		plr.initLock[handleIdx].Lock()
		plr.notNull[handleIdx] = false
		plr.initLock[handleIdx].Unlock()
	}
}

func TestPoller1(t *testing.T) {
	const GO_ROUTINE_NUM = 64
	const REQUEST_NUM = 10000000
	plr := new(Poller)
	req := Request("1")
	var sum atomic.Int64
	action := func(r *Request) {
		i, err := strconv.Atoi(string(*r))
		if err != nil {
			return
		}
		sum.Add(int64(i))
	}
	addDone := make(chan struct{})
	go func() {
		for range REQUEST_NUM {
			for !plr.TryAddRequest(&req) {
			}
		}
		close(addDone)
	}()
	handleCtr := 0
	ctx, cancelFunc := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, HookName("beforeVisit"), func() {
		handleCtr++
	})
	for range GO_ROUTINE_NUM {
		go plr.Poll(ctx, action)
	}
	<-addDone
	for {
		if handleCtr == REQUEST_NUM {
			break
		}
	}
	cancelFunc()
	if sum.Load() != REQUEST_NUM {
		t.Errorf("want: %d, but: %d", REQUEST_NUM, sum.Load())
	}
}

func ChanPoll(ctx context.Context, in chan *Request, action func(*Request)) {
	wgDone := func() {}
	wgVal := ctx.Value(HookName("wg"))
	if wgVal != nil {
		if f, ok := wgVal.(func()); ok {
			wgDone = f
		}
	}
	for req := range in {
		action(req)
	}
	wgDone()
}

func TestPoller2(t *testing.T) {
	const GO_ROUTINE_NUM = 64
	const REQUEST_NUM = 10000000
	req := Request("1")
	var sum atomic.Int64
	action := func(r *Request) {
		i, err := strconv.Atoi(string(*r))
		if err != nil {
			return
		}
		sum.Add(int64(i))
	}
	in := make(chan *Request, BUFFER_SIZE)
	go func() {
		for range REQUEST_NUM {
			in <- &req
		}
		close(in)
	}()
	var wg sync.WaitGroup
	wg.Add(GO_ROUTINE_NUM)
	for range GO_ROUTINE_NUM {
		go ChanPoll(
			context.WithValue(context.Background(), HookName("wg"), wg.Done),
			in, action)
	}
	wg.Wait()
	if sum.Load() != REQUEST_NUM {
		t.Errorf("want: %d, but: %d", REQUEST_NUM, sum.Load())
	}
}
