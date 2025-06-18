package assert_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/vaguevoid/cloud-cli/internal/lib/httpx"
	"github.com/vaguevoid/cloud-cli/internal/test/assert"
)

func TestAssertionsForCodeCoverage(t *testing.T) {
	assert.True(t, true)
	assert.False(t, false)
	assert.Nil(t, nil)
	assert.NotNil(t, 42)
	assert.Empty(t, "")
	assert.NotEmpty(t, "yolo")
	assert.Error(t, "uh oh", errors.New("uh oh"))
	assert.Equal(t, 42, 42)
	assert.Regexp(t, "hello", "hello")
	assert.Length(t, 2, []int{1, 2})
}

func TestAssertionsForRequests(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.RequestMethodEqual(t, http.MethodGet, r)
		assert.RequestPathEqual(t, "/path/to/route", r)
		assert.RequestHeaderEqual(t, "unit testing", httpx.HeaderUserAgent, r)
		assert.RequestBodyEqual(t, "hello world", r)
		httpx.RespondOk("ping", w)
	}))

	url := fmt.Sprintf("%s/path/to/route", mockServer.URL)
	body := strings.NewReader("hello world")
	req, err := http.NewRequest(http.MethodGet, url, body)
	assert.Nil(t, err)
	assert.NotNil(t, req)
	req.Header.Set(httpx.HeaderUserAgent, "unit testing")

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	defer resp.Body.Close()

	assert.ResponseStatusCode(t, http.StatusOK, resp)
	assert.ResponseBodyEqual(t, "ping", resp)
	assert.ResponseHeaderEqual(t, httpx.ContentTypeText, httpx.HeaderContentType, resp)
}
