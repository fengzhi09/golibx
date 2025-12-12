## jsonx
### 1. 为什么需要jsonx
- 如果你厌烦了到处使用obj.(int/string/base_type)之类的代码
- 如果你想写出深层解析对象，又无需关心是否中间某层级为null
- 如果你想写一个向后兼容的数据建模或者业务架构
- 你希望获得如同使用原生类型一样去操作json对象的体验

你满足以上情况之一我建议你使用本包

### 2.使用效果

#### 2.1基本使用
```golang
jsonStr := "{\"arr\":[{\"f1\":\"a string\",\"f2\":\"1.2\",\"f3\":\"2002-01-01\",\"f4\":\"20020101\"}]}"
root := jsonx.ParseObj([]byte(jsonStr))
arr := root.GetJArr("arr")
f1 := a.GetObj(0).GetStr("f1")
```
#### 2.2类型转换
假设有个json对象需要解析，要按指定的类型读取字段，可能会有如下代码
```
a := map[string]any{}
_ = json.Unmarshal([]byte(jsonStr),&a)
var f2 float64 = 0.0
var f3 time.Time = time.Unix(0,0)
var f4 time.Time = time.Unix(0,0)
if val,hit := a["f2"]; hit{
    switch val.(type){
    case float64:
        f2=val.(float64)
        break
    }
}
if val,hit := a["f3"]; hit{
    switch val.(type){
    case string:
        f3=time.Parse("2006-01-02",val.(string))
        break
    }
}
if val,hit := a["f4"]; hit{
    switch val.(type){
    case string:
        f4=time.Parse("2006-01-02",val.(string))
        break
    case int:
        v:=val.(int)
        year,month,date:=v/10000,v/100%100,v%100
        f4=time.Date(year, time.Month(month), date, 0, 0, 0, 0, time.Local)
        break
    }
}
```
现在仅需要4行即可
```golang
// 无需转换代码，直接获取你想要的类型
root := jsonx.ParseObj([]byte(jsonStr))
arr := root.GetJArr("arr")
f2 := arr.GetObj(0).GetStr("f2").ToDouble()
f3 := arr.GetObj(0).GetStr("f3").ToTime()
f4 := arr.GetObj(0).GetStr("f4").ToTime()
// 或者更简单
root := jsonx.ParseObj([]byte(jsonStr))
f3_ := root.GetJArr("arr").GetObj(0).GetTime("f3")
f2_ := root.GetJArr("arr").GetObj(0).GetDouble("f2")
f4_ := root.GetJArr("arr").GetObj(0).GetTime("f4")
```
#### 2.3字段写入
```golang
// 使用类型指定的API
a := jsonx.JObj{}
a.PutInt("int",1)
a.PutLong("long",1e12)
a.PutString("string","a string")
a.PutObj("obj",jsonx.JObj{})
a.PutArr("arr",jsonx.JArr{})
// 或者更简单
a.Put("arr",[]int{1,2,3})
```
### 3.技术细节
- 原生集合类型，不会占用更多内存
```golang
type JInt int64
type JNum float64
type JObj map[string]any
type JStr string
type JBool bool
type JArr []any
type JNull struct {}
```
- 好用的转换API
```golang
type JValue interface {
    Type() JType         // 数据类型
    IsNull() bool        // 是否空
    ToInt() int          // 强转为int
    ToLong() int64       // 强转为int64
    ToTime() time.Time   // 强转为时间 time.Time
    ToDouble() float64   // 强转为float64
    ToFloat() float32    // 强转为float32
    ToBool() bool        // 强转为bool
    String() string      // 强转为string, 等效fmt.Sprint(v)
    ToObj() JObj         // 强转为对象，JStr 会当做json字符串解析
    ToObjPtr() *JObj     // 强转为对象指针
    ToArr() JArr         // 强转为数组，JStr 会当做json字符串解析
    ToJDoc() JDoc        // 安装json字符串格式输出，JStr 会转义加""
    ToJVal() JValue      // 返回接口操作对象
    ToGVal() any // 返回底层对象
}
```
- 更多请自行探索

