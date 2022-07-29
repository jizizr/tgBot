package funcs

import (
	"bot/botTool"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode/utf8"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var KH = [...]string{"(", ")", "[", "]", "{", "}", "<", ">", "（", "）", "《", "》", "【", "】", "［", "］", "（", "）", "｛", "｝", "＜", "＞", "『", "』", "「", "」", "«", "»"}
var ZKH = [...]string{"(", "[", "{", "<", "（", "《", "【", "［", "（", "｛", "＜", "『", "「", "«"}
var YKH = [...]string{")", "]", "}", ">", "）", "》", "】", "］", "）", "｝", "＞", "』", "」", "»"}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Fix(update *tgbotapi.Update) {
	var st []string
	for _, i := range update.Message.Text {
		i := string(i)
		if find(i, ZKH) {
			st = append(st, YKH[index(i, ZKH)])
		} else if find(i, YKH) {
			if len(st) != 0 && st[len(st)-1] == i {
				st = st[:len(st)-1]
			} else {
				st = append(st, ZKH[index(i, YKH)])
			}
		}
	}
	text := strings.Join(st, "")
	if text == "" {
		return
	}
	//random
	n := utf8.RuneCountInString(text)
	if n > 8 && rand.Intn(10) < 5 {
		text1 := strings.Repeat("这", n)
		text = fmt.Sprintf("%s\n\n你TM拿%s么多括号给宁妈上坟吗？！", text, text1)
	}
	_, err := botTool.SendMessage(update, &text, true)
	if err != nil {
		str := "宁个寄吧，发了这么多括号，您拿他吃饭吗？"
		botTool.SendMessage(update, &str, true)
	}
}
