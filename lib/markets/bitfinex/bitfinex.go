package bitfinexClient

import (
	"github.com/bitfinexcom/bitfinex-api-go/v1"
	"github.com/alexshemesh/claptrap/lib/types"
	"github.com/alexshemesh/claptrap/lib/logs"
	"strconv"
)

type Bitfinex struct {
	apikey string
	secret string
	log logs.Logger
}


func NewBitfinexClient(logPar logs.Logger, settings types.Settings)( retVal *Bitfinex ){
	retVal = &Bitfinex{log:logPar }

	retVal.apikey,_ = settings.GetValue("bitfinex/apikey")
	retVal.secret,_ = settings.GetValue("bitfinex/apisecret")
	return retVal
}

func (this Bitfinex) GetBalance( ) (retVal types.Balance, err error )  {
	client := bitfinex.NewClient().Auth(this.apikey, this.secret)
	var RawBalances []bitfinex.WalletBalance

	RawBalances,err = client.Balances.All()
	if err == nil {
		for _,item := range(RawBalances) {
			amountAsFloat,_ := strconv.ParseFloat(item.Amount, 64)
			newVal := types.WalletPosition{Name: item.Currency,  Amount: amountAsFloat  }
			retVal.Positions = append(retVal.Positions, newVal )
		}
	}
	return retVal, err
}





