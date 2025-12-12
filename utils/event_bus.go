package utils

import (
	"context"
	"sync"
	"time"

	"github.com/fengzhi09/golibx/dbx"
	"github.com/fengzhi09/golibx/gox"
	"github.com/fengzhi09/golibx/jsonx"
	"github.com/fengzhi09/golibx/logx"
	"github.com/fengzhi09/golibx/op"

	"go.uber.org/atomic"
)

type (
	EventType  = string
	StepHandle = func(ctx context.Context, step EventStep, event *MemEvent) error
	EventStep  = string
	MemEvent   struct {
		Id   dbx.OID
		Type EventType
		Data any
	}
)

type handler_ struct {
	name       string
	accept     func(EventType) bool
	stepHandle StepHandle
}

type MemEventBus struct {
	name      string
	lifetime  []EventStep
	queue     []*MemEvent
	handlers  []*handler_
	queueSize int64
	running   *atomic.Bool
	debug     Debugger
	*sync.Mutex
}

//goland:noinspection GoUnusedConst
const defaultEventStep = "" //nolint
type Debugger func(ctx context.Context, format string, args ...any)

func NewEventBus(name string, queueSize int64, debug Debugger, steps ...EventStep) *MemEventBus {
	if len(steps) == 0 {
		steps = []EventStep{defaultEventStep}
	}
	bus := &MemEventBus{
		name:      name,
		debug:     debug,
		lifetime:  steps,
		queue:     []*MemEvent{},
		handlers:  []*handler_{},
		queueSize: queueSize,
		running:   atomic.NewBool(true),
		Mutex:     &sync.Mutex{},
	}
	return bus
}

func (re *MemEventBus) AsyncEnable() bool {
	re.Lock()
	defer re.Unlock()
	return re.queueSize > 0
}

func AcceptTypes(types ...EventType) Acceptor {
	return func(got EventType) bool {
		return gox.IndexOf(types, got) >= 0
	}
}

func AcceptAny() Acceptor {
	return func(got EventType) bool {
		return true
	}
}

func AcceptPattern(patterns ...EventType) Acceptor {
	return func(got EventType) bool {
		pass, _, _ := op.WordIncld.Accept(got, patterns, jsonx.JObj{})
		return pass
	}
}

type Acceptor = func(EventType) bool

func (re *MemEventBus) Sub(subName string, handle StepHandle, accept Acceptor) {
	re.Lock()
	defer re.Unlock()
	re.doDel(subName)
	handler := &handler_{name: subName, accept: accept, stepHandle: handle}
	re.handlers = append(re.handlers, handler)
	re.Debugf(context.Background(), "sub bus[%v] handler[%v]", re.name, subName)
}

func (re *MemEventBus) Debugf(ctx context.Context, format string, args ...any) {
	if re.debug != nil {
		re.debug(ctx, format, re.name, args)
	}
}

func (re *MemEventBus) Del(subName string) {
	re.Lock()
	defer re.Unlock()
	re.doDel(subName)
	re.Debugf(context.Background(), "sub bus[%v] handler[%v] released", re.name, subName)
}

func (re *MemEventBus) doDel(subName string) {
	handlers := make([]*handler_, 0)
	for _, handler := range re.handlers {
		if handler.name != subName {
			handlers = append(handlers, handler)
		}
	}
	re.handlers = handlers
}

func (re *MemEventBus) findHandlers(event *MemEvent) []*handler_ {
	re.Lock()
	defer re.Unlock()

	handlers := make([]*handler_, 0)
	for _, handler := range re.handlers {
		if handler.accept != nil && handler.accept(event.Type) {
			handlers = append(handlers, handler)
		}
	}

	return handlers
}

func (re *MemEventBus) PubSync(ctx context.Context, event *MemEvent) {
	if !re.running.Load() {
		logx.Errorf(ctx, "event-bus is stopped,skip event:=%v", event.Id.Hex())
		return
	}

	re.doPub(ctx, event)
}

func (re *MemEventBus) PubAsync(ctx context.Context, event *MemEvent) {
	if !re.running.Load() {
		logx.Errorf(ctx, "event-bus is stopped,skip event:=%v", event.Id.Hex())
		return
	}

	re.Lock()
	defer re.Unlock()
	re.queue = append(re.queue, event)
}

func (re *MemEventBus) doPub(ctx context.Context, event *MemEvent) {
	handlers := re.findHandlers(event)
	if len(handlers) == 0 {
		logx.Warnf(ctx, "event=%v skipped on no handler", event.Id.Hex())
	}
	for _, handler := range handlers {
		for _, step := range re.lifetime {
			err := handler.stepHandle(ctx, step, event)
			if err != nil && err.Error() != "filtered" {
				logx.Debugf(ctx, "event=%v failed on handler %v step %v,err:%v", event.Id.Hex(), handler.name, step, err)
				break
			} else if err != nil {
				logx.Warnf(ctx, "event=%v failed on handler %v step %v,err:%v", event.Id.Hex(), handler.name, step, err)
				break
			}
		}
	}
}

func (re *MemEventBus) Close() {
	re.running.Store(false)
	ctx := context.Background()
	time.Sleep(3 * time.Second)
	for _, event := range re.queue {
		re.doPub(ctx, event)
		time.Sleep(100 * time.Millisecond)
	}
}

func (re *MemEventBus) Start() {
	if re.AsyncEnable() {
		popHead := func(max int) []*MemEvent {
			re.Lock()
			defer re.Unlock()
			head := make([]*MemEvent, 0)
			i := 0
			for i < max && i < len(re.queue) {
				head = append(head, re.queue[i])
				i++
			}
			re.queue = re.queue[i:]
			return head
		}
		runOnce := func() {
			start := time.Now()
			ctx, head := context.Background(), popHead(10)
			for _, event := range head {
				re.doPub(ctx, event)
				if time.Since(start).Milliseconds() >= 100 {
					break
				}
			}
		}
		go func() {
			time.Sleep(3 * time.Second)
			for re.running.Load() {
				runOnce()
				time.Sleep(100 * time.Millisecond)
			}
		}()
	}
}
