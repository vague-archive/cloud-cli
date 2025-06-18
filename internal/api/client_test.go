package api_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/vaguevoid/cloud-cli/internal/api"
	"github.com/vaguevoid/cloud-cli/internal/lib/httpx"
	"github.com/vaguevoid/cloud-cli/internal/test/assert"
	"github.com/vaguevoid/cloud-cli/internal/test/mock"
)

const TestToken = "header.payload.signature"

//-------------------------------------------------------------------------------------------------

func TestClientGet(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/action/route", r.URL.Path)
		assert.Equal(t, "Bearer header.payload.signature", r.Header.Get(httpx.HeaderAuthorization))
		w.WriteHeader(http.StatusOK)
	}))

	api, err := api.NewClient(mockServer.URL, TestToken)
	assert.Nil(t, err)
	assert.NotNil(t, api)

	resp, err := api.Get("action/route")
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()
}

//-------------------------------------------------------------------------------------------------

func TestClientPost(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/action/route", r.URL.Path)
		assert.Equal(t, "Bearer header.payload.signature", r.Header.Get(httpx.HeaderAuthorization))
		w.WriteHeader(http.StatusOK)
	}))

	api, err := api.NewClient(mockServer.URL, TestToken)
	assert.Nil(t, err)
	assert.NotNil(t, api)

	resp, err := api.Post("action/route", nil)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()
}

//-------------------------------------------------------------------------------------------------

func TestClientPostJSON(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/action/route", r.URL.Path)
		assert.Equal(t, "Bearer header.payload.signature", r.Header.Get(httpx.HeaderAuthorization))
		assert.Equal(t, httpx.ContentTypeJSON, r.Header.Get(httpx.HeaderContentType))
		bodyBytes, err := io.ReadAll(r.Body)
		assert.Nil(t, err)
		assert.Equal(t, `["foo","bar"]`, string(bodyBytes))
		w.WriteHeader(http.StatusOK)
	}))

	api, err := api.NewClient(mockServer.URL, TestToken)
	assert.Nil(t, err)
	assert.NotNil(t, api)

	resp, err := api.PostJSON("action/route", []string{"foo", "bar"})
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()
}

//-------------------------------------------------------------------------------------------------

func TestClientPostFILE(t *testing.T) {
	content := "Hello World"
	path := "path/to/hello.txt"

	tmp := mock.TempDir(t)
	tmp.AddTextFile(t, path, content)

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/action/route", r.URL.Path)
		assert.Equal(t, "Bearer header.payload.signature", r.Header.Get(httpx.HeaderAuthorization))
		assert.Equal(t, httpx.ContentTypeBytes, r.Header.Get(httpx.HeaderContentType))
		assert.RequestBodyEqual(t, content, r)
		w.WriteHeader(http.StatusOK)
	}))

	api, err := api.NewClient(mockServer.URL, TestToken)
	assert.Nil(t, err)
	assert.NotNil(t, api)

	resp, err := api.PostFILE("action/route", filepath.Join(tmp.Dir, path))
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()
}

//-------------------------------------------------------------------------------------------------
