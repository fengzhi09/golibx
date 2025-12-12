package utils

import (
	"fmt"
	"time"
)

func DelayRun(fn func(), delay time.Duration) {
	timer := time.NewTimer(delay)
	select {
	case <-timer.C:
		fn()
	}
}

type CondFunc func() bool

func WaitCond(cond CondFunc, interval, timeout time.Duration) error {
	return wartFor(NewWatch(interval), cond, timeout)
}

func Wait(period time.Duration) {
	start := time.Now()
	_ = wartFor(NewWatch(time.Millisecond), func() bool {
		return time.Since(start) >= period
	}, period)
}

func NewWatch(interval time.Duration) watch {
	return func(done <-chan struct{}) <-chan struct{} {
		ch := make(chan struct{})
		go func() {
			defer close(ch)
			tick := time.NewTicker(interval)
			defer tick.Stop()
			for {
				select {
				case <-tick.C: //间隔到点
					select {
					case ch <- struct{}{}:
					default:
					}
				case <-done:
					return
				}
			}
		}()
		return ch
	}
}

type watch func(done <-chan struct{}) <-chan struct{}

func wartFor(watcher watch, cond CondFunc, timeout time.Duration) error {
	stopCH := make(chan struct{})
	defer close(stopCH)
	listener := watcher(stopCH)

	var after <-chan time.Time
	if timeout != 0 {
		timer := time.NewTimer(timeout)
		after = timer.C
		defer timer.Stop()
	}
	for {
		select {
		case _, open := <-listener:
			if cond() {
				return nil
			}
			if !open {
				return fmt.Errorf("closed")
			}
		case <-after: // 超时到点
			return fmt.Errorf("timeout")
		}
	}
}
