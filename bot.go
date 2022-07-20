package main

import (
	"bot/botTool"
	. "bot/funcs"

	"fmt"

	. "bot/config"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron"
)

var KEY string

func init() {
	if len(os.Args) > 1 {
		KEY = TOKEN2
	} else {
		KEY = TOKEN1
	}
}



var h *botTool.Handler

// type hhh *botTool.Tgbotapi.Update
func main() {
	err := botTool.Init(KEY)
	if err != nil {
		panic(err)
	}
	c := cron.New()
	c.AddFunc("0 0 4,12,18,22 * * ?", ScheduleTask)
	c.AddFunc("0 02 4 * * ?", Clear)
	c.Start()
	h = botTool.NewHandler()
	h.HandleFunc("/start", Start, "开始！")
	h.HandleFunc("/history", History, "历史上的今天")
	h.HandleFunc("/his", History)
	h.HandleFunc("/mingyan", Quote, "名人名言")
	h.HandleFunc("/btc", Btc, "获取当前btc价格")
	h.HandleFunc("/xmr", Xmr, "获取当前xmr价格")
	h.HandleFunc("/eth", Eth, "获取当前eth价格")
	h.HandleFunc("/id", GetId, "获取自己的id")
	h.HandleFunc("/help", help, "获取帮助")
	h.HandleFunc("/curl", Curl, "curl")
	h.HandleFunc("/weather", Weather, "获取城市天气")
	h.HandleFunc("/test", test, "测试")
	h.HandleFunc("/s", Short)
	h.HandleFunc("/t", Translate)
	h.HandleFunc("/translate", Translate, "翻译")
	h.HandleFunc("/short", Short, "生成短链接")
	h.HandleFunc("/make", MakePic, "生产词云")
	h.HandleFunc("/tp", Ping, "测试连通性")
	h.HandleFunc("/user", User)
	h.HandleFunc("/music", Music, "搜索音乐")
	h.HandleFunc("/ban", Ban, "")
	//	h.HandleFunc("/html", Html, "Html To Pic")
	h.HandleFunc("^(一言|morning|早上好){1}$", Quote)
	h.HandleFunc(`^[\s\S]*(\(|[|{|<|（|《|【|（|［|｛|＜|『|「|«|\)|]|}|>|）|》|】|］|）|｝|＞|』|」|»)$`, Fix)
	h.HandleFunc("", TextManager)
	h.HandleFunc("/", GetMessgae)
	h.HandleFunc("", GetMessgae)
	h.Polling(BOT_CONFIG)
}

func test(update *tgbotapi.Update) {
	if update.Message.From.ID == 1456780662 {
		fmt.Printf("%+v\n", update.Message.ReplyToMessage.From.ID)
	}
}

func help(update *tgbotapi.Update) {
	var text string
	for i, j := range h.CommandHandler.Msgs {
		text += fmt.Sprintf("%s   %s\n", i, j)
	}
	botTool.SendMessage(update, &text, true)
}
