package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/mbarrin/gwarr/internal/pkg/cache"
	"github.com/mbarrin/gwarr/internal/pkg/data"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

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
	OK    bool   `json:"ok,omitempty"`
	TS    string `json:"ts,omitempty"`
	Error string `json:"error,omitempty"`
}

// Client defines a slack client and associated cache
type Client struct {
	url     string
	channel string
	token   string
	client  http.Client
	redis   *redis.Client
}

// New creates a new Slack client
func New(channel string, token string, redisAddr string) (*Client, error) {
	cache, err := cache.New(redisAddr)
	if err != nil {
		return nil, err
	}

	sc := Client{
		url:     "https://slack.com/api/",
		channel: channel,
		token:   "Bearer " + token,
		client:  *http.DefaultClient,
		redis:   cache,
	}

	slog.With("package", "slack").Info("Slack client initialised")
	return &sc, nil
}

func (sc *Client) newRequest(b []byte, m string) *http.Request {
	r, _ := http.NewRequest(http.MethodPost, sc.url+m, bytes.NewBuffer(b))

	r.Header.Add("Content-Type", "application/json; charset=utf-8")
	r.Header.Add("Authorization", sc.token)

	return r
}

// Post posts a Radarr webhook formatted to a Slack message
func (sc *Client) Post(d data.Data) error {
	ts, err := sc.redis.HGet(ctx, d.Service(), fmt.Sprint(d.ID())).Result()
	if err != nil {
		slog.Debug(fmt.Sprintf("Could not find timestamp for ID: %d for %s", d.ID(), d.Service()))
	}
	slog.Debug(ts)

	var r *http.Request
	switch d.Type() {
	case "MovieAdded":
		b := onAddInfo(sc.channel, d, ts)
		jb, err := json.Marshal(b)
		if err != nil {
			return err
		}

		if ts == "" {
			r = sc.newRequest(jb, "chat.postMessage")
		} else {
			r = sc.newRequest(jb, "chat.update")
		}
	case "Grab":
		b := onGrabInfo(sc.channel, d, ts)
		jb, err := json.Marshal(b)
		if err != nil {
			return err
		}

		if ts == "" {
			r = sc.newRequest(jb, "chat.postMessage")
		} else {
			r = sc.newRequest(jb, "chat.update")
		}
	case "Download":
		b := onDownloadInfo(sc.channel, d, ts)
		jb, err := json.Marshal(b)
		if err != nil {
			return err
		}

		if ts == "" {
			r = sc.newRequest(jb, "chat.postMessage")
		} else {
			r = sc.newRequest(jb, "chat.update")
		}
	case "MovieDelete":
		b := onDeleteInfo(sc.channel, d)
		jb, err := json.Marshal(b)
		if err != nil {
			return err
		}

		r = sc.newRequest(jb, "chat.postMessage")
	default:
		b := unhandled(sc.channel, d)
		jb, err := json.Marshal(b)
		if err != nil {
			return err
		}

		r = sc.newRequest(jb, "chat.postMessage")
	}

	resp, err := sc.client.Do(r)
	if err != nil {
		return err
	}
	defer func() {
		err := r.Body.Close()
		if err != nil {
			slog.With("package", "slack").Error("Failed to close body")
		}
	}()

	body, _ := io.ReadAll(resp.Body)

	response := response{}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return errors.New("Message sent, but response could not be decoded. Err: " + err.Error())
	}

	if response.OK {
		if d.Type() == "MovieDelete" || d.Type() == "Download" {
			err := sc.redis.HDel(ctx, d.Service(), fmt.Sprint(d.ID())).Err()
			if err != nil {
				slog.Error(err.Error())
			}
		} else {
			err := sc.redis.HSet(ctx, d.Service(), d.ID(), response.TS).Err()
			if err != nil {
				slog.Error(err.Error())
			}
		}
	} else {
		slog.Error(response.Error)
		slog.Debug(sc.token)
	}

	return nil
}

func onGrabInfo(c string, d data.Data, ts string) body {
	b := base(c, d)
	b.TS = ts
	b.Blocks[0].Text.Text = fmt.Sprintf(":large_orange_circle: Grabbed: %s", d.Title())
	b.Blocks = append(b.Blocks,
		block{
			Type: "section",
			Fields: &[]text{
				{Type: "mrkdwn", Text: "*Quality:*\n" + d.Quality()},
				{Type: "mrkdwn", Text: "*Release Group:*\n" + d.ReleaseGroup()},
			},
		},
	)
	return b
}

func onDownloadInfo(c string, d data.Data, ts string) body {
	b := base(c, d)
	b.TS = ts
	b.Blocks[0].Text.Text = fmt.Sprintf(":large_green_circle: Downloaded: %s", d.Title())
	b.Blocks = append(b.Blocks,
		block{
			Type: "section",
			Fields: &[]text{
				{Type: "mrkdwn", Text: "*Quality:*\n" + d.Quality()},
				{Type: "mrkdwn", Text: "*Release Group:*\n" + d.ReleaseGroup()},
			},
		},
	)
	return b
}

func onAddInfo(c string, d data.Data, ts string) body {
	b := base(c, d)
	b.TS = ts
	b.Blocks[0].Text.Text = fmt.Sprintf(":large_green_circle: Added: %s", d.Title())
	return b
}

func onDeleteInfo(c string, d data.Data) body {
	b := base(c, d)
	b.Blocks[0].Text.Text = fmt.Sprintf(":red_circle: Delete: %s", d.Title())
	return b
}

func unhandled(c string, d data.Data) body {
	unhandledData, _ := json.Marshal(d)
	return body{
		Channel: c,
		Blocks: []block{
			{
				Type: "header",
				Text: &text{Type: "plain_text", Text: "unhandled", Emoji: true},
			},
			{
				Type: "section",
				Text: &text{
					Type: "plain_text",
					Text: string(unhandledData),
				},
			},
		},
	}
}

func base(c string, d data.Data) body {
	return body{
		Channel: c,
		Blocks: []block{
			{
				Type: "header",
				Text: &text{Type: "plain_text", Emoji: true},
			},
			{
				Type: "section",
				Text: &text{
					Type: "mrkdwn",
					Text: d.URL(),
				},
			},
			{
				Type: "section",
				Fields: &[]text{
					{Type: "mrkdwn", Text: "*Release Date:*\n" + d.ReleaseDate()},
					{Type: "mrkdwn", Text: "*IMDB:*\nhttps://imdb.com/title/" + d.IMDBID()},
				},
			},
		},
	}
}
