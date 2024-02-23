package leetgode

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"leetgode/internal/db/postgres"
	"leetgode/internal/leetgode/models"
	"leetgode/internal/leetgode/services"
)

type Leetgode struct {
	bot                 *tgbotapi.BotAPI
	notificationService *services.NotificationService
	stateService        *services.StateService
	leetcodeService     *services.LeetcodeService
	mainMenuService     *services.MainMenuService
}

func NewLeetgode(token string, userRepository *postgres.UserRepository, notificationRepository *postgres.NotificationRepository, leetcodeRepository *postgres.LeetcodeRepository) *Leetgode {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}

	h := &Leetgode{
		notificationService: services.NewNotificationService(bot, notificationRepository),
		stateService:        services.NewStateService(userRepository),
		leetcodeService:     services.NewLeetcodeService(bot, leetcodeRepository),
		mainMenuService:     services.NewMainMenuService(bot),
		bot:                 bot,
	}

	go h.notificationService.CheckNotifications()

	return h
}

func (h *Leetgode) Run() {
	h.bot.Debug = true
	updateConfig := tgbotapi.NewUpdate(0)

	updateConfig.Timeout = 30

	updates := h.bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		h.HandleUpdate(&update)
	}
}

func (h *Leetgode) HandleUpdate(update *tgbotapi.Update) {
	if update.Message == nil {
		return
	} else {
		h.HandleMessage(update.Message)
	}
}

func (h *Leetgode) HandleMessage(msg *tgbotapi.Message) {
	state, err := h.stateService.GetUserState(msg.From.ID)
	if err != nil {
		h.ReplyMessage(msg, "Something went wrong with your account state", nil)
	}
	if state == "AuthLeetcode" {
		h.leetcodeService.HandleAuthLeetcode(msg)
		if !h.stateService.TryChangeStateTo(msg.From.ID, "MainMenu") {
			h.SendImpossibleToChangeState(msg.Chat.ID, "MainMenu")
		}
	} else if state == "WaitingAddReminderTime" {
		h.notificationService.HandleWaitingAddReminderTime(msg)
		if !h.stateService.TryChangeStateTo(msg.From.ID, "RemindersMenu") {
			h.SendImpossibleToChangeState(msg.Chat.ID, "RemindersMenu")
		}
	} else if state == "WaitingDeleteReminderTime" {
		h.notificationService.HandleWaitingDeleteReminderTime(msg)
		if !h.stateService.TryChangeStateTo(msg.From.ID, "RemindersMenu") {
			h.SendImpossibleToChangeState(msg.Chat.ID, "RemindersMenu")
		}
	} else if msg.IsCommand() {
		switch msg.Command() {
		case "start":
			h.mainMenuService.HandleShowMainMenu(msg)
			if !h.stateService.TryChangeStateTo(msg.From.ID, "MainMenu") {
				h.SendImpossibleToChangeState(msg.Chat.ID, "MainMenu")
			}
		}
	} else {
		switch msg.Text {
		case "📊 Statistics":
			if !h.stateService.TryChangeStateTo(msg.From.ID, "ShowMyStats") {
				h.SendImpossibleToChangeState(msg.Chat.ID, "ShowMyStats")
				return
			}
			h.leetcodeService.HandleShowMyStats(msg)
			if !h.stateService.TryChangeStateTo(msg.From.ID, "MainMenu") {
				h.SendImpossibleToChangeState(msg.Chat.ID, "MainMenu")
			}
		case "🔥 Show daily":
			if !h.stateService.TryChangeStateTo(msg.From.ID, "ShowTodayProblem") {
				h.SendImpossibleToChangeState(msg.Chat.ID, "ShowTodayProblem")
				return
			}
			h.leetcodeService.HandleShowTodayProblem(msg)
			if !h.stateService.TryChangeStateTo(msg.From.ID, "MainMenu") {
				h.SendImpossibleToChangeState(msg.Chat.ID, "MainMenu")
			}
		case "🔑 Set my Leetcode username":
			if !h.stateService.TryChangeStateTo(msg.From.ID, "PleaseAuthLeetcode") {
				h.SendImpossibleToChangeState(msg.Chat.ID, "PleaseAuthLeetcode")
				return
			}
			h.leetcodeService.HandlePleaseAuthLeetcode(msg)
			if !h.stateService.TryChangeStateTo(msg.From.ID, "AuthLeetcode") {
				h.SendImpossibleToChangeState(msg.Chat.ID, "AuthLeetcode")
			}
		case "⏰ Reminders menu":
			if !h.stateService.TryChangeStateTo(msg.From.ID, "RemindersMenu") {
				h.SendImpossibleToChangeState(msg.Chat.ID, "RemindersMenu")
				return
			}
			h.notificationService.HandleRemindersMenu(msg)
		case "➕ Add reminder":
			if !h.stateService.TryChangeStateTo(msg.From.ID, "AddReminder") {
				h.SendImpossibleToChangeState(msg.Chat.ID, "AddReminder")
				return
			}
			h.notificationService.HandleAddReminder(msg)
			if !h.stateService.TryChangeStateTo(msg.From.ID, "WaitingAddReminderTime") {
				h.SendImpossibleToChangeState(msg.Chat.ID, "WaitingAddReminderTime")
			}
		case "🗑️ Delete reminder":
			if !h.stateService.TryChangeStateTo(msg.From.ID, "DeleteReminder") {
				h.SendImpossibleToChangeState(msg.Chat.ID, "DeleteReminder")
				return
			}
			h.notificationService.HandleDeleteReminder(msg)
			if !h.stateService.TryChangeStateTo(msg.From.ID, "WaitingDeleteReminderTime") {
				h.SendImpossibleToChangeState(msg.Chat.ID, "WaitingDeleteReminderTime")
				return
			}
		case "👀 Show all of my reminders":
			if !h.stateService.TryChangeStateTo(msg.From.ID, "ShowReminders") {
				h.SendImpossibleToChangeState(msg.Chat.ID, "ShowReminders")
				return
			}
			h.notificationService.HandleShowReminders(msg)
			if !h.stateService.TryChangeStateTo(msg.From.ID, "RemindersMenu") {
				h.SendImpossibleToChangeState(msg.Chat.ID, "RemindersMenu")
			}
		case "🔙 Back to main menu":
			h.mainMenuService.HandleShowMainMenu(msg)
			if !h.stateService.TryChangeStateTo(msg.From.ID, "MainMenu") {
				h.SendImpossibleToChangeState(msg.Chat.ID, "MainMenu")
			}
		default:
			h.HandleUnknownCommand(msg)
		}
	}
}

func (h *Leetgode) HandleUnknownCommand(msg *tgbotapi.Message) {
	h.ReplyMessage(msg, "Unknown command", &models.MainMenuKeyboard)
	if !h.stateService.TryChangeStateTo(msg.From.ID, "MainMenu") {
		h.SendImpossibleToChangeState(msg.Chat.ID, "MainMenu")
	}
}

func (h *Leetgode) ReplyMessage(replyTo *tgbotapi.Message, text string, markup *tgbotapi.ReplyKeyboardMarkup) {
	msg := tgbotapi.NewMessage(replyTo.Chat.ID, text)
	msg.ReplyMarkup = markup
	msg.ReplyToMessageID = replyTo.MessageID
	msg.ParseMode = tgbotapi.ModeHTML
	_, err := h.bot.Send(msg)
	if err != nil {
		panic(err)
	}
}

func (h *Leetgode) SendImpossibleToChangeState(chatID int64, newState string) {
	//msg := tgbotapi.NewMessage(chatID, "Impossible to change state to "+newState)
	//msg.ParseMode = tgbotapi.ModeHTML
	//_, err := h.bot.Send(msg)
	//if err != nil {
	//	panic(err)
	//}
}
