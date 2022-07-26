package funcs

import (
	"bot/botTool"
	"fmt"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var hisMatch = regexp.MustCompile(`</em>\.(.*?)</dt>`)

func History(update *tgbotapi.Update) {
	arr := strings.Split(update.Message.Text, " ")
	var his []byte
	if len(arr) == 1 {
		getHistory(&his)
	} else {
		if len(arr[1]) == 1 {
			arr[1] = "0" + arr[1]
		}
		if len(arr[2]) == 1 {
			arr[2] = "0" + arr[2]
		}

		getHistory(&his, arr[1], arr[2])
	}
	// str := string(his)
	// regexp.FindAllStringSubmatch()
	strpool := hisMatch.FindAllStringSubmatch(string(his), -1)
	var str string
	if len(arr) == 1 {
		str = "今天历史上发生了：\n"
	} else if len(arr) == 3 && len(strpool) > 0 {
		str = fmt.Sprintf("%s月%s日历史上发生了：\n", arr[1], arr[2])
	} else {
		str = "检查输入日期是否正确，如果不正确，请输入格式为：\n/history [月] [日]\n如果不输入日期，则默认为今天"
	}
	for i := range strpool {
		str += fmt.Sprintf("%d. %s\n", i+1, strpool[i][1])
	}
	if len(arr) == 3 && arr[1] == "06" && arr[2] == "04" {
		str += "11. 1989年-中国发生天安门学生运动"
	}
	botTool.SendMessage(update, &str, true)
}
