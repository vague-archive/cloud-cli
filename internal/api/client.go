package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/vaguevoid/cloud-cli/internal/lib/httpx"
)

type Client struct {
	Endpoint *url.URL
	Token    string
	httpc    *http.Client
}

//-------------------------------------------------------------------------------------------------

func NewClient(server string, token string) (*Client, error) {
	url, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	return &Client{
		Endpoint: url,
		Token:    token,
		httpc:    &http.Client{},
	}, nil
}

//-------------------------------------------------------------------------------------------------

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set(httpx.HeaderAuthorization, fmt.Sprintf("Bearer %s", c.Token))
	return c.httpc.Do(req)
}

//-------------------------------------------------------------------------------------------------

func (c *Client) Get(route string) (*http.Response, error) {
	url := c.URL(route)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

//-------------------------------------------------------------------------------------------------

func (c *Client) Post(route string, content io.Reader) (*http.Response, error) {
	url := c.URL(route)
	req, err := http.NewRequest(http.MethodPost, url, content)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

//-------------------------------------------------------------------------------------------------

func (c *Client) PostJSON(route string, content any) (*http.Response, error) {
	url := c.URL(route)

	data, err := json.Marshal(content)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set(httpx.HeaderContentType, httpx.ContentTypeJSON)

	return c.Do(req)
}

//-------------------------------------------------------------------------------------------------

func (c *Client) PostFILE(route string, filepath string) (*http.Response, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	url := c.URL(route)
	req, err := http.NewRequest(http.MethodPost, url, f)
	if err != nil {
		return nil, err
	}
	req.ContentLength = fi.Size()
	req.Header.Set(httpx.HeaderContentType, httpx.ContentTypeBytes)

	return c.Do(req)
}

//-------------------------------------------------------------------------------------------------

func (c *Client) URL(route string) string {
	url := *c.Endpoint
	url.Path = path.Join("api", route)
	return url.String()
}

func (c *Client) Route(parts ...any) string {
	strs := make([]string, len(parts))
	for i, part := range parts {
		strs[i] = fmt.Sprint(part)
	}
	return path.Join(strs...)
}

//-------------------------------------------------------------------------------------------------
