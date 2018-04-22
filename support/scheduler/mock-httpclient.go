//
// Copyright (c) 2018 Tencent
//
// SPDX-License-Identifier: Apache-2.0
//

package scheduler

type MockHttpClient struct {
}

func NewMockHttpClient() HttpClient {
	return MockHttpClient{}
}

func (hc MockHttpClient) Post(url string, contentType string, contentLength int, timeoutStr string) ([]byte, int, error) {
	return nil, 0, nil
}
