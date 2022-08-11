package funcs

import (
	"bot/botTool"
	"io"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Curl(update *tgbotapi.Update) {
	arr := strings.Split(update.Message.Text, " ")
	var url string
	var msg *tgbotapi.Message
	if len(arr) == 1 {
		str := "Usage: curl [url]"
		botTool.SendMessage(update, &str, true)
		return
	} else {
		str := "正在请求中..."
		msg, _ = botTool.SendMessage(update, &str, true)
		url = arr[1]
	}
	if url[:4] != "http" {
		url = "https://" + url
	}
	resp, err := http.Get(url)
	if err != nil {
		str := err.Error()
		botTool.Edit(msg, &str)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		str := err.Error()
		botTool.Edit(msg, &str)
		return
	}
	if len(body) < 2000 {
		str := string(body)
		_, err = botTool.SendMessage(update, &str, true)
	} else {
		contentTypeArr := strings.Split(resp.Header.Get("Content-type"), "/")
		var contentType string
		if len(contentTypeArr) < 2 {
			contentType = "text"
		} else {
			contentType = contentTypeArr[len(contentTypeArr)-1]
			contentType = strings.Split(contentType, ";")[0]
		}
		_, err = botTool.SendDocument(update, body, "curl."+contentType, true, "结果太长，请下载")
	}
	if err != nil {
		str := err.Error()
		botTool.Edit(msg, &str)
	} else {
		str := "获取成功"
		botTool.Edit(msg, &str)
	}
}
