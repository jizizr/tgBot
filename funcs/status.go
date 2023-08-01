package funcs

import (
	"bot/botTool"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Status(update *tgbotapi.Update, message *tgbotapi.Message) {
	str := "ping..."
	startTime := time.Now()
	msg, _ := botTool.SendMessage(message, str, true)
	str = fmt.Sprintf("Pong!\n响应时间: %.2f ms", time.Since(startTime).Seconds()*1000)
	botTool.Edit(msg, str)
}
