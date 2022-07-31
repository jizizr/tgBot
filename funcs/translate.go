package funcs

import (
	"bot/botTool"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var lantype = map[string]struct{}{"sq": {}, "ar": {}, "am": {}, "az": {}, "ga": {}, "et": {}, "or": {}, "eu": {}, "be": {}, "bg": {}, "is": {}, "pl": {}, "bs": {}, "fa": {}, "af": {}, "tt": {}, "da": {}, "de": {}, "ru": {}, "fr": {}, "tl": {}, "fi": {}, "fy": {}, "km": {}, "ka": {}, "gu": {}, "kk": {}, "ht": {}, "ko": {}, "ha": {}, "nl": {}, "ky": {}, "gl": {}, "ca": {}, "cs": {}, "kn": {}, "co": {}, "hr": {}, "ku": {}, "la": {}, "lv": {}, "lo": {}, "lt": {}, "lb": {}, "rw": {}, "ro": {}, "mg": {}, "mt": {}, "mr": {}, "ml": {}, "ms": {}, "mk": {}, "mi": {}, "mn": {}, "bn": {}, "my": {}, "hmn": {}, "xh": {}, "zu": {}, "ne": {}, "no": {}, "pa": {}, "pt": {}, "ps": {}, "ny": {}, "ja": {}, "sv": {}, "sm": {}, "sr": {}, "st": {}, "si": {}, "eo": {}, "sk": {}, "sl": {}, "sw": {}, "gd": {}, "ceb": {}, "so": {}, "tg": {}, "te": {}, "ta": {}, "th": {}, "tr": {}, "tk": {}, "cy": {}, "ug": {}, "ur": {}, "uk": {}, "uz": {}, "es": {}, "iw": {}, "el": {}, "haw": {}, "sd": {}, "hu": {}, "sn": {}, "hy": {}, "ig": {}, "it": {}, "yi": {}, "hi": {}, "su": {}, "id": {}, "jw": {}, "en": {}, "yo": {}, "vi": {}, "zh-TW": {}, "zh-CN": {}, "zh-cn": {}, "zh-tw": {}}

func splitString(r rune) bool {
	return r == '.' || r == '。'
}

func split(text string) (target []string) {
	n := 0
	temp := []string{}
	arr := strings.FieldsFunc(text, splitString)

	for _, s := range arr {
		if n+len(s) > 1000 {
			target = append(target, strings.Join(temp, "."))
			temp = []string{s}
			n = len(s) + 1
		} else {
			temp = append(temp, s)
			n += len(s) + 1
		}
	}
	target = append(target, strings.Join(temp, "."))
	return
}

func Translate(update *tgbotapi.Update) {
	var text string
	var lan string
	user := strings.Split(update.Message.Text, " ")
	if update.Message.ReplyToMessage != nil {
		text = update.Message.ReplyToMessage.Text
		if text == "" {
			text = update.Message.ReplyToMessage.Caption
		}
		if len(user) > 1 {
			lan = user[1]
		} else {
			lan = "zh-cn"
		}
	}
	if text == "" {
		if len(user) == 1 {
			str := "请输入要翻译的内容:\nUsage:()内为可选参数[]为必选参数\n1.回复需要翻译的语句：\n/translate (目标语言)\n2.翻译自己的句子\n/translate [文本内容] (目标语言)\n"
			botTool.SendMessage(update, &str, true)
			return
		} else {
			if botTool.Contains(lantype, user[len(user)-1]) {
				lan = user[len(user)-1]
				text = strings.Join(user[1:len(user)-1], " ")
			} else {
				lan = "zh-cn"
				text = strings.Join(user[1:], " ")
			}
		}
	}
	var source1, source2 []string
	var arr []string
	if len(text) > 1000 {
		arr = split(text)
	} else {
		arr = []string{text}
	}

	for _, text := range arr {
		text = url.QueryEscape(text + ".")
		resp, err := http.Get(fmt.Sprintf("https://translate.googleapis.com/translate_a/single?client=gtx&sl=auto&tl=%s&dt=at&q=%s", lan, text))
		if err != nil {
			return
		}
		var res []interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		if res[5] == nil {
			resp, err := http.Get(fmt.Sprintf("https://translate.googleapis.com/translate_a/single?client=gtx&sl=auto&tl=en&dt=at&q=%s", text))
			if err != nil {
				return
			}
			json.NewDecoder(resp.Body).Decode(&res)
		}
		defer resp.Body.Close()
		if res[5] == nil {
			str := "翻译失败: 超过接口长度限制"
			botTool.SendMessage(update, &str, true)
			return
		}
		res = res[5].([]interface{})
		n := len(res)

		for i := 0; i < n; i++ {
			source := res[i].([]interface{})[2]
			if source == nil {
				continue
			}
			stemp := source.([]interface{})
			source1 = append(source1, stemp[0].([]interface{})[0].(string)+"\n")
			if len(stemp) == 1 {
				continue
			}
			source2 = append(source2, source.([]interface{})[1].([]interface{})[0].(string)+"\n")
		}
	}
	source1Str := strings.Join(source1, "")
	source2Str := strings.Join(source2, "")
	bodyStr := fmt.Sprintf("接口1:\n`%s`\n\n接口2:\n`%s`", source2Str, source1Str)
	botTool.SendMessage(update, &bodyStr, true, "Markdown")
}
