package jsonx

import (
	"time"
)

type JNum float64

func (re JNum) Type() JType       { return JNUM }
func (re JNum) IsNull() bool      { return false }
func (re JNum) ToInt() int        { return int(re.ToLong()) }
func (re JNum) ToTime() time.Time { return AsTime(re.ToLong()) }
func (re JNum) ToLong() int64     { return int64(re) }
func (re JNum) ToDouble() float64 { return float64(re) }
func (re JNum) ToFloat() float32  { return float32(re.ToDouble()) }
func (re JNum) ToBool() bool      { return re.ToInt() > 0 }
func (re JNum) String() string    { return AsStr(re.ToDouble()) }
func (re JNum) Pretty() string    { return re.String() }
func (re JNum) ToObj() JObj       { return JObj{} }
func (re JNum) ToObjPtr() *JObj   { obj := re.ToObj(); return &obj }
func (re JNum) ToArr() JArr       { return JArr{} }
func (re JNum) ToJDoc() JDoc      { return JDoc(re.String()) }
func (re JNum) ToJVal() JValue    { return re }
func (re JNum) ToGVal() any       { return float64(re) }

func NewJNum(val float64) JNum { return JNum(val) }
