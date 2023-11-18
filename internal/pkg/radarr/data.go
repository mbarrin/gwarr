package radarr

import (
	"encoding/json"
	"fmt"
	"log/slog"
)

type ParseError struct{}

func (pe *ParseError) Error() string {
	return "Unable to parse webhook"
}

type Movie struct {
	ID          int    `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Year        int    `json:"year,omitempty"`
	ReleaseDate string `json:"releaseDate,omitempty"`
	FolderPath  string `json:"folderPath,omitempty"`
	TMDBID      int    `json:"tmdbId,omitempty"`
	IMDBID      string `json:"imdbId,omitempty"`
}

type RemoteMovie struct {
	TMDBID int    `json:"tmdbId,omitempty"`
	IMDBID string `json:"imdbId,omitempty"`
	Title  string `json:"title,omitempty"`
	Year   int    `json:"year,omitempty"`
}

type Release struct {
	Quality        string `json:"quality,omitempty"`
	QualityVersion int    `json:"qualityVersion,omitempty"`
	ReleaseGroup   string `json:"releaseGroup,omitempty"`
	ReleaseTitle   string `json:"releaseTitle,omitempty"`
	Indexer        string `json:"indexer,omitempty"`
	Size           int    `json:"size,omitempty"`
}

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

type OnHealthIssue struct {
	Level     string `json:"level,omitempty"`
	Message   string `json:"message,omitempty"`
	Type      string `json:"type,omitempty"`
	WikiURL   string `json:"wikiUrl,omitempty"`
	EventType string `json:"eventType,omitempty"`
}

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

func ParseWebhook(body []byte) (*Data, error) {
	d := Data{}
	fmt.Println(string(body))

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
