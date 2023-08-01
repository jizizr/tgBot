package funcs

import (
	"bot/botTool"
	. "bot/config"
	"fmt"
	"net/url"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func AutoReply(update *tgbotapi.Update, message *tgbotapi.Message) {
	if message.ReplyToMessage != nil {
		if message.ReplyToMessage.From.ID != botTool.Bot.Self.ID {
			return
		}
	} else {
		if !strings.HasPrefix(message.Text, "Allen") && !strings.HasPrefix(message.Text, "allen") {
			return
		}
	}
	url := fmt.Sprintf("http://api.a20safe.com/api.php?api=51&key=%s&text=%s", API_TOKEN, url.QueryEscape(message.Text))
	replyRaw := getToMap(url)
	if replyRaw["code"].(float64) != 0 {
		return
	}
	reply := replyRaw["data"].([]interface{})[0].(map[string]interface{})["reply"].(string)
	// fmt.Println(reply)
	botTool.SendMessage(message, strings.ReplaceAll(reply, "\\n", "\n"), true)
}
