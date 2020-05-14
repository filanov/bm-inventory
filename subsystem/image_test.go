package subsystem

import (
	"context"
	"io/ioutil"
	"log"
	"os"

	"github.com/filanov/bm-inventory/client/installer"
	"github.com/filanov/bm-inventory/models"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("system-test image tests", func() {
	ctx := context.Background()
	var cluster *installer.RegisterClusterCreated
	var clusterID strfmt.UUID

	AfterEach(func() {
		clearDB()
	})

	BeforeEach(func() {
		var err error
		cluster, err = bmclient.Installer.RegisterCluster(ctx, &installer.RegisterClusterParams{
			NewClusterParams: &models.ClusterCreateParams{
				Name:             swag.String("test cluster"),
				OpenshiftVersion: swag.String("4.4"),
			},
		})
		Expect(err).NotTo(HaveOccurred())
		clusterID = *cluster.GetPayload().ID
	})

	It("create_and_get_image", func() {
		file, err := ioutil.TempFile("", "tmp")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(file.Name())

		_, err = bmclient.Installer.GenerateClusterISO(ctx, &installer.GenerateClusterISOParams{
			ClusterID:         clusterID,
			ImageCreateParams: &models.ImageCreateParams{},
		})
		Expect(err).NotTo(HaveOccurred())
		_, err = bmclient.Installer.DownloadClusterISO(ctx, &installer.DownloadClusterISOParams{
			ClusterID: clusterID,
		}, file)
		Expect(err).NotTo(HaveOccurred())
		s, err := file.Stat()
		Expect(err).NotTo(HaveOccurred())
		Expect(s.Size()).ShouldNot(Equal(0))
	})
})

var _ = Describe("image tests", func() {
	ctx := context.Background()
	var file *os.File
	var err error

	AfterEach(func() {
		clearDB()
		os.Remove(file.Name())
	})

	BeforeEach(func() {
		file, err = ioutil.TempFile("", "tmp")
		Expect(err).To(BeNil())
	})

	It("download_non_existing_cluster", func() {
		_, err = bmclient.Installer.DownloadClusterISO(ctx, &installer.DownloadClusterISOParams{ClusterID: *strToUUID(uuid.New().String())}, file)
		Expect(err).Should(HaveOccurred())
	})
})
