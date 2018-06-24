package kuna

import (
	"github.com/alexshemesh/claptrap/lib/logs"
	"github.com/alexshemesh/claptrap/lib/http"
	"net/url"
	"path"
	"github.com/alexshemesh/claptrap/lib/types"
	"encoding/json"
)

type KunaClient struct{
	log logs.Logger
	ServerURL string
}

func NewKunaClient(logPar logs.Logger)(retVal *KunaClient){
	retVal = &KunaClient{ log:logPar ,ServerURL:"https://kuna.io"}
	return retVal
}

func (this KunaClient) GetOrdersBook()( retVal types.KunaOrdersBook, err error ){
	this.log.Log("Executing GetOrdersBook for Kuna market")
	httpClient := httpClient.NewHttpExecutor()
	var u *url.URL
	u, err = url.Parse( path.Join("api/v2/order_book?market=btcuah") )
	u.Scheme = "https"
	u.Host = "kuna.io"
	q := u.Query()
	u.RawQuery = q.Encode()
	var response []byte
	response,err  = httpClient.Get().Execute(u.String(),nil,nil)
	if err == nil {
		err = json.Unmarshal(response, &retVal)
	}
	return retVal, err
}