package host

import (
	"context"

	"github.com/filanov/bm-inventory/models"
	"github.com/filanov/bm-inventory/pkg/requestid"
)

func (m *Manager) HostMonitoring() {
	var (
		hosts     []*models.Host
		requestID = requestid.NewID()
		ctx       = requestid.ToContext(context.Background(), requestID)
		log       = requestid.RequestIDLogger(m.log, requestID)
	)

	monitorStates := []string{HostStatusDiscovering, HostStatusKnown, HostStatusDisconnected, HostStatusInsufficient, HostStatusPendingForInput}
	if err := m.db.Where("status IN (?)", monitorStates).Find(&hosts).Error; err != nil {
		log.WithError(err).Errorf("failed to get hosts")
		return
	}
	for _, host := range hosts {
		err := m.RefreshStatus(ctx, host, m.db)
		if err != nil {
			m.log.WithError(err).Errorf("failed to refresh host %s state", host.ID)
		}
	}
}
