package etherscan

import (

	"io/ioutil"
	"path"
	"testing"
	"github.com/alexshemesh/claptrap/lib/logs"
	"github.com/stretchr/testify/assert"
	"github.com/Jeffail/gabs"
	"github.com/alexshemesh/claptrap/lib/vault"
)

func Test_GetTransactions(t *testing.T){
	log := logs.NewLogger("Test_GetTransactions")
	settings := vault.NewVaultTestKit()
	client := NewEtherScanClient(log,settings)

	fileContent,err := ioutil.ReadFile(path.Join("..","..","testdata","etherscandata.json"))
	assert.Nil(t, err)
	client.ParseData(fileContent)

}

func Test_parseSingleTaransacton(t *testing.T){
	fileContent,err := ioutil.ReadFile(path.Join("..","..","testdata","eths_singleentry.json"))
	assert.Nil(t, err)

	gabsData,err := gabs.ParseJSON(fileContent)
	transaction,err := parseSingleTransaction(gabsData)
	assert.Nil(t, err)
	assert.Equal(t,transaction.Currency,"ETH")
	assert.Equal(t,0.62986797, transaction.Value )
	assert.Equal(t,"0x1151314c646ce4e0efd76d1af4760ae66a9fe30f", transaction.FromAddress )
}