package server

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

func DefaultFlags(envPrefix string, append ...[]cli.Flag) []cli.Flag {
	prefixWrapper := func(k string) string {
		return fmt.Sprintf("%s_%s", envPrefix, k)
	}

	return []cli.Flag{
		&cli.StringFlag{
			Name:    "metrics",
			Usage:   "the prometheus metrics endpoint, empty means disabled",
			EnvVars: []string{prefixWrapper("METRICS_PATH")},
		},

		&cli.StringFlag{
			Name:     "root",
			Usage:    "the document root full path",
			EnvVars:  []string{prefixWrapper("DOCUMENT_ROOT")},
			Required: true,
		},

		&cli.StringFlag{
			Name:     "entrypoint",
			Usage:    "the default php entrypoint script",
			EnvVars:  []string{"ENTRYPOINT"},
			Required: true,
		},

		&cli.StringFlag{
			Name:    "http",
			Usage:   "the http listen address in the form of [address]:port",
			EnvVars: []string{prefixWrapper("HTTP_LISTEN_ADDR")},
			Value:   ":8000",
		},

		&cli.BoolFlag{
			Name:    "logs",
			Usage:   "whether to enable access logs or not",
			EnvVars: []string{prefixWrapper("ENABLE_ACCESS_LOGS")},
			Value:   false,
		},

		&cli.StringFlag{
			Name:    "fpm",
			Usage:   "the php-fpm binary filename",
			EnvVars: []string{"PHP_FPM"},
			Value:   "php-fpm",
		},
	}
}
