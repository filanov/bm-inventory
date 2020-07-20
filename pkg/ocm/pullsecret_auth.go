package ocm

import (
	"context"
)

type OCMAuthentication interface {
	AuthenticatePullSecret(ctx context.Context, pullSecret string) (allowed bool, err error)
}

type authentication service

var _ OCMAuthentication = &authentication{}

func (a authentication) AuthenticatePullSecret(ctx context.Context, pullSecret string) (allowed bool, err error) {
	return true, nil
}
