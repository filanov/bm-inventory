package hardware

import (
	"encoding/json"
	"testing"

	"github.com/alecthomas/units"
	"github.com/filanov/bm-inventory/models"
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestValidator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Subsystem Suite")
}

var _ = Describe("hardware_validator", func() {
	var (
		hwvalidator Validator
		host        *models.Host
	)
	BeforeEach(func() {
		var cfg ValidatorCfg
		Expect(envconfig.Process("myapp", &cfg)).ShouldNot(HaveOccurred())
		hwvalidator = NewValidator(cfg)
		id := strfmt.UUID(uuid.New().String())
		host = &models.Host{Base: models.Base{ID: &id}, ClusterID: strfmt.UUID(uuid.New().String())}
	})

	It("sufficient_hw", func() {
		hwInfo := &models.Introspection{
			CPU:    &models.CPU{Cpus: 16},
			Memory: []*models.Memory{{Name: "Mem", Total: int64(32 * units.GiB)}},
			BlockDevices: []*models.BlockDevice{
				{DeviceType: "loop", Fstype: "squashfs", MajorDeviceNumber: 7, MinorDeviceNumber: 0, Mountpoint: "/sysroot", Name: "loop0", ReadOnly: true, RemovableDevice: 1, Size: 746217472},
				{DeviceType: "disk", Fstype: "iso9660", MajorDeviceNumber: 11, Mountpoint: "/test", Name: "sdb", RemovableDevice: 1, Size: 822083584}},
		}
		hw, err := json.Marshal(&hwInfo)
		Expect(err).NotTo(HaveOccurred())
		host.HardwareInfo = string(hw)

		roles := []string{"", "master", "worker"}
		for _, role := range roles {
			host.Role = role
			sufficient(hwvalidator.IsSufficient(host))
		}
	})

	It("insufficient_minimal_hw_requirements", func() {
		hwInfo := &models.Introspection{
			CPU:    &models.CPU{Cpus: 1},
			Memory: []*models.Memory{{Name: "Mem", Total: int64(3 * units.GiB)}},
			BlockDevices: []*models.BlockDevice{
				{DeviceType: "loop", Fstype: "squashfs", MajorDeviceNumber: 7, MinorDeviceNumber: 0, Mountpoint: "/sysroot", Name: "loop0", ReadOnly: true, RemovableDevice: 1, Size: 746217472},
				{DeviceType: "disk", Fstype: "iso9660", MajorDeviceNumber: 11, Mountpoint: "/test", Name: "sdb", RemovableDevice: 1, Size: 822083584}},
		}
		hw, err := json.Marshal(&hwInfo)
		Expect(err).NotTo(HaveOccurred())
		host.HardwareInfo = string(hw)

		roles := []string{"", "master", "worker"}
		for _, role := range roles {
			host.Role = role
			insufficient(hwvalidator.IsSufficient(host))
		}
	})

	It("insufficient_master_but_valid_worker", func() {
		hwInfo := &models.Introspection{
			CPU:    &models.CPU{Cpus: 8},
			Memory: []*models.Memory{{Name: "Mem", Total: int64(8 * units.GiB)}},
			BlockDevices: []*models.BlockDevice{
				{DeviceType: "loop", Fstype: "squashfs", MajorDeviceNumber: 7, MinorDeviceNumber: 0, Mountpoint: "/sysroot", Name: "loop0", ReadOnly: true, RemovableDevice: 1, Size: 746217472},
				{DeviceType: "disk", Fstype: "iso9660", MajorDeviceNumber: 11, Mountpoint: "/test", Name: "sdb", RemovableDevice: 1, Size: 822083584}},
		}
		hw, err := json.Marshal(&hwInfo)
		Expect(err).NotTo(HaveOccurred())
		host.HardwareInfo = string(hw)
		host.Role = "master"
		insufficient(hwvalidator.IsSufficient(host))
		host.Role = "worker"
		sufficient(hwvalidator.IsSufficient(host))
	})

	It("insufficient_number_of disks", func() {
		hwInfo := &models.Introspection{
			CPU:    &models.CPU{Cpus: 8},
			Memory: []*models.Memory{{Name: "Mem", Total: int64(8 * units.GiB)}},
			BlockDevices: []*models.BlockDevice{
				{DeviceType: "loop", Fstype: "squashfs", MajorDeviceNumber: 7, MinorDeviceNumber: 0, Mountpoint: "/sysroot", Name: "loop0", ReadOnly: true, RemovableDevice: 1, Size: 746217472}},
		}
		hw, err := json.Marshal(&hwInfo)
		Expect(err).NotTo(HaveOccurred())
		host.HardwareInfo = string(hw)
		insufficient(hwvalidator.IsSufficient(host))
	})

	It("invalid_hw_info", func() {
		host.HardwareInfo = "not a valid json"
		roles := []string{"", "master", "worker"}
		for _, role := range roles {
			host.Role = role
			reply, err := hwvalidator.IsSufficient(host)
			Expect(err).To(HaveOccurred())
			Expect(reply).To(BeNil())
		}
	})
})

func sufficient(reply *IsSufficientReply, err error) {
	Expect(err).NotTo(HaveOccurred())
	Expect(reply.IsSufficient).To(BeTrue())
	Expect(reply.Reason).Should(Equal(""))
}

func insufficient(reply *IsSufficientReply, err error) {
	Expect(err).NotTo(HaveOccurred())
	Expect(reply.IsSufficient).To(BeFalse())
	Expect(reply.Reason).ShouldNot(Equal(""))
}
