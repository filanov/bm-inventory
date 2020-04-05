package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/filanov/bm-inventory/internal/bminventory"
	"github.com/filanov/bm-inventory/models"
	"github.com/filanov/bm-inventory/restapi"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func init() {
	strfmt.MarshalFormat = strfmt.ISO8601LocalTime
}

var Options struct {
	BMConfig bminventory.Config
	DBHost   string `envconfig:"DB_HOST" default:"mariadb"`
	DBPort   string `envconfig:"DB_PORT" default:"3306"`
}

func main() {
	err := envconfig.Process("myapp", &Options)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	port := flag.String("port", "8090", "define port that the service will listen to")
	flag.Parse()

	logrus.Println("Starting bm service")

	db, err := gorm.Open("mysql",
		fmt.Sprintf("admin:admin@tcp(%s:%s)/installer?charset=utf8&parseTime=True&loc=Local",
			Options.DBHost, Options.DBPort))

	if err != nil {
		logrus.Fatal("Fail to connect to DB, ", err)
	}
	defer db.Close()

	scheme := runtime.NewScheme()
	if err = clientgoscheme.AddToScheme(scheme); err != nil {
		logrus.Fatal()
	}

	kclient, err := client.New(config.GetConfigOrDie(), client.Options{Scheme: scheme})
	if err != nil {
		logrus.Fatal("failed to create client:", err)
	}

	if err = db.AutoMigrate(&models.Host{}, &models.Cluster{}).Error; err != nil {
		logrus.Fatal("failed to auto migrate, ", err)
	}

	bm, err := bminventory.NewBareMetalInventory(db, kclient, Options.BMConfig)
	if err != nil {
		logrus.Fatalf("Error creating baremetal inventory: %s", err.Error())
	}
	h, err := restapi.Handler(restapi.Config{
		InventoryAPI: bm,
		Logger:       logrus.Printf,
	})
	if err != nil {
		logrus.Fatal("Failed to init rest handler,", err)
	}

	logrus.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", swag.StringValue(port)), h))
}
