package main

import (
	"github.com/CryptoQuantX/chain_monitor/models"
)

func main() {
	//go monitor.MonitorWhale()
	//time.Sleep(1 * time.Second)
	for {
		ethblock := &models.BlockStruct{}

		eblock, cli := ethblock.ParseEthBlock()

	}
	//watchToken.WatchUSDC()

}
