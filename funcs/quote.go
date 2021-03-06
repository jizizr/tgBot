package funcs

import (
	"bot/botTool"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Quote(update *tgbotapi.Update) {
	res := getToMap("https://international.v1.hitokoto.cn")
	var zuozhe string
	if res["from_who"] == "None" {
		zuozhe = ""
	}
	text := fmt.Sprintf("%s\n ââ %sã%sã", res["hitokoto"], zuozhe, res["from"])
	botTool.SendMessage(update, &text, true, "Markdown")
}
