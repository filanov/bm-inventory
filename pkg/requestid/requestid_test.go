package requestid

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/filanov/bm-inventory/client"
	installer2 "github.com/filanov/bm-inventory/client/installer"
	"github.com/google/uuid"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/runtime/security"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"

	"github.com/filanov/bm-inventory/restapi"
	"github.com/filanov/bm-inventory/restapi/operations/installer"
	"github.com/sirupsen/logrus"
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
	agentKey := "X-Secret-Key"
	userKey := "X-User-Key"
	agentKeyValue := "SecretKey"
	userKeyValue := "userKey"

	t.Parallel()
	tests := []struct {
		name                   string
		tokenKey               string
		expectedTokenValue     string
		isListOperation        bool
		enableAuth             bool
		addHeaders             bool
		expectedRequestSuccess bool
	}{
		/*
			{
				name:                   "agent auth",
				tokenKey:               agentKey,
				expectedTokenValue:     agentKeyValue,
			    isListOperation:         false,
				enableAuth:             true,
				addHeaders:             true,
				expectedRequestSuccess: true,
			},
		*/
		{
			name:                   "user auth",
			tokenKey:               userKey,
			expectedTokenValue:     userKeyValue,
			isListOperation:        true,
			enableAuth:             true,
			addHeaders:             true,
			expectedRequestSuccess: true,
		},
		{
			name:                   "Fail auth without headers",
			tokenKey:               agentKey,
			expectedTokenValue:     agentKeyValue,
			isListOperation:        false,
			enableAuth:             true,
			addHeaders:             false,
			expectedRequestSuccess: false,
		},
		{
			name:                   "Ignore auth if disabled",
			tokenKey:               userKey,
			expectedTokenValue:     userKeyValue,
			isListOperation:        true,
			enableAuth:             false,
			addHeaders:             false,
			expectedRequestSuccess: true,
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
				AuthAgentAuth:       authAgentAuth,
				AuthUserAuth:        authUserAuth,
				APIKeyAuthenticator: createAuthenticator(tt.enableAuth),
				InstallerAPI:        fakeInventory{},
				EventsAPI:           nil,
				Logger:              logrus.Printf,
				VersionsAPI:         nil,
				ManagedDomainsAPI:   nil,
				InnerMiddleware:     nil,
			})

			clientAuth := func() runtime.ClientAuthInfoWriter {
				return runtime.ClientAuthInfoWriterFunc(func(r runtime.ClientRequest, _ strfmt.Registry) error {
					return r.SetHeaderParam(tt.tokenKey, tt.expectedTokenValue)
				})
			}

			cfg := client.Config{
				URL: &url.URL{
					Scheme: client.DefaultSchemes[0],
					Host:   "localhost:8081",
					Path:   client.DefaultBasePath,
				},
			}
			if tt.addHeaders {
				cfg.AuthInfo = clientAuth()
			}
			bmclient := client.New(cfg)

			server := &http.Server{Addr: "localhost:8081", Handler: h}
			go server.ListenAndServe()
			defer server.Close()

			expectedStatusCode := 401
			if tt.expectedRequestSuccess {
				expectedStatusCode = 200
			}

			var e error
			if tt.isListOperation {
				_, e = bmclient.Installer.ListClusters(context.TODO(), &installer2.ListClustersParams{})
			} else {
				id := uuid.New()
				_, e = bmclient.Installer.GetCluster(context.TODO(), &installer2.GetClusterParams{
					ClusterID: strfmt.UUID(id.String()),
				})
			}
			if expectedStatusCode == 200 {
				assert.Nil(t, e)
			} else {
				apierr := e.(*runtime.APIError)
				assert.Equal(t, apierr.Code, expectedStatusCode)

			}
		})
	}
}

func createAuthenticator(isEnabled bool) func(name, in string, authenticate security.TokenAuthentication) runtime.Authenticator {
	return func(name string, _ string, authenticate security.TokenAuthentication) runtime.Authenticator {
		getToken := func(r *http.Request) string { return r.Header.Get(name) }

		return security.HttpAuthenticator(func(r *http.Request) (bool, interface{}, error) {
			if !isEnabled {
				return true, "", nil
			}
			token := getToken(r)
			if token == "" {
				return false, nil, nil
			}

			p, err := authenticate(token)
			return true, p, err
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
