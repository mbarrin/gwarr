package slack

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mbarrin/gwarr/internal/pkg/data"
	"github.com/mbarrin/gwarr/internal/pkg/radarr"
	"github.com/mbarrin/gwarr/internal/pkg/sonarr"
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

var slackRadarrOnGrab = body{
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
			Text: &text{Type: "mrkdwn", Text: "http://localhost/movie/55"},
		},
		{
			Type: "section",
			Fields: &[]text{
				{Type: "mrkdwn", Text: "*Release Date:*\n1970-01-01"},
				{Type: "mrkdwn", Text: "*IMDB:*\nhttps://imdb.com/title/tt8415836"},
			},
		},
		{
			Type: "section",
			Fields: &[]text{
				{Type: "mrkdwn", Text: "*Quality:*\n1080p"},
				{Type: "mrkdwn", Text: "*Release Group:*\nlegit"},
			},
		},
	},
}

var slackRadarrUpdateOnGrab = body{
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
			Text: &text{Type: "mrkdwn", Text: "http://localhost/movie/55"},
		},
		{
			Type: "section",
			Fields: &[]text{
				{Type: "mrkdwn", Text: "*Release Date:*\n1970-01-01"},
				{Type: "mrkdwn", Text: "*IMDB:*\nhttps://imdb.com/title/tt8415836"},
			},
		},
		{
			Type: "section",
			Fields: &[]text{
				{Type: "mrkdwn", Text: "*Quality:*\n1080p"},
				{Type: "mrkdwn", Text: "*Release Group:*\nlegit"},
			},
		},
	},
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
	EventType:      "Download",
	ApplicationURL: "http://localhost",
}

var slackRadarrOnDownload = body{
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
			Text: &text{Type: "mrkdwn", Text: "http://localhost/movie/55"},
		},
		{
			Type: "section",
			Fields: &[]text{
				{Type: "mrkdwn", Text: "*Release Date:*\n1970-01-01"},
				{Type: "mrkdwn", Text: "*IMDB:*\nhttps://imdb.com/title/tt8415836"},
			},
		},
		{
			Type: "section",
			Fields: &[]text{
				{Type: "mrkdwn", Text: "*Quality:*\n1080p"},
				{Type: "mrkdwn", Text: "*Release Group:*\nlegit"},
			},
		},
	},
}

var sonarrOnDownload = sonarr.Data{
	Series: sonarr.Series{
		ID:     123,
		Title:  "Name Of Show!",
		IMDBID: "tt10574558",
	},
	Episodes: []sonarr.Episode{
		{ID: 555, Title: "title", SeasonNumber: 4, EpisodeNumber: 1, AirDate: "1970-01-01"},
	},
	EpisodeFile:    &sonarr.EpisodeFile{Quality: "1080p", ReleaseGroup: "legit"},
	EventType:      "Download",
	ApplicationURL: "http://localhost",
}

var slackSonarrOnDownload = body{
	TS:      "123",
	Channel: "c123",
	Blocks: []block{
		{
			Type: "header",
			Text: &text{
				Type:  "plain_text",
				Text:  ":large_green_circle: Downloaded: Name Of Show! - 4x01 - title",
				Emoji: true,
			},
		},
		{
			Type: "section",
			Text: &text{Type: "mrkdwn", Text: "http://localhost/series/name-of-show"},
		},
		{
			Type: "section",
			Fields: &[]text{
				{Type: "mrkdwn", Text: "*Release Date:*\n1970-01-01"},
				{Type: "mrkdwn", Text: "*IMDB:*\nhttps://imdb.com/title/tt10574558"},
			},
		},
		{
			Type: "section",
			Fields: &[]text{
				{Type: "mrkdwn", Text: "*Quality:*\n1080p"},
				{Type: "mrkdwn", Text: "*Release Group:*\nlegit"},
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
				token:   "Bearer xoxb-123",
				client:  *http.DefaultClient,
				cache:   state.New("", true, true),
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
		data     data.Data
		update   bool
		expected body
	}{
		"new movie grab": {
			channel:  "c123",
			data:     &radarrOnGrab,
			update:   false,
			expected: slackRadarrOnGrab,
		},
		"updated movie grab": {
			channel:  "c123",
			data:     &radarrOnGrab,
			update:   true,
			expected: slackRadarrUpdateOnGrab,
		},
	}

	for _, tc := range tests {
		ts := ""
		if tc.update {
			ts = "1234"
		}
		actual := onGrabInfo(tc.channel, tc.data, ts)
		assert.Equal(t, tc.expected, actual)
	}
}

func TestOnDownloadBody(t *testing.T) {
	tests := map[string]struct {
		channel  string
		data     data.Data
		ts       string
		expected body
	}{
		"new movie download": {
			channel:  "c123",
			data:     &radarrOnDownload,
			ts:       "123",
			expected: slackRadarrOnDownload,
		},
		"new episode download": {
			channel:  "c123",
			data:     &sonarrOnDownload,
			ts:       "123",
			expected: slackSonarrOnDownload,
		},
	}

	for _, tc := range tests {
		actual := onDownloadInfo(tc.channel, tc.data, tc.ts)
		assert.Equal(t, tc.expected, actual)
	}
}
