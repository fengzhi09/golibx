package jsonx

import (
	"strconv"
	"time"
)

type JStr string

func (re JStr) Type() JType       { return JSTR }
func (re JStr) IsNull() bool      { return false }
func (re JStr) ToInt() int        { val, _ := strconv.Atoi(re.String()); return val }
func (re JStr) ToTime() time.Time { return AsTime(re.String()) }
func (re JStr) ToLong() int64     { return int64(re.ToDouble()) }
func (re JStr) ToDouble() float64 { val, _ := strconv.ParseFloat(re.String(), 64); return val }
func (re JStr) ToFloat() float32  { return float32(re.ToDouble()) }
func (re JStr) ToBool() bool      { return re.ToInt() > 0 }
func (re JStr) String() string    { return string(re) }
func (re JStr) Pretty() string    { return re.String() }
func (re JStr) ToObj() JObj       { return ParseJObj(re.String()) }
func (re JStr) ToObjPtr() *JObj   { obj := re.ToObj(); return &obj }
func (re JStr) ToArr() JArr       { return ParseJArr(re.String()) }
func (re JStr) ToJDoc() JDoc      { return JDoc(UnsafeMarshalString(re)) }
func (re JStr) ToJVal() JValue    { return re }
func (re JStr) ToGVal() any       { return string(re) }

func NewJStr(str string) JStr { return JStr(str) }
