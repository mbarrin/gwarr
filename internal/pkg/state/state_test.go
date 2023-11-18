package state

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	dir, _ := os.MkdirTemp("", "gwarr")
	defer os.RemoveAll(dir)

	path := dir + "/cache.json"

	tests := map[string]struct {
		expected     State
		fileContents []byte
	}{
		"no existing cache": {
			expected:     State{Path: path, Cache: cache{"radarr": entries{}}},
			fileContents: nil,
		},
		"existing cache": {
			expected:     State{Path: path, Cache: cache{"radarr": entries{777: {TS: "14"}}}},
			fileContents: []byte("{\"radarr\":{\"777\":{\"ts\":\"14\"}}}"),
		},
	}

	for name, tc := range tests {
		if tc.fileContents != nil {
			os.WriteFile(path, tc.fileContents, 0600)
		}
		actual := New(path)
		assert.Equal(t, tc.expected, actual, name)
	}
}

func TestSet(t *testing.T) {
	dir, _ := os.MkdirTemp("", "gwarr")
	defer os.RemoveAll(dir)

	path := dir + "/cache.json"

	tests := map[string]struct {
		state        State
		expected     State
		fileContents []byte
	}{
		"add key": {
			state: New(path),
			expected: State{
				Path:  path,
				Cache: cache{"radarr": entries{666: {TS: "123", State: "Grab"}}},
			},
			fileContents: []byte("{\"radarr\":{\"666\":{\"ts\":\"123\",\"state\":\"Grab\"}}}"),
		},
		"add existing key": {
			state: State{
				Path:  path,
				Cache: cache{"radarr": entries{666: {TS: "123", State: "Grab"}}},
			},
			expected: State{
				Path:  path,
				Cache: cache{"radarr": entries{666: {TS: "123", State: "Grab"}}},
			},
			fileContents: []byte("{\"radarr\":{\"666\":{\"ts\":\"123\",\"state\":\"Grab\"}}}"),
		},
		"change value": {
			state: State{
				Path:  path,
				Cache: cache{"radarr": entries{666: {TS: "456", State: "Download"}}},
			},
			expected: State{
				Path:  path,
				Cache: cache{"radarr": entries{666: {TS: "123", State: "Grab"}}},
			},
			fileContents: []byte("{\"radarr\":{\"666\":{\"ts\":\"123\",\"state\":\"Grab\"}}}"),
		},
	}

	for name, tc := range tests {
		tc.state.Set("radarr", 666, "123", "Grab")
		assert.Equal(t, tc.expected, tc.state, name)

		f, err := os.ReadFile(path)
		if err != nil {
			assert.Fail(t, "Unable to read file")
		}
		assert.Equal(t, f, tc.fileContents, name)
	}
}

func TestDelete(t *testing.T) {
	dir, _ := os.MkdirTemp("", "gwarr")
	defer os.RemoveAll(dir)

	path := dir + "/cache.json"

	tests := map[string]struct {
		state        State
		expected     State
		fileContents []byte
	}{
		"exists": {
			state:        State{Path: path, Cache: cache{"radarr": entries{666: {TS: "123", State: "Grab"}}}},
			expected:     State{Path: path, Cache: cache{"radarr": entries{}}},
			fileContents: []byte("{\"radarr\":{}}"),
		},
		"multiple": {
			state:        State{Path: path, Cache: cache{"radarr": entries{666: {TS: "123", State: "Grab"}, 123: {TS: "456", State: "Grab"}}}},
			expected:     State{Path: path, Cache: cache{"radarr": entries{123: {TS: "456", State: "Grab"}}}},
			fileContents: []byte("{\"radarr\":{\"123\":{\"ts\":\"456\",\"state\":\"Grab\"}}}"),
		},
		"does not exist": {
			state:        New(path),
			expected:     New(path),
			fileContents: []byte("{\"radarr\":{}}"),
		},
	}

	for name, tc := range tests {
		tc.state.Delete("radarr", 666)
		assert.Equal(t, tc.state, tc.expected)
		f, err := os.ReadFile(path)
		if err != nil {
			assert.Fail(t, "Unable to read file")
		}
		assert.Equal(t, f, tc.fileContents, name)
	}
}
