package main

import (
	"github.com/labstack/echo/v4"
	"github.com/yookoala/gofast"
	"os"
	"path"
	"strings"
)

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
