package gox

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsNan(t *testing.T) {
	assert.True(t, IsNan(NaN))
	assert.False(t, IsNan(0))
	assert.False(t, IsNan(1.0))
	assert.False(t, IsNan(PosInf))
	assert.False(t, IsNan(NegInf))
}

func TestIsNanF(t *testing.T) {
	assert.True(t, IsNanF(float32(NaN)))
	assert.False(t, IsNanF(float32(0)))
	assert.False(t, IsNanF(float32(1.0)))
}

func TestIsInf(t *testing.T) {
	assert.True(t, IsInf(PosInf))
	assert.True(t, IsInf(NegInf))
	assert.False(t, IsInf(0))
	assert.False(t, IsInf(1.0))
	assert.False(t, IsInf(NaN))
}

func TestIsInfPos(t *testing.T) {
	assert.True(t, IsInfPos(PosInf))
	assert.False(t, IsInfPos(NegInf))
	assert.False(t, IsInfPos(0))
}

func TestIsInfNeg(t *testing.T) {
	assert.True(t, IsInfNeg(NegInf))
	assert.False(t, IsInfNeg(PosInf))
	assert.False(t, IsInfNeg(0))
}

func TestEqF(t *testing.T) {
	assert.True(t, EqF(float32(1.0), float32(1.0)))
	assert.True(t, EqF(float32(1.0), float32(1.0)+float32(EpsilonD/2)))
	assert.False(t, EqF(float32(1.0), float32(2.0)))
}

func TestEqD(t *testing.T) {
	assert.True(t, EqD(1.0, 1.0))
	assert.True(t, EqD(1.0, 1.0+EpsilonD/2))
	assert.False(t, EqD(1.0, 2.0))
}

func TestMinS(t *testing.T) {
	// 测试字符串最小值
	assert.Equal(t, "a", MinS("a", "b", "c"))
	assert.Equal(t, "", MinS("", "a", "b"))
	assert.Equal(t, "z", MinS("z"))
}

func TestMaxS(t *testing.T) {
	// 测试字符串最大值
	assert.Equal(t, "c", MaxS("a", "b", "c"))
	assert.Equal(t, "b", MaxS("", "a", "b"))
	assert.Equal(t, "z", MaxS("z"))
}

func TestMinMaxS(t *testing.T) {
	// 测试字符串的MinMaxS
	result := MinMaxS("", "apple", "banana", "cherry")
	assert.NotNil(t, result)
	assert.Equal(t, "", result.MinVal)
	assert.Equal(t, "cherry", result.MaxVal)
	assert.Equal(t, 0, result.MinIdx)
	assert.Equal(t, 3, result.MaxIdx)

	// 测试单个元素
	result = MinMaxS("single")
	assert.NotNil(t, result)
	assert.Equal(t, "single", result.MinVal)
	assert.Equal(t, "single", result.MaxVal)
	assert.Equal(t, 0, result.MinIdx)
	assert.Equal(t, 0, result.MaxIdx)

	// 测试空切片
	result = MinMaxS()
	assert.Equal(t, 0, result.Cnt)
	assert.Equal(t, "", result.MinVal)
	assert.Equal(t, "", result.MaxVal)
}

func TestMinABC(t *testing.T) {
	toStr := func(i int) string { return string(rune('a' + i)) }
	val, idx := MinABC(toStr, 2, 1, 3)
	assert.Equal(t, 1, val)
	assert.Equal(t, 1, idx)
}

func TestMaxABC(t *testing.T) {
	toStr := func(i int) string { return string(rune('a' + i)) }
	val, idx := MaxABC(toStr, 2, 1, 3)
	assert.Equal(t, 3, val)
	assert.Equal(t, 2, idx)
}

func TestMinMaxABC(t *testing.T) {
	toStr := func(i int) string { return string(rune('a' + i)) }
	result := MinMaxABC(toStr, 2, 1, 3)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.MinVal)
	assert.Equal(t, 3, result.MaxVal)
	assert.Equal(t, 1, result.MinIdx)
	assert.Equal(t, 2, result.MaxIdx)

	// 测试空切片
	result = MinMaxABC(toStr)
	assert.Nil(t, result)
}

func TestMinMaxBy(t *testing.T) {
	// 使用自定义比较函数
	cmp := func(a, b int) int {
		return a - b
	}
	result := MinMaxBy(cmp, 3, 1, 4, 1, 5, 9)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.MinVal)
	assert.Equal(t, 9, result.MaxVal)
	assert.Equal(t, 1, result.MinIdx)
	assert.Equal(t, 5, result.MaxIdx)

	// 测试空切片
	result = MinMaxBy(cmp)
	assert.Nil(t, result)
}

func TestKeyCnt(t *testing.T) {
	keyGen := func(s string) string { return string(s[0]) }
	result := KeyCnt(keyGen, "apple", "banana", "cherry", "date")
	assert.Equal(t, 4, len(result))
	assert.Equal(t, 1, result["a"])
	assert.Equal(t, 1, result["b"])
	assert.Equal(t, 1, result["c"])
	assert.Equal(t, 1, result["d"])

	// 测试重复键
	result = KeyCnt(keyGen, "apple", "apricot", "banana")
	assert.Equal(t, 2, len(result))
	assert.Equal(t, 2, result["a"])
	assert.Equal(t, 1, result["b"])

	// 测试空切片
	result = KeyCnt(keyGen)
	assert.Equal(t, 0, len(result))
}

func TestMinN(t *testing.T) {
	// 测试int类型
	assert.Equal(t, 1, MinN(3, 1, 4, 2))
	assert.Equal(t, -10, MinN(5, -10, 3))
	assert.Equal(t, 42, MinN(42))

	// 测试float64类型
	assert.Equal(t, 1.5, MinN(3.0, 1.5, 4.2))

	assert.Equal(t, 0, MinN[int]())
}

func TestMaxN(t *testing.T) {
	// 测试int类型
	assert.Equal(t, 4, MaxN(3, 1, 4, 2))
	assert.Equal(t, 5, MaxN(5, -10, 3))
	assert.Equal(t, 42, MaxN(42))

	// 测试float64类型
	assert.Equal(t, 4.2, MaxN(3.0, 1.5, 4.2))

	// 测试空切片应该返回0
	assert.Equal(t, 0, MaxN[int]())
}

func TestSumN(t *testing.T) {
	// 测试int类型
	assert.Equal(t, 10, SumN(1, 2, 3, 4))
	assert.Equal(t, -2, SumN(-1, -1, 0))
	assert.Equal(t, 42, SumN(42))

	// 测试float64类型
	assert.True(t, EqD(10.5, SumN(1.5, 2.0, 3.0, 4.0)))
}

func TestAvgN(t *testing.T) {
	// 测试int类型
	assert.Equal(t, 2.5, AvgN(1, 2, 3, 4))
	assert.Equal(t, 42.0, AvgN(42))

	// 测试float64类型
	assert.True(t, math.Abs(2.5-AvgN(1.0, 2.0, 3.0, 4.0)) < EpsilonD)
}

func TestMinMaxN(t *testing.T) {
	// 测试int类型
	minV, maxV := MinMaxN(3, 1, 4, 2)
	assert.Equal(t, 1, minV)
	assert.Equal(t, 4, maxV)

	// 测试float64类型
	minFloat, maxFloat := MinMaxN(3.0, 1.5, 4.2)
	assert.Equal(t, 1.5, minFloat)
	assert.Equal(t, 4.2, maxFloat)

	// 测试单个元素
	minSingle, maxSingle := MinMaxN(42)
	assert.Equal(t, 42, minSingle)
	assert.Equal(t, 42, maxSingle)
	minV, maxV = MinMaxN[int]()
	assert.Equal(t, 0, minV)
	assert.Equal(t, 0, maxV)
}

func TestSummaryN(t *testing.T) {
	// 测试int类型
	stats := SummaryN(1, 2, 3, 4)
	assert.Equal(t, 1, stats.Min)
	assert.Equal(t, 4, stats.Max)
	assert.Equal(t, 10, stats.Sum)
	assert.Equal(t, 2.5, stats.Avg)

	// 测试float64类型
	statsFloat := SummaryN(1.5, 2.0, 3.0, 4.0)
	assert.Equal(t, 1.5, statsFloat.Min)
	assert.Equal(t, 4.0, statsFloat.Max)
	assert.True(t, EqD(10.5, statsFloat.Sum))
	assert.True(t, EqD(2.625, statsFloat.Avg))

	// 测试单个元素
	statsSingle := SummaryN(42)
	assert.Equal(t, 42, statsSingle.Min)
	assert.Equal(t, 42, statsSingle.Max)
	assert.Equal(t, 42, statsSingle.Sum)
	assert.Equal(t, 42.0, statsSingle.Avg)

	stats = SummaryN[int]()
	assert.Equal(t, 0, stats.Cnt)
	assert.Equal(t, 0, stats.Min)
	assert.Equal(t, 0, stats.Max)
	assert.Equal(t, 0, stats.Sum)
	assert.Equal(t, float64(0), stats.Avg)
}
