package main

import (
	"os"

	"github.com/urfave/cli"
	"go.stevenxie.me/gopkg/cmdutil"
	"go.stevenxie.me/warriorrec/internal"
)

const name = "warriorrec"

func main() {
	app := cli.NewApp()
	app.Name = name
	app.Usage = "Check the Warrior Recreation activity schedule."
	// app.UsageText = fmt.Sprintf("%s [command]", name)
	app.Version = internal.Version
	app.Commands = []cli.Command{
		{
			Name:  "categories",
			Usage: "List available activity categories.",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "include-urls",
					Usage: "Include reference URLs.",
				},
			},
			Action: categories,
		},
		{
			Name:      "schedule",
			Usage:     "List activity schedule for a date (format: YYYY-MM-DD).",
			ArgsUsage: "[<date>]",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "include-cancelled",
					Usage: "Include cancelled activities.",
				},
			},
			Action: schedule,
		},
	}

	if err := app.Run(os.Args); err != nil {
		cmdutil.Fatalf("Error: %+v\n", err)
	}
}
