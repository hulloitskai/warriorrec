package main

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/urfave/cli"
	"go.stevenxie.me/warriorrec/innosoft"
)

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
