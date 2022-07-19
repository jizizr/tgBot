package botTool

import (
	"log"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"runtime/debug"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Bot *tgbotapi.BotAPI
var Token string

type Updates struct {
	*tgbotapi.Update
}

type HandleFunc func(*tgbotapi.Update)

type IHandler interface {
	match(*tgbotapi.Update)
}

type CommandHandler struct {
	funcs       map[string]HandleFunc
	Msgs        map[string]string
	defultFuncs []HandleFunc
}

func (c *CommandHandler) call(name string, update *tgbotapi.Update) {
	f, ok := c.funcs[name]
	if !ok {
		return
	}
	f(update)
}

func (c *CommandHandler) handle(command string, callback HandleFunc, msg ...string) {
	if reflect.ValueOf(callback).Type().NumIn() != 1 {
		panic("too less value!")
	}
	c.funcs[command] = safe(callback)
	if len(msg) == 1 {
		c.Msgs[command] = msg[0]
	}
}

func (h *CommandHandler) match(update *tgbotapi.Update) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	var data string
	if update.Message != nil {
		data = update.Message.Text
	} else if update.CallbackQuery != nil {
		data = update.CallbackQuery.Data
	}
	arr := strings.Split(data, "@")
	if len(arr) == 2 {
		arr1 := strings.Split(arr[1], " ")
		if arr1[0] != Bot.Self.UserName {
			return
		}
	} else {
		arr = strings.Split(data, " ")
	}
	h.call(arr[0], update)
	for _, f := range h.defultFuncs {
		go f(update)
	}
}

func (t *CommandHandler) defultHandle(callback HandleFunc) {
	if reflect.ValueOf(callback).Type().NumIn() != 1 {
		panic("too less value!")
	}
	t.defultFuncs = append(t.defultFuncs, safe(callback))
}

type matchPair struct {
	key *regexp.Regexp
	fn  HandleFunc
}

type TextHandler struct {
	funcs       []*matchPair
	defultFuncs []HandleFunc
}

func (c *TextHandler) handle(command string, callback HandleFunc) {
	if reflect.ValueOf(callback).Type().NumIn() != 1 {
		panic("too less value!")
	}
	c.funcs = append(c.funcs, &matchPair{
		key: regexp.MustCompile(command),
		fn:  safe(callback),
	})
}

func (t *TextHandler) defultHandle(callback HandleFunc) {
	if reflect.ValueOf(callback).Type().NumIn() != 1 {
		panic("too less value!")
	}
	t.defultFuncs = append(t.defultFuncs, safe(callback))
}

func (h *TextHandler) match(update *tgbotapi.Update) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	data := update.Message.Text
	for _, f := range h.funcs {
		if f.key.MatchString(data) {
			go f.fn(update)
		}
	}
	for _, f := range h.defultFuncs {
		go f(update)
	}
}

func NewCommandHandler() *CommandHandler {
	return &CommandHandler{
		funcs:       map[string]HandleFunc{},
		Msgs:        map[string]string{},
		defultFuncs: []HandleFunc{},
	}
}

func NewTextHandler() *TextHandler {
	return &TextHandler{
		funcs:       []*matchPair{},
		defultFuncs: []HandleFunc{},
	}
}

func Init(KEY string) (err error) {
	Bot, err = tgbotapi.NewBotAPI(KEY)
	Token = KEY
	return
}

type Handler struct {
	CommandHandler *CommandHandler
	TextHandler    *TextHandler
}

func NewHandler() *Handler {
	c := NewCommandHandler()
	t := NewTextHandler()
	return &Handler{CommandHandler: c, TextHandler: t}
}

func (h *Handler) HandleFunc(command string, callback HandleFunc, msg ...string) {
	if command == "" {
		h.TextHandler.defultHandle(callback)
	} else if command[0] == '/' {
		if len(command) == 1 {
			h.CommandHandler.defultHandle(callback)
		} else {
			h.CommandHandler.handle(command, callback, msg...)
		}
	} else {
		h.TextHandler.handle(command, callback)
	}
}

func (h *Handler) Polling(CONFIG string) {
	Bot.Debug = false
	var updates tgbotapi.UpdatesChannel
	if len(os.Args) > 1 {
		updateConfig := tgbotapi.NewUpdate(0)
		updateConfig.Timeout = 60
		updates = Bot.GetUpdatesChan(updateConfig)
	} else {

		wh, _ := tgbotapi.NewWebhook(CONFIG + Token)

		Bot.Request(wh)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// info, err := Bot.GetWebhookInfo()
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// if info.LastErrorDate != 0 {
		// 	log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
		// }

		updates = Bot.ListenForWebhook("/" + Token)
		go http.ListenAndServeTLS("0.0.0.0:8443", "cert.pem", "key.pem", nil)
	}

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	for update := range updates {
		if update.Message != nil {
			// text := update.Message.Text
			// if len(text) == 0 {
			// 	continue
			// }
			if update.Message.IsCommand() {
				go h.CommandHandler.match(&update)
			} else {
				go h.TextHandler.match(&update)
			}
		} else if update.CallbackQuery != nil {
			go h.CommandHandler.match(&update)
		}
	}
}

func safe(f HandleFunc) HandleFunc {
	return func(update *tgbotapi.Update) {
		defer func() {
			if e, ok := recover().(error); ok {
				log.Printf("Panic in %s,error:%v\n", runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(), e)
				log.Println(string(debug.Stack()))
			}
		}()
		// log.Printf("call %v", f)
		f(update)
	}
}
