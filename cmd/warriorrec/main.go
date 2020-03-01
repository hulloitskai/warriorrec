package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/urfave/cli"
	"go.stevenxie.me/gopkg/cmdutil"
	"go.stevenxie.me/warriorrec"
	"go.stevenxie.me/warriorrec/innosoft"
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
			Action:    today,
		},
	}

	if err := app.Run(os.Args); err != nil {
		cmdutil.Fatalf("Error: %+v\n", err)
	}
}

func today(c *cli.Context) error {
	client := innosoft.NewClient(nil)
	schedule, err := client.GetSchedule(context.Background())
	if err != nil {
		return err
	}

	// Compute dayStart and dayEnd.
	var dayStart time.Time
	if c.Args().Present() {
		dayStart, err = time.ParseInLocation("2006-01-02", c.Args().First(), time.Local)
		if err != nil {
			return errors.Wrap(err, "parse date")
		}
	} else {
		now := time.Now()
		dayStart = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	}
	dayEnd := dayStart.Add(24 * time.Hour)

	// Initialize tabwriter.
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintf(w, "START\tEND\tNAME\tCATEGORY\tLOCATION\n")

	// Oragnize categories by ID.
	categories := make(map[string]*warriorrec.ActivityCategory)
	for _, c := range schedule.Categories {
		categories[c.ID] = c
	}

	var activities []*warriorrec.Activity
	for _, a := range schedule.Activities {
		if a.Start.Before(dayStart) || a.End.After(dayEnd) {
			continue
		}
		activities = append(activities, a)
	}

	// Sort activities by start and end time.
	sort.Slice(activities, func(i, j int) bool {
		var (
			iStart = activities[i].Start
			jStart = activities[j].Start
		)
		if iStart.Equal(jStart) {
			return activities[i].End.Before(activities[j].End)
		}
		return iStart.Before(jStart)
	})

	// Print activities.
	const format = "3:04 PM"
	for _, a := range activities {
		// Shorten common naming redundancies.
		name := strings.TrimPrefix(a.Name, "Winter 2020 - ")
		name = strings.TrimPrefix(name, "CIF Fitness Centre Hours - ")
		name = strings.TrimPrefix(name, "Swimming - Fitness and Rec Swim - ")
		name = strings.TrimPrefix(name, "Open Rec Warrior Zone High Performance Center - ")

		// Fix common misspellings.
		if strings.HasSuffix(name, "Performanc") {
			name = strings.TrimSuffix(name, "Performanc") + "Performance)"
		}

		// Trim name.
		if len(name) > 50 {
			name = name[:47] + "..."
		}

		fmt.Fprintf(
			w,
			"%s\t%s\t%s\t%s\t%s\n",
			a.Start.Format(format),
			a.End.Format(format),
			name,
			categories[a.CategoryID].Name,
			a.Location,
		)
	}
	return w.Flush()
}

func categories(c *cli.Context) error {
	client := innosoft.NewClient(nil)
	schedule, err := client.GetSchedule(context.Background())
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	urls := c.Bool("include-urls")
	if urls {
		fmt.Fprint(w, "ID\tNAME\tURL\n")
	} else {
		fmt.Fprint(w, "ID\tNAME\n")
	}

	for _, c := range schedule.Categories {
		if urls {
			fmt.Fprintf(w, "%s\t%s\t%s\n", c.ID, c.Name, c.URL)
		} else {
			fmt.Fprintf(w, "%s\t%s\n", c.ID, c.Name)
		}
	}
	return w.Flush()
}
