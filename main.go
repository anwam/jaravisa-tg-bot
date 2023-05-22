package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

type Envs struct {
	TelegramBotToken string
	NotionSecret     string
	NotionDatabaseID string
	Port             string
}

func loadEnv() Envs {
	// print full directory path
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(dir)

	// read .env file if present
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("no .env file found, loading from ../.env\n")
		if err = godotenv.Load("/app/.env"); err != nil {
			log.Printf("no ../.env file found, stop application\n")
			panic(err)
		}
	}
	log.Println("env loaded")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}
	// get envs
	return Envs{
		TelegramBotToken: getEnv("TELEGRAM_BOT_TOKEN"),
		NotionSecret:     getEnv("NOTION_SECRET"),
		NotionDatabaseID: getEnv("NOTION_DATABASE_ID"),
		Port:             port,
	}
}

func getEnv(name string) string {
	log.Printf("get env %s\n", name)
	return os.Getenv(name)
}

func main() {
	envs := loadEnv()

	bot, err := tgbotapi.NewBotAPI(envs.TelegramBotToken)
	if err != nil {
		log.Panic(err.Error())
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	service := &Service{
		notion: NewNotion(envs.NotionDatabaseID, envs.NotionSecret),
		tgBot:  bot,
	}

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
			service.handleWebhook(tgUpdate)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	http.ListenAndServe(":"+envs.Port, r)

}

type Service struct {
	notion *Notion
	tgBot  *tgbotapi.BotAPI
}

func (s *Service) handleWebhook(update *tgbotapi.Update) error {
	if update.Message == nil { // ignore non-Message updates
		return nil
	}

	if update.Message != nil { // If we got a message
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		amount, category, title := extract(update.Message.Text)
		if amount != 0 && category != "" {
			err := s.notion.Add(amount, category, title)
			if err != nil {
				log.Println(err)
			}

			messageText := fmt.Sprintf("฿ %.2f in %s added.", amount, category)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)
			msg.ReplyToMessageID = update.Message.MessageID
			s.tgBot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I don't understand this command")
			msg.ReplyToMessageID = update.Message.MessageID
			s.tgBot.Send(msg)
		}
	}
	return nil
}

func extract(s string) (amount float64, category, title string) {
	reg := regexp.MustCompile(`^(\d+(?:\.\d{1,2})?)([ftcgbm])\s*(.+)?$`)
	subMatchs := reg.FindStringSubmatch(s)
	if len(subMatchs) < 3 {
		return 0, "", ""
	}
	amountStr := subMatchs[1]
	category = getCategory(subMatchs[2])
	if len(subMatchs) > 3 && subMatchs[3] != "" {
		title = subMatchs[3]
	}

	amount, _ = strconv.ParseFloat(amountStr, 64)
	return amount, category, title
}

func getCategory(s string) string {
	switch s {
	case "b":
		return "beverage"
	case "f":
		return "food"
	case "t":
		return "transport"
	case "c":
		return "clothes"
	case "g":
		return "grocery"
	case "m":
		return "misc"
	default:
		return "unknown"
	}
}
