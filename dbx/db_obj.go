package dbx

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/fengzhi09/golibx/jsonx"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
)

//goland:noinspection SpellCheckingInspection
type Obj struct {
	ID        OID            `bson:"_id" form:"id" json:"id" gorm:"column:id;type:string"`
	CreatedAt time.Time      `bson:"created_at" form:"created_at" json:"created_at" gorm:"column:created_at;index"`
	UpdatedAt time.Time      `bson:"update_at" form:"update_at" json:"update_at" gorm:"column:update_at;index"`
	DeletedAt gorm.DeletedAt `bson:"deleted_at" form:"deleted_at" json:"-"  gorm:"column:deleted_at;index"`
	ExtAttrs  Json           `bson:"ext_attrs" json:"ext_attrs" gorm:"column:ext_attrs;"` // 扩展属性
}

func (o *Obj) GetAttrs() Json {
	if o.ExtAttrs.JObj == nil {
		o.ExtAttrs.JObj = jsonx.JObj{}
	}
	return o.ExtAttrs
}

// OID sql自定义类型:json对象,数据库存为string，解析为jsonx.JObj;另参见JSONArr
type OID struct {
	primitive.ObjectID
}

func NewOID() OID {
	return OID{primitive.NewObjectID()}
}

func HexOID(oid string) OID {
	id, _ := primitive.ObjectIDFromHex(oid)
	return OID{id}
}

func (j *OID) Scan(value any) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("not json str: %v", value)
	}
	oid, err := primitive.ObjectIDFromHex(str)
	if err == nil {
		j.ObjectID = oid
	}
	return err
}

func (j *OID) IsValid() bool {
	_, err := primitive.ObjectIDFromHex(j.Hex())
	return err == nil
}

func (j *OID) IsEmpty() bool {
	return j.Hex() == primitive.NilObjectID.Hex()
	// return j.ObjectID==primitive.NilObjectID
}

func (j OID) GormDataType() string {
	return CustomerIDOnSql
}

func (j OID) Value() (driver.Value, error) {
	return j.Hex(), nil
}

func (j OID) MarshalJSON() ([]byte, error) {
	return j.ObjectID.MarshalJSON()
}

func (j *OID) UnmarshalJSON(data []byte) error {
	return j.ObjectID.UnmarshalJSON(data)
}

func (j OID) MarshalBSON() ([]byte, error) {
	return []byte(j.ObjectID.String()), nil
}

func (j *OID) UnmarshalBSON(data []byte) error {
	if len(data) == 12 {
		var oid [12]byte
		copy(oid[:], data[:])
		j.ObjectID = oid
		return nil
	}
	return fmt.Errorf("bad oid len:%v", len(data))
}

// Json sql自定义类型:json对象,数据库存为string，解析为jsonx.JObj;另参见JSONArr
type Json struct {
	jsonx.JObj
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 JSONB
func (j *Json) Scan(value any) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("not json str: %v", value)
	}
	obj := jsonx.JObj{}
	if len(str) > 0 {
		obj = jsonx.ParseJObj(str)
	}
	*j = Json{obj}
	return nil
}

func (j Json) GormDataType() string {
	return CustomerTypeOnSql
}

// 实现 driver.Valuer 接口，Data 返回 JSONB value
func (j Json) Value() (driver.Value, error) {
	if j.JObj == nil || j.Size() == 0 {
		return nil, nil
	}
	return j.String(), nil
}

func (j Json) MarshalJSON() ([]byte, error) {
	return []byte(j.JObj.String()), nil
}

func (j *Json) UnmarshalJSON(data []byte) error {
	obj := jsonx.JObj{}
	if len(data) > 0 {
		obj = jsonx.ParseJObj(string(data))
	}
	j.JObj = obj
	return nil
}

func (j Json) MarshalBSON() ([]byte, error) {
	return j.MarshalJSON()
}

func (j *Json) UnmarshalBSON(data []byte) error {
	return j.UnmarshalJSON(data)
}

// JsonArr sql自定义类型:json对象,数据库存为string，解析为jsonx.JArr;另参见JSON
type JsonArr struct {
	jsonx.JArr
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 JSONB
func (j *JsonArr) Scan(value any) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("not json str: %v", value)
	}
	*j = JsonArr{jsonx.ParseJArr(str)}
	return nil
}

func (j JsonArr) GormDataType() string {
	return CustomerTypeOnSql
}

// 实现 driver.Valuer 接口，Data 返回 JSONB value
func (j JsonArr) Value() (driver.Value, error) {
	if j.JArr == nil || j.Size() == 0 {
		return nil, nil
	}
	return j.String(), nil
}

func (j JsonArr) MarshalJSON() ([]byte, error) {
	return []byte(j.JArr.String()), nil
}

func (j *JsonArr) UnmarshalJSON(data []byte) error {
	arr := jsonx.JArr{}
	if len(data) > 0 {
		arr = jsonx.ParseJArr(string(data))
	}
	j.JArr = arr
	return nil
}

func (j JsonArr) MarshalBSON() ([]byte, error) {
	return j.MarshalJSON()
}

func (j *JsonArr) UnmarshalBSON(data []byte) error {
	return j.UnmarshalJSON(data)
}

type StrArr []string

func NewStrArr() StrArr {
	return StrArr{}
}

func (sa *StrArr) Add(str string) {
	*sa = append(*sa, str)
}

func (sa *StrArr) ToJArr() jsonx.JArr {
	arr := make([]any, 0)
	for _, v := range *sa {
		arr = append(arr, v)
	}
	return jsonx.NewJArr(arr...)
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 JSONB
func (sa *StrArr) Scan(value any) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("not json str: %v", value)
	}
	return sa.UnmarshalJSON([]byte(str))
}

func (sa StrArr) GormDataType() string {
	return CustomerTypeOnSql
}

// 实现 driver.Valuer 接口，Data 返回 JSONB value
func (sa StrArr) Value() (driver.Value, error) {
	if len(sa) == 0 {
		return nil, nil
	}
	return jsonx.UnsafeMarshalString(sa), nil
}

func (sa StrArr) MarshalJSON() ([]byte, error) {
	if len(sa) == 0 {
		return []byte("[]"), nil
	}
	str := sa.ToJArr().String()
	return []byte(str), nil
}

func (sa *StrArr) UnmarshalJSON(data []byte) error {
	arr := jsonx.JArr{}
	if len(data) > 0 {
		arr = jsonx.ParseJArr(string(data))
	}
	*sa = arr.ToStrArr()
	return nil
}

func (sa StrArr) MarshalBSON() ([]byte, error) {
	return sa.MarshalJSON()
}

func (sa *StrArr) UnmarshalBSON(data []byte) error {
	return sa.UnmarshalJSON(data)
}

type (
	IntArr  = LongArr
	LongArr []int64
)

func NewIntArr() LongArr {
	return LongArr{}
}

func (ia *LongArr) Add(elem int64) {
	*ia = append(*ia, elem)
}

func (ia *LongArr) ToJArr() jsonx.JArr {
	arr := make([]any, 0)
	for _, v := range *ia {
		arr = append(arr, v)
	}
	return jsonx.NewJArr(arr...)
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 JSONB
func (ia *LongArr) Scan(value any) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("not json str: %v", value)
	}
	return ia.UnmarshalJSON([]byte(str))
}

func (ia LongArr) GormDataType() string {
	return CustomerTypeOnSql
}

// 实现 driver.Valuer 接口，Data 返回 JSONB value
func (ia LongArr) Value() (driver.Value, error) {
	if len(ia) == 0 {
		return nil, nil
	}
	return jsonx.UnsafeMarshalString(ia), nil
}

func (ia LongArr) MarshalJSON() ([]byte, error) {
	if len(ia) == 0 {
		return []byte("[]"), nil
	}
	str := ia.ToJArr().String()
	return []byte(str), nil
}

func (ia *LongArr) UnmarshalJSON(data []byte) error {
	arr := jsonx.JArr{}
	if len(data) > 0 {
		arr = jsonx.ParseJArr(string(data))
	}
	*ia = arr.ToLongArr()
	return nil
}

func (ia LongArr) MarshalBSON() ([]byte, error) {
	return ia.MarshalJSON()
}

func (ia *LongArr) UnmarshalBSON(data []byte) error {
	return ia.UnmarshalJSON(data)
}

type (
	FloatArr  = DoubleArr
	DoubleArr []float64
)

func NewDoubleArr() DoubleArr {
	return DoubleArr{}
}

func (re *DoubleArr) Add(elem float64) {
	*re = append(*re, elem)
}

func (re *DoubleArr) ToJArr() jsonx.JArr {
	arr := make([]any, 0)
	for _, v := range *re {
		arr = append(arr, v)
	}
	return jsonx.NewJArr(arr...)
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 JSONB
func (re *DoubleArr) Scan(value any) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("not json str: %v", value)
	}
	return re.UnmarshalJSON([]byte(str))
}

func (re DoubleArr) GormDataType() string {
	return CustomerTypeOnSql
}

// 实现 driver.Valuer 接口，Data 返回 JSONB value
func (re DoubleArr) Value() (driver.Value, error) {
	if len(re) == 0 {
		return nil, nil
	}
	return jsonx.UnsafeMarshalString(re), nil
}

func (re DoubleArr) MarshalJSON() ([]byte, error) {
	if len(re) == 0 {
		return []byte("[]"), nil
	}
	str := re.ToJArr().String()
	return []byte(str), nil
}

func (re *DoubleArr) UnmarshalJSON(data []byte) error {
	arr := jsonx.JArr{}
	if len(data) > 0 {
		arr = jsonx.ParseJArr(string(data))
	}
	*re = arr.ToDoubleArr()
	return nil
}

func (re DoubleArr) MarshalBSON() ([]byte, error) {
	return re.MarshalJSON()
}

func (re *DoubleArr) UnmarshalBSON(data []byte) error {
	return re.UnmarshalJSON(data)
}

type MapI2S map[string]string

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 JSONB
func (j *MapI2S) Scan(value any) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("not json str: %v", value)
	}
	err := j.UnmarshalJSON([]byte(str))
	return err
}

func (j MapI2S) GormDataType() string {
	return CustomerTypeOnSql
}

// 实现 driver.Valuer 接口，Data 返回 JSONB value
func (j MapI2S) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return j.MarshalJSON()
}

func (j MapI2S) MarshalJSON() ([]byte, error) {
	if len(j) == 0 {
		return []byte("{}"), nil
	}
	str := "{"
	for i, v := range j {
		if len(str) > 1 {
			str += ","
		}
		str += fmt.Sprintf("\"%v\":\"%v\"", i, v)
	}
	str += "}"
	return []byte(str), nil
}

func (j *MapI2S) UnmarshalJSON(data []byte) error {
	obj := jsonx.JObj{}
	if len(data) > 0 {
		obj = jsonx.ParseJObj(string(data))
	}
	tmp := map[string]string{}
	obj.Foreach(func(k string, v jsonx.JValue) bool {
		tmp[k] = v.String()
		return true
	})
	*j = tmp
	return nil
}

func (j MapI2S) MarshalBSON() ([]byte, error) {
	return j.MarshalJSON()
}

func (j *MapI2S) UnmarshalBSON(data []byte) error {
	return j.UnmarshalJSON(data)
}

const (
	CustomerIDOnSql   = "varchar(24)"
	CustomerTypeOnSql = "text"
)
