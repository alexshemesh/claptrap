package claymore

import (
	"testing"
	"path"
	"github.com/stretchr/testify/assert"

	"io/ioutil"
	"github.com/alexshemesh/claptrap/lib/logs"
	"github.com/alexshemesh/claptrap/lib/vault"
)

func Test_ParseStraight(t *testing.T){
	fileContent,err := ioutil.ReadFile (path.Join("..","..","testdata","claymore_output.html"))

	assert.Nil(t,err)

	SplitTable(string(fileContent))
}

func Test_ParseCards(t *testing.T){
	fileContent,err := ioutil.ReadFile (path.Join("..","..","testdata","claymore_output.html"))
	assert.Nil(t,err)
	results := SplitTable(string(fileContent))
	assert.NotEqual(t, 0, len(results))

	secondCard := results["gemb66"]
	assert.NotEqual(t,0 , len(secondCard.Cards))
}

func INtegrationTest_GetMinersData(t *testing.T){
//	fileContent,err := ioutil.ReadFile (path.Join("..","..","testdata","claymore_output.html"))


	settings := vault.NewVaultTestKit()

	log := logs.NewLogger("GetMinersData test")


	client := NewClaymoreManagerClient(*log, settings)
	res,err := client.GetMinersData()
	assert.NoError(t,err)
	assert.NotNil(t, res)

}

func Test_CompareLessCards( t *testing.T){
	settings := vault.NewVaultTestKit()
	log := logs.NewLogger("GetMinersData test")
	client := NewClaymoreManagerClient(*log, settings)

	oldMiners,err := client.LoadMinersFromFile(path.Join("..","..","testdata","miners.json"))
	assert.NoError(t,err)

	newMiners,err := client.LoadMinersFromFile(path.Join("..","..","testdata","miners_lesscards.json"))
	assert.NoError(t,err)

	res,reasons, err := client.CompareSetOfMiners(oldMiners,newMiners)
	assert.NoError(t,err)
	assert.Equal(t, false,res)
	assert.NotEqual(t, "",reasons)
}

func Test_CompareHashrate( t *testing.T){
	settings := vault.NewVaultTestKit()
	log := logs.NewLogger("GetMinersData test")
	client := NewClaymoreManagerClient(*log, settings)

	oldMiners,err := client.LoadMinersFromFile(path.Join("..","..","testdata","miners.json"))
	assert.NoError(t,err)

	newMiners,err := client.LoadMinersFromFile(path.Join("..","..","testdata","miners_lowhashrate.json"))
	assert.NoError(t,err)

	res,reasons, err := client.CompareSetOfMiners(oldMiners,newMiners)
	assert.NoError(t,err)
	assert.Equal(t, false,res)
	assert.NotEqual(t, "",reasons)
}