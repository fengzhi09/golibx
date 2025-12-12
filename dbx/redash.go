package dbx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/fengzhi09/golibx/gox"
	"github.com/fengzhi09/golibx/logx"
)

// Redash Redash API客户端
type Redash struct {
	redashURL string
	apiKey    string
	client    *http.Client
}

// NewRedash 创建Redash客户端实例
func NewRedash(redashURL, apiKey string) *Redash {
	return &Redash{
		redashURL: strings.TrimSuffix(redashURL, "/"),
		apiKey:    apiKey,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Events 获取事件列表
func (r *Redash) Events(page, pageSize int) (map[string]any, error) {
	params := url.Values{}
	params.Add("page", fmt.Sprintf("%d", page))
	params.Add("page_size", fmt.Sprintf("%d", pageSize))
	params.Add("api_key", r.apiKey)

	return r.getJSON("api/events", params)
}

// Queries 获取查询列表
func (r *Redash) Queries(page, pageSize int, onlyFavorites bool) (map[string]any, error) {
	targetURL := "api/queries"
	if onlyFavorites {
		targetURL = "api/queries/favorites"
	}

	params := url.Values{}
	params.Add("page", fmt.Sprintf("%d", page))
	params.Add("page_size", fmt.Sprintf("%d", pageSize))
	params.Add("api_key", r.apiKey)

	return r.getJSON(targetURL, params)
}

// CreateFavorite 创建收藏
func (r *Redash) CreateFavorite(_type string, id int) (map[string]any, error) {
	var urlPath string
	switch _type {
	case "dashboard":
		urlPath = fmt.Sprintf("api/dashboards/%d/favorite", id)
	case "query":
		urlPath = fmt.Sprintf("api/queries/%d/favorite", id)
	default:
		return nil, fmt.Errorf("unsupported type: %s", _type)
	}

	return r.postJSON(urlPath, nil)
}

// GetQuery 获取查询详情
func (r *Redash) GetQuery(queryID int) (map[string]any, error) {
	return r.getJSON(fmt.Sprintf("api/queries/%d", queryID), nil)
}

// Users 获取用户列表
func (r *Redash) Users(page, pageSize int, onlyDisabled bool) (map[string]any, error) {
	params := url.Values{}
	params.Add("page", fmt.Sprintf("%d", page))
	params.Add("page_size", fmt.Sprintf("%d", pageSize))
	params.Add("disabled", fmt.Sprintf("%v", onlyDisabled))
	params.Add("api_key", r.apiKey)

	return r.getJSON("api/users", params)
}

// DisableUser 禁用用户
func (r *Redash) DisableUser(userID int) (map[string]any, error) {
	params := url.Values{}
	params.Add("api_key", r.apiKey)

	return r.postJSON(fmt.Sprintf("api/users/%d/disable", userID), nil, params)
}

// Dashboards 获取仪表板列表
func (r *Redash) Dashboards(page, pageSize int, onlyFavorites bool) (map[string]any, error) {
	targetURL := "api/dashboards"
	if onlyFavorites {
		targetURL = "api/dashboards/favorites"
	}

	params := url.Values{}
	params.Add("page", fmt.Sprintf("%d", page))
	params.Add("page_size", fmt.Sprintf("%d", pageSize))
	params.Add("api_key", r.apiKey)

	return r.getJSON(targetURL, params)
}

// GetDashboard 获取仪表板详情
func (r *Redash) GetDashboard(id any) (map[string]any, error) {
	return r.getJSON(fmt.Sprintf("api/dashboards/%v", id), nil)
}

// GetDataSources 获取数据源列表
func (r *Redash) GetDataSources() (map[string]any, error) {
	return r.getJSON("api/data_sources", nil)
}

// GetDataSource 获取数据源详情
func (r *Redash) GetDataSource(id int) (map[string]any, error) {
	return r.getJSON(fmt.Sprintf("api/data_sources/%d", id), nil)
}

// CreateDataSource 创建数据源
func (r *Redash) CreateDataSource(name, _type string, options map[string]any) (map[string]any, error) {
	payload := map[string]any{
		"name":    name,
		"type":    _type,
		"options": options,
	}

	return r.postJSON("api/data_sources", payload)
}

// CreateQuery 创建查询
func (r *Redash) CreateQuery(ds string, sql string, params map[string]any) (map[string]any, error) {
	if params == nil {
		params = make(map[string]any)
	}

	queryJSON := map[string]any{
		"query":      sql,
		"parameters": params,
	}

	return r.postJSON("api/queries", queryJSON)
}

// CreateDashboard 创建仪表板
func (r *Redash) CreateDashboard(name string) (map[string]any, error) {
	payload := map[string]any{
		"name": name,
	}

	return r.postJSON("api/dashboards", payload)
}

// UpdateDashboard 更新仪表板
func (r *Redash) UpdateDashboard(dashboardID int, properties map[string]any) (map[string]any, error) {
	return r.postJSON(fmt.Sprintf("api/dashboards/%d", dashboardID), properties)
}

// CreateWidget 创建小部件
func (r *Redash) CreateWidget(dashboardID, visualizationID int, text string, options map[string]any) (map[string]any, error) {
	data := map[string]any{
		"dashboard_id":     dashboardID,
		"visualization_id": visualizationID,
		"text":             text,
		"options":          options,
		"width":            1,
	}

	return r.postJSON("api/widgets", data)
}

// DuplicateDashboard 复制仪表板
func (r *Redash) DuplicateDashboard(slug string, newName string) (map[string]any, error) {
	currentDashboard, err := r.GetDashboard(slug)
	if err != nil {
		return nil, err
	}

	if newName == "" {
		if name, ok := currentDashboard["name"].(string); ok {
			newName = "Copy of: " + name
		} else {
			newName = "Copy of Dashboard"
		}
	}

	newDashboard, err := r.CreateDashboard(newName)
	if err != nil {
		return nil, err
	}

	if tags, ok := currentDashboard["tags"].([]any); ok && len(tags) > 0 {
		dashboardID, ok := newDashboard["id"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid dashboard ID")
		}
		r.UpdateDashboard(int(dashboardID), map[string]any{"tags": tags})
	}

	if widgets, ok := currentDashboard["widgets"].([]any); ok {
		dashboardID, ok := newDashboard["id"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid dashboard ID")
		}

		for _, widget := range widgets {
			if w, ok := widget.(map[string]any); ok {
				var visualizationID int
				if viz, ok := w["visualization"].(map[string]any); ok {
					if id, ok := viz["id"].(float64); ok {
						visualizationID = int(id)
					}
				}

				text := ""
				if t, ok := w["text"].(string); ok {
					text = t
				}

				options := make(map[string]any)
				if opt, ok := w["options"].(map[string]any); ok {
					options = opt
				}

				r.CreateWidget(int(dashboardID), visualizationID, text, options)
			}
		}
	}

	return newDashboard, nil
}

// DuplicateQuery 复制查询
func (r *Redash) DuplicateQuery(queryID int, newName string) (map[string]any, error) {
	response, err := r.postJSON(fmt.Sprintf("api/queries/%d/fork", queryID), nil)
	if err != nil {
		return nil, err
	}

	if newName == "" {
		return response, nil
	}

	newQuery := make(map[string]any)
	for k, v := range response {
		newQuery[k] = v
	}
	newQuery["name"] = newName

	if id, ok := newQuery["id"].(float64); ok {
		return r.UpdateQuery(int(id), newQuery)
	}

	return nil, fmt.Errorf("invalid query ID")
}

// ScheduledQueries 获取计划查询列表
func (r *Redash) ScheduledQueries() ([]any, error) {
	allQueries, err := r.Paginate(func(page, pageSize int) (map[string]any, error) {
		return r.Queries(page, pageSize, false)
	})
	if err != nil {
		return nil, err
	}

	scheduled := make([]any, 0)
	for _, query := range allQueries {
		if q, ok := query.(map[string]any); ok {
			if q["schedule"] != nil {
				scheduled = append(scheduled, query)
			}
		}
	}

	return scheduled, nil
}

// UpdateQuery 更新查询
func (r *Redash) UpdateQuery(queryID int, data map[string]any) (map[string]any, error) {
	return r.postJSON(fmt.Sprintf("api/queries/%d", queryID), data)
}

// UpdateVisualization 更新可视化
func (r *Redash) UpdateVisualization(vizID int, data map[string]any) (map[string]any, error) {
	return r.postJSON(fmt.Sprintf("api/visualizations/%d", vizID), data)
}

// Alerts 获取告警列表
func (r *Redash) Alerts() (map[string]any, error) {
	return r.getJSON("api/alerts", nil)
}

// GetAlert 获取告警详情
func (r *Redash) GetAlert(alertID int) (map[string]any, error) {
	return r.getJSON(fmt.Sprintf("api/alerts/%d", alertID), nil)
}

// CreateAlert 创建告警
func (r *Redash) CreateAlert(name string, options map[string]any, queryID int) (map[string]any, error) {
	payload := map[string]any{
		"name":     name,
		"options":  options,
		"query_id": queryID,
	}

	return r.postJSON("api/alerts", payload)
}

// UpdateAlert 更新告警
func (r *Redash) UpdateAlert(id int, name *string, options map[string]any, queryID *int, rearm *int) (map[string]any, error) {
	payload := make(map[string]any)

	if name != nil {
		payload["name"] = *name
	}
	if options != nil {
		payload["options"] = options
	}
	if queryID != nil {
		payload["query_id"] = *queryID
	}
	if rearm != nil {
		payload["rearm"] = *rearm
	}

	return r.postJSON(fmt.Sprintf("api/alerts/%d", id), payload)
}

// Paginate 分页获取所有资源
func (r *Redash) Paginate(resource func(int, int) (map[string]any, error)) ([]any, error) {
	page := 1
	pageSize := 100
	allItems := make([]any, 0)

	for {
		response, err := resource(page, pageSize)
		if err != nil {
			return nil, err
		}

		if results, ok := response["results"].([]any); ok {
			allItems = append(allItems, results...)
		}

		pageVal, pageOk := response["page"].(float64)
		pageSizeVal, pageSizeOk := response["page_size"].(float64)
		count, countOk := response["count"].(float64)

		if !pageOk || !pageSizeOk || !countOk {
			break
		}

		if int(pageVal)*int(pageSizeVal) >= int(count) {
			break
		}

		page++
	}

	return allItems, nil
}

// get 发送GET请求
func (r *Redash) get(path string, params url.Values) (*http.Response, error) {
	return r.request(http.MethodGet, path, nil, params)
}

// post 发送POST请求
func (r *Redash) post(path string, body any, params url.Values) (*http.Response, error) {
	return r.request(http.MethodPost, path, body, params)
}

// delete 发送DELETE请求
func (r *Redash) delete(path string, params url.Values) (*http.Response, error) {
	return r.request(http.MethodDelete, path, nil, params)
}

// request 发送HTTP请求
func (r *Redash) request(method string, path string, body any, params url.Values) (*http.Response, error) {
	requestURL := fmt.Sprintf("%s/%s", r.redashURL, path)

	// 构建请求体
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %v", err)
		}
		bodyReader = bytes.NewBuffer(data)
	}

	// 创建请求
	req, err := http.NewRequest(method, requestURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// 设置请求头
	req.Header.Set("Authorization", fmt.Sprintf("Key %s", r.apiKey))
	req.Header.Set("Content-Type", "application/json")

	// 添加查询参数
	if params != nil {
		req.URL.RawQuery = params.Encode()
	}

	// 发送请求
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}

	// 检查响应状态码
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		bodyBytes, _ := io.ReadAll(resp.Body)
		return resp, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return resp, nil
}

// getJSON 发送GET请求并解析JSON响应
func (r *Redash) getJSON(path string, params url.Values) (map[string]any, error) {
	resp, err := r.get(path, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return r.parseJSONResponse(resp)
}

// postJSON 发送POST请求并解析JSON响应
func (r *Redash) postJSON(path string, body any, params ...url.Values) (map[string]any, error) {
	var p url.Values
	if len(params) > 0 {
		p = params[0]
	}

	resp, err := r.post(path, body, p)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return r.parseJSONResponse(resp)
}

// parseJSONResponse 解析JSON响应
func (r *Redash) parseJSONResponse(resp *http.Response) (map[string]any, error) {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v, body: %s", err, string(bodyBytes))
	}

	return result, nil
}

var spaceRegex = gox.MustCompile(`\s+`, gox.ReV2)

// SearchQuery 搜索Redash中的查询语句
func SearchQuery(ctx context.Context, url string, appKey string, filter map[string]any, want int, pnum int, psize int, limit int) ([]map[string]any, error) {
	redash := NewRedash(url, appKey)
	got := 0
	scanned := 0
	records := make([]map[string]any, 0)

	for got < want && scanned < limit {
		resp, err := redash.Events(pnum, psize)
		if err != nil {
			logx.ErrorfM(ctx, "Redash", "获取事件失败: %v", err)
			break
		}

		if results, ok := resp["results"].([]any); ok {
			logx.DebugfM(ctx, "Redash", "已扫描 %d 个事件，获取 %d 个项目，页码: %d, 每页大小: %d", scanned, got, pnum, psize)
			scanned += len(results)

			for _, item := range results {
				event, ok := item.(map[string]any)
				if !ok {
					continue
				}

				// 提取SQL和数据源
				if details, ok := event["details"].(map[string]any); ok {
					if query, ok := details["query"].(string); ok {
						event["sql"] = spaceRegex.ReplaceAll(strings.TrimSpace(query), " ")
					}
					if ds, ok := details["data_source"].(string); ok {
						event["ds"] = ds
					}
				}

				// 应用过滤条件
				match := true
				for k, v := range filter {
					curKOk := false
					if eventVal, exists := event[k]; exists {
						if listVal, isList := v.([]any); isList {
							for _, listItem := range listVal {
								if fmt.Sprintf("%v", eventVal) == fmt.Sprintf("%v", listItem) {
									curKOk = true
									break
								}
							}
						} else if fmt.Sprintf("%v", eventVal) == fmt.Sprintf("%v", v) {
							curKOk = true
						}
					}
					if !curKOk {
						match = false
						break
					}
				}

				if match {
					got++
					records = append(records, event)
				}
			}
		}

		pnum++
	}

	logx.DebugfM(ctx, "Redash", "扫描了 %d 个事件后获取了 %d 个结果", scanned, got)
	return records, nil
}

// ExecQuery 执行Redash中的查询语句
func ExecQuery(url string, appKey string, ds int, sql string) (map[string]any, error) {
	redash := NewRedash(url, appKey)
	return redash.CreateQuery(fmt.Sprintf("%d", ds), sql, nil)
}
