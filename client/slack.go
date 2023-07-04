package client

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/sixt/hardcoded-slack-webhook-alerter/model"
)

type Data struct {
	Channel string
	Repo    string
	URL     string
	Message string
}

// SlackClient is used to send messages to slack
type SlackClient struct {
	httpClient *http.Client
	dryRun     bool
	payload    *template.Template
	message    string
}

// NewSlackClient returns a new slack client
func NewSlackClient(dryRun bool, message string) *SlackClient {
	return &SlackClient{
		httpClient: &http.Client{},
		dryRun:     dryRun,
		payload:    template.Must(template.New("payload").Parse(payload)),
		message:    message,
	}
}

// SendMessage sends a message to the given webhook
func (s *SlackClient) SendMessage(result model.Result) error {
	var err error

	if len(result.Channels) == 0 {
		err = s.send(result.Webhook, Data{Channel: "", Repo: result.Repo, URL: result.URL})
	} else {
		//for each slack channel we found, send a notification
		for channel := range result.Channels {
			err = s.send(result.Webhook, Data{Channel: channel, Repo: result.Repo, URL: result.URL})
		}
	}

	return err
}

func (s *SlackClient) send(webhook string, data Data) error {
	if data.Channel == "" {
		log.Printf("Sending request to hook %s for %s\n", webhook, data.Repo)
	} else {
		log.Printf("Sending request to %s for hook %s for %s\n", data.Channel, webhook, data.Repo)
	}

	if s.dryRun {
		return nil
	}

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	s.payload.Execute(w, Data{Channel: data.Channel, Repo: data.Repo, URL: data.URL, Message: s.message})
	w.Flush()

	req, _ := http.NewRequest(http.MethodPost, webhook, bufio.NewReader(&b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error sending message to slack: %s", resp.Status)
	}
	return nil
}

const payload = `{
	{{if .Channel}}"channel": "{{.Channel}}",{{end}}
	"attachments":[
	   {
		  "fallback":"<!channel> This message was sent using a hardcoded webhook found in {{.Repo}}",
		  "pretext":"<!channel>",
		  "title":"[Alert] HARD-CODED SLACK WEBHOOK",
		  "title_link":"https://www.sixt.tech/slack-webhook-security",
		  "text":"This message was sent using a hardcoded webhook found in {{.Repo}}: {{.URL}}\n\n{{.Message}}",
		  "color":"#d63232",
	   }
	]
 }`
