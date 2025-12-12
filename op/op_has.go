package op

import (
	"strings"

	"github.com/fengzhi09/golibx/gox"
	"github.com/fengzhi09/golibx/jsonx"
)

type HasOp string

//goland:noinspection GoNameStartsWithPackageName
const (
	HasNone  HasOp = "none"
	HasAny   HasOp = "any"
	HasOne   HasOp = "one"
	HasAll   HasOp = "all"
	UnknownH HasOp = "" // 未知
)

var HasOps = []HasOp{HasNone, HasAny, HasOne, HasAll}

func AsHasOp(opStr string) HasOp {
	opStr = "has/" + strings.ToLower(opStr)
	for _, hasOp := range HasOps {
		if hasOp.String() == opStr {
			return hasOp
		}
	}
	return UnknownH
}

func (re HasOp) String() string {
	if re == "" {
		return ""
	}
	return "has/" + string(re)
}

func (re HasOp) Accept(val string, opts []string) (bool, string) {
	switch re {
	default:
		panic("unknown has hasOp:" + gox.AsStr(re))
	case HasAny:
		return val != "", val
	case HasOne:
		idx := gox.IndexOf(opts, val)
		return idx >= 0, jsonx.GoV2JV(opts).ToArr().GetStr(idx)
	case HasNone:
		idx := gox.IndexOf(opts, val)
		return idx < 0, jsonx.GoV2JV(opts).ToArr().GetStr(idx)
	case HasAll:
		for _, opt := range opts {
			if !strings.Contains(val, opt) {
				return false, opt
			}
		}
		return true, val
	}
}

func (re HasOp) AcceptStrict(val string, opts []string) (bool, string) {
	switch re {
	default:
		panic("unknown has hasOp:" + gox.AsStr(re))
	case HasAny:
		return val != "", val
	case HasOne:
		idx := gox.IndexOf(opts, val)
		return idx >= 0, gox.At(idx, opts...)
	case HasNone:
		idx := gox.IndexOf(opts, val)
		return idx < 0, gox.At(idx, opts...)
	case HasAll:
		for _, opt := range opts {
			if val != opt {
				return false, opt
			}
		}
		return true, val
	}
}

func (re HasOp) AcceptArr(values []string, opts []string) (bool, string) {
	opts = gox.ArrUniq(opts)
	switch re {
	default:
		panic("unknown has hasOp:" + gox.AsStr(re))
	case HasAny:
		return len(values) > 0, jsonx.GoV2JV(values).ToArr().First().String()
	case HasOne:
		for _, val := range values {
			idx := gox.IndexOf(opts, val)
			if idx >= 0 {
				return true, gox.At(idx, opts...)
			}
		}
		return false, "none match"
	case HasNone:
		for _, val := range values {
			idx := gox.IndexOf(opts, val)
			if idx >= 0 {
				return false, gox.At(idx, opts...)
			}
		}
		return true, ""
	case HasAll:
		for _, opt := range opts {
			hit := false
			for _, val := range values {
				idx := strings.Index(val, opt)
				if idx >= 0 {
					hit = true
					break
				}
			}
			if !hit {
				return false, "not match " + opt
			}
		}
		return true, ""
	}
}
