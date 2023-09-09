package main

import (
	"log/slog"
	"os"
)

type Envs struct {
	TelegramBotToken string
	NotionSecret     string
	NotionDatabaseID string
	Port             string
}

func loadEnv(logger *slog.Logger) Envs {
	tgBotToken, ok := os.LookupEnv("TELEGRAM_BOT_TOKEN")
	if !ok {
		logger.Error("TELEGRAM_BOT_TOKEN env not set")
		os.Exit(1)
	}
	notionSecret, ok := os.LookupEnv("NOTION_SECRET")
	if !ok {
		logger.Error("NOTION_SECRET env not set")
		os.Exit(1)
	}

	notionDatabase, ok := os.LookupEnv("NOTION_DATABASE_ID")
	if !ok {
		logger.Error("NOTION_DATABASE_ID env not set")
		os.Exit(1)
	}

	port, portOK := os.LookupEnv("PORT")
	if !portOK {
		port = "8080"
	}

	logger.Info("Env loaded successfully")

	// get envs
	return Envs{
		TelegramBotToken: tgBotToken,
		NotionSecret:     notionSecret,
		NotionDatabaseID: notionDatabase,
		Port:             port,
	}
}
