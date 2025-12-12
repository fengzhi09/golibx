package gox

import (
	"math"
	"strings"

	"golang.org/x/exp/constraints"
)

// Runtime-safe integer limits
const (
	MaxInt64  int64   = 1<<63 - 1
	MinInt64  int64   = -1 << 63
	MaxInt            = int(^uint(0) >> 1)
	MinInt            = -MaxInt - 1
	MaxDouble float64 = 1.7976931348623157e+308
	MinDouble float64 = -1.7976931348623157e+308
	EpsilonD  float64 = 4.94065645841247e-324
	MaxFloat  float32 = 3.402823e+38
	MinFloat  float32 = -3.402823e+38
	EpsilonF  float32 = 1.401298e-45
	NaN       float64 = float64(0x7FF8000000000001)
	PosInf    float64 = float64(0x7FF0000000000000)
	NegInf    float64 = float64(0xFFF0000000000000)
)

type Number interface {
	constraints.Integer | constraints.Float
}

func IsNanF(a float32) bool {
	return float64(a) == NaN
}

func IsNan(a float64) bool {
	return a == NaN
}

func IsInf(a float64) bool {
	return IsInfPos(a) || IsInfNeg(a)
}

func IsInfPos(a float64) bool {
	return a == PosInf
}

func IsInfNeg(a float64) bool {
	return a == NegInf
}

func EqF(a, b float32) bool {
	return float32(math.Abs(float64(a-b))) <= EpsilonF
}

func EqD(a, b float64) bool {
	return math.Abs(a-b) <= EpsilonD
}

// MinS returns the minimum value in a slice of numeric values.
func MinS(values ...string) string {
	minmax := MinMaxS(values...)
	return minmax.MinVal
}

// MaxS returns the maximum value in a slice of strings.
func MaxS(values ...string) string {
	minmax := MinMaxS(values...)
	return minmax.MaxVal
}

type MinMax[T any] struct {
	MinIdx int
	MaxIdx int
	MinVal T
	MaxVal T
	Cnt    int
}

// MinMaxS returns the indices of the min and max values by str comparison.
func MinMaxS(strs ...string) *MinMax[string] {
	if len(strs) == 0 {
		return &MinMax[string]{MinIdx: 0, MaxIdx: 0, MinVal: "", MaxVal: "", Cnt: 0}
	}
	minV, maxV := 0, 0
	for idx, v := range strs[1:] {
		actualIdx := idx + 1 // 因为我们从strs[1:]开始遍历，所以实际索引是idx+1
		if v < strs[minV] {
			minV = actualIdx
		}
		if v > strs[maxV] {
			maxV = actualIdx
		}
	}
	return &MinMax[string]{MinIdx: minV, MaxIdx: maxV, MinVal: strs[minV], MaxVal: strs[maxV], Cnt: len(strs)}
}

// MinABC returns the minimum value in a slice of numeric values.
func MinABC[T any](toStr func(T) string, values ...T) (T, int) {
	minmax := MinMaxABC(toStr, values...)
	return minmax.MinVal, minmax.MinIdx
}

// MaxABC returns the maximum value in a slice of strings.
func MaxABC[T any](toStr func(T) string, values ...T) (T, int) {
	minmax := MinMaxABC(toStr, values...)
	return minmax.MaxVal, minmax.MaxIdx
}

// MinMaxABC returns the indices of the min and max values by abc
func MinMaxABC[T any](toStr func(T) string, values ...T) *MinMax[T] {
	return MinMaxBy(func(a, b T) int {
		return strings.Compare(toStr(a), toStr(b))
	}, values...)
}

// MinMaxBy returns the indices of the min and max values by comparison.
func MinMaxBy[T any](cmp func(T, T) int, values ...T) *MinMax[T] {
	if len(values) == 0 {
		return nil
	}
	minV, maxV := 0, 0
	for idx, v := range values[1:] {
		actualIdx := idx + 1 // 因为我们从values[1:]开始遍历，所以实际索引是idx+1
		if cmp(v, values[minV]) < 0 {
			minV = actualIdx
		}
		if cmp(v, values[maxV]) > 0 {
			maxV = actualIdx
		}
	}
	return &MinMax[T]{
		MinIdx: minV,
		MaxIdx: maxV,
		MinVal: values[minV],
		MaxVal: values[maxV],
	}
}

// KeyCnt returns count stats of values by keyGen.
func KeyCnt[T any](keyGen func(T) string, values ...T) map[string]int {
	cnt := make(map[string]int)
	for _, v := range values {
		cnt[keyGen(v)]++
	}
	return cnt
}

// Min returns the minimum value in a slice of ordered values.
func MinN[T Number](values ...T) T {
	stats := SummaryN(values...)
	return stats.Min
}

func LimitIn[T Number](val, min, max T) T {
	val = IfElse(val > max, max, val).(T)
	val = IfElse(val < min, min, val).(T)
	return val
}

// Max returns the maximum value in a slice of ordered values.
// Panics if values is empty.
func MaxN[T Number](values ...T) T {
	stats := SummaryN(values...)
	return stats.Max
}

func DivN[T Number](a, b T) float64 {
	if b == 0 {
		return 0
	}
	return float64(a) / float64(b)
}

// Sum returns the sum of a slice of numeric values.
func SumN[T Number](values ...T) T {
	stats := SummaryN(values...)
	return stats.Sum
}

// Avg returns the average of a slice of numeric values.
func AvgN[T Number](values ...T) float64 {
	stats := SummaryN(values...)
	return stats.Avg
}

// MinMax returns both the minimum and maximum values in a slice of ordered values.
// Panics if values is empty.
// MinMaxN is an alias for MinMax
func MinMaxN[T constraints.Integer | constraints.Float](values ...T) (min, max T) {
	stats := SummaryN(values...)
	return stats.Min, stats.Max
}

type Stats[T Number] struct {
	Min T
	Max T
	Sum T
	Avg float64
	Cnt int
}

func SummaryN[T Number](values ...T) Stats[T] {
	cnt := len(values)
	if cnt == 0 {
		return Stats[T]{T(0), T(0), T(0), 0, cnt}
	}
	minV, maxV, sum := values[0], values[0], values[0]
	for _, v := range values[1:] {
		sum += v

		if v < minV {
			minV = v
		}
		if v > maxV {
			maxV = v
		}
	}
	return Stats[T]{Min: minV, Max: maxV, Sum: sum, Avg: float64(sum) / float64(cnt), Cnt: cnt}
}
