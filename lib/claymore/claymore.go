package claymore

import (
	"errors"
	"strings"
	"regexp"

	"encoding/json"
	"github.com/ghodss/yaml"
	"fmt"
	"net/url"
	"github.com/alexshemesh/claptrap/lib/http"
	"bytes"
	"path"
	"github.com/alexshemesh/claptrap/lib/types"
	"github.com/alexshemesh/claptrap/lib/logs"
)

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

func SplitTable(tableText string) (retVal []types.MinerEntry) {
	lines := strings.Split(tableText, "<tr>")
	for _, val := range (lines) {
		miner, err := parseMinerInfoFromRawData(val)
		if err == nil {
			fmt.Printf("Miner: %s, RunningTime: %s, HashRate1: %s for pool1:%s, HashRate2: %s, for pool2: %s", miner.MinerName, miner.RunningTime, miner.Hashrate[0], miner.MiningPool[0], miner.Hashrate[1], miner.MiningPool[1])
			println("========================================")
			retVal = append(retVal, *miner)
		}
	}
	return retVal
}

func (this ClaymoreManagerClient) GetMinersData() (retVal string, err error) {
	//admin
	//statuscheck
	httpClient := httpClient.NewHttpExecutor().WithBasicAuth(this.username, this.password)
	var u *url.URL
	u, err = url.Parse(path.Join(""))
	u.Scheme = "http"
	u.Host = this.address +":8193"

	q := u.Query()
	u.RawQuery = q.Encode()
	var response []byte
	body := `{"id":0,"jsonrpc":"2.0","method":"miner_getstat"}`
	response, err = httpClient.Post().Execute(u.String(), nil, []byte(body))

	miners := SplitTable(string(response))
	var buffer bytes.Buffer
	for _, miner := range (miners) {
		minerString := miner.ShortDesc()
		buffer.WriteString(minerString)
	}
	retVal = string(buffer.Bytes())
	return retVal, err
}

//https://api.etherscan.io/api?module=logs&action=getLogs&fromBlock=0&toBlock=latest&address=0xd25e81504e0aea1e1343d61581296bc3c460978b&topic0=0x7bed72de3ef4a9b073e620379f5c21d965578a70

//http://api.etherscan.io/api?module=account&action=txlist&address=0x7bed72de3ef4a9b073e620379f5c21d965578a70&sort=asc
