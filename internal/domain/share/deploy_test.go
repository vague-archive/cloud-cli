package share_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/vaguevoid/cloud-cli/internal/api"
	"github.com/vaguevoid/cloud-cli/internal/domain/share"
	"github.com/vaguevoid/cloud-cli/internal/lib/crypto"
	"github.com/vaguevoid/cloud-cli/internal/lib/httpx"
	"github.com/vaguevoid/cloud-cli/internal/test/assert"
	"github.com/vaguevoid/cloud-cli/internal/test/mock"
)

//-------------------------------------------------------------------------------------------------

const (
	TestDeployID   = 42
	TestDeploySlug = "latest"
	TestDeployURL  = "https://test.void.dev/void/snakes/latest"
	TestServer     = "https://test.void.dev/"
	TestOrg        = "void"
	TestGame       = "snakes"
	TestLabel      = "latest"
	TestPath       = "."
	TestToken      = "personal-access-token"
	FirstContent   = "first"
	SecondContent  = "second"
	ThirdContent   = "third"
	FirstPath      = "path/to/first.txt"
	SecondPath     = "path/to/second.txt"
	ThirdPath      = "path/to/third.txt"
)

func makeAPI(t *testing.T) *api.Client {
	api, err := api.NewClient(TestServer, TestToken)
	assert.Nil(t, err)
	assert.NotNil(t, api)
	return api
}

//-------------------------------------------------------------------------------------------------

func TestLoginMissingApi(t *testing.T) {
	_, err := share.Deploy(&share.DeployCommand{
		Org:  TestOrg,
		Game: TestGame,
		Path: TestPath,
	})
	assert.NotNil(t, err)
	assert.Error(t, "missing api client", err)
}

//-------------------------------------------------------------------------------------------------

func TestLoginMissingOrg(t *testing.T) {
	api := makeAPI(t)
	_, err := share.Deploy(&share.DeployCommand{
		API:  api,
		Game: TestGame,
		Path: TestPath,
	})
	assert.NotNil(t, err)
	assert.Error(t, "missing organization", err)
}

//-------------------------------------------------------------------------------------------------

func TestLoginMissingGame(t *testing.T) {
	api := makeAPI(t)
	_, err := share.Deploy(&share.DeployCommand{
		API:  api,
		Org:  TestOrg,
		Path: TestPath,
	})
	assert.NotNil(t, err)
	assert.Error(t, "missing game", err)
}

//-------------------------------------------------------------------------------------------------

func TestLoginMissingPath(t *testing.T) {
	api := makeAPI(t)
	_, err := share.Deploy(&share.DeployCommand{
		API:  api,
		Org:  TestOrg,
		Game: TestGame,
	})
	assert.NotNil(t, err)
	assert.Error(t, "missing path", err)
}

//-------------------------------------------------------------------------------------------------

func TestLoginPathNotFound(t *testing.T) {
	api := makeAPI(t)
	_, err := share.Deploy(&share.DeployCommand{
		API:  api,
		Org:  TestOrg,
		Game: TestGame,
		Path: "path/to/unknown/file.txt",
	})
	assert.NotNil(t, err)
	assert.Error(t, "directory not found path/to/unknown/file.txt", err)
}

//-------------------------------------------------------------------------------------------------

func TestFullDeploy(t *testing.T) {

	mockDir := mock.TempDir(t)
	mockDir.AddTextFile(t, FirstPath, FirstContent)
	mockDir.AddTextFile(t, SecondPath, SecondContent)
	mockDir.AddTextFile(t, ThirdPath, ThirdContent)

	expectedManifest := []share.DeployEntry{
		{
			Path:          FirstPath,
			Blake3:        crypto.Blake3(FirstContent),
			ContentLength: len(FirstContent),
		},
		{
			Path:          SecondPath,
			Blake3:        crypto.Blake3(SecondContent),
			ContentLength: len(SecondContent),
		},
		{
			Path:          ThirdPath,
			Blake3:        crypto.Blake3(ThirdContent),
			ContentLength: len(ThirdContent),
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.RequestMethodEqual(t, http.MethodPost, r)
		assert.RequestHeaderEqual(t, "Bearer personal-access-token", httpx.HeaderAuthorization, r)
		if r.URL.Path == "/api/void/snakes/deploy" {
			manifest := assert.RequestJSON[[]share.DeployEntry](t, r)
			assert.Equal(t, expectedManifest, manifest)
			w.Header().Add(httpx.HeaderXDeployID, fmt.Sprintf("%d", TestDeployID))
			httpx.RespondAccepted(manifest, w)
		} else if r.URL.Path == "/api/void/snakes/deploy/42/activate" {
			httpx.RespondOk(&share.DeployResult{
				DeployID: TestDeployID,
				Slug:     TestDeploySlug,
				URL:      TestDeployURL,
				Manifest: expectedManifest,
			}, w)
		} else if strings.HasPrefix(r.URL.Path, "/api/void/snakes/deploy/42") {
			httpx.RespondOk("ok", w)
		} else {
			httpx.RespondBadRequest(fmt.Sprintf("unexpected %s", r.URL.Path), w)
		}
	}))

	api, err := api.NewClient(mockServer.URL, TestToken)
	assert.NoError(t, err)

	uploads := make([]string, 0)
	onUpload := func(deployID int64, path string) {
		uploads = append(uploads, path)
	}

	result, err := share.Deploy(&share.DeployCommand{
		API:      api,
		Org:      TestOrg,
		Game:     TestGame,
		Path:     mockDir.Dir,
		OnUpload: onUpload,
	})

	assert.NoError(t, err)
	assert.Equal(t, TestDeployID, result.DeployID)
	assert.Equal(t, TestDeploySlug, result.Slug)
	assert.Equal(t, TestDeployURL, result.URL)
	assert.Equal(t, expectedManifest, result.Manifest)

	assert.Equal(t, []string{
		"path/to/first.txt",
		"path/to/second.txt",
		"path/to/third.txt",
	}, uploads)
}

//-------------------------------------------------------------------------------------------------

func TestIncrementalDeploy(t *testing.T) {

	mockDir := mock.TempDir(t)
	mockDir.AddTextFile(t, FirstPath, FirstContent)
	mockDir.AddTextFile(t, SecondPath, SecondContent)
	mockDir.AddTextFile(t, ThirdPath, ThirdContent)

	expectedManifest := []share.DeployEntry{
		{
			Path:          FirstPath,
			Blake3:        crypto.Blake3(FirstContent),
			ContentLength: len(FirstContent),
		},
		{
			Path:          SecondPath,
			Blake3:        crypto.Blake3(SecondContent),
			ContentLength: len(SecondContent),
		},
		{
			Path:          ThirdPath,
			Blake3:        crypto.Blake3(ThirdContent),
			ContentLength: len(ThirdContent),
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.RequestMethodEqual(t, http.MethodPost, r)
		assert.RequestHeaderEqual(t, "Bearer personal-access-token", httpx.HeaderAuthorization, r)
		if r.URL.Path == "/api/void/snakes/deploy" {
			manifest := assert.RequestJSON[[]share.DeployEntry](t, r)
			assert.Equal(t, expectedManifest, manifest)
			w.Header().Add(httpx.HeaderXDeployID, fmt.Sprintf("%d", TestDeployID))
			httpx.RespondAccepted(manifest[2:], w) // NOTE: slice off first and second - assume we already have them
		} else if r.URL.Path == "/api/void/snakes/deploy/42/activate" {
			httpx.RespondOk(&share.DeployResult{
				DeployID: TestDeployID,
				Slug:     TestDeploySlug,
				URL:      TestDeployURL,
				Manifest: expectedManifest,
			}, w)
		} else if strings.HasPrefix(r.URL.Path, "/api/void/snakes/deploy/42") {
			httpx.RespondOk("ok", w)
		} else {
			httpx.RespondBadRequest(fmt.Sprintf("unexpected %s", r.URL.Path), w)
		}
	}))

	api, err := api.NewClient(mockServer.URL, TestToken)
	assert.NoError(t, err)

	uploads := make([]string, 0)
	onUpload := func(deployID int64, path string) {
		uploads = append(uploads, path)
	}

	result, err := share.Deploy(&share.DeployCommand{
		API:      api,
		Org:      TestOrg,
		Game:     TestGame,
		Path:     mockDir.Dir,
		OnUpload: onUpload,
	})

	assert.NoError(t, err)
	assert.Equal(t, TestDeployID, result.DeployID)
	assert.Equal(t, TestDeploySlug, result.Slug)
	assert.Equal(t, TestDeployURL, result.URL)
	assert.Equal(t, expectedManifest, result.Manifest)

	assert.Equal(t, []string{
		"path/to/third.txt",
	}, uploads)
}

//-------------------------------------------------------------------------------------------------

func TestDeployIgnoresDisallowedFiles(t *testing.T) {

	mockDir := mock.TempDir(t)
	mockDir.AddTextFile(t, ".ssh", "secret")
	mockDir.AddTextFile(t, ".env", "secret")
	mockDir.AddTextFile(t, ".git", "secret")
	mockDir.AddTextFile(t, FirstPath, FirstContent)

	expectedManifest := []share.DeployEntry{
		{
			Path:          FirstPath,
			Blake3:        crypto.Blake3(FirstContent),
			ContentLength: len(FirstContent),
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.RequestMethodEqual(t, http.MethodPost, r)
		assert.RequestHeaderEqual(t, "Bearer personal-access-token", httpx.HeaderAuthorization, r)
		if r.URL.Path == "/api/void/snakes/deploy" {
			manifest := assert.RequestJSON[[]share.DeployEntry](t, r)
			assert.Equal(t, expectedManifest, manifest)
			w.Header().Add(httpx.HeaderXDeployID, fmt.Sprintf("%d", TestDeployID))
			httpx.RespondAccepted(manifest, w)
		} else if r.URL.Path == "/api/void/snakes/deploy/42/activate" {
			httpx.RespondOk(&share.DeployResult{
				DeployID: TestDeployID,
				Slug:     TestDeploySlug,
				URL:      TestDeployURL,
				Manifest: expectedManifest,
			}, w)
		} else if strings.HasPrefix(r.URL.Path, "/api/void/snakes/deploy/42") {
			httpx.RespondOk("ok", w)
		} else {
			httpx.RespondBadRequest(fmt.Sprintf("unexpected %s", r.URL.Path), w)
		}
	}))

	api, err := api.NewClient(mockServer.URL, TestToken)
	assert.NoError(t, err)

	uploads := make([]string, 0)
	onUpload := func(deployID int64, path string) {
		uploads = append(uploads, path)
	}

	result, err := share.Deploy(&share.DeployCommand{
		API:      api,
		Org:      TestOrg,
		Game:     TestGame,
		Path:     mockDir.Dir,
		OnUpload: onUpload,
	})

	assert.NoError(t, err)
	assert.Equal(t, TestDeployID, result.DeployID)
	assert.Equal(t, TestDeploySlug, result.Slug)
	assert.Equal(t, TestDeployURL, result.URL)
	assert.Equal(t, expectedManifest, result.Manifest)

	assert.Equal(t, []string{
		"path/to/first.txt",
	}, uploads)
}

//-------------------------------------------------------------------------------------------------

func TestDeployFailed(t *testing.T) {

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.RequestMethodEqual(t, http.MethodPost, r)
		assert.RequestPathEqual(t, "/api/void/snakes/deploy", r)
		assert.RequestHeaderEqual(t, "Bearer personal-access-token", httpx.HeaderAuthorization, r)
		assert.RequestBodyEqual(t, "[]", r)
		httpx.RespondBadRequest("uh oh, manifest was empty", w)
	}))
	api, err := api.NewClient(mockServer.URL, TestToken)
	assert.Nil(t, err)
	assert.NotNil(t, api)

	tmp := mock.TempDir(t)

	_, err = share.Deploy(&share.DeployCommand{
		API:  api,
		Org:  TestOrg,
		Game: TestGame,
		Path: tmp.Dir,
	})

	assert.NotNil(t, err)
	assert.Error(t, "unexpected status code 400: uh oh, manifest was empty", err)
}

//-------------------------------------------------------------------------------------------------
