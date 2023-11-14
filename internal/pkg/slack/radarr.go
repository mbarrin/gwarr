package slack

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/mbarrin/gwarr/internal/pkg/radarr"
)

func (sc *Client) RadarrPost(rdr radarr.Data) error {
	ts := ""
	e, ok := sc.cache.GetRadarr(rdr)
	if ok {
		ts = e.Timestamp()
	}
	var r *http.Request
	switch rdr.EventType {
	case "MovieAdded":
		b := RadarrOnMovieAddBody(sc.channel, rdr, ts)
		jb, err := json.Marshal(b)
		if err != nil {

			return err
		}

		if ts == "" {
			slog.Info("Added film")
			r = sc.newRequest(jb, "chat.postMessage")
		} else {
			slog.Info("Updating add?")
			r = sc.newRequest(jb, "chat.update")
		}
	case "Grab":
		b := RadarrOnGrabBody(sc.channel, rdr, ts)
		jb, err := json.Marshal(b)
		if err != nil {

			return err
		}

		if ts == "" {
			slog.Info("Posting new grab")
			r = sc.newRequest(jb, "chat.postMessage")
		} else {
			slog.Info("Updating grab")
			r = sc.newRequest(jb, "chat.update")
		}
	case "Download":
		b := RadarrOnDownloadBody(sc.channel, rdr, ts)
		jb, err := json.Marshal(b)
		if err != nil {
			return err
		}

		if ts == "" {
			slog.Info("Couldn't find existing message")
			r = sc.newRequest(jb, "chat.postMessage")
		} else {
			slog.Info("Updating to download")
			r = sc.newRequest(jb, "chat.update")
		}
	case "MovieDelete":
		b := RadarrOnMovieDeleteBody(sc.channel, rdr)
		jb, err := json.Marshal(b)
		if err != nil {
			return err
		}

		slog.Info("Posting Delete")
		r = sc.newRequest(jb, "chat.postMessage")
	default:
		return errors.New("Unable to handle event: " + rdr.EventType)
	}

	resp, err := sc.client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	response := response{}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return errors.New("Message sent, but response could not be decoded. Err: " + err.Error())
	}

	if response.OK {
		if rdr.EventType == "MovieDelete" || rdr.EventType == "Download" {
			sc.cache.DeleteRadarr(rdr)
		} else {
			sc.cache.SetRadarr(rdr, response.TS)
		}
	}

	return nil
}

func RadarrOnGrabBody(c string, r radarr.Data, ts string) body {
	b := radarrMetadata(c, r)
	b.TS = ts
	b.Blocks[0].Text.Text = fmt.Sprintf(":large_orange_circle: Grabbed: %s (%d)", r.Movie.Title, r.Movie.Year)
	b.Blocks = append(b.Blocks,
		block{
			Type: "section",
			Fields: &[]text{
				{Type: "mrkdwn", Text: "*Quality:*\n" + r.Release.Quality},
				{Type: "mrkdwn", Text: "*Release Group:*\n" + r.Release.ReleaseGroup},
			},
		},
	)
	return b
}

func RadarrOnDownloadBody(c string, r radarr.Data, ts string) body {
	b := radarrMetadata(c, r)
	b.TS = ts
	b.Blocks[0].Text.Text = fmt.Sprintf(":large_green_circle: Downloaded: %s (%d)", r.Movie.Title, r.Movie.Year)
	b.Blocks = append(b.Blocks,
		block{
			Type: "section",
			Fields: &[]text{
				{Type: "mrkdwn", Text: "*Quality:*\n" + r.MovieFile.Quality},
				{Type: "mrkdwn", Text: "*Release Group:*\n" + r.MovieFile.ReleaseGroup},
			},
		},
	)
	return b
}

func RadarrOnMovieAddBody(c string, r radarr.Data, ts string) body {
	b := radarrMetadata(c, r)
	b.TS = ts
	b.Blocks[0].Text.Text = fmt.Sprintf(":large_green_circle: Added: %s (%d)", r.Movie.Title, r.Movie.Year)
	return b
}

func RadarrOnMovieDeleteBody(c string, r radarr.Data) body {
	b := radarrMetadata(c, r)
	b.Blocks[0].Text.Text = fmt.Sprintf(":red_circle: Delete: %s (%d)", r.Movie.Title, r.Movie.Year)
	return b
}

func radarrMetadata(c string, r radarr.Data) body {
	return body{
		Channel: c,
		Blocks: []block{
			{
				Type: "header",
				Text: &text{Type: "plain_text", Emoji: true},
			},
			{
				Type: "section",
				Text: &text{
					Type: "mrkdwn",
					Text: fmt.Sprintf("%s/movie/%d", os.Getenv("GWARR_RADARR_URL"), r.Movie.ID),
				},
			},
			{
				Type: "section",
				Fields: &[]text{
					{Type: "mrkdwn", Text: "*Release Date:*\n" + r.Movie.ReleaseDate},
					{Type: "mrkdwn", Text: "*IMDB:*\nhttps://imdb.com/title/" + r.Movie.IMDBID},
				},
			},
		},
	}
}
