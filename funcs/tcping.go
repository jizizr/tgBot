package funcs

import (
	"bot/botTool"
	. "bot/config"
	"fmt"
	"regexp"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func splitfunc(r rune) bool {
	return r == ' ' || r == ':'
}

var ipMatch = regexp.MustCompile(`(\s|^|https?://)([^:\./\s]+\.)+[^\./:\s]+(:([1-5]\d{4}|6[0-4]\d{3}|65[0-4]\d{2}|655[0-2]\d|6553[0-5]|\d{1,4}))?`)

func Tcping(update *tgbotapi.Update) {
	str := "正在测试，plz wait..."
	var ip, port, url string
	msg, _ := botTool.SendMessage(update, &str, true)
	if update.Message.ReplyToMessage != nil {
		url = update.Message.ReplyToMessage.Text
		if url == "" {
			url = update.Message.ReplyToMessage.Caption
		}
		url = ipMatch.FindString(url)
		if url == "" {
			str = "请回复包含ip的文本"
			botTool.Edit(msg, &str)
			return
		} else {
			url = strings.TrimPrefix(strings.TrimPrefix(url, "http://"), "https://")
		}
		arr := strings.FieldsFunc(url, splitfunc)
		if len(arr) == 1 {
			ip = arr[0]
			port = "80"
		} else {
			ip = arr[0]
			port = arr[1]
		}
	} else {
		arr := strings.FieldsFunc(update.Message.Text, splitfunc)
		if len(arr) == 2 {
			ip = arr[1]
			port = "80"
		} else if len(arr) != 3 {
			str = "请输入正确的格式，例如：\n/tp 91.121.210.56:54343\n/tp 91.121.210.56 54343"
			botTool.Edit(msg, &str)
			return
		} else {
			ip = arr[1]
			port = arr[2]
		}
	}
	var a, b string
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { a = goHttp(TP_IP1, ip, port, &wg) }()
	go func() { b = goHttp(TP_IP2, ip, port, &wg) }()
	wg.Wait()
	str = fmt.Sprintf("CN:%s\nUK:%s", a, b)
	botTool.Edit(msg, &str)
}
