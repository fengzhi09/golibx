# golibx

一个全面的Golang实用工具库集合，为常见编程任务提供广泛的辅助函数和工具。

## 包

### dbx
一个用于Golang的数据库操作库，为不同的数据库（包括PostgreSQL、MySQL和Doris）提供统一的接口，支持向量数据库（Milvus、PGVector、Qdrant）和缓存功能。

[README](dbx/README.cn.md)

### excelx
一个功能强大的Golang Excel文件处理库，支持XLSX、XLS和CSV格式，具有统一的API、流处理和文件转换功能。

[README](excelx/README.cn.md)

### gox
一个全面的Golang实用工具库，提供了数组操作、比较、转换、文件操作、JSON处理、时间工具等多种辅助函数。

[README](gox/README.cn.md)

### httpx
一个功能强大的Golang HTTP客户端库，提供低级和高级API，支持钩子机制、超时配置和文件上传/下载功能。

[README](httpx/README.cn.md)

### jsonx
一个用于Golang的JSON处理库，简化了JSON解析和操作，提供了用户友好的API，无需类型断言即可访问和转换JSON值。

[README](jsonx/README.cn.md)

### logx
一个灵活且强大的Golang日志库，支持多种日志级别、基于模块的日志记录和各种输出格式，具有上下文支持和panic恢复功能。

[README](logx/README.cn.md)

### op
一个用于实时类型计算的库，包括数值比较、选项查找、时间比较、文本相似度计算和条件合并。

[README](op/README.cn.md)

### utils
一个Golang实用工具函数和工具的集合，提供了事件总线、时间工具、URL查询字符串处理、Viper配置工具等功能。

[README](utils/README.cn.md)

## 测试覆盖率

### 覆盖率摘要

| 包 | 覆盖率 |
| --- | --- |
| dbx | 0.0% |
| dbx/dbx_vec | 0.0% |
| excelx | 20.6% |
| gox | 57.1% |
| httpx | 72.2% |
| jsonx | 53.7% |
| logx | 61.2% |
| op | 88.9% |
| utils | 79.4% |

要查看详细的测试覆盖率报告，请运行以下命令：

```bash
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## 许可证

MIT
