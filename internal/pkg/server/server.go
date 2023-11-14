package server

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/mbarrin/gwarr/internal/pkg/radarr"
	"github.com/mbarrin/gwarr/internal/pkg/slack"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var sc slack.Client

func Start(port int, client slack.Client) {
	sc = client

	http.HandleFunc("/radarr", radarrWebhook)
	http.Handle("/metrics", promhttp.Handler())

	p := fmt.Sprintf(":%d", port)
	err := http.ListenAndServe(p, nil)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func radarrWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid Method", 405)
		return
	}

	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	radarr, err := radarr.ParseWebhook(body)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Invalid Content", 400)
		return
	}

	err = sc.RadarrPost(*radarr)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
}
