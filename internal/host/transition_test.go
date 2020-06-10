package host

import (
	"context"

	"github.com/pkg/errors"

	"github.com/golang/mock/gomock"

	"github.com/go-openapi/swag"

	. "github.com/onsi/gomega"

	"github.com/filanov/bm-inventory/models"

	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("RegisterHost", func() {
	var (
		ctx                     = context.Background()
		hapi                    API
		db                      *gorm.DB
		ctrl                    *gomock.Controller
		mockExternalValidations *MockExternalValidations
		hostId, clusterId       strfmt.UUID
	)

	BeforeEach(func() {
		db = prepareDB()
		ctrl = gomock.NewController(GinkgoT())
		mockExternalValidations = NewMockExternalValidations(ctrl)
		hapi = NewManager(getTestLog(), db, nil, nil, mockExternalValidations)
		hostId = strfmt.UUID(uuid.New().String())
		clusterId = strfmt.UUID(uuid.New().String())
		Expect(db.Create(&models.Cluster{ID: &clusterId}).Error).ShouldNot(HaveOccurred())
	})

	conditionSuccess := func() {
		mockExternalValidations.EXPECT().AcceptRegistration(gomock.Any()).Return(nil).Times(1)
	}

	conditionFailed := func() {
		mockExternalValidations.EXPECT().AcceptRegistration(gomock.Any()).
			Return(errors.Errorf("error")).Times(1)
	}

	It("register_new", func() {
		conditionSuccess()
		Expect(hapi.RegisterHost(ctx, &models.Host{ID: &hostId, ClusterID: clusterId})).ShouldNot(HaveOccurred())
		h := getHost(hostId, clusterId, db)
		Expect(swag.StringValue(h.Status)).Should(Equal(HostStatusDiscovering))
	})

	Context("register during installation put host in error", func() {
		tests := []struct {
			name     string
			srcState string
		}{
			{
				name:     "discovering",
				srcState: HostStatusInstalling,
			},
			{
				name:     "insufficient",
				srcState: HostStatusInstallingInProgress,
			},
		}

		AfterEach(func() {
			h := getHost(hostId, clusterId, db)
			Expect(swag.StringValue(h.Status)).Should(Equal(HostStatusError))
			Expect(h.Role).Should(Equal(RoleMaster))
			Expect(h.HardwareInfo).Should(Equal(defaultHwInfo))
		})

		for i := range tests {
			t := tests[i]

			It(t.name, func() {
				Expect(db.Create(&models.Host{
					ID:           &hostId,
					ClusterID:    clusterId,
					Role:         RoleMaster,
					HardwareInfo: defaultHwInfo,
					Status:       swag.String(t.srcState),
				}).Error).ShouldNot(HaveOccurred())

				Expect(hapi.RegisterHost(ctx, &models.Host{
					ID:        &hostId,
					ClusterID: clusterId,
					Status:    swag.String(t.srcState),
				})).ShouldNot(HaveOccurred())
			})
		}
	})

	Context("host already exists register success", func() {
		tests := []struct {
			name     string
			srcState string
		}{
			{
				name:     "discovering",
				srcState: HostStatusDiscovering,
			},
			{
				name:     "insufficient",
				srcState: HostStatusInsufficient,
			},
			{
				name:     "disconnected",
				srcState: HostStatusDisconnected,
			},
			{
				name:     "known",
				srcState: HostStatusKnown,
			},
		}

		AfterEach(func() {
			h := getHost(hostId, clusterId, db)
			Expect(swag.StringValue(h.Status)).Should(Equal(HostStatusDiscovering))
			Expect(h.Role).Should(Equal(""))
			Expect(h.HardwareInfo).Should(Equal(""))
		})

		for i := range tests {
			t := tests[i]
			It(t.name, func() {
				conditionSuccess()
				Expect(db.Create(&models.Host{
					ID:           &hostId,
					ClusterID:    clusterId,
					Role:         RoleMaster,
					HardwareInfo: defaultHwInfo,
					Status:       swag.String(t.srcState),
				}).Error).ShouldNot(HaveOccurred())

				Expect(hapi.RegisterHost(ctx, &models.Host{
					ID:        &hostId,
					ClusterID: clusterId,
					Status:    swag.String(t.srcState),
				})).ShouldNot(HaveOccurred())
			})
		}
	})

	Context("host already exist registration fail", func() {
		tests := []struct {
			name        string
			srcState    string
			targetState string
			condition   func()
		}{
			{
				name:     "disabled",
				srcState: HostStatusDisabled,
			},
			{
				name:     "error",
				srcState: HostStatusError,
			},
			{
				name:     "installed",
				srcState: HostStatusInstalled,
			},
			{
				name:      "discovering",
				srcState:  HostStatusDiscovering,
				condition: conditionFailed,
			},
			{
				name:      "insufficient",
				srcState:  HostStatusInsufficient,
				condition: conditionFailed,
			},
			{
				name:      "disconnected",
				srcState:  HostStatusDisconnected,
				condition: conditionFailed,
			},
			{
				name:      "known",
				srcState:  HostStatusKnown,
				condition: conditionFailed,
			},
		}

		for i := range tests {
			t := tests[i]
			It(t.name, func() {
				if t.condition != nil {
					t.condition()
				}
				Expect(db.Create(&models.Host{
					ID:           &hostId,
					ClusterID:    clusterId,
					Role:         RoleMaster,
					HardwareInfo: defaultHwInfo,
					Status:       swag.String(t.srcState),
				}).Error).ShouldNot(HaveOccurred())

				Expect(hapi.RegisterHost(ctx, &models.Host{
					ID:        &hostId,
					ClusterID: clusterId,
					Status:    swag.String(t.srcState),
				})).Should(HaveOccurred())

				h := getHost(hostId, clusterId, db)
				Expect(swag.StringValue(h.Status)).Should(Equal(t.srcState))
				Expect(h.Role).Should(Equal(RoleMaster))
				Expect(h.HardwareInfo).Should(Equal(defaultHwInfo))
			})
		}
	})

	AfterEach(func() {
		ctrl.Finish()
		db.Close()
	})
})

var _ = Describe("HostInstallationFailed", func() {
	var (
		ctx               = context.Background()
		hapi              API
		db                *gorm.DB
		hostId, clusterId strfmt.UUID
		host              models.Host
	)

	BeforeEach(func() {
		db = prepareDB()
		hapi = NewManager(getTestLog(), db, nil, nil, nil)
		hostId = strfmt.UUID(uuid.New().String())
		clusterId = strfmt.UUID(uuid.New().String())
		host = getTestHost(hostId, clusterId, "")
		host.Status = swag.String(HostStatusInstalling)
		Expect(db.Create(&host).Error).ShouldNot(HaveOccurred())
	})

	It("handle_installation_error", func() {
		Expect(hapi.HandleInstallationFailure(ctx, &host)).ShouldNot(HaveOccurred())
		h := getHost(hostId, clusterId, db)
		Expect(swag.StringValue(h.Status)).Should(Equal(HostStatusError))
		Expect(swag.StringValue(h.StatusInfo)).Should(Equal("installation command failed"))
	})

	AfterEach(func() {
		db.Close()
	})
})
