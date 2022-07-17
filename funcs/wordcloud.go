package funcs

import (
	"bot/botTool"
	"bot/dbManager"
	group "bot/wdCloud"
	"fmt"

	// "regexp"
	. "bot/config"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yanyiwu/gojieba"
)

// var re3, _ = regexp.Compile(`[\p{P}\s]*`)
var jieba = gojieba.NewJieba()
var db = dbManager.InitMysql("data", DB_TOKEN, "data")
var cx = map[string]struct{}{"eng": {}, "v": {}, "l": {}, "x": {}, "n": {}, "nr": {}, "a": {}, "vd": {}, "nz": {}, "PER": {}, "f": {}, "ns": {}, "LOC": {}, "s": {}, "nt": {}, "ORG": {}, "nw": {}, "vn": {}}

func TextManager(update *tgbotapi.Update) {
	if update.Message.From.IsBot || update.Message.From.ID == 777000 || update.Message.IsCommand() {
		return
	}
	text := update.Message.Text
	userId := fmt.Sprint(update.Message.From.ID)
	chatId := fmt.Sprint(update.Message.Chat.ID)
	name := getName(update)
	db.AddUser(chatId, userId, name)
	// text = re3.ReplaceAllString(text, "")
	// config.AddGroup(chatId, update.Message.Chat.UserName, update.Message.Chat.Title,fmt.Sprint(update.Message.From.ID),update.Message.From.UserName,getName(update))
	if len(text) < 2 {
		return
	} else if len(text) < 7 {
		text = strings.Join(jieba.CutForSearch(text, true), " ")
	}
	word := jieba.Tag(text)
	for _, v := range word {
		w := strings.Split(v, "/")
		// log.Println(w)
		if len(w[0]) > 3 && botTool.Contains(cx, w[1]) {
			// log.Println(w[0])
			// log.Println(len("中"))
			db.AddMessage(chatId, w[0])
		}
	}

}

func getPic(chatId string) {
	chatId2 := fmt.Sprintf("%sGroup", chatId)
	result := db.GetAllWords(&chatId2)
	botTool.SendPhoto(chatId, group.Rank(result))
}

func Clear() {
	db.Clear()
}

func ScheduleTask() {
	groups := make([]string, 0)
	db.TableInfo(&groups)
	for _, v := range groups {
		getPic(v)
		getUsers(v)
	}
}

func getUsers(chatId string) {
	result := db.GetAllUsers(&chatId)
	users := result[1]
	times := result[0]
	top5Users := make([]string, 0)
	for i := 0; i < len(users); i++ {
		user := users[i]
		if len(user) > 5 {
			user = strings.TrimSpace(strings.Split(user, "|")[0])
		}
		if len(user) > 5 {
			user = strings.TrimSpace(strings.Split(user, " ")[0])
		}
		if len(user) > 5 {
			user = user[:6]
		}
		top5Users = append(top5Users, fmt.Sprintf("\t\t🎖`%s` 呱唧了:`%s`句\n", user, times[i]))
	}
	text := fmt.Sprintf(`🏵 今日活跃用户排行榜 🏵
  📅 %s
  ⏱ 截至今天 %s

%s
  感谢这些朋友的哔哔赖赖! 👏 
  遇到问题,向他们请教说不定会吃ban呢😃`, time.Now().Format("`2006-01-02`"), time.Now().Format("`15:04`"), strings.Join(top5Users, ""))
	id, _ := strconv.ParseInt(chatId, 10, 64)
	msg := tgbotapi.NewMessage(id, text)
	msg.ParseMode = "Markdown"
	botTool.Bot.Send(msg)
}
