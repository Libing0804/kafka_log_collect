package influxDB

import (
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/sirupsen/logrus"
)

var Cli client.Client

func InitInfluxDB() (err error) {
	Cli, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://127.0.0.1:8086",
		Username: "admin",
		Password: "",
	})
	logrus.Info("influx init success!")
	return
}
