package radarr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var grabJSON = []byte(`{
	"movie": {
		"id": 686,
		"title": "Film",
		"year": 1970,
		"releaseDate": "1970-01-01",
		"folderPath": "/path/to/",
		"tmdbId": 123,
		"imdbId": "tt456"
	},
	"remoteMovie": {
		"tmdbId": 123,
		"imdbId": "tt456",
		"title": "Film",
		"year": 1970
	},
	"release": {
		"quality": "Bluray-1080p",
		"qualityVersion": 1,
		"releaseGroup": "legit",
		"releaseTitle": "Film.1970.1080p.BluRay.x265-legit",
		"indexer": "usenet",
		"size": 1234578
	},
	"downloadClient": "sab",
	"downloadClientType": "usenet",
	"eventType": "Grab"
	}`)

var downloadJSON = []byte(`{
	"movie": {
		"id": 686,
		"title": "Film",
		"year": 1970,
		"releaseDate": "1970-01-01",
		"folderPath": "/path/to/",
		"tmdbId": 123,
		"imdbId": "tt456"
	},
	"remoteMovie": {
		"tmdbId": 123,
		"imdbId": "tt456",
		"title": "Film",
		"year": 1970
	},
	"movieFile": {
		"id": 36745,
		"relativePath": "newFileName",
		"path": "/path/to/file",
		"quality": "WEBDL-1080p",
		"qualityVersion": 1,
		"releaseGroup": "legit",
		"sceneName": "Film.1970.1080p.BluRay.x265-legit",
		"indexerFlags": "0",
		"size": 1234578
	},
	"isUpgrade": false,
	"downloadClient": "sab",
	"downloadClientType": "usenet",
	"downloadId": "sab_Film.1970.1080p.BluRay.x265-legit_1234",
	"eventType": "Download"
	}`)

var grabRadarr = &Data{
	Movie: Movie{
		ID:          686,
		Title:       "Film",
		Year:        1970,
		ReleaseDate: "1970-01-01",
		FolderPath:  "/path/to/",
		TMDBID:      123,
		IMDBID:      "tt456",
	},
	RemoteMovie: &RemoteMovie{
		TMDBID: 123,
		IMDBID: "tt456",
		Title:  "Film",
		Year:   1970,
	},
	Release: &Release{
		Quality:        "Bluray-1080p",
		QualityVersion: 1,
		ReleaseGroup:   "legit",
		ReleaseTitle:   "Film.1970.1080p.BluRay.x265-legit",
		Indexer:        "usenet",
		Size:           1234578,
	},
	DownloadClient:     "sab",
	DownloadClientType: "usenet",
	EventType:          "Grab",
}

var downloadRadarr = &Data{
	Movie: Movie{
		ID:          686,
		Title:       "Film",
		Year:        1970,
		ReleaseDate: "1970-01-01",
		FolderPath:  "/path/to/",
		TMDBID:      123,
		IMDBID:      "tt456",
	},
	RemoteMovie: &RemoteMovie{
		TMDBID: 123,
		IMDBID: "tt456",
		Title:  "Film",
		Year:   1970,
	},
	MovieFile: &MovieFile{
		ID:             36745,
		RelativePath:   "newFileName",
		Path:           "/path/to/file",
		Quality:        "WEBDL-1080p",
		QualityVersion: 1,
		ReleaseGroup:   "legit",
		SceneName:      "Film.1970.1080p.BluRay.x265-legit",
		IndexerFlags:   "0",
		Size:           1234578,
	},
	IsUpgrade:          false,
	DownloadClient:     "sab",
	DownloadClientType: "usenet",
	DownloadID:         "sab_Film.1970.1080p.BluRay.x265-legit_1234",
	EventType:          "Download",
}

func TestParseWebhook(t *testing.T) {
	tests := map[string]struct {
		input        []byte
		expectedData *Data
		expectedErr  error
	}{
		"grabbed":    {input: grabJSON, expectedData: grabRadarr, expectedErr: nil},
		"downloaded": {input: downloadJSON, expectedData: downloadRadarr, expectedErr: nil},
		"malformed":  {input: []byte("}"), expectedData: nil, expectedErr: &ParseError{}},
		"invalid":    {input: []byte("{}"), expectedData: nil, expectedErr: &ParseError{}},
	}

	for _, tc := range tests {
		actualData, actualErr := ParseWebhook(tc.input)
		assert.Equal(t, tc.expectedData, actualData)
		if actualErr == nil {
			assert.Equal(t, tc.expectedErr, actualErr)
		} else {
			assert.EqualError(t, actualErr, tc.expectedErr.Error())
		}

	}
}
