package funcs

import (
	"bot/botTool"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Curl(update *tgbotapi.Update) {
	arr := strings.Split(update.Message.Text, " ")
	var url string
	if len(arr) == 1 {
		replyMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Usage: curl [url]")
		botTool.Bot.Send(replyMsg)
		return
	} else {
		url = arr[1]
	}
	if url[:4] != "http" {
		url = "https://" + url
	}
	resp, err := http.Get(url)
	if err != nil {
		str := err.Error()
		botTool.SendMessage(update, &str, true)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		str := err.Error()
		botTool.SendMessage(update, &str, true)
		return
	}
	// text := string(body)
	if len(body) < 2000 {
		text := string(body)
		botTool.SendMessage(update, &text, true)
	} else {
		_, err = botTool.SendDocument(update, body, "curl.txt", true, "结果太长，请下载")
	}
	if err != nil {
		log.Println(err)
	}
}
