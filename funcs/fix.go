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

// var KH = [...]rune{'(', ')', '[', ']', '{', '}', '（', '）', '《', '》', '【', '】', '［', '］', '｛', '｝', '＜', '＞', '『', '』', '「', '」', '«', '»'}
var ZKH = "([{（《【［｛＜『「«"
var YKH = ")]}）》】］｝＞』」»"

// var ZKHINDEX = [...]rune{'(', '[', '{', '（', '《', '【', '［', '｛', '＜', '『', '「', '«'}
// var YKHINDEX = [...]rune{')', ']', '}', '）', '》', '】', '］', '｝', '＞', '』', '」', '»'}
var ZKHMAP = map[rune]struct{}{'(': {}, '[': {}, '{': {}, '（': {}, '《': {}, '【': {}, '［': {}, '｛': {}, '＜': {}, '『': {}, '「': {}, '«': {}}
var YKHMAP = map[rune]struct{}{')': {}, ']': {}, '}': {}, '）': {}, '》': {}, '】': {}, '］': {}, '｝': {}, '＞': {}, '』': {}, '」': {}, '»': {}}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Fix(update *tgbotapi.Update, message *tgbotapi.Message) {
	var st []rune
	var s rune
	for _, i := range message.Text {
		if find(i, ZKHMAP) {
			s, _ = utf8.DecodeRuneInString(YKH[strings.IndexRune(ZKH, i):])
			st = append(st, s)
		} else if find(i, YKHMAP) {
			if len(st) != 0 && st[len(st)-1] == i {
				st = st[:len(st)-1]
				continue
			}
			s, _ = utf8.DecodeRuneInString(ZKH[strings.IndexRune(YKH, i):])
			st = append(st, s)
		}
	}

	text := string(st)
	if text == "" {
		return
	}
	//random
	n := utf8.RuneCountInString(text)
	if n > 8 && rand.Intn(10) < 5 {
		text1 := strings.Repeat("这", n)
		text = fmt.Sprintf("%s\n\n你TM拿%s么多括号给宁妈上坟吗？！", text, text1)
	}
	_, err := botTool.SendMessage(message, text, true)
	if err != nil {
		str := "宁个寄吧，发了这么多括号，您拿他吃饭吗？"
		botTool.SendMessage(message, str, true)
	}
}
