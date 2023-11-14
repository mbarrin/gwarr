package main

import (
	"log/slog"
	"os"

	"github.com/mbarrin/gwarr/internal/pkg/server"
	"github.com/mbarrin/gwarr/internal/pkg/slack"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("GWARR is starting")

	sc := slack.New(os.Getenv("GWARR_SLACK_CHANNEL_ID"), os.Getenv("GWARR_SLACK_BOT_TOKEN"))

	server.Start(31337, sc)
}
