package funcs

import (
	"bot/botTool"
	"fmt"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var guozaoMatch = regexp.MustCompile(`^/[a-zA-Z0-9_/]*$`)

func Guozao(update *tgbotapi.Update) {
	if update.Message.Text == "" {
		return
	}
	arr := strings.Split(update.Message.Text, " ")
	if guozaoMatch.MatchString(arr[0]) {
		return
	}
	var str string
	var player1, player2 string

	player1 = getAt(update)
	if update.Message.ReplyToMessage != nil {
		player2 = getReplyAt(update)
	} else {
		player2 = fmt.Sprintf("[自己](tg://user?id=%d)", update.Message.From.ID)
	}

	if len(arr) == 1 {
		str = fmt.Sprintf("%s %s了 %s！", player1, arr[0][1:], player2)
	} else {
		str = fmt.Sprintf("%s %s %s %s！", player1, arr[0][1:], player2, strings.Join(arr[1:], " "))
	}
	str = strings.Replace(strings.Replace(str, "$from", player1, -1), "$to", player2, -1)
	botTool.SendMessage(update, &str, true, "Markdown")
}
