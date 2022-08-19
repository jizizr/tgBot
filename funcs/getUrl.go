package funcs

import (
	"bot/botTool"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetFileUrl(update *tgbotapi.Update) {
	var str string
	var err error
	if update.Message.From.ID != 1456780662 || update.Message.ReplyToMessage == nil {
		return
	}
	fmt.Println(update.Message.ReplyToMessage.Video)
	if update.Message.ReplyToMessage.Video != nil {
		str, err = botTool.Test.GetFileDirectURL(update.Message.ReplyToMessage.Video.FileID)
	} else if update.Message.ReplyToMessage.Document != nil {
		str, err = botTool.Test.GetFileDirectURL(update.Message.ReplyToMessage.Document.FileID)
	} else {
		return
	}
	if err != nil {
		str = err.Error()
	}
	botTool.SendMessage(update, &str, false)
}
