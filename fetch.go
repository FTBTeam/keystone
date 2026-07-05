package keystone

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Headers map[string][]string

// MakeRequest makes an HTTP request with the given method, url, body, and headers.
func MakeRequest(method, url string, body []byte, headers Headers) ([]byte, error) {
	client := http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	requestHeaders := Headers{
		"Accept": {"application/json"},
	}

	for key, values := range headers {
		requestHeaders[key] = values
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, errors.New(fmt.Sprintf("Request failed %s with status code %d and status %s", url, resp.StatusCode, resp.Status))
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// GetAsJsonBasic makes a GET request to the given URL and unmarshal the JSON response into the specified type V.
func GetAsJsonBasic[V any](url string) (V, error) {
	return GetAsJson[V](url, nil)
}

// GetAsJson makes a GET request to the given URL with the specified headers and unmarshals the JSON response into the specified type V.
func GetAsJson[V any](url string, headers Headers) (V, error) {
	var nothing V

	body, err := MakeRequest("GET", url, nil, headers)
	if err != nil {
		return nothing, err
	}

	var data V
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nothing, err
	}

	return data, nil
}
