package logx

import (
	"context"
	"fmt"
	"sync"

	"github.com/fengzhi09/golibx/gox"
)

type cmdWriter struct {
	sync.Mutex
	level LogLevel
	timed bool
}

func CmdWriter(level LogLevel, timed bool) LogWriter {
	return &cmdWriter{level: level, timed: timed}
}

func (c *cmdWriter) SetLevelLogx(level LogLevel) {
	c.Lock()
	defer c.Unlock()
	c.level = level
}

func (c *cmdWriter) Write(event *LogEvent) error {
	c.Lock()
	defer c.Unlock()
	if event.Level >= c.level {
		timeStr := event.Time.Format(gox.ISODateTimeMs)
		for len(timeStr) < 23 {
			timeStr += "0"
		}
		fmt.Printf("%v [%v] %v\n", gox.IfElse(c.timed, timeStr, ""), event.Level.String(), event.Msg)

	}
	return nil
}

func (c *cmdWriter) SetTimed(timed bool) {
	c.timed = timed
}

func (c *cmdWriter) Close(ctx context.Context) error {
	return nil
}
