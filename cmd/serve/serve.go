package serve

import (
	"errors"
	"github.com/alash3al/phoo/internals/fpm"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
)

func Action() cli.ActionFunc {
	return func(ctx *cli.Context) error {
		if err := os.RemoveAll(ctx.String("data")); err != nil {
			return err
		}

		if err := os.MkdirAll(ctx.String("data"), 0755); err != nil {
			return err
		}

		if ctx.String("entrypoint") == "" {
			entrypoints := []string{
				filepath.Join(ctx.String("root"), "index.php"),
				filepath.Join(ctx.String("root"), "app.php"),
				filepath.Join(ctx.String("root"), "server.php"),
			}

			detectedEntrypoint := ""

			for _, entrypoint := range entrypoints {
				stat, err := os.Stat(entrypoint)
				if err != nil {
					continue
				}

				if stat.IsDir() {
					continue
				}

				detectedEntrypoint = entrypoint
			}

			if detectedEntrypoint == "" {
				return errors.New("unable to auto-detect the entrypoint script, try to put it yourself")
			}

			if err := ctx.Set("entrypoint", detectedEntrypoint); err != nil {
				return err
			}
		}

		fpmProcess := &fpm.Process{
			BinFilename:           ctx.String("fpm"),
			PIDFilename:           filepath.Join(ctx.String("data"), "fpm.pid"),
			ConfigFilename:        filepath.Join(ctx.String("data"), "fpm.ini"),
			SocketFilename:        filepath.Join(ctx.String("data"), "fpm.sock"),
			INI:                   ctx.StringSlice("ini"),
			WorkerCount:           ctx.Int("workers"),
			WorkerMaxRequestCount: ctx.Int("requests"),
			WorkerMaxRequestTime:  ctx.Int("timeout"),
		}

		if err := fpmProcess.Start(); err != nil {
			return err
		}

		e := echo.New()
		e.HideBanner = true

		if ctx.Bool("logs") {
			e.Use(middleware.Logger())
		}

		e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
			Level: ctx.Int("gzip"),
		}))

		e.Use(middleware.Recover())
		e.Use(serveStaticFilesOnlyMiddleware(ctx.String("root")))
		e.Use(serveFastCGIMiddleware(
			ctx.String("entrypoint"),
			"unix",
			fpmProcess.SocketFilename,
		))

		return e.Start(ctx.String("http"))
	}
}
