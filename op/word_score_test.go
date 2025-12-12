package op

import (
	"github.com/fengzhi09/golibx/gox"
	"github.com/fengzhi09/golibx/jsonx"
	"math"
	"testing"
	"unicode/utf8"

	"github.com/agnivade/levenshtein"
)

type wordScoreTest struct {
	name    string
	content string
	word    string
	options jsonx.JObj
	want    float64
}

var baseTests = []*wordScoreTest{
	{
		name: "both_empty", options: jsonx.JObj{"threshold": 0.4},
		want: 1, word: "", content: "",
	},
	{
		name: "word_empty", options: jsonx.JObj{"threshold": 0.4},
		want: 0, word: "", content: "1",
	},
	{
		name: "content_empty", options: jsonx.JObj{"threshold": 0.4},
		want: 0, word: "1", content: "",
	},
	{
		name: "eq_num", options: jsonx.JObj{"threshold": 0.4},
		want: 1, word: "1", content: "1",
	},
	{
		name: "eq_word", options: jsonx.JObj{"threshold": 0.4},
		want: 1, word: "中", content: "中",
	},
	{
		name: "eq_char", options: jsonx.JObj{"threshold": 0.4},
		want: 1, word: "A", content: "A",
	},
	{
		name: "key[12345]", options: jsonx.JObj{"threshold": 0.4},
		want: 0.5, word: "12345", content: "1234567890",
	},
	{
		name: "key[150]", options: jsonx.JObj{"threshold": 0.6},
		want: 0.3, word: "150", content: "1234567890",
	},
	{
		name: "key[145]", options: jsonx.JObj{"threshold": 0.6},
		want: 0.3, word: "145", content: "1234567890",
	},
	{
		name: "key[125]", options: jsonx.JObj{"threshold": 0.6},
		want: 0.3, word: "125", content: "1234567890",
	},
	{
		name: "key[124]", options: jsonx.JObj{"threshold": 0.6},
		want: 0.3, word: "124", content: "1234567890",
	},
	{
		name: "key[235]", options: jsonx.JObj{"threshold": 0.6},
		want: 0.3, word: "235", content: "1234567890",
	},
	{
		name: "key[234]", options: jsonx.JObj{"threshold": 0.6},
		want: 0.3, word: "234", content: "1234567890",
	},
	{
		name: "key[156]", options: jsonx.JObj{"threshold": 0.6},
		want: 0.3, word: "156", content: "1234567890",
	},
	{
		name: "key[12A]", options: jsonx.JObj{"threshold": 0.6},
		want: 0.2, word: "12A", content: "1234567890",
	},
	{
		name: "key[15A]", options: jsonx.JObj{"threshold": 0.6},
		want: 0.2, word: "15A", content: "1234567890",
	},
	{
		name: "key[16A]", options: jsonx.JObj{"threshold": 0.6},
		want: 0.2, word: "16A", content: "1234567890",
	},
	{
		name: "key[24A]", options: jsonx.JObj{"threshold": 0.6},
		want: 0.2, word: "24A", content: "1234567890",
	},
	{
		name: "key[123]", options: jsonx.JObj{"threshold": 0.6},
		want: 0.3, word: "123", content: "1234567890",
	},
	{
		name: "key[135]", options: jsonx.JObj{"threshold": 0.6},
		want: 0.3, word: "135", content: "1234567890",
	},
	{
		name: "key[10]", options: jsonx.JObj{"threshold": 0.4},
		want: 0.2, word: "10", content: "1234567890",
	},
	{
		name: "key[16]", options: jsonx.JObj{"threshold": 0.4},
		want: 0.2, word: "16", content: "1234567890",
	},
	{
		name: "key[15]", options: jsonx.JObj{"threshold": 0.4},
		want: 0.2, word: "15", content: "1234567890",
	},
	{
		name: "key[14]", options: jsonx.JObj{"threshold": 0.4},
		want: 0.2, word: "14", content: "1234567890",
	},
	{
		name: "key[13]", options: jsonx.JObj{"threshold": 0.4},
		want: 0.2, word: "13", content: "1234567890",
	},
	{
		name: "key[12]", options: jsonx.JObj{"threshold": 0.4},
		want: 0.2, word: "13", content: "1234567890",
	},
	{
		name: "key_no", options: jsonx.JObj{"threshold": 0.4},
		want: 0, word: "AB", content: "1234567890",
	},
}

func Test_LevenshteinDistance(t *testing.T) {
	type ldCase struct {
		txt  string
		key  string
		want int
	}
	tests := []*ldCase{
		{txt: "", key: "", want: 0},
		{txt: "1", key: "", want: 1},
		{txt: "12", key: "21", want: 2},
		{txt: "123", key: "231", want: 2},
		{txt: "1_23", key: "123", want: 1},
		{txt: "喜欢", key: "欢喜", want: 2},
		{txt: "24A", key: "234", want: 2},
		{txt: "24A", key: "12345", want: 3},
		{txt: "243", key: "234", want: 2},
	}
	for _, test := range tests {
		got := levenshtein.ComputeDistance(test.txt, test.key)
		if got != test.want {
			t.Errorf("case[%v,%v] dis=%v but want=%v", test.txt, test.key, got, test.want)
		} else {
			score := LevenshteinScore(test.txt, test.key, []rune{}, 0, jsonx.JObj{"threshold": 0.4})
			t.Logf("case[%v,%v] dis=%v score=%.2f", test.txt, test.key, got, score)
		}
	}
}

func Test_LevenshteinScorer(t *testing.T) {
	for _, test := range baseTests {
		t.Run(test.name, func(tt *testing.T) {
			runWordScoreTest(tt, test, test.want, LevenshteinScore)
		})
	}
	windowScoreSpec := map[string]float64{
		"key[12345]": 1, "key[12]": 0.5, "key[13]": 0.5, "key[14]": 0.5, "key[15]": 0.5, "key[16]": 0.5, "key[10]": 0.5,
		"key[123]": 1, "key[124]": 0.67, "key[150]": 0.33, "key[156]": 0.67, "key[145]": 0.67, "key[125]": 0.67,
		"key[135]": 0.60, "key[235]": 0.67, "key[234]": 1, "key[12A]": 0.67, "key[15A]": 0.33, "key[16A]": 0.33, "key[24A]": 0.40,
	}
	for _, test := range baseTests {
		t.Run(test.name+"[window]", func(tt *testing.T) {
			spec, hit := windowScoreSpec[test.name]
			want := gox.IfElse(hit, spec, test.want).(float64)
			runWordScoreTest(tt, test, want, LevenshteinWindowScore)
		})
	}
}

func Test_LevenshteinCostCompile(t *testing.T) {
	//scorer0, score1 := LevenshteinScore, LevenshteinWindowScore
}

func runWordScoreTest(t *testing.T, wst *wordScoreTest, want float64, scorer SimilarScorer) {
	runes := []rune(wst.content)
	length := utf8.RuneCountInString(wst.content)
	got := scorer(wst.content, wst.word, runes, length, wst.options)
	if got < 0 || got > 1 {
		t.Errorf("%v got=%.2f not in [0,1]", wst.name, got)
	}
	if want != 0 && math.Abs(got-want) > 0.01 {
		t.Errorf("%v got=%.2f but want=%.2f", wst.name, got, want)
	}
	t.Logf("%v got=%.2f, want=%.2f word=%v content=%v", wst.name, got, want, wst.word, wst.content)
}
