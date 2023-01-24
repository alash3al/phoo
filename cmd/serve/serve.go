package serve

import (
	"fmt"
	"github.com/NYTimes/gziphandler"
	"github.com/alash3al/phoo/pkg/fastcgi"
	"github.com/alash3al/phoo/pkg/fpm"
	"github.com/alash3al/phoo/pkg/symbols"
	"github.com/labstack/gommon/log"
	"github.com/urfave/cli/v2"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:        "serve",
		Description: "start the http server",
		Action:      listenAndServe(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     symbols.FlagNameHTTPListenAddr,
				Usage:    "the http address to listen on",
				EnvVars:  []string{symbols.EnvKeyListen},
				Required: true,
				Category: symbols.AppName,
			},
			&cli.StringFlag{
				Name:     symbols.FlagNameDocumentRoot,
				Usage:    "the document root",
				EnvVars:  []string{symbols.EnvKeyDocumentRoot},
				Required: true,
				Category: symbols.AppName,
			},
			&cli.BoolFlag{
				Name:     symbols.FlagNameEnableLogs,
				Usage:    "whether to enable/disable access log",
				EnvVars:  []string{symbols.EnvKeyEnableLogs},
				Value:    true,
				Category: symbols.AppName,
			},
			&cli.StringFlag{
				Name:     symbols.FlagNamePHPFPMBinary,
				Usage:    "the PHP-FPM binary",
				EnvVars:  []string{symbols.EnvKeyFPMBin},
				Value:    "php-fpm",
				Category: "php",
			},
			&cli.StringFlag{
				Name:     symbols.FlagNamePHPINI,
				Usage:    "additional PHP INI settings separated with semicolon (;)",
				EnvVars:  []string{symbols.EnvKeyPHPINI},
				Category: "php",
			},
			&cli.StringFlag{
				Name:     symbols.FlagNamePHPUser,
				Usage:    "the user who will PHP-FPM listen as",
				EnvVars:  []string{symbols.EnvKeyPHPUser},
				Value:    "www-data",
				Category: "php",
			},
			&cli.StringFlag{
				Name:     symbols.FlagNamePHPGroup,
				Usage:    "the group who will PHP-FPM listen as",
				EnvVars:  []string{symbols.EnvKeyPHPGroup},
				Value:    "www-data",
				Category: "php",
			},
			&cli.Int64Flag{
				Name:     symbols.FlagNameWorkersCount,
				Usage:    "the PHP workers count",
				EnvVars:  []string{symbols.EnvKeyWorkersCount},
				Value:    int64(runtime.NumCPU()),
				Category: "php",
			},
			&cli.Int64Flag{
				Name:     symbols.FlagNameWorkerMaxRequests,
				Usage:    "the PHP worker max requests (the worker will be restarted after reaching this value)",
				EnvVars:  []string{symbols.EnvKeyWorkerMaxRequests},
				Value:    500,
				Category: "php",
			},
			&cli.StringFlag{
				Name:     symbols.FlagNameRequestTimeout,
				Usage:    "the request timeout",
				EnvVars:  []string{symbols.EnvKeyRequestTimeout},
				Value:    "0",
				Category: "php",
			},
			&cli.StringFlag{
				Name:     symbols.FlagNameDefaultScript,
				Usage:    "the default script used as router",
				EnvVars:  []string{symbols.EnvKeyRouter},
				Required: true,
				Category: "php",
			},
		},
	}
}

func listenAndServe() cli.ActionFunc {
	return func(cliCtx *cli.Context) error {
		config := Config{
			HTTPListenAddr: cliCtx.String(symbols.FlagNameHTTPListenAddr),
			DocumentRoot:   cliCtx.String(symbols.FlagNameDocumentRoot),
			EnableLogs:     cliCtx.Bool(symbols.FlagNameEnableLogs),
			FPM: fpm.Config{
				ConfigFilename:    ".fpm.config.ini",
				PIDFilename:       ".fpm.pid",
				SocketFilename:    ".fpm.sock",
				User:              cliCtx.String(symbols.FlagNamePHPUser),
				Group:             cliCtx.String(symbols.FlagNamePHPGroup),
				Bin:               cliCtx.String(symbols.FlagNamePHPFPMBinary),
				RequestTimeout:    cliCtx.String(symbols.FlagNameRequestTimeout),
				WorkerMaxRequests: cliCtx.Int64(symbols.FlagNameWorkerMaxRequests),
				WorkersCount:      cliCtx.Int64(symbols.FlagNameWorkersCount),
				INI:               strings.Split(cliCtx.String(symbols.FlagNamePHPINI), ";"),
				Stdout:            os.Stdout,
				Stderr:            os.Stderr,
			},
			FastCGI: fastcgi.Config{
				DefaultScript:          cliCtx.String(symbols.FlagNameDefaultScript),
				RestrictDotFilesAccess: true,
				FastCGIParams: map[string]string{
					"SERVER_SOFTWARE": fmt.Sprintf("%s/%s", symbols.AppName, symbols.AppVersion),
				},
			},
		}

		if err := config.Verify(); err != nil {
			return err
		}

		fastCGIHandler, err := fastcgi.New(config.FastCGI)
		if err != nil {
			return err
		}

		mainHandler, err := assetsCacheMiddleware(&config, recoverMiddleware(
			fastCGIHandler.ServeHTTP,
		))

		if err != nil {
			return err
		}

		runner, err := fpm.New(config.FPM)
		if err != nil {
			return err
		}

		go (func() {
			if err := runner.Wait(); err != nil {
				log.Fatal(err.Error())
			}
		})()

		for runner.Process == nil {
			log.Info("waiting php-fpm to be ready")
			time.Sleep(5 * time.Second)
		}

		log.Infoj(map[string]interface{}{
			"message": "configurations",
			"configs": config,
			"fpm-cmd": runner.String(),
		})

		return http.ListenAndServe(
			config.HTTPListenAddr,
			gziphandler.GzipHandler(loggerMiddleware(
				config.EnableLogs,
				mainHandler,
			)),
		)
	}
}
