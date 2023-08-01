package funcs

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Btc(update *tgbotapi.Update, message *tgbotapi.Message) {
	getCoin(update, message, "BTC")
}

func Xmr(update *tgbotapi.Update, message *tgbotapi.Message) {
	getCoin(update, message, "XMR")
}

func Eth(update *tgbotapi.Update, message *tgbotapi.Message) {
	getCoin(update, message, "ETH")
}
