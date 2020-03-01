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

	"go.stevenxie.me/warriorrec"
	"go.stevenxie.me/warriorrec/innosoft"
)

func schedule(c *cli.Context) error {
	includeCancelled := c.Bool("include-cancelled")

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
		if a.Cancelled && !includeCancelled {
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
	const layout = "3:04 PM"
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

		// Print raw text.
		text := fmt.Sprintf(
			"%s\t%s\t%s\t%s\t%s",
			a.Start.Format(layout),
			a.End.Format(layout),
			name,
			categories[a.CategoryID].Name,
			a.Location,
		)
		fmt.Fprintf(w, "%s\n", text)
	}
	return w.Flush()
}
