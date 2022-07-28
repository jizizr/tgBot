package funcs

import (
	"bot/botTool"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetId(update *tgbotapi.Update) {
	if update.Message.ReplyToMessage != nil {
		text := fmt.Sprintf("宁回复这个人的的id是: `%d`", update.Message.ReplyToMessage.From.ID)
		botTool.SendMessage(update, &text, true, "Markdown")
		return
	}
	text := fmt.Sprintf("宁的id是: `%d`", update.Message.From.ID)
	botTool.SendMessage(update, &text, true, "Markdown")
}
