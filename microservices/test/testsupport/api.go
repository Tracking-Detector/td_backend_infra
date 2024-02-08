package testsupport

import (
	"bytes"
	"fmt"
	"io"

	"net/http"
)

type Response struct {
	Body       string
	StatusCode int
}

func Get(url string) (*Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body failed: %v", err)
	}

	return &Response{
		Body:       string(body),
		StatusCode: resp.StatusCode,
	}, nil
}

func Post(url, body, contentType string) (*Response, error) {
	resp, err := http.Post(url, contentType, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return nil, fmt.Errorf("POST request failed: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body failed: %v", err)
	}

	return &Response{
		Body:       string(responseBody),
		StatusCode: resp.StatusCode,
	}, nil
}
