## jsonx
### 1. Why do we need jsonx
- If you are tired of using obj.(int/string/base_type) everywhere
- If you want to write deep parsing objects without worrying about whether any intermediate level is null
- If you want to write a backward-compatible data model or business architecture
- You want to get the experience of operating json objects as if you were using native types

If you meet one of the above conditions, I recommend you use this package

### 2. Usage effect

#### 2.1 Basic usage
```golang
jsonStr := "{\"arr\":[{\"f1\":\"a string\",\"f2\":\"1.2\",\"f3\":\"2002-01-01\",\"f4\":\"20020101\"}]}"
root := jsonx.ParseObj([]byte(jsonStr))
arr := root.GetJArr("arr")
f1 := a.GetObj(0).GetStr("f1")
```
#### 2.2 Type conversion
Suppose there is a json object that needs to be parsed, and you want to read fields according to specified types, you may have the following code
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
Now it only takes 4 lines
```golang
// convert code free，just get directly to the type you just want
root := jsonx.ParseObj([]byte(jsonStr))
arr := root.GetJArr("arr")
f2 := arr.GetObj(0).GetStr("f2").ToDouble()
f3 := arr.GetObj(0).GetStr("f3").ToTime()
f4 := arr.GetObj(0).GetStr("f4").ToTime()
// or simply
root := jsonx.ParseObj([]byte(jsonStr))
f3_ := root.GetJArr("arr").GetObj(0).GetTime("f3")
f2_ := root.GetJArr("arr").GetObj(0).GetDouble("f2")
f4_ := root.GetJArr("arr").GetObj(0).GetTime("f4")
```
#### 2.3 Field writing
```golang
// use api for type-spec
a :=jsonx.JObj{}
a.PutInt("int",1)
a.PutLong("long",1e12)
a.PutString("string","a string")
a.PutObj("obj",jsonx.JObj{})
a.PutArr("arr",jsonx.JArr{})
// or simply
a.Put("arr",[]int{1,2,3})
```
### 3. Technical details
- Native collection types, no additional memory usage
```golang
type JInt int64
type JNum float64
type JObj map[string]any
type JStr string
type JBool bool
type JArr []any
type JNull struct {}
```
- Easy-to-use conversion API
```golang
type JValue interface {
    Type() JType         // Data type
    IsNull() bool        // Whether null
    ToInt() int          // Convert to int
    ToLong() int64       // Convert to int64
    ToTime() time.Time   // Convert to time.Time
    ToDouble() float64   // Convert to float64
    ToFloat() float32    // Convert to float32
    ToBool() bool        // Convert to bool
    String() string      // Convert to string, equivalent to fmt.Sprint(v)
    ToObj() JObj         // Convert to object, JStr will be parsed as json string
    ToObjPtr() *JObj     // Convert to object pointer
    ToArr() JArr         // Convert to array, JStr will be parsed as json string
    ToJDoc() JDoc        // Output as json string format, JStr will be escaped with ""
    ToJVal() JValue      // Return interface operation object
    ToGVal() any // Return underlying object
}
```
- Explore more by yourself


### 4. UT&benchmark
 //to be move to this package


### 5. Where used
- 明察质检xiaoduo-qc/qc-engine
