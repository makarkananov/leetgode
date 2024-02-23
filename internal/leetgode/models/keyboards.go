package models

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var MainMenuKeyboard = tgbotapi.ReplyKeyboardMarkup{
	Keyboard: [][]tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🔥 Show daily"),
			tgbotapi.NewKeyboardButton("📊 Statistics"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🔑 Set my Leetcode username"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("⏰ Reminders menu"),
		)},
}

var ReminderMenuKeyboard = tgbotapi.ReplyKeyboardMarkup{
	Keyboard: [][]tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("➕ Add reminder"),
			tgbotapi.NewKeyboardButton("🗑️ Delete reminder"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("👀 Show all of my reminders"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🔙 Back to main menu"),
		),
	},
}
