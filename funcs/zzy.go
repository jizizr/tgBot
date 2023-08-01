package funcs

import (
	"bot/botTool"
	"fmt"
	"math/rand"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var baseDao = "钟先生今天%s哦，%s!"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Zzy(update *tgbotapi.Update, message *tgbotapi.Message) {
	str := "Allen 闻到了野生寄远的气息！\n寄远兄，你别自爆啦！"
	botTool.SendMessage(message, str, true)
}

func Dao(update *tgbotapi.Update, message *tgbotapi.Message) {
	var str string
	if rand.Intn(2) == 0 {
		str = fmt.Sprintf(baseDao, "导了", "掌声鼓励")
	} else {
		str = fmt.Sprintf(baseDao, "没导", "快导快导")
	}
	botTool.SendMessage(message, str, true)
}
