package httpx

import (
	"os"

	"github.com/fengzhi09/golibx/jsonx"
)

type ApiX interface {
	Get(path string, body any, opts ...HttpOpt) (int, *jsonx.JObj, error)
	PostTxt(path string, body any, opts ...HttpOpt) (int, *jsonx.JObj, error)
	PostForm(path string, body any, opts ...HttpOpt) (int, *jsonx.JObj, error)
	PostJson(path string, body any, opts ...HttpOpt) (int, *jsonx.JObj, error)
	Upload(path string, file *os.File, opts ...HttpOpt) (int, *jsonx.JObj, error)
	Download(urlPath string, savePath string, opts ...HttpOpt) (int, *jsonx.JObj, error)
}
