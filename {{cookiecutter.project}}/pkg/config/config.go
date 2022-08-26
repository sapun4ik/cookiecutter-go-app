package config

import (
	"os"
	"time"

	"{{cookiecutter.module_path}}/pkg/probes"
	"{{cookiecutter.module_path}}/pkg/tracing"

	"github.com/spf13/viper"
)

const (
	defaultHTTPPort               = "8000"
	defaultHTTPRWTimeout          = 10 * time.Second
	defaultHTTPMaxHeaderMegabytes = 1

	ProfileDev = "dev"
	Prod       = "prod"
)

var (
	ignoreLogUrls = []string{"metrics", "swagger"}
)

// App config struct.
type Config struct {
	Logger      LoggerConfig    `mapstructure:"logger"`
	Application AppConfig       `mapstructure:"application"`
	HTTP        HTTPConfig      `mapstructure:"http"`
	Probes      probes.Config   `mapstructure:"probes"`
	Postgres    PostgresConfig  `mapstructure:"postgres"`
	Jaeger      *tracing.Config `mapstructure:"jaeger"`
}

// Server config struct.
type AppConfig struct {
	Version string `mapstructure:"version"`
	Name    string `mapstructure:"name"`
	Profile string `mapstructure:"profile"`
}

type HTTPConfig struct {
	Port               string        `mapstructure:"port"`
	ReadTimeout        time.Duration `mapstructure:"readTimeout"`
	WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
	MaxHeaderMegabytes int           `mapstructure:"maxHeaderMegabytes"`
	IgnoreLogUrls      []string      `mapstructure:"ignoreLogUrls"`
}

// LoggerConfig config.
type LoggerConfig struct {
	Development       bool   `mapstructure:"development"`
	DisableCaller     bool   `mapstructure:"disableCaller"`
	DisableStacktrace bool   `mapstructure:"disableStacktrace"`
	DevMode           bool   `mapstructure:"devMode"`
	Encoding          string `mapstructure:"encoding"`
	Level             string `mapstructure:"level"`
}

// Postgresql config.
type PostgresConfig struct {
	Host            string        `mapstructure:"host"`
	Port            string        `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	DBName          string        `mapstructure:"dbName"`
	SSLMode         bool          `mapstructure:"sslMode"`
	PGDriver        string        `mapstructure:"pgDriver"`
	MaxOpenConns    int           `mapstructure:"maxOpenConns"`
	ConnMaxLifetime time.Duration `mapstructure:"connMaxLifetime"`
	MaxIdleConns    int           `mapstructure:"maxIdleConns"`
	ConnMaxIdleTime time.Duration `mapstructure:"connMaxIdleTime"`
}

// Init populates Config struct with values from config file
// located at filepath and environment variables.
func Init(configsDir string) (*Config, error) {
	populateDefaults()

	if err := parseConfigFile(configsDir, os.Getenv("APP_PROFILE")); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	setFromEnv(&cfg)

	return &cfg, nil
}

func parseConfigFile(folder, profile string) error {
	viper.AddConfigPath(folder)
	viper.SetConfigName("default")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	viper.SetConfigName(profile)

	return viper.MergeInConfig()
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("http", &cfg.HTTP); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("logger", &cfg.Logger); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("postgres", &cfg.Postgres); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("jaeger", &cfg.Jaeger); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("probes", &cfg.Probes); err != nil {
		return err
	}

	return viper.UnmarshalKey("application", &cfg.Application)
}

func setFromEnv(cfg *Config) {
	if os.Getenv("HTTP_PORT") != "" {
		cfg.HTTP.Port = os.Getenv("HTTP_PORT")
	}

	cfg.Application.Profile = os.Getenv("APP_PROFILE")
}

func populateDefaults() {
	viper.SetDefault("http.port", defaultHTTPPort)
	viper.SetDefault("http.maxHeaderMegabytes", defaultHTTPMaxHeaderMegabytes)
	viper.SetDefault("http.ignoreLogUrls", ignoreLogUrls)
	viper.SetDefault("http.readTimeout", defaultHTTPRWTimeout)
	viper.SetDefault("http.writeTimeout", defaultHTTPRWTimeout)
}
