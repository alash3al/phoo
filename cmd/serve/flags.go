package serve

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"runtime"
)

func DefaultFlags(envPrefix string) []cli.Flag {
	prefixWrapper := func(k string) string {
		return fmt.Sprintf("%s_%s", envPrefix, k)
	}

	return []cli.Flag{
		&cli.StringFlag{
			Name:     "root",
			Usage:    "the document root full path",
			EnvVars:  []string{prefixWrapper("DOCUMENT_ROOT")},
			Required: true,
		},

		&cli.StringFlag{
			Name:    "entrypoint",
			Usage:   "the default php entrypoint script",
			EnvVars: []string{"ENTRYPOINT"},
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
			EnvVars: []string{prefixWrapper("PHP_FPM")},
			Value:   "php-fpm",
		},

		&cli.StringFlag{
			Name:    "data",
			Usage:   "the directory to store phoo related files in",
			EnvVars: []string{prefixWrapper("DATA_DIR")},
			Value:   "./.phoo",
		},

		&cli.StringSliceFlag{
			Name:    "ini",
			Usage:   "php ini settings in the form of key=value, this flag could be repeated for multiple ini settings",
			EnvVars: []string{prefixWrapper("PHP_INI")},
		},

		&cli.IntFlag{
			Name:    "workers",
			Usage:   "php fpm workers, this is the maximum requests to be served at the same time",
			EnvVars: []string{prefixWrapper("WORKER_COUNT")},
			Value:   runtime.NumCPU(),
		},

		&cli.IntFlag{
			Name:    "requests",
			Usage:   "php fpm max requests per worker, if a worker reached this number, it would be recycled",
			EnvVars: []string{prefixWrapper("WORKER_MAX_REQUEST_COUNT")},
			Value:   runtime.NumCPU() * 100,
		},

		&cli.IntFlag{
			Name:    "timeout",
			Usage:   "php fpm max request time in seconds per worker, if a worker reached this number, it would be terminated, 0 means 'Disabled'",
			EnvVars: []string{prefixWrapper("WORKER_MAX_REQUEST_TIME")},
			Value:   300,
		},

		&cli.StringFlag{
			Name:    "metrics",
			Usage:   "the prometheus metrics endpoint, empty means disabled",
			EnvVars: []string{prefixWrapper("METRICS_PATH")},
		},

		&cli.IntFlag{
			Name:    "gzip",
			Usage:   "gzip level, 0 means 'Disabled",
			EnvVars: []string{prefixWrapper("GZIP_LEVEL")},
			Value:   5,
		},
	}
}
