package sonarr

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
)

// ParseError defines a custom error type for failing to turn
// a webhook into a sonarr.Data struct
type ParseError struct{}

func (pe *ParseError) Error() string {
	return "Unable to parse webhook"
}

type Data struct {
	ApplicationURL     string       `json:"applicationUrl,omitempty"`
	DownloadClient     string       `json:"downloadClient,omitempty"`
	DownloadClientType string       `json:"downloadClientType,omitempty"`
	DownloadID         string       `json:"downloadId,omitempty"`
	EpisodeFile        *EpisodeFile `json:"episodeFile,omitempty"`
	Episodes           []Episode    `json:"episodes,omitempty"`
	EventType          string       `json:"eventType,omitempty"`
	Release            Release      `json:"release,omitempty"`
	Series             Series       `json:"series,omitempty"`
}

type Series struct {
	ID       int    `json:"id,omitempty"`
	Title    string `json:"title,omitempty"`
	Path     string `json:"path,omitempty"`
	TVDBID   int    `json:"tvdbId,omitempty"`
	TVMazeID int    `json:"tvMazeId,omitempty"`
	IMDBID   string `json:"imdbId,omitempty"`
	Type     string `json:"type,omitempty"`
}

type Episode struct {
	ID            int    `json:"id,omitempty"`
	EpisodeNumber int    `json:"episodeNumber,omitempty"`
	SeasonNumber  int    `json:"seasonNumber,omitempty"`
	Title         string `json:"title,omitempty"`
	AirDate       string `json:"airDate,omitempty"`
	AirDateUTC    string `json:"airDateUtc,omitempty"`
}

// EpisodeFile defines metadata about a local movie file
type EpisodeFile struct {
	ID             int    `json:"id,omitempty"`
	RelativePath   string `json:"relativePath,omitempty"`
	Path           string `json:"path,omitempty"`
	Quality        string `json:"quality,omitempty"`
	QualityVersion int    `json:"qualityVersion,omitempty"`
	ReleaseGroup   string `json:"releaseGroup,omitempty"`
	SceneName      string `json:"sceneName,omitempty"`
	IndexerFlags   string `json:"indexerFlags,omitempty"`
	Size           int    `json:"size,omitempty"`
}

// Release defines metadata about an episode release
type Release struct {
	Quality        string `json:"quality,omitempty"`
	QualityVersion int    `json:"qualityVersion,omitempty"`
	ReleaseGroup   string `json:"releaseGroup,omitempty"`
	ReleaseTitle   string `json:"releaseTitle,omitempty"`
	Indexer        string `json:"indexer,omitempty"`
	Size           int    `json:"size,omitempty"`
}

// ParseWebhook takes a webhook and turns it into a struct
func ParseWebhook(body []byte) (*Data, error) {
	d := Data{}

	err := json.Unmarshal(body, &d)
	if err != nil {
		fmt.Println(err)
		slog.With("package", "sonarr").Error("Invalid JSON")
		return nil, &ParseError{}
	}

	if d.Series.ID == 0 {
		fmt.Println(err)
		slog.With("package", "sonarr").Error("Bad Webhook")
		return nil, &ParseError{}
	}

	slog.Debug(fmt.Sprintf("%#v", d))

	return &d, nil
}

func (d *Data) urlID() string {
	lowered := strings.ToLower(d.Series.Title)
	lowered = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(lowered, "")
	split := strings.Split(lowered, " ")
	joined := strings.Join(split, "-")
	return joined
}

func (d *Data) IMDBID() string  { return d.Series.IMDBID }
func (d *Data) Type() string    { return d.EventType }
func (d *Data) URL() string     { return fmt.Sprintf("%s/series/%s", d.ApplicationURL, d.urlID()) }
func (d *Data) Service() string { return "sonarr" }

func (d *Data) ID() int {
	if d.EventType == "SeriesDelete" || d.EventType == "SeriesAdd" {
		return d.Series.ID
	}
	return d.Episodes[0].ID
}

func (d *Data) ReleaseDate() string {
	if d.EventType == "SeriesDelete" || d.EventType == "SeriesAdd" {
		return "N/A"
	}
	return d.Episodes[0].AirDate
}

func (d *Data) Year() string {
	if d.EventType == "SeriesDelete" || d.EventType == "SeriesAdd" {
		return "N/A"
	}
	return d.Episodes[0].AirDate
}

func (d *Data) Title() string {
	ep := d.Episodes[0]
	return fmt.Sprintf("%s - %dx%02d - %s", d.Series.Title, ep.SeasonNumber, ep.EpisodeNumber, ep.Title)
}

func (d *Data) Quality() string {
	if d.EventType == "Grab" {
		return d.Release.Quality
	}
	return d.EpisodeFile.Quality
}

func (d *Data) ReleaseGroup() string {
	if d.EventType == "Grab" {
		return d.Release.ReleaseGroup
	}
	return d.EpisodeFile.ReleaseGroup
}
