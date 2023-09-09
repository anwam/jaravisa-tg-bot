package app

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	ExpenseCommand = "expense"
)

type Service struct {
	notion NotionHandler
	tgBot  *tgbotapi.BotAPI
}

type NotionHandler interface {
	Add(amount float64, category, title string) error
}

func NewService(notion NotionHandler, tgBot *tgbotapi.BotAPI) *Service {
	return &Service{
		notion: notion,
		tgBot:  tgBot,
	}
}

func (s *Service) HandleWebhook(update *tgbotapi.Update) error {
	if update.Message == nil { // ignore non-Message updates
		return nil
	}

	if update.Message != nil { // If we got a message
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		cmd := whichCommand(update.Message.Text)
		switch cmd {
		case ExpenseCommand:

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
			}
		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I don't understand this command")
			msg.ReplyToMessageID = update.Message.MessageID
			s.tgBot.Send(msg)
		}
	}
	return nil
}

func whichCommand(s string) string {
	if checkIsExpenseCommand(s) {
		return ExpenseCommand
	}

	return "unknown"
}

func checkIsExpenseCommand(s string) bool {
	reg := regexp.MustCompile(`^(\d+(?:\.\d{1,2})?)([ftcgbm])\s*(.+)?$`)
	subMatchs := reg.FindStringSubmatch(s)
	return len(subMatchs) > 2
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
