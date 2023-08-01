package funcs

import (
	"bot/botTool"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Quote(update *tgbotapi.Update, message *tgbotapi.Message) {
	res := getToMap("https://international.v1.hitokoto.cn")
	var zuozhe string
	if res["from_who"] == "None" {
		zuozhe = ""
	}
	text := fmt.Sprintf("%s\n —— %s《%s》", res["hitokoto"], zuozhe, res["from"])
	botTool.SendMessage(message, text, true, "Markdown")
}
