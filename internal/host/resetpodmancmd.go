package host

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/filanov/bm-inventory/models"
)

type resetPodmanCmd struct {
	baseCmd
}

func NewResetPodmanCmd(log logrus.FieldLogger) *resetPodmanCmd {
	return &resetPodmanCmd{
		baseCmd: baseCmd{log: log},
	}
}

func (h *resetPodmanCmd) GetStep(ctx context.Context, host *models.Host) (*models.Step, error) {
	step := &models.Step{
		StepType: models.StepTypeResetPodman,
		Command:  "podman",
		Args: []string{
			"container", "rm", "--all", "-f",
		},
	}
	return step, nil
}
