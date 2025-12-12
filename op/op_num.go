package op

import (
	"strings"

	"github.com/fengzhi09/golibx/gox"
)

type NumOp string

//goland:noinspection GoNameStartsWithPackageName
const (
	NumEq        NumOp = "=="     // 相等
	NumNe        NumOp = "!="     // 不等于
	NumLt        NumOp = "<"      // 小于
	NumIn        NumOp = "[a,b]"  // 范围内(含双边界)
	NumInBorderL NumOp = "[a,b)"  // 范围内(含左边界)
	NumInBorderR NumOp = "(a,b]"  // 范围内(含右边界)
	NumInBorderN NumOp = "(a,b)"  // 范围内(不含边界)
	NumLe        NumOp = "<="     // 不高于
	NumGt        NumOp = ">"      // 大于
	NumGe        NumOp = ">="     // 不低于
	IsNan        NumOp = "isNan"  // 取值是否nan
	IsInf        NumOp = "isInf"  // 取值是否inf
	IsNegInf     NumOp = "negInf" // 取值是否-inf

	NumSet NumOp = "="   // 赋值
	NumSub NumOp = "-"   // 减法
	NumAdd NumOp = "+"   // 加法
	NumDiv NumOp = "/"   // 除法
	NumMul NumOp = "*"   // 乘法
	NumMax NumOp = "max" // 取较大值
	NumMin NumOp = "min" // 取较小值
	// OpIfThen NumOp = "if-than" // 如果……就……
	// OpIfElse NumOp = "if-else" // 如果……就……否则……

	UnknownN NumOp = "" // 未知
)

var (
	CmpOps = []NumOp{NumEq, NumNe, NumLt, NumIn, NumInBorderL, NumInBorderR, NumInBorderN, NumLe, NumGt, NumGe, IsNan, IsInf, IsNegInf}
	CalOps = []NumOp{NumSet, NumSub, NumAdd, NumDiv, NumMul, NumMax, NumMin}
	NumOps = append(CalOps, CalOps...)
)

func AsCmpOp(opStr string) NumOp {
	opStr = "num/cmp/" + strings.ToLower(opStr)
	for _, numOp := range CmpOps {
		if numOp.String() == opStr {
			return numOp
		}
	}
	return UnknownN
}

func AsNumOp(opStr string) NumOp {
	for _, numOp := range NumOps {
		i := strings.LastIndex(numOp.String(), "/")
		if numOp.String()[i+1:] == opStr {
			return numOp
		}
	}
	return UnknownN
}

func (re NumOp) Type() string {
	for i := 0; i < len(CmpOps); i++ {
		if re == CmpOps[i] {
			return "cmp"
		}
	}
	for i := 0; i < len(CalOps); i++ {
		if re == CalOps[i] {
			return "cal"
		}
	}
	return "unknown"
}

func (re NumOp) String() string {
	if re == "" {
		return ""
	}

	return "num/" + re.Type() + "/" + string(re)
}

func (re NumOp) Accept(src float64, args []float64) any {
	switch re {
	case IsInf, IsNan:
		return re.accept1(src) //.(bool)
	case NumGe, NumGt, NumLe, NumLt, NumEq, NumNe:
		return re.accept2(src, args[0]) //.(bool)
	case NumIn, NumInBorderL, NumInBorderR, NumInBorderN:
		if len(args) < 2 {
			return false
		}
		return re.accept3(src, args[0], args[1]) //.(bool)
	case NumAdd, NumSub, NumMul, NumDiv, NumMax, NumMin:
		return re.acceptX(src, args) //.(float64)
	}
	panic("op not support: " + re.String())
}

func (re NumOp) accept1(src float64) bool {
	switch re {
	case IsInf:
		return gox.IsInf(src)
	case IsNan:
		return gox.IsNan(src)
	}
	panic("op not calcScore: " + re.String())
}

func (re NumOp) accept2(a, b float64) bool {
	switch re {
	case NumGe:
		return a >= b
	case NumGt:
		return a > b
	case NumLe:
		return a <= b
	case NumLt:
		return a < b
	case NumEq:
		return gox.EqD(a, b)
	case NumNe:
		return !gox.EqD(a, b)
	default:
		panic("op not calcScore: " + re.String())
	}
}

func (re NumOp) accept3(src, min, max float64) bool {
	switch re {
	case NumIn:
		return src >= min && src <= max
	case NumInBorderL:
		return src >= min && src < max
	case NumInBorderR:
		return src > min && src <= max
	case NumInBorderN:
		return src > min && src < max
	default:
		panic("op not accept3: " + re.String())
	}
}

func (re NumOp) acceptX(src float64, args []float64) float64 {
	res := src
	for _, arg := range args {
		switch re {
		case NumAdd:
			res += arg
		case NumSub:
			res -= arg
		case NumMul:
			res *= arg
		case NumDiv:
			res /= arg
		case NumMax:
			if res < arg {
				res = arg
			}
		case NumMin:
			if res > arg {
				res = arg
			}
		default:
			panic("op not acceptX: " + re.String())
		}
	}
	return res
}
