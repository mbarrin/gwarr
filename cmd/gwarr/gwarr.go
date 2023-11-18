package main

import (
	"errors"
	"flag"
	"log/slog"
	"os"

	"github.com/mbarrin/gwarr/internal/pkg/server"
	"github.com/mbarrin/gwarr/internal/pkg/slack"
)

func main() {
	cachePath := flag.String("cache-path", "cache.json", "path to where the cache is stored")
	port := flag.Int64("port", 31337, "run server on this port")
	radarr := flag.Bool("radarr", true, "run the radarr endpoint")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	channelID, slackBotToken, err := checkEnv()
	if err != nil {
		os.Exit(1)
	}

	slog.With("package", "main").Info("GWARR is starting")

	sc := slack.New(channelID, slackBotToken, *cachePath)

	server.Start(*port, sc, *radarr)
}

func checkEnv() (string, string, error) {
	channelID, channelIDExists := os.LookupEnv("GWARR_SLACK_CHANNEL_ID")
	if !channelIDExists {
		slog.With("package", "main").Error("Missing $GWARR_SLACK_CHANNEL_ID")
	}

	slackBotToken, slackBotTokenExists := os.LookupEnv("GWARR_SLACK_BOT_TOKEN")
	if !slackBotTokenExists {
		slog.With("package", "main").Error("Missing $GWARR_SLACK_BOT_TOKEN")
	}

	if channelIDExists && slackBotTokenExists {
		return channelID, slackBotToken, nil
	}

	return "", "", errors.New("missing required tokens. Check logs")

}
