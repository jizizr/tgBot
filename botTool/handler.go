//This is a function to handle the message
//and to judge the message is a command or not
//Then decide to call the right function

package botTool

import (
	. "bot/config"
	"bytes"
	"encoding/json"
	"fmt"
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
type HandleFunc func(*tgbotapi.Update, *tgbotapi.Message)

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
func (c *CommandHandler) call(name string, update *tgbotapi.Update, message *tgbotapi.Message) {
	f, ok := c.funcs[name]
	if !ok {
		return
	}
	f(update, message)
}

func (c *CommandHandler) handle(command string, callback HandleFunc, msg ...string) {
	if reflect.ValueOf(callback).Type().NumIn() != 2 {
		panic("too less value!")
	}
	c.funcs[command] = safe(callback)
	if len(msg) == 1 {
		c.Msgs[command] = msg[0]
	}
}

func (h *CommandHandler) match(update *tgbotapi.Update, message *tgbotapi.Message) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("in handler.go/match", err)
		}
	}()
	var data string
	if message != nil {
		data = message.Text
	} else {
		data = update.CallbackData()
	}
	command := strings.SplitN(data, " ", 2)[0]
	arr := strings.SplitN(command, "@", 2)
	if len(arr) == 2 && arr[1] != Bot.Self.UserName {
		return
	}
	h.call(arr[0], update, message)
	for _, f := range h.defultFuncs {
		go f(update, message)
	}
}

// To call handle function which must be handled when any message is received
func (t *CommandHandler) defultHandle(callback HandleFunc) {
	if reflect.ValueOf(callback).Type().NumIn() != 2 {
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
	if reflect.ValueOf(callback).Type().NumIn() != 2 {
		panic("too less value!")
	}
	c.funcs = append(c.funcs, &matchPair{
		key: regexp.MustCompile(command),
		fn:  safe(callback),
	})
}

func (t *TextHandler) defultHandle(callback HandleFunc) {
	if reflect.ValueOf(callback).Type().NumIn() != 2 {
		panic("too less value!")
	}
	t.defultFuncs = append(t.defultFuncs, safe(callback))
}

func (h *TextHandler) match(update *tgbotapi.Update, message *tgbotapi.Message) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("in handler.go/match", err)
		}
	}()
	for _, f := range h.funcs {
		if f.key.MatchString(message.Text) {
			go f.fn(update, message)
		}
	}
	for _, f := range h.defultFuncs {
		go f(update, message)
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

func (h *Handler) handleUpdate(update *tgbotapi.Update) {
	var message *tgbotapi.Message
	defer func() {
		if e, ok := recover().(error); ok {
			bs, _ := json.Marshal(update)
			var out bytes.Buffer
			json.Indent(&out, bs, "", "    ")
			str := fmt.Sprintf("%v", out.String())
			errMsg := fmt.Sprintf("Panic in Polling,error:%v\n\n%s\n\n%s", e, string(debug.Stack()), str)
			if len(os.Args) > 1 {
				log.Println(errMsg)
			} else {
				Bot.Send(tgbotapi.NewMessage(CHAT_ID, errMsg))
				log.Println(errMsg)
			}
		}
	}()
	if update.Message != nil {
		message = update.Message
	} else if update.EditedMessage != nil {
		message = update.EditedMessage
	} else if update.CallbackQuery != nil {
		go h.CommandHandler.match(update, nil)
		return
	} else {
		return
	}
	if message.Text != "" && message.Text[0] == '/' {
		go h.CommandHandler.match(update, message)
	} else {
		go h.TextHandler.match(update, message)
	}
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

	for {
		update := <-updates
		go h.handleUpdate(&update)
	}
}

// wrapped function to make sure the function is safes
func safe(f HandleFunc) HandleFunc {
	return func(update *tgbotapi.Update, message *tgbotapi.Message) {
		defer func() {
			if e, ok := recover().(error); ok {
				errMsg := fmt.Sprintf("Panic in %s,error:%v\n\n%s\n\nin @%s %d %s\n%s\n%d @%s", runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(), e, string(debug.Stack()), message.Chat.UserName, message.Chat.ID, message.Chat.Title, message.Text, message.From.ID, message.From.UserName)
				if len(os.Args) > 1 {
					log.Println(errMsg)
				} else {
					Bot.Send(tgbotapi.NewMessage(CHAT_ID, errMsg))
					log.Println(errMsg)
					bs, _ := json.Marshal(message)
					var out bytes.Buffer
					json.Indent(&out, bs, "", "    ")
					errMsg = fmt.Sprintf("%v", out.String())
					Bot.Send(tgbotapi.NewMessage(CHAT_ID, errMsg))
				}
			}
		}()
		// log.Printf("call %v", f)
		f(update, message)
	}
}
