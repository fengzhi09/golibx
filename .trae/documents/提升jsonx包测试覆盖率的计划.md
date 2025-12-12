# 提升jsonx包测试覆盖率的计划

## 问题分析

通过分析jsonx包的测试覆盖率报告，发现以下问题：

1. **整体覆盖率低**：仅8.2%，大部分函数未被测试
2. **测试范围窄**：现有测试主要集中在时间转换和对象比较上
3. **核心功能缺乏测试**：
   - 数组操作（arr.go）：0%覆盖率
   - 布尔值处理（bool.go）：0%覆盖率
   - 文档解析（doc.go）：大部分函数未测试
   - 整数处理（int.go）：0%覆盖率
   - 空值处理（null.go）：0%覆盖率
   - 数值处理（num.go）：部分函数未测试
   - 对象操作（obj.go）：大部分函数未测试
   - 字符串处理（str.go）：大部分函数未测试
   - 值转换（val.go）：0%覆盖率
   - 扩展库（xlibs.go）：部分函数未测试

## 测试计划

### 1. 添加数组操作测试（arr.go）
- 测试数组基本操作：Size、IsEmpty、Contains、First、Last
- 测试数组遍历：Foreach
- 测试数组访问：GetVal、GetInt、GetTime、GetLong、GetDouble、GetFloat、GetBool、GetStr、GetObj、GetJArr
- 测试数组转换：ToIntArr、ToTimeArr、ToLongArr、ToDoubleArr、ToFloatArr、ToBoolArr、ToStrArr、ToObjArr、ToJArrArr
- 测试数组修改：Add、AddVal、AddStr、AddDT、AddOID、AddInt、AddTS、AddLong、AddFloat、AddDouble、AddObj、AddArr
- 测试数组查找：IndexOfVal、IndexOf
- 测试数组创建：NewJArr、NewJArrPtr

### 2. 添加布尔值处理测试（bool.go）
- 测试布尔值转换：ToInt、ToTime、ToLong、ToDouble、ToFloat、ToBool、ToString
- 测试布尔值创建：NewJBool、ItoJBool、FtoJBool、DtoJBool、AtoJBool

### 3. 完善文档解析测试（doc.go）
- 测试数组解析：ParseJArr
- 测试值转换：JV2GOV、JV2Struct
- 测试JSON输出：DumpJObj、DumpJArr
- 测试转换注册：RegJOConv、RegGoConv
- 测试值转换：GoV2JV

### 4. 添加整数处理测试（int.go）
- 测试整数转换：ToInt、ToTime、ToLong、ToDouble、ToFloat、ToBool、ToString
- 测试整数创建：NewJInt

### 5. 添加空值处理测试（null.go）
- 测试空值转换：ToInt、ToTime、ToLong、ToDouble、ToFloat、ToBool、ToString
- 测试空值创建：NewJNull、AtoJNull

### 6. 完善数值处理测试（num.go）
- 测试数值转换：ToFloat、ToBool、Pretty、ToObj、ToObjPtr、ToArr、ToJDoc、ToJVal、ToGVal

### 7. 完善对象操作测试（obj.go）
- 测试对象基本操作：ToMap、Values、Size、IsEmpty、Contains、Foreach、ForeachObj
- 测试对象合并：Merge、MergePtr
- 测试对象访问：GetValIgnore、GetInt、GetLong、GetDouble、GetFloat、GetBool、GetOrDefault
- 测试对象快照：Snap、Clone
- 测试对象访问：GetObj、GetObjPtr、GetJArr、GetIntArr、GetTimeArr、GetLongArr、GetDoubleArr、GetFloatArr、GetBoolArr、GetStrArr、GetObjArr、GetJArrArr
- 测试对象绑定：BindStruct、ToStruct、ToStructUnsafe
- 测试对象修改：Put、PutInt、PutTS、PutDT、PutLong、PutFloat、PutDouble、PutBool、PutStr、PutObj、PutArr
- 测试对象获取：GetOr
- 测试对象输出：ToLine、ToCSV、ToCells
- 测试对象创建：NewObj、NewObjPtr、ParseObjPtr、ParseObj

### 8. 完善字符串处理测试（str.go）
- 测试字符串转换：ToInt、ToLong、ToDouble、ToFloat、ToBool、Pretty、ToObj、ToObjPtr、ToArr、ToJDoc、ToJVal、ToGVal

### 9. 添加值转换测试（val.go）
- 测试值比较：JValEqJVal、JValEqGVal

### 10. 完善扩展库测试（xlibs.go）
- 测试tryAsBsonId函数的完整逻辑

## 测试策略

1. **编写单元测试**：为每个函数编写独立的测试用例
2. **覆盖边界情况**：测试正常值、边界值、异常值
3. **保持测试独立性**：每个测试用例只测试一个功能点
4. **不修改源码**：仅添加测试用例，不修改任何源码或bug
5. **使用现有测试框架**：继续使用标准库的testing包
6. **保持测试简洁**：每个测试用例尽量简洁明了

## 预期结果

通过添加上述测试用例，预期将jsonx包的测试覆盖率提升到70%以上，覆盖所有核心功能和主要API，同时保持源码不变。

## 实施步骤

1. 创建新的测试文件或扩展现有测试文件
2. 按照上述计划逐步添加测试用例
3. 运行测试验证覆盖率提升情况
4. 优化测试用例，确保覆盖所有关键路径
5. 最终生成新的覆盖率报告

## 注意事项

- 严格遵守不修改源码的原则
- 测试用例应具有可维护性和可读性
- 确保测试用例能够稳定通过
- 测试用例应覆盖所有主要功能和边界情况