package prometheus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type Client struct {
	BaseURL string
}

type QueryResponse struct {
	Status string `json:"status"`
	Data   struct {
		Result []struct {
			Value []interface{} `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

func (c *Client) Query(query string) (*QueryResponse, error) {
	endpoint := fmt.Sprintf(
		"%s/api/v1/query?query=%s",
		c.BaseURL,
		url.QueryEscape(query),
	)

	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result QueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *QueryResponse) FirstValue() string {
	if len(r.Data.Result) == 0 {
		return ""
	}

	if len(r.Data.Result[0].Value) < 2 {
		return ""
	}

	value, ok := r.Data.Result[0].Value[1].(string)
	if !ok {
		return ""
	}

	return value
}

func (r *QueryResponse) FirstValueFloat() float64 {
	value := r.FirstValue()
	if value == "" {
		return 0
	}

	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0
	}

	return f
}
