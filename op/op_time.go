package op

import (
	"github.com/fengzhi09/golibx/gox"
	"strings"
	"time"
)

type TimeOp string

//goland:noinspection GoNameStartsWithPackageName
const (
	TimeLt        TimeOp = "<"     // 早于(不含边界)
	TimeLte       TimeOp = "<="    // 不晚于(含边界)
	TimeGte       TimeOp = ">="    // 不早于(含边界)
	TimeGt        TimeOp = ">"     // 晚于(不含边界)
	TimeIn        TimeOp = "[a,b]" // 范围内(含双边界)
	TimeInBorderL TimeOp = "[a,b)" // 范围内(含左边界)
	TimeInBorderR TimeOp = "(a,b]" // 范围内(含右边界)
	TimeInBorderN TimeOp = "(a,b)" // 范围内(不含边界)
	UnknownT      TimeOp = ""      // 未知
)

var TimeOps = []TimeOp{TimeLt, TimeLte, TimeGte, TimeGt, TimeIn, TimeInBorderL, TimeInBorderR, TimeInBorderN}

func AsTimeOp(opStr string) TimeOp {
	opStr = "time/" + strings.ToLower(opStr)
	for _, timeOp := range TimeOps {
		if timeOp.String() == opStr {
			return timeOp
		}
	}
	return UnknownT
}

func (re TimeOp) String() string {
	if re == "" {
		return ""
	}
	return "time/" + string(re)
}

func (re TimeOp) Accept(src time.Time, args []time.Time) bool {
	switch re {
	case TimeGte, TimeGt, TimeLte, TimeLt:
		return re.accept2(gox.UnixMilli(src), gox.UnixMilli(args[0]))
	case TimeIn, TimeInBorderL, TimeInBorderR, TimeInBorderN:
		return re.accept3(gox.UnixMilli(src), gox.UnixMilli(args[0]), gox.UnixMilli(args[1]))
	}
	panic("op not support: " + re.String())
}

func (re TimeOp) accept2(a, b int64) bool {
	switch re {
	case TimeGte:
		return a >= b
	case TimeGt:
		return a > b
	case TimeLte:
		return a <= b
	case TimeLt:
		return a < b
	default:
		panic("op not calcScore: " + re.String())
	}
}

func (re TimeOp) accept3(src, min, max int64) bool {
	switch re {
	case TimeIn:
		return src >= min && src <= max
	case TimeInBorderL:
		return src >= min && src < max
	case TimeInBorderR:
		return src > min && src <= max
	case TimeInBorderN:
		return src > min && src < max
	default:
		panic("op not accept3: " + re.String())
	}
}
