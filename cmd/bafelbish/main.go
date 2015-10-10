package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/moul/bafelbish"
)

func main() {
	app := cli.NewApp()
	app.Name = "bafelbish"
	app.Email = "https://github.com/moul"
	app.Usage = "Translate YAML, JSON, TOML, ..."
	app.Action = Action

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "input, input-format, i",
			Value: "json",
			Usage: "input format: yaml, json, toml",
		},
		cli.StringFlag{
			Name:  "output, output-format, o",
			Value: "json",
			Usage: "output format: yaml, json, toml",
		},
	}

	app.Run(os.Args)
}

func Action(c *cli.Context) {
	fish := bafelbish.NewFish()

	fish.SetInputFormat(c.String("input-format"))
	fish.SetOutputFormat(c.String("output-format"))

	if len(c.Args()) > 0 {
		// open file
	} else {
		if err := fish.Parse(os.Stdin, os.Stdout); err != nil {
			logrus.Fatalf("fish.Parse error: %v", err)
		}
	}
}
