package services

import (
	"github.com/dustyRAIN/leetcode-api-go/leetcodeapi"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"leetgode/internal/db/postgres"
	"leetgode/internal/leetgode/models"
	"strconv"
	"strings"
)

type LeetcodeService struct {
	bot                *tgbotapi.BotAPI
	leetcodeRepository *postgres.LeetcodeRepository
}

func NewLeetcodeService(bot *tgbotapi.BotAPI, leetcodeRepository *postgres.LeetcodeRepository) *LeetcodeService {
	return &LeetcodeService{
		bot:                bot,
		leetcodeRepository: leetcodeRepository,
	}
}

func (ls *LeetcodeService) HandleShowMyStats(msg *tgbotapi.Message) {
	username, err := ls.leetcodeRepository.GetLeetcodeUser(msg.Chat.ID)
	if err != nil {
		ls.ReplyMessage(msg, "⚠️ Firstly set your <b>Leetcode</b> username using the button below.", &models.MainMenuKeyboard)
		return
	}

	profile, err := leetcodeapi.GetUserPublicProfile(username)
	if err != nil || profile.Username == "" {
		ls.ReplyMessage(msg, "Your <b>Leetcode</b> username is probably invalid. You can set another one using the button below.", &models.MainMenuKeyboard)
		return
	}

	solved, err := leetcodeapi.GetUserSolveCountByDifficulty(username)
	if err != nil {
		ls.ReplyMessage(msg, "Your <b>Leetcode</b> username is probably invalid. You can set another one using the button below.", &models.MainMenuKeyboard)
		return
	}
	var builder strings.Builder
	builder.WriteString("<b>" + profile.Username + "</b>\n\n")
	builder.WriteString("<b>Real name: </b> " + profile.Profile.RealName + "\n")
	builder.WriteString("<b>Solved: </b> " + strconv.Itoa(solved.SolveCount.SubmitStatsGlobal.AcSubmissionNum[0].Count) + "\n")
	builder.WriteString("	- <b>Easy: </b> " + strconv.Itoa(solved.SolveCount.SubmitStatsGlobal.AcSubmissionNum[1].Count) + "\n")
	builder.WriteString("	- <b>Medium: </b> " + strconv.Itoa(solved.SolveCount.SubmitStatsGlobal.AcSubmissionNum[2].Count) + "\n")
	builder.WriteString("	- <b>Hard: </b> " + strconv.Itoa(solved.SolveCount.SubmitStatsGlobal.AcSubmissionNum[3].Count) + "\n")
	builder.WriteString("<b>Ranking: </b> " + strconv.Itoa(profile.Profile.Ranking) + "\n")
	builder.WriteString("<b>Public solutions: </b> " + strconv.Itoa(profile.Profile.SolutionCount) + "\n")
	builder.WriteString("<b>Link: </b>leetcode.com/" + profile.Username + "\n")

	ls.ReplyMessage(msg, builder.String(), &models.MainMenuKeyboard)
}

func (ls *LeetcodeService) HandleShowTodayProblem(msg *tgbotapi.Message) {
	type TodayProblemResponse struct {
		Data struct {
			ActiveDailyCodingChallengeQuestion struct {
				Date     string `json:"date"`
				Link     string `json:"link"`
				Question struct {
					Difficulty string `json:"difficulty"`
					Title      string `json:"title"`
				} `json:"question"`
			} `json:"activeDailyCodingChallengeQuestion"`
		} `json:"data"`
	}

	var responseBody TodayProblemResponse
	payload := `{
    "query": "query questionOfToday { activeDailyCodingChallengeQuestion { date link question { difficulty title} }}",
    "variables": {}
  	}`
	err := (&leetcodeapi.Util{}).MakeGraphQLRequest(payload, &responseBody)
	if err != nil {
		panic(err)
	}

	var builder strings.Builder
	builder.WriteString("<b>" + responseBody.Data.ActiveDailyCodingChallengeQuestion.Question.Title + "</b>\n\n")
	builder.WriteString("<b>Date: </b> " + responseBody.Data.ActiveDailyCodingChallengeQuestion.Date + "\n")
	builder.WriteString("<b>Difficulty: </b> " + responseBody.Data.ActiveDailyCodingChallengeQuestion.Question.Difficulty + "\n")
	builder.WriteString("<b>Link: </b>leetcode.com" + responseBody.Data.ActiveDailyCodingChallengeQuestion.Link + "\n")

	ls.ReplyMessage(msg, builder.String(), &models.MainMenuKeyboard)
}

func (ls *LeetcodeService) HandlePleaseAuthLeetcode(msg *tgbotapi.Message) {
	ls.ReplyMessage(msg, "Please send me your <b>Leetcode</b> username.", nil)
}

func (ls *LeetcodeService) HandleAuthLeetcode(msg *tgbotapi.Message) {
	profile, err := leetcodeapi.GetUserPublicProfile(msg.Text)
	if err != nil || profile.Username == "" {
		ls.ReplyMessage(msg, "No such <b>Leetcode</b> username found. You can try to set another one using the button below.", &models.MainMenuKeyboard)
		return
	}

	err = ls.leetcodeRepository.AddLeetodeUser(msg.From.ID, msg.Text)
	if err != nil {
		ls.ReplyMessage(msg, "Error setting <b>Leetcode</b> username.", nil)
		return
	}

	ls.ReplyMessage(msg, "✅ Successful login.", &models.MainMenuKeyboard)

}

func (ls *LeetcodeService) ReplyMessage(replyTo *tgbotapi.Message, text string, markup *tgbotapi.ReplyKeyboardMarkup) {
	msg := tgbotapi.NewMessage(replyTo.Chat.ID, text)
	msg.ReplyMarkup = markup
	msg.ReplyToMessageID = replyTo.MessageID
	msg.ParseMode = tgbotapi.ModeHTML
	_, err := ls.bot.Send(msg)
	if err != nil {
		panic(err)
	}
}
