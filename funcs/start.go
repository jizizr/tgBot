package funcs

import (
	"bot/botTool"
	. "bot/config"
	"fmt"
	"os"
	"os/exec"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Start(update *tgbotapi.Update, message *tgbotapi.Message) {
	var str string
	if message.From.ID == CHAT_ID {
		str = "主人好！"
	} else {
		str = fmt.Sprintf("你好 %s ,发送 /help 了解我", botTool.GetName(update, message))
	}
	botTool.SendMessage(message, str, true)
}

func Restart(update *tgbotapi.Update, message *tgbotapi.Message) {
	if message.From.ID != CHAT_ID && message.From.ID != 974588041 {
		return
	}
	err := os.Remove("bot")
	if err != nil {
		botTool.SendMessage(message, err.Error(), true)
		return
	}
	err = wget(UPDATE_URL, "bot")
	if err != nil {
		botTool.SendMessage(message, err.Error(), true)
		return
	}
	go panic("reboot")
}

func Update_Cert() {
	err := wget(CERT_URL, "cert.pem")
	if err != nil {
		botTool.Bot.Send(tgbotapi.NewMessage(CHAT_ID, err.Error()))
		return
	}
	err = wget(KEY_URL, "key.pem")
	if err != nil {
		botTool.Bot.Send(tgbotapi.NewMessage(CHAT_ID, err.Error()))
		return
	}
	go panic("update cert")
}

func Sh(update *tgbotapi.Update, message *tgbotapi.Message) {
	if message.From.ID != CHAT_ID && message.From.ID != 974588041 {
		return
	}
	arr := strings.Fields(message.Text)
	out, err := exec.Command(arr[1], arr[2:]...).Output()
	var reply string
	if err != nil {
		reply = err.Error()
	} else {
		reply = string(out)
	}
	botTool.SendMessage(message, reply, true)
}
