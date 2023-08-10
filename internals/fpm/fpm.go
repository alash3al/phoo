package fpm

import (
	_ "embed"
	"fmt"
	"github.com/valyala/fasttemplate"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

//go:embed php-fpm.conf
var configTemplate string

type Process struct {
	BinFilename    string
	PIDFilename    string
	ConfigFilename string
	SocketFilename string
	INI            []string

	WorkerCount           int
	WorkerMaxRequestCount int
	WorkerMaxRequestTime  int
}

func (p *Process) Start() error {
	paths := []*string{
		&p.ConfigFilename,
		&p.SocketFilename,
		&p.PIDFilename,
	}

	for _, path := range paths {
		abs, err := filepath.Abs(*path)
		if err != nil {
			return err
		}

		*path = abs
	}

	fpmConfigFileContents := fasttemplate.ExecuteString(configTemplate, "{{", "}}", map[string]any{
		"files.pid":                p.PIDFilename,
		"files.socket":             p.SocketFilename,
		"worker.count":             fmt.Sprintf("%v", p.WorkerCount),
		"worker.request.max_count": fmt.Sprintf("%v", p.WorkerMaxRequestCount),
		"worker.request.max_time":  fmt.Sprintf("%v", p.WorkerMaxRequestTime),
	})

	if err := os.WriteFile(p.ConfigFilename, []byte(fpmConfigFileContents), 0755); err != nil {
		return err
	}

	return p.execAndWait()
}

func (p *Process) execAndWait() error {
	cmd := exec.Command(p.BinFilename, "-F", "-O", "-y", p.ConfigFilename)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	for _, v := range p.INI {
		cmd.Args = append(cmd.Args, "-d", v)
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	for {
		if _, err := os.Stat(p.SocketFilename); err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		break
	}

	sig := make(chan os.Signal)

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for s := range sig {
			_ = cmd.Process.Signal(s)
			os.Exit(0)
		}
	}()

	return nil
}
