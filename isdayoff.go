package isdayoff

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client for requests to isdayoff.ru
type Client struct {
	httpClient *http.Client
}

// New initiates client with default http client
func New() *Client {
	return NewWithClient(http.DefaultClient)
}

// NewWithClient initiates client with custom http client
func NewWithClient(client *http.Client) *Client {
	return &Client{client}
}

// IsLeap checks if year is leap
func (c *Client) IsLeap(year int) (bool, error) {
	url := fmt.Sprintf("https://isdayoff.ru/api/isleap?year=%d", year)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return false, fmt.Errorf("http.NewRequest failed: %w", err)
	}
	req.Header.Set("User-Agent", "isdayoff-golang-lib/1.0.0 (https://github.com/kotopheiop)")
	res, err := c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("client.Do(req) failed: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return false, fmt.Errorf("io.ReadAll failed: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return false, parseAPIError(res.StatusCode, body)
	}

	return YearType(string(body)) == YearTypeLeap, nil
}

var boolToStr = map[bool]string{
	false: "0",
	true:  "1",
}

// APIError represents an error returned by the API
type APIError struct {
	Code    ErrorCode
	Message string
	Status  int
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error %s (HTTP %d): %s", e.Code, e.Status, e.Message)
}

// parseAPIError parses API error response according to API documentation
func parseAPIError(statusCode int, body []byte) error {
	bodyStr := strings.TrimSpace(string(body))

	// Check for known error codes
	switch bodyStr {
	case string(ErrorCodeWrongDate):
		return &APIError{
			Code:    ErrorCodeWrongDate,
			Message: "Ошибка в дате",
			Status:  statusCode,
		}
	case string(ErrorCodeNotFound):
		return &APIError{
			Code:    ErrorCodeNotFound,
			Message: "Данные не найдены",
			Status:  statusCode,
		}
	case string(ErrorCodeInternalError):
		return &APIError{
			Code:    ErrorCodeInternalError,
			Message: "Ошибка сервиса",
			Status:  statusCode,
		}
	default:
		return fmt.Errorf("unexpected status code %d: %s", statusCode, bodyStr)
	}
}

// Params contains various filters for request
type Params struct {
	Year        int
	Month       *time.Month
	Day         *int
	CountryCode *CountryCode
	Pre         *bool // помечать сокращённые рабочие дни цифрой 2
	Covid       *bool // помечать рабочие дни цифрой 4 (в связи с пандемией COVID-19)
	SixDayWeek  *bool // считать, что неделя шестидневная (sd)
	TZ          *string
}

// GetBy Get data by particular params
func (c *Client) GetBy(params Params) ([]DayType, error) {
	baseURL := "https://isdayoff.ru/api/getdata"
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	q := u.Query()
	q.Set("year", fmt.Sprintf("%d", params.Year))

	if params.Month != nil {
		q.Set("month", fmt.Sprintf("%02d", *params.Month))
	}
	if params.Day != nil {
		q.Set("day", fmt.Sprintf("%02d", *params.Day))
	}
	if params.CountryCode != nil {
		q.Set("cc", string(*params.CountryCode))
	}
	if params.Pre != nil {
		q.Set("pre", boolToStr[*params.Pre])
	}
	if params.Covid != nil {
		q.Set("covid", boolToStr[*params.Covid])
	}
	if params.SixDayWeek != nil {
		q.Set("sd", boolToStr[*params.SixDayWeek])
	}
	if params.TZ != nil {
		q.Set("tz", *params.TZ)
	}

	u.RawQuery = q.Encode()
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest failed: %w", err)
	}

	req.Header.Set("User-Agent", "isdayoff-golang-lib/1.0.0 (https://github.com/kotopheiop)")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client.Do(req) failed: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll failed: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, parseAPIError(res.StatusCode, body)
	}
	result := []DayType{}

	bodyStr := string(body)
	for _, char := range bodyStr {
		result = append(result, DayType(string(char)))
	}

	return result, nil
}

// GetByPeriod Get data for arbitrary period (date1 to date2)
// Maximum 366 days can be requested
// date1 and date2 should be in format YYYYMMDD
func (c *Client) GetByPeriod(date1, date2 string, params Params) ([]DayType, error) {
	baseURL := "https://isdayoff.ru/api/getdata"
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	q := u.Query()
	q.Set("date1", date1)
	q.Set("date2", date2)

	if params.CountryCode != nil {
		q.Set("cc", string(*params.CountryCode))
	}
	if params.Pre != nil {
		q.Set("pre", boolToStr[*params.Pre])
	}
	if params.Covid != nil {
		q.Set("covid", boolToStr[*params.Covid])
	}
	if params.SixDayWeek != nil {
		q.Set("sd", boolToStr[*params.SixDayWeek])
	}

	u.RawQuery = q.Encode()
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest failed: %w", err)
	}

	req.Header.Set("User-Agent", "isdayoff-golang-lib/1.0.0 (https://github.com/kotopheiop)")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client.Do(req) failed: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll failed: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, parseAPIError(res.StatusCode, body)
	}
	result := []DayType{}

	bodyStr := string(body)
	for _, char := range bodyStr {
		result = append(result, DayType(string(char)))
	}

	return result, nil
}

// Today get data for today by particular params
func (c *Client) Today(params Params) (*DayType, error) {
	return c.aliasRequest("today", params)
}

// Tomorrow get data for tomorrow by particular params
func (c *Client) Tomorrow(params Params) (*DayType, error) {
	return c.aliasRequest("tomorrow", params)
}

func (c *Client) aliasRequest(alias string, params Params) (*DayType, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://isdayoff.ru/%s", alias), nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest failed: %w", err)
	}

	q := req.URL.Query()
	if params.CountryCode != nil {
		q.Add("cc", string(*params.CountryCode))
	}
	if params.Pre != nil {
		q.Add("pre", boolToStr[*params.Pre])
	}
	if params.Covid != nil {
		q.Add("covid", boolToStr[*params.Covid])
	}
	if params.SixDayWeek != nil {
		q.Add("sd", boolToStr[*params.SixDayWeek])
	}
	if params.TZ != nil {
		q.Add("tz", *params.TZ)
	}

	req.URL.RawQuery = q.Encode()

	req.Header.Set("User-Agent", "isdayoff-golang-lib/1.0.2 (https://github.com/kotopheiop)")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client.Do(req) failed: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll failed: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, parseAPIError(res.StatusCode, body)
	}

	bodyStr := strings.TrimSpace(string(body))
	result := DayType(bodyStr)

	return &result, nil
}
