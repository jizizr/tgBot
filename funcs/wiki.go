package funcs

import (
	"bot/botTool"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var wikiUrl = "https://zh.wikipedia.org/w/api.php?action=query&list=search&format=json&srlimit=1&srsearch=%s"
var wikiRe = regexp.MustCompile(`<span class="searchmatch">|</span>`)

func Wiki(update *tgbotapi.Update, message *tgbotapi.Message) {
	arr := strings.SplitN(message.Text, " ", 2)
	var queryText string
	var msg *tgbotapi.Message
	if len(arr) == 1 {
		str := "Usage: /wiki [keyword]"
		botTool.SendMessage(message, str, true)
		return
	} else {
		str := "正在查询中..."
		msg, _ = botTool.SendMessage(message, str, true)
		queryText = url.QueryEscape(strings.TrimSpace(strings.Join(arr[1:]," ")))
	}
	result := getToMap(fmt.Sprintf(wikiUrl, queryText))["query"].(map[string]interface{})
	if result["searchinfo"].(map[string]interface{})["totalhits"].(float64) == 0 {
	    str := "没有查询到结果"
	 	botTool.Edit(msg, str)
	 	return
	 }
	resultmap := result["search"].([]interface{})[0].(map[string]interface{})
	pos := resultmap["title"].(string)
	snippet := wikiRe.ReplaceAllString(resultmap["snippet"].(string), "")
	str := fmt.Sprintf("链接: https://zh.wikipedia.org/wiki/%s\n镜像: https://wiki.774.gs/wiki/%s\n\n概要: %s", pos, pos, snippet)
	botTool.Edit(msg, str)
}
