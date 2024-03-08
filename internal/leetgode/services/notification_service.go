package services

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"leetgode/internal/db/postgres"
	"leetgode/internal/leetgode/models"
	"strconv"
	"strings"
	"time"
)

type NotificationService struct {
	bot                    *tgbotapi.BotAPI
	notificationRepository *postgres.NotificationRepository
}

func NewNotificationService(bot *tgbotapi.BotAPI, notificationRepository *postgres.NotificationRepository) *NotificationService {
	return &NotificationService{
		bot:                    bot,
		notificationRepository: notificationRepository,
	}
}

func (ns *NotificationService) CheckNotifications() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ns.notifyUsers()
		}
	}
}

func (ns *NotificationService) notifyUsers() {
	currentTime := time.Now()
	fmt.Println("Sending reminders...")
	usersWithNotifications, err := ns.notificationRepository.GetUsersWithNotification(time.Date(0, 0, 0, currentTime.Hour(), currentTime.Minute(), 0, 0, currentTime.Location()))
	if err != nil {
		fmt.Println("Error getting users with notifications:", err)
		return
	}

	for _, user := range usersWithNotifications {
		ns.SendMessage(user.UserID, "A wise man once said to me: <i>«Consistency is the key».</i>\n\nIt's time to do your daily <b>Leetcode</b> routine! 🔥", nil)
	}
}

func (ns *NotificationService) HandleRemindersMenu(msg *tgbotapi.Message) {
	ns.ReplyMessage(msg, "You can set reminders that will force you to do your daily leetcode routine.", &models.ReminderMenuKeyboard)
}

func (ns *NotificationService) HandleAddReminder(msg *tgbotapi.Message) {
	ns.ReplyMessage(msg, "Please send me the UTC time in which you want to receive a reminder in the following format: \n\nHH:MM", nil)
}

func (ns *NotificationService) HandleWaitingAddReminderTime(msg *tgbotapi.Message) {
	layout := "15:04"

	resultTime, err := time.Parse(layout, msg.Text)
	if err != nil {
		ns.ReplyMessage(msg, "Invalid time format. Please use HH:MM format.", &models.ReminderMenuKeyboard)
		return
	}

	err = ns.notificationRepository.AddNotificationTimeIfNotExists(msg.From.ID, resultTime)
	if err != nil {
		ns.ReplyMessage(msg, "Notification time already exists.", &models.ReminderMenuKeyboard)
		return
	}
	ns.ReplyMessage(msg, "✅ Reminder added successfully.", &models.ReminderMenuKeyboard)
}

func (ns *NotificationService) HandleDeleteReminder(msg *tgbotapi.Message) {
	ns.ReplyMessage(msg, "Please send me the UTC time of reminder which you want to delete in the following format: \n\nHH:MM", nil)
}

func (ns *NotificationService) HandleWaitingDeleteReminderTime(msg *tgbotapi.Message) {
	layout := "15:04"

	resultTime, err := time.Parse(layout, msg.Text)
	if err != nil {
		ns.ReplyMessage(msg, "Invalid time format. Please use HH:MM format next time.", nil)
		return
	}
	err = ns.notificationRepository.RemoveNotificationTime(msg.From.ID, resultTime)
	if err != nil {
		ns.ReplyMessage(msg, "Error while removing notification time.", nil)
		return
	}
	ns.ReplyMessage(msg, "✅ Reminder deleted successfully.", &models.ReminderMenuKeyboard)
}

func (ns *NotificationService) HandleShowReminders(msg *tgbotapi.Message) {
	reminders, err := ns.notificationRepository.GetNotificationTimes(msg.From.ID)
	if err != nil {
		ns.ReplyMessage(msg, "Error while getting reminders. Try again later.", &models.ReminderMenuKeyboard)
		return
	}
	if len(reminders) == 0 {
		ns.ReplyMessage(msg, "You have no reminders yet. Try adding a new one.", nil)
		return
	}
	var builder strings.Builder
	builder.WriteString("<b> Available reminders: </b>\n\n")
	for i, rem := range reminders {
		builder.WriteString(strconv.Itoa(i+1) + ") " + rem.Format("15:04") + " \n\n")
	}
	ns.ReplyMessage(msg, builder.String(), nil)
}

func (ns *NotificationService) ReplyMessage(replyTo *tgbotapi.Message, text string, markup *tgbotapi.ReplyKeyboardMarkup) {
	msg := tgbotapi.NewMessage(replyTo.Chat.ID, text)
	msg.ReplyMarkup = markup
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyToMessageID = replyTo.MessageID
	_, err := ns.bot.Send(msg)
	if err != nil {
		panic(err)
	}
}

func (ns *NotificationService) SendMessage(chatID int64, text string, markup *tgbotapi.ReplyKeyboardMarkup) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = markup
	msg.ParseMode = tgbotapi.ModeHTML
	_, err := ns.bot.Send(msg)
	if err != nil {
		panic(err)
	}
}
