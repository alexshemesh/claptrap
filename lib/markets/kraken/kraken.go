package kraken

import (
	"github.com/alexshemesh/claptrap/lib/logs"
	"github.com/alexshemesh/claptrap/lib/types"

	krakenApi "github.com/beldur/kraken-go-api-client"
)


// List of valid public methods



type KrakenClient struct{

	log logs.Logger
	apikey string
	secret string

}

func NewKrakenClient(logPar logs.Logger, settingPar types.Settings )(retVal *KrakenClient){
	retVal = &KrakenClient{  log: *logPar.SubLogger("kraken") }
	retVal.apikey,_  = settingPar.GetValue("kraken/key")
	retVal.secret,_  = settingPar.GetValue("kraken/secret")
	return retVal
}

func (this KrakenClient) GetBalance( ) (retVal types.Balance, err error )  {
	client := krakenApi.New(this.apikey, this.secret)
	var balanceRaw *krakenApi.BalanceResponse
	balanceRaw,err = client.Balance()

	retVal.Positions = append(retVal.Positions, types.WalletPosition{Name: "BCH", Amount: float64(balanceRaw.BCH)})
	retVal.Positions = append(retVal.Positions, types.WalletPosition{Name: "BCH", Amount: float64(balanceRaw.DASH)})
	retVal.Positions = append(retVal.Positions, types.WalletPosition{Name: "BCH", Amount: float64(balanceRaw.EOS)})
	retVal.Positions = append(retVal.Positions, types.WalletPosition{Name: "BCH", Amount: float64(balanceRaw.GNO)})
	retVal.Positions = append(retVal.Positions, types.WalletPosition{Name: "BCH", Amount: float64(balanceRaw.KFEE)})
	retVal.Positions = append(retVal.Positions, types.WalletPosition{Name: "BCH", Amount: float64(balanceRaw.USDT)})

	return retVal, err
}

