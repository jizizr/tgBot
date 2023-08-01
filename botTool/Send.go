package botTool

import (
	// "bufio"
	// "log"

	"math/rand"
	"strconv"

	// "os"
	"time"
	"unsafe"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var src = rand.NewSource(time.Now().UnixNano())

const (
	// 6 bits to represent a letter index
	letterIdBits = 6
	// All 1-bits as many as letterIdBits
	letterIdMask = 1<<letterIdBits - 1
	letterIdMax  = 63 / letterIdBits
)

func randStr(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdMax letters!
	for i, cache, remain := n-1, src.Int63(), letterIdMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdMax
		}
		if idx := int(cache & letterIdMask); idx < len(letters) {
			b[i] = letters[idx]
			i--
		}
		cache >>= letterIdBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}

func SendMessage(message *tgbotapi.Message, text string, reply bool, mode ...string) (*tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	if reply {
		msg.ReplyToMessageID = message.MessageID
	}
	if len(mode) > 0 {
		msg.ParseMode = mode[0]
	}
	msgConfig, error := Bot.Send(msg)
	return &msgConfig, error
}

func Edit(msg *tgbotapi.Message, text string, mode ...string) (*tgbotapi.Message, error) {
	editMessage := tgbotapi.EditMessageTextConfig{
		BaseEdit:              tgbotapi.BaseEdit{ChatID: msg.Chat.ID, MessageID: msg.MessageID},
		Text:                  text,
		Entities:              []tgbotapi.MessageEntity{},
		DisableWebPagePreview: false,
	}
	if len(mode) > 0 {
		editMessage.ParseMode = mode[0]
	}
	msgConfig, error := Bot.Send(editMessage)
	return &msgConfig, error
}

func SendDocument(message *tgbotapi.Message, text []byte, name string, reply bool, caption ...string) (*tgbotapi.Message, error) {
	// newName := randStr(10) + name
	// file, err := os.OpenFile(newName, os.O_WRONLY|os.O_CREATE, 0666)
	// if err != nil {
	// 	log.Println("文件打开失败", err)
	// }
	// //及时关闭file句柄
	// defer file.Close()
	// //写入文件时，使用带缓存的 *Writer
	// write := bufio.NewWriter(file)
	// write.WriteString(*text)
	// //Flush将缓存的文件真正写入到文件中
	// write.Flush()
	// f, err := os.Open(newName)
	// if err != nil {
	// 	log.Println("文件打开失败", err)
	// }
	// updateFile := tgbotapi.FileReader{Name: name, Reader: f}
	// document := tgbotapi.NewDocument(message.Chat.ID, updateFile)

	// // log.Println(document.Caption)
	config := tgbotapi.FileBytes{
		Name:  name,
		Bytes: text,
	}
	document := tgbotapi.NewDocument(message.Chat.ID, config)
	if reply {
		document.ReplyToMessageID = message.MessageID
	}
	if len(caption) > 0 {
		document.Caption = caption[0]
	}
	msg, err := Bot.Send(document)
	// defer f.Close()
	// os.Remove(newName)
	return &msg, err
}

func SendPhoto(chatId string, data []byte) {
	config := tgbotapi.FileBytes{
		Name:  "",
		Bytes: data,
	}
	id, _ := strconv.ParseInt(chatId, 10, 64)
	photo := tgbotapi.NewPhoto(id, config)
	Bot.Send(photo)
}

func SendForward(chatId int64, fromChatID int64, msgID int) {
	forward := tgbotapi.NewForward(chatId, fromChatID, msgID)
	Bot.Send(forward)
}

func SendFile(message *tgbotapi.Message, url string, reply bool, caption ...string) {
	file := tgbotapi.FileURL(url)
	photo := tgbotapi.NewPhoto(message.Chat.ID, file)
	if reply {
		photo.ReplyToMessageID = message.MessageID
	}
	if len(caption) == 1 {
		photo.Caption = caption[0]
	} else if len(caption) == 2 {
		photo.Caption = caption[0]
		photo.ParseMode = caption[1]
	}
	Bot.Send(photo)
}
