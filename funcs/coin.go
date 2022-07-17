package funcs

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Btc(update *tgbotapi.Update) {
	getCoin(update, "btc")
}

func Xmr(update *tgbotapi.Update) {
	getCoin(update, "xmr")
}

func Eth(update *tgbotapi.Update) {
	getCoin(update, "eth")
}
