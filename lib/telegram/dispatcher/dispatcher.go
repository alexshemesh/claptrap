package dispatcher

import (
	"github.com/alexshemesh/claptrap/lib/logs"
	"github.com/alexshemesh/claptrap/lib/telegram"
	"gopkg.in/telegram-bot-api.v4"

	"time"
	"fmt"
	"github.com/alexshemesh/claptrap/lib/markets/kuna"
	"github.com/alexshemesh/claptrap/lib/markets/bitfinex"
	"github.com/alexshemesh/claptrap/lib/types"

	"strings"
	"github.com/alexshemesh/claptrap/lib/claymore"

)

type messageHandler func( types.TGCommand)(string,error)

var ProgramVersion string

type Dispatcher struct {
	log logs.Logger

	auth types.Auth
}

func NewDispatcher(logPar logs.Logger, authPar types.Auth)(retVal *Dispatcher) {
	retVal = &Dispatcher{log: logPar,auth:authPar }
	return retVal
}

func (this Dispatcher) CheckAuth(userName string)(retVal types.Settings, err error){
	retVal,_,err = this.auth.LogInCached(userName)
	return retVal,err
}



func (this Dispatcher)Dispatch(bot telegram.TelegramBot, cmd types.TGCommand)(err error){

	if cmd.Msg.Command() == "kuna_orders_book" {
		err = this.cmdHandler(this.handleKunaOrdersBookCmd,bot, cmd)
	}else if cmd.Msg.Command() == "bitfinex_account" {
		err = this.cmdHandler(this.handleBitfinexAccountCmd, bot, cmd)
	} else if cmd.Msg.Command() == "login" {
		err = this.cmdHandler(this.handleLoginCmd, bot, cmd)
		} else if cmd.Msg.Command() == "miners" {
			err = this.cmdHandler(this.handleMinerCmd, bot, cmd)
	} else {
		err = this.cmdHandler(this.handleUsageCmd, bot, cmd)
	}
	return err
}

func (this Dispatcher)cmdHandler(functionToRun messageHandler,bot telegram.TelegramBot, cmd types.TGCommand )(err error) {
	startTime := cmd.Msg.Time()
	var responseText string
	if cmd.Msg.Command() != "login" {

		cmd.Settings, err = this.CheckAuth( cmd.Msg.Chat.UserName )
		if err != nil {
			if err.Error() == "Authentication Failed" {
				responseText = "You are not authenticated to use system.\n Use /login <your password> command to log in into system"
			} else {
				responseText = err.Error()
			}
		}
	}

	if	err == nil {
		responseText, err = functionToRun( cmd )
	}
	duration := time.Since(startTime)
	responseText = responseText + fmt.Sprintf( "\nTook %s to execute" , duration.String())
	err = sendReply(bot,responseText,cmd.Msg.Chat.ID,cmd.Msg.MessageID )
	return err
}

func (this Dispatcher)handleUsageCmd( cmd types.TGCommand )(response string, err error){

	responseText := fmt.Sprintf( "Hi there. This is Clpatrap (v %s) talking. This is the list of available commands\n", ProgramVersion)
	responseText = responseText + "/miners\n"
	//responseText = responseText + "/kuna_orders_book\n"
	//responseText = responseText + "/bitfinex_account\n"


	return responseText, err
}

func (this Dispatcher)handleKunaOrdersBookCmd(cmd types.TGCommand )(responseText string,err error){


	kunaClient := kuna.NewKunaClient(this.log)

	book,err := kunaClient.GetOrdersBook()

	if err == nil {
		responseText = "Asks\n"
		for i, order := range (book.Asks){
			responseText = responseText + fmt.Sprintf("\t%d %s %s %s %s\n", i, order.Side, order.OrdType, order.Price,  order.Volume )
		}
		responseText = responseText + "Bids\n"
		for i, order := range (book.Bids) {
			responseText = responseText + fmt.Sprintf("\t%d %s %s %s %s\n", i, order.Side, order.OrdType, order.Price, order.Volume)
		}

	}else{
		responseText = responseText + "Error: " + err.Error()
	}

	return responseText,err
}

func (this Dispatcher)handleBitfinexAccountCmd(cmd types.TGCommand )(response string,err error){

	var responseText string

	bitfinexClnt := bitfinexClient.NewBitfinexClient(this.log, cmd.Settings)
	var balance types.Balance
	balance,err = bitfinexClnt.GetBalance()
	if err == nil {
		responseText = responseText + "Bitfinex Balances:\r\n"
		for _, pos := range(balance.Positions){
			responseText = responseText + fmt.Sprintf("\t (%s)%s Amount %f Value %f\n", pos.Type, pos.Name, pos.Amount, pos.Value)
		}
	}else{
		responseText = responseText + "Error: " + err.Error()
	}

	return responseText,err
}

func (this Dispatcher)handleLoginCmd(cmd types.TGCommand )(response string,err error){

	var params []string
	params = strings.Split(cmd.Msg.Text," ")
	var localAuth types.Auth
	_, localAuth,err = this.auth.LogIn(cmd.Msg.Chat.UserName, params[1])

	responseText := "OK"
	if err == nil {
		if localAuth.GetUserName() == "" {
			responseText = "Unknown error"
		}else{
			this.auth.SaveCachedToken( localAuth.GetUserName(), localAuth.GetUserToken())
		}
	}else{
		this.log.Error(err)
		responseText = err.Error()
	}
	return responseText, err
}

func (this Dispatcher)handleMinerCmd(cmd types.TGCommand )(response string,err error){
	claymoreClient := claymore.NewClaymoreManagerClient(this.log,cmd.Settings)
	minersData,err := claymoreClient.GetMinersDataAsString()
	if minersData == "" {
		minersData,err = claymoreClient.GetMinersDataAsString()
	}
	return minersData, err
}


func sendReply(bot telegram.TelegramBot, text string,chatID int64,  replyTo int)(err error){
	if len(text) < 4096 {
		newMsg := tgbotapi.NewMessage(chatID, text)
		newMsg.ReplyToMessageID = replyTo
		err = bot.Send(newMsg)
	}else{
		fileBytes := tgbotapi.FileBytes{Name: "Result",Bytes: []byte(text)}
		document := tgbotapi.NewDocumentUpload(chatID, fileBytes)
		err = bot.SendFile(document)
	}
	return err
}


