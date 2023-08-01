package funcs

// import (
// 	"bot/botTool"
// 	. "bot/config"
// 	"bot/dbManager"
// 	"fmt"
// 	"regexp"
// 	"strings"

// 	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
// )

// var rulesDb = dbManager.InitMysql("groupRules", RULES_TOKEN, "grouprules")

// type groupRule struct {
// 	rule  *regexp.Regexp
// 	reply string
// }

// type groupRules map[int64][]*groupRule

// func newGroupRule(regexpStr string, replyStr string) (*groupRule, error) {
// 	re, err := regexp.Compile(regexpStr)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &groupRule{
// 		rule:  re,
// 		reply: replyStr,
// 	}, err
// }

// var rules = make(groupRules, 0)

// func (gr *groupRules) getrule(chatid int64) []*groupRule {
// 	r, ok := rules[chatid]
// 	if !ok {
// 		rls := rulesDb.GetAllRules(chatid)
// 		if rls == nil {
// 			rules[chatid] = make([]*groupRule, 0)
// 		} else {
// 			for _, v := range rls {
// 				rule, _ := newGroupRule(v[0], v[1])
// 				rules[chatid] = append(rules[chatid], rule)
// 			}
// 		}
// 		r = rules[chatid]
// 	}
// 	return r
// }

// func Add(update *tgbotapi.Update, message *tgbotapi.Message) {
// 	if message.Chat.ID != -1001738281858 && message.Chat.ID != -1001386067483 && !checkAdmin(message) && message.From.ID != 1456780662 {
// 		botTool.SendMessage(message, "您不是管理,无权添加规则。", true)
// 		return
// 	}
// 	arr := strings.Split(message.Text, "^$==")
// 	if len(arr) != 2 {
// 		botTool.SendMessage(message, "请包含两个要素,并使用`^$==`分隔：\\[正则表达式\\] ^$== \\[回复内容\\]", true, "Markdown")
// 		return
// 	}
// 	regexpStr := strings.Split(arr[0], " ")[1]
// 	rule, err := newGroupRule(regexpStr, arr[1])
// 	if err != nil {
// 		botTool.SendMessage(message, "正则表达式书写错误", true)
// 		return
// 	}
// 	rules[message.Chat.ID] = append(rules[message.Chat.ID], rule)
// 	rulesDb.AddRules(message.Chat.ID, regexpStr, arr[1])
// 	text := fmt.Sprintf("正则: %s\n\n回复: %s", regexpStr, arr[1])
// 	botTool.SendMessage(message, text, true)
// }

// func MatchMessage(update *tgbotapi.Update, message *tgbotapi.Message) {
// 	r := rules.getrule(message.Chat.ID)
// 	for _, rule := range r {
// 		if rule.rule.MatchString(message.Text) {
// 			botTool.SendMessage(message, rule.reply, true)
// 		}
// 	}
// }
