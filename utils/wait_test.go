package utils

import (
	"fmt"
	"testing"
	"time"

	"github.com/fengzhi09/golibx/gox"

	"go.uber.org/atomic"
)

func TestWaitCondOk(t *testing.T) {
	round, start := atomic.NewInt32(0), time.Now()
	cond := func() bool {
		fmt.Printf("round %v run at %v\n", round.Load(), time.Now().Format(gox.ISODateTimeMs))
		round.Inc()
		return round.Load() >= 5
	}
	err := WaitCond(cond, time.Second*2, time.Second*12)
	if err != nil {
		t.Errorf("wart for met err: %v", err)
	}
	if round.Load() != 5 {
		t.Errorf("round is too low: %v", round.Load())
	}
	t.Logf("%v rounds cost %v sec", round.Load(), time.Since(start).Seconds())
}

func TestWaitCondTimeout(t *testing.T) {
	round, start := atomic.NewInt32(0), time.Now()
	cond := func() bool {
		fmt.Printf("round %v run at %v\n", round.Load(), time.Now().Format(gox.ISODateTimeMs))
		round.Inc()
		return round.Load() >= 5
	}
	err := WaitCond(cond, time.Second*2, time.Second*6)
	if err == nil {
		t.Errorf("wart for should timeout")
	}
	if round.Load() >= 5 {
		t.Errorf("round is too high: %v", round)
	}
	t.Logf("%v rounds cost %v sec", round.Load(), time.Since(start).Seconds())
}
