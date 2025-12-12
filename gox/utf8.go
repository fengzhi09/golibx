package gox

import (
	"io"
	"unicode/utf8"
)

func IndexUtf8(target, keyword string) (int, int) {
	words, keys := []rune(target), []rune(keyword)
	for i := range words {
		match := true
		for j, key := range keys {
			if key != words[i+j] {
				match = false
				break
			}
		}
		if match {
			return i, len(keys)
		}
	}
	return -1, len(keys)
}

func LenUnt8(text string) int {
	return utf8.RuneCountInString(text)
}

func ReadRunes(r io.RuneReader) ([]rune, error) {
	runes := make([]rune, 0)
	for {
		r, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		runes = append(runes, r)
	}
	return runes, nil
}

func ToRunes(text string) []rune {
	return []rune(text)
}

func SubStr(content string, start, end int) string {
	words, substr := []uint8(content), make([]uint8, 0)
	for i, w := range words {
		if i >= start && (i < end || end < 0) {
			substr = append(substr, w)
		}
	}
	return string(substr)
}

func SubStrUtf8(content string, start, end int) string {
	words, substr := []rune(content), make([]rune, 0, len(content))
	for i, w := range words {
		if i >= start && (i < end || end < 0) {
			substr = append(substr, w)
		}
	}
	return string(substr)
}

func SubStrRune(content []rune, length, start, end int) string {
	if start >= length || (start >= end && end >= 0) {
		return ""
	}
	if end < 0 || end > length {
		end = length
	}

	substr := string(content[start:end])
	return substr
}
