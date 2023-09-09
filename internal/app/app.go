package app

import (
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	ExpenseCommand = "expense"
	GPTCommand     = "gpt"
)

type Service struct {
	notion NotionHandler
	tgBot  *tgbotapi.BotAPI
	logger *slog.Logger
}

type NotionHandler interface {
	Add(amount float64, category, title string) error
}

func NewService(notion NotionHandler, tgBot *tgbotapi.BotAPI, logger *slog.Logger) *Service {
	logger.Info("Service created")
	return &Service{
		notion: notion,
		tgBot:  tgBot,
		logger: logger,
	}
}

func (s *Service) HandleWebhook(update *tgbotapi.Update) error {
	if update.Message == nil { // ignore non-Message updates
		return nil
	}

	if update.Message != nil { // If we got a message
		cmd := whichCommand(update.Message.Text)
		switch cmd {
		case ExpenseCommand:
			{
				amount, category, title := s.extract(update.Message.Text)
				if amount != 0 && category != "" {
					err := s.notion.Add(amount, category, title)
					if err != nil {
						s.logger.Error(
							"Error when adding expense to Notion",
							slog.String("errorMessage", err.Error()),
							slog.String("command", "expense"),
							slog.String("commandMessage", update.Message.Text),
						)
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Something happens with Notion. Check system log in Cloud Run Logs.")
						msg.ReplyToMessageID = update.Message.MessageID
						s.tgBot.Send(msg)
						return nil
					}

					messageText := fmt.Sprintf("฿ %.2f in %s added.", amount, category)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)
					msg.ReplyToMessageID = update.Message.MessageID
					s.tgBot.Send(msg)
				}
			}
		case GPTCommand:
			{
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "abdul command is not implemented yet.")
				msg.ReplyToMessageID = update.Message.MessageID
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
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
	} else if checkIsGPTCommand(s) {
		return GPTCommand
	}
	return "unknown"
}

func checkIsExpenseCommand(s string) bool {
	reg := regexp.MustCompile(`^(\d+(?:\.\d{1,2})?)([ftcgbm])\s*(.+)?$`)
	subMatchs := reg.FindStringSubmatch(s)
	return len(subMatchs) > 2
}

func checkIsGPTCommand(s string) bool {
	reg := regexp.MustCompile(`^(abdul)\s(.+)$`)
	subMaches := reg.FindStringSubmatch(s)
	slog.Info("subMatches", strings.Join(subMaches, ","))
	isMatch := reg.MatchString(s)
	return isMatch
}

func (srv *Service) extract(s string) (amount float64, category, title string) {
	reg := regexp.MustCompile(`^(\d+(?:\.\d{1,2})?)([ftcgbm])\s*(.+)?$`)
	subMatches := reg.FindStringSubmatch(s)
	if len(subMatches) < 3 {
		srv.logger.Info("Cannot extract amount and category from command message",
			slog.String("subMatches", strings.Join(subMatches, ",")),
		)
		return 0, "", ""
	}
	amountStr := subMatches[1]
	category = getCategory(subMatches[2])
	if len(subMatches) > 3 && subMatches[3] != "" {
		title = subMatches[3]
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
