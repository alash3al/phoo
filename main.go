package main

import (
	"github.com/labstack/echo/v4/middleware"
	"log"
)

func main() {
	if *flagEnableLogging {
		server.Use(middleware.Logger())
	}

	server.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: *flagGzipLevel,
	}))

	server.Use(middleware.Recover())
	server.Use(serveStaticFilesOnlyMiddleware(*flagDocumentRoot))
	server.Use(serveFastCGIMiddleware(
		*flagRouterFileName,
		"unix",
		fpmSocketFilename,
	))

	if err := startPHPFPM(); err != nil {
		server.Logger.Fatal(err.Error())
	}

	log.Fatal(server.Start(*flagHttpListenAddr))
}
