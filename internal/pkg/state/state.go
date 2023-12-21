/*
Package state handles storing data about messages sent to Slack
This allows updating of existing messages
*/
package state

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
)

// State defines the location and cache format
type State struct {
	Path  string
	Cache cache
}

type entry struct {
	TS    string `json:"ts,omitempty"`
	State string `json:"state,omitempty"`
}

type entries map[int]entry
type cache map[string]entries

// New creates a new in memory cache that is stored as json on disk
func New(cachePath string, radarr, sonarr bool) State {
	s := State{
		Path:  cachePath,
		Cache: make(cache),
	}
	if radarr {
		s.Cache["radarr"] = make(entries)
	}
	if sonarr {
		s.Cache["sonarr"] = make(entries)
	}

	err := s.read()
	if err != nil {
		slog.With("package", "state").Error("Could not read config: " + err.Error())
	}

	slog.With("package", "state").Debug(fmt.Sprintf("%v", s))

	return s
}

// Set stores the Slack message ID and its state
func (s State) Set(t string, id int, ts string, event string) {
	s.Cache[t].set(id, ts, event)
	s.write()
}

// Delete removes the info about a Slack message ID
func (s State) Delete(t string, id int) {
	s.Cache[t].delete(id)
	s.write()
}

// Timestamp returns a Slack message ID for a release
func (s State) Timestamp(t string, id int) string {
	e, _ := s.Cache[t].get(id)
	return e.TS
}

func (s State) write() {
	j, err := json.Marshal(s.Cache)
	if err != nil {
		slog.Error("Unable to marshall json: " + err.Error())
	}
	err = os.WriteFile(s.Path, j, 0600)
	if err != nil {
		slog.Error("Unable to write cache: " + err.Error())
	}
}

func (s State) read() error {
	f, err := os.ReadFile(s.Path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(f, &s.Cache)
	if err != nil {
		return err
	}

	return nil
}

func (c entries) set(key int, ts, state string) {
	c[key] = entry{
		TS:    ts,
		State: state,
	}
}

func (c entries) get(key int) (entry, bool) {
	if value, exist := c[key]; exist {
		return value, true
	}
	return entry{}, false
}

func (c entries) delete(key int) {
	delete(c, key)
}
