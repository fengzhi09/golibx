package utils

import (
	"bytes"
	"testing"
)

func TestReaders(t *testing.T) {
	rawGBK := []byte{0xC4, 0xE3, 0xBA, 0xC3, 0xD6, 0xD0, 0xB9, 0xFA} // 你好中国 gbk
	rawUTF8 := []byte("你好中国")                                        // 你好中国 utf8

	t.Logf("gbk :=%v", string(rawGBK))
	t.Logf("utf8 :=%v", string(rawUTF8))
	t.Logf("gbk read gbk:=%v", Reader2String(GbkReader(bytes.NewReader(rawGBK))))
	t.Logf("gbk read utf8:=%v", Reader2String(GbkReader(bytes.NewReader(rawUTF8))))
	t.Logf("utf8 read gbk:=%v", Reader2String(Utf8Reader(bytes.NewReader(rawGBK))))
	t.Logf("utf8 read utf8:=%v", Reader2String(Utf8Reader(bytes.NewReader(rawUTF8))))
	utf82gbk, err := Utf82Gbk(string(rawUTF8))
	t.Logf("utf82gbk:=%v,%v", utf82gbk, err)
	gbk2utf8, err := Gbk2Utf8(string(rawGBK))
	t.Logf("gbk2utf8:=%v,%v", gbk2utf8, err)
}
