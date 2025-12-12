package jsonx

import (
	"time"
)

type JInt int64

func (re JInt) Type() JType  { return JINT }
func (re JInt) IsNull() bool { return false }
func (re JInt) ToInt() int   { return int(re) }
func (re JInt) ToTime() time.Time {
	return AsTime(re.ToLong())
}
func (re JInt) ToLong() int64     { return int64(re) }
func (re JInt) ToDouble() float64 { return float64(re) }
func (re JInt) ToFloat() float32  { return float32(re.ToDouble()) }
func (re JInt) ToBool() bool      { return re.ToInt() > 0 }
func (re JInt) String() string    { return AsStr(re.ToLong()) }
func (re JInt) Pretty() string    { return re.String() }
func (re JInt) ToObj() JObj       { return JObj{} }
func (re JInt) ToObjPtr() *JObj   { obj := re.ToObj(); return &obj }
func (re JInt) ToArr() JArr       { return JArr{} }
func (re JInt) ToJDoc() JDoc      { return JDoc(re.String()) }
func (re JInt) ToJVal() JValue    { return re }
func (re JInt) ToGVal() any       { return float64(re) }

func NewJInt(val int64) JInt { return JInt(val) }
