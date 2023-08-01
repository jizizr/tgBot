package funcs

import (
	"bot/botTool"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetFileUrl(update *tgbotapi.Update, message *tgbotapi.Message) {
	var str string
	var err error
	if message.From.ID != 1456780662 || message.ReplyToMessage == nil {
		return
	}
	fmt.Println(message.ReplyToMessage.Video)
	if message.ReplyToMessage.Video != nil {
		str, err = botTool.Test.GetFileDirectURL(message.ReplyToMessage.Video.FileID)
	} else if message.ReplyToMessage.Document != nil {
		str, err = botTool.Test.GetFileDirectURL(message.ReplyToMessage.Document.FileID)
	} else {
		return
	}
	if err != nil {
		str = err.Error()
	}
	botTool.SendMessage(message, str, false)
}
