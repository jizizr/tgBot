package funcs

import (
	"bot/botTool"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	. "bot/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Weather(update *tgbotapi.Update) {
	arr := strings.Split(update.Message.Text, " ")
	if len(arr) == 1 {
		// log.Println(arr, 1)
		replyMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Usage: /weather [城市]")
		botTool.Bot.Send(replyMsg)
		return
	}
	l := arr[1]
	tmp := getToMap(fmt.Sprintf("https://geoapi.qweather.com/v2/city/lookup?location=%s&key=%s",WEATHER_TOKEN, url.QueryEscape(l)))["location"]
	if tmp == nil {
		str := "未找到此地点"
		botTool.SendMessage(update, &str, true)
		return
	}
	location := tmp.([]interface{})[0].(map[string]interface{})
	country := location["country"].(string)
	adm1 := location["adm1"].(string)
	adm2 := location["adm2"].(string)
	name := location["name"].(string)
	l = location["id"].(string)
	res := getToMap(fmt.Sprintf("https://devapi.qweather.com/v7/weather/now?location=%s&key=%s",WEATHER_TOKEN, l))
	time := res["updateTime"]
	link := res["fxLink"]
	res = res["now"].(map[string]interface{})
	temp := res["temp"].(string)
	weather := res["text"].(string)
	windDir := res["windDir"].(string)
	feelsLike := res["feelsLike"].(string)
	windSpeed := res["windSpeed"].(string)
	w, _ := strconv.ParseFloat(windSpeed, 32)
	windSpeed = fmt.Sprintf("%.2f", w/3.6)
	windScale := res["windScale"].(string)
	precip := res["precip"].(string)
	vis := res["vis"].(string)
	text := fmt.Sprintf("*城市：%s  [%s %s %s]\n链接：%s\n天气：%s\n温度：%s℃    体感温度：%s℃ \n风力：%s, %s级, %sm/s\n当前小时累计降水：%s\n能见度：%skm\n数据更新时间：%s*", name, adm1, adm2, country, link, weather, temp, feelsLike, windDir, windScale, windSpeed, precip, vis, time)
	botTool.SendMessage(update, &text, true, "Markdown")
}
