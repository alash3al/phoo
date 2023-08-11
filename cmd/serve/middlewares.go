package serve

import (
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/yookoala/gofast"
	"os"
	"path"
	"strings"
)

func servePrometheusMetricsMiddleware(metricsPath string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			currentPath := strings.Trim(c.Request().URL.Path, "/")
			metricsPath = strings.Trim(metricsPath, "/")

			if metricsPath == "" {
				return next(c)
			}

			if strings.EqualFold(currentPath, metricsPath) {
				return echoprometheus.NewHandler()(c)
			}

			return next(c)
		}
	}
}

func serveStaticFilesOnlyMiddleware(documentRoot string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			filename := path.Join(documentRoot, path.Clean(c.Request().URL.Path))
			ext := strings.ToLower(path.Ext(filename))

			stat, err := os.Stat(filename)
			if err != nil || stat.IsDir() || (ext == "php") || strings.HasPrefix(path.Base(filename), ".") {
				return next(c)
			}

			return c.File(filename)
		}
	}
}

func serveFastCGIMiddleware(routerFilename, fastcgiServerNetwork, fastcgiServerAddr string) echo.MiddlewareFunc {
	connFactory := gofast.SimpleConnFactory(fastcgiServerNetwork, fastcgiServerAddr)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return echo.WrapHandler(gofast.NewHandler(
			gofast.NewFileEndpoint(routerFilename)(gofast.BasicSession),
			gofast.SimpleClientFactory(connFactory),
		))
	}
}
