package serve

import (
	"github.com/alash3al/phoo/pkg/fastcgi"
	"github.com/alash3al/phoo/pkg/fpm"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
)

type Config struct {
	FPM            fpm.Config
	FastCGI        fastcgi.Config
	HTTPListenAddr string
	EnableLogs     bool
}

func (c *Config) Verify() error {
	paths := []*string{
		&(c.FPM.ConfigFilename),
		&(c.FPM.PIDFilename),
		&(c.FPM.SocketFilename),
		&(c.FastCGI.DefaultScript),
	}

	for _, path := range paths {
		abs, err := filepath.Abs(*path)
		if err != nil {
			return err
		}
		*path = abs
	}

	if _, err := exec.LookPath(c.FPM.Bin); err != nil {
		return err
	}

	c.FastCGI.FastCGIServerURL = "unix://" + c.FPM.SocketFilename

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		_ = <-signalChannel
		c.FPM.Clean()
		os.Exit(0)
	}()

	return nil
}
