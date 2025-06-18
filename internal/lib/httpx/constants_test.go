package httpx_test

import (
	"testing"

	"github.com/vaguevoid/cloud-cli/internal/lib/httpx"
	"github.com/vaguevoid/cloud-cli/internal/test/assert"
)

//-------------------------------------------------------------------------------------------------

func TestParams(t *testing.T) {
	assert.Equal(t, "cli", httpx.ParamCLI)
	assert.Equal(t, "jwt", httpx.ParamJWT)
	assert.Equal(t, "origin", httpx.ParamOrigin)
}

//-------------------------------------------------------------------------------------------------

func TestHttpHeaders(t *testing.T) {
	assert.Equal(t, "Accept", httpx.HeaderAccept)
	assert.Equal(t, "Access-Control-Allow-Headers", httpx.HeaderAccessControlAllowHeaders)
	assert.Equal(t, "Access-Control-Allow-Methods", httpx.HeaderAccessControlAllowMethods)
	assert.Equal(t, "Access-Control-Allow-Origin", httpx.HeaderAccessControlAllowOrigin)
	assert.Equal(t, "Authorization", httpx.HeaderAuthorization)
	assert.Equal(t, "X-CSRF-Token", httpx.HeaderCSRFToken)
	assert.Equal(t, "Cache-Control", httpx.HeaderCacheControl)
	assert.Equal(t, "Content-Disposition", httpx.HeaderContentDisposition)
	assert.Equal(t, "Content-Encoding", httpx.HeaderContentEncoding)
	assert.Equal(t, "Content-Language", httpx.HeaderContentLanguage)
	assert.Equal(t, "Content-Length", httpx.HeaderContentLength)
	assert.Equal(t, "Content-Location", httpx.HeaderContentLocation)
	assert.Equal(t, "Content-MD5", httpx.HeaderContentMD5)
	assert.Equal(t, "Content-Range", httpx.HeaderContentRange)
	assert.Equal(t, "Content-Type", httpx.HeaderContentType)
	assert.Equal(t, "Cookie", httpx.HeaderCookie)
	assert.Equal(t, "Cross-Origin-Embedder-Policy", httpx.HeaderCrossOriginEmbedderPolicy)
	assert.Equal(t, "Cross-Origin-Opener-Policy", httpx.HeaderCrossOriginOpenerPolicy)
	assert.Equal(t, "Cross-Origin-Resource-Policy", httpx.HeaderCrossOriginResourcePolicy)
	assert.Equal(t, "X-Command", httpx.HeaderCustomCommand)
	assert.Equal(t, "ETag", httpx.HeaderETag)
	assert.Equal(t, "Expires", httpx.HeaderExpires)
	assert.Equal(t, "HX-Redirect", httpx.HeaderHxRedirect)
	assert.Equal(t, "HX-Refresh", httpx.HeaderHxRefresh)
	assert.Equal(t, "HX-Request", httpx.HeaderHxRequest)
	assert.Equal(t, "HX-Retarget", httpx.HeaderHxRetarget)
	assert.Equal(t, "If-Modified-Since", httpx.HeaderIfModifiedSince)
	assert.Equal(t, "Last-Modified", httpx.HeaderLastModified)
	assert.Equal(t, "Location", httpx.HeaderLocation)
	assert.Equal(t, "User-Agent", httpx.HeaderUserAgent)
	assert.Equal(t, "X-Deploy-ID", httpx.HeaderXDeployID)
	assert.Equal(t, "X-Deploy-Label", httpx.HeaderXDeployLabel)
	assert.Equal(t, "X-Deploy-Password", httpx.HeaderXDeployPassword)
	assert.Equal(t, "X-Deploy-Pinned", httpx.HeaderXDeployPinned)
	assert.Equal(t, "X-Forwarded-For", httpx.HeaderXForwardedFor)
	assert.Equal(t, "X-Forwarded-Host", httpx.HeaderXForwardedHost)
	assert.Equal(t, "X-Forwarded-Prefix", httpx.HeaderXForwardedPrefix)
	assert.Equal(t, "X-Forwarded-Proto", httpx.HeaderXForwardedProto)
	assert.Equal(t, "X-Frame-Options", httpx.HeaderXFrameOptions)
}

//-------------------------------------------------------------------------------------------------

func TestHttpContentTypes(t *testing.T) {
	assert.Equal(t, "application/octet-stream", httpx.ContentTypeBytes)
	assert.Equal(t, "text/css", httpx.ContentTypeCSS)
	assert.Equal(t, "application/x-www-form-urlencoded", httpx.ContentTypeForm)
	assert.Equal(t, "application/gzip", httpx.ContentTypeGzip)
	assert.Equal(t, "text/html", httpx.ContentTypeHTML)
	assert.Equal(t, "text/javascript", httpx.ContentTypeJavascript)
	assert.Equal(t, "application/json", httpx.ContentTypeJSON)
	assert.Equal(t, "text/markdown", httpx.ContentTypeMarkdown)
	assert.Equal(t, "application/pdf", httpx.ContentTypePdf)
	assert.Equal(t, "image/png", httpx.ContentTypePng)
	assert.Equal(t, "text/plain", httpx.ContentTypeText)
	assert.Equal(t, "text/plain; charset=utf-8", httpx.ContentTypeTextUtf8)
	assert.Equal(t, "application/wasm", httpx.ContentTypeWasm)
	assert.Equal(t, "text/xml", httpx.ContentTypeXML)
}

//-------------------------------------------------------------------------------------------------
