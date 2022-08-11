package funcs

import (
	"bot/botTool"
	. "bot/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var urlMatch = regexp.MustCompile(`(\s|^|https?://)([^:\./\s]+\.)+[^\./:\s]+(:(\d{1,4}|[1-5]\d{4}|6[0-4]\d{3}|65[0-4]\d{2}|655[0-2]\d|6553[0-5]))?(/\S*)*`)

func Short(update *tgbotapi.Update) {
	str := "正在生成短链..."
	msg, _ := botTool.SendMessage(update, &str, true)
	var text string
	var arr []string
	var form url.Values
	var user = strings.Split(update.Message.Text, " ")
	if update.Message.ReplyToMessage != nil {
		text = update.Message.ReplyToMessage.Text
		if text == "" {
			text = update.Message.ReplyToMessage.Caption
		}
		text = urlMatch.FindString(text)
		if text == "" {
			str = "请回复包含链接的文本"
			botTool.Edit(msg, &str)
			return
		}
		if len(user) == 1 {
			arr = append(arr, "", text)
		} else {
			arr = append(arr, "", text, user[1])
		}
	}
	if text == "" {
		arr = strings.Split(update.Message.Text, " ")
	}
	if len(arr) == 2 {
		form = url.Values{"url": {httpfix(arr[1])}}
	} else if len(arr) == 3 {
		form = url.Values{"url": {httpfix(arr[1])}, "shorturl": {arr[2]}}
	} else {
		str = "用法：/short [url] (shorturl)"
		botTool.Edit(msg, &str)
		return
	}
	resp, err := http.PostForm(SHORT_IP, form)
	if err != nil {
		str = "api请求失败"
		botTool.Edit(msg, &str)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		str = err.Error()
		botTool.Edit(msg, &str)
		return
	}
	urlMsg := map[string]string{}
	json.Unmarshal(body, &urlMsg)
	code := urlMsg["code"]
	if code == "200" {
		str = fmt.Sprintf("短链来咯：\n原链接:%s\n短链1: https://774.gs/%s\n短链2: https://ntt.gay/%s", form["url"][0], urlMsg["shorturl"], urlMsg["shorturl"])
	} else if code == "1001" {
		str = fmt.Sprintf("%s\n不符合url规则", form["url"][0])
	} else if code == "2003" {
		str = fmt.Sprintf("%s\n后缀已被使用", urlMsg["shorturl"])
	}
	// log.Println(str)
	botTool.Edit(msg, &str)
}
