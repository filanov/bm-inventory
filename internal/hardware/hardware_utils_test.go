package hardware

import (
	"encoding/json"
	"testing"

	"github.com/alecthomas/units"
	"github.com/filanov/bm-inventory/models"
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHardwareUtils(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Subsystem Suite")
}

var _ = Describe("hardware_utils_tests", func() {
	var (
		host *models.Host
	)
	BeforeEach(func() {
		id := strfmt.UUID(uuid.New().String())
		host = &models.Host{Base: models.Base{ID: &id}, ClusterID: strfmt.UUID(uuid.New().String())}
	})

	It("disk list, verify ordering", func() {
		var hwInfo = &models.Introspection{
			CPU:    &models.CPU{Cpus: 16},
			Memory: []*models.Memory{{Name: "Mem", Total: int64(32 * units.GiB)}},
			BlockDevices: []*models.BlockDevice{
				{DeviceType: "loop", Fstype: "squashfs", MajorDeviceNumber: 7, MinorDeviceNumber: 0, Mountpoint: "/sysroot", Name: "loop0", ReadOnly: true, RemovableDevice: 1, Size: 746217472},
				{DeviceType: "disk", Fstype: "iso9660", MajorDeviceNumber: 11, Mountpoint: "/test", Name: "sdb", RemovableDevice: 1, Size: 822083584},
				{DeviceType: "disk", Fstype: "iso9660", MajorDeviceNumber: 11, Mountpoint: "/test", Name: "vda", RemovableDevice: 1, Size: 822083584},
				{DeviceType: "disk", Fstype: "iso9660", MajorDeviceNumber: 11, Mountpoint: "/test", Name: "sda", RemovableDevice: 1, Size: 822083584}},
		}
		hw, err := json.Marshal(&hwInfo)
		Expect(err).NotTo(HaveOccurred())
		host.HardwareInfo = string(hw)
		reply, err := GetHostDisks(host)
		Expect(err).NotTo(HaveOccurred())
		Expect(reply[0].Name).Should(Equal("sda"))
		Expect(len(reply)).Should(Equal(3))
	})
	It("host with no disks", func() {
		var hwInfo = &models.Introspection{
			CPU:    &models.CPU{Cpus: 16},
			Memory: []*models.Memory{{Name: "Mem", Total: int64(32 * units.GiB)}},
			BlockDevices: []*models.BlockDevice{
				{DeviceType: "loop", Fstype: "squashfs", MajorDeviceNumber: 7, MinorDeviceNumber: 0, Mountpoint: "/sysroot", Name: "loop0", ReadOnly: true, RemovableDevice: 1, Size: 746217472},
				{DeviceType: "smth", Fstype: "iso9660", MajorDeviceNumber: 11, Mountpoint: "/test", Name: "sdb", RemovableDevice: 1, Size: 822083584},
				{DeviceType: "smth", Fstype: "iso9660", MajorDeviceNumber: 11, Mountpoint: "/test", Name: "vda", RemovableDevice: 1, Size: 822083584},
				{DeviceType: "smth", Fstype: "iso9660", MajorDeviceNumber: 11, Mountpoint: "/test", Name: "sda", RemovableDevice: 1, Size: 822083584}},
		}
		hw, err := json.Marshal(&hwInfo)
		Expect(err).NotTo(HaveOccurred())
		host.HardwareInfo = string(hw)
		_, err = GetHostDisks(host)
		Expect(err).Should(HaveOccurred())
	})
})
