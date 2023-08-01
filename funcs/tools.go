package funcs

import (
	"bot/botTool"
	. "bot/config"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetProxies(update *tgbotapi.Update, message *tgbotapi.Message) {
	msgConfig, _ := botTool.SendMessage(message, "正在获取代理...", true)
	proxiesRaw := getToMap(fmt.Sprintf("http://api.a20safe.com/api.php?api=28&key=%s", API_TOKEN))
	if proxiesRaw["code"].(float64) != 0 {
		botTool.Edit(msgConfig, "接口失效！")
		return
	}
	proxies := proxiesRaw["data"].([]interface{})[0].(map[string]interface{})
	total := proxies["total"].(string)
	updateTime := proxies["time"].(string)
	proxiesListRaw := proxies["list"].([]interface{})
	var proxiesList []string
	for _, v := range proxiesListRaw {
		proxiesList = append(proxiesList, v.(string))
	}
	replyText := fmt.Sprintf("代理%s\n%s\n\n代理列表:\n`%s`", total, updateTime, strings.Join(proxiesList, "\n"))
	botTool.Edit(msgConfig, replyText, "Markdown")
}

func RandomPic(update *tgbotapi.Update, message *tgbotapi.Message) {
	picConfig := getToMap(fmt.Sprintf("http://api.a20safe.com/api.php?api=9&lx=dongman&key=%s", API_TOKEN))
	if picConfig["code"].(float64) != 0 {
		botTool.SendMessage(message, "接口失效！", true)
		return
	}
	pic := picConfig["data"].([]interface{})[0].(map[string]interface{})
	width := pic["width"].(string)
	height := pic["height"].(string)
	url := pic["imgurl"].(string)
	botTool.SendFile(message, url, true, fmt.Sprintf("Allen 给你送来了一张图片，赶紧看一下吧！\n图片参数：`%s x %s`", width, height), "Markdown")
}

func Test(update *tgbotapi.Update, message *tgbotapi.Message)  {
	botTool.SendFile(message,"https://down.a20safe.com/png/166583863028472.png",true)
}