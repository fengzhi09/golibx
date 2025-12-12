package gox

import (
	"fmt"
	"io"
	"regexp"
	"sync"

	"github.com/dlclark/regexp2"
)

func Match(pattern string, b []byte) (matched bool, err error) {
	re, err := Compile(pattern, ReV2)
	if err != nil {
		return false, err
	}
	return re.Match(b), nil
}

func MatchReader(pattern string, r io.RuneReader) (matched bool, err error) {
	re, err := Compile(pattern, ReV2)
	if err != nil {
		return false, err
	}
	return re.MatchReader(r), nil
}

func MatchString(pattern string, s string) (matched bool, err error) {
	re, err := Compile(pattern, ReV2)
	if err != nil {
		return false, err
	}
	return re.MatchString(s), nil
}

func QuoteMeta(s string) string {
	return regexp.QuoteMeta(s)
}

type (
	ReVer string
	ReOpt = regexp2.RegexOptions
)

const (
	ReV1 ReVer = "v1"
	ReV2 ReVer = "v2"
)

type Regexp struct {
	sync.Mutex
	ver  ReVer
	expr string
	api  any
	opt  ReOpt
}

func NewRegexp(expr string, ver ReVer, opt ReOpt) *Regexp {
	return &Regexp{
		ver:  ver,
		expr: expr,
		opt:  opt,
	}
}

func (re *Regexp) Copy() *Regexp {
	return &Regexp{
		ver:  re.ver,
		expr: re.expr,
		api:  re.api,
		opt:  re.opt,
	}
}

func (re *Regexp) init() error {
	return re.Use(re.ver, re.opt, true)
}

func (re *Regexp) Use(ver ReVer, opt ReOpt, force bool) error {
	re.Lock()
	defer re.Unlock()
	var err error
	if re.ver != ver || re.opt != opt || force {
		re.ver = ver
		re.opt = opt
		switch re.ver {
		case ReV1:
			re.api, err = regexp.Compile(re.expr)
		case ReV2:
			re.api, err = regexp2.Compile(re.expr, re.opt)
		default:
			re.api, err = regexp2.Compile(re.expr, re.opt)
		}
	}
	return err
}

func (re *Regexp) POSIX(ver ReVer) error {
	re.Lock()
	defer re.Unlock()
	var err error
	if re.ver != ver {
		re.ver = ver
		switch re.ver {
		case ReV1:
			re.api, err = regexp.CompilePOSIX(re.expr)
		case ReV2:
			re.api, err = regexp2.Compile(re.expr, re.opt)
		default:
			re.api, err = regexp2.Compile(re.expr, re.opt)
		}
	}
	return err
}

func Compile(expr string, ver ReVer) (*Regexp, error) {
	re := NewRegexp(expr, ReVer(ver), regexp2.None)
	return re, re.init()
}

func CompilePOSIX(expr string, ver ReVer) (*Regexp, error) {
	re := NewRegexp(expr, ver, regexp2.None)
	return re, re.init()
}

func MustCompile(str string, ver ReVer) *Regexp {
	re, err := Compile(str, ver)
	if err != nil {
		panic(err)
	}
	return re
}

func MustCompilePOSIX(str string, ver ReVer) *Regexp {
	re, err := CompilePOSIX(str, ver)
	if err != nil {
		panic(err)
	}
	return re
}

func (re *Regexp) AppendText(b []byte) ([]byte, error) {
	return append(b, re.String()...), nil
}

var (
	reVerNotFound = fmt.Errorf("不支持的正则引擎版本")
	reV2NotFound  = fmt.Errorf("regexp2 不支持")
)

func (re *Regexp) Match(b []byte) bool {
	switch api := re.api.(type) {
	case *regexp.Regexp:
		return api.Match(b)
	case *regexp2.Regexp:
		return re.MatchString(string(b))
	default:
		panic(reVerNotFound)
	}
}

func (re *Regexp) MatchReader(r io.RuneReader) bool {
	switch api := re.api.(type) {
	case *regexp.Regexp:
		return api.MatchReader(r)
	case *regexp2.Regexp:
		runes, err := ReadRunes(r)
		if err != nil {
			return false
		}
		matched, err := api.MatchRunes(runes)
		return err == nil && matched
	default:
		panic(reVerNotFound)
	}
}

func (re *Regexp) MatchString(s string) bool {
	switch api := re.api.(type) {
	case *regexp.Regexp:
		return api.MatchString(s)
	case *regexp2.Regexp:
		matched, err := api.MatchString(s)
		return err == nil && matched
	default:
		panic(reVerNotFound)
	}
}

func (re *Regexp) FindString(s string) string {
	switch api := re.api.(type) {
	case *regexp.Regexp:
		return api.FindString(s)
	case *regexp2.Regexp:
		res, err := api.FindStringMatch(s)
		if err != nil {
			return ""
		}
		return res.String()
	default:
		panic(reVerNotFound)
	}
}

// FindAll 从[字符串]中查找所有匹配的[字符串]位置, from(原串下标)=0 表示查找所有匹配项
func (re *Regexp) FindStrings(s string, from int) []string {
	strs := make([]string, 0)
	re.ForEach([]byte(s), from, func(idx int, length int, itm string) {
		strs = append(strs, itm)
	})
	return strs
}

// FindAllBytes 从[字节组]中查找所有匹配的[字节组]位置, from(原字节下标)=0 表示查找所有匹配项
func (re *Regexp) FindBytes(b []byte, from int) [][]byte {
	bytes := make([][]byte, 0)
	re.ForEach(b, from, func(idx int, length int, itm string) {
		bytes = append(bytes, []byte(itm))
	})
	return bytes
}

// FindIndexes 从字符串中查找所有匹配的索引位置, from=0 表示查找所有匹配项
func (re *Regexp) FindIndexes(s string, from int) []int {
	idxs := make([]int, 0)
	re.ForEach([]byte(s), from, func(idx int, length int, itm string) {
		idxs = append(idxs, idx)
	})
	return idxs
}

// FindLocations 从[字符串]中查找所有匹配的[起止点], from=0 表示查找所有匹配项
func (re *Regexp) FindLocations(s string, from int) [][]int {
	idxs := make([][]int, 0)
	re.ForEach([]byte(s), from, func(idx int, length int, itm string) {
		idxs = append(idxs, []int{idx, idx + length})
	})
	return idxs
}

// ForEach 从[字节组]中查找所有匹配项[起始点、长度、字符串]，from=0 表示查找所有匹配项
func (re *Regexp) ForEach(raw []byte, from int, cb func(idx int, length int, itm string)) error {
	switch api := re.api.(type) {
	case *regexp.Regexp:
		idxs := api.FindAllStringIndex(string(raw), from)
		for _, idx := range idxs {
			cb(idx[0], idx[1]-idx[0], string(raw[idx[0]:idx[1]]))
		}
		return nil
	case *regexp2.Regexp:
		return foreachV2(api, string(raw), from, cb)
	default:
		panic(reVerNotFound)
	}
}

func foreachV2(re *regexp2.Regexp, str string, from int, cb func(idx int, length int, itm string)) error {
	m, e := re.FindStringMatchStartingAt(str, from)
	if e != nil {
		return e
	}
	for m != nil {
		cb(m.Index, m.Length, m.String())
		m, e = re.FindNextMatch(m)
		if e != nil {
			return e
		}
		if m == nil {
			break
		}
	}
	return nil
}

func (re *Regexp) MarshalText() ([]byte, error) {
	switch api := re.api.(type) {
	case *regexp.Regexp:
		return api.MarshalText()
	case *regexp2.Regexp:
		return api.MarshalText()
	}
	return nil, reVerNotFound
}

func (re *Regexp) ReplaceAll(src, repl string) string {
	switch api := re.api.(type) {
	case *regexp.Regexp:
		return api.ReplaceAllString(src, repl)
	case *regexp2.Regexp:
		res, err := api.Replace(string(src), string(repl), -1, -1)
		if err != nil {
			return ""
		}
		return res
	default:
		panic(reVerNotFound)
	}
}

func (re *Regexp) ReplaceBytes(src, repl []byte) []byte {
	switch api := re.api.(type) {
	case *regexp.Regexp:
		return api.ReplaceAll(src, repl)
	default:
		return []byte(re.ReplaceAll(string(src), string(repl)))
	}
}

func (re *Regexp) String() string {
	switch api := re.api.(type) {
	case *regexp.Regexp:
		return api.String()
	case *regexp2.Regexp:
		return api.String()
	default:
		panic(reVerNotFound)
	}
}

func (re *Regexp) UnmarshalText(text []byte) error {
	switch api := re.api.(type) {
	case *regexp.Regexp:
		return api.UnmarshalText(text)
	case *regexp2.Regexp:
		return api.UnmarshalText(text)
	default:
		panic(reVerNotFound)
	}
}

func (re *Regexp) GroupNameCnt() int {
	switch api := re.api.(type) {
	case *regexp.Regexp:
		return api.NumSubexp()
	case *regexp2.Regexp:
		return len(api.GetGroupNames())
	default:
		panic(reVerNotFound)
	}
}

func (re *Regexp) GroupNameIndex(name string) int {
	switch api := re.api.(type) {
	case *regexp.Regexp:
		return api.SubexpIndex(name)
	case *regexp2.Regexp:
		return api.GroupNumberFromName(name)
	default:
		panic(reVerNotFound)
	}
}

func (re *Regexp) GroupNames() []string {
	switch api := re.api.(type) {
	case *regexp.Regexp:
		return api.SubexpNames()
	case *regexp2.Regexp:
		return api.GetGroupNames()
	default:
		panic(reVerNotFound)
	}
}

func (re *Regexp) GroupMatch(s string) map[string][]string {
	panic(reVerNotFound)
	// switch api := re.api.(type) {
	// case *regexp.Regexp:
	// 	matches := api.FindAllStringSubmatch(s)
	// 	if matches == nil {
	// 		return nil
	// 	}
	// 	groupMap := make(map[string][]string)
	// 	groupNames := re.GroupNames()
	// 	for i, name := range groupNames {
	// 		groupMap[name] =
	// 	}

	// 	return groupMap
	// case *regexp2.Regexp:
	// 	matches, err := api.FindStringSubmatch(s)
	// 	if err != nil {
	// 		return nil
	// 	}
	// 	return matches
	// default:
	// 	panic(reVerNotFound)
	// }
}
