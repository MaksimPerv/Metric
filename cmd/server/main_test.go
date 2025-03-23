package main

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	req.Header.Add("Content-Type", "text/plain")
	assert.NoError(t, err)
	resp, err := ts.Client().Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	return resp, string(respBody)
}

func TestServer(t *testing.T) {
	type want struct {
		statusCode int
		statusType string
		body       string
	}
	ts := httptest.NewServer(run())
	defer ts.Close()

	tests := []struct {
		url    string
		method string
		want   want
	}{
		{
			url:    "/update/gauge/asd/123",
			method: "POST",
			want: want{
				statusCode: 200,
				statusType: "text/plain; charset=utf-8",
				body:       "Metric updated\n",
			},
		},
		{
			url:    "/update/gauge/asd/123",
			method: "GET",
			want: want{
				statusCode: 405,
			},
		},
	}

	for _, test := range tests {
		resp, respBody := testRequest(t, ts, test.method, test.url)
		assert.Equal(t, resp.StatusCode, test.want.statusCode)
		assert.Equal(t, resp.Header.Get("Content-Type"), test.want.statusType)
		assert.Equal(t, respBody, test.want.body)
	}
}
