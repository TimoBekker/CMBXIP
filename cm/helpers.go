package cm

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
)

func (c *Client) newAPIRequest(method, apiuri string, body io.Reader) (*http.Request, error) {
	if len(apiuri) > 0 && apiuri[0] == '/' {
		apiuri = apiuri[1:]
	}
	req, err := http.NewRequest(method, c.apiURL+"/"+apiuri, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Basic "+c.auth)
	return req, nil
}

// newAPIRequestBody returns status code + body data
func (c *Client) newAPIRequestBody(method, apiuri string, body io.Reader) (int, []byte, error) {
	req, err := c.newAPIRequest(method, apiuri, body)
	if err != nil {
		return -1, nil, err
	}
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
	if json.Unmarshal(data, errResponse) == nil && errResponse.Message != "" {
		return errors.New(errResponse.Message)
	}
	return fmt.Errorf("response status code: %d", statusCode)
}

var regexpIDFromAddress = regexp.MustCompile(`\\[0-9]([0-9A-Fa-f]*\:[0-9A-Fa-f]*)`)

func IDFromAddress(address string) string {
	address, _ = url.PathUnescape(address)
	address, _ = url.PathUnescape(address)
	submatch := regexpIDFromAddress.FindStringSubmatch(address)
	if len(submatch) == 2 {
		return submatch[1]
	}
	return ""
}
