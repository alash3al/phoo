package serve

import (
	"fmt"
	"github.com/labstack/echo/v4/middleware"
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
			Aliases:  []string{"r"},
			EnvVars:  []string{prefixWrapper("DOCUMENT_ROOT")},
			Required: true,
		},

		&cli.StringFlag{
			Name:    "entrypoint",
			Usage:   "the default php entrypoint script",
			Aliases: []string{"e"},
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
			Aliases: []string{"i"},
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

		&cli.BoolFlag{
			Name:    "cors",
			Usage:   "whether to enable/disable the cors-* features/flags",
			EnvVars: []string{prefixWrapper("CORS_ENABLED")},
			Value:   false,
		},

		&cli.StringSliceFlag{
			Name:    "cors-origin",
			Usage:   "this flag adds the specified origin to the list of allowed cors origins, you can call it multiple times to add multiple origins",
			EnvVars: []string{prefixWrapper("CORS_ORIGINS")},
			Value:   cli.NewStringSlice("*"),
		},

		&cli.StringSliceFlag{
			Name:    "cors-methods",
			Usage:   "this flag adds the specified methods to the list of allowed cors methods, you can call it multiple times to add multiple methods",
			EnvVars: []string{prefixWrapper("CORS_METHODS")},
			Value:   cli.NewStringSlice(middleware.DefaultCORSConfig.AllowMethods...),
		},

		&cli.StringSliceFlag{
			Name:    "cors-headers",
			Usage:   "this flag adds the specified headers to the list of allowed cors headers the client can send, you can call it multiple times to add multiple headers",
			EnvVars: []string{prefixWrapper("CORS_HEADERS")},
		},

		&cli.StringSliceFlag{
			Name:    "cors-expose",
			Usage:   "this flag adds the specified headers to the list of allowed headers the client can access, you can call it multiple times to add multiple headers",
			EnvVars: []string{prefixWrapper("CORS_EXPOSE")},
		},

		&cli.BoolFlag{
			Name:    "cors-credentials",
			Usage:   "this flag indicates whether or not the actual request can be made using credentials",
			EnvVars: []string{prefixWrapper("CORS_CREDENTIALS")},
			Value:   false,
		},

		&cli.IntFlag{
			Name:    "cors-age",
			Usage:   "the cors max_age in seconds",
			EnvVars: []string{prefixWrapper("CORS_AGE")},
			Value:   0,
		},

		&cli.StringFlag{
			Name:    "user",
			Usage:   "run the fpm www pool user",
			EnvVars: []string{prefixWrapper("USER")},
			Value:   "www-data",
		},

		&cli.StringFlag{
			Name:    "group",
			Usage:   "run the fpm www pool group",
			EnvVars: []string{prefixWrapper("GROUP")},
			Value:   "www-data",
		},
	}
}
