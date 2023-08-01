package funcs

import (
	"bot/botTool"
	. "bot/config"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var lantype = map[string]struct{}{"sq": {}, "ar": {}, "am": {}, "az": {}, "ga": {}, "et": {}, "or": {}, "eu": {}, "be": {}, "bg": {}, "is": {}, "pl": {}, "bs": {}, "fa": {}, "af": {}, "tt": {}, "da": {}, "de": {}, "ru": {}, "fr": {}, "tl": {}, "fi": {}, "fy": {}, "km": {}, "ka": {}, "gu": {}, "kk": {}, "ht": {}, "ko": {}, "ha": {}, "nl": {}, "ky": {}, "gl": {}, "ca": {}, "cs": {}, "kn": {}, "co": {}, "hr": {}, "ku": {}, "la": {}, "lv": {}, "lo": {}, "lt": {}, "lb": {}, "rw": {}, "ro": {}, "mg": {}, "mt": {}, "mr": {}, "ml": {}, "ms": {}, "mk": {}, "mi": {}, "mn": {}, "bn": {}, "my": {}, "hmn": {}, "xh": {}, "zu": {}, "ne": {}, "no": {}, "pa": {}, "pt": {}, "ps": {}, "ny": {}, "ja": {}, "sv": {}, "sm": {}, "sr": {}, "st": {}, "si": {}, "eo": {}, "sk": {}, "sl": {}, "sw": {}, "gd": {}, "ceb": {}, "so": {}, "tg": {}, "te": {}, "ta": {}, "th": {}, "tr": {}, "tk": {}, "cy": {}, "ug": {}, "ur": {}, "uk": {}, "uz": {}, "es": {}, "iw": {}, "el": {}, "haw": {}, "sd": {}, "hu": {}, "sn": {}, "hy": {}, "ig": {}, "it": {}, "yi": {}, "hi": {}, "su": {}, "id": {}, "jw": {}, "en": {}, "yo": {}, "vi": {}, "zh-TW": {}, "zh-CN": {}, "zh-cn": {}, "zh-tw": {}}
var client = http.Client{}
var t = NewTranslator(&client)

type Translator struct {
	client *http.Client
}

func NewTranslator(client *http.Client) Translator {
	return Translator{client: client}
}

func (t Translator) Translate(text, sourceLang, targetLang string) (string, string, error) {
	var result []interface{}
	var translated []string

	urlStr := fmt.Sprintf(
		"https://translate.googleapis.com/translate_a/single?client=gtx&sl=%s&tl=%s&dt=t&q=%s",
		sourceLang,
		targetLang,
		url.QueryEscape(text),
	)

	req, _ := http.NewRequest(http.MethodGet, urlStr, nil)
	res, err := t.client.Do(req)

	if err != nil {
		return "err", "", errors.New("error getting translate.googleapis.com")
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return "err", "", errors.New("error reading response body")
	}

	if res.StatusCode != 200 {
		return "err", "", errors.New("translation failed")
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return "err", "", errors.New("error unmarshaling body")
	}
	if len(result) > 0 {
		lang := result[len(result)-1].([]interface{})
		trueSourceLang := lang[len(lang)-1].([]interface{})[0].(string)
		if strings.ToLower(trueSourceLang) == targetLang {
			return "", "", nil
		}
		data := result[0]
		for _, slice := range data.([]interface{}) {
			for _, translatedText := range slice.([]interface{}) {
				translated = append(translated, fmt.Sprintf("%v", translatedText))
				break
			}
		}
		return strings.Join(translated, ""), trueSourceLang, nil
	}
	return "err", "", errors.New("translation not found")
}

func Translator2(text string) string {
	url1 := fmt.Sprintf("http://api.a20safe.com/api.php?api=30&key=%s&text=%s", API_TOKEN, url.QueryEscape(text))
	// fmt.Println(url1)
	tRaw := getToMap(url1)
	if tRaw["code"].(float64) != 0 {
		return ""
	}
	t := tRaw["data"].([]interface{})[0].(map[string]interface{})["result"].(string)
	return strings.ReplaceAll(t, "<br>", "\n")
}

func Translate(update *tgbotapi.Update, message *tgbotapi.Message) {
	str := "正在翻译..."
	msg, _ := botTool.SendMessage(message, str, true)
	var text string
	var lan string
	user := strings.Split(message.Text, " ")
	if message.ReplyToMessage != nil {
		text = message.ReplyToMessage.Text
		if text == "" {
			text = message.ReplyToMessage.Caption
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
			botTool.Edit(msg, str)
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
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				return
			}
		}()
		defer wg.Done()
		bodyStr, lang, _ := t.Translate(text, "auto", lan)
		if bodyStr == "" {
			lan = "en"
			bodyStr, lang, _ = t.Translate(text, "auto", lan)
		}
		bodyStr = fmt.Sprintf("接口1:\n*Translate from*  `%s`  *to*  `%s`\n*Result:*\n`%s`", lang, lan, bodyStr)
		botTool.Edit(msg, bodyStr, "Markdown")
	}()
	go func() {
		defer func() {
			if err := recover(); err != nil {
				return
			}
		}()
		defer wg.Done()
		bodyStr := Translator2(text)
		if bodyStr == "" {
			return
		}
		bodyStr = fmt.Sprintf("接口2:\n*Result:*\n`%s`", bodyStr)
		botTool.SendMessage(message, bodyStr, true, "Markdown")
	}()
	wg.Wait()
}
