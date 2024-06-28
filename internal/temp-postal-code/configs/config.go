package configs

import "github.com/spf13/viper"

type Conf struct {
	OtelExporterEndpoint string `mapstructure:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	ServiceName          string `mapstructure:"SERVICE_NAME"`
	ViaCepApiUrl         string `mapstructure:"VIACEP_API_URL"`
	WeatherApiUrl        string `mapstructure:"WEATHER_API_URL"`
	WeatherApiKey        string `mapstructure:"WEATHER_API_KEY"`
	WebServerUrl         string `mapstructure:"WEBSERVER_URL"`
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
	viper.BindEnv("VIACEP_API_URL")
	viper.BindEnv("WEATHER_API_URL")
	viper.BindEnv("WEATHER_API_KEY")
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
