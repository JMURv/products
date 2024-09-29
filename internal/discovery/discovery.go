package discovery

import (
	"bytes"
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"go.uber.org/zap"
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
		zap.L().Debug("Error marshalling request", zap.Error(err))
		return err
	}

	post, err := http.Post(fmt.Sprintf("%v/register", d.url), "application/json", bytes.NewBuffer(req))
	if err != nil {
		zap.L().Debug("Error registering service", zap.Error(err))
		return err
	}

	if err != nil || post.StatusCode != http.StatusCreated {
		zap.L().Debug("Error registering service", zap.Error(err))
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
		zap.L().Debug("Error marshalling request", zap.Error(err))
		return err
	}

	post, err := http.Post(fmt.Sprintf("%v/deregister", d.url), "application/json", bytes.NewBuffer(req))
	if err != nil {
		zap.L().Debug("Error deregistering service", zap.Error(err))
		return err
	}

	if err != nil || post.StatusCode != http.StatusOK {
		zap.L().Debug("Error deregistering service", zap.Error(err))
		return err
	}

	return nil
}

func (d *Discovery) FindServiceByName(ctx context.Context, name string) (string, error) {
	req, err := json.Marshal(map[string]string{
		"name": name,
	})
	if err != nil {
		zap.L().Debug("Error marshalling request", zap.Error(err))
		return "", err
	}

	post, err := http.Post(fmt.Sprintf("%v/find", d.url), "application/json", bytes.NewBuffer(req))
	if err != nil {
		zap.L().Debug("Error service", zap.Error(err))
		return "", err
	}

	if err != nil || post.StatusCode != http.StatusOK {
		zap.L().Debug("Error service", zap.Error(err))
		return "", err
	}

	res := struct {
		Address string `json:"address"`
	}{}
	if err := json.NewDecoder(post.Body).Decode(&res); err != nil {
		zap.L().Debug("Error service", zap.Error(err))
		return "", err
	}

	return res.Address, nil
}
