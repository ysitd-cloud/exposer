package client

import (
	"net/http"
)

func New(c http.Client, serviceName string) *client {
	return &client{
		Client:      c,
		ServiceName: serviceName,
	}
}
