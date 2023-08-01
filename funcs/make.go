package funcs

import (
	"bot/botTool"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func MakePic(update *tgbotapi.Update, message *tgbotapi.Message) {
	chatId := fmt.Sprint(message.Chat.ID)
	getPic(chatId, botTool.GetName(update, message))
}
