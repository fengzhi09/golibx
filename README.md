# golibx

A comprehensive collection of utility libraries for Golang, providing a wide range of helper functions and tools for common programming tasks.

## Packages

### dbx
A database operation library for Golang, providing a unified interface for different databases including PostgreSQL, MySQL, and Doris, with vector database support (Milvus, PGVector, Qdrant) and cache functionality.

[README](dbx/README.md)

### excelx
A powerful Excel file processing library for Golang, supporting XLSX, XLS, and CSV formats with a unified API, stream processing, and file conversion capabilities.

[README](excelx/README.md)

### gox
A comprehensive utility library for Golang, providing a wide range of helper functions for array operations, comparison, conversion, file operations, JSON handling, time utilities, and more.

[README](gox/README.md)

### httpx
A powerful HTTP client library for Golang, providing both low-level and high-level APIs with hooks support, timeout configuration, and file upload/download capabilities.

[README](httpx/README.md)

### jsonx
A JSON processing library for Golang that simplifies JSON parsing and manipulation, providing a user-friendly API for accessing and converting JSON values without type assertions.

[README](jsonx/README.MD)

### logx
A flexible and powerful logging library for Golang, supporting multiple log levels, module-based logging, and various output formats with context support and panic recovery.

[README](logx/README.md)

### op
A library for real-time type calculations, including numerical comparison, option lookup, time comparison, text similarity calculation, and condition merging.

[README](op/README.MD)

### utils
A collection of utility functions and tools for Golang, providing event bus, time utilities, URL query string handling, Viper configuration utilities, and more.

[README](utils/README.md)

## Test Coverage

### Coverage Summary

| Package | Coverage |
| --- | --- |
| dbx | 0.0% |
| dbx/dbx_vec | 0.0% |
| excelx | 20.6% |
| gox | 57.1% |
| httpx | 72.2% |
| jsonx | 8.2% |
| logx | 61.2% |
| op | 88.9% |
| utils | 79.4% |

For detailed test coverage reports, please run the following commands:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## License

MIT

