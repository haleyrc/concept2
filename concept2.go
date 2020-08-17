package concept2

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const Version = "0.0.1"

const DevelopmentURL = "https://log-dev.concept2.com"

const AcceptHeader = "application/vnd.c2logbook.v1+json"
const ContentTypeHeader = "application/json"

const DefaultPMNumber int64 = 5

type Response struct {
	Data json.RawMessage `json:"data"`
	Meta struct {
		Total       int64 `json:"total"`
		Count       int64 `json:"count"`
		PerPage     int64 `json:"per_page"`
		CurrentPage int64 `json:"current_page"`
		TotalPages  int64 `json:"total_pages"`
		Links       struct {
			Next     string `json:"next"`
			Previous string `json:"prev"`
		} `json:"links"`
	} `json:"meta"`
}

func NewClient() *Client {
	return &Client{
		PMNumber: DefaultPMNumber,
		c:        &http.Client{Timeout: 5 * time.Second},
	}
}

type Client struct {
	PMNumber int64

	c     *http.Client
	token string
}

type Result struct {
	ID            int64         `json:"id"`
	UserID        int64         `json:"user_id"`
	Date          time.Time     `json:"date"`
	Timezone      string        `json:"timezone"`
	DateUTC       time.Time     `json:"date_utc"`
	Distance      int64         `json:"distance"`
	Type          string        `json:"type"`
	Time          int64         `json:"time"`
	TimeFormatted time.Duration `json:"time_formatted"`
	WorkoutType   string        `json:"workout_type"`
	Source        string        `json:"source"`
	WeightClass   string        `json:"weight_class"`
	Verified      bool          `json:"verified"`
	Ranked        bool          `json:"ranked"`
	Comments      string        `json:"comments"`
}

type ResultQuery struct {
	From         string
	To           string
	Type         string
	UpdatedAfter string
}

func (c *Client) Results(ctx context.Context, userID string, query ResultQuery) ([]Result, error) {
	path := "/api/users/%s/results"
	req, err := c.makeRequest(c.url(path, userID), http.MethodGet)
	if err != nil {
		return nil, fmt.Errorf("Client.Results: %w", err)
	}

	params := url.Values{}
	if query.From != "" {
		params.Set("from", query.From)
	}
	if query.To != "" {
		params.Set("to", query.To)
	}
	if query.Type != "" {
		params.Set("type", query.Type)
	}
	if query.UpdatedAfter != "" {
		params.Set("updated_after", query.UpdatedAfter)
	}

	req.URL.RawQuery = params.Encode()

	var data []Result
	if err := c.do(req, &data); err != nil {
		return nil, fmt.Errorf("Client.Results: %w", err)
	}

	return data, nil
}

func (c *Client) do(r *http.Request, dest interface{}) error {
	resp, err := c.c.Do(r)
	if err != nil {
		return fmt.Errorf("Client.do: %w", err)
	}

	type response struct {
	}
}

// TODO (RCH): Make this use the production url if dev mode not true
func (c *Client) url(path string, args ...interface{}) string {
	path = fmt.Sprintf(path, args...)
	return DevelopmentURL + path
}

func (c *Client) makeRequest(url, method string) (*http.Request, error) {
	req, err := http.NewRequest(url, method, nil)
	if err != nil {
		return nil, fmt.Errorf("Client.makeRequest: %w", err)
	}

	c.applyStandardHeaders(req)

	return req, nil
}

func (c *Client) applyStandardHeaders(r *http.Request) {
	r.Header.Set("Accept", AcceptHeader)
	r.Header.Set("Content-Type", ContentTypeHeader)
	r.Header.Set("X-PM-Version", strconv.FormatInt(c.PMNumber, 10))
	r.Header.Set("X-Client-Version", Version)
	if c.token != "" {
		r.Header.Set("Authorization", "Bearer "+c.token)
	}
}
