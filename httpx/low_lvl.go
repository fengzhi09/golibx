package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type (
	HttpOpt    = func(*httpx) *httpx
	WebReqHook = func(client *httpx, req *http.Request) error
	WebRspHook = func(client *httpx, req *http.Request, rsp *http.Response) error
	WebHook    = func(client *httpx, startAt time.Time, req *http.Request, rsp *http.Response, err error)
)

type Httpx interface {
	WithOpts(opts ...HttpOpt) Httpx
	WithReqHooks(hooks ...WebReqHook) Httpx
	WithRspHooks(hooks ...WebRspHook) Httpx
	WithHooks(hooks ...WebHook) Httpx
	Do(method string, path string, body any, params url.Values) (*http.Response, error)
}

type httpx struct {
	client   *http.Client
	reqHooks []WebReqHook
	rspHooks []WebRspHook
	webHooks []WebHook
	baseURL  string
	headers  map[string]string
	timeout  int
}

func NewHttp(timeout int) Httpx {
	return &httpx{
		client:   &http.Client{},
		reqHooks: []WebReqHook{},
		rspHooks: []WebRspHook{},
		webHooks: []WebHook{},
		headers:  make(map[string]string),
		timeout:  timeout,
	}
}

// WithOpts 应用HTTP选项
func (h *httpx) WithOpts(opts ...HttpOpt) Httpx {
	n := &httpx{
		client: h.client, timeout: h.timeout,
		baseURL: h.baseURL, headers: h.headers,
		reqHooks: h.reqHooks, rspHooks: h.rspHooks, webHooks: h.webHooks,
	}

	// 复制headers
	for k, v := range h.headers {
		n.headers[k] = v
	}

	// 应用选项
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		n = opt(n)
	}

	return n
}

// WithReqHooks 添加请求钩子
func (h *httpx) WithReqHooks(hooks ...WebReqHook) Httpx {
	return &httpx{
		client:   h.client,
		reqHooks: append(h.reqHooks, hooks...),
		rspHooks: h.rspHooks,
		webHooks: h.webHooks,
		baseURL:  h.baseURL,
		headers:  h.headers,
		timeout:  h.timeout,
	}
}

// WithRspHooks 添加响应钩子
func (h *httpx) WithRspHooks(hooks ...WebRspHook) Httpx {
	return &httpx{
		client:   h.client,
		reqHooks: h.reqHooks,
		rspHooks: append(h.rspHooks, hooks...),
		webHooks: h.webHooks,
		baseURL:  h.baseURL,
		headers:  h.headers,
		timeout:  h.timeout,
	}
}

// WithHooks 添加通用钩子
func (h *httpx) WithHooks(hooks ...WebHook) Httpx {
	return &httpx{
		client:   h.client,
		reqHooks: h.reqHooks,
		rspHooks: h.rspHooks,
		webHooks: append(h.webHooks, hooks...),
		baseURL:  h.baseURL,
		headers:  h.headers,
		timeout:  h.timeout,
	}
}

// Do 执行HTTP请求
func (h *httpx) Do(method string, path string, body any, params url.Values) (*http.Response, error) {
	// 构建完整URL
	urlStr := path
	if h.baseURL != "" {
		// 确保baseURL和path正确连接
		if h.baseURL[len(h.baseURL)-1] == '/' && len(path) > 0 && path[0] == '/' {
			urlStr = h.baseURL + path[1:]
		} else if h.baseURL[len(h.baseURL)-1] != '/' && len(path) > 0 && path[0] != '/' {
			urlStr = h.baseURL + "/" + path
		} else {
			urlStr = h.baseURL + path
		}
	}

	// 添加查询参数
	if params != nil {
		parsedURL, err := url.Parse(urlStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse URL: %v", err)
		}

		query := parsedURL.Query()
		for k, v := range params {
			for _, val := range v {
				query.Add(k, val)
			}
		}
		parsedURL.RawQuery = query.Encode()
		urlStr = parsedURL.String()
	}
	var resp *http.Response
	var req *http.Request
	var err error
	startAt := time.Now()
	defer func() {
		// 执行钩子
		for _, hook := range h.webHooks {
			hook(h, startAt, req, resp, err)
		}
	}()

	// 处理请求体
	var bodyReader io.Reader
	var bodyCopy *bytes.Buffer
	var jsonData []byte
	if body != nil {
		// 如果body是string或[]byte，直接使用
		switch b := body.(type) {
		case string:
			bodyReader = bytes.NewBufferString(b)
			bodyCopy = bytes.NewBufferString(b)
		case []byte:
			bodyReader = bytes.NewBuffer(b)
			bodyCopy = bytes.NewBuffer(b)
		default:
			// 尝试JSON序列化
			jsonData, err = json.Marshal(body)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal body to JSON: %v", err)
			}
			bodyReader = bytes.NewBuffer(jsonData)
			bodyCopy = bytes.NewBuffer(jsonData)
		}
	}

	// 创建请求
	req, err = http.NewRequest(method, urlStr, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// 设置默认headers
	for k, v := range h.headers {
		req.Header.Set(k, v)
	}

	// 设置默认Content-Type（如果需要）
	if body != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// 对于需要记录日志的情况，创建一个可重复读取的Body
	if bodyCopy != nil {
		req.Body = io.NopCloser(bodyCopy)
	}

	// 执行请求钩子
	for _, hook := range h.reqHooks {
		if err := hook(h, req); err != nil {
			return nil, fmt.Errorf("request hook failed: %v", err)
		}
	}

	// 执行请求
	resp, err = h.client.Do(req.WithContext(context.Background()))
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}

	// 包装响应体以便它可以被多次读取
	if resp.Body != nil {
		// 使用一个缓冲区来存储响应体的内容
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("failed to read response body: %v", err)
		}
		resp.Body.Close()

		// 创建一个新的可读取的Body
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	// 执行响应钩子
	for _, hook := range h.rspHooks {
		if err := hook(h, req, resp); err != nil {
			// 关闭响应体
			resp.Body.Close()
			return nil, fmt.Errorf("response hook failed: %v", err)
		}
	}

	return resp, nil
}

// 常用的HTTP选项函数
func WithBaseURL(baseURL string) HttpOpt {
	return func(h *httpx) *httpx {
		h.baseURL = baseURL
		return h
	}
}

func WithHeader(key, value string) HttpOpt {
	return func(h *httpx) *httpx {
		h.headers[key] = value
		return h
	}
}

func WithHeaders(headers map[string]string) HttpOpt {
	return func(h *httpx) *httpx {
		for k, v := range headers {
			h.headers[k] = v
		}
		return h
	}
}

func WithTimeout(timeout int) HttpOpt {
	return func(h *httpx) *httpx {
		h.timeout = timeout
		return h
	}
}

func WithClient(client *http.Client) HttpOpt {
	return func(h *httpx) *httpx {
		h.client = client
		return h
	}
}

func WithMetric(writer func(method, path string, statusCode int, elapsedMs int64, err error)) WebHook {
	return func(client *httpx, startAt time.Time, req *http.Request, rsp *http.Response, err error) {
		if rsp != nil {
			writer(req.Method, req.URL.Path, rsp.StatusCode, time.Since(startAt).Milliseconds(), err)
		}
	}
}

func WithLog(writer func(method, path string, input string, output string)) WebHook {
	return func(client *httpx, startAt time.Time, req *http.Request, rsp *http.Response, err error) {
		// 读取请求体
		var inputData []byte
		if req.Body != nil {
			bodyBytes, err := io.ReadAll(req.Body)
			if err == nil {
				inputData = bodyBytes
			}
		}

		// 读取响应体
		var outputData []byte
		if rsp.Body != nil {
			bodyBytes, err := io.ReadAll(rsp.Body)
			if err == nil {
				outputData = bodyBytes
			}
		}

		// 调用writer函数记录日志
		writer(req.Method, req.URL.Path, string(inputData), string(outputData))
	}
}
