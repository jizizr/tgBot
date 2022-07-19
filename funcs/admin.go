package funcs

import (
	"bot/botTool"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Ban(update *tgbotapi.Update) {
	var sec int64
	if update.Message.ReplyToMessage == nil {
		str := "请回复消息"
		botTool.SendMessage(update, &str, true)
		return
	}
	arr := strings.Split(update.Message.Text, " ")
	if len(arr) < 2 {
		sec = 60
	} else {
		sec, _ = strconv.ParseInt(arr[1], 10, 64)
	}
	err := botTool.BanMember(update, update.Message.Chat.ID, update.Message.ReplyToMessage.From.ID, sec)
	if err != nil {
		return
	}

	str := "[" + botTool.GetReplyName(update) + "](tg://user?id=" + fmt.Sprint(update.Message.ReplyToMessage.From.ID) + ") 已禁言" + fmt.Sprint(sec) + "秒"
	botTool.SendMessage(update, &str, true, "Markdown")
}

func BanPlayer(update *tgbotapi.Update) {
	gid := update.Message.Chat.ID

	uid := update.Message.ReplyToMessage.From.ID
	botme, _ := botTool.Bot.GetChatMember(tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: gid,
			UserID: botTool.Bot.Self.ID,
		},
	})
	var str string
	if botme.CanRestrictMembers {
		botTool.BanMember(update, gid, uid, 60)
		str = "[" + botTool.GetName(update) + "](tg://user?id=" + fmt.Sprint(uid) + ")乱玩管理员命令,禁言一分钟"
	} else {
		str = "[" + botTool.GetName(update) + "](tg://user?id=" + fmt.Sprint(uid) + ")不要乱玩管理员命令"
	}
	botTool.SendMessage(update, &str, true, "Markdown")
}
