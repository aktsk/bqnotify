package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/slack-go/slack"
)

// Slack has configurations for notifying to Slack
// Fields that have `json:"-"` are not included in the payload to Slack
type Slack struct {
	URL             string `json:"-"`
	Title           string `json:"title"`
	Color           string `json:"color"`
	UploadChannelID string `json:"-" yaml:"upload_channel_id"`
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

func (s Slack) Upload(csvBuffer *bytes.Buffer) error {
	api := slack.New(os.Getenv("SLACK_BOT_TOKEN"))

	_, err := api.UploadFileV2Context(context.Background(), slack.UploadFileV2Parameters{
		FileSize: len(csvBuffer.String()),
		Reader:   csvBuffer,
		Filename: fmt.Sprintf("%s.csv", s.Title),
		Title:    s.Title,
		Channel:  s.UploadChannelID,
	})

	return err
}
