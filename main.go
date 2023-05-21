package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

type Envs struct {
	TelegramBotToken string
	NotionSecret     string
	NotionDatabaseID string
}

func loadEnv() Envs {
	// read .env file if present
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	// get envs
	return Envs{
		TelegramBotToken: getEnv("TELEGRAM_BOT_TOKEN"),
		NotionSecret:     getEnv("NOTION_SECRET"),
		NotionDatabaseID: getEnv("NOTION_DATABASE_ID"),
	}
}

func getEnv(name string) string {
	return os.Getenv(name)
}

func main() {
	envs := loadEnv()
	bot, err := tgbotapi.NewBotAPI(envs.TelegramBotToken)
	if err != nil {
		log.Panic(err)
	}

	notion := NewNotion(envs.NotionDatabaseID, envs.NotionSecret)
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil { // ignore non-Message updates
			continue
		}

		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			amount, category, title := extract(update.Message.Text)
			if amount != 0 && category != "" {
				err := notion.Send(amount, category, title)
				if err != nil {
					log.Println(err)
				}
			}
			messageText := fmt.Sprintf("฿ %.2f in %s added.", amount, category)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}
}

func extract(s string) (amount float64, category, title string) {
	reg := regexp.MustCompile(`^(\d+(?:\.\d{1,2})?)([ftcgbm])\s*(.+)?$`)
	// capture group 0: amount
	subMatchs := reg.FindStringSubmatch(s)
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
