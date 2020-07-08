package requestid

import (
	"context"
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/filanov/bm-inventory/restapi"
	"github.com/filanov/bm-inventory/restapi/operations/installer"
	"github.com/sirupsen/logrus"
	"gotest.tools/assert"

	"github.com/stretchr/testify/mock"
)

type mockTransport struct {
	mock.Mock
}

func (m *mockTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	m.Called(r)
	return nil, nil
}

func TestTransport(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		prepare func(t *testing.T, tr *mockTransport) *http.Request
	}{
		{
			name: "happy flow",
			prepare: func(t *testing.T, tr *mockTransport) *http.Request {
				const requestID = "1234"

				match := mock.MatchedBy(func(req *http.Request) bool {
					return req.Header.Get(headerKey) == requestID
				})

				tr.On("RoundTrip", match).Return(nil, nil).Once()

				ctx := context.WithValue(context.Background(), ctxKey, requestID)
				req := httptest.NewRequest(http.MethodGet, "http://example.org", nil)
				req = req.WithContext(ctx)
				return req
			},
		},
		{
			name: "no request id in context",
			prepare: func(t *testing.T, tr *mockTransport) *http.Request {
				match := mock.MatchedBy(func(req *http.Request) bool {
					return req.Header.Get(headerKey) == ""
				})

				tr.On("RoundTrip", match).Return(nil, nil).Once()

				req := httptest.NewRequest(http.MethodGet, "http://example.org", nil)
				return req
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tr mockTransport
			defer tr.AssertExpectations(t)
			req := tt.prepare(t, &tr)
			_, _ = Transport(&tr).RoundTrip(req)
		})
	}
}

func TestAuth(t *testing.T) {
	url1 := installer.ListClustersURL{}
	url2 := installer.GetClusterURL{ClusterID: "2faecba1-4903-4e2f-a994-fb58bd770066"}

	listClustersUrl := url1.String()
	getClusterUrl := url2.String()
	agentKey := "X-Secret-Key"
	userKey := "X-User-Key"
	agentKeyValue := "SecretKey"
	userKeyValue := "userKey"

	t.Parallel()
	tests := []struct {
		name               string
		tokenKey           string
		expectedTokenValue string
		url                string
	}{
		{
			name:               "user auth",
			tokenKey:           userKey,
			expectedTokenValue: userKeyValue,
			url:                listClustersUrl,
		},
		{
			name:               "agent auth",
			tokenKey:           agentKey,
			expectedTokenValue: agentKeyValue,
			url:                getClusterUrl,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			type p struct {
				user string
			}
			authAgentAuth := func(token string) (interface{}, error) {
				assert.Equal(t, tt.expectedTokenValue, token)
				assert.Equal(t, tt.tokenKey, agentKey)
				return "user2", nil
			}

			// AuthUserAuth for basic authentication
			authUserAuth := func(token string) (interface{}, error) {
				assert.Equal(t, tt.expectedTokenValue, token)
				assert.Equal(t, tt.tokenKey, userKey)
				return "user1", nil
			}
			h, _ := restapi.Handler(restapi.Config{
				AuthAgentAuth:     authAgentAuth,
				AuthUserAuth:      authUserAuth,
				InstallerAPI:      fakeInventory{},
				EventsAPI:         nil,
				Logger:            logrus.Printf,
				VersionsAPI:       nil,
				ManagedDomainsAPI: nil,
				InnerMiddleware:   nil,
			})

			// create a mock request to use
			req := httptest.NewRequest("GET", tt.url, nil)
			req.Header.Set(tt.tokenKey, tt.expectedTokenValue)

			// call the handler using a mock response recorder (we'll not use that anyway)
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			fmt.Println(rec)
			assert.Equal(t, rec.Code, 200)
		})
	}
}

type fakeInventory struct{}

func (f fakeInventory) CancelInstallation(ctx context.Context, params installer.CancelInstallationParams) middleware.Responder {
	panic("implement me")
}

func (f fakeInventory) DeregisterCluster(ctx context.Context, params installer.DeregisterClusterParams) middleware.Responder {
	panic("implement me")
}

func (f fakeInventory) DeregisterHost(ctx context.Context, params installer.DeregisterHostParams) middleware.Responder {
	panic("implement me")
}

func (f fakeInventory) DisableHost(ctx context.Context, params installer.DisableHostParams) middleware.Responder {
	panic("implement me")
}

func (f fakeInventory) DownloadClusterFiles(ctx context.Context, params installer.DownloadClusterFilesParams) middleware.Responder {
	panic("implement me")
}

func (f fakeInventory) DownloadClusterISO(ctx context.Context, params installer.DownloadClusterISOParams) middleware.Responder {
	panic("implement me")
}

func (f fakeInventory) DownloadClusterKubeconfig(ctx context.Context, params installer.DownloadClusterKubeconfigParams) middleware.Responder {
	panic("implement me")
}

func (f fakeInventory) EnableHost(ctx context.Context, params installer.EnableHostParams) middleware.Responder {
	panic("implement me")
}

func (f fakeInventory) GenerateClusterISO(ctx context.Context, params installer.GenerateClusterISOParams) middleware.Responder {
	panic("implement me")
}

func (f fakeInventory) GetCluster(ctx context.Context, params installer.GetClusterParams) middleware.Responder {
	return installer.NewGetClusterOK()
}

func (f fakeInventory) GetCredentials(ctx context.Context, params installer.GetCredentialsParams) middleware.Responder {
	panic("implement me")
}

func (f fakeInventory) GetFreeAddresses(ctx context.Context, params installer.GetFreeAddressesParams) middleware.Responder {
	panic("implement me")
}

func (f fakeInventory) GetHost(ctx context.Context, params installer.GetHostParams) middleware.Responder {
	panic("implement me")
}

func (f fakeInventory) GetNextSteps(ctx context.Context, params installer.GetNextStepsParams) middleware.Responder {
	panic("implement me")
}

func (f fakeInventory) InstallCluster(ctx context.Context, params installer.InstallClusterParams) middleware.Responder {
	panic("implement me")
}

func (f fakeInventory) ListClusters(ctx context.Context, params installer.ListClustersParams) middleware.Responder {
	return installer.NewListClustersOK()
}

func (f fakeInventory) ListHosts(ctx context.Context, params installer.ListHostsParams) middleware.Responder {
	panic("implement me")
}

func (f fakeInventory) PostStepReply(ctx context.Context, params installer.PostStepReplyParams) middleware.Responder {
	panic("implement me")
}

func (f fakeInventory) RegisterCluster(ctx context.Context, params installer.RegisterClusterParams) middleware.Responder {
	panic("implement me")
}

func (f fakeInventory) RegisterHost(ctx context.Context, params installer.RegisterHostParams) middleware.Responder {
	panic("implement me")
}

func (f fakeInventory) ResetCluster(ctx context.Context, params installer.ResetClusterParams) middleware.Responder {
	panic("implement me")
}

func (f fakeInventory) SetDebugStep(ctx context.Context, params installer.SetDebugStepParams) middleware.Responder {
	panic("implement me")
}

func (f fakeInventory) UpdateCluster(ctx context.Context, params installer.UpdateClusterParams) middleware.Responder {
	panic("implement me")
}

func (f fakeInventory) UpdateHostInstallProgress(ctx context.Context, params installer.UpdateHostInstallProgressParams) middleware.Responder {
	panic("implement me")
}

func (f fakeInventory) UploadClusterIngressCert(ctx context.Context, params installer.UploadClusterIngressCertParams) middleware.Responder {
	panic("implement me")
}

var _ restapi.InstallerAPI = fakeInventory{}
