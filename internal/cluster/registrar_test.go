package cluster

import (
	context "context"

	"github.com/filanov/bm-inventory/models"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("registrar", func() {
	var (
		ctx             = context.Background()
		registerManager RegistrationAPI
		db              *gorm.DB
		currentState    = clusterStatusInsufficient
		id              strfmt.UUID
		updateReply     *UpdateReply
		updateErr       error
		cluster         models.Cluster
		host            models.Host
	)

	BeforeEach(func() {
		db = prepareDB()
		//state = &Manager{insufficient: NewInsufficientState(getTestLog(), db)}
		registerManager = NewRegistrar(getTestLog(), db)

		id = strfmt.UUID(uuid.New().String())
		cluster = models.Cluster{
			Base: models.Base{
				ID: &id,
			},
			Status: swag.String(currentState),
		}

		//register cluster
		updateReply, updateErr = registerManager.RegisterCluster(ctx, &cluster)
		Expect(updateErr).Should(BeNil())
		Expect(updateReply.State).Should(Equal(clusterStatusInsufficient))
		c := geCluster(*cluster.ID, db)
		Expect(swag.StringValue(c.Status)).Should(Equal(clusterStatusInsufficient))
	})

	Context("register cluster", func() {
		It("register a registered cluster", func() {
			updateReply, updateErr = registerManager.RegisterCluster(ctx, &cluster)
			Expect(updateErr).Should(HaveOccurred())
			Expect(updateReply.State).Should(Equal(clusterStatusInsufficient))
		})
	})

	Context("deregister", func() {
		It("unregister a registered cluster", func() {
			updateReply, updateErr = registerManager.DeregisterCluster(ctx, &cluster)
			Expect(updateErr).Should(BeNil())
			Expect(updateReply.State).Should(Equal("unregistered"))

			Expect(db.First(&cluster, "id = ?", cluster.ID).Error).Should(HaveOccurred())
			Expect(db.First(&host, "cluster_id = ?", cluster.ID).Error).Should(HaveOccurred())

		})
	})

	AfterEach(func() {

		db.Close()
		updateReply = nil
		updateErr = nil
	})
})
