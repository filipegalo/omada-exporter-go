package internal

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"sync"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	LogLevel string `env:"LOG_LEVEL" envDefault:"error"`
	Omada    struct {
		OmadaURL      string `env:"OMADA_URL,required"`
		SiteName      string `env:"OMADA_SITE_NAME,required"`
		ClientID      string `env:"OMADA_CLIENT_ID,required"`
		ClientSecret  string `env:"OMADA_CLIENT_SECRET,required"`
		Username      string `env:"OMADA_USERNAME,required"`
		Password      string `env:"OMADA_PASSWORD,required"`
		SkipTLSVerify bool   `env:"OMADA_SKIP_TLS_VERIFY" envDefault:"true"`
	}
	Prometheus struct {
		MetricsPath string `env:"METRICS_PATH" envDefault:"/metrics"`
		MetricsPort string `env:"METRICS_PORT" envDefault:"8080"`
	}
}

var (
	instance *Config
	once     sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		var err error
		instance, err = loadConfig()
		if err != nil {
			log.Fatalf("Failed to load configuration: %v", err)
		}
	})
	return instance
}

func loadConfig() (*Config, error) {
	// Load .env if present; ignore error so real env vars (e.g. Docker) still work.
	_ = godotenv.Load()

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) validate() error {
	u, err := url.ParseRequestURI(c.Omada.OmadaURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return fmt.Errorf("OMADA_URL must be a valid http(s) URL, got %q", c.Omada.OmadaURL)
	}

	port, err := strconv.Atoi(c.Prometheus.MetricsPort)
	if err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("METRICS_PORT must be an integer between 1 and 65535, got %q", c.Prometheus.MetricsPort)
	}

	return nil
}
