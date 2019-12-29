package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const DoorIDEmpty string = ""

type Client struct {
	HostURI string
	APIKey  string
	client  http.Client
}

func NewClient() *Client {
	c := new(Client)
	c.client.Timeout = 5 * time.Second

	return c
}

func (p *Client) Call(method, path, id string, dest interface{}) error {
	req := p.makeRequest(method, path, id)
	return p.sendRequest(req, dest)
}

func (p *Client) makeRequest(method, path, id string) *http.Request {
	url := p.HostURI
	if id != DoorIDEmpty {
		url = fmt.Sprintf("%s%s%s", p.HostURI, path, id)
	} else {
		url = fmt.Sprintf("%s%s", p.HostURI, path)
	}

	req, _ := http.NewRequest(method, url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", p.APIKey))

	return req
}

func (p *Client) sendRequest(r *http.Request, target interface{}) error {
	res, err := p.client.Do(r)
	if err != nil {
		return err
	}

	dec := json.NewDecoder(res.Body)
	if err := dec.Decode(&target); err != nil {
		return err
	}

	return nil
}
