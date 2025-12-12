package gox

import (
	"testing"
)

func TestRound(t *testing.T) {
	tCases := map[int][]string{
		-2: {"1.2", "0", "18.4", "0", "100", "1"},
		-1: {"1.2", "0", "18.4", "2", "100", "10", "14.4", "1", "10.8", "1"},
		0:  {"1.2", "1", "18.4", "18", "10.8", "11"},
		1:  {"1.2", "1.2", "18.4", "18.4", "10.89", "10.9"},
		2:  {"1.2", "1.20", "18.4", "18.40", "10.89", "10.89"},
	}
	for pre, strs := range tCases {
		for i := 0; i < len(strs)/2; i++ {
			raw, want := strs[2*i], strs[2*i+1]
			got := RoundS(raw, pre)

			if want != got {
				t.Errorf("%v rounds %v got %v but want=%v", raw, pre, got, want)
			}
			t.Logf("%v rounds %v is %v", raw, pre, got)
		}
	}
	for pre, strs := range tCases {
		for i := 0; i < len(strs)/2; i++ {
			raw, want := strs[2*i], strs[2*i+1]
			got := RoundD(AsDouble(raw), pre)

			if AsDouble(want) != got {
				t.Errorf("%v roundv %v got %v but want=%v", raw, pre, got, want)
			}
			t.Logf("%v roundv %v is %v", raw, pre, got)
		}
	}
}

func TestFill(t *testing.T) {
	tCases := map[int][]string{
		5: {"123", "00123", "12300"},
		6: {"0123", "000123", "012300"},
		7: {"0.123", "000.123", "0.12300"},
	}
	for pre, strs := range tCases {
		raw, lWant, rWant := strs[0], strs[1], strs[2]
		lGot, rGot := FillLeft(raw, "0", pre), FillRight(raw, "0", pre)
		if lWant != lGot {
			t.Errorf("%v FillLeft got %v but want=%v", raw, lGot, lWant)
		}
		if rWant != rGot {
			t.Errorf("%v FillRight got %v but want=%v", raw, rGot, rWant)
		}
		t.Logf("%v fill got left:%v right:%v", raw, lGot, rGot)
	}
}

func TestAsMoney(t *testing.T) {
	tCases := map[string]int64{
		"18.4":       1840,
		"148.98":     14898,
		"123,148.98": 12314898,
		"1.84e1":     1840,
		"1.1184e3":   111840,
		".0184e3":    1840,
	}
	for raw, want := range tCases {
		gotBefore := int64(AsDouble(raw) * 100)
		gotAfter := AsMoney(raw, 2)
		if gotAfter != want {
			t.Errorf("%v got %v(before:%v) but want=%v", raw, gotAfter, gotBefore, want)
		} else {
			t.Logf("%v got %v(before:%v) equals want=%v", raw, gotAfter, gotBefore, want)
		}
	}
}
