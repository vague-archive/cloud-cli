package account_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/vaguevoid/cloud-cli/internal/domain/account"
	"github.com/vaguevoid/cloud-cli/internal/lib/httpx"
	"github.com/vaguevoid/cloud-cli/internal/test/assert"
	"github.com/vaguevoid/cloud-cli/internal/test/mock"
)

const TestServer = "https://test.void.dev/"

//-------------------------------------------------------------------------------------------------

func TestLoginMissingServer(t *testing.T) {
	user, err := account.Login(&account.LoginCommand{})
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.Equal(t, "missing server", err.Error())
}

//-------------------------------------------------------------------------------------------------

func TestLoginMissingRuntime(t *testing.T) {
	user, err := account.Login(&account.LoginCommand{
		Server: TestServer,
	})
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.Equal(t, "missing runtime", err.Error())
}

//-------------------------------------------------------------------------------------------------

func TestLoginMissingKeyring(t *testing.T) {
	user, err := account.Login(&account.LoginCommand{
		Server:  TestServer,
		Runtime: mock.Runtime(),
	})
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.Equal(t, "missing keyring", err.Error())
}

//-------------------------------------------------------------------------------------------------

func TestLoginSuccess(t *testing.T) {

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer header.payload.signature", r.Header.Get(httpx.HeaderAuthorization))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":100,"name":"TestLoginSuccess"}`))
	}))

	runtime := mock.Runtime()
	keyring := mock.Keyring()
	resultChannel := make(chan *account.User, 1)

	assert.False(t, keyring.Has(httpx.ParamJWT), "preconditions")

	go func() {
		user, err := account.Login(&account.LoginCommand{
			Server:  mockServer.URL,
			Runtime: runtime,
			Keyring: keyring,
		})
		assert.Nil(t, err)
		assert.NotNil(t, user)
		resultChannel <- user
	}()
	briefPause()

	assert.Regexp(t, `http://127.0.0.1:\d*/login\?cli=true&origin=http%3A%2F%2F127.0.0.1%3A\d*%2Fcallback`, runtime.OpenedURL)
	opened, err := url.Parse(runtime.OpenedURL)
	assert.Nil(t, err)
	origin := opened.Query().Get("origin")

	client := &http.Client{Timeout: 10 * time.Millisecond}
	resp, err := client.Get(fmt.Sprintf("%s?%s=header.payload.signature", origin, httpx.ParamJWT))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	user := <-resultChannel
	assert.Equal(t, 100, user.ID)
	assert.Equal(t, t.Name(), user.Name)

	savedJwt, ok := keyring.Get(httpx.ParamJWT)
	assert.True(t, ok)
	assert.Equal(t, "header.payload.signature", savedJwt)
}

//-------------------------------------------------------------------------------------------------

func TestLoginJWTAlreadyInKeyring(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer existing.value", r.Header.Get(httpx.HeaderAuthorization))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":200,"name":"TestLoginJWTAlreadyInKeyring"}`))
	}))

	runtime := mock.Runtime()
	keyring := mock.Keyring()
	resultChannel := make(chan *account.User, 1)

	keyring.Set(httpx.ParamJWT, "existing.value")

	go func() {
		user, err := account.Login(&account.LoginCommand{
			Server:  mockServer.URL,
			Runtime: runtime,
			Keyring: keyring,
		})
		assert.Nil(t, err)
		resultChannel <- user
	}()

	user := <-resultChannel
	assert.Equal(t, 200, user.ID)
	assert.Equal(t, t.Name(), user.Name)
	assert.Empty(t, runtime.OpenedURL, "browser was never opened")
}

//-------------------------------------------------------------------------------------------------

func TestLoginInvalidJWTAlreadyInKeyring(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get(httpx.HeaderAuthorization)
		if auth == "Bearer old.jwt" {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			assert.Equal(t, "Bearer new.jwt", auth)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id":300,"name":"TestLoginInvalidJWTAlreadyInKeyring"}`))
		}
	}))

	runtime := mock.Runtime()
	keyring := mock.Keyring()
	resultChannel := make(chan *account.User, 1)

	keyring.Set(httpx.ParamJWT, "old.jwt")

	go func() {
		user, err := account.Login(&account.LoginCommand{
			Server:  mockServer.URL,
			Runtime: runtime,
			Keyring: keyring,
		})
		assert.Nil(t, err)
		resultChannel <- user
	}()
	briefPause()

	assert.Regexp(t, `http://127.0.0.1:\d*/login\?cli=true&origin=http%3A%2F%2F127.0.0.1%3A\d*%2Fcallback`, runtime.OpenedURL)
	opened, err := url.Parse(runtime.OpenedURL)
	assert.Nil(t, err)
	origin := opened.Query().Get("origin")

	client := &http.Client{Timeout: 10 * time.Millisecond}
	resp, err := client.Get(fmt.Sprintf("%s?%s=new.jwt", origin, httpx.ParamJWT))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	user := <-resultChannel
	assert.Equal(t, 300, user.ID)
	assert.Equal(t, t.Name(), user.Name)

	savedJwt, ok := keyring.Get(httpx.ParamJWT)
	assert.True(t, ok)
	assert.Equal(t, "new.jwt", savedJwt)
}

//-------------------------------------------------------------------------------------------------

func TestLoginInvalidJWTReturnedFromServer(t *testing.T) {

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer header.payload.signature", r.Header.Get(httpx.HeaderAuthorization))
		w.WriteHeader(http.StatusUnauthorized)
	}))

	runtime := mock.Runtime()
	keyring := mock.Keyring()
	errorChannel := make(chan error, 1)

	assert.False(t, keyring.Has(httpx.ParamJWT), "preconditions")

	go func() {
		user, err := account.Login(&account.LoginCommand{
			Server:  mockServer.URL,
			Runtime: runtime,
			Keyring: keyring,
		})
		assert.Nil(t, user)
		assert.NotNil(t, err)
		errorChannel <- err
	}()
	briefPause()

	assert.Regexp(t, `http://127.0.0.1:\d*/login\?cli=true&origin=http%3A%2F%2F127.0.0.1%3A\d*%2Fcallback`, runtime.OpenedURL)
	opened, err := url.Parse(runtime.OpenedURL)
	assert.Nil(t, err)
	origin := opened.Query().Get("origin")

	client := &http.Client{Timeout: 10 * time.Millisecond}
	resp, err := client.Get(fmt.Sprintf("%s?%s=header.payload.signature", origin, httpx.ParamJWT))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	err = <-errorChannel
	assert.Equal(t, "unauthorized", err.Error())
	assert.False(t, keyring.Has(httpx.ParamJWT))
}

//-------------------------------------------------------------------------------------------------

func TestLoginInvalidStatusCodeReturnedFromServer(t *testing.T) {

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer header.payload.signature", r.Header.Get(httpx.HeaderAuthorization))
		w.WriteHeader(http.StatusInternalServerError)
	}))

	runtime := mock.Runtime()
	keyring := mock.Keyring()
	errorChannel := make(chan error, 1)

	assert.False(t, keyring.Has(httpx.ParamJWT), "preconditions")

	go func() {
		user, err := account.Login(&account.LoginCommand{
			Server:  mockServer.URL,
			Runtime: runtime,
			Keyring: keyring,
		})
		assert.Nil(t, user)
		assert.NotNil(t, err)
		errorChannel <- err
	}()
	briefPause()

	assert.Regexp(t, `http://127.0.0.1:\d*/login\?cli=true&origin=http%3A%2F%2F127.0.0.1%3A\d*%2Fcallback`, runtime.OpenedURL)
	opened, err := url.Parse(runtime.OpenedURL)
	assert.Nil(t, err)
	origin := opened.Query().Get("origin")

	client := &http.Client{Timeout: 10 * time.Millisecond}
	resp, err := client.Get(fmt.Sprintf("%s?%s=header.payload.signature", origin, httpx.ParamJWT))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	err = <-errorChannel
	assert.Equal(t, "unexpected status code 500", err.Error())
	assert.False(t, keyring.Has(httpx.ParamJWT))
}

//-------------------------------------------------------------------------------------------------

func TestLoginInvalidJSONReturnedFromServer(t *testing.T) {

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer header.payload.signature", r.Header.Get(httpx.HeaderAuthorization))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid JSON`))
	}))

	runtime := mock.Runtime()
	keyring := mock.Keyring()
	errorChannel := make(chan error, 1)

	assert.False(t, keyring.Has(httpx.ParamJWT), "preconditions")

	go func() {
		user, err := account.Login(&account.LoginCommand{
			Server:  mockServer.URL,
			Runtime: runtime,
			Keyring: keyring,
		})
		assert.Nil(t, user)
		assert.NotNil(t, err)
		errorChannel <- err
	}()
	briefPause()

	assert.Regexp(t, `http://127.0.0.1:\d*/login\?cli=true&origin=http%3A%2F%2F127.0.0.1%3A\d*%2Fcallback`, runtime.OpenedURL)
	opened, err := url.Parse(runtime.OpenedURL)
	assert.Nil(t, err)
	origin := opened.Query().Get("origin")

	client := &http.Client{Timeout: 10 * time.Millisecond}
	resp, err := client.Get(fmt.Sprintf("%s?%s=header.payload.signature", origin, httpx.ParamJWT))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	err = <-errorChannel
	assert.Equal(t, "unexpected JSON response: invalid character 'i' looking for beginning of value", err.Error())
	assert.False(t, keyring.Has(httpx.ParamJWT))
}

//-------------------------------------------------------------------------------------------------

func TestLoginTimeout(t *testing.T) {
	runtime := mock.Runtime()
	keyring := mock.Keyring()
	timeout := 5 * time.Millisecond
	errorChannel := make(chan error, 1)

	assert.False(t, keyring.Has(httpx.ParamJWT), "preconditions")

	go func() {
		jwt, err := account.Login(&account.LoginCommand{
			Server:  TestServer,
			Runtime: runtime,
			Keyring: keyring,
			Timeout: timeout,
		})
		assert.Empty(t, jwt)
		assert.NotNil(t, err)
		errorChannel <- err
	}()
	briefPause()

	assert.Regexp(t, `https://test.void.dev/login\?cli=true&origin=http%3A%2F%2F127.0.0.1%3A\d*%2Fcallback`, runtime.OpenedURL)

	err := <-errorChannel
	assert.Equal(t, "login timed out", err.Error())

	assert.False(t, keyring.Has(httpx.ParamJWT))
}

//-------------------------------------------------------------------------------------------------

func TestLoginMissingJWT(t *testing.T) {
	runtime := mock.Runtime()
	keyring := mock.Keyring()
	errorChannel := make(chan error, 1)

	assert.False(t, keyring.Has(httpx.ParamJWT), "preconditions")

	go func() {
		jwt, err := account.Login(&account.LoginCommand{
			Server:  TestServer,
			Runtime: runtime,
			Keyring: keyring,
		})
		assert.NotNil(t, err)
		assert.Empty(t, jwt)
		errorChannel <- err
	}()
	briefPause()

	assert.Regexp(t, `https://test.void.dev/login\?cli=true&origin=http%3A%2F%2F127.0.0.1%3A\d*%2Fcallback`, runtime.OpenedURL)
	opened, err := url.Parse(runtime.OpenedURL)
	assert.Nil(t, err)
	origin := opened.Query().Get("origin")

	client := &http.Client{Timeout: 10 * time.Millisecond}
	resp, err := client.Get(origin)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	defer resp.Body.Close()

	err = <-errorChannel
	assert.Equal(t, "missing jwt in callback", err.Error())

	assert.False(t, keyring.Has(httpx.ParamJWT))
}

//-------------------------------------------------------------------------------------------------

func briefPause() {
	time.Sleep(10 * time.Millisecond) // little bit sketch, but need time for Login() to spin up it's local http server
}

//-------------------------------------------------------------------------------------------------
