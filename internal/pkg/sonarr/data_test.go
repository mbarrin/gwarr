package sonarr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testJSON = []byte(`{
	"series": {
		"id": 1,
		"title": "Test Title",
		"path": "/path/to/file", 
		"tvdbId": 1234,
		"tvMazeId": 0, 
		"type": "standard"
	},
	"episodes": [{"id": 123, "episodeNumber": 1, "seasonNumber": 1, "title": "Test title"}],
	"eventType": "Test",
	"applicationUrl": "http://localhost"
}`)

var testSonarr = &Data{
	Series: Series{
		ID:       1,
		Title:    "Test Title",
		Path:     "/path/to/file",
		TVDBID:   1234,
		TVMazeID: 0,
		Type:     "standard",
	},
	Episodes: []Episode{
		{
			ID:            123,
			EpisodeNumber: 1,
			SeasonNumber:  1,
			Title:         "Test title",
		},
	},
	EventType:      "Test",
	ApplicationURL: "http://localhost",
}

var grabJSON = []byte(`{
	"series": {
		"id": 1,
		"title": "Test Title",
		"path": "/path/to/show",
	        "tvdbId": 123,
	        "tvMazeId": 456,
	        "imdbId": "tt12345",
	        "type": "standard"
	},
	"episodes": [{
		"id": 123142441,
		"episodeNumber": 1,
		"seasonNumber": 1,
		"title": "Blah",
		"airDate": "1970-01-01",
	        "airDateUtc": "1970-01-01T00:00:00Z"
	}],
	"release": {
		"quality": "WEBDL-1080p", 
		"qualityVersion": 1,
		"releaseGroup": "legit",
		"releaseTitle": "Test.Title.S01E01.Multi.1080p.WEB-DL.legit.mkv",  
	        "indexer": "usenet",
		"size": 345712323295
	}, 
	"downloadClient": "sab",
	"downloadClientType": "usenet",
	"downloadId": "SABnzbd_nzo_dsionq_f", 
	"eventType": "Grab",
	"applicationUrl": "http://localhost"
}`)

var grabSonarr = &Data{
	Series: Series{
		ID:       1,
		Title:    "Test Title",
		Path:     "/path/to/show",
		TVDBID:   123,
		TVMazeID: 456,
		IMDBID:   "tt12345",
		Type:     "standard",
	},
	Episodes: []Episode{
		{
			ID:            123142441,
			EpisodeNumber: 1,
			SeasonNumber:  1,
			Title:         "Blah",
			AirDate:       "1970-01-01",
			AirDateUTC:    "1970-01-01T00:00:00Z",
		},
	},
	Release: Release{
		Quality:        "WEBDL-1080p",
		QualityVersion: 1,
		ReleaseGroup:   "legit",
		ReleaseTitle:   "Test.Title.S01E01.Multi.1080p.WEB-DL.legit.mkv",
		Indexer:        "usenet",
		Size:           345712323295,
	},
	DownloadClient:     "sab",
	DownloadClientType: "usenet",
	DownloadID:         "SABnzbd_nzo_dsionq_f",
	EventType:          "Grab",
	ApplicationURL:     "http://localhost",
}

var downloadJSON = []byte(`{ 
	"series": {
		"id": 1,
		"title": "Test Title",
		"path": "/path/to/show",
	        "tvdbId": 123,
	        "tvMazeId": 456,
	        "imdbId": "tt12345",
	        "type": "standard"
	},
	"episodes": [{
		"id": 123142441,
		"episodeNumber": 1,
		"seasonNumber": 1,
		"title": "Blah",
		"airDate": "1970-01-01",
	        "airDateUtc": "1970-01-01T00:00:00Z"
	}], 
	"episodeFile": {
		"id": 657591, 
		"relativePath": "Season 1/Test.Title.S01E01.Multi.1080p.WEB-DL.legit.mkv", 
		"path": "/file/go/here/Test.Title.S01E01.Multi.1080p.WEB-DL.legit.mkv/Test.Title.S01E01.Multi.1080p.WEB-DL.legit.mkv", 
		"quality": "WEBDL-1080p", 
		"qualityVersion": 1, 
		"releaseGroup": "legit",
		"sceneName": "Test.Title.S01E01.Multi.1080p.WEB-DL.legit.mkv", 
		"size": 3316751093
	},
	"isUpgrade": false,
	"downloadClient": "sab",
	"downloadClientType": "usenet",
	"downloadId": "SABnzbd_nzo_dsionq_f",
	"eventType": "Download",
	"applicationUrl": "http://localhost"
}`)

var downloadSonarr = &Data{
	Series: Series{
		ID:       1,
		Title:    "Test Title",
		Path:     "/path/to/show",
		TVDBID:   123,
		TVMazeID: 456,
		IMDBID:   "tt12345",
		Type:     "standard",
	},
	Episodes: []Episode{
		{
			ID:            123142441,
			EpisodeNumber: 1,
			SeasonNumber:  1,
			Title:         "Blah",
			AirDate:       "1970-01-01",
			AirDateUTC:    "1970-01-01T00:00:00Z",
		},
	},
	EpisodeFile: &EpisodeFile{
		ID:             657591,
		RelativePath:   "Season 1/Test.Title.S01E01.Multi.1080p.WEB-DL.legit.mkv",
		Path:           "/file/go/here/Test.Title.S01E01.Multi.1080p.WEB-DL.legit.mkv/Test.Title.S01E01.Multi.1080p.WEB-DL.legit.mkv",
		Quality:        "WEBDL-1080p",
		QualityVersion: 1,
		ReleaseGroup:   "legit",
		SceneName:      "Test.Title.S01E01.Multi.1080p.WEB-DL.legit.mkv",
		Size:           3316751093,
	},
	DownloadClient:     "sab",
	DownloadClientType: "usenet",
	DownloadID:         "SABnzbd_nzo_dsionq_f",
	EventType:          "Download",
	ApplicationURL:     "http://localhost",
}

func TestParseWebhook(t *testing.T) {
	tests := map[string]struct {
		input        []byte
		expectedData *Data
		expectedErr  error
	}{
		"test":       {input: testJSON, expectedData: testSonarr, expectedErr: nil},
		"grabbed":    {input: grabJSON, expectedData: grabSonarr, expectedErr: nil},
		"downloaded": {input: downloadJSON, expectedData: downloadSonarr, expectedErr: nil},
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

func TestURLID(t *testing.T) {
	tests := map[string]struct {
		data     Data
		expected string
	}{
		"single word":           {data: Data{ApplicationURL: "http://localhost", Series: Series{Title: "Show"}}, expected: "http://localhost/series/show"},
		"single word with year": {data: Data{ApplicationURL: "http://localhost", Series: Series{Title: "Show 2016"}}, expected: "http://localhost/series/show-2016"},
		"multi word":            {data: Data{ApplicationURL: "http://localhost", Series: Series{Title: "Show Name Here"}}, expected: "http://localhost/series/show-name-here"},
		"multi word with year":  {data: Data{ApplicationURL: "http://localhost", Series: Series{Title: "Show Name Here 2020"}}, expected: "http://localhost/series/show-name-here-2020"},
		"! symbol":              {data: Data{ApplicationURL: "http://localhost", Series: Series{Title: "Bang!"}}, expected: "http://localhost/series/bang"},
		": symbol":              {data: Data{ApplicationURL: "http://localhost", Series: Series{Title: "Show: Title"}}, expected: "http://localhost/series/show-title"},
		"multi symbol":          {data: Data{ApplicationURL: "http://localhost", Series: Series{Title: "lol` wow`> 1970"}}, expected: "http://localhost/series/lol-wow-1970"},
	}

	for _, tc := range tests {
		actual := tc.data.URL()
		assert.Equal(t, actual, tc.expected)
	}
}
