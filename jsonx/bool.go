package jsonx

import (
	"fmt"
	"strings"
	"time"
)

type JBool bool

func (re JBool) Type() JType       { return JBOOL }
func (re JBool) IsNull() bool      { return false }
func (re JBool) ToInt() int        { return int(re.ToDouble()) }
func (re JBool) ToTime() time.Time { return AsTime(nil) }
func (re JBool) ToLong() int64     { return int64(re.ToDouble()) }
func (re JBool) ToDouble() float64 { return IfElse(re.ToBool(), 1, 0).(float64) }
func (re JBool) ToFloat() float32  { return float32(re.ToDouble()) }
func (re JBool) ToBool() bool      { return bool(re) }
func (re JBool) String() string    { return strings.ToLower(fmt.Sprint(re.ToBool())) }
func (re JBool) Pretty() string    { return re.String() }
func (re JBool) ToObj() JObj       { return JObj{} }
func (re JBool) ToObjPtr() *JObj   { obj := re.ToObj(); return &obj }
func (re JBool) ToArr() JArr       { return JArr{} }
func (re JBool) ToJDoc() JDoc      { return JDoc(re.String()) }
func (re JBool) ToJVal() JValue    { return re }
func (re JBool) ToGVal() any       { return bool(re) }

func NewJBool(val bool) JBool    { jv := JBool(val); return jv }
func ItoJBool(val int) JBool     { jv := JBool(val > 0); return jv }
func FtoJBool(val float32) JBool { jv := JBool(val > 0); return jv }
func DtoJBool(val float64) JBool { jv := JBool(val > 0); return jv }
func AtoJBool(val string) JBool  { return strings.ToLower(val) == "true" || AsLong(val) > 0 }
