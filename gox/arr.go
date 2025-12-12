package gox

import (
	"strings"
)

func ArrUniq[T comparable](arrs ...[]T) []T {
	tmp := map[T]bool{}
	res := make([]T, 0, len(arrs)*10)
	for _, arr := range arrs {
		for _, ele := range arr {
			if _, ok := tmp[ele]; !ok {
				tmp[ele] = true
				res = append(res, ele)
			}
		}
	}
	return res
}

func ArrSub[T comparable](items, skips []T) []T {
	tmp, res := map[T]bool{}, []T{}
	for _, skip := range skips {
		tmp[skip] = true
	}
	for _, item := range items {
		if _, hit := tmp[item]; hit {
			continue
		}
		res = append(res, item)
	}

	return res
}

func In[T any](val T, arr ...T) bool {
	return IndexOf(arr, val) >= 0
}

func IndexOf[T any](a []T, b T) int {
	alen := len(a)
	for i := 0; i < alen; i++ {
		if XEq(a[i], b) {
			return i
		}
	}
	return -1
}

func ArrStarts[T any](a, b []T) bool {
	alen, blen := len(a), len(b)
	if alen < blen {
		return false
	}

	cnt := 0
	for j := 0; j < blen; j++ {
		if !XEq(a[j], b[j]) {
			break
		}
		cnt = cnt + 1
	}
	return cnt == blen
}

func ArrEnds[T any](a, b []T) bool {
	alen, blen := len(a), len(b)
	if alen < blen {
		return false
	}

	cnt := 0
	for j := 0; j < blen; j++ {
		if !XEq(a[alen-blen+j], b[j]) {
			break
		}
		cnt = cnt + 1
	}
	return cnt == blen
}

func ArrIn[T any](a, b []T) bool {
	return ArrIndexOf(a, b) >= 0
}

func ArrIndexOf[T any](a, b []T) int {
	alen, blen := len(a), len(b)
	if alen < blen {
		return -1
	}

	for i := 0; i < alen; i++ {
		cnt := 0
		for j := 0; j < blen; j++ {
			if !XEq(a[i+j], b[j]) {
				break
			}
			cnt = cnt + 1
		}
		if cnt == blen {
			return i
		}
	}
	return -1
}

func WordIn(content string, words []string) int {
	for i, word := range words {
		if strings.Contains(content, word) {
			return i
		}
	}
	return -1
}

func StrAt(a string, b int) string {
	alen := len(a)
	if b >= 0 && b < alen {
		return string(a[b])
	}
	return ""
}

func StrAtMatch(dts string, idx int, opts ...string) bool {
	s := StrAt(dts, idx)
	for _, opt := range opts {
		if s == opt {
			return true
		}
	}
	return false
}
