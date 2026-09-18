package calendar

import (
	"cocon/internal/models"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	ical "github.com/arran4/golang-ical"
)

type Client struct {
	EventsURL string

	mu        sync.Mutex
	cached    []models.Event
	fetchedAt time.Time
	ttl       time.Duration
}

func NewClient(eventsURL string) *Client {
	return &Client{
		EventsURL: eventsURL,
		ttl:       15 * time.Minute,
	}
}

func (c *Client) EventsBetween(start, end time.Time) ([]models.Event, error) {
	if c.EventsURL == "" {
		return nil, nil
	}

	events, err := c.fetch()
	if err != nil {
		return nil, err
	}

	var result []models.Event
	for _, ev := range events {
		if !ev.Start.Before(start) && !ev.Start.After(end) {
			result = append(result, ev)
		}
	}
	return result, nil
}

func (c *Client) fetch() ([]models.Event, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cached != nil && time.Since(c.fetchedAt) < c.ttl {
		return c.cached, nil
	}

	resp, err := http.Get(c.EventsURL)
	if err != nil {
		return nil, fmt.Errorf("fetch ics: %w", err)
	}
	defer resp.Body.Close()

	cal, err := ical.ParseCalendar(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parse ics: %w", err)
	}

	paris, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		return nil, fmt.Errorf("load timezone Europe/Paris: %w", err)
	}

	var events []models.Event

	for _, comp := range cal.Events() {
		summary := comp.GetProperty(ical.ComponentPropertySummary)

		start, errStart := comp.GetStartAt()
		if summary == nil || errStart != nil {
			continue
		}

		end, err := comp.GetEndAt()
		if err != nil {
			continue
		}

		events = append(events, models.Event{
			Title:  summary.Value,
			Start:  start.In(paris),
			End:    end.In(paris),
			Author: eventAuthor(comp),
		})
	}

	c.cached = events
	c.fetchedAt = time.Now()

	return events, nil
}

func eventAuthor(comp *ical.VEvent) string {
	organizer := comp.GetProperty(ical.ComponentPropertyOrganizer)
	if organizer == nil {
		return ""
	}
	if cn := organizer.ICalParameters["CN"]; len(cn) > 0 {
		return cn[0]
	}
	return strings.TrimPrefix(organizer.Value, "mailto:")
}
