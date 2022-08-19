package funcs

import (
	"bot/botTool"
	"bot/funcs/KG"
	"bot/funcs/KW"
	"encoding/base64"
	"net/http"
	"strconv"
	"time"

	StringSplit "github.com/UallenQbit/GoLang-SplitString"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Client = http.Client{
	Timeout: time.Second * 8,
}

func Music(update *tgbotapi.Update) {
	var SplitData []string
	var SplitCount int
	var ChatID int64
	var MessageID int
	var ID string
	var CallbackText string

	if update.Message == nil {
		if update.EditedMessage == nil {
			if update.CallbackQuery != nil {
				if SplitData, SplitCount = StringSplit.SplitString(update.CallbackQuery.Data, " "); SplitCount != 0 {
					ChatID = update.CallbackQuery.Message.Chat.ID
					MessageID = update.CallbackQuery.Message.MessageID
					ID = update.CallbackQuery.ID
					CallbackText = update.CallbackQuery.Message.Text
				}
			}
		} else {
			if SplitData, SplitCount = StringSplit.SplitString(update.EditedMessage.Text, " "); SplitCount != 0 {
				ChatID = update.EditedMessage.Chat.ID
				MessageID = update.EditedMessage.MessageID
			}
		}
	} else {
		if SplitData, SplitCount = StringSplit.SplitString(update.Message.Text, " "); SplitCount != 0 {
			ChatID = update.Message.Chat.ID
			MessageID = update.Message.MessageID
		}
	}
	if SplitCount > 1 {
		if ID == "" {
			var Name string

			for Index := 1; Index < SplitCount; Index++ {
				if Index == (SplitCount - 1) {
					Name += SplitData[Index]
				} else {
					Name += SplitData[Index] + " "
				}
			}

			NameBase64 := base64.StdEncoding.EncodeToString([]byte(Name))

			SendText := "请选择 <code>" + Name + "</code> 音乐搜索平台"
			SendMessage := tgbotapi.NewMessage(ChatID, SendText)
			SendMessage.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("酷我", "/music KW "+NameBase64+" "+strconv.Itoa(MessageID)),
					tgbotapi.NewInlineKeyboardButtonData("酷狗", "/music KG "+NameBase64+" "+strconv.Itoa(MessageID)),
				),
			)
			SendMessage.ReplyToMessageID = MessageID
			SendMessage.ParseMode = "HTML"

			if _, Error := botTool.Bot.Send(SendMessage); Error != nil {
				SendMessage := tgbotapi.NewMessage(ChatID, "搜索音乐名太长")
				SendMessage.ReplyToMessageID = MessageID
				SendMessage.ParseMode = "HTML"
			}
		} else {
			Type := SplitData[1]
			if Type == "QQ" || Type == "KG" || Type == "KW" {
				if NameByte, Error := base64.StdEncoding.DecodeString(SplitData[2]); Error == nil {
					Name := string(NameByte)
					OriginMessageID := SplitData[3]

					if Type == "KG" {
						KG.List(botTool.Bot, ChatID, MessageID, ID, Name, OriginMessageID)
					} else if Type == "KW" {
						KW.List(botTool.Bot, ChatID, MessageID, ID, Name, OriginMessageID)
					}
				} else {
					botTool.Bot.Request(tgbotapi.CallbackConfig{
						CallbackQueryID: ID,
						Text:            "数据异常",
						ShowAlert:       true,
					})
				}
			} else {
				if Type == "KGLink" {
					KG.Link(botTool.Bot, ChatID, MessageID, ID, SplitData)
				} else if Type == "KWLink" {
					KW.Link(botTool.Bot, ChatID, MessageID, ID, SplitData, CallbackText)
				}
			}
		}
	} else {
		SendText := "请输入查询音乐名\n"
		SendText += "例: <code>/music 鹿 be free</code>"

		SendMessage := tgbotapi.NewMessage(ChatID, SendText)
		SendMessage.ReplyToMessageID = MessageID
		SendMessage.ParseMode = "HTML"
	}
}
