package funcs

import (
	"bot/botTool"
	"bytes"
	"encoding/json"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Json(update *tgbotapi.Update) {
	bs, _ := json.Marshal(update.Message)
	var out bytes.Buffer
	json.Indent(&out, bs, "", "    ")
	str := fmt.Sprintf("%v", out.String())
	botTool.SendMessage(update, &str, true)
}
