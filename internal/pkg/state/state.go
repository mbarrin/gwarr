package state

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
)

type entry struct {
	TS    string `json:"ts,omitempty"`
	State string `json:"state,omitempty"`
}

type State struct {
	Path  string
	Cache cache
}

type entries map[int]entry
type cache map[string]entries

func New(cachePath string) State {
	s := State{
		Path:  cachePath,
		Cache: make(cache),
	}
	s.Cache["radarr"] = make(entries)

	err := s.read()
	if err != nil {
		slog.With("package", "state").Error("Could not read config: " + err.Error())
	}

	fmt.Println(s)

	return s
}

func (s State) Set(t string, id int, ts string, event string) {
	s.Cache[t].set(id, ts, event)
	s.write()
}

func (s State) Delete(t string, id int) {
	s.Cache[t].delete(id)
	s.write()
}

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
