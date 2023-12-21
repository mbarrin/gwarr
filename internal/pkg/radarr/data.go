/*
Package radarr defines the structure of a radarr webhook
*/
package radarr

import (
	"encoding/json"
	"fmt"
	"log/slog"
)

// ParseError defines a custom error type for failing to turn
// a webhook into a radarr.Data struct
type ParseError struct{}

func (pe *ParseError) Error() string {
	return "Unable to parse webhook"
}

// Data defines the structure of a Radarr webhook
type Data struct {
	DeleteReason       string               `json:"deleteReason,omitempty"`
	DeletedFiles       bool                 `json:"deletedFiles,omitempty"`
	DownloadClient     string               `json:"downloadClient,omitempty"`
	DownloadClientType string               `json:"downloadClientType,omitempty"`
	DownloadID         string               `json:"downloadId,omitempty"`
	EventType          string               `json:"eventType,omitempty"`
	InstanceName       string               `json:"instanceName,omitempty"`
	IsUpgrade          bool                 `json:"isUpgrade,omitempty"`
	Movie              Movie                `json:"movie"`
	MovieFile          *MovieFile           `json:"movieFile,omitempty"`
	Release            *Release             `json:"release,omitempty"`
	RemoteMovie        *RemoteMovie         `json:"remoteMovie,omitempty"`
	RenamedMovieFiles  []*RenamedMovieFiles `json:"renamedMovieFiles,omitempty"`
	ApplicationURL     string               `json:"applicationUrl,omitempty"`
}

// Movie defines a movie
type Movie struct {
	ID          int    `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Year        int    `json:"year,omitempty"`
	ReleaseDate string `json:"releaseDate,omitempty"`
	FolderPath  string `json:"folderPath,omitempty"`
	TMDBID      int    `json:"tmdbId,omitempty"`
	IMDBID      string `json:"imdbId,omitempty"`
}

// RemoteMovie defines external data about a movie
type RemoteMovie struct {
	TMDBID int    `json:"tmdbId,omitempty"`
	IMDBID string `json:"imdbId,omitempty"`
	Title  string `json:"title,omitempty"`
	Year   int    `json:"year,omitempty"`
}

// Release defines metadata about a movie release
type Release struct {
	Quality        string `json:"quality,omitempty"`
	QualityVersion int    `json:"qualityVersion,omitempty"`
	ReleaseGroup   string `json:"releaseGroup,omitempty"`
	ReleaseTitle   string `json:"releaseTitle,omitempty"`
	Indexer        string `json:"indexer,omitempty"`
	Size           int    `json:"size,omitempty"`
}

// MovieFile defines metadata about a local movie file
type MovieFile struct {
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

// RenamedMovieFiles defines metadata about a movie file rename
type RenamedMovieFiles struct {
	PreviousRelativePath string `json:"previousRelativePath,omitempty"`
	PreviousPath         string `json:"previousPath,omitempty"`
	ID                   int    `json:"id,omitempty"`
	RelativePath         string `json:"relativePath,omitempty"`
	Quality              string `json:"quality,omitempty"`
	QualityVersion       int    `json:"qualityVersion,omitempty"`
	IndexerFlags         string `json:"indexerFlags,omitempty"`
	Size                 int    `json:"size,omitempty"`
}

// OnHealthIssue defines information about a Radarr health issue
type OnHealthIssue struct {
	Level     string `json:"level,omitempty"`
	Message   string `json:"message,omitempty"`
	Type      string `json:"type,omitempty"`
	WikiURL   string `json:"wikiUrl,omitempty"`
	EventType string `json:"eventType,omitempty"`
}

// ParseWebhook takes a webhook and turns it into a struct
func ParseWebhook(body []byte) (*Data, error) {
	d := Data{}

	err := json.Unmarshal(body, &d)
	if err != nil {
		slog.With("package", "radarr").Error("Invalid JSON")
		return nil, &ParseError{}
	}

	if d.Movie.ID == 0 {
		slog.With("package", "radarr").Error("Bad Webhook")
		return nil, &ParseError{}
	}

	return &d, nil
}

func (d *Data) ID() int             { return d.Movie.ID }
func (d *Data) IMDBID() string      { return d.Movie.IMDBID }
func (d *Data) ReleaseDate() string { return d.Movie.ReleaseDate }
func (d *Data) Service() string     { return "radarr" }
func (d *Data) Title() string       { return fmt.Sprintf("%s (%d)", d.Movie.Title, d.Movie.Year) }
func (d *Data) Type() string        { return d.EventType }
func (d *Data) URL() string         { return fmt.Sprintf("%s/movie/%d", d.ApplicationURL, d.Movie.TMDBID) }

func (d *Data) Quality() string {
	if d.EventType == "Grab" {
		return d.Release.Quality
	} else {
		return d.MovieFile.Quality
	}
}

func (d *Data) ReleaseGroup() string {
	if d.EventType == "Grab" {
		return d.Release.ReleaseGroup
	} else {
		return d.MovieFile.ReleaseGroup
	}
}
