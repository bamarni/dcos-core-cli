package networking

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dcos/dcos-cli/pkg/httpclient"
)

// Client is a networking client for DC/OS.
type Client struct {
	http *httpclient.Client
}

// NewClient creates a new networking client.
func NewClient(baseClient *httpclient.Client) *Client {
	return &Client{
		http: baseClient,
	}
}

// Nodes returns the nodes in a cluster.
func (c *Client) Nodes() ([]Node, error) {
	resp, err := c.http.Get("/net/v1/nodes")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200:
		var nodes []Node
		err = json.NewDecoder(resp.Body).Decode(&nodes)
		if err != nil {
			return nil, err
		}
		return nodes, nil
	default:
		return nil, httpResponseToError(resp)
	}
}

func httpResponseToError(resp *http.Response) error {
	if resp.StatusCode < 400 {
		return fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}
	return &httpclient.HTTPError{
		Response: resp,
	}
}
