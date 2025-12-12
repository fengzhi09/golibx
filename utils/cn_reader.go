package utils

import (
	"io"

	"github.com/NuoMinMin/mahonia"
)

func GbkReader(reader io.Reader) io.Reader {
	return mahonia.NewDecoder("gbk").NewReader(reader)
}

func Utf8Reader(reader io.Reader) io.Reader {
	return mahonia.NewDecoder("utf-8").NewReader(reader)
}

func Utf82Gbk(str string) (string, error) {
	return convert2byte(str, "utf-8", "gbk")
}

func Gbk2Utf8(str string) (string, error) {
	return convert2byte(str, "gbk", "utf-8")
}

func convert2byte(src string, fromEncode, toEncode string) (string, error) {
	from := mahonia.NewDecoder(fromEncode)
	to := mahonia.NewDecoder(toEncode)
	_, translared, err := to.Translate([]byte(from.ConvertString(src)), true)
	return string(translared), err
}

func Utf82GbkUnsafe(str string) string {
	s, _ := Utf82Gbk(str)
	return s
}

func Gbk2Utf8Unsafe(str string) string {
	s, _ := Gbk2Utf8(str)
	return s
}

func Reader2String(reader io.Reader) string {
	s, e := io.ReadAll(reader)
	if e != nil {
		return e.Error()
	}
	return string(s)
}
