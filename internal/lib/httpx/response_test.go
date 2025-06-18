package httpx_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vaguevoid/cloud-cli/internal/lib/httpx"
	"github.com/vaguevoid/cloud-cli/internal/test/assert"
)

//-------------------------------------------------------------------------------------------------

func TestRespondOk(t *testing.T) {

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpx.RespondOk("ping", w)
	}))

	resp, err := http.Get(mockServer.URL)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	defer resp.Body.Close()

	assert.ResponseStatusCode(t, http.StatusOK, resp)
	assert.ResponseBodyEqual(t, "ping", resp)
	assert.ResponseHeaderEqual(t, httpx.ContentTypeText, httpx.HeaderContentType, resp)

}

//-------------------------------------------------------------------------------------------------

func TestRespondOkJson(t *testing.T) {

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}{
			ID:   123,
			Name: "Jake",
		}
		httpx.RespondOk(user, w)
	}))

	resp, err := http.Get(mockServer.URL)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	defer resp.Body.Close()

	assert.ResponseStatusCode(t, http.StatusOK, resp)
	assert.ResponseBodyEqual(t, `{"id":123,"name":"Jake"}`, resp)
	assert.ResponseHeaderEqual(t, httpx.ContentTypeJSON, httpx.HeaderContentType, resp)

}

//-------------------------------------------------------------------------------------------------
