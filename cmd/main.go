package main

import (
	"github.com/alash3al/phoo/cmd/serve"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)
import _ "github.com/joho/godotenv/autoload"

func main() {
	app := &cli.App{
		Name:                 "phoo",
		Description:          "php modern applications serve that utilizes the bullet-proof php-fpm under-the-hood",
		Suggest:              true,
		Version:              "v3.x",
		EnableBashCompletion: true,
		SliceFlagSeparator:   ";",
	}

	app.Commands = append(app.Commands, &cli.Command{
		Name:   "serve",
		Flags:  serve.DefaultFlags("PHOO"),
		Before: serve.Before(),
		Action: serve.Action(),
	})

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err.Error())
	}
}
