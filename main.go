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

var envs Envs

func init() {
	// read .env file if present
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	// get envs
	envs = Envs{
		TelegramBotToken: getEnv("TELEGRAM_BOT_TOKEN"),
		NotionSecret:     getEnv("NOTION_SECRET"),
		NotionDatabaseID: getEnv("NOTION_DATABASE_ID"),
	}
}
func getEnv(name string) string {
	return os.Getenv(name)
}

func main() {
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

			amount, category := extract(update.Message.Text)
			if amount != 0 && category != "" {
				err := notion.Send(amount, category)
				if err != nil {
					log.Println(err)
				}
			}
			messageText := fmt.Sprintf("%.2f THB in %s would be add to notion expenses database", amount, category)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}
}

func extract(s string) (float64, string) {
	reg := regexp.MustCompile(`[\d]+[ftcg]`)
	//  return 1234, g
	amount := reg.FindString(s)
	amountInt, _ := strconv.ParseFloat(amount[:len(amount)-1], 64)
	category := category(amount[len(amount)-1:])
	return amountInt, category
}

func category(s string) string {
	switch s {
	case "f":
		return "food"
	case "t":
		return "transport"
	case "c":
		return "clothes"
	case "g":
		return "grocery"
	default:
		return "misc"
	}
}
