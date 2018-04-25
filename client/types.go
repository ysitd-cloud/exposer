package client

import (
	"context"
)

type Client interface {
	List(ctx context.Context) (serviceMap map[string]*Service, err error)
	Get(ctx context.Context, hostname string) (service *Service, err error)
	Connect(ctx context.Context, hostname string, service *Service) (success bool, err error)
	Migrate(ctx context.Context, hostname string, service *Service) (success bool, err error)
	Disconnect(ctx context.Context, hostname string) (success bool, err error)
}
