package auth

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/filanov/bm-inventory/internal/common"
	"github.com/filanov/bm-inventory/restapi"
	"github.com/sirupsen/logrus"
)

type contextKey string

const (
	bareMetalClusterResource string     = "BareMetalCluster"
	amsActionCreate          string     = "create"
	clustersPathRegex        string     = "clusters/.+"
	capabilityName           string     = "bare_metal_admin"
	capabilityType           string     = "Account"
	adminKey                 contextKey = "Admin"
)

// CreateAuthorizer returns Authorizer if auth is enabled
func CreateAuthorizer(isEnabled bool) func(*http.Request) error {
	if !isEnabled {
		return func(*http.Request) error {
			return nil
		}
	}

	return Authorizer
}

// Authorizer is used to authorize a request after the Auth function was called using the "Auth*" functions
// and the principal was stored in the context in the "AuthKey" context value.
func Authorizer(request *http.Request) error {
	authPayload := PayloadFromContext(request.Context())
	if authPayload == nil || authPayload.Username == "" {
		// auth is disabled
		return nil
	}
	username := authPayload.Username
	allowed, err := allowedToUseAssistedInstaller(username)
	if err != nil {
		logrus.Errorf("Failed to authorize user: %v", err)
	} else if allowed {
		if authPayload.IsAdmin {
			// All endpoints are allowed for admins.
			return nil
		}

		// If endpoint contains a cluster ID, ensures it's owned by the user.
		if clusterID := getClusterIDFromPath(request.URL.Path); clusterID != "" {
			if isClusterOwnedByUser(clusterID, username) {
				// User allowed to manipulate only owned cluster.
				return nil
			}
		} else {
			// Other API endpoints are allowed for all authorized users.
			return nil
		}
	}

	return fmt.Errorf("method is not allowed")
}

// Ensure that the user has authorization to use the bare metal installer service.
// For now the indication is simply "create BareMetalCluster" permission,
// which is allowed for users with BareMetalInstallerUser role.
func allowedToUseAssistedInstaller(username string) (bool, error) {
	return authHandler.client.Authorization.AccessReview(
		context.Background(), username, amsActionCreate, bareMetalClusterResource)
}

// Extracts cluster ID for path if available
func getClusterIDFromPath(path string) string {
	re := regexp.MustCompile(clustersPathRegex)
	if re.MatchString(path) {
		res := strings.Split(re.FindString(path), "/")
		return res[1]
	}
	return ""
}

// Checks whether the cluster owned by the user
func isClusterOwnedByUser(clusterID, username string) bool {
	if err := authHandler.db.First(&common.Cluster{}, "id = ? and user_name = ?", clusterID, username).Error; err == nil {
		// Cluster owned by the user
		return true
	}
	return false
}

// IsAdmin checks whether user is an admin
func IsAdmin(username string) (bool, error) {
	return authHandler.client.Authorization.CapabilityReview(
		context.Background(), fmt.Sprint(username), capabilityName, capabilityType)
}

// PayloadFromContext returns auth payload from the specified context
func PayloadFromContext(ctx context.Context) *AuthPayload {
	payload := ctx.Value(restapi.AuthKey)
	if payload == nil {
		return nil
	}
	return payload.(*AuthPayload)
}

// UserNameFromContext returns username from the specified context
func UserNameFromContext(ctx context.Context) string {
	payload := PayloadFromContext(ctx)
	if payload == nil {
		return ""
	}
	return payload.Username
}
