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

	"github.com/mbarrin/gwarr/internal/pkg/data"
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
		http.HandleFunc("/radarr", webhook)
	}

	if sonarr {
		http.HandleFunc("/sonarr", webhook)
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

func webhook(w http.ResponseWriter, r *http.Request) {
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

	var data data.Data
	var err error

	if r.URL.Path == "/sonarr" {
		data, err = sonarr.ParseWebhook(body)
	} else if r.URL.Path == "/radarr" {
		data, err = radarr.ParseWebhook(body)
	}

	slog.With("package", "server").Debug(string(body))
	if err != nil {
		slog.With("package", "server").Error(err.Error())
		http.Error(w, "Invalid Content", 400)
		return
	}

	err = sc.Post(data)
	if err != nil {
		slog.With("package", "server").Error(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
}
