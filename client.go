package tushare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"time"
)

var (
	Endpoint = "http://api.tushare.pro"
)

type DataApi struct {
	token  string
	client *http.Client
}

func NewApi(token string) *DataApi {
	return &DataApi{
		token: token,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func NewWithClient(token string, client *http.Client) *DataApi {
	return &DataApi{
		token:  token,
		client: client,
	}
}

func (api *DataApi) Query(apiName string, params map[string]string, fields string) (*ProResponse, error) {
	reqBody := map[string]interface{}{
		"api_name": apiName,
		"token":    api.token,
		"params":   params,
		"fields":   fields,
	}
	req, err := api.buildRequest(http.MethodPost, Endpoint, reqBody)
	if err != nil {
		return nil, err
	}
	proRsp, err := api.doRequest(req)
	if err != nil {
		return proRsp, err
	}
	return proRsp, nil
}

func (api *DataApi) buildRequest(method string, url string, body interface{}) (*http.Request, error) {
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (api *DataApi) doRequest(req *http.Request) (*ProResponse, error) {

	req.Header.Set("Content-Type", "application/json")

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read request
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	//Handle network error
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s", string(body))
	}

	// Check mime type of response
	mimeType, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}
	if mimeType != "application/json" {
		return nil, fmt.Errorf("Could not execute request (%s)", fmt.Sprintf("Response Content-Type is '%s', but should be 'application/json'.", mimeType))
	}

	// Parse Request
	var dataResp ProResponse

	if err = json.Unmarshal(body, &dataResp); err != nil {
		return nil, err
	}
	if dataResp.Code != 0 {
		return &dataResp, ApiError{dataResp.Code, dataResp.Msg}
	}
	return &dataResp, nil
}
