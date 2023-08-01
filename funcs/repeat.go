package funcs

import (
	"bot/botTool"
	"container/list"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var previousMessage sync.Map

type messageConfig struct {
	message string
	userID  int64
	msgID   int
}

func search(mc *list.List, msg *messageConfig) bool {
	for e := mc.Front(); e != nil; e = e.Next() {
		v := e.Value.(*messageConfig)
		if v.message == msg.message {
			mc.Remove(e)
			if v.userID != msg.userID {
				return true
			}
		}
	}
	return false
}

func Repeat(update *tgbotapi.Update, message *tgbotapi.Message) {
	var str string
	if message.Text != "" {
		str = message.Text
	} else if message.Sticker != nil {
		str = message.Sticker.FileUniqueID
	} else if message.Caption != "" {
		str = message.Caption
	} else {
		return
	}
	updateMsg := &messageConfig{message: str, userID: message.From.ID, msgID: message.MessageID}
	v, ok := previousMessage.Load(message.Chat.ID)
	if ok {
		mc := v.(*list.List)
		ok := search(mc, updateMsg)
		if ok {
			botTool.SendForward(message.Chat.ID, message.Chat.ID, updateMsg.msgID)
		} else {
			mc.PushBack(updateMsg)
			if mc.Len() > 3 {
				mc.Remove(mc.Front())
			}
		}
	} else {
		mc := list.New()
		mc.PushBack(updateMsg)
		previousMessage.Store(message.Chat.ID, mc)
	}
}
