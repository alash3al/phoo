package config

import (
	_ "embed"
	"fmt"
	"github.com/valyala/fasttemplate"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

var (
	//go:embed php-fpm.conf
	fpmConfContent string
)

type Config struct {
	DocumentRoot     string `yaml:"document_root"`
	DefaultScript    string `yaml:"default_script"`
	GZIPLevel        int    `yaml:"gzip_level"`
	HTTPListenAddr   string `yaml:"http_listen_addr"`
	EnableAccessLogs bool   `yaml:"enable_access_logs"`
	DataDir          string `yaml:"data_dir"`

	Env map[string]string `yaml:"env"`

	INI map[string]string `yaml:"ini"`

	FPM struct {
		Bin                string    `yaml:"bin"`
		WorkerCount        int       `yaml:"worker_count"`
		WorkerMaxRequests  int       `yaml:"worker_max_requests"`
		TerminationTimeout int       `yaml:"termination_timeout"`
		ConfigContent      string    `yaml:"-"`
		ConfigFilename     string    `yaml:"-"`
		PIDFilename        string    `yaml:"-"`
		SocketFilename     string    `yaml:"-"`
		Command            *exec.Cmd `yaml:"-"`
	} `yaml:"fpm"`
}

func LoadFile(filename string) (*Config, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	content = []byte(os.ExpandEnv(string(content)))

	var c Config

	if err := yaml.Unmarshal(content, &c); err != nil {
		return nil, err
	}

	if err := c.verify(); err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *Config) verify() error {
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

	_ = os.RemoveAll(c.DataDir)

	fpmDir := filepath.Join(c.DataDir, "fpm")

	for _, p := range []string{c.DataDir, fpmDir} {
		if err := os.MkdirAll(p, 0755); err != nil {
			return err
		}
	}

	if c.FPM.WorkerCount < 1 {
		c.FPM.WorkerCount = runtime.NumCPU()
	}

	c.FPM.ConfigFilename = filepath.Join(fpmDir, "fpm.ini")
	c.FPM.PIDFilename = filepath.Join(fpmDir, "fpm.pid")
	c.FPM.SocketFilename = filepath.Join(fpmDir, "fpm.sock")

	fpmConfContent = fasttemplate.ExecuteString(fpmConfContent, "{{", "}}", map[string]any{
		"files.pid":           c.FPM.PIDFilename,
		"files.socket":        c.FPM.SocketFilename,
		"worker.count":        fmt.Sprintf("%v", c.FPM.WorkerCount),
		"worker.max_requests": fmt.Sprintf("%v", c.FPM.WorkerMaxRequests),
		"termination_timeout": fmt.Sprintf("%v", c.FPM.TerminationTimeout),
	})

	if err := os.WriteFile(c.FPM.ConfigFilename, []byte(fpmConfContent), 0755); err != nil {
		return err
	}

	return nil
}
