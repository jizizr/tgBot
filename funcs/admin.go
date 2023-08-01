package funcs

import (
	"bot/botTool"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Ban(update *tgbotapi.Update, message *tgbotapi.Message) {
	var sec int64
	if !checkAdmin(message) {
		BanPlayer(update, message)
		return
	}
	if message.ReplyToMessage == nil {
		str := "请回复消息"
		botTool.SendMessage(message, str, true)
		return
	}
	arr := strings.Fields(message.Text)
	if len(arr) < 2 {
		sec = 60
	} else {
		sec, _ = strconv.ParseInt(arr[1], 10, 64)
	}
	err := botTool.BanMember(update, message, message.ReplyToMessage.From.ID, sec)
	if err != nil {
		return
	}

	str := fmt.Sprintf("%s 已被禁言 %d 秒", getReplyAt(update, message), sec)
	botTool.SendMessage(message, str, true, "Markdown")
}

func BanPlayer(update *tgbotapi.Update, message *tgbotapi.Message) {
	gid := message.Chat.ID

	uid := message.From.ID
	botme, _ := botTool.Bot.GetChatMember(tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: gid,
			UserID: botTool.Bot.Self.ID,
		},
	})
	var str string
	if botme.CanRestrictMembers {
		botTool.BanMember(update, message, uid, 60)
		str = "[" + botTool.GetName(update, message) + "](tg://user?id=" + fmt.Sprint(uid) + ")乱玩管理员命令,禁言一分钟"
	} else {
		str = "[" + botTool.GetName(update, message) + "](tg://user?id=" + fmt.Sprint(uid) + ")不要乱玩管理员命令"
	}
	botTool.SendMessage(message, str, true, "Markdown")
}
