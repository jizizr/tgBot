package funcs

import (
	"bot/botTool"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var KH = [...]string{"(", ")", "[", "]", "{", "}", "<", ">", "（", "）", "《", "》", "【", "】", "［", "］", "（", "）", "｛", "｝", "＜", "＞", "『", "』", "「", "」", "«", "»"}
var ZKH = [...]string{"(", "[", "{", "<", "（", "《", "【", "［", "（", "｛", "＜", "『", "「", "«"}
var YKH = [...]string{")", "]", "}", ">", "）", "》", "】", "］", "）", "｝", "＞", "』", "」", "»"}

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
	botTool.SendMessage(update, &text, true)
}
