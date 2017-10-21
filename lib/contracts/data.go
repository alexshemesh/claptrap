package contracts
import "gopkg.in/telegram-bot-api.v4"

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