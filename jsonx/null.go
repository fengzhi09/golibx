package jsonx

import (
	"time"
)

type JNull struct {
}

func (re JNull) Type() JType       { return JNULL }
func (re JNull) IsNull() bool      { return true }
func (re JNull) ToInt() int        { return 0 }
func (re JNull) ToTime() time.Time { return AsTime(nil) }
func (re JNull) ToLong() int64     { return 0 }
func (re JNull) ToDouble() float64 { return 0 }
func (re JNull) ToFloat() float32  { return 0 }
func (re JNull) ToBool() bool      { return false }
func (re JNull) String() string    { return "" } //返回null会导致很多bug
func (re JNull) Pretty() string    { return re.String() }
func (re JNull) ToObj() JObj       { return JObj{} }
func (re JNull) ToObjPtr() *JObj   { return &JObj{} }
func (re JNull) ToArr() JArr       { return JArr{} }
func (re JNull) ToArrPtr() *JArr   { return &JArr{} }
func (re JNull) ToJDoc() JDoc      { return "null" }
func (re JNull) ToJVal() JValue    { return re }
func (re JNull) ToGVal() any       { return nil }

func NewJNull() JNull { return JNull{} }
func AtoJNull(val string) JNull {
	if val == "" || val == "null" || val == "nil" {
		return JNull{}
	}
	panic("bad for null:" + val)
}
