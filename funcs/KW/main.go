package KW

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Client = http.Client{
	Timeout: time.Second * 8,
}

func List(Bot *tgbotapi.BotAPI, ChatID int64, MessageID int, ID string, Name string, OriginMessageID string) {
	var AnswerCallbackQueryText string
	var Token string = GetToken()

	if Token == "" {
		AnswerCallbackQueryText = "Token获取失败"
	} else {
		var Url string = "https://www.kuwo.cn/api/www/search/searchMusicBykeyWord?rn=8&key=" + url.QueryEscape(Name)

		if ResponseData, Error := http.NewRequest("GET", Url, nil); Error == nil {
			ResponseData.Header.Set("User-Agent", "Mozilla/5.0")
			ResponseData.Header.Set("Referer", "https://www.kuwo.cn/api/www/search/searchMusicBykeyWord")
			ResponseData.Header.Set("Csrf", Token)
			ResponseData.Header.Set("Cookie", "kw_token="+Token)
			if Response, Error := Client.Do(ResponseData); Error == nil {
				defer Response.Body.Close()

				if Body, Error := ioutil.ReadAll(Response.Body); Error == nil {
					if InlineKeyboardButton, OK := GetListInlineKeyboardButton(Body, OriginMessageID); OK == "OK" {
						EditText := "请选择 <code>" + Name + "</code> 音乐"
						EditMessage := tgbotapi.NewEditMessageTextAndMarkup(ChatID, MessageID, EditText, InlineKeyboardButton)
						EditMessage.ParseMode = "HTML"
						if _, Error := Bot.Send(EditMessage); Error != nil {
							AnswerCallbackQueryText = "音乐列表发送失败"
						}
					} else {
						AnswerCallbackQueryText = "音乐信息搜索失败"
					}
				} else {
					AnswerCallbackQueryText = "获取响应内容错误"
				}
			} else {
				AnswerCallbackQueryText = "请求错误"
			}
		} else {
			AnswerCallbackQueryText = "创建请求错误"
		}
	}

	if AnswerCallbackQueryText != "" {
		Bot.Request(tgbotapi.CallbackConfig{
			CallbackQueryID: ID,
			Text:            AnswerCallbackQueryText,
			ShowAlert:       true,
		})
	}
}

func GetToken() string {
	var Token string
	var Url string = "https://www.kuwo.cn"

	if ResponseData, Error := http.NewRequest("GET", Url, nil); Error == nil {
		ResponseData.Header.Set("User-Agent", "Mozilla/5.0")
		if Response, Error := Client.Do(ResponseData); Error == nil {
			defer Response.Body.Close()

			Token = strings.Split(Response.Header.Get("Set-Cookie")[9:], ";")[0]
		}
	}

	return Token
}

func GetListInlineKeyboardButton(Body []byte, OriginMessageID string) (tgbotapi.InlineKeyboardMarkup, string) {
	defer func() {
		if Error := recover(); Error != nil {
			log.Println("GetData", "FatalRrror", Error)
		}
	}()

	var InlineKeyboardButton [][]tgbotapi.InlineKeyboardButton

	var JsonData map[string]interface{}
	if Error := json.Unmarshal(Body, &JsonData); Error == nil {

		Count := len(JsonData["data"].(map[string]interface{})["list"].([]interface{}))
		if Count > 8 {
			Count = 8
		}

		for Index := 0; Index < Count; Index++ {
			Name := JsonData["data"].(map[string]interface{})["list"].([]interface{})[Index].(map[string]interface{})["name"].(string)
			IDFloat64 := JsonData["data"].(map[string]interface{})["list"].([]interface{})[Index].(map[string]interface{})["rid"].(float64)

			Name = strings.Replace(Name, "&nbsp;", " ", -1)
			ID := strconv.FormatFloat(IDFloat64, 'f', 10, 64)

			Key := "[LQ]" + Name
			Value := "/music KWLink " + ID + " " + OriginMessageID

			NewInlineKeyboardRow := tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					Key,
					Value,
				),
			)

			InlineKeyboardButton = append(InlineKeyboardButton, NewInlineKeyboardRow)
		}
	}

	InlineKeyboardMarkup := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: InlineKeyboardButton,
	}

	return InlineKeyboardMarkup, "OK"
}

func Link(Bot *tgbotapi.BotAPI, ChatID int64, MessageID int, ID string, OrderSplit []string, CallbackText string) {
	var AnswerCallbackQueryText string
	var EditText string
	var RID string = OrderSplit[2]

	Url := "https://antiserver.kuwo.cn/anti.s?type=convert_url&format=mp3&response=url&rid=" + RID

	if ResponseData, Error := http.NewRequest("GET", Url, nil); Error == nil {
		ResponseData.Header.Set("User-Agent", "Mozilla/5.0")
		ResponseData.Header.Set("RHost", "https://wwwapi.kugou.com")
		if Response, Error := Client.Do(ResponseData); Error == nil {
			defer Response.Body.Close()

			if Body, Error := ioutil.ReadAll(Response.Body); Error == nil {
				if OriginMessageID, Error := strconv.Atoi(OrderSplit[3]); Error == nil {
					Name := strings.Replace(strings.Replace(CallbackText, "请选择 ", "", -1), " 音乐", "", -1)
					EditText = "已获取到 <code>" + Name + "</code> 音乐文件正在上传..."
					EditMessage := tgbotapi.NewMessage(ChatID, EditText)
					EditMessage.ParseMode = "HTML"
					msg, _ := Bot.Send(EditMessage)
					MessageID = msg.MessageID
					SendAudio := tgbotapi.NewAudio(ChatID, tgbotapi.FileURL(string(Body)))
					SendAudio.ReplyToMessageID = OriginMessageID
					SendAudio.Caption = "[LQ]" + Name

					if _, Error := Bot.Send(SendAudio); Error == nil {
						if _, Error := Bot.Send(tgbotapi.NewDeleteMessage(ChatID, MessageID)); Error == nil {
							EditText = ""
						} else {
							EditText = "等等信息删除错误"
						}
					} else {
						AnswerCallbackQueryText = "音乐文件发送失败"
					}
				} else {
					EditText = "原信息ID错误"
				}
			} else {
				AnswerCallbackQueryText = "获取音乐响应内容错误"

			}
		} else {
			AnswerCallbackQueryText = "请求音乐错误"
		}
	} else {
		AnswerCallbackQueryText = "获取音乐响应内容错误"
	}

	if AnswerCallbackQueryText != "" {
		Bot.Request(tgbotapi.CallbackConfig{
			CallbackQueryID: ID,
			Text:            AnswerCallbackQueryText,
			ShowAlert:       true,
		})
	}

	if EditText != "" {
		EditMessage := tgbotapi.NewEditMessageText(ChatID, MessageID, EditText)
		EditMessage.ParseMode = "HTML"
	}
}

func GetDataHash(Body []byte) string {
	defer func() {
		if Error := recover(); Error != nil {
			log.Println("GetDataHash", "FatalRrror", Error)
		}
	}()

	var Hash string

	var JsonData map[string]interface{}
	if Error := json.Unmarshal(Body, &JsonData); Error == nil {
		Hash = JsonData["data"].(map[string]interface{})["info"].([]interface{})[0].(map[string]interface{})["hash"].(string)
	}

	return Hash
}
