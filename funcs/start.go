package funcs

import (
	"bot/botTool"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Start(update *tgbotapi.Update) {
	var str string
	if update.Message.From.ID == 1456780662 {
		str = "主人好！"
	} else {
		str = fmt.Sprintf("你好 %s ,实现完整功能,请给我读取消息权限", getName(update))
	}
	botTool.SendMessage(update, &str, true)
}
