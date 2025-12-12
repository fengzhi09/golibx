package op

import (
	"math"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/fengzhi09/golibx/gox"
	"github.com/fengzhi09/golibx/jsonx"

	"github.com/agnivade/levenshtein"
)

// SimilarScorer 比较的是源文本和单个目标词，得分范围应该在[0,1]之间;0表示无关/反面，1表示等价
type SimilarScorer func(txt string, key string, runes []rune, length int, options jsonx.JObj) float64

type SimilarAlgo = string

var similarScorers = map[SimilarAlgo]SimilarScorer{
	WordAlgoDefault:           scoreByCnt,
	WordAlgoReg:               scoreByReg,
	WordAlgoLevenshtein:       LevenshteinScore,
	WordAlgoLevenshteinWindow: LevenshteinWindowScore,
}

func GetSimScorer(algo SimilarAlgo) SimilarScorer {
	if scorer, ok := similarScorers[algo]; ok {
		return scorer
	}
	if algo == WordAlgoReg {
		return scoreByReg
	}
	return scoreByCnt
}

//goland:noinspection GoUnusedExportedFunction
func RegSimilarAlgo(method SimilarAlgo, calculator SimilarScorer) {
	similarScorers[method] = calculator
}

const (
	WordAlgoDefault SimilarAlgo = "default"
	WordAlgoReg     SimilarAlgo = "reg"
)

func scoreByCnt(a, b string, runes []rune, length int, options jsonx.JObj) float64 {
	return gox.IfElse(strings.Count(a, b) > 0, MaxScore, MinScore).(float64)
}

func scoreByReg(a, b string, runes []rune, length int, options jsonx.JObj) float64 {
	pattern := regexp.MustCompile(b)
	matches := pattern.FindAllString(a, -1)
	return gox.IfElse(len(matches) > 0, MaxScore, MinScore).(float64)
}

const WordAlgoLevenshtein SimilarAlgo = "levenshtein"

func preLevenshtein(content, keyword string) (string, string, int, int) {
	target := strings.Join(strings.Fields(content), "") // 去重空格
	word := strings.Join(strings.Fields(keyword), "")   // 去重空格
	tLen := utf8.RuneCountInString(target)
	wLen := utf8.RuneCountInString(word)
	return target, word, tLen, wLen
}

// get rune length and skip space
func TrimSpaceAndRuneCount(s string) (string, int) {
	str := strings.Join(strings.Fields(s), "")
	return str, utf8.RuneCountInString(str)
}

func LevenshteinScore(content, keyword string, runes []rune, length int, options jsonx.JObj) float64 {
	target, word, tLen, wLen := preLevenshtein(content, keyword)
	if target == word {
		return 1.0
	}
	dist := levenshtein.ComputeDistance(word, target)
	// 归一化分数公是:  1- 编辑距离/max(关键词长度, 目标词长度)
	return 1.0 - float64(dist)/math.Max(float64(tLen), float64(wLen))
}

const WordAlgoLevenshteinWindow SimilarAlgo = "levenshtein-window"

func LevenshteinWindowScore(content string, keyword string, target []rune, tLen int, options jsonx.JObj) float64 {
	word, wLen := TrimSpaceAndRuneCount(keyword)
	if content == word {
		return 1
	}
	ratio := gox.LimitIn(options.GetDouble("threshold"), 1e-6, 1.0)
	windowSizes := getLevenshteinWindowSize(tLen, wLen, ratio)
	scores := []float64{0.0}
	for _, window := range windowSizes {
		tmpScore := calcLevenshteinWindow(target, word, tLen, wLen, window)
		scores = append(scores, tmpScore)
	}
	return gox.MaxN(scores...)
}

/*
getLevenshteinWindowSize 计算窗口大小
输入：tLen 为目标长, wLen 为关键字长, ratio 为容忍比例(0-1)
输出：在以下视窗大小中去重，并舍弃超过tLen的

	标准视窗(wLen)
	最小视窗(wLen*ratio,向上取整)
	最大视窗(wLen/ratio,向上取整)
	无视窗(tLen)
*/
func getLevenshteinWindowSize(tLen, wLen int, ratio float64) []int {
	minWindow := math.Ceil(float64(wLen) * ratio)
	maxWindow := math.Ceil(float64(wLen) / ratio)
	preSizes := []int{wLen, int(minWindow), int(maxWindow), tLen}

	windowSizes := make([]int, 0)
	for _, preSize := range preSizes {
		if preSize <= tLen && preSize >= 1 {
			added := false
			for _, addedSize := range windowSizes {
				if addedSize == preSize {
					added = true
					break
				}
			}
			if !added {
				windowSizes = append(windowSizes, preSize)
			}
		}
	}
	return windowSizes
}

/*
calcLevenshteinWindow
按照window进行剪裁获取winTarget,同关键词计算编辑距离，返回归一化得分
*/
func calcLevenshteinWindow(target []rune, word string, tLen, wLen, window int) float64 {
	scores := []float64{0.0}
	for i := 0; i <= tLen-window; i++ {
		winTarget := gox.SubStrRune(target, tLen, i, i+window)
		dist := levenshtein.ComputeDistance(word, winTarget)
		if winTarget == word {
			dist = 0
		}
		tmpScore := 1.0 - float64(dist)/math.Max(float64(wLen), float64(window))
		scores = append(scores, tmpScore)
	}
	return gox.MaxN(scores...)
}
