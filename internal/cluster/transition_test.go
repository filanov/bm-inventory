package cluster

import (
	"context"

	"github.com/filanov/bm-inventory/internal/common"

	"github.com/go-openapi/swag"

	. "github.com/onsi/gomega"

	"github.com/filanov/bm-inventory/models"

	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("CancelInstallation", func() {
	var (
		ctx       = context.Background()
		capi      API
		db        *gorm.DB
		clusterId strfmt.UUID
	)

	BeforeEach(func() {
		db = prepareDB()
		capi = NewManager(getTestLog(), db, nil)
		clusterId = strfmt.UUID(uuid.New().String())
	})

	It("cancel_installation", func() {
		c := common.Cluster{
			Cluster: models.Cluster{ID: &clusterId, Status: swag.String(clusterStatusInstalling)},
		}
		Expect(db.Create(&c).Error).ShouldNot(HaveOccurred())
		Expect(capi.CancelInstallation(ctx, &c, "", db)).ShouldNot(HaveOccurred())
		Expect(swag.StringValue(c.Status)).Should(Equal(clusterStatusError))
	})

	AfterEach(func() {
		db.Close()
	})
})
