package main

import (
	"github.com/alash3al/phoo/cmd/serve"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name:        "phoo",
		Description: "php modern applications serve that utilizes the bullet-proof php-fpm under-the-hood",
	}

	app.Commands = append(app.Commands, &cli.Command{
		Name:   "serve",
		Flags:  serve.DefaultFlags("PHOO"),
		Action: serve.Action(),
	})

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err.Error())
	}
}
