package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/anwam/jaravisa-tg-bot/internal/app"
	"github.com/anwam/jaravisa-tg-bot/internal/notion"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	log.Println("Starting the application...")
	envs := loadEnv()

	bot, err := tgbotapi.NewBotAPI(envs.TelegramBotToken)
	if err != nil {
		log.Panic(err.Error())
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	service := app.NewService(
		notion.NewNotion(envs.NotionDatabaseID, envs.NotionSecret),
		bot,
	)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	r.Post("/webhooks", func(w http.ResponseWriter, r *http.Request) {
		tgUpdate := new(tgbotapi.Update)
		if err := json.NewDecoder(r.Body).Decode(tgUpdate); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		if tgUpdate != nil && tgUpdate.Message != nil {
			service.HandleWebhook(tgUpdate)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})
	http.ListenAndServe(":"+envs.Port, r)
}
