package innosoft

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cockroachdb/errors"
	"go.stevenxie.me/warriorrec"
)

// A Client can load information from the Innosoft API server.
type Client struct{ c *http.Client }

// NewClient creates a new Client.
func NewClient(c *http.Client) *Client {
	if c == nil {
		c = http.DefaultClient
	}
	return &Client{c: c}
}

// A Schedule contains a list of Activities and Categories.
type Schedule struct {
	Activities []*warriorrec.Activity         `json:"activities"`
	Categories []*warriorrec.ActivityCategory `json:"category"`
}

const (
	baseURL     = "https://innosoftfusiongo.com/schools/school20"
	scheduleURL = baseURL + "/schedule.json"
	timezone    = "America/Toronto"
)

// GetSchedule gets the complete Warriors Recreaction schedule for the current
// school term.
func (c *Client) GetSchedule(ctx context.Context) (*Schedule, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, errors.Wrap(err, "load timezone")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, scheduleURL, nil)
	if err != nil {
		return nil, err
	}
	res, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var data struct {
		Categories []struct {
			ID   string `json:"id"`
			Name string `json:"category"`
			Days []struct {
				Date       string `json:"date"`
				Activities []struct {
					ID          string `json:"activityID"`
					DetailURL   string `json:"detailUrl"`
					Name        string `json:"activity"`
					Description string `json:"description"`
					Location    string `json:"location"`
					Start       string `json:"startTime"`
					End         string `json:"endTime"`
					Cancelled   string `json:"isCancelled"`
					Spots       int    `json:"availableSpots"`
				} `json:"scheduled_activities"`
			} `json:"days"`
		} `json:"categories"`
	}
	if err = json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, errors.Wrap(err, "decode body")
	}

	schedule := new(Schedule)
	for _, c := range data.Categories {
		category := &warriorrec.ActivityCategory{
			ID:   c.ID,
			Name: c.Name,
		}
		for _, day := range c.Days {
			for i, a := range day.Activities {
				if i == 0 {
					category.URL = a.DetailURL
				}

				parse := func(t string) (time.Time, error) {
					value := fmt.Sprintf("%sT%s", day.Date, t)
					return time.ParseInLocation("2006-01-02T15:04:05", value, loc)
				}

				start, err := parse(a.Start)
				if err != nil {
					return nil, errors.Wrap(err, "parse start time")
				}
				end, err := parse(a.End)
				if err != nil {
					return nil, errors.Wrap(err, "parse end time")
				}
				if end.Hour() == 0 {
					end = end.AddDate(0, 0, 1)
				}

				schedule.Activities = append(schedule.Activities, &warriorrec.Activity{
					ID:         a.ID,
					CategoryID: c.ID,

					Name:        a.Name,
					Description: a.Description,
					Location:    a.Location,
					Spots:       a.Spots,

					Start:     start,
					End:       end,
					Cancelled: a.Cancelled != "0",
				})
			}
		}
		schedule.Categories = append(schedule.Categories, category)
	}

	return schedule, nil
}
