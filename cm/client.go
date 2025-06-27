package cm

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"nikeron/cmbxip/config"
	"strings"
)

type Client struct {
	httpClient *http.Client
	apiURL     string
	auth       string
}

func NewClient(apiURL, auth string) *Client {
	return &Client{
		httpClient: &http.Client{},
		apiURL:     apiURL,
		auth:       auth}
}

func (c *Client) SetAuth(auth string) *Client {
	c.auth = auth
	return c
}

func (c *Client) SetAPIURL(url string) *Client {
	c.apiURL = url
	return c
}

func (c *Client) GetURI(uri string) (*http.Response, error) {
	config.DebugLogger().Printf("cm.GetURI: %s", uri)
	if parsedURI, err := url.Parse(uri); err == nil {
		if strings.Contains(parsedURI.Scheme, "http") {
			req, err := http.NewRequest(http.MethodGet, uri, nil)
			if err != nil {
				return nil, err
			}
			config.DebugLogger().Printf("cm.GetURI: http scheme found, req: %+v", req)
			req.Header.Set("Authorization", "Basic "+c.auth)
			return c.httpClient.Do(req)
		}
	}
	req, err := c.newAPIRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	config.DebugLogger().Printf("cm.GetURI: req: %+v", req)
	return c.httpClient.Do(req)
}

func (c *Client) ComponentVersions() (*ComponentVersionsResponse, error) {
	statusCode, data, err := c.newAPIRequestBody(http.MethodGet, "component-versions", nil)
	if err != nil {
		return nil, err
	}
	if statusCode != 200 {
		return nil, c.errorFromStatusCodeAndData(statusCode, data)
	}
	jsonDecoder := json.NewDecoder(bytes.NewReader(data))
	jsonDecoder.DisallowUnknownFields()
	retValue := &ComponentVersionsResponse{}
	if err = jsonDecoder.Decode(retValue); err != nil {
		return nil, err
	}
	return retValue, nil
}

func (c *Client) FromID(id string) (*Document, error) {
	statusCode, data, err := c.newAPIRequestBody(http.MethodGet, "ids/"+id, nil)
	if err != nil {
		return nil, err
	}
	if statusCode != 200 {
		return nil, c.errorFromStatusCodeAndData(statusCode, data)
	}
	retValue := &Document{}
	if err = json.Unmarshal(data, retValue); err != nil {
		return nil, err
	}
	return retValue, nil
}

func (c *Client) ExecutionHierarchy(id string) (*ExecutionResponse, error) {
	statusCode, data, err := c.newAPIRequestBody(http.MethodGet, "execution/hierarchy/"+id, nil)
	if err != nil {
		return nil, err
	}
	if statusCode != 200 {
		return nil, c.errorFromStatusCodeAndData(statusCode, data)
	}
	parsedResponse := &ExecutionResponse{}
	if err = json.Unmarshal(data, parsedResponse); err != nil {
		return nil, err
	}
	return parsedResponse, nil
}

func (c *Client) Executors(id string) ([]string, error) {
	executionHierarchy, err := c.ExecutionHierarchy(id)
	if err != nil {
		return nil, err
	}
	retMapLikeSet := map[string]interface{}{}
	for _, v := range executionHierarchy.Entry {
		for k := range c.executorsFromEntry(v) {
			retMapLikeSet[k] = nil
		}
	}
	retValue := []string{}
	for k := range retMapLikeSet {
		retValue = append(retValue, k)
	}
	return retValue, nil
}

func (c *Client) executorsFromEntry(entry ExecutionEntry) map[string]interface{} {
	retMapLikeSet := map[string]interface{}{}
	for _, executor := range entry.Value.Executor {
		retMapLikeSet[executor.Executor.FullName] = nil
	}
	for _, innerEntry := range entry.Value.Execution {
		for k := range c.executorsFromEntry(innerEntry) {
			retMapLikeSet[k] = nil
		}
	}
	return retMapLikeSet
}
