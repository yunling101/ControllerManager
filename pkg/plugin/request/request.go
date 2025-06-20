package request

import (
	"io"
	"net/http"
)

type client struct {
	method             string
	url                string
	body               io.Reader
	header             map[string]string
	readBody           bool
	username, password string
}

type response struct {
	Error error
	Body  []byte
}

func New(url string) *client {
	return &client{method: "GET", url: url, readBody: false}
}

func (c *client) SetMethod(method string) *client {
	c.method = method
	return c
}

func (c *client) SetBody(body io.Reader) *client {
	c.body = body
	return c
}

func (c *client) SetHeader(header map[string]string) *client {
	c.header = header
	return c
}

func (c *client) SetReadBody(readBody bool) *client {
	c.readBody = readBody
	return c
}

func (c *client) SetBasicAuth(username, password string) *client {
	c.username = username
	c.password = password
	return c
}

func (c *client) Do() (r response) {
	request, err := http.NewRequest(c.method, c.url, c.body)
	if err != nil {
		r.Error = err
		return
	}

	if len(c.header) != 0 {
		for k, v := range c.header {
			request.Header.Add(k, v)
		}
	}

	if c.username != "" && c.password != "" {
		request.SetBasicAuth(c.username, c.password)
	}

	reqBody, err := http.DefaultClient.Do(request)
	if err != nil {
		r.Error = err
		return
	}
	defer func() {
		_ = reqBody.Body.Close()
	}()

	if c.readBody {
		r.Body, r.Error = io.ReadAll(reqBody.Body)
	}
	return
}
