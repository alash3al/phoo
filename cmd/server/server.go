package server

import (
	"fmt"
	"github.com/alash3al/phoo/internals/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"os/exec"
	"time"
)

func Serve(conf *config.Config) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		e := echo.New()
		e.HideBanner = true

		if conf.EnableAccessLogs {
			e.Use(middleware.Logger())
		}

		e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
			Level: conf.GZIPLevel,
		}))

		e.Use(middleware.Recover())
		e.Use(serveStaticFilesOnlyMiddleware(conf.DocumentRoot))
		e.Use(serveFastCGIMiddleware(
			conf.DefaultScript,
			"unix",
			conf.FPM.SocketFilename,
		))

		if err := startPHPFPM(conf); err != nil {
			e.Logger.Fatal(err.Error())
		}

		return e.Start(conf.HTTPListenAddr)
	}
}

func startPHPFPM(conf *config.Config) error {
	cmd := exec.Command(conf.FPM.Bin, "-F", "-O", "-y", conf.FPM.ConfigFilename)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	for k, v := range conf.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	cmd.Env = append(cmd.Env, os.Environ()...)

	for k, v := range conf.INI {
		cmd.Args = append(cmd.Args, "-d", fmt.Sprintf("%s=%s", k, v))
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	for {
		if _, err := os.Stat(conf.FPM.SocketFilename); err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		break
	}

	go (func() {
		if err := cmd.Wait(); err != nil {
			log.Fatal(err.Error())
		}
	})()

	return nil
}
