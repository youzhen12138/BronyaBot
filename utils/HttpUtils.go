package utils

import (
	"BronyaBot/global"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var client = &http.Client{}

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
