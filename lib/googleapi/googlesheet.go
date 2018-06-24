package googleapi

import (
	"github.com/alexshemesh/claptrap/lib/logs"
	"github.com/alexshemesh/claptrap/lib/types"
)

type GoogleApi struct{
	log *logs.Logger

}

func NewGoogleApi(log *logs.Logger )(retVal *GoogleApi ){
	retVal = &GoogleApi{}
	return retVal
}

func (this GoogleApi)GetWallets () (retVal []types.CoinConfig,err error){
	return retVal,err
}