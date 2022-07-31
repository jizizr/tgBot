package funcs

import (
	"bot/botTool"
	. "bot/config"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var ocrurl = "http://api.ocr.space/parse/imageurl?apikey=%s&filetype=%s&language=%s&OCREngine=%s&detectOrientation=true&isCreateSearchablePdf=%s&url=%s"
var lan = map[string]struct{}{"ara": {}, "bul": {}, "chs": {}, "cht": {}, "hrv": {}, "cze": {}, "dan": {}, "dut": {}, "eng": {}, "fin": {}, "fre": {}, "ger": {}, "gre": {}, "hun": {}, "kor": {}, "ita": {}, "jpn": {}, "pol": {}, "por": {}, "rus": {}, "slv": {}, "spa": {}, "swe": {}, "tur": {}}

func Ocr(update *tgbotapi.Update) {
	if update.Message.ReplyToMessage == nil {
		str := "请回复一条信息\nUsage: /ocr [lan] [any]\n若any存在任意参数则会附带一个可搜索pdf文件\n语言列表:\n阿拉伯语 = `ara`\n保加利亚语 = `bul`\n中文（简体） = `chs`\n中文（繁体） = `cht`\n克罗地亚语 = `hrv`\n捷克语 = `cze`\n丹麦语 = `dan`\n荷兰语 = `dut`\n英语 = `eng`\n芬兰语 = `fin`\n法语 = `fre`\n德语  = `ger`\n希腊语  = `gre`\n匈牙利语  = `hun`\n韩语  = `kor`\n意大利语  = `ita`\n日语  = `jpn`\n波兰语  = `pol`\n葡萄牙语  = `por`\n俄语  = `rus`\n斯洛文尼亚语  = `slv`\n西班牙语  = `spa`\n瑞典语  = `swe`\n土耳其语  = `tur`"
		botTool.SendMessage(update, &str, true)
		return
	} else if update.Message.ReplyToMessage.Photo == nil && update.Message.ReplyToMessage.Document == nil {
		str := "请回复一张图片"
		botTool.SendMessage(update, &str, true)
		return
	}
	var fileID string
	var filetype string
	var language string
	var number string
	if update.Message.ReplyToMessage.Photo != nil {
		photoFileID := update.Message.ReplyToMessage.Photo
		fileID, _ = botTool.Bot.GetFileDirectURL(photoFileID[len(photoFileID)-1].FileID)
		filetype = "jpg"
	} else {
		fileID, _ = botTool.Bot.GetFileDirectURL(update.Message.ReplyToMessage.Document.FileID)
		arr := strings.Split(update.Message.ReplyToMessage.Document.FileName, ".")
		filetype = arr[len(arr)-1]
	}

	isCreateSearchablePdf := "false"
	arr := strings.Split(update.Message.Text, " ")
	if len(arr) == 1 {
		language = "chs"
	} else if len(arr) == 2 {
		if botTool.Contains(lan, arr[1]) {
			language = arr[1]
		} else {
			str := "语言列表:\n阿拉伯语 = `ara`\n保加利亚语 = `bul`\n中文（简体） = `chs`\n中文（繁体） = `cht`\n克罗地亚语 = `hrv`\n捷克语 = `cze`\n丹麦语 = `dan`\n荷兰语 = `dut`\n英语 = `eng`\n芬兰语 = `fin`\n法语 = `fre`\n德语  = `ger`\n希腊语  = `gre`\n匈牙利语  = `hun`\n韩语  = `kor`\n意大利语  = `ita`\n日语  = `jpn`\n波兰语  = `pol`\n葡萄牙语  = `por`\n俄语  = `rus`\n斯洛文尼亚语  = `slv`\n西班牙语  = `spa`\n瑞典语  = `swe`\n土耳其语  = `tur`"
			botTool.SendMessage(update, &str, true, "Markdown")
			return
		}
	} else {
		if botTool.Contains(lan, arr[1]) {
			language = arr[1]
			isCreateSearchablePdf = "true"
		} else {
			str := "语言列表:\n阿拉伯语 = `ara`\n保加利亚语 = `bul`\n中文（简体） = `chs`\n中文（繁体） = `cht`\n克罗地亚语 = `hrv`\n捷克语 = `cze`\n丹麦语 = `dan`\n荷兰语 = `dut`\n英语 = `eng`\n芬兰语 = `fin`\n法语 = `fre`\n德语  = `ger`\n希腊语  = `gre`\n匈牙利语  = `hun`\n韩语  = `kor`\n意大利语  = `ita`\n日语  = `jpn`\n波兰语  = `pol`\n葡萄牙语  = `por`\n俄语  = `rus`\n斯洛文尼亚语  = `slv`\n西班牙语  = `spa`\n瑞典语  = `swe`\n土耳其语  = `tur`"
			botTool.SendMessage(update, &str, true, "Markdown")
			return
		}
	}
	if language == "chs" || language == "cht" || language == "eng" {
		number = "5"
	} else {
		number = "3"
	}
	str := "正在识别中，请等候..."
	msg, _ := botTool.SendMessage(update, &str, true)
	url := fmt.Sprintf(ocrurl, OCR_TOKEN, filetype, language, number, isCreateSearchablePdf, fileID)
	resp, err := http.Get(url)
	if err != nil {
		str := "请求失败"
		botTool.SendMessage(update, &str, true)
		return
	}
	defer resp.Body.Close()
	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	if _, ok := res["ParsedResults"]; !ok {
		str := res["ErrorMessage"].([]string)[0]
		botTool.SendMessage(update, &str, true)
		return
	}
	if isCreateSearchablePdf == "false" {
		str = fmt.Sprintf("识别结果：\n`%s`", res["ParsedResults"].([]interface{})[0].(map[string]interface{})["ParsedText"].(string))
	} else {
		str = fmt.Sprintf("识别结果：\n`%s`\n\n[可搜索pdf下载链接](%s)", res["ParsedResults"].([]interface{})[0].(map[string]interface{})["ParsedText"].(string), res["SearchablePDFURL"].(string))
	}
	botTool.Edit(msg, &str, "Markdown")
}
