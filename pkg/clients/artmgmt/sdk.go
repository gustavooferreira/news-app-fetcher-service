package artmgmt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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

func (c *Client) AddArticles(articles Articles) (err error) {
	requestBody, err := json.Marshal(articles)
	if err != nil {
		return err
	}

	rawURL := c.baseURL + "/articles:batch"

	req, err := http.NewRequest("POST", rawURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		return fmt.Errorf("status code returned different than 200 - status code [%d]", resp.StatusCode)
	}

	return nil
}
