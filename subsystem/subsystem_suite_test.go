package subsystem

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/go-openapi/runtime"

	"github.com/filanov/bm-inventory/client"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/kelseyhightower/envconfig"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"

	"github.com/filanov/bm-inventory/pkg/auth"
)

var db *gorm.DB
var agentBMClient, userBMClient *client.AssistedInstall
var log *logrus.Logger

var Options struct {
	DBHost        string `envconfig:"DB_HOST"`
	DBPort        string `envconfig:"DB_PORT"`
	EnableAuth    bool   `envconfig:"ENABLE_AUTH"`
	InventoryHost string `envconfig:"INVENTORY"`
	TestToken     string `envconfig:"TEST_TOKEN"`
}

func clientcfg(authInfo runtime.ClientAuthInfoWriter) client.Config {
	cfg := client.Config{
		URL: &url.URL{
			Scheme: client.DefaultSchemes[0],
			Host:   Options.InventoryHost,
			Path:   client.DefaultBasePath,
		},
	}
	if Options.EnableAuth {
		log.Info("API Key authentication enabled for subsystem tests")
		cfg.AuthInfo = authInfo
	}
	return cfg
}

func init() {
	var err error
	log = logrus.New()
	log.SetReportCaller(true)
	err = envconfig.Process("subsystem", &Options)
	if err != nil {
		log.Fatal(err.Error())
	}
	userClientCfg := clientcfg(auth.UserAuthHeaderWriter("bearer " + Options.TestToken))
	AgentClientCfg := clientcfg(auth.AgentAuthHeaderWriter("fake_pull_secret"))
	userBMClient = client.New(userClientCfg)
	agentBMClient = client.New(AgentClientCfg)

	db, err = gorm.Open("postgres",
		fmt.Sprintf("host=%s port=%s user=admin dbname=installer password=admin sslmode=disable",
			Options.DBHost, Options.DBPort))
	if err != nil {
		logrus.Fatal("Fail to connect to DB, ", err)
	}
}

func TestSubsystem(t *testing.T) {
	RegisterFailHandler(Fail)
	clearDB()
	RunSpecs(t, "Subsystem Suite")
}
