package hyper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func AcceptLanguage(spec string) func(*http.Request) {
	return func(r *http.Request) {
		r.Header.Set(HeaderAcceptLanguage, spec)
	}
}

func Accept(typ string) func(*http.Request) {
	return func(r *http.Request) {
		r.Header.Set(HeaderAccept, typ)
	}
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{},
	}
}

type Client struct {
	httpClient *http.Client
}

func (c *Client) Fetch(url string, opts ...func(*http.Request)) (Item, error) {
	opts = append([]func(*http.Request){Accept(ContentTypeHyperItem)}, opts...)
	_, data, err := c.FetchRaw(url, opts...)
	res := Item{}
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(data, &res)
	if err != nil {
		return res, fmt.Errorf("decode: %v", err)
	}
	return res, nil
}

func (c *Client) FetchRaw(url string, opts ...func(*http.Request)) (http.Header, []byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("create: %v", err)
	}
	for _, opt := range opts {
		opt(req)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("do: %v", err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	return resp.Header, data, err
}

func (c *Client) Submit(a Action, args Arguments) (*http.Response, error) {
	if args == nil {
		args = Arguments{}
	}
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(args)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(a.Method, a.Href, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set(HeaderContentType, a.Encoding)
	return c.httpClient.Do(req)
}
