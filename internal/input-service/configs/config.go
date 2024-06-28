package configs

import (
	"github.com/spf13/viper"
)

type Conf struct {
	OtelExporterEndpoint  string `mapstructure:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	ServiceName           string `mapstructure:"SERVICE_NAME"`
	TemperatureServiceUrl string `mapstructure:"TEMPERATURE_SERVICE_URL"`
	WebServerUrl          string `mapstructure:"WEBSERVER_URL"`
}

func LoadConfig() (*Conf, error) {
	var cfg *Conf

	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	viper.BindEnv("OTEL_EXPORTER_OTLP_ENDPOINT")
	viper.BindEnv("SERVICE_NAME")
	viper.BindEnv("TEMPERATURE_SERVICE_URL")
	viper.BindEnv("WEBSERVER_URL")
	err := viper.ReadInConfig()

	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	return cfg, err
}
