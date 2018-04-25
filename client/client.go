package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type client struct {
	Client      http.Client `inject:""`
	ServiceName string      `inject:"exposer service"`
}

func (c *client) List(ctx context.Context) (serviceMap map[string]*Service, err error) {
	url := fmt.Sprintf("http://%s/", c.ServiceName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req = req.WithContext(ctx)

	res, err := c.Client.Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()

	content, err := ioutil.ReadAll(res.Body)

	connections := make([]connection, 0)

	if err := json.Unmarshal(content, &connections); err != nil {
		return
	}

	serviceMap = make(map[string]*Service, 0)

	for _, conn := range connections {
		serviceMap[conn.Hostname] = &conn.Service
	}

	return
}

func (c *client) Get(ctx context.Context, hostname string) (service *Service, err error) {
	url := fmt.Sprintf("http://%s/%s", c.ServiceName, hostname)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req = req.WithContext(ctx)

	res, err := c.Client.Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()

	content, err := ioutil.ReadAll(res.Body)

	var instance connection

	if err := json.Unmarshal(content, &instance); err != nil {
		return nil, err
	}

	service = &instance.Service

	return
}

func (c *client) Connect(ctx context.Context, hostname string, service *Service) (success bool, err error) {
	url := fmt.Sprintf("http://%s/%s", c.ServiceName, hostname)

	body, err := json.Marshal(service)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return
	}
	req = req.WithContext(ctx)

	res, err := c.Client.Do(req)
	if err != nil {
		return
	}

	return res.StatusCode == http.StatusCreated, nil
}

func (c *client) Migrate(ctx context.Context, hostname string, service *Service) (success bool, err error) {
	url := fmt.Sprintf("http://%s/%s", c.ServiceName, hostname)

	body, err := json.Marshal(service)
	if err != nil {
		return
	}

	req, err := http.NewRequest("PUT", url, bytes.NewReader(body))
	if err != nil {
		return
	}
	req = req.WithContext(ctx)

	res, err := c.Client.Do(req)
	if err != nil {
		return
	}

	return res.StatusCode == http.StatusOK, nil
}

func (c *client) Disconnect(ctx context.Context, hostname string) (success bool, err error) {
	url := fmt.Sprintf("http://%s/%s", c.ServiceName, hostname)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return
	}
	req = req.WithContext(ctx)

	res, err := c.Client.Do(req)
	if err != nil {
		return
	}

	return res.StatusCode == http.StatusOK, nil
}
