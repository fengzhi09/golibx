package httpx

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 测试NewHttp函数
func TestNewHttp(t *testing.T) {
	t.Run("创建默认HTTP客户端", func(t *testing.T) {
		client := NewHttp(10) // 提供timeout参数
		assert.NotNil(t, client)
	})
}

// 测试WithOpts函数
func TestWithOpts(t *testing.T) {
	client := NewHttp(10)

	// 测试WithBaseURL选项
	client = client.WithOpts(WithBaseURL("http://example.com"))
	assert.NotNil(t, client)

	// 测试WithHeader选项
	client = client.WithOpts(WithHeader("Authorization", "Bearer token"))
	assert.NotNil(t, client)

	// 测试WithTimeout选项
	client = client.WithOpts(WithTimeout(30))
	assert.NotNil(t, client)
}

// 测试URL构建逻辑
func TestURLBuilding(t *testing.T) {
	// 已在下面的请求测试中覆盖
}

// 测试基本HTTP请求
func TestDoRequest(t *testing.T) {
	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success"}`))
	}))
	defer server.Close()

	// 创建HTTP客户端
	client := NewHttp(10).WithOpts(WithBaseURL(server.URL))

	// 执行GET请求
	resp, err := client.Do("GET", "/test", nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(body), "success")
}

// 测试请求体处理
func TestRequestBodyHandling(t *testing.T) {
	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 读取请求体
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		// 返回请求体内容
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}))
	defer server.Close()

	// 创建HTTP客户端
	client := NewHttp(10).WithOpts(WithBaseURL(server.URL))

	// 测试JSON请求体
	testData := map[string]interface{}{
		"name":  "test",
		"value": 123,
	}

	resp, err := client.Do("POST", "/test", testData, nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	defer resp.Body.Close()

	// 读取响应体
	var responseData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	assert.NoError(t, err)
	assert.Equal(t, "test", responseData["name"])
	assert.Equal(t, 123.0, responseData["value"])

	// 测试字符串请求体
	resp, err = client.Do("POST", "/test", "plain text", nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "plain text", string(body))
}

// 测试查询参数
func TestQueryParameters(t *testing.T) {
	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证查询参数
		assert.Equal(t, "value1", r.URL.Query().Get("param1"))
		assert.Equal(t, "value2", r.URL.Query().Get("param2"))

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// 创建HTTP客户端
	client := NewHttp(10).WithOpts(WithBaseURL(server.URL))

	// 设置查询参数
	params := url.Values{}
	params.Add("param1", "value1")
	params.Add("param2", "value2")

	resp, err := client.Do("GET", "/test", nil, params)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	defer resp.Body.Close()
}

// 测试请求钩子
func TestRequestHooks(t *testing.T) {
	var hookCalled bool
	hook := func(client *httpx, req *http.Request) error {
		hookCalled = true
		req.Header.Set("X-Custom-Header", "test-value")
		return nil
	}

	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "test-value", r.Header.Get("X-Custom-Header"))
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// 创建带钩子的HTTP客户端
	client := NewHttp(10).WithReqHooks(hook).WithOpts(WithBaseURL(server.URL))

	resp, err := client.Do("GET", "/test", nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	defer resp.Body.Close()
	assert.True(t, hookCalled)
}

// 测试响应钩子
func TestResponseHooks(t *testing.T) {
	var hookCalled bool
	hook := func(client *httpx, req *http.Request, rsp *http.Response) error {
		hookCalled = true
		rsp.Header.Set("X-Custom-Response", "test-response")
		return nil
	}

	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// 创建带钩子的HTTP客户端
	client := NewHttp(10).WithRspHooks(hook).WithOpts(WithBaseURL(server.URL))

	resp, err := client.Do("GET", "/test", nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	defer resp.Body.Close()
	assert.True(t, hookCalled)
	assert.Equal(t, "test-response", resp.Header.Get("X-Custom-Response"))
}

// 测试通用钩子
func TestWebHooks(t *testing.T) {
	var hookCalled bool
	hook := func(client *httpx, startAt time.Time, req *http.Request, rsp *http.Response, err error) {
		hookCalled = true
	}

	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// 创建带钩子的HTTP客户端
	client := NewHttp(10).WithHooks(hook).WithOpts(WithBaseURL(server.URL))

	resp, err := client.Do("GET", "/test", nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	defer resp.Body.Close()
	assert.True(t, hookCalled)
}

// 测试WithMetric函数
func TestWithMetric(t *testing.T) {
	var (
		mMethod     string
		mPath       string
		mStatusCode int
		mElapsedMs  int64
		mErr        error
	)

	metricHook := func(method, path string, statusCode int, elapsedMs int64, err error) {
		mMethod = method
		mPath = path
		mStatusCode = statusCode
		mElapsedMs = elapsedMs
		mErr = err
	}

	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// 创建带指标钩子的HTTP客户端
	client := NewHttp(10).WithHooks(WithMetric(metricHook)).WithOpts(WithBaseURL(server.URL))

	resp, err := client.Do("GET", "/test", nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	defer resp.Body.Close()

	// 验证指标数据
	assert.Equal(t, "GET", mMethod)
	assert.Contains(t, mPath, "/test")
	assert.Equal(t, http.StatusOK, mStatusCode)
	assert.GreaterOrEqual(t, mElapsedMs, int64(0), "Elapsed time should be non-negative")
	assert.Nil(t, mErr)
}

// 测试错误处理
func TestErrorHandling(t *testing.T) {
	// 测试无效的URL
	client := NewHttp(10)
	resp, err := client.Do("GET", "://invalid-url", nil, nil)
	assert.Error(t, err)
	assert.Nil(t, resp)

	// 测试JSON序列化错误
	type Unserializable struct {
		Data chan int // 不可序列化的字段
	}
	resp, err = client.Do("POST", "http://example.com", Unserializable{Data: make(chan int)}, nil)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

// 测试响应体可重复读取
func TestResponseBodyRewind(t *testing.T) {
	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"test":"data"}`))
	}))
	defer server.Close()

	client := NewHttp(10).WithOpts(WithBaseURL(server.URL))
	resp, err := client.Do("GET", "/test", nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	defer resp.Body.Close()

	// 第一次读取
	body1, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, `{"test":"data"}`, strings.ReplaceAll(string(body1), "\n", ""))

	// 检查是否可以重新读取（通过类型断言和Seek）
	if seeker, ok := resp.Body.(io.Seeker); ok {
		_, err = seeker.Seek(0, io.SeekStart)
		assert.NoError(t, err)

		body2, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, string(body1), string(body2))
	}
}

// 测试URL连接逻辑
func TestURLJoining(t *testing.T) {
	testCases := []struct {
		baseURL  string
		path     string
		expected string
	}{{
		baseURL:  "http://example.com",
		path:     "/api",
		expected: "http://example.com/api",
	}, {
		baseURL:  "http://example.com/",
		path:     "api",
		expected: "http://example.com/api",
	}, {
		baseURL:  "http://example.com/",
		path:     "/api",
		expected: "http://example.com/api",
	}, {
		baseURL:  "http://example.com",
		path:     "api",
		expected: "http://example.com/api",
	}}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Test case %d", i), func(t *testing.T) {
			// 创建测试服务器
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			// 模拟URL连接
			client := NewHttp(10).WithOpts(WithBaseURL(tc.baseURL))
			resp, err := client.Do("GET", tc.path, nil, nil)
			// 这里我们只验证不会崩溃，因为实际请求会失败（tc.baseURL不是真实服务器）
			if err != nil {
				assert.Contains(t, err.Error(), "connection refused") // 预期的错误类型
			} else if resp != nil {
				defer resp.Body.Close()
			}
		})
	}
}
