package cluster

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/jinzhu/gorm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/filanov/bm-inventory/models"
)

var _ = Describe("insufficient_state", func() {
	var (
		ctx          = context.Background()
		state        API
		db           *gorm.DB
		currentState = clusterStatusInsufficient
		id           *strfmt.UUID
		updateReply  *UpdateReply
		updateErr    error
		cluster      models.Cluster
	)

	BeforeEach(func() {
		db = prepareDB()
		state = &Manager{insufficient: NewInsufficientState(getTestLog(), db)}
		registerManager := NewRegistrar(getTestLog(), db)

		cluster = models.Cluster{
			Status: swag.String(currentState),
		}

		c, replyErr := registerManager.RegisterCluster(ctx, &cluster)
		id = c.ID
		Expect(replyErr).Should(BeNil())
		cluster = geCluster(*cluster.ID, db)

		Expect(swag.StringValue(cluster.Status)).Should(Equal(clusterStatusInsufficient))
		Expect(swag.StringValue(c.Status)).Should(Equal(clusterStatusInsufficient))
	})

	Context("refresh_state", func() {
		It("not answering requirement to be ready", func() {
			updateReply, updateErr = state.RefreshStatus(ctx, &cluster, db)
			Expect(updateErr).Should(BeNil())
			Expect(updateReply.State).Should(Equal(clusterStatusInsufficient))
			c := geCluster(*cluster.ID, db)
			Expect(swag.StringValue(c.Status)).Should(Equal(clusterStatusInsufficient))
		})

		It("answering requirement to be ready", func() {
			addInstallationRequirements(*id, db)
			updateReply, updateErr = state.RefreshStatus(ctx, &cluster, db)
			Expect(updateErr).Should(BeNil())
			Expect(updateReply.State).Should(Equal(clusterStatusReady))
			c := geCluster(*cluster.ID, db)
			Expect(swag.StringValue(c.Status)).Should(Equal(clusterStatusReady))

		})
	})

	AfterEach(func() {
		db.Close()
		updateReply = nil
		updateErr = nil
	})
})
