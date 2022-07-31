package funcs

import (
	"bot/botTool"
	"container/list"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var previousMessage sync.Map

type messageConfig struct {
	message *string
	userID  int64
	msgID   int
}

func search(mc *list.List, msg *messageConfig) bool {
	for e := mc.Front(); e != nil; e = e.Next() {
		v := e.Value.(*messageConfig)
		if *v.message == *msg.message {
			mc.Remove(e)
			if v.userID != msg.userID {
				return true
			}
		}
	}
	return false
}

func Repeat(update *tgbotapi.Update) {
	var str *string
	if update.Message.Text != "" {
		str = &update.Message.Text
	} else if update.Message.Sticker != nil {
		str = &update.Message.Sticker.FileUniqueID
	} else if update.Message.Caption != "" {
		str = &update.Message.Caption
	} else {
		return
	}
	updateMsg := &messageConfig{message: str, userID: update.Message.From.ID, msgID: update.Message.MessageID}
	v, ok := previousMessage.Load(update.Message.Chat.ID)
	if ok {
		mc := v.(*list.List)
		ok := search(mc, updateMsg)
		if ok {
			botTool.SendForward(update.Message.Chat.ID, update.Message.Chat.ID, updateMsg.msgID)
		} else {
			mc.PushBack(updateMsg)
			if mc.Len() > 3 {
				mc.Remove(mc.Front())
			}
		}
	} else {
		mc := list.New()
		mc.PushBack(updateMsg)
		previousMessage.Store(update.Message.Chat.ID, mc)
	}
}
