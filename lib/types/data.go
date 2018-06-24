package types

import (
	"gopkg.in/telegram-bot-api.v4"

	"time"
)

type WalletPosition struct {
	Type string
	Name string
	Amount float64
	Value float64
}

type Balance struct {
	Positions []WalletPosition
	Gain float64
	Deposits float64
}

type TGCommand struct {
	Msg *tgbotapi.Message
	Settings Settings
}

type TransactionEntry struct{
	Time time.Time
	Currency string
	FromAddress string
	ToAddress string
	Value float64
}