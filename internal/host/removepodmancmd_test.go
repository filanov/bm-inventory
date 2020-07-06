package host

import (
	"context"

	"github.com/filanov/bm-inventory/models"
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("remove-podman", func() {
	ctx := context.Background()
	var host models.Host
	var db *gorm.DB
	var rmCmd *removePodmanCmd
	var id, clusterId strfmt.UUID
	var stepReply *models.Step
	var stepErr error

	BeforeEach(func() {
		db = prepareDB()
		rmCmd = NewRemovePodmanCmd(getTestLog())

		id = strfmt.UUID(uuid.New().String())
		clusterId = strfmt.UUID(uuid.New().String())
		host = getTestHost(id, clusterId, HostStatusResetting)
		Expect(db.Create(&host).Error).ShouldNot(HaveOccurred())
	})

	It("get_step", func() {
		stepReply, stepErr = rmCmd.GetStep(ctx, &host)
		Expect(stepReply.StepType).To(Equal(models.StepTypeRemovePodman))
		Expect(stepErr).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		// cleanup
		db.Close()
		stepReply = nil
		stepErr = nil
	})
})
