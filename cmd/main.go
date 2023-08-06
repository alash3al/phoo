package main

import (
	"github.com/alash3al/phoo/cmd/server"
	"github.com/alash3al/phoo/internals/config"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

var (
	globalConfig config.Config
)

func main() {
	app := &cli.App{
		Name:        "phoo",
		Description: "php modern applications server that utilizes the bullet-proof php-fpm under-the-hood",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				EnvVars: []string{"PHOO_CONFIG"},
				Value:   "phoo.yaml",
			},
		},
		Before: func(ctx *cli.Context) error {
			c, err := config.LoadFile(ctx.String("config"))
			if err != nil {
				return err
			}

			globalConfig = *c

			return nil
		},
	}

	app.Commands = append(app.Commands, &cli.Command{
		Name:    "serve",
		Aliases: []string{"s"},
		Action:  server.Serve(&globalConfig),
	})

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err.Error())
	}
}
