package funcs

import (
	"bot/botTool"
	"fmt"
	"io"
	"net/http"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// var ocrurl = "http://api.ocr.space/parse/imageurl?apikey=%s&filetype=%s&language=%s&OCREngine=%s&detectOrientation=true&isCreateSearchablePdf=%s&url=%s"
var ocrurl = "http://127.0.0.1:1222/ocr?url=%s"
var lan = map[string]struct{}{"ara": {}, "bul": {}, "chs": {}, "cht": {}, "hrv": {}, "cze": {}, "dan": {}, "dut": {}, "eng": {}, "fin": {}, "fre": {}, "ger": {}, "gre": {}, "hun": {}, "kor": {}, "ita": {}, "jpn": {}, "pol": {}, "por": {}, "rus": {}, "slv": {}, "spa": {}, "swe": {}, "tur": {}}
var ocrClient = http.Client{
	Timeout: 60 * time.Second,
}

func Ocr(update *tgbotapi.Update, message *tgbotapi.Message) {
	if message.ReplyToMessage == nil {
		str := "请回复一条信息\nUsage: /ocr [lan] [any]\n若any存在任意参数则会附带一个可搜索pdf文件\n语言列表:\n阿拉伯语 = `ara`\n保加利亚语 = `bul`\n中文（简体） = `chs`\n中文（繁体） = `cht`\n克罗地亚语 = `hrv`\n捷克语 = `cze`\n丹麦语 = `dan`\n荷兰语 = `dut`\n英语 = `eng`\n芬兰语 = `fin`\n法语 = `fre`\n德语  = `ger`\n希腊语  = `gre`\n匈牙利语  = `hun`\n韩语  = `kor`\n意大利语  = `ita`\n日语  = `jpn`\n波兰语  = `pol`\n葡萄牙语  = `por`\n俄语  = `rus`\n斯洛文尼亚语  = `slv`\n西班牙语  = `spa`\n瑞典语  = `swe`\n土耳其语  = `tur`"
		botTool.SendMessage(message, str, true)
		return
	} else if message.ReplyToMessage.Photo == nil && message.ReplyToMessage.Document == nil {
		str := "请回复一张图片"
		botTool.SendMessage(message, str, true)
		return
	}
	var fileID string
	if message.ReplyToMessage.Photo != nil {
		photoFileID := message.ReplyToMessage.Photo
		fileID, _ = botTool.Bot.GetFileDirectURL(photoFileID[len(photoFileID)-1].FileID)
	} else {
		fileID, _ = botTool.Bot.GetFileDirectURL(message.ReplyToMessage.Document.FileID)
	}
	str := "正在识别中，请等候..."
	msg, _ := botTool.SendMessage(message, str, true)
	url := fmt.Sprintf(ocrurl, fileID)
	resp, err := ocrClient.Get(url)
	if err != nil {
		str := "请求失败"
		botTool.SendMessage(message, str, true)
		return
	}
	defer resp.Body.Close()
	text,_:= io.ReadAll(resp.Body)
	str = fmt.Sprintf("识别结果：\n`%s`",text)
	botTool.Edit(msg, str, "Markdown")
}
