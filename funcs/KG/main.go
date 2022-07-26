package KG

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Client = http.Client{
	Timeout: time.Second * 8,
}

func List(Bot *tgbotapi.BotAPI, ChatID int64, MessageID int, ID string, Name string, OriginMessageID string) {
	var AnswerCallbackQueryText string
	var Url string = "https://tooltt.com/send/api/v3/search/song?keyword=" + url.QueryEscape(Name)

	if ResponseData, Error := http.NewRequest("GET", Url, nil); Error == nil {
		ResponseData.Header.Set("User-Agent", "Mozilla/5.0")
		ResponseData.Header.Set("RHost", "https://mobilecdn.kugou.com")
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

	if AnswerCallbackQueryText != "" {
		Bot.Request(tgbotapi.CallbackConfig{
			CallbackQueryID: ID,
			Text:            AnswerCallbackQueryText,
			ShowAlert:       true,
		})
	}
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

		Count := len(JsonData["data"].(map[string]interface{})["info"].([]interface{}))
		if Count > 8 {
			Count = 8
		}

		for Index := 0; Index < Count; Index++ {
			Name := JsonData["data"].(map[string]interface{})["info"].([]interface{})[Index].(map[string]interface{})["filename"].(string)
			Hash := JsonData["data"].(map[string]interface{})["info"].([]interface{})[Index].(map[string]interface{})["hash"].(string)

			Key := "[LQ]" + Name
			Value := "/music KGLink " + Hash + " " + OriginMessageID

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

func Link(Bot *tgbotapi.BotAPI, ChatID int64, MessageID int, ID string, OrderSplit []string) {
	var AnswerCallbackQueryText string
	var EditText string
	var Hash string = OrderSplit[2]
	var Name string
	var AudioID string

	Url := "https://tooltt.com/send/yy/index.php?r=play/getdata&mid=" + Hash + "&hash=" + Hash

	for Index := 0; Index < 2; Index++ {
		if ResponseData, Error := http.NewRequest("GET", Url, nil); Error == nil {
			ResponseData.Header.Set("User-Agent", "Mozilla/5.0")
			ResponseData.Header.Set("RHost", "https://wwwapi.kugou.com")
			if Response, Error := Client.Do(ResponseData); Error == nil {
				defer Response.Body.Close()

				if Body, Error := ioutil.ReadAll(Response.Body); Error == nil {
					if Index == 0 {
						if Name, AudioID = GetDataID(Body); Name == "" || AudioID == "" {
							AnswerCallbackQueryText = "音乐文件ID获取失败"
						} else {
							Url += "&album_audio_id=" + AudioID

							EditText = "正在获取 <code>" + Name + "</code> 音乐ID..."
							EditMessage := tgbotapi.NewMessage(ChatID, EditText)
							EditMessage.ParseMode = "HTML"
							msg,_:=Bot.Send(EditMessage)
							MessageID = msg.MessageID
						}
					} else {
						Caption, Link := GetData(Body)
						if Caption == "" || Link == "" {
							EditText = "音乐文件获取失败"
						} else {
							if OriginMessageID, Error := strconv.Atoi(OrderSplit[3]); Error == nil {
								EditText = "已获取到 <code>" + Name + "</code> 音乐文件正在上传..."
								EditMessage := tgbotapi.NewEditMessageText(ChatID, MessageID, EditText)
								EditMessage.ParseMode = "HTML"
								Bot.Send(EditMessage)

								SendAudio := tgbotapi.NewAudio(ChatID, tgbotapi.FileURL(Link))
								SendAudio.ReplyToMessageID = OriginMessageID
								SendAudio.Caption = "[LQ]" + Caption

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
						}
					}
				} else {
					if Index == 0 {
						AnswerCallbackQueryText = "获取ID响应内容错误"
					} else {
						AnswerCallbackQueryText = "获取音乐响应内容错误"
					}
				}
			} else {
				if Index == 0 {
					AnswerCallbackQueryText = "请求ID错误"
				} else {
					AnswerCallbackQueryText = "请求音乐错误"
				}
			}
		} else {
			if Index == 0 {
				AnswerCallbackQueryText = "创建ID请求错误"
			} else {
				AnswerCallbackQueryText = "创建音乐请求错误"
			}
		}
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

func GetDataID(Body []byte) (string, string) {
	defer func() {
		if Error := recover(); Error != nil {
			log.Println("GetDataID", "FatalRrror", Error)
		}
	}()

	var Name string
	var ID string

	var JsonData map[string]interface{}
	if Error := json.Unmarshal(Body, &JsonData); Error == nil {
		Name = JsonData["data"].(map[string]interface{})["audio_name"].(string)
		IDFloat64 := JsonData["data"].(map[string]interface{})["album_audio_id"].(float64)
		ID = strconv.FormatFloat(IDFloat64, 'f', 0, 64)
	}

	return Name, ID
}

func GetData(Body []byte) (string, string) {
	defer func() {
		if Error := recover(); Error != nil {
			log.Println("GetData", "FatalRrror", Error)
		}
	}()

	var Caption string
	var Link string

	var JsonData map[string]interface{}
	if Error := json.Unmarshal(Body, &JsonData); Error == nil {
		Caption = JsonData["data"].(map[string]interface{})["audio_name"].(string)
		Link = JsonData["data"].(map[string]interface{})["play_backup_url"].(string)
	}

	return Caption, Link
}
