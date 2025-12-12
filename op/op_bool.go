package op

import (
	"fmt"
	"strings"
)

type BoolOp string

//goland:noinspection GoNameStartsWithPackageName
const (
	And      BoolOp = "and" // 同时满足多个条件
	Or       BoolOp = "or"  // 满足多个条件中的任意一个
	Not      BoolOp = "not" // 不满足多个条件中的任意一个
	UnknownB BoolOp = ""    // 未知
)

var BoolOps = []BoolOp{And, Or, Not}

func AsBoolOp(opStr string) BoolOp {
	opStr = "bool/" + strings.ToLower(opStr)
	for _, hasOp := range BoolOps {
		if hasOp.String() == opStr {
			return hasOp
		}
	}
	return UnknownB
}

func (re BoolOp) String() string {
	if re == "" {
		return ""
	}
	return "bool/" + string(re)
}

type (
	BoolFuncE func() (bool, error)
	BoolFunc  func() bool
)

func (re BoolOp) Accept(hits []bool) bool {
	switch re {
	case And:
		return andValue(hits)
	case Or:
		return orValue(hits)
	case Not:
		return notValue(hits)
	}
	return false
}

func (re BoolOp) AcceptFunc(hits []BoolFunc) bool {
	switch re {
	case And:
		return andFunc(hits)
	case Or:
		return orFunc(hits)
	case Not:
		return notFunc(hits)
	}
	return false
}

func (re BoolOp) AcceptFuncE(hits []BoolFuncE) (bool, error) {
	switch re {
	case And:
		return andFuncE(hits)
	case Or:
		return orFuncE(hits)
	case Not:
		return notFuncE(hits)
	}
	return false, fmt.Errorf("unknown op: %s", re.String())
}

func orFunc(hits []BoolFunc) bool {
	for _, hit := range hits {
		if hit() {
			return true
		}
	}
	return false
}

func andFunc(hits []BoolFunc) bool {
	for _, hit := range hits {
		if !hit() {
			return false
		}
	}
	return true
}

func notFunc(hits []BoolFunc) bool {
	for _, hit := range hits {
		if hit() {
			return false
		}
	}
	return true
}

func orFuncE(hits []BoolFuncE) (bool, error) {
	for _, hit := range hits {
		if ok, err := hit(); ok && err == nil {
			return true, nil
		} else if err != nil {
			return false, err
		}
	}
	return false, nil
}

func andFuncE(hits []BoolFuncE) (bool, error) {
	for _, hit := range hits {
		if ok, err := hit(); !ok || err != nil {
			return false, err
		}
	}
	return true, nil
}

func notFuncE(hits []BoolFuncE) (bool, error) {
	for _, hit := range hits {
		if ok, err := hit(); ok && err == nil {
			return false, nil
		} else if err != nil {
			return false, err
		}
	}
	return true, nil
}

func orValue(hits []bool) bool {
	for _, hit := range hits {
		if hit {
			return true
		}
	}
	return false
}

func andValue(hits []bool) bool {
	for _, hit := range hits {
		if !hit {
			return false
		}
	}
	return true
}

func notValue(hits []bool) bool {
	for _, hit := range hits {
		if hit {
			return false
		}
	}
	return true
}
