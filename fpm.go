package main

import (
	"os"
	"os/exec"
	"strings"
	"time"
)

func startPHPFPM() error {
	cmd := exec.Command(*flagFPMCommand, "-F", "-O", "-y", fpmConfigFilename)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	for _, entry := range strings.Split(*flagPHPINIVars, ";") {
		entry = strings.TrimSpace(entry)
		if entry != "" {
			cmd.Args = append(cmd.Args, "-d", entry)
		}
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	for {
		if _, err := os.Stat(fpmSocketFilename); err != nil {
			server.Logger.Warn("waiting till php-fpm starts ...")
			time.Sleep(1 * time.Second)
			continue
		}

		server.Logger.Info("the php-fpm process has been started ...")
		break
	}

	go (func() {
		if err := cmd.Wait(); err != nil {
			server.Logger.Fatal(err.Error())
		}
	})()

	return nil
}
