package funcs

import (
	"bot/botTool"
	. "bot/config"
	"bot/dbManager"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var adrule = regexp.MustCompile("((免|中文|询|我).*(电报|tg|telegram))|((电报|tg|telegram).*(免|vpn))")

const (
	locb = 0x80 // 1000 0000
	hicb = 0xBF // 1011 1111
	xx   = 0xF1 // invalid: size 1
	as   = 0xF0 // ASCII: size 1
	s1   = 0x02 // accept 0, size 2
	s2   = 0x13 // accept 1, size 3
	s3   = 0x03 // accept 0, size 3
	s4   = 0x23 // accept 2, size 3
	s5   = 0x34 // accept 3, size 4
	s6   = 0x04 // accept 0, size 4
	s7   = 0x44 // accept 4, size 4
)

var first = [256]uint8{
	//   1   2   3   4   5   6   7   8   9   A   B   C   D   E   F
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, // 0x00-0x0F
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, // 0x10-0x1F
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, // 0x20-0x2F
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, // 0x30-0x3F
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, // 0x40-0x4F
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, // 0x50-0x5F
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, // 0x60-0x6F
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, // 0x70-0x7F
	//   1   2   3   4   5   6   7   8   9   A   B   C   D   E   F
	xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, // 0x80-0x8F
	xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, // 0x90-0x9F
	xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, // 0xA0-0xAF
	xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, // 0xB0-0xBF
	xx, xx, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, // 0xC0-0xCF
	s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, // 0xD0-0xDF
	s2, s3, s3, s3, s3, s3, s3, s3, s3, s3, s3, s3, s3, s4, s3, s3, // 0xE0-0xEF
	s5, s6, s6, s6, s7, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, // 0xF0-0xFF
}

type acceptRange struct {
	lo uint8 // lowest value for second byte.
	hi uint8 // highest value for second byte.
}

var acceptRanges = [16]acceptRange{
	0: {locb, hicb},
	1: {0xA0, hicb},
	2: {locb, 0x9F},
	3: {0x90, hicb},
	4: {locb, 0x8F},
}

var config = dbManager.InitMysql("config", CONFIG_TOKEN, "config")

func Getmessage(update *tgbotapi.Update, message *tgbotapi.Message) {
	if message != nil {
		if message.Text == "" && message.NewChatMembers != nil {
			newMemberMsg := message.NewChatMembers[0]
			if adrule.MatchString(newMemberMsg.FirstName) {
				botTool.BanMember(update, message, newMemberMsg.ID, -1)
				botTool.SendMessage(message, fmt.Sprintf("Allen觉得 [%d](tg://user?id=%d) 是来发广告的,所以禁言了他", newMemberMsg.ID, newMemberMsg.ID), false, "Markdown")
			}
		}
		config.AddGroup(update, message, fmt.Sprint(message.Chat.ID), message.Chat.UserName, message.Chat.Title, fmt.Sprint(message.From.ID), message.From.UserName, botTool.GetName(update, message))
	}
}

type RegionPing struct {
	Ping   string
	Region string
}

func getHistory(body *[]byte, date ...string) {
	var resp *http.Response
	if len(date) == 0 {
		resp, _ = http.Get("http://hao.360.cn/histoday")
	} else {
		resp, _ = http.Get(fmt.Sprintf("http://hao.360.cn/histoday/%s%s.html", date[0], date[1]))
	}
	defer resp.Body.Close()
	*body, _ = io.ReadAll(resp.Body)
}

func httpfix(url string) string {
	url = strings.TrimSpace(url)
	if url[0:4] != "http" || (url[5:8] != "://" && url[4:7] != "://") {
		url = "http://" + url
	}
	return url
}

func wget(url, file string) (err error) {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.OpenFile(file, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	defer f.Close()
	context, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	_, err = f.Write(context)
	return
}

func Dial(network, address string) (*rpc.Client, error) {
	conn, err := net.DialTimeout(network, address, 2*time.Second)
	if err != nil {
		return nil, err
	}
	return rpc.NewClient(conn), nil
}

func RpcCall(conn *rpc.Client, replys chan RegionPing, region string, target string, m map[string]*rpc.Client) {
	var reply string
	var err error
	if conn == nil {
		// replys <- RegionPing{Ping: "Failed(Api Error)", Region: region}
		// go func() {
		// 	defer func() {
		// 		if err := recover(); err != nil {
		// 			return
		// 		}
		// 	}()
		// 	conn, _ := rpc.Dial("tcp", TP_URLS[region])
		// 	m[region] = conn
		// }()
		// return
		conn, err = Dial("tcp", TP_URLS[region])
		if err != nil {
			replys <- RegionPing{Ping: "Failed(Api Error)", Region: region}
			return
		}
		m[region] = conn
	}
	call := conn.Go("GfwTest.Tcping", target, &reply, make(chan *rpc.Call, 1))
	select {
	case call = <-call.Done:
		err := call.Error
		if err != nil {
			conn, err = Dial("tcp", TP_URLS[region])
			if err != nil {
				replys <- RegionPing{Ping: "Failed(Api Error)", Region: region}
				m[region] = nil
				return
			}
			m[region] = conn
			err = conn.Call("GfwTest.Tcping", target, &reply)
		}
		if err != nil {
			reply = fmt.Sprintf("Failed(%s)", err)
		} else {
			reply = fmt.Sprintf("Succeeded(%s)", reply)
		}
		replys <- RegionPing{Ping: reply, Region: region}
	case <-time.After(2500 * time.Millisecond):
		replys <- RegionPing{Ping: "Failed(Api Error)", Region: region}
	}
}

// func goHttp(url string, target string, wg ...*sync.WaitGroup) string {
// 	defer func() {
// 		if err := recover(); err != nil {
// 			log.Println(err)
// 			log.Println(string(debug.Stack()))
// 		}
// 	}()
// 	if wg != nil {
// 		defer wg[0].Done()
// 	}
// 	a, err := rpc.Dial("tcp", url)
// 	if err != nil {
// 		fmt.Println(err)
// 		return "Api Error"
// 	}
// 	var reply string
// 	err = a.Call("GfwTest.Tcping", target, &reply)
// 	if err != nil {
// 		reply = fmt.Sprintf("Failed(%s)", err)
// 	} else {
// 		reply = fmt.Sprintf("Succeeded(%s)", reply)
// 	}
// 	return reply
// }

func runeIndexInString(s string, n int) (int, bool) {
	var i int
	for ; n > 0 && i < len(s); n-- {
		if s[i] < utf8.RuneSelf {
			// ASCII fast path
			i++
			continue
		}

		x := first[s[i]]
		if x == xx {
			i++ // invalid.
			continue
		}

		size := int(x & 7)
		if i+size > len(s) {
			i++ // Short or invalid.
			continue
		}
		accept := acceptRanges[x>>4]
		if c := s[i+1]; c < accept.lo || accept.hi < c {
			size = 1
		} else if size == 2 {
		} else if c := s[i+2]; c < locb || hicb < c {
			size = 1
		} else if size == 3 {
		} else if c := s[i+3]; c < locb || hicb < c {
			size = 1
		}
		i += size
	}

	return i, n <= 0
}

func find(item rune, slice map[rune]struct{}) bool {
	_, ok := slice[item]
	return ok
}

func getToMap(url string) (res map[string]interface{}) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println("GetTomap", err)
		return
	}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&res)
	return
}

func checkAdmin(message *tgbotapi.Message) bool {
	ChatMember, Error := botTool.Bot.GetChatMember(tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: message.Chat.ID,
			UserID: message.From.ID,
		},
	})
	return Error == nil && (ChatMember.IsCreator() || ChatMember.IsAdministrator())
}

func getCoin(update *tgbotapi.Update, message *tgbotapi.Message, coinType string) {
	var text map[string]interface{}
	url := fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%sUSDT", coinType)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("GetTocoin", err)
		return
	}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&text)
	msg := "正在查询，plz wait..."
	msgConfig, _ := botTool.SendMessage(message, msg, true)
	price, _ := strconv.ParseFloat(text["price"].(string), 64)
	msg = fmt.Sprintf("啊哈哈哈哈哈哈\n价格来咯！\n1.0 %s = %.2f USD", strings.ToUpper(coinType), price)
	botTool.Edit(msgConfig, msg)
}

func getReplyAt(update *tgbotapi.Update, message *tgbotapi.Message) string {
	return fmt.Sprintf("[%s](tg://user?id=%d)", botTool.GetReplyName(update, message), message.ReplyToMessage.From.ID)
}

func getAt(update *tgbotapi.Update, message *tgbotapi.Message) string {
	return fmt.Sprintf("[%s](tg://user?id=%d)", botTool.GetName(update, message), message.From.ID)
}

// func urlToQr(url string){

// }
