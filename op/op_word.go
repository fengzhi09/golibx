package op

import (
	"strings"

	"github.com/fengzhi09/golibx/gox"
	"github.com/fengzhi09/golibx/jsonx"
)

type WordOp string

//goland:noinspection GoNameStartsWithPackageName
const (
	WordSet   WordOp  = "="       // 赋值
	WordEmpty WordOp  = "empty"   // 空串
	WordAny   WordOp  = "any"     // 不空
	WordEq    WordOp  = "=="      // a与b相等
	WordNe    WordOp  = "!="      // a与b不等
	WordIncld WordOp  = "include" // a包含b
	WordExcld WordOp  = "exclude" // a不含b
	WordOne   WordOp  = "one"     // 包含某个关键字
	WordNone  WordOp  = "none"    // 不含任何关键字
	WordAll   WordOp  = "all"     // 包含所有关键字
	UnknownW  WordOp  = ""        // 未知
	MaxScore  float64 = 1.0
	MinScore  float64 = 0.0
)

var WordOps = []WordOp{WordSet, WordEmpty, WordAny, WordEq, WordNe, WordIncld, WordExcld, WordOne, WordNone, WordAll}

func AsWordOp(opStr string) WordOp {
	opStr = "word/" + strings.ToLower(opStr)
	for _, timeOp := range WordOps {
		if timeOp.String() == opStr {
			return timeOp
		}
	}
	return UnknownW
}

func (re WordOp) String() string {
	if re == "" {
		return ""
	}
	return "word/" + string(re)
}

func (re WordOp) Accept(src string, args []string, options jsonx.JObj) (bool, []string, []float64) {
	scorer := re.scorer(options.GetStr("method"))
	threshold := options.GetDouble("threshold")
	minScore, maxScore, maxArgCnt, bestHit, leastHit := -1.0, -1.0, re.maxArgs(), "", ""
	// 一次计算减少计算消耗
	content, length := TrimSpaceAndRuneCount(src)
	target := []rune(content)
	for i, key := range args {
		if maxArgCnt >= 0 && i >= maxArgCnt {
			break
		}
		score := scorer(content, key, target, length, options)
		if maxScore < 0 || score > maxScore {
			maxScore = score
			bestHit = key
		}
		if minScore < 0 || score < minScore {
			minScore = score
			leastHit = key
		}
	}
	pass := re.judgeScore(minScore, maxScore, threshold)
	hits := []string{leastHit, bestHit}
	scores := []float64{minScore, maxScore}
	return pass, hits, scores
}

func (re WordOp) maxArgs() int {
	switch re {
	case WordEmpty, WordAny:
		return 1 // 一元运算符
	case WordEq, WordIncld, WordNe, WordExcld:
		return 2 //　二元运算符
	case WordAll, WordOne, WordNone:
		return -1 //　不限制
	}
	return 0
}

func (re WordOp) judgeScore(minScore, maxScore, threshold float64) bool {
	threshold = gox.LimitIn(threshold, 1e-6, 1.0) // 门限值, 最低为0.0001 即万分之一；最高为1
	switch re {
	case WordEmpty, WordAny,
		WordNe, WordExcld, WordEq, WordIncld: // 只算1次,匹配就通过 等价于 最大值超过阈值
		return maxScore >= threshold
	case WordOne: // 多次计算有一次匹配就通过 等价于 最大值超过阈值
		return maxScore >= threshold
	case WordNone: // 多次计算所有都不匹配才通过(匹配计0,不匹配计1) 等价于 最小值不低于阈值
		return minScore >= threshold
	case WordAll: // 多次计算所有全都匹配才通过(匹配计1,不匹配计0) 等价于 最小值不低于阈值
		return minScore >= threshold
	}
	panic("op not accept1: " + re.String())
}

func (re WordOp) IsNeg() bool {
	return strings.HasPrefix(string(re), "no")
}

func (re WordOp) scorer(algo string) SimilarScorer {
	// 此时比较的是源文本和单个目标词
	return func(a string, b string, runes []rune, length int, options jsonx.JObj) float64 {
		switch re {
		case WordAny:
			return gox.IfElse(a != "", MaxScore, MinScore).(float64)
		case WordEmpty:
			return gox.IfElse(a == "", MaxScore, MinScore).(float64)
		case WordEq:
			return gox.IfElse(a == b, MaxScore, MinScore).(float64)
		case WordNe:
			return gox.IfElse(a != b, MaxScore, MinScore).(float64)
		case WordIncld:
			return scoreByCnt(a, b, runes, length, options) // 同算法类型无关
		case WordExcld:
			return 1 - scoreByCnt(a, b, runes, length, options) // 同算法类型无关
		case WordOne, WordAll:
			return GetSimScorer(algo)(a, b, runes, length, options)
		case WordNone: // 不含任一 取1-score；1的含义是认为一致
			return 1 - GetSimScorer(algo)(a, b, runes, length, options)
		}
		panic("op not accept1: " + re.String())
	}
}
