package config

import (
	_ "embed"
	"fmt"
	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
	"github.com/valyala/fasttemplate"
	"os"
	"path/filepath"
	"runtime"
)

var (
	//go:embed php-fpm.conf
	fpmConfContent string
)

type Config struct {
	PrometheusMetricsEnabled bool              `env:"PROMETHEUS_METRICS_ENABLED"`
	PrometheusMetricsPath    string            `env:"PROMETHEUS_METRICS_PATH"`
	DocumentRoot             string            `env:"DOCUMENT_ROOT"`
	DefaultScript            string            `env:"DEFAULT_SCRIPT"`
	GZIPEnabled              bool              `env:"GZIP_ENABLED"`
	GZIPLevel                int               `env:"GZIP_LEVEL"`
	HTTPListenAddr           string            `env:"HTTP_LISTEN_ADDR"`
	EnableAccessLogs         bool              `env:"ENABLE_ACCESS_LOGS"`
	DataDir                  string            `env:"DATA_DIR"`
	INI                      map[string]string `env:"PHP_INI" envSeparator:";"`
	FPMBinFilename           string            `env:"FPM_BIN_PATH"`
	WorkerCount              int               `env:"WORKER_COUNT"`
	WorkerMaxRequests        int               `env:"WORKER_MAX_REQUEST_COUNT"`
	WorkerTerminationTimeout int               `env:"WORKER_MAX_REQUEST_TIMEOUT"`

	fpmConfigContent  string
	fpmConfigFilename string
	fpmPIDFilename    string
	fpmSocketFilename string
}

func LoadFile(filename string) (*Config, error) {
	if err := godotenv.Overload(filename); err != nil {
		return nil, err
	}

	var c Config

	if err := env.Parse(&c); err != nil {
		return nil, err
	}

	if err := c.ensure(); err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *Config) ensure() error {
	if _, err := os.Stat(c.DefaultScript); err != nil {
		return err
	}

	for _, p := range []*string{&(c.DefaultScript), &(c.DataDir)} {
		if abs, err := filepath.Abs(*p); err != nil {
			return err
		} else {
			*p = abs
		}
	}

	if c.WorkerCount < 1 {
		c.WorkerCount = runtime.NumCPU()
	}

	_ = os.RemoveAll(c.DataDir)

	fpmDir := filepath.Join(c.DataDir, "fpm")

	for _, p := range []string{c.DataDir, fpmDir} {
		if err := os.MkdirAll(p, 0755); err != nil {
			return err
		}
	}

	c.fpmConfigFilename = filepath.Join(fpmDir, "fpm.ini")
	c.fpmSocketFilename = filepath.Join(fpmDir, "fpm.sock")
	c.fpmPIDFilename = filepath.Join(fpmDir, "fpm.pid")

	fpmConfigFileContents := fasttemplate.ExecuteString(fpmConfContent, "{{", "}}", map[string]any{
		"files.pid":              c.fpmPIDFilename,
		"files.socket":           c.fpmSocketFilename,
		"worker.count":           fmt.Sprintf("%v", c.WorkerCount),
		"worker.max_requests":    fmt.Sprintf("%v", c.WorkerMaxRequests),
		"worker.request_timeout": fmt.Sprintf("%v", c.WorkerTerminationTimeout),
	})

	if err := os.WriteFile(fpmConfigFileContents, []byte(fpmConfigFileContents), 0755); err != nil {
		return err
	}

	return nil
}
