package bx

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Client struct {
	httpClient *http.Client
	apiURL     string
}

func NewClient(apiURL string) *Client {
	return &Client{
		httpClient: &http.Client{},
		apiURL:     apiURL}
}

func (c *Client) SetAPIURL(url string) *Client {
	c.apiURL = url
	return c
}

func (c *Client) SearchUser(query string, isEmployee bool) ([]*User, error) {
	reqURL := "user.search.json?FIND=" + url.QueryEscape(query)
	if isEmployee {
		reqURL += "&USER_TYPE=employee"
	}
	statusCode, data, err := c.newAPIRequestBody(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	if statusCode != 200 {
		return nil, c.errorFromStatusCodeAndData(statusCode, data)
	}
	retValue := &UserSearchResponse{}
	if err = json.Unmarshal(data, retValue); err != nil {
		return nil, err
	}
	return retValue.Result, nil
}

func (c *Client) CurrentUser() (*User, error) {
	statusCode, data, err := c.newAPIRequestBody(http.MethodGet, "user.current.json", nil)
	if err != nil {
		return nil, err
	}
	if statusCode != 200 {
		return nil, c.errorFromStatusCodeAndData(statusCode, data)
	}
	retValue := &UserCurrentResponse{}
	if err = json.Unmarshal(data, retValue); err != nil {
		return nil, err
	}
	return retValue.Result, nil
}

func (c *Client) AddTask(addTaskRequest *AddTaskRequest) (*AddTaskResponse, error) {
	addTaskJSONPayload, err := json.Marshal(addTaskRequest)
	if err != nil {
		return nil, err
	}
	req, err := c.newAPIRequest(http.MethodPost, "tasks.task.add.json", bytes.NewReader(addTaskJSONPayload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	statusCode, data, err := c.doReqBodyResp(req)
	if err != nil {
		return nil, err
	}
	if statusCode != 200 {
		return nil, c.errorFromStatusCodeAndData(statusCode, data)
	}
	retValue := &AddTaskResponse{}
	if err = json.Unmarshal(data, retValue); err != nil {
		return nil, err
	}
	return retValue, nil
}

func (c *Client) TaskAttachFile(taskID, fileID string) (*TaskAttachFileResponse, error) {
	statusCode, data, err := c.newAPIRequestBody(http.MethodGet, "tasks.task.files.attach.json?fileId="+fileID+"&taskId="+taskID, nil)
	if err != nil {
		return nil, err
	}
	if statusCode != 200 {
		return nil, c.errorFromStatusCodeAndData(statusCode, data)
	}
	retValue := &TaskAttachFileResponse{}
	if err = json.Unmarshal(data, retValue); err != nil {
		return nil, err
	}
	return retValue, nil
}

func (c *Client) DiskStorageList() (*DiskStorageListResponse, error) {
	statusCode, data, err := c.newAPIRequestBody(http.MethodGet, "disk.storage.getlist.json", nil)
	if err != nil {
		return nil, err
	}
	if statusCode != 200 {
		return nil, c.errorFromStatusCodeAndData(statusCode, data)
	}
	retValue := &DiskStorageListResponse{}
	if err = json.Unmarshal(data, retValue); err != nil {
		return nil, err
	}
	return retValue, nil
}

func (c *Client) DiskStorageListOfUser(id string) (*DiskStorageListResponse, error) {
	req, err := c.newAPIRequest(http.MethodPost, "disk.storage.getlist.json", bytes.NewReader(
		[]byte(fmt.Sprintf("{\"filter\":{\"ENTITY_TYPE\":\"user\",\"ENTITY_ID\":\"%s\"}}", id)),
	))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	statusCode, data, err := c.doReqBodyResp(req)
	if err != nil {
		return nil, err
	}
	if statusCode != 200 {
		return nil, c.errorFromStatusCodeAndData(statusCode, data)
	}
	retValue := &DiskStorageListResponse{}
	if err = json.Unmarshal(data, retValue); err != nil {
		return nil, err
	}
	return retValue, nil
}

func (c *Client) DiskStorageChildren(id string) (*DiskStorageChildrenResponse, error) {
	statusCode, data, err := c.newAPIRequestBody(http.MethodGet, "disk.storage.getchildren.json?id="+id, nil)
	if err != nil {
		return nil, err
	}
	if statusCode != 200 {
		return nil, c.errorFromStatusCodeAndData(statusCode, data)
	}
	retValue := &DiskStorageChildrenResponse{}
	if err = json.Unmarshal(data, retValue); err != nil {
		return nil, err
	}
	return retValue, nil
}

func (c *Client) DiskFolderUploadFile(folderID, fileName string, fileContent []byte, generateUniqueName bool) (*DiskFolderUploadResponse, error) {
	jsonPayload, err := json.Marshal(&DiskFolderUploadRequest{
		ID: folderID,
		Data: struct {
			Name string `json:"NAME"`
		}{Name: fileName},
		FileContent:        []string{fileName, base64.StdEncoding.EncodeToString(fileContent)},
		GenerateUniqueName: generateUniqueName,
	})
	if err != nil {
		return nil, err
	}
	req, err := c.newAPIRequest(http.MethodPost, "disk.folder.uploadfile.json", bytes.NewReader(jsonPayload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	statusCode, data, err := c.doReqBodyResp(req)
	if err != nil {
		return nil, err
	}
	if statusCode != 200 {
		return nil, c.errorFromStatusCodeAndData(statusCode, data)
	}
	retValue := &DiskFolderUploadResponse{}
	if err = json.Unmarshal(data, retValue); err != nil {
		return nil, err
	}
	return retValue, nil
}
