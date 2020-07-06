package host

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/filanov/bm-inventory/models"
)

type stopPodmanCmd struct {
	baseCmd
}

func NewStopPodmanCmd(log logrus.FieldLogger) *stopPodmanCmd {
	return &stopPodmanCmd{
		baseCmd: baseCmd{log: log},
	}
}

func (h *stopPodmanCmd) GetStep(ctx context.Context, host *models.Host) (*models.Step, error) {
	step := &models.Step{
		StepType: models.StepTypeStopPodman,
		Command:  "/usr/bin/podman",
		Args: []string{
			"kill", "--all",
		},
	}
	return step, nil
}
