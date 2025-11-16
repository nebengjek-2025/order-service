package apm

import (
	"order-service/src/internal/config"
	"os"

	"github.com/spf13/viper"
	"go.elastic.co/apm"
	"go.elastic.co/apm/transport"
)

func InitConnection() {
	viper := config.GetConfig()
	os.Setenv("ELASTIC_APM_SERVER_URL", viper.GetString("apm.host"))
	os.Setenv("ELASTIC_APM_SECRET_TOKEN", viper.GetString("apm.token"))

	if _, err := transport.InitDefault(); err != nil {
		panic(err)
	}
}

func GetTracer() *apm.Tracer {
	tracer, _ := apm.NewTracer(viper.GetString("app.name"), viper.GetString("app.version"))
	return tracer
}
