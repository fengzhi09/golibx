package utils

import (
	"context"
	"fmt"
	"github.com/fengzhi09/golibx/dbx"
	"testing"
	"time"
)

func TestEventBus(t *testing.T) {
	const (
		step0 EventStep = "step0"
		step1 EventStep = "step1"
		step2 EventStep = "step2"
		step3 EventStep = "step3"
	)
	const (
		type1 EventType = "type1"
		type2 EventType = "type2"
	)
	bus1 := NewEventBus("test", 1000, nil, step1, step2, step3)
	bus1.Sub("handler[step1]", func(ctx context.Context, step EventStep, event *MemEvent) error {
		if step == step1 {
			if event.Type != type1 {
				t.Errorf("got bad type:%v vs %v", event.Type, type1)
				return fmt.Errorf("got bad type")
			}
			t.Logf("handler[step1] got %v-%v-%v", event.Id.Hex(), event.Type, event.Data)
			event.Data = step
		}
		if step == step2 {
			t.Logf("handler[step2] got %v-%v-%v", event.Id.Hex(), event.Type, event.Data)
			event.Data = step
		}
		if step == step3 {
			t.Logf("handler[step3] got %v-%v-%v", event.Id.Hex(), event.Type, event.Data)
			return fmt.Errorf("some err")
		}
		return nil
	}, AcceptTypes(type1))

	events := []*MemEvent{
		{dbx.NewOID(), type1, step0},
		{dbx.NewOID(), type2, step0},
	}
	ctx := context.Background()
	t.Logf("PubSync begin")
	for _, event := range events {
		bus1.PubSync(ctx, event)
	}
	bus1.Start()
	t.Logf("PubAsync begin")
	for _, event := range events {
		bus1.PubAsync(ctx, event)
	}
	time.Sleep(5 * time.Second)
	t.Logf("PubAsync done")
	bus1.Close()
	t.Logf("All done")
}
