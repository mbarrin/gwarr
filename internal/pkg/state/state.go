package state

import (
	"encoding/json"
	"os"

	"github.com/mbarrin/gwarr/internal/pkg/radarr"
)

type entry struct {
	TS    string `json:"ts,omitempty"`
	State string `json:"state,omitempty"`
}

type Cache map[int]entry

func New() Cache {
	return Cache{}
}

func (c Cache) SetRadarr(r radarr.Data, ts string) {
	c[r.Movie.ID] = entry{
		TS:    ts,
		State: r.EventType,
	}
	c.writeToDisk()
}

func (c Cache) GetRadarr(r radarr.Data) (*entry, bool) {
	if value, exist := c[r.Movie.ID]; exist {
		return &value, true
	}
	return nil, false
}

func (c Cache) DeleteRadarr(r radarr.Data) {
	delete(c, r.Movie.ID)
	c.writeToDisk()
}

func (c Cache) List() Cache {
	return c
}

func (c Cache) writeToDisk() {
	j, _ := json.Marshal(c)
	os.WriteFile("cache.json", j, 0600)
}

func ReadFromDisk() (Cache, error) {
	c := Cache{}

	f, err := os.ReadFile("cache.json")
	if err != nil {
		return c, err
	}

	err = json.Unmarshal(f, &c)
	if err != nil {
		return c, err
	}

	return c, nil
}

func (e *entry) Timestamp() string {
	return e.TS
}
