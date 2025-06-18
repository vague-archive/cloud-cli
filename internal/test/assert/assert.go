package assert

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
)

func Fail(t *testing.T, failureMessage string, msgAndArgs ...any) {
	require.Fail(t, failureMessage, msgAndArgs)
}

func True(t *testing.T, value bool, msgAndArgs ...any) {
	require.True(t, value, msgAndArgs)
}

func False(t *testing.T, value bool, msgAndArgs ...any) {
	require.False(t, value, msgAndArgs)
}

func Nil(t *testing.T, value any, msgAndArgs ...any) {
	require.Nil(t, value, msgAndArgs)
}

func NotNil(t *testing.T, value any, msgAndArgs ...any) {
	require.NotNil(t, value, msgAndArgs)
}

func Empty(t *testing.T, value any, msgAndArgs ...any) {
	require.Empty(t, value, msgAndArgs)
}

func NotEmpty(t *testing.T, value any, msgAndArgs ...any) {
	require.NotEmpty(t, value, msgAndArgs)
}

func Error(t *testing.T, expected string, actual error) {
	require.Equal(t, expected, actual.Error())
}

func NoError(t *testing.T, actual error) {
	if actual != nil {
		Fail(t, fmt.Sprintf("Expected nil, but got Error: %s", actual.Error()))
	}
}

func Equal[T any](t *testing.T, expected T, actual T, msgAndArgs ...any) {
	require.Equal(t, expected, actual, msgAndArgs)
}

func Regexp(t *testing.T, rx any, str any, msgAndArgs ...any) {
	require.Regexp(t, rx, str, msgAndArgs)
}

func Length(t *testing.T, expected int, actual any, msgAndArgs ...any) {
	require.Len(t, actual, expected, msgAndArgs)
}

//=================================================================================================
// HTTP REQUEST ASSERTIONS
//=================================================================================================

func RequestBody(t *testing.T, r *http.Request) string {
	bytes, err := io.ReadAll(r.Body)
	assert.Nil(t, err)
	assert.NotNil(t, bytes)
	return string(bytes)
}

func RequestJSON[T any](t *testing.T, r *http.Request) T {
	var result T
	err := json.NewDecoder(r.Body).Decode(&result)
	assert.Nil(t, err)
	return result
}

func RequestBodyEqual(t *testing.T, expected string, r *http.Request) {
	assert.Equal(t, expected, RequestBody(t, r))
}

func RequestJSONEqual[T any](t *testing.T, expected T, r *http.Request) {
	assert.Equal(t, expected, RequestJSON[T](t, r))
}

func RequestMethodEqual(t *testing.T, expected string, r *http.Request) {
	assert.Equal(t, expected, r.Method)
}

func RequestPathEqual(t *testing.T, expected string, r *http.Request) {
	assert.Equal(t, expected, r.URL.Path)
}

func RequestHeaderEqual(t *testing.T, expected string, header string, r *http.Request) {
	assert.Equal(t, expected, r.Header.Get(header))
}

//=================================================================================================
// HTTP RESPONSE ASSERTIONS
//=================================================================================================

func ResponseStatusCode(t *testing.T, expected int, r *http.Response) {
	assert.Equal(t, expected, r.StatusCode)
}

func ResponseBody(t *testing.T, r *http.Response) string {
	bytes, err := io.ReadAll(r.Body)
	assert.Nil(t, err)
	assert.NotNil(t, bytes)
	return string(bytes)
}

func ResponseJSON[T any](t *testing.T, r *http.Response) T {
	var result T
	err := json.NewDecoder(r.Body).Decode(&result)
	assert.Nil(t, err)
	return result
}

func ResponseBodyEqual(t *testing.T, expected string, r *http.Response) {
	assert.Equal(t, expected, ResponseBody(t, r))
}

func ResponseJSONEqual[T any](t *testing.T, expected T, r *http.Response) {
	assert.Equal(t, expected, ResponseJSON[T](t, r))
}

func ResponseHeaderEqual(t *testing.T, expected string, header string, r *http.Response) {
	assert.Equal(t, expected, r.Header.Get(header))
}

//-------------------------------------------------------------------------------------------------
