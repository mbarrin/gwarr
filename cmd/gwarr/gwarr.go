/*
Package main is the main package for gwarr
*/
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
	port := flag.Int64("port", 31337, "run server on this port")
	radarr := flag.Bool("radarr", true, "run the radarr endpoint")
	sonarr := flag.Bool("sonarr", true, "run the sonarr endpoint")
	debug := flag.Bool("debug", false, "enable debug logging")
	redisAddr := flag.String("redis-addr", "localhost:6379", "override the redis address")
	flag.Parse()

	logLevel := slog.LevelInfo
	if *debug {
		logLevel = slog.LevelDebug
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	slog.SetDefault(logger)

	channelID, slackBotToken, err := checkEnv()
	if err != nil {
		slog.With("package", "main").Error(err.Error())
		os.Exit(1)
	}

	sc, err := slack.New(channelID, slackBotToken, *redisAddr)
	if err != nil {
		os.Exit(1)
	}

	err = server.Start(*port, *sc, *radarr, *sonarr)
	if err != nil {
		os.Exit(1)
	}

	slog.With("package", "main").Info("GWARR is running")
}

func checkEnv() (string, string, error) {
	channelID, channelIDExists := os.LookupEnv("GWARR_SLACK_CHANNEL_ID")
	if !channelIDExists {
		slog.With("package", "main").Error("Missing GWARR_SLACK_CHANNEL_ID")
	}

	slackBotToken, slackBotTokenExists := os.LookupEnv("GWARR_SLACK_BOT_TOKEN")
	if !slackBotTokenExists {
		slog.With("package", "main").Error("Missing GWARR_SLACK_BOT_TOKEN")
	}

	if channelIDExists && slackBotTokenExists {
		return channelID, slackBotToken, nil
	}

	return "", "", errors.New("missing required tokens. Check logs")
}
