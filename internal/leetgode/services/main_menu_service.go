package services

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"leetgode/internal/leetgode/models"
)

type MainMenuService struct {
	bot *tgbotapi.BotAPI
}

func NewMainMenuService(bot *tgbotapi.BotAPI) *MainMenuService {
	return &MainMenuService{
		bot: bot,
	}
}

func (mms *MainMenuService) ReplyMessage(replyTo *tgbotapi.Message, text string, markup *tgbotapi.ReplyKeyboardMarkup) {
	msg := tgbotapi.NewMessage(replyTo.Chat.ID, text)
	msg.ReplyMarkup = markup
	msg.ReplyToMessageID = replyTo.MessageID
	msg.ParseMode = tgbotapi.ModeHTML
	_, err := mms.bot.Send(msg)
	if err != nil {
		panic(err)
	}
}

func (mms *MainMenuService) HandleShowMainMenu(msg *tgbotapi.Message) {
	mms.ReplyMessage(msg, "Here is what I can do", &models.MainMenuKeyboard)
}
