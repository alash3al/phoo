package main

import (
	"github.com/alash3al/xcgi/cmd/serve"
	"github.com/alash3al/xcgi/pkg/symbols"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = symbols.AppName
	app.Version = symbols.AppVersion
	app.Before = func(ctxCli *cli.Context) error {
		filename := ctxCli.String(symbols.FlagNameEnvFilename)
		if len(filename) < 1 {
			return nil
		}
		return godotenv.Load(filename)
	}

	app.Flags = append(app.Flags, &cli.StringFlag{
		Name:  symbols.FlagNameEnvFilename,
		Usage: "if provided, the configuration values will be loaded from it",
	})

	app.Commands = append(app.Commands, serve.Command())

	if err := app.Run(os.Args); err != nil {
		log.Error(err.Error())
	}
}
