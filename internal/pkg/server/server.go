/*
Package server defines the server to receive webhooks
*/
package server

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/mbarrin/gwarr/internal/pkg/radarr"
	"github.com/mbarrin/gwarr/internal/pkg/slack"
	"github.com/mbarrin/gwarr/internal/pkg/sonarr"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var sc slack.Client

// Start starts a server to receive webhooks
func Start(port int64, client slack.Client, radarr, sonarr bool) {
	sc = client

	if radarr {
		http.HandleFunc("/radarr", radarrWebhook)
	}

	if sonarr {
		http.HandleFunc("/sonarr", sonarrWebhook)
	}

	http.Handle("/metrics", promhttp.Handler())

	p := fmt.Sprintf(":%d", port)
	slog.With("package", "server").Info("Server running on: " + p)
	err := http.ListenAndServe(p, nil)
	if err != nil {
		slog.With("package", "server").Error(err.Error())
		os.Exit(1)
	}
}

func radarrWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid Method", 405)
		return
	}

	body, _ := io.ReadAll(r.Body)
	defer func() {
		err := r.Body.Close()
		if err != nil {
			slog.With("package", "server").Error("Failed to close body")
		}
	}()

	slog.With("package", "server").Debug(string(body))
	radarr, err := radarr.ParseWebhook(body)
	if err != nil {
		slog.With("package", "server").Error(err.Error())
		http.Error(w, "Invalid Content", 400)
		return
	}

	err = sc.Post(radarr)
	if err != nil {
		slog.With("package", "server").Error(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
}

func sonarrWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid Method", 405)
		return
	}

	body, _ := io.ReadAll(r.Body)
	defer func() {
		err := r.Body.Close()
		if err != nil {
			slog.With("package", "server").Error("Failed to close body")
		}
	}()

	slog.With("package", "server").Debug(string(body))
	sonarr, err := sonarr.ParseWebhook(body)
	if err != nil {
		slog.With("package", "server").Error(err.Error())
		http.Error(w, "Invalid Content", 400)
		return
	}

	err = sc.Post(sonarr)
	if err != nil {
		slog.With("package", "server").Error(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
}
