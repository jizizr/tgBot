package botTool

import (
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Contains(slice map[string]struct{}, item string) bool {
	_, ok := slice[item]
	return ok
}

func BanMember(update *tgbotapi.Update, gid int64, uid int64, sec int64) error {
	if sec <= 0 {
		sec = 9999999999999
	}
	chatuserconfig := tgbotapi.ChatMemberConfig{ChatID: gid, UserID: uid}
	restricconfig := tgbotapi.RestrictChatMemberConfig{
		ChatMemberConfig: chatuserconfig,
		UntilDate:        time.Now().Unix() + sec,
	}
	_, err := Bot.Request(restricconfig)
	if err != nil {
		var str string
		if strings.Contains(err.Error(), "can't restrict self") {
			str = "你想禁言我？"
		} else if strings.Contains(err.Error(), "user is an administrator of the chat") {
			str = "对面是管理！"
		} else {
			str = "无权禁言！"
		}
		SendMessage(update, &str, true)
	}
	return err
}

func GetName(update *tgbotapi.Update) (name string) {
	user := update.Message.From
	name = user.FirstName + " " + user.LastName
	if name != " " {
		return
	}
	name = user.UserName
	if name != "" {
		return
	}
	name = string(rune(user.ID))
	return
}

func GetReplyName(update *tgbotapi.Update) (name string) {
	user := update.Message.ReplyToMessage.From
	name = user.FirstName + " " + user.LastName
	if name != " " {
		return
	}
	name = user.UserName
	if name != "" {
		return
	}
	name = string(rune(user.ID))
	return
}
