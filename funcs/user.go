package funcs

import (
	"bot/botTool"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func User(update *tgbotapi.Update) {
	if update.Message.From.ID != 1456780662 {
		return
	}
	var str string
	arr := strings.Split(update.Message.Text, " ")
	if len(arr) == 1 {
		str = "Usage: /user [userId]"
		botTool.SendMessage(update, &str, true)
		return
	}
	result := config.CheckId2User(arr[1])
	if result[0] == "" && result[1] == "" {
		str = "User not found"
	} else {
		str = fmt.Sprintf("User found:\nUserName: @%s\nNickName: %s", result[0], result[1])
	}
	botTool.SendMessage(update, &str, true)
}
