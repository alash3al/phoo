package fpm

import (
	"bytes"
	_ "embed"
	"io"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

//go:embed config.ini
var configTemplate string

type Config struct {
	PIDFilename       string
	SocketFilename    string
	ConfigFilename    string
	User              string
	Group             string
	Bin               string
	WorkersCount      int64
	WorkerMaxRequests int64
	RequestTimeout    string
	INI               []string
	Stdout            io.Writer
	Stderr            io.Writer
}

func (c *Config) Clean() {
	os.Remove(c.ConfigFilename)
	os.Remove(c.PIDFilename)
	os.Remove(c.SocketFilename)
}

func New(config Config) (*exec.Cmd, error) {
	tpl, err := template.New("template").Parse(configTemplate)
	if err != nil {
		return nil, err
	}

	var finalConfig bytes.Buffer

	if err := tpl.Execute(&finalConfig, config); err != nil {
		return nil, err
	}

	if err := os.WriteFile(config.ConfigFilename, finalConfig.Bytes(), 0755); err != nil {
		return nil, err
	}

	cmd := exec.Command(config.Bin, "-F", "-O", "-y", config.ConfigFilename)
	cmd.Stdout = config.Stdout
	cmd.Stderr = config.Stderr

	for _, entry := range config.INI {
		entry = strings.TrimSpace(entry)
		if entry != "" {
			cmd.Args = append(cmd.Args, "-d", entry)
		}
	}

	return cmd, cmd.Start()
}
