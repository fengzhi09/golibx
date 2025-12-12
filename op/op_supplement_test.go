package op

import (
	"testing"
	"time"

	"github.com/fengzhi09/golibx/jsonx"

	"github.com/stretchr/testify/assert"
)

// 补充测试op_time.go中的accept2和accept3函数
func TestTimeOpAccept2(t *testing.T) {
	// 测试TimeGte
	a := int64(100)
	b := int64(100)
	assert.True(t, TimeGte.accept2(a, b))
	b = int64(99)
	assert.True(t, TimeGte.accept2(a, b))
	b = int64(101)
	assert.False(t, TimeGte.accept2(a, b))

	// 测试TimeGt
	a = int64(100)
	b = int64(100)
	assert.False(t, TimeGt.accept2(a, b))
	b = int64(99)
	assert.True(t, TimeGt.accept2(a, b))
	b = int64(101)
	assert.False(t, TimeGt.accept2(a, b))

	// 测试TimeLte
	a = int64(100)
	b = int64(100)
	assert.True(t, TimeLte.accept2(a, b))
	b = int64(99)
	assert.False(t, TimeLte.accept2(a, b))
	b = int64(101)
	assert.True(t, TimeLte.accept2(a, b))

	// 测试TimeLt
	a = int64(100)
	b = int64(100)
	assert.False(t, TimeLt.accept2(a, b))
	b = int64(99)
	assert.False(t, TimeLt.accept2(a, b))
	b = int64(101)
	assert.True(t, TimeLt.accept2(a, b))
}

func TestTimeOpAccept3(t *testing.T) {
	src := int64(100)
	min := int64(90)
	max := int64(110)

	// 测试TimeIn (含双边界)
	assert.True(t, TimeIn.accept3(src, min, max))
	assert.True(t, TimeIn.accept3(min, min, max))
	assert.True(t, TimeIn.accept3(max, min, max))
	assert.False(t, TimeIn.accept3(89, min, max))
	assert.False(t, TimeIn.accept3(111, min, max))

	// 测试TimeInBorderL (含左边界)
	assert.True(t, TimeInBorderL.accept3(src, min, max))
	assert.True(t, TimeInBorderL.accept3(min, min, max))
	assert.False(t, TimeInBorderL.accept3(max, min, max))
	assert.False(t, TimeInBorderL.accept3(89, min, max))

	// 测试TimeInBorderR (含右边界)
	assert.True(t, TimeInBorderR.accept3(src, min, max))
	assert.False(t, TimeInBorderR.accept3(min, min, max))
	assert.True(t, TimeInBorderR.accept3(max, min, max))
	assert.False(t, TimeInBorderR.accept3(111, min, max))

	// 测试TimeInBorderN (不含边界)
	assert.True(t, TimeInBorderN.accept3(src, min, max))
	assert.False(t, TimeInBorderN.accept3(min, min, max))
	assert.False(t, TimeInBorderN.accept3(max, min, max))
	assert.False(t, TimeInBorderN.accept3(89, min, max))
	assert.False(t, TimeInBorderN.accept3(111, min, max))
}

// 测试TimeOp的完整Accept功能
func TestTimeOpAccept(t *testing.T) {
	// 创建测试时间
	time1 := time.Unix(100, 0)
	time2 := time.Unix(200, 0)
	time3 := time.Unix(300, 0)

	// 测试二元操作符
	assert.True(t, TimeGte.Accept(time2, []time.Time{time1}))
	assert.True(t, TimeGt.Accept(time3, []time.Time{time2}))
	assert.True(t, TimeLte.Accept(time1, []time.Time{time2}))
	assert.True(t, TimeLt.Accept(time1, []time.Time{time2}))

	// 测试三元操作符
	assert.True(t, TimeIn.Accept(time2, []time.Time{time1, time3}))
	assert.True(t, TimeInBorderL.Accept(time2, []time.Time{time1, time3}))
	assert.True(t, TimeInBorderR.Accept(time2, []time.Time{time1, time3}))
	assert.True(t, TimeInBorderN.Accept(time2, []time.Time{time1, time3}))
}

// 补充测试op_word.go中的String函数
func TestWordOpString(t *testing.T) {
	// 测试有效操作符
	assert.Equal(t, "word/=", WordSet.String())
	assert.Equal(t, "word/empty", WordEmpty.String())
	assert.Equal(t, "word/any", WordAny.String())
	assert.Equal(t, "word/==", WordEq.String())
	assert.Equal(t, "word/!=", WordNe.String())
	assert.Equal(t, "word/include", WordIncld.String())
	assert.Equal(t, "word/exclude", WordExcld.String())
	assert.Equal(t, "word/one", WordOne.String())
	assert.Equal(t, "word/none", WordNone.String())
	assert.Equal(t, "word/all", WordAll.String())

	// 测试未知操作符
	assert.Equal(t, "", UnknownW.String())
}

// 补充测试op_word.go中的IsNeg函数
func TestWordOpIsNeg(t *testing.T) {
	// 测试否定操作符 (只检查以"no"开头的)
	assert.True(t, WordNone.IsNeg())

	// 测试非否定操作符
	assert.False(t, WordNe.IsNeg())
	assert.False(t, WordExcld.IsNeg())
	assert.False(t, WordOne.IsNeg())
	assert.False(t, WordAll.IsNeg())
	assert.False(t, WordEq.IsNeg())
	assert.False(t, WordIncld.IsNeg())
	assert.False(t, WordEmpty.IsNeg())
	assert.False(t, WordAny.IsNeg())
	assert.False(t, WordSet.IsNeg())
}

// 补充测试word_score.go中的RegSimilarAlgo函数
func TestRegSimilarAlgo(t *testing.T) {
	// 保存原始状态以便恢复
	defer func() {
		similarScorers["test_algo"] = nil
	}()

	// 定义一个自定义评分器
	customScorer := func(txt, key string, runes []rune, length int, options jsonx.JObj) float64 {
		return 0.75
	}

	// 注册自定义评分器
	RegSimilarAlgo("test_algo", customScorer)

	// 验证注册是否成功
	scorer := GetSimScorer("test_algo")
	result := scorer("test", "key", nil, 0, jsonx.JObj{})
	assert.Equal(t, 0.75, result)
}

// 补充测试word_score.go中的getLevenshteinWindowSize函数
func TestGetLevenshteinWindowSize(t *testing.T) {
	// 测试正常情况
	sizes := getLevenshteinWindowSize(20, 5, 0.5)
	expectedSizes := []int{5, 3, 10, 20} // 5*0.5=2.5→3, 5/0.5=10
	assert.ElementsMatch(t, expectedSizes, sizes)

	// 测试tLen小于某些窗口大小的情况
	sizes = getLevenshteinWindowSize(8, 5, 0.5)
	expectedSizes = []int{5, 3, 8} // 10超过了tLen(8)，所以被过滤掉
	assert.ElementsMatch(t, expectedSizes, sizes)

	// 测试tLen等于窗口大小的情况
	sizes = getLevenshteinWindowSize(5, 5, 0.5)
	expectedSizes = []int{5, 3} // 10超过了tLen(5)，tLen(5)已经包含在窗口大小中
	assert.ElementsMatch(t, expectedSizes, sizes)

	// 测试ratio为1.0的情况
	sizes = getLevenshteinWindowSize(20, 5, 1.0)
	expectedSizes = []int{5, 5, 5, 20} // 应该去重
	assert.ElementsMatch(t, []int{5, 20}, sizes)

	// 测试ratio为0.1的情况（极小的ratio）
	sizes = getLevenshteinWindowSize(20, 5, 0.1)
	expectedSizes = []int{5, 1, 50, 20} // 50超过了tLen(20)，所以被过滤掉
	assert.ElementsMatch(t, []int{5, 1, 20}, sizes)
}

// 补充测试TrimSpaceAndRuneCount函数
func TestTrimSpaceAndRuneCount(t *testing.T) {
	// 测试空字符串
	str, count := TrimSpaceAndRuneCount("")
	assert.Equal(t, "", str)
	assert.Equal(t, 0, count)

	// 测试无空格字符串
	str, count = TrimSpaceAndRuneCount("hello")
	assert.Equal(t, "hello", str)
	assert.Equal(t, 5, count)

	// 测试有空格字符串
	str, count = TrimSpaceAndRuneCount(" hello world ")
	assert.Equal(t, "helloworld", str)
	assert.Equal(t, 10, count)

	// 测试包含多个空格和换行符的字符串
	str, count = TrimSpaceAndRuneCount("hello\nworld  test")
	assert.Equal(t, "helloworldtest", str)
	assert.Equal(t, 14, count) // "helloworldtest" 是14个字符

	// 测试中文字符串
	str, count = TrimSpaceAndRuneCount("你好 世界")
	assert.Equal(t, "你好世界", str)
	assert.Equal(t, 4, count)
}

// 补充测试preLevenshtein函数
func TestPreLevenshtein(t *testing.T) {
	// 测试普通字符串
	target, word, tLen, wLen := preLevenshtein("hello world", "test")
	assert.Equal(t, "helloworld", target)
	assert.Equal(t, "test", word)
	assert.Equal(t, 10, tLen)
	assert.Equal(t, 4, wLen)

	// 测试有多个空格的字符串
	target, word, tLen, wLen = preLevenshtein("  hello   world  ", "  test  ")
	assert.Equal(t, "helloworld", target)
	assert.Equal(t, "test", word)
	assert.Equal(t, 10, tLen)
	assert.Equal(t, 4, wLen)

	// 测试中文字符串
	target, word, tLen, wLen = preLevenshtein("你好 世界", "测试")
	assert.Equal(t, "你好世界", target)
	assert.Equal(t, "测试", word)
	assert.Equal(t, 4, tLen)
	assert.Equal(t, 2, wLen)
}

// 补充测试calcLevenshteinWindow函数
func TestCalcLevenshteinWindow(t *testing.T) {
	// 准备测试数据
	target := []rune("abcdefg")
	tLen := 7
	word := "def"
	wLen := 3

	// 测试窗口大小为3
	score := calcLevenshteinWindow(target, word, tLen, wLen, 3)
	assert.Equal(t, 1.0, score) // 完全匹配

	// 测试窗口大小为4
	score = calcLevenshteinWindow(target, word, tLen, wLen, 4)
	assert.Equal(t, 1.0-score, 0.25) // 应该有约0.75的分数

	// 测试不存在的模式
	score = calcLevenshteinWindow(target, "xyz", tLen, 3, 3)
	assert.Equal(t, 0.0, score) // 完全不匹配

	// 测试中文字符
	target = []rune("你好世界")
	tLen = 4
	score = calcLevenshteinWindow(target, "世界", tLen, 2, 2)
	assert.Equal(t, 1.0, score) // 完全匹配
}

// 测试WordOp的maxArgs函数的所有分支
func TestWordOpMaxArgs(t *testing.T) {
	assert.Equal(t, 1, WordEmpty.maxArgs())
	assert.Equal(t, 1, WordAny.maxArgs())
	assert.Equal(t, 2, WordEq.maxArgs())
	assert.Equal(t, 2, WordIncld.maxArgs())
	assert.Equal(t, 2, WordNe.maxArgs())
	assert.Equal(t, 2, WordExcld.maxArgs())
	assert.Equal(t, -1, WordAll.maxArgs())
	assert.Equal(t, -1, WordOne.maxArgs())
	assert.Equal(t, -1, WordNone.maxArgs())
	assert.Equal(t, 0, WordSet.maxArgs())
}

// 测试AsTimeOp函数的边界情况
func TestAsTimeOp(t *testing.T) {
	// 测试已知操作符
	assert.Equal(t, TimeLt, AsTimeOp("<"))
	assert.Equal(t, TimeLt, AsTimeOp("<"))
	assert.Equal(t, TimeGt, AsTimeOp(">"))
	assert.Equal(t, TimeIn, AsTimeOp("[a,b]"))

	// 测试大小写不敏感
	assert.Equal(t, TimeLt, AsTimeOp("<"))
	assert.Equal(t, TimeIn, AsTimeOp("[A,B]"))

	// 测试未知操作符
	assert.Equal(t, UnknownT, AsTimeOp("unknown"))
	assert.Equal(t, UnknownT, AsTimeOp(""))
}
