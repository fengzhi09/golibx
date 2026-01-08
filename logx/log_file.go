package logx

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fengzhi09/golibx/gox"
)

type fileWriter struct {
	sync.Mutex
	dir        string
	interval   int
	maxMb      int
	filePrefix string
	suffixLen  int
	rotateAt   string
	level      LogLevel
	cur        *os.File
}

type RotateConf struct {
	Dir    string
	Enable bool
	Level  LogLevel
	Rotate string
	MaxMb  int
}

func NewRotate() *RotateConf {
	return &RotateConf{
		Dir:    "log/",
		Enable: true,
		Level:  INFO,
		Rotate: fmt.Sprintf("%dd", 1),
		MaxMb:  -1,
	}
}

func (r *RotateConf) SetDir(dir string) *RotateConf {
	r.Dir = dir
	return r
}

func (r *RotateConf) SetLevel(level LogLevel) *RotateConf {
	r.Level = level
	return r
}

func (r *RotateConf) ByDays(days int) *RotateConf {
	r.Rotate = fmt.Sprintf("%dd", days)
	return r
}

func (r *RotateConf) ByHours(hours int) *RotateConf {
	r.Rotate = fmt.Sprintf("%dh", hours)
	return r
}

func (r *RotateConf) ByMb(mb int) *RotateConf {
	r.MaxMb = mb
	return r
}

func RotateWriter(name string, conf *RotateConf) LogWriter {
	interval, unit := gox.ParseInterval(conf.Rotate)
	suffixMap := map[string]int{"d": 6, "h": 8, "m": 10, "s": 11}
	w := &fileWriter{
		dir:        conf.Dir,
		interval:   gox.LimitIn(interval, int(gox.UnitMinute), int(gox.UnitMonth)),
		filePrefix: name,
		suffixLen:  gox.GetOrDefault(suffixMap, unit, 6),
		level:      conf.Level,
		maxMb:      conf.MaxMb,
	}
	gox.MkDir(w.dir)
	err := w.tryRotate()
	if err != nil {
		panic(err)
	}
	return w
}

func (c *fileWriter) SetLevelLogx(level LogLevel) {
	c.Lock()
	defer c.Unlock()
	c.level = level
}

func suffixSize(t time.Time, size int) string {
	return gox.SubStr(t.Format(gox.YYYYMMDDHHMMSS), 2, 2+size)
}

func (c *fileWriter) tryRotate() error {
	var err error
	now := gox.Now()
	if c.rotateAt != "" {
		last := gox.AsTime(c.rotateAt)
		if int(gox.Now().Sub(last).Seconds()) > c.interval {
			c.rotateAt = suffixSize(now, c.suffixLen)
			_ = c.cur.Close()
		}
	} else {
		c.rotateAt = suffixSize(now, c.suffixLen)
	}

	// 检查文件大小
	info, _ := c.cur.Stat()
	if c.maxMb > 0 && info != nil && int(info.Size()/gox.MB) > c.maxMb {
		c.rotateAt = suffixSize(now, c.suffixLen)
	}
	nxt, cur := c.nxtPath(), c.curPath()
	if nxt != cur {
		c.cur, err = os.OpenFile(nxt, os.O_APPEND|os.O_CREATE, 0o644)
	}
	return err
}

func (c *fileWriter) nxtPath() string {
	return strings.ReplaceAll(fmt.Sprintf("%s/%s-%s.log", c.dir, c.filePrefix, c.rotateAt), "\\", "/")
}

func (c *fileWriter) curPath() string {
	if c.cur != nil {
		return strings.ReplaceAll(c.cur.Name(), "\\", "/")
	}
	return ""
}

func (c *fileWriter) SetInterval(interval int) {
	c.Lock()
	defer c.Unlock()
	c.interval = interval
}

func (c *fileWriter) SetRootDir(dir string) {
	c.Lock()
	defer c.Unlock()
	c.dir = gox.MkDir(dir)
}

func (c *fileWriter) SetName(name string) {
	c.Lock()
	defer c.Unlock()
	c.filePrefix = name
}

func (c *fileWriter) Write(event *LogEvent) error {
	c.Lock()
	defer c.Unlock()
	err := c.tryRotate()
	if err != nil {
		return err
	}

	if event.Level >= c.level {
		timeStr := event.Time.Format(gox.ISODateTimeMs)
		for len(timeStr) < 23 {
			timeStr += "0"
		}
		msg := fmt.Sprintf("%v [%v] %v\n", timeStr, event.Level.String(), event.Msg)
		_, err = c.cur.Write([]byte(msg))
	}
	return err
}

func (c *fileWriter) Close(ctx context.Context) error {
	if c.cur != nil {
		return c.cur.Close()
	}
	return nil
}
