package funcs

import (
	"bot/botTool"
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Admin(update *tgbotapi.Update, message *tgbotapi.Message) {
	if message.From.ID != 1456780662 {
		str := fmt.Sprintf("%s\tYou are not @zrcccc!", getAt(update, message))
		botTool.SendMessage(message, str, true, "Markdown")
		return
	}
	if message.ReplyToMessage == nil {
		return
	}
	sqlStr := "INSERT IGNORE INTO `admin` (userid) values(?)"
	result, err := config.Db.Exec(sqlStr, message.ReplyToMessage.From.ID)
	if err != nil {
		log.Println(err)
	}
	_, err = result.RowsAffected()
	if err != nil {
		log.Println(err)
	}
	str := fmt.Sprintf("%s\tYou are admin now!", getReplyAt(update, message))
	botTool.SendMessage(message, str, true, "Markdown")
}

func User(update *tgbotapi.Update, message *tgbotapi.Message) {
	if !config.IsAdmin(message.From.ID) {
		str := fmt.Sprintf("%s\tYou are not admin!", getAt(update, message))
		botTool.SendMessage(message, str, true, "Markdown")
		return
	}
	var str string
	arr := strings.Split(message.Text, " ")
	if len(arr) == 1 {
		str = "Usage: /user [userId]"
		botTool.SendMessage(message, str, true)
		return
	}
	if len(arr[1]) > 1 && arr[1][0] == '@' {
		result := config.CheckUser2Id(arr[1][1:])
		if result[0] == "" {
			str = "User not found"
		} else {
			str = fmt.Sprintf("User found:\nId: %s\nNickName: %s\nLast Message Time: %s", result[0], result[1], result[2])
		}
	} else {
		result := config.CheckId2User(arr[1])
		if result[0] == "" && result[1] == "" {
			str = "User not found"
		} else {
			str = fmt.Sprintf("User found:\nUserName: @%s\nNickName: %s\nLast Message Time: %s", result[0], result[1], result[2])
		}
	}
	botTool.SendMessage(message, str, true)
}
