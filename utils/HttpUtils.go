package utils

import (
	"BronyaBot/global"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"io"
	"net/http"
	"strings"
	"time"
)

type HttpResponse struct {
	Response *resty.Response
	Error    error
}

type HttpUtils struct {
	Client *resty.Client
}

var client = &http.Client{}

func (c *HttpUtils) AsyncPostWithTimeout(
	url string,
	headers map[string]string,
	params interface{},
	timeout time.Duration,
) (*HttpResponse, error) {
	resultChan := make(chan HttpResponse, 1)
	r := c.Client.R()
	go func() {
		if headers != nil {
			r.SetHeaders(headers)
		}
		if params != nil {
			r.SetBody(params)
		}
		r.SetCookie(&http.Cookie{})
		resp, err := r.Post(url)
		resultChan <- HttpResponse{Response: resp, Error: err}
	}()
	select {
	case result := <-resultChan:
		return &result, nil
	case <-time.After(timeout):
		return nil, errors.New("request timed out")
	}
}

func SendRequest(method, url string, payload map[string]string, header http.Header) ([]byte, error) {
	marshal, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %v", err)
	}
	request, err := http.NewRequest(method, url, io.NopCloser(strings.NewReader(string(marshal))))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	request.Header = header

	resp, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			global.Log.Info("failed to close response body")
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}
	return body, nil
}
