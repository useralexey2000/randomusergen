package apiclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"randomusergen/domain"
)

var url = "https://randomuser.me/api/1.4/"

type Requester interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client interface {
	Get(ctx context.Context, query string) (*domain.Response, error)
}

type client struct {
	requester Requester
}

func New(r Requester) Client {
	return &client{
		requester: r,
	}
}

func (c *client) Get(ctx context.Context, query string) (*domain.Response, error) {
	url := fmt.Sprintf("%v?%v", url, query)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.requester.Do(req)
	if err != nil {
		return nil, err
	}

	var resApi domain.Response
	err = json.NewDecoder(res.Body).Decode(&resApi)
	if err != nil {
		return nil, err
	}

	return &resApi, nil
}
