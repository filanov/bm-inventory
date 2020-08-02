package auth

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/filanov/bm-inventory/internal/common"
	"github.com/filanov/bm-inventory/models"
	"github.com/filanov/bm-inventory/pkg/ocm"
	"github.com/filanov/bm-inventory/restapi"
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// #nosec
const ()

type mockOCMAuthorization struct {
	ocm.OCMAuthorization
}

var accessReviewMock func(ctx context.Context, username, action, resourceType string) (allowed bool, err error)

func (m *mockOCMAuthorization) AccessReview(ctx context.Context, username, action, resourceType string) (allowed bool, err error) {
	return accessReviewMock(ctx, username, action, resourceType)
}

func TestValidator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "authorizer_test")
}

var _ = Describe("Authorizer", func() {
	var (
		ctx         = context.Background()
		dbName      = "authorizer"
		db          *gorm.DB
		allowedUser bool
		clustersAPI = "/api/assisted-install/v1/clusters/"
	)

	BeforeEach(func() {
		db = common.PrepareTestDB(dbName)
		authHandler.db = db
		accessReviewMock = func(ctx context.Context, username, action, resourceType string) (allowed bool, err error) {
			return allowedUser, nil
		}
	})

	AfterEach(func() {
		common.DeleteTestDB(db, dbName)
	})

	Context("Unauthorized User", func() {
		It("Empty context", func() {
			ctx = context.WithValue(ctx, restapi.AuthKey, nil)
			err := Authorizer(getRequestWithContext(ctx, ""))

			Expect(err).Should(BeNil())
		})
		It("Empty payload", func() {
			ctx = context.WithValue(ctx, restapi.AuthKey, &AuthPayload{})
			err := Authorizer(getRequestWithContext(ctx, ""))

			Expect(err).Should(BeNil())
		})
		It("User unallowed to access installer", func() {
			mockOCMClient()
			allowedUser = false

			payload := &AuthPayload{}
			payload.Username = "unallowed@user"
			ctx = context.WithValue(ctx, restapi.AuthKey, payload)

			err := Authorizer(getRequestWithContext(ctx, ""))
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).Should(Equal("method is not allowed"))
		})
		It("User unallowed to access cluster", func() {
			mockOCMClient()
			allowedUser = true

			payload := &AuthPayload{}
			payload.Username = "unallowed@user"
			ctx = context.WithValue(ctx, restapi.AuthKey, payload)

			req := getRequestWithContext(ctx, clustersAPI+uuid.New().String())
			err := Authorizer(req)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).Should(Equal("method is not allowed"))
		})
	})

	Context("Authorized User", func() {
		It("User allowed to access owned cluster", func() {
			mockOCMClient()
			allowedUser = true

			payload := &AuthPayload{}
			payload.Username = "allowed@user"
			ctx = context.WithValue(ctx, restapi.AuthKey, payload)

			clusterID := strfmt.UUID(uuid.New().String())
			req := getRequestWithContext(ctx, clustersAPI+clusterID.String())

			err := db.Create(&common.Cluster{Cluster: models.Cluster{
				ID:       &clusterID,
				UserName: payload.Username,
			}}).Error
			Expect(err).ShouldNot(HaveOccurred())

			err = Authorizer(req)
			Expect(err).ToNot(HaveOccurred())
		})
		It("User allowed to access non cluster context API", func() {
			mockOCMClient()
			allowedUser = true

			payload := &AuthPayload{}
			payload.Username = "allowed@user"
			ctx = context.WithValue(ctx, restapi.AuthKey, payload)

			req := getRequestWithContext(ctx, "/api/assisted-install/v1/events/")

			err := Authorizer(req)
			Expect(err).ToNot(HaveOccurred())
		})
		It("Admin allowed all endpoints", func() {
			mockOCMClient()
			allowedUser = true

			payload := &AuthPayload{}
			payload.Username = "admin@user"
			payload.IsAdmin = true
			ctx = context.WithValue(ctx, restapi.AuthKey, payload)

			clusterID := strfmt.UUID(uuid.New().String())
			req := getRequestWithContext(ctx, clustersAPI+clusterID.String())

			err := db.Create(&common.Cluster{Cluster: models.Cluster{
				ID:       &clusterID,
				UserName: "nonadmin@user",
			}}).Error
			Expect(err).ShouldNot(HaveOccurred())

			err = Authorizer(req)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})

func getRequestWithContext(ctx context.Context, urlPath string) *http.Request {
	req := &http.Request{}
	req.URL = &url.URL{}
	req.URL.Path = urlPath
	return req.WithContext(ctx)
}

func mockOCMClient() {
	authHandler.client = ocm.Client{}
	authHandler.client.Authorization = &mockOCMAuthorization{}
}

func mockAccessReview(allowed bool) {

}
