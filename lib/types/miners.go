package types

import (
	"strings"
	"regexp"

	"strconv"
	"fmt"
)

type CardEntry struct{
	temp int
	load int
}

type MinerEntry struct {
	MinerName   string
	Currency    string
	RunningTime string
	Hashrate    []string
	MiningPool  []string
	Cards				[]CardEntry
}

func parseCardsData(data string)(retVal []CardEntry){
	retVal = make([]CardEntry,0)
	r := regexp.MustCompile(`\(([0-9]*)C:([0-9]*)% \)`)

	matches := r.FindAllString(data, -1)

	for _, nextMatch := range(matches){
		results := r.FindStringSubmatch(nextMatch)
		temp,_ := strconv.Atoi(results[1])
		load,_ := strconv.Atoi(results[2])
		newCard := CardEntry{temp,load}
		retVal  =append(retVal,newCard)
	}
	return retVal
}

func (this *MinerEntry) SetFieldForIndex(data string, index int) {
	if index == 0 {
		this.MinerName = data
	} else if index == 2 {
		this.RunningTime = data
	} else if index == 3 {
		this.Hashrate[0] = data
	} else if index == 4 {
		this.Hashrate[1] = data
	} else if index == 6 {
		minerpools := strings.Split(data, ";")
		this.MiningPool[0] = minerpools[0]
		if len(minerpools) >1 {
			this.MiningPool[1] = minerpools[1]
		}
	}else if index == 5{
		this.Cards = parseCardsData(data)
	}

}

func (miner MinerEntry) ShortDesc() (retVal string) {
	retVal = fmt.Sprintf("%s %s Num Of Cards: %d\r\t%s, %s\r\t%s, %s\r", miner.MinerName,miner.RunningTime, len(miner.Cards), miner.MiningPool[0],miner.Hashrate[0],miner.MiningPool[1], miner.Hashrate[1])
	return retVal
}