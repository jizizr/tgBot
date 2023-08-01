package funcs

import (
	"bot/botTool"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetId(update *tgbotapi.Update, message *tgbotapi.Message) {
	if message.ReplyToMessage != nil {
		text := fmt.Sprintf("宁回复这个人的的id是: `%d`", message.ReplyToMessage.From.ID)
		botTool.SendMessage(message, text, true, "Markdown")
		return
	}
	text := fmt.Sprintf("宁的id是: `%d`", message.From.ID)
	botTool.SendMessage(message, text, true, "Markdown")
}
