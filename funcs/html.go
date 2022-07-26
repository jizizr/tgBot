package funcs

import (
	"bot/botTool"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Html(update *tgbotapi.Update) {
	arr := strings.Split(update.Message.Text, " ")
	if len(arr) == 1 {
		str := "Usage: /html [url]"
		botTool.SendMessage(update, &str, true)
		return
	}
	resp, _ := http.Get(fmt.Sprintf("http://ping.774.gs/pic?url=%s", arr[1]))

	fmt.Printf("http://ping.774.gs/pic?url=%s\n", arr[1])
	// out, _ := os.Create("1.jpg")

	// io.Copy(out, resp.Body)

	body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println(len(body))
	if len(body) == 0 {
		str := "请检查网址是否正确"
		botTool.SendMessage(update, &str, true)
		return
	}
	base64.StdEncoding.Decode(body, body)
	botTool.SendPhoto(fmt.Sprint(update.Message.Chat.ID), body)
}
