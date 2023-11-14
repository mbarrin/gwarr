package state

import (
	"testing"

	"github.com/mbarrin/gwarr/internal/pkg/radarr"
	"github.com/stretchr/testify/assert"
)

var radarrData = radarr.Data{
	Movie: radarr.Movie{
		ID:          666,
		Title:       "Film",
		Year:        1970,
		ReleaseDate: "1970-01-01",
		IMDBID:      "tt8415836",
	},
	Release: &radarr.Release{
		Quality:      "1080p",
		ReleaseGroup: "legit",
	},
	EventType: "Grab",
}

func TestNew(t *testing.T) {
	tests := map[string]struct {
		expected Cache
	}{
		"creates a cache": {expected: Cache{}},
	}

	for _, tc := range tests {
		actual := New()
		assert.Equal(t, tc.expected, actual)
	}
}
func TestSetRadar(t *testing.T) {
	tests := map[string]struct {
		cache    Cache
		radarr   radarr.Data
		expected entry
	}{
		"add key":          {cache: Cache{}, radarr: radarrData, expected: entry{TS: "123", State: "Grab"}},
		"add existing key": {cache: Cache{666: {TS: "123", State: "Grab"}}, radarr: radarrData, expected: entry{TS: "123", State: "Grab"}},
		"change value":     {cache: Cache{666: {TS: "456", State: "Download"}}, radarr: radarrData, expected: entry{TS: "123", State: "Grab"}},
	}

	for _, tc := range tests {
		tc.cache.SetRadarr(tc.radarr, "123")
		assert.Equal(t, tc.expected, tc.cache[666])
	}
}
func TestGetRadarr(t *testing.T) {
	tests := map[string]struct {
		cache         Cache
		radarr        radarr.Data
		expectedEntry *entry
		expectedBool  bool
	}{
		"sucess":      {cache: Cache{666: {TS: "123", State: "Grab"}}, radarr: radarrData, expectedEntry: &entry{TS: "123", State: "Grab"}, expectedBool: true},
		"emtpy cache": {cache: Cache{}, radarr: radarrData, expectedEntry: nil, expectedBool: false},
	}

	for _, tc := range tests {
		actual, ok := tc.cache.GetRadarr(tc.radarr)
		assert.Equal(t, tc.expectedEntry, actual)
		assert.Equal(t, tc.expectedBool, ok)
	}
}
func TestDeleteRadarr(t *testing.T) {
	tests := map[string]struct {
		cache    Cache
		radarr   radarr.Data
		expected Cache
	}{
		"exists":         {cache: Cache{666: {TS: "123", State: "Grab"}}, radarr: radarrData, expected: Cache{}},
		"multiple":       {cache: Cache{666: {TS: "123", State: "Grab"}, 777: {TS: "456", State: "Grab"}}, radarr: radarrData, expected: Cache{777: {TS: "456", State: "Grab"}}},
		"does not exist": {cache: Cache{}, radarr: radarrData, expected: Cache{}},
	}

	for _, tc := range tests {
		tc.cache.DeleteRadarr(tc.radarr)
		assert.Equal(t, tc.cache, tc.expected)
	}
}

func TestList(t *testing.T) {
	tests := map[string]struct {
		cache    Cache
		expected Cache
	}{
		"empty": {cache: Cache{}, expected: Cache{}},
	}

	for _, tc := range tests {
		actual := tc.cache.List()
		assert.Equal(t, actual, tc.expected)
	}
}
