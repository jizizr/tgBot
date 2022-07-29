package funcs

import (
	"bot/botTool"
	. "bot/config"
	"fmt"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func splitfunc(r rune) bool {
	return r == ' ' || r == ':'
}

func Ping(update *tgbotapi.Update) {
	str := "正在测试，plz wait..."
	var ip, port string
	msg, _ := botTool.SendMessage(update, &str, true)
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
	var a, b string
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { a = goHttp(TP_IP1, ip, port, &wg) }()
	go func() { b = goHttp(TP_IP2, ip, port, &wg) }()
	wg.Wait()
	str = fmt.Sprintf("CN:%s\nUK:%s", a, b)
	botTool.Edit(msg, &str)
}
