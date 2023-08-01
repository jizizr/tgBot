package funcs

import (
	"bot/botTool"
	. "bot/config"
	"errors"
	"fmt"
	"net/rpc"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var connections map[string]*rpc.Client

func splitfunc(r rune) bool {
	return r == ' ' || r == ':'
}

func init() {
	connections = make(map[string]*rpc.Client)
	for i, j := range TP_URLS {
		conn, err := Dial("tcp", j)
		if err != nil {
			connections[i] = nil
		} else {
			connections[i] = conn
		}
	}
}

var ipMatch = regexp.MustCompile(`(https?://|\s|^)([^:\./\s]+\.)+[^\./:\s"]+((:|\s)([1-5]\d{4}|6[0-4]\d{3}|65[0-4]\d{2}|655[0-2]\d|6553[0-5]|\d{1,4}))?`)

type urls struct {
	pre  string
	host string
}

func (u *urls) toUrlStr() string {
	return u.pre + u.host
}

func parseUrl(textRaw string) (urls, error) {
	textRaw = strings.ToLower(strings.TrimRight(textRaw, "/"))
	textMatch := ipMatch.FindString(textRaw)
	if textMatch == "" {
		return urls{}, errors.New("文本未包含可测试url")
	}
	textMatch = strings.ReplaceAll(strings.TrimSpace(textMatch), " ", ":")
	if strings.HasPrefix(textMatch, "https://") {
		return urls{pre: "https://", host: textMatch[8:]}, nil
	} else if strings.HasPrefix(textMatch, "http://") {
		return urls{pre: "http://", host: textMatch[7:]}, nil
	} else {
		return urls{pre: "tcp://", host: textMatch}, nil
	}
}

func parseText(textRaw string) (urlMsg string, err error) {
	textRaw = strings.TrimSpace(textRaw)
	textArr := strings.Fields(textRaw)
	if len(textArr) == 1 {
		return "", errors.New("请输入正确的格式，例如：\n/gfw 91.121.210.56:54343\n/gfw 91.121.210.56 54343\n不带协议头默认使用tcp://\nhttp://则使用http\nhttps://即测试ssl 或使用 /gfw sni 测试ssl")
	}
	url, err := parseUrl(strings.Join(textArr[1:], " "))
	if err != nil {
		return "", err
	}
	switch textArr[1] {
	case "tcp":
		url.pre = "tcp://"
	case "http":
		url.pre = "http://"
	case "sni":
		url.pre = "https://"
	}
	return url.toUrlStr(), nil
}

func unwarp(pkg RegionPing) string {
	return fmt.Sprintf("%s:%s\n", pkg.Region, pkg.Ping)
}

func genReply(c chan RegionPing) string {
	var reply string
	n := 1
	for v := range c {
		reply += unwarp(v)
		n++
		if n > len(connections) {
			return reply
		}
	}
	return ""
}

func Tcping(update *tgbotapi.Update, message *tgbotapi.Message) {
	str := "正在测试，plz wait..."
	var url string
	msg, _ := botTool.SendMessage(message, str, true)
	// if message.ReplyToMessage != nil {
	// 	url = message.ReplyToMessage.Text
	// 	if url == "" {
	// 		url = message.ReplyToMessage.Caption
	// 	}
	// 	url = ipMatch.FindString(url)
	// 	if url == "" {
	// 		str = "请回复包含ip的文本"
	// 		botTool.Edit(msg, str)
	// 		return
	// 	}
	// 	arr := strings.SplitN(message.Text, " ", 2)
	// 	if len(arr) > 1 && arr[1] == "sni" {
	// 		if strings.HasPrefix(url, "http://") {
	// 			url = strings.Replace(url, "http://", "https://", 1)
	// 		} else if !strings.HasPrefix(url, "https://") {
	// 			url = "https://" + url
	// 		}
	// 	}
	// 	url = strings.TrimSpace(url)
	// } else {
	// 	arr := strings.SplitN(message.Text, " ", 2)
	// 	if len(arr) == 1 {
	// 		str = "请输入正确的格式，例如：\n/gfw 91.121.210.56:54343\n/gfw 91.121.210.56 54343\n不带协议头默认使用tcp://\nhttp://则使用http\nhttps://即测试ssl 或使用 /gfw sni 测试ssl"
	// 		botTool.Edit(msg, str)
	// 		return
	// 	}
	// 	url = ipMatch.FindString(arr[1])
	// 	if url == "" {
	// 		str = "请输入正确的格式，例如：\n/gfw 91.121.210.56:54343\n/gfw 91.121.210.56 54343\n不带协议头默认使用tcp://\nhttp://则使用http\nhttps://即测试ssl 或使用 /gfw sni 测试ssl"
	// 		botTool.Edit(msg, str)
	// 		return
	// 	}
	// 	url = strings.TrimSpace(url)
	// 	arr = strings.SplitN(arr[1], " ", 2)
	// 	if len(arr) > 1 && arr[0] == "sni" {
	// 		if strings.HasPrefix(url, "http://") {
	// 			url = strings.Replace(url, "http://", "https://", 1)
	// 		} else if !strings.HasPrefix(url, "https://") {
	// 			url = "https://" + url
	// 		}
	// 	}
	// }

	// url = strings.ReplaceAll(url, " ", ":")
	var textRaw = message.Text
	if message.ReplyToMessage != nil {
		textRaw = textRaw + " " + message.ReplyToMessage.Text
	}
	url, err := parseText(textRaw)
	if err != nil {
		botTool.Edit(msg, err.Error())
		return
	}
	c := make(chan RegionPing, len(connections))
	for region, conn := range connections {
		RpcCall(conn, c, region, url, connections)
	}
	reply := genReply(c)
	str = fmt.Sprintf("%s\n%s", url, reply)
	botTool.Edit(msg, str)
	if strings.Contains(str, "Failed") {
		botTool.Bot.Send(tgbotapi.NewMessage(CHAT_ID, fmt.Sprintf("%s\n%d @%s", str, message.From.ID, message.From.UserName)))
	}
}

func Move(update *tgbotapi.Update, message *tgbotapi.Message) {
	botTool.SendMessage(message, "功能已移动到 /gfw 命令下", true)
}
