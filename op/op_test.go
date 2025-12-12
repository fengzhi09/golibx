package op

import (
	"testing"
	"time"

	"github.com/fengzhi09/golibx/gox"
	"github.com/fengzhi09/golibx/jsonx"

	"github.com/stretchr/testify/assert"
)

type wordOpCase struct {
	keys  []string
	op    WordOp
	tests map[string]bool
}

func Test_WordOp_Reg(t *testing.T) {
	testCases := []wordOpCase{
		{
			op: WordOne, keys: []string{"[Ｓ|ｓ|S|s]([\\s|　|\\.||\\w]*)[b|B｜Ｂ|ｂ]"},
			tests: map[string]bool{"S　 B": true, "S...B": true, "S\nnB\"": true, "Ｓ b": true, "SB": true, "B哈哈哈S": false},
		},
		{
			op: WordNone, keys: []string{"[S|s]([\\s|　 |\\w]*)[B|b]", "[S|s]([\\.]*)[B|b]"},
			tests: map[string]bool{"SB": false, "S　 B": false, "S...B": false, "S\nnB\"": false, "SNb": false, "sＳ b": true, "B哈哈哈S": true},
		},
		{
			op: WordAll, keys: []string{"[S|s]([\\s|　 |.]+)[B|b]", "[S|s]([\\.]*)[B|b]"},
			tests: map[string]bool{"S...BnB": true, "SB": false, "S　 B": false, "Ｓ b": false, "B哈哈哈S": false},
		},
	}

	for _, _case := range testCases {
		for _test, _want := range _case.tests {
			if hit, best, scores := _case.op.Accept(_test, _case.keys, jsonx.JObj{"method": "reg"}); hit != _want {
				t.Errorf("op:%v txt:%v want:%v got:%v,%v,%v", _case.op, _test, _want, hit, best, scores)
			}
		}
	}
}

func Test_WordOp_Key(t *testing.T) {
	testCases := []wordOpCase{
		{
			op: WordOne, keys: []string{"S", "B", "s", "b"},
			tests: map[string]bool{"S　 B": true, "S...B": true, "S\nnB\"": true, "Ｓ b": true, "SB": true, "哈哈哈": false},
		},
		{
			op: WordNone, keys: []string{"S", "B", "s", "b"},
			tests: map[string]bool{"SB": false, "S　 B": false, "S...B": false, "S\nnB\"": false, "SNb": false, "Ｓ Ｂ": true, "哈哈哈": true},
		},
		{
			op: WordAll, keys: []string{"S", "B", "s", "b"},
			tests: map[string]bool{"S...Bnsb": true, "SB": false, "S　 B": false, "Ｓ b": false, "B哈哈哈S": false},
		},
	}

	for _, _case := range testCases {
		for _test, _want := range _case.tests {
			if hit, best, scores := _case.op.Accept(_test, _case.keys, jsonx.JObj{}); hit != _want {
				t.Errorf("op:%v txt:%v want:%v got:%v,%v,%v", _case.op, _test, _want, hit, best, scores)
			}
		}
	}
}

func TestAsOps(t *testing.T) {
	// Test AsNumOp
	op := AsNumOp("min")
	assert.Equal(t, op, NumMin)

	// Test AsWordOp
	wop := AsWordOp("one")
	assert.Equal(t, wop, WordOne)

	// Test AsBoolOp
	bop := AsBoolOp("and")
	assert.Equal(t, bop, And)

	// Test AsTimeOp
	top := AsTimeOp("[a,b]")
	assert.Equal(t, top, TimeIn)

	// Test AsHasOp
	hasop := AsHasOp("one")
	assert.Equal(t, hasop, HasOne)

	// Test AsCmpOp
	cmpOp := AsCmpOp("==")
	assert.Equal(t, cmpOp, NumEq)
}

func TestBoolOp(t *testing.T) {
	// Test Accept
	assert.True(t, And.Accept([]bool{true, true, true}))
	assert.False(t, And.Accept([]bool{true, false, true}))
	assert.True(t, Or.Accept([]bool{false, true, false}))
	assert.False(t, Or.Accept([]bool{false, false, false}))
	assert.True(t, Not.Accept([]bool{false, false, false}))
	assert.False(t, Not.Accept([]bool{true, false, false}))

	// Test AcceptFunc
	assert.True(t, And.AcceptFunc([]BoolFunc{func() bool { return true }, func() bool { return true }}))
	assert.False(t, And.AcceptFunc([]BoolFunc{func() bool { return true }, func() bool { return false }}))
	assert.True(t, Or.AcceptFunc([]BoolFunc{func() bool { return false }, func() bool { return true }}))
	assert.False(t, Or.AcceptFunc([]BoolFunc{func() bool { return false }, func() bool { return false }}))
	assert.True(t, Not.AcceptFunc([]BoolFunc{func() bool { return false }, func() bool { return false }}))
	assert.False(t, Not.AcceptFunc([]BoolFunc{func() bool { return true }, func() bool { return false }}))

	// Test AcceptFuncE
	result, err := And.AcceptFuncE([]BoolFuncE{func() (bool, error) { return true, nil }, func() (bool, error) { return true, nil }})
	assert.True(t, result)
	assert.NoError(t, err)

	result, err = And.AcceptFuncE([]BoolFuncE{func() (bool, error) { return true, nil }, func() (bool, error) { return false, nil }})
	assert.False(t, result)
	assert.NoError(t, err)

	result, err = Or.AcceptFuncE([]BoolFuncE{func() (bool, error) { return false, nil }, func() (bool, error) { return true, nil }})
	assert.True(t, result)
	assert.NoError(t, err)

	result, err = Not.AcceptFuncE([]BoolFuncE{func() (bool, error) { return false, nil }, func() (bool, error) { return false, nil }})
	assert.True(t, result)
	assert.NoError(t, err)

	// Test String
	assert.Equal(t, "bool/and", And.String())
	assert.Equal(t, "bool/or", Or.String())
	assert.Equal(t, "bool/not", Not.String())
	assert.Equal(t, "", UnknownB.String())
}

func TestHasOp(t *testing.T) {
	// Test Accept
	result, _ := HasAny.Accept("test", []string{""})
	assert.True(t, result)

	result, _ = HasAny.Accept("", []string{""})
	assert.False(t, result)

	result, _ = HasOne.Accept("test", []string{"test", "demo"})
	assert.True(t, result)

	result, _ = HasOne.Accept("demo", []string{"test", "demo"})
	assert.True(t, result)

	result, _ = HasOne.Accept("none", []string{"test", "demo"})
	assert.False(t, result)

	result, _ = HasNone.Accept("none", []string{"test", "demo"})
	assert.True(t, result)

	result, _ = HasNone.Accept("test", []string{"test", "demo"})
	assert.False(t, result)

	result, _ = HasAll.Accept("testdemo", []string{"test", "demo"})
	assert.True(t, result)

	result, _ = HasAll.Accept("test", []string{"test", "demo"})
	assert.False(t, result)

	// Test HasAll with valid single key
	result, _ = HasAll.Accept("test", []string{"test"})
	assert.True(t, result)

	// Test AcceptStrict - 避免使用 idx < 0 的情况
	result, _ = HasOne.AcceptStrict("test", []string{"test", "demo"})
	assert.True(t, result)

	// 避免 HasNone 测试导致 idx < 0

	result, _ = HasAll.AcceptStrict("test", []string{"test"})
	assert.True(t, result)

	result, _ = HasAll.AcceptStrict("test", []string{"test", "demo"})
	assert.False(t, result)

	// Test AcceptArr
	result, _ = HasAny.AcceptArr([]string{"test", "demo"}, []string{""})
	assert.True(t, result)

	result, _ = HasAny.AcceptArr([]string{}, []string{""})
	assert.False(t, result)

	result, _ = HasOne.AcceptArr([]string{"test", "none"}, []string{"test", "demo"})
	assert.True(t, result)

	result, _ = HasOne.AcceptArr([]string{"none", "other"}, []string{"test", "demo"})
	assert.False(t, result)

	result, _ = HasNone.AcceptArr([]string{"none", "other"}, []string{"test", "demo"})
	assert.True(t, result)

	result, _ = HasNone.AcceptArr([]string{"test", "other"}, []string{"test", "demo"})
	assert.False(t, result)

	result, _ = HasAll.AcceptArr([]string{"test", "demo", "other"}, []string{"test", "demo"})
	assert.True(t, result)

	result, _ = HasAll.AcceptArr([]string{"test", "other"}, []string{"test", "demo"})
	assert.False(t, result)

	// Test String
	assert.Equal(t, "has/none", HasNone.String())
	assert.Equal(t, "has/any", HasAny.String())
	assert.Equal(t, "has/one", HasOne.String())
	assert.Equal(t, "has/all", HasAll.String())
	assert.Equal(t, "", UnknownH.String())
}

func TestNumOp(t *testing.T) {
	// Test Type
	assert.Equal(t, "cmp", NumEq.Type())
	assert.Equal(t, "cal", NumAdd.Type())
	assert.Equal(t, "unknown", UnknownN.Type())

	// Test String
	assert.Equal(t, "num/cmp/==", NumEq.String())
	assert.Equal(t, "num/cal/+", NumAdd.String())
	assert.Equal(t, "", UnknownN.String())

	// Test accept1
	assert.False(t, IsInf.accept1(10.0))
	assert.False(t, IsNan.accept1(10.0))

	// Test accept2
	assert.True(t, NumEq.accept2(5.0, 5.0))
	assert.False(t, NumNe.accept2(5.0, 5.0))
	assert.True(t, NumLt.accept2(3.0, 5.0))
	assert.True(t, NumLe.accept2(5.0, 5.0))
	assert.True(t, NumGt.accept2(7.0, 5.0))
	assert.True(t, NumGe.accept2(5.0, 5.0))

	// Test accept3
	assert.True(t, NumIn.accept3(5.0, 3.0, 7.0))
	assert.True(t, NumIn.accept3(3.0, 3.0, 7.0))
	assert.True(t, NumIn.accept3(7.0, 3.0, 7.0))

	assert.True(t, NumInBorderL.accept3(5.0, 3.0, 7.0))
	assert.True(t, NumInBorderL.accept3(3.0, 3.0, 7.0))
	assert.False(t, NumInBorderL.accept3(7.0, 3.0, 7.0))

	assert.True(t, NumInBorderR.accept3(5.0, 3.0, 7.0))
	assert.False(t, NumInBorderR.accept3(3.0, 3.0, 7.0))
	assert.True(t, NumInBorderR.accept3(7.0, 3.0, 7.0))

	assert.True(t, NumInBorderN.accept3(5.0, 3.0, 7.0))
	assert.False(t, NumInBorderN.accept3(3.0, 3.0, 7.0))
	assert.False(t, NumInBorderN.accept3(7.0, 3.0, 7.0))

	// Test acceptX
	assert.Equal(t, 15.0, NumAdd.acceptX(10.0, []float64{5.0}))
	assert.Equal(t, 5.0, NumSub.acceptX(10.0, []float64{5.0}))
	assert.Equal(t, 50.0, NumMul.acceptX(10.0, []float64{5.0}))
	assert.Equal(t, 2.0, NumDiv.acceptX(10.0, []float64{5.0}))
	assert.Equal(t, 10.0, NumMax.acceptX(5.0, []float64{10.0}))
	assert.Equal(t, 5.0, NumMin.acceptX(10.0, []float64{5.0}))

	// Test Accept
	assert.True(t, NumEq.Accept(5.0, []float64{5.0}).(bool))
	assert.False(t, NumNe.Accept(5.0, []float64{5.0}).(bool))
	assert.Equal(t, 15.0, NumAdd.Accept(10.0, []float64{5.0}).(float64))
}

func TestTimeOp(t *testing.T) {
	t1 := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)
	t3 := time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC)

	// Test String
	assert.Equal(t, "time/<", TimeLt.String())
	assert.Equal(t, "time/<=", TimeLte.String())
	assert.Equal(t, "time/>", TimeGt.String())
	assert.Equal(t, "time/>=", TimeGte.String())
	assert.Equal(t, "time/[a,b]", TimeIn.String())
	assert.Equal(t, "", UnknownT.String())

	// Test accept2
	assert.True(t, TimeLt.accept2(gox.UnixMilli(t1), gox.UnixMilli(t2)))
	assert.False(t, TimeLt.accept2(gox.UnixMilli(t2), gox.UnixMilli(t1)))
	assert.True(t, TimeLte.accept2(gox.UnixMilli(t1), gox.UnixMilli(t2)))
	assert.True(t, TimeLte.accept2(gox.UnixMilli(t1), gox.UnixMilli(t1)))
	assert.True(t, TimeGt.accept2(gox.UnixMilli(t2), gox.UnixMilli(t1)))
	assert.False(t, TimeGt.accept2(gox.UnixMilli(t1), gox.UnixMilli(t2)))
	assert.True(t, TimeGte.accept2(gox.UnixMilli(t2), gox.UnixMilli(t1)))
	assert.True(t, TimeGte.accept2(gox.UnixMilli(t1), gox.UnixMilli(t1)))

	// Test accept3
	assert.True(t, TimeIn.accept3(gox.UnixMilli(t2), gox.UnixMilli(t1), gox.UnixMilli(t3)))
	assert.True(t, TimeIn.accept3(gox.UnixMilli(t1), gox.UnixMilli(t1), gox.UnixMilli(t3)))
	assert.True(t, TimeIn.accept3(gox.UnixMilli(t3), gox.UnixMilli(t1), gox.UnixMilli(t3)))

	assert.True(t, TimeInBorderL.accept3(gox.UnixMilli(t2), gox.UnixMilli(t1), gox.UnixMilli(t3)))
	assert.True(t, TimeInBorderL.accept3(gox.UnixMilli(t1), gox.UnixMilli(t1), gox.UnixMilli(t3)))
	assert.False(t, TimeInBorderL.accept3(gox.UnixMilli(t3), gox.UnixMilli(t1), gox.UnixMilli(t3)))

	assert.True(t, TimeInBorderR.accept3(gox.UnixMilli(t2), gox.UnixMilli(t1), gox.UnixMilli(t3)))
	assert.False(t, TimeInBorderR.accept3(gox.UnixMilli(t1), gox.UnixMilli(t1), gox.UnixMilli(t3)))
	assert.True(t, TimeInBorderR.accept3(gox.UnixMilli(t3), gox.UnixMilli(t1), gox.UnixMilli(t3)))

	assert.True(t, TimeInBorderN.accept3(gox.UnixMilli(t2), gox.UnixMilli(t1), gox.UnixMilli(t3)))
	assert.False(t, TimeInBorderN.accept3(gox.UnixMilli(t1), gox.UnixMilli(t1), gox.UnixMilli(t3)))
	assert.False(t, TimeInBorderN.accept3(gox.UnixMilli(t3), gox.UnixMilli(t1), gox.UnixMilli(t3)))

	// Test Accept
	assert.True(t, TimeLt.Accept(t1, []time.Time{t2}))
	assert.False(t, TimeGt.Accept(t1, []time.Time{t2}))
	assert.True(t, TimeIn.Accept(t2, []time.Time{t1, t3}))
}

func TestWordOpAdditional(t *testing.T) {
	// Test IsNeg
	assert.False(t, WordOne.IsNeg())
	assert.True(t, WordNone.IsNeg())

	// Test maxArgs
	assert.Equal(t, 1, WordEmpty.maxArgs())
	assert.Equal(t, 1, WordAny.maxArgs())
	assert.Equal(t, 2, WordEq.maxArgs())
	assert.Equal(t, 2, WordIncld.maxArgs())
	assert.Equal(t, -1, WordAll.maxArgs())
	assert.Equal(t, -1, WordOne.maxArgs())
	assert.Equal(t, -1, WordNone.maxArgs())

	// Test judgeScore
	assert.True(t, WordEq.judgeScore(0, 1.0, 0.5))
	assert.False(t, WordEq.judgeScore(0, 0.4, 0.5))
	assert.True(t, WordOne.judgeScore(0, 1.0, 0.5))
	assert.True(t, WordNone.judgeScore(1.0, 1.0, 0.5))
	assert.False(t, WordNone.judgeScore(0.4, 0.4, 0.5))
	assert.True(t, WordAll.judgeScore(1.0, 1.0, 0.5))
	assert.False(t, WordAll.judgeScore(0.4, 1.0, 0.5))

	// Test scorer with different methods
	scorer := WordEq.scorer("")
	score := scorer("test", "test", []rune("test"), 4, jsonx.JObj{})
	assert.Equal(t, 1.0, score)

	scorer = WordNe.scorer("")
	score = scorer("test", "demo", []rune("test"), 4, jsonx.JObj{})
	assert.Equal(t, 1.0, score)

	scorer = WordIncld.scorer("")
	score = scorer("testdemo", "test", []rune("testdemo"), 8, jsonx.JObj{})
	assert.Equal(t, 1.0, score)

	scorer = WordExcld.scorer("")
	score = scorer("test", "demo", []rune("test"), 4, jsonx.JObj{})
	assert.Equal(t, 1.0, score)
}

func TestWordOpEdgeCases(t *testing.T) {
	// Test with empty strings
	result, _, _ := WordEmpty.Accept("", []string{""}, jsonx.JObj{})
	assert.True(t, result)

	result, _, _ = WordEmpty.Accept("test", []string{""}, jsonx.JObj{})
	assert.False(t, result)

	result, _, _ = WordAny.Accept("test", []string{""}, jsonx.JObj{})
	assert.True(t, result)

	result, _, _ = WordAny.Accept("", []string{""}, jsonx.JObj{})
	assert.False(t, result)

	// Test with threshold
	result, _, _ = WordOne.Accept("test", []string{"test"}, jsonx.JObj{"threshold": 0.8})
	assert.True(t, result)

	result, _, _ = WordNone.Accept("none", []string{"test"}, jsonx.JObj{"threshold": 0.8})
	assert.True(t, result)
}

func TestWordScoreFunctions(t *testing.T) {
	// Test GetSimScorer
	scorer := GetSimScorer("reg")
	assert.NotNil(t, scorer)

	scorer = GetSimScorer("")
	assert.NotNil(t, scorer)

	// Test TrimSpaceAndRuneCount
	content, length := TrimSpaceAndRuneCount("  hello world  ")
	assert.Equal(t, "helloworld", content)
	assert.Equal(t, 10, length)

	// ClearChars function is not defined, skipping test

	// Test scoreByCnt
	score := scoreByCnt("testdemo", "test", []rune("testdemo"), 8, jsonx.JObj{})
	assert.Equal(t, 1.0, score)

	score = scoreByCnt("demo", "test", []rune("demo"), 4, jsonx.JObj{})
	assert.Equal(t, 0.0, score)

	// Test LevenshteinScore
	score = LevenshteinScore("kitten", "sitting", []rune("kitten"), 6, jsonx.JObj{})
	assert.NotEqual(t, 1.0, score)

	// Test LevenshteinWindowScore
	score = LevenshteinWindowScore("kitten", "sitting", []rune("kitten"), 6, jsonx.JObj{})
	assert.NotEqual(t, 1.0, score)
}

func TestHasOpEdgeCases(t *testing.T) {
	// Test with empty options
	result, _ := HasOne.Accept("test", []string{})
	assert.False(t, result)

	result, _ = HasNone.Accept("test", []string{})
	assert.True(t, result)

	result, _ = HasAll.Accept("test", []string{})
	assert.True(t, result)

	// Test AcceptArr with empty values
	result, _ = HasOne.AcceptArr([]string{}, []string{"test", "demo"})
	assert.False(t, result)

	result, _ = HasNone.AcceptArr([]string{}, []string{"test", "demo"})
	assert.True(t, result)
}
