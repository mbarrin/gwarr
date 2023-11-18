package slack

import (
	"bytes"
	"log/slog"
	"net/http"

	"github.com/mbarrin/gwarr/internal/pkg/state"
)

type body struct {
	Channel string  `json:"channel,omitempty"`
	Text    string  `json:"text,omitempty"`
	TS      string  `json:"ts,omitempty"`
	Blocks  []block `json:"blocks,omitempty"`
}

type block struct {
	Type   string  `json:"type,omitempty"`
	Text   *text   `json:"text,omitempty"`
	Fields *[]text `json:"fields,omitempty"`
}

type text struct {
	Type  string `json:"type,omitempty"`
	Text  string `json:"text,omitempty"`
	Emoji bool   `json:"emoji,omitempty"`
}

type response struct {
	OK bool   `json:"ok,omitempty"`
	TS string `json:"ts,omitempty"`
}

type Client struct {
	url     string
	channel string
	token   string
	client  http.Client
	cache   state.State
}

func New(channel string, token string, cachePath string) Client {
	sc := Client{
		url:     "https://slack.com/api/",
		channel: channel,
		token:   token,
		client:  *http.DefaultClient,
		cache:   state.New(cachePath),
	}

	slog.With("package", "slack").Info("Slack client initialised")
	return sc
}

func (sc *Client) newRequest(b []byte, m string) *http.Request {
	r, _ := http.NewRequest(http.MethodPost, sc.url+m, bytes.NewBuffer(b))

	r.Header.Add("Content-Type", "application/json; charset=utf-8")
	r.Header.Add("Authorization", "Bearer "+sc.token)

	return r
}
