package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/valyala/fasttemplate"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var (
	flagHttpListenAddr       = flag.String("listen", ":8000", "the http server listen address in the form of [host]:port")
	flagEnableLogging        = flag.Bool("logging", false, "whether to enable logging or not")
	flagGzipLevel            = flag.Int("gzip", 0, "gzip level, 0 means disable")
	flagDocumentRoot         = flag.String("docroot", "./", "the webserver document root")
	flagRouterFileName       = flag.String("script", "./index.php", "the default script filename")
	flagFPMCommand           = flag.String("fpm.cmd", "php-fpm", "the php-fpm command")
	flagFPMWorkerCount       = flag.Int("fpm.worker.count", runtime.NumCPU(), "the fpm workers count")
	flagFPMWorkerMaxRequests = flag.Int("fpm.worker.max_requests", 500, "the maximum number of requests that a single worker will restart after reaching it")
	flagFPMRequestTimeout    = flag.Int("fpm.request.timeout", -1, "the request timeout")
	flagPHPINIVars           = flag.String("php.ini", "", "semicolon delimited string of ini key=value pairs")
)

var (
	//go:embed php-fpm.conf
	fpmConfTemplate string

	fpmPIDFilename    string
	fpmConfigFilename string
	fpmSocketFilename string

	server = echo.New()
)

func init() {
	flag.Parse()

	server.HideBanner = true

	for _, p := range []*string{flagDocumentRoot, flagRouterFileName} {
		if abs, err := filepath.Abs(*p); err != nil {
			log.Fatal(err.Error())
		} else {
			*p = abs
		}
	}

	_ = os.RemoveAll(".phoo")

	if err := os.MkdirAll(".phoo/fpm", 0775); err != nil && !os.IsExist(err) {
		log.Fatal(err.Error())
	}

	phooDir, err := filepath.Abs(".phoo")
	if err != nil {
		server.Logger.Fatal(err.Error())
	}

	fpmConfigFilename = filepath.Join(phooDir, "fpm", "config")
	fpmPIDFilename = filepath.Join(phooDir, "fpm", "pid")
	fpmSocketFilename = filepath.Join(phooDir, "fpm", "sock")

	fpmConfTemplate = fasttemplate.ExecuteString(fpmConfTemplate, "{{", "}}", map[string]any{
		"files.pid":           fpmPIDFilename,
		"files.socket":        fpmSocketFilename,
		"worker.count":        fmt.Sprintf("%d", *flagFPMWorkerCount),
		"worker.max_requests": fmt.Sprintf("%d", *flagFPMWorkerMaxRequests),
		"timeout":             fmt.Sprintf("%d", *flagFPMRequestTimeout),
	})

	if err := os.WriteFile(fpmConfigFilename, []byte(fpmConfTemplate), 0775); err != nil {
		server.Logger.Fatal(err.Error())
	}
}
