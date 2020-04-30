package hardware

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/filanov/bm-inventory/models"
)

func GetHostDisks(host *models.Host) ([]*models.BlockDevice, error) {
	var disks []*models.BlockDevice
	var hwInfo models.Introspection
	if err := json.Unmarshal([]byte(host.HardwareInfo), &hwInfo); err != nil {
		return nil, err
	}
	for _, blockDevice := range hwInfo.BlockDevices {
		if blockDevice.DeviceType == "disk" {
			disks = append(disks, blockDevice)
		}
	}
	// Sorting list by name
	sort.Slice(disks, func(i, j int) bool {
		return disks[i].Name < disks[j].Name
	})
	if len(disks) == 0 {
		return nil, fmt.Errorf("host %s doesn't have disks", host.HostID)
	}
	return disks, nil
}
