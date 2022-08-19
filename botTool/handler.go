//This is a function to handle the message
//and to judge the message is a command or not
//Then decide to call the right function

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

// HandleFunc is a type which is users must implement for a text or a command
type HandleFunc func(*tgbotapi.Update)

type IHandler interface {
	match(*tgbotapi.Update)
}

// A type to handle the command
type CommandHandler struct {
	funcs       map[string]HandleFunc
	Msgs        map[string]string
	defultFuncs []HandleFunc
}

// To call the function
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

func (h *CommandHandler) match(update tgbotapi.Update) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("in handler.go/match", err)
		}
	}()
	var data string
	if update.Message != nil {
		data = update.Message.Text
	} else if update.CallbackQuery != nil {
		data = update.CallbackQuery.Data
	}
	command := strings.SplitN(data, " ", 2)[0]
	arr := strings.SplitN(command, "@", 2)
	if len(arr) == 2 && arr[1] != Bot.Self.UserName {
		return
	}
	h.call(arr[0], &update)
	for _, f := range h.defultFuncs {
		go f(&update)
	}
}

// To call handle function which must be handled when any message is received
func (t *CommandHandler) defultHandle(callback HandleFunc) {
	if reflect.ValueOf(callback).Type().NumIn() != 1 {
		panic("too less value!")
	}
	t.defultFuncs = append(t.defultFuncs, safe(callback))
}

// textHandler Regexp pattern matching
type matchPair struct {
	key *regexp.Regexp
	fn  HandleFunc
}

// To handle the Text
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

func (h *TextHandler) match(update tgbotapi.Update) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("in handler.go/match", err)
		}
	}()
	data := update.Message.Text
	for _, f := range h.funcs {
		if f.key.MatchString(data) {
			go f.fn(&update)
		}
	}
	for _, f := range h.defultFuncs {
		go f(&update)
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

var Test *tgbotapi.BotAPI

// Init Bot
func Init(KEY string, TEST ...string) (err error) {
	Bot, err = tgbotapi.NewBotAPI(KEY)
	if len(TEST) == 1 {
		Test, _ = tgbotapi.NewBotAPI(TEST[0])
	}
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

// Provide a function to automatically add function to the handler
// if first character is / it will add to commmandHandler
// otherwise add to textHandler
// if it is a blank it will be a TextHandler.defultHandler
// if it is a / only it will be a CommandHandler.defultHandler
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

type serverErrorLogWriter struct{}

func (*serverErrorLogWriter) Write(p []byte) (int, error) {
	m := string(p)
	// https://github.com/golang/go/issues/26918
	if strings.HasPrefix(m, "http: TLS handshake error") {
		return 0, nil
	}
	return len(p), nil
}

func newServerErrorLog() *log.Logger {
	return log.New(&serverErrorLogWriter{}, "", 0)
}

// Polling to get the message
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
		s := &http.Server{
			Addr:     "0.0.0.0:8443",
			ErrorLog: newServerErrorLog(),
		}
		go s.ListenAndServeTLS("cert.pem", "key.pem")
	}

	defer func() {
		if err := recover(); err != nil {
			log.Println("in handler.go/Polling", err)
		}
	}()
	for update := range updates {
		if update.Message != nil {
			// text := update.Message.Text
			// if len(text) == 0 {
			// 	continue
			// }
			if update.Message.IsCommand() || (update.Message.Text != "" && update.Message.Text[0] == '/') {
				go h.CommandHandler.match(update)
			} else {
				go h.TextHandler.match(update)
			}
		} else if update.CallbackQuery != nil {
			go h.CommandHandler.match(update)
		}
	}
}

// wrapped function to make sure the function is safes
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
