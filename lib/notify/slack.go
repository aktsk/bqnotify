package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Slack has configurations for notifying to Slack
type Slack struct {
	URL   string `json:"-"`
	Title string `json:"title"`
	Color string `json:"color"`
}

type payload struct {
	Slack
	Attachments []attachment `json:"attachments"`
}

type attachment struct {
	Text     string   `json:"text"`
	Color    string   `json:"color"`
	MrkdwnIn []string `json:"mrkdwn_in"`
}

// Notify notifies BigQuery query results to Slack
func (s Slack) Notify(message string) error {
	color := s.Color

	a := attachment{
		Text:     fmt.Sprintf("*%s*\n%s", s.Title, message),
		Color:    color,
		MrkdwnIn: []string{"text"},
	}

	p := payload{Slack: s, Attachments: []attachment{a}}
	j, _ := json.Marshal(p)
	buf := bytes.NewBuffer(j)

	resp, err := http.Post(s.URL, "application/json", buf)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
