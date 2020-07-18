package host

import (
	"context"

	"github.com/filanov/bm-inventory/models"
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("reset", func() {
	ctx := context.Background()
	var host models.Host
	var rstCmd *resetInstallationCmd
	var id, clusterId strfmt.UUID
	var stepReply *models.Step
	var stepErr error

	BeforeEach(func() {
		rstCmd = NewResetInstallationCmd(getTestLog())

		id = strfmt.UUID(uuid.New().String())
		clusterId = strfmt.UUID(uuid.New().String())
		host = getTestHost(id, clusterId, HostStatusResetting)
		Expect(db.Create(&host).Error).ShouldNot(HaveOccurred())
	})

	It("get_step", func() {
		stepReply, stepErr = rstCmd.GetStep(ctx, &host)
		Expect(stepReply.StepType).To(Equal(models.StepTypeResetInstallation))
		Expect(stepErr).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		// cleanup
		stepReply = nil
		stepErr = nil
	})
})
