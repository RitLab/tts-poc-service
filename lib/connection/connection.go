package connection

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"tts-poc-service/lib/baselogger"
)

type httpConnection struct {
	client *http.Client
	logger *baselogger.Logger
}

type HttpConnectionInterface interface {
	Get(uri string, queries map[string]string) (status int, response []byte, err error)
	GetWithContext(ctx context.Context, uri string, headers map[string]string, queries map[string]string) (status int, response []byte, err error)
	Post(uri string, request []byte) (status int, response []byte, err error)
	PostWithContext(ctx context.Context, uri string, headers map[string]string, request []byte) (status int, response []byte, err error)
}

func NewHttpConnection(client *http.Client, logger *baselogger.Logger) HttpConnectionInterface {
	return &httpConnection{
		client: client,
		logger: logger,
	}
}

func (h *httpConnection) SubmitRequest(req *http.Request) (status int, response []byte, err error) {
	var res *http.Response

	if res, err = h.client.Do(req); err != nil {
		h.logger.Error(fmt.Printf("Exception caught %s\n", err.Error()))
		return
	}

	defer res.Body.Close()

	status = res.StatusCode

	response, err = io.ReadAll(res.Body)

	return
}

func (h *httpConnection) Get(uri string, queries map[string]string) (status int, response []byte, err error) {
	status, response, err = h.GetWithContext(context.Background(), uri, nil, queries)
	return
}

func (h *httpConnection) GetWithContext(ctx context.Context, uri string, headers map[string]string, queries map[string]string) (status int, response []byte, err error) {
	if uri == "" {
		err = fmt.Errorf("Unable to resolve uri")
		h.logger.Error(err)
		return
	}

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return
	}

	if queries != nil {
		q := req.URL.Query()
		for k, v := range queries {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	status, response, err = h.SubmitRequest(req)

	return
}

func (h *httpConnection) Post(uri string, request []byte) (status int, response []byte, err error) {
	status, response, err = h.PostWithContext(context.Background(), uri, nil, request)

	return
}

func (h *httpConnection) PostWithContext(ctx context.Context, uri string, headers map[string]string, request []byte) (status int, response []byte, err error) {
	if uri == "" {
		err = fmt.Errorf("Unable to resolve uri")
		return
	}

	var req *http.Request

	req, err = http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(request))
	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/json")

	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	status, response, err = h.SubmitRequest(req)

	return
}
