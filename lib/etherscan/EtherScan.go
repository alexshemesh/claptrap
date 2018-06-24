package etherscan

import (
	"time"
	"github.com/alexshemesh/claptrap/lib/types"
	"github.com/alexshemesh/claptrap/lib/logs"
	"github.com/Jeffail/gabs"
	"strconv"
	"github.com/alexshemesh/claptrap/lib/http"

)

type EtherScanClient struct{
	log *logs.Logger
	blockChainURL string
	address string
}


func NewEtherScanClient(logPar *logs.Logger,  settings types.Settings) (retVal *EtherScanClient){
	etherAddr,_  := settings.GetValue("eth/address")

	retVal = &EtherScanClient{
		log: logPar,
		blockChainURL: "http://api.etherscan.io/api?module=account&action=txlist&address=",
		address: etherAddr,
	}
	return retVal
}

func parseSingleTransaction(entry *gabs.Container)(retVal types.TransactionEntry,err error){
	retVal = types.TransactionEntry{

	}
	contract := entry.S("contractAddress")
	if contract.Data().(string) == "" {
		retVal.Currency = "ETH"
	}
	value := entry.S("value")
	valueStr := value.Data().(string)
	valueInt, err := strconv.Atoi(valueStr)
	retVal.Value = float64(valueInt) / 10e+17
	retVal.FromAddress = entry.S("from").Data().(string)

	value = entry.S("timeStamp")
	valueStr = value.Data().(string)
	valueInt,err = strconv.Atoi(valueStr)
	retVal.Time = time.Unix(int64(valueInt), 0)
	return retVal,err
}

func (this EtherScanClient)ParseData(rawData []byte)( retVal []types.TransactionEntry, err error){
	jsonParsed, err := gabs.ParseJSON(rawData)

	children, err := jsonParsed.S("result").Children()
	for _, child := range children {
		transaction,err := parseSingleTransaction(child)
		if err != nil {
			this.log.Error(err)
		}
		retVal = append(retVal,transaction)
	}
	return retVal,err
}

func (this EtherScanClient)GetTransactions(from time.Time, till time.Time)(retVal []types.TransactionEntry,err error){

	httpClient := httpClient.NewHttpExecutor()
///http://api.etherscan.io/api?module=account&action=txlist&address=0x7bed72de3ef4a9b073e620379f5c21d965578a70&sort=asc


	resp,err := httpClient.Get().Execute(this.blockChainURL+this.address,nil,nil)
	list,err:= this.ParseData(resp)
	for _,val := range(list){
		if val.Time.After(from) && val.Time.Before(till) {
			retVal = append(retVal, val)
		}
	}
	return retVal,err
}