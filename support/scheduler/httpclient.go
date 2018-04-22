//
// Copyright (c) 2018 Tencent
//
// SPDX-License-Identifier: Apache-2.0
//

package scheduler

import (
	"io/ioutil"
	"net/http"
	"time"
)

type HttpClient interface {
	Post(url string, contentType string, contentLength int, timeoutStr string) ([]byte, int, error)
}

type DefaultHttpClient struct {
}

func NewDefaultHttpClient() HttpClient {
	return DefaultHttpClient{}
}

func (dhc DefaultHttpClient) Post(url string, contentType string, contentLength int, timeoutStr string) ([]byte, int, error) {
	req, err := http.NewRequest(http.MethodPost, url, nil)

	req.Header.Set(ContentTypeKey, contentType)
	req.Header.Set(ContentLengthKey, string(contentLength))

	if err != nil {
		loggingClient.Error("create new request occurs error : " + err.Error())
	}

	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		loggingClient.Error("parse timeout duration error : " + err.Error())
	}

	client := &http.Client{
		Timeout: timeout,
	}

	return makeRequestAndGetResponse(client, req)
}

func makeRequestAndGetResponse(client *http.Client, req *http.Request) ([]byte, int, error) {
	resp, err := client.Do(req)

	if err != nil {
		return []byte{}, 500, err
	}

	defer resp.Body.Close()
	resp.Close = true

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, 500, err
	}

	return bodyBytes, resp.StatusCode, nil
}
