package subsystem

import (
	"context"
	"strings"

	"github.com/filanov/bm-inventory/internal/bminventory"

	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"

	"github.com/filanov/bm-inventory/client/inventory"
	"github.com/filanov/bm-inventory/models"
	"github.com/go-openapi/swag"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Host tests", func() {
	ctx := context.Background()
	var cluster *inventory.RegisterClusterCreated
	var clusterID strfmt.UUID

	AfterEach(func() {
		clearDB()
	})

	BeforeEach(func() {
		var err error
		cluster, err = bmclient.Inventory.RegisterCluster(ctx, &inventory.RegisterClusterParams{
			NewClusterParams: &models.ClusterCreateParams{
				Name: swag.String("test cluster"),
			},
		})
		Expect(err).NotTo(HaveOccurred())
	})

	JustBeforeEach(func() {
		clusterID = *cluster.GetPayload().ID
	})

	It("host CRUD", func() {
		host := registerHost(clusterID)
		host = getHost(clusterID, *host.ID)
		Expect(*host.Status).Should(Equal("discovering"))

		list, err := bmclient.Inventory.ListHosts(ctx, &inventory.ListHostsParams{ClusterID: clusterID})
		Expect(err).NotTo(HaveOccurred())
		Expect(len(list.GetPayload())).Should(Equal(1))

		_, err = bmclient.Inventory.DeregisterHost(ctx, &inventory.DeregisterHostParams{
			ClusterID: clusterID,
			HostID:    *host.ID,
		})
		Expect(err).NotTo(HaveOccurred())
		list, err = bmclient.Inventory.ListHosts(ctx, &inventory.ListHostsParams{ClusterID: clusterID})
		Expect(err).NotTo(HaveOccurred())
		Expect(len(list.GetPayload())).Should(Equal(0))

		_, err = bmclient.Inventory.GetHost(ctx, &inventory.GetHostParams{
			ClusterID: clusterID,
			HostID:    *host.ID,
		})
		Expect(err).Should(HaveOccurred())
	})

	It("next step", func() {
		host := registerHost(clusterID)
		steps := getNextSteps(clusterID, *host.ID)
		_, ok := getStepByIdPrefix(steps, bminventory.HardwareInfo)
		Expect(ok).Should(Equal(true))
	})

	It("disable enable", func() {
		host := registerHost(clusterID)
		_, err := bmclient.Inventory.DisableHost(ctx, &inventory.DisableHostParams{
			ClusterID: clusterID,
			HostID:    *host.ID,
		})
		Expect(err).NotTo(HaveOccurred())
		host = getHost(clusterID, *host.ID)
		Expect(*host.Status).Should(Equal("disabled"))
		Expect(len(getNextSteps(clusterID, *host.ID))).Should(Equal(0))

		_, err = bmclient.Inventory.EnableHost(ctx, &inventory.EnableHostParams{
			ClusterID: clusterID,
			HostID:    *host.ID,
		})
		Expect(err).NotTo(HaveOccurred())
		host = getHost(clusterID, *host.ID)
		Expect(*host.Status).Should(Equal("discovering"))
		Expect(len(getNextSteps(clusterID, *host.ID))).ShouldNot(Equal(0))
	})

	It("debug", func() {
		host1 := registerHost(clusterID)
		host2 := registerHost(clusterID)
		// set debug to host1
		_, err := bmclient.Inventory.SetDebugStep(ctx, &inventory.SetDebugStepParams{
			ClusterID: clusterID,
			HostID:    *host1.HostID,
			Step:      &models.DebugStep{Command: swag.String("echo hello")},
		})
		Expect(err).NotTo(HaveOccurred())

		var step *models.Step
		var ok bool
		// debug should be only for host1
		_, ok = getStepByIdPrefix(getNextSteps(clusterID, *host2.ID), string(models.StepTypeExecute))
		Expect(ok).Should(Equal(false))

		step, ok = getStepByIdPrefix(getNextSteps(clusterID, *host1.ID), string(models.StepTypeExecute))
		Expect(ok).Should(Equal(true))
		Expect(step.Command).Should(Equal("bash"))
		Expect(step.Args).Should(Equal([]string{"-c", "echo hello"}))

		// debug executed only once
		_, ok = getStepByIdPrefix(getNextSteps(clusterID, *host1.ID), string(models.StepTypeExecute))
		Expect(ok).Should(Equal(false))

		_, err = bmclient.Inventory.PostStepReply(ctx, &inventory.PostStepReplyParams{
			ClusterID: clusterID,
			HostID:    *host1.ID,
			Reply: &models.StepReply{
				ExitCode: 0,
				Output:   "hello",
				StepID:   step.StepID,
			},
		})
		Expect(err).NotTo(HaveOccurred())
	})

	It("register same host id", func() {
		hostID := strToUUID(uuid.New().String())
		// register to cluster1
		_, err := bmclient.Inventory.RegisterHost(context.Background(), &inventory.RegisterHostParams{
			ClusterID: clusterID,
			NewHostParams: &models.HostCreateParams{
				HostID: hostID,
			},
		})
		Expect(err).NotTo(HaveOccurred())

		cluster2, err := bmclient.Inventory.RegisterCluster(ctx, &inventory.RegisterClusterParams{
			NewClusterParams: &models.ClusterCreateParams{
				Name: swag.String("another cluster"),
			},
		})
		Expect(err).NotTo(HaveOccurred())

		// register to cluster2
		_, err = bmclient.Inventory.RegisterHost(ctx, &inventory.RegisterHostParams{
			ClusterID: *cluster2.GetPayload().ID,
			NewHostParams: &models.HostCreateParams{
				HostID: hostID,
			},
		})
		Expect(err).NotTo(HaveOccurred())

		// successfully get from both clusters
		_ = getHost(clusterID, *hostID)
		_ = getHost(*cluster2.GetPayload().ID, *hostID)

		_, err = bmclient.Inventory.DeregisterHost(ctx, &inventory.DeregisterHostParams{
			ClusterID: clusterID,
			HostID:    *hostID,
		})
		Expect(err).NotTo(HaveOccurred())
		h := getHost(*cluster2.GetPayload().ID, *hostID)

		// register again to cluster 2 and expect it to be in discovery status
		Expect(db.Model(h).Update("status", "known").Error).NotTo(HaveOccurred())
		h = getHost(*cluster2.GetPayload().ID, *hostID)
		Expect(swag.StringValue(h.Status)).Should(Equal("known"))
		_, err = bmclient.Inventory.RegisterHost(ctx, &inventory.RegisterHostParams{
			ClusterID: *cluster2.GetPayload().ID,
			NewHostParams: &models.HostCreateParams{
				HostID: hostID,
			},
		})
		Expect(err).NotTo(HaveOccurred())
		h = getHost(*cluster2.GetPayload().ID, *hostID)
		Expect(swag.StringValue(h.Status)).Should(Equal("discovering"))
	})
})

func getStepByIdPrefix(steps models.Steps, idPrefix string) (*models.Step, bool) {
	for _, step := range steps {
		if strings.HasPrefix(step.StepID, idPrefix) {
			return step, true
		}
	}
	return nil, false
}

func getNextSteps(clusterID, hostID strfmt.UUID) models.Steps {
	steps, err := bmclient.Inventory.GetNextSteps(context.Background(), &inventory.GetNextStepsParams{
		ClusterID: clusterID,
		HostID:    hostID,
	})
	Expect(err).NotTo(HaveOccurred())
	return steps.GetPayload()
}
