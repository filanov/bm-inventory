package host

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/filanov/bm-inventory/models"
)

type removePodmanCmd struct {
	baseCmd
}

func NewRemovePodmanCmd(log logrus.FieldLogger) *removePodmanCmd {
	return &removePodmanCmd{
		baseCmd: baseCmd{log: log},
	}
}

func (h *removePodmanCmd) GetStep(ctx context.Context, host *models.Host) (*models.Step, error) {
	step := &models.Step{
		StepType: models.StepTypeRemovePodman,
		Command:  "/usr/bin/podman",
		Args: []string{
			"rm", "--all", "-f",
		},
	}
	return step, nil
}
