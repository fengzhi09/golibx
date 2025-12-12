package gox

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

type CallerInfo struct {
	Dep  int
	Pkg  string
	File string
	Line int
	Func string
}

func (c *CallerInfo) String() string {
	return fmt.Sprintf("%s:%d %s()", c.File, c.Line, c.Func)
}

func CallerDep(dep int) *CallerInfo {
	return findCaller(func(caller *CallerInfo) bool {
		if dep == 0 {
			return true
		}
		if caller != nil {
			return strings.Contains(caller.Func, "CallerDep")
		}
		return false
	}, dep)
}

func CallerLike(name string) *CallerInfo {
	return findCaller(func(caller *CallerInfo) bool {
		if caller != nil {
			return strings.Contains(caller.Func, name) || strings.Contains(caller.File, name)
		}
		return false
	}, 0)
}
func CallerSkip(name string) *CallerInfo {
	return findCaller(func(caller *CallerInfo) bool {
		if caller != nil {
			return strings.Contains(caller.Func, name) || strings.Contains(caller.File, name)
		}
		return false
	}, 1)
}

var NotFound = &CallerInfo{-1, "", "", -1, ""}

func findCaller(starter func(caller *CallerInfo) bool, dep int) *CallerInfo {
	if starter == nil {
		starter = func(caller *CallerInfo) bool {
			return true
		}
	}
	if dep < 0 {
		return NotFound
	}
	start := -1
	for i := 0; ; i++ { // Skip the expected number of frames
		pc, path, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		pkg, name := location(path)
		c := &CallerInfo{i - start, pkg, name, line, string(function(pc))}
		if start >= 0 && i-start == dep {
			return c
		}
		if start < 0 && starter(c) {
			start = i
		}
		if start >= 0 && i-start == dep {
			return c
		}
	}
	return NotFound
}

func Stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := os.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

func location(path string) (string, string) {
	path = strings.ReplaceAll(path, "\\", "/")
	args := strings.Split(path, "/")
	size := len(args)
	fileIdx := MaxN(0, size-1)
	pkgIdx := MaxN(0, fileIdx-2)
	return strings.Join(args[pkgIdx:fileIdx], "."), args[fileIdx]
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}
