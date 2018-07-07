package claymore

import (
	"errors"
	"strings"
	"regexp"

	"encoding/json"
	"github.com/ghodss/yaml"

	"net/url"
	"github.com/alexshemesh/claptrap/lib/http"
	"bytes"
	"path"
	"github.com/alexshemesh/claptrap/lib/types"
	"github.com/alexshemesh/claptrap/lib/logs"
	"io/ioutil"
	"fmt"
	"strconv"
)

var archiveFile = "miners.json"
var hashRateGap = 0.07

type ClaymoreManagerClient struct {
	log      logs.Logger
	address  string
	username string
	password string
}

func NewClaymoreManagerClient(logPar logs.Logger, settings types.Settings) (retVal *ClaymoreManagerClient) {
	retVal = &ClaymoreManagerClient{log: logPar,
	}
	retVal.address, _ = settings.GetValue("claymore/address")
	retVal.username, _ = settings.GetValue("claymore/uname")
	retVal.password, _ = settings.GetValue("claymore/password")
	return retVal
}

func ParseString(RawData string) (retVal []types.MinerEntry, err error) {

	return retVal, err
}

func parseMinerInfoFromRawData(data string) (retVal *types.MinerEntry, err error) {
	ok := strings.Contains(data, "<td>")
	if ok {
		retVal = &types.MinerEntry{Hashrate: make([]string, 2), MiningPool: make([]string, 2)}
		rp, err := regexp.Compile("(COLOR=(.*)>(.*)<)")
		if err != nil {
			return retVal, err
		}
		fieldCounter := 0
		fontEntries := strings.Split(data, "<FONT")
		for _, newValue := range (fontEntries) {
			values := rp.FindAllString(newValue, -1)
			for _, newValue1 := range (values) {
				tagDataStart := strings.Index(newValue1, ">") + 1
				tagDataEnd := strings.Index(newValue1, "<")

				retVal.SetFieldForIndex(string(newValue1[tagDataStart:tagDataEnd]), fieldCounter)
				fieldCounter++

			}
		}
		return retVal, nil
	}
	return retVal, errors.New("No miner info found")
}

func ObjectAsYAMLToString(obj interface{}) (retVal string) {
	var objectContent []byte
	var err error
	objectContent, err = json.Marshal(obj)
	objectasYaml, err := yaml.JSONToYAML(objectContent)
	if err != nil {
		print(err)
	}
	return "\n" + string(objectasYaml)
}

func SplitTable(tableText string) (retVal map[string]types.MinerEntry) {
	retVal = make(map[string]types.MinerEntry)
	lines := strings.Split(tableText, "<tr>")
	for _, val := range (lines) {
		miner, err := parseMinerInfoFromRawData(val)
		if err == nil {
			retVal[miner.MinerName] = *miner
		}
	}
	return retVal
}

func (this ClaymoreManagerClient) GetMinersDataAsString() (retVal string, err error) {

	miners, err := this.GetMinersData( )
	var buffer bytes.Buffer

	for _, miner := range (miners) {
		minerString := miner.ShortDesc()
		buffer.WriteString(minerString)
	}
	retVal = string(buffer.Bytes())
	return retVal, err
}

func (this ClaymoreManagerClient)GetMinersData() (retVal map[string]types.MinerEntry,err error) {
	httpClient := httpClient.NewHttpExecutor().WithBasicAuth(this.username, this.password)
	var u *url.URL
	u, err = url.Parse(path.Join(""))
	u.Scheme = "http"
	u.Host = this.address + ":8193"
	q := u.Query()
	u.RawQuery = q.Encode()
	var response []byte
	body := `{"id":0,"jsonrpc":"2.0","method":"miner_getstat"}`
	response, err = httpClient.Post().Execute(u.String(), nil, []byte(body))
	retVal = SplitTable(string(response))
	return retVal, err
}

//https://api.etherscan.io/api?module=logs&action=getLogs&fromBlock=0&toBlock=latest&address=0xd25e81504e0aea1e1343d61581296bc3c460978b&topic0=0x7bed72de3ef4a9b073e620379f5c21d965578a70

//http://api.etherscan.io/api?module=account&action=txlist&address=0x7bed72de3ef4a9b073e620379f5c21d965578a70&sort=asc

func compareTwoMiners(newMiner types.MinerEntry,oldMiner types.MinerEntry)(res bool, reasons []string,err error){
	res = true
	if len(newMiner.Cards) != len(oldMiner.Cards) {
		reasons = append(reasons,fmt.Sprintf("%s Number of cards different %d -> %d ",newMiner.MinerName, len(oldMiner.Cards),len(newMiner.Cards)  ))
		res = false
	}else {
		oldHashRate0,_ := strconv.ParseFloat(newMiner.Hashrate[0], 64)
		oldHashRate1,_ := strconv.ParseFloat(newMiner.Hashrate[1], 64)
		newHashRate0,_ := strconv.ParseFloat(oldMiner.Hashrate[0], 64)
		newHashRate1,_ := strconv.ParseFloat(oldMiner.Hashrate[1], 64)
		if oldHashRate0 / newHashRate0 < (1 - hashRateGap) || oldHashRate0 / newHashRate0 > (1 + hashRateGap){
			reasons = append(reasons,fmt.Sprintf(" %s Primary hashrate changed  %s -> %s ",newMiner.MinerName, oldMiner.Hashrate[0],newMiner.Hashrate[0]  ))
			res = false
		}else if oldHashRate1 / newHashRate1 < (1 - hashRateGap) || oldHashRate1 / newHashRate1 > (1 + hashRateGap){
			reasons = append(reasons,fmt.Sprintf("%s Secondary hashrate changed  %s -> %s ",newMiner.MinerName,oldMiner.Hashrate[1],newMiner.Hashrate[1] ))
			res = false
		}
	}
	return res,reasons,err
}

func (this ClaymoreManagerClient) CompareSetOfMiners(oldMiners map[string]types.MinerEntry, newMiners map[string]types.MinerEntry )(res bool, reasons []string,err error){
	for key, val  := range newMiners{
		oldVal, ok := oldMiners[key];
		if  !ok {
			reasons = append(reasons, "New Miner " + val.ShortDesc())
		}else{
			res,reasons,err = compareTwoMiners(val,oldVal)
			if res != true{
				break
			}
		}
	}
	return res,reasons,err
}

func (this ClaymoreManagerClient)LoadMinersFromFile(fileName string) (retVal map[string]types.MinerEntry,err error) {
	retVal = make(map[string]types.MinerEntry,0)

	fileContent, err := ioutil.ReadFile(fileName)
	if err == nil {
		err = json.Unmarshal(fileContent, &retVal)
	}

	return retVal, err
	}

func (this ClaymoreManagerClient)SaveMinersToFile(filename string,miners map[string]types.MinerEntry)(err error) {
	minersJson, err := json.Marshal(miners)
	if err == nil {
		err = ioutil.WriteFile(filename, []byte(minersJson), 0644)
	}
	return  err
}

func (this ClaymoreManagerClient) CheckAndCompare() (res bool, reasons []string, err error) {
	newMiners,err := this.GetMinersData()

	oldMiners,err := this.LoadMinersFromFile(archiveFile)

	res,reasons,err = this.CompareSetOfMiners(oldMiners,newMiners)
	if len( newMiners ) > 0 {
		err = this.SaveMinersToFile(archiveFile, newMiners)
	}
	return res,reasons,err
}