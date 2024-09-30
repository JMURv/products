package discovery

import (
	"bytes"
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"net/http"
)

type Discovery struct {
	url  string
	name string
	addr string
}

func New(url, name, addr string) *Discovery {
	return &Discovery{
		url:  url,
		name: name,
		addr: addr,
	}
}

func (d *Discovery) Register() error {
	req, err := json.Marshal(map[string]string{
		"name":    d.name,
		"address": d.addr,
	})
	if err != nil {
		return err
	}

	post, err := http.Post(fmt.Sprintf("%v/register", d.url), "application/json", bytes.NewBuffer(req))
	if err != nil {
		return err
	}

	if err != nil || post.StatusCode != http.StatusCreated {
		return err
	}

	return nil
}

func (d *Discovery) Deregister() error {
	req, err := json.Marshal(map[string]string{
		"name":    d.name,
		"address": d.addr,
	})
	if err != nil {
		return err
	}

	post, err := http.Post(fmt.Sprintf("%v/deregister", d.url), "application/json", bytes.NewBuffer(req))
	if err != nil {
		return err
	}

	if err != nil || post.StatusCode != http.StatusOK {
		return err
	}

	return nil
}

func (d *Discovery) FindServiceByName(ctx context.Context, name string) (string, error) {
	req, err := json.Marshal(map[string]string{
		"name": name,
	})
	if err != nil {
		return "", err
	}

	post, err := http.Post(fmt.Sprintf("%v/find", d.url), "application/json", bytes.NewBuffer(req))
	if err != nil {
		return "", err
	}

	if err != nil || post.StatusCode != http.StatusOK {
		return "", err
	}

	res := struct {
		Address string `json:"address"`
	}{}
	if err := json.NewDecoder(post.Body).Decode(&res); err != nil {
		return "", err
	}

	return res.Address, nil
}
