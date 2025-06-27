package bx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) newAPIRequest(method, apiuri string, body io.Reader) (*http.Request, error) {
	if len(apiuri) > 0 && apiuri[0] == '/' {
		apiuri = apiuri[1:]
	}
	req, err := http.NewRequest(method, c.apiURL+"/"+apiuri, body)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// newAPIRequestBody returns status code + body data
func (c *Client) newAPIRequestBody(method, apiuri string, body io.Reader) (int, []byte, error) {
	req, err := c.newAPIRequest(method, apiuri, body)
	if err != nil {
		return -1, nil, err
	}
	return c.doReqBodyResp(req)
}

func (c *Client) doReqBodyResp(req *http.Request) (int, []byte, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return -1, nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1, nil, err
	}
	return resp.StatusCode, data, nil
}

func (c *Client) errorFromStatusCodeAndData(statusCode int, data []byte) error {
	errResponse := &ErrorResponse{}
	if json.Unmarshal(data, errResponse) == nil && errResponse.Description != "" {
		return errors.New(errResponse.Description)
	}
	return fmt.Errorf("response status code: %d", statusCode)
}
