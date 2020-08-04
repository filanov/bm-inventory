package generator

import (
	"context"

	"github.com/filanov/bm-inventory/internal/common"
	"github.com/filanov/bm-inventory/internal/events"
)

type ISOGenerator interface {
	GenerateISO(ctx context.Context, cluster common.Cluster, jobName string, imageName string, ignitionConfig string, eventsHandler events.Handler) error
}

type InstallConfigGenerator interface {
	GenerateInstallConfig(ctx context.Context, cluster common.Cluster, cfg []byte) error
	AbortInstallConfig(ctx context.Context, cluster common.Cluster) error
}

type ISOInstallConfigGenerator interface {
	ISOGenerator
	InstallConfigGenerator
}
