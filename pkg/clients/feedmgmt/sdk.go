package feedmgmt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

// NOTE: This client (SDK) needs improvement, like better error reporting in logs and what not
// And ideally this would be on a library of its own with all the RPCs supported.

type Client struct {
	httpClient *http.Client
	baseURL    string
}

func NewClient(host string, port int, httpClient *http.Client) *Client {
	c := &Client{httpClient: httpClient, baseURL: fmt.Sprintf("http://%s:%d/api/v1", host, port)}
	return c
}

func (c *Client) GetFeeds(provider string, category string, enabled bool) (feeds Feeds, err error) {

	rawURL := c.baseURL + "/feeds"

	v := url.Values{}
	v.Set("enabled", strconv.FormatBool(enabled))
	if provider != "" {
		v.Set("provider", provider)
	}
	if category != "" {
		v.Set("category", category)
	}
	queryParams := v.Encode()

	if len(v) != 0 {
		rawURL += "?" + queryParams
	}

	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code returned different than 200 - status code [%d]", resp.StatusCode)
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var responseBodyData Feeds

	err = json.Unmarshal(responseBody, &responseBodyData)
	if err != nil {
		return nil, err
	}

	return responseBodyData, nil
}
