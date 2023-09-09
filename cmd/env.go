package main

import (
	"log"
	"os"
)

type Envs struct {
	TelegramBotToken string
	NotionSecret     string
	NotionDatabaseID string
	Port             string
}

func loadEnv() Envs {
	tgBotToken, ok := os.LookupEnv("TELEGRAM_BOT_TOKEN")
	if !ok {
		panic("TELEGRAM_BOT_TOKEN env not set")
	}
	notionSecret, ok := os.LookupEnv("NOTION_SECRET")
	if !ok {
		panic("NOTION_SECRET env not set")
	}

	notionDatabase, ok := os.LookupEnv("NOTION_DATABASE_ID")
	if !ok {
		panic("NOTION_DATABASE_ID env not set")
	}

	port, portOK := os.LookupEnv("PORT")
	if !portOK {
		port = "8080"
	}

	log.Println("env loaded")

	// get envs
	return Envs{
		TelegramBotToken: tgBotToken,
		NotionSecret:     notionSecret,
		NotionDatabaseID: notionDatabase,
		Port:             port,
	}
}
