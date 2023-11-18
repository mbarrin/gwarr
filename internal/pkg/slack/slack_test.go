package slack

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mbarrin/gwarr/internal/pkg/radarr"
	"github.com/mbarrin/gwarr/internal/pkg/state"
)

var radarrOnGrab = radarr.Data{
	Movie: radarr.Movie{
		Title:       "Film",
		Year:        1970,
		ReleaseDate: "1970-01-01",
		IMDBID:      "tt8415836",
		TMDBID:      55,
	},
	Release: &radarr.Release{
		Quality:      "1080p",
		ReleaseGroup: "legit",
	},
	EventType:      "Grab",
	ApplicationURL: "http://localhost",
}

var radarrOnDownload = radarr.Data{
	Movie: radarr.Movie{
		Title:       "Film",
		Year:        1970,
		ReleaseDate: "1970-01-01",
		IMDBID:      "tt8415836",
		TMDBID:      55,
	},
	MovieFile: &radarr.MovieFile{
		Quality:      "1080p",
		ReleaseGroup: "legit",
	},
	EventType:      "Grab",
	ApplicationURL: "http://localhost",
}

var slackOnGrab = body{
	TS:      "",
	Channel: "c123",
	Blocks: []block{
		{
			Type: "header",
			Text: &text{
				Type:  "plain_text",
				Text:  ":large_orange_circle: Grabbed: Film (1970)",
				Emoji: true,
			},
		},
		{
			Type: "section",
			Text: &text{
				Type: "mrkdwn",
				Text: "http://localhost/movie/55",
			},
		},
		{
			Type: "section",
			Fields: &[]text{
				{
					Type: "mrkdwn",
					Text: "*Release Date:*\n1970-01-01",
				},
				{
					Type: "mrkdwn",
					Text: "*IMDB:*\nhttps://imdb.com/title/tt8415836",
				},
			},
		},
		{
			Type: "section",
			Fields: &[]text{
				{
					Type: "mrkdwn",
					Text: "*Quality:*\n1080p",
				},
				{
					Type: "mrkdwn",
					Text: "*Release Group:*\nlegit",
				},
			},
		},
	},
}

var slackOnDownload = body{
	TS:      "123",
	Channel: "c123",
	Blocks: []block{
		{
			Type: "header",
			Text: &text{
				Type:  "plain_text",
				Text:  ":large_green_circle: Downloaded: Film (1970)",
				Emoji: true,
			},
		},
		{
			Type: "section",
			Text: &text{
				Type: "mrkdwn",
				Text: "http://localhost/movie/55",
			},
		},
		{
			Type: "section",
			Fields: &[]text{
				{
					Type: "mrkdwn",
					Text: "*Release Date:*\n1970-01-01",
				},
				{
					Type: "mrkdwn",
					Text: "*IMDB:*\nhttps://imdb.com/title/tt8415836",
				},
			},
		},
		{
			Type: "section",
			Fields: &[]text{
				{
					Type: "mrkdwn",
					Text: "*Quality:*\n1080p",
				},
				{
					Type: "mrkdwn",
					Text: "*Release Group:*\nlegit",
				},
			},
		},
	},
}

var slackUpdateOnGrab = body{
	TS:      "1234",
	Channel: "c123",
	Blocks: []block{
		{
			Type: "header",
			Text: &text{
				Type:  "plain_text",
				Text:  ":large_orange_circle: Grabbed: Film (1970)",
				Emoji: true,
			},
		},
		{
			Type: "section",
			Text: &text{
				Type: "mrkdwn",
				Text: "http://localhost/movie/55",
			},
		},
		{
			Type: "section",
			Fields: &[]text{
				{
					Type: "mrkdwn",
					Text: "*Release Date:*\n1970-01-01",
				},
				{
					Type: "mrkdwn",
					Text: "*IMDB:*\nhttps://imdb.com/title/tt8415836",
				},
			},
		},
		{
			Type: "section",
			Fields: &[]text{
				{
					Type: "mrkdwn",
					Text: "*Quality:*\n1080p",
				},
				{
					Type: "mrkdwn",
					Text: "*Release Group:*\nlegit",
				},
			},
		},
	},
}

func TestNew(t *testing.T) {
	tests := map[string]struct {
		channel  string
		token    string
		expected Client
	}{
		"creates client": {
			channel: "c1234",
			token:   "xoxb-123",
			expected: Client{
				url:     "https://slack.com/api/",
				channel: "c1234",
				token:   "xoxb-123",
				client:  *http.DefaultClient,
				cache:   state.New(""),
			},
		},
	}

	for _, tc := range tests {
		actual := New("c1234", "xoxb-123", "")
		assert.Equal(t, tc.expected, actual)
	}
}

func TestOnGrabBody(t *testing.T) {
	tests := map[string]struct {
		channel  string
		data     radarr.Data
		update   bool
		expected body
	}{
		"new grab": {
			channel:  "c123",
			data:     radarrOnGrab,
			update:   false,
			expected: slackOnGrab,
		},
		"updated grab": {
			channel:  "c123",
			data:     radarrOnGrab,
			update:   true,
			expected: slackUpdateOnGrab,
		},
	}

	for _, tc := range tests {
		ts := ""
		if tc.update {
			ts = "1234"
		}
		actual := RadarrOnGrabBody(tc.channel, tc.data, ts)
		assert.Equal(t, tc.expected, actual)
	}
}

func TestOnDownloadBody(t *testing.T) {
	tests := map[string]struct {
		channel  string
		data     radarr.Data
		ts       string
		expected body
	}{
		"new download": {
			channel:  "c123",
			data:     radarrOnDownload,
			ts:       "123",
			expected: slackOnDownload,
		},
	}

	for _, tc := range tests {
		actual := RadarrOnDownloadBody(tc.channel, tc.data, tc.ts)
		assert.Equal(t, tc.expected, actual)
	}
}
