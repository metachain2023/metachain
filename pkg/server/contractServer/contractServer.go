package contractServer

import (
	"log"
	"net/http"
	"os"

	"metechain/pkg/blockchain"
	"metechain/pkg/config"
	"metechain/pkg/logger"
	"metechain/pkg/server/contractServer/api"
	"metechain/pkg/txpool"
)

func RunmetemaskServer(bc blockchain.Blockchains, tp *txpool.Pool, cfg *config.CfgInfo) {
	s := api.NewmetemaskServer(bc, tp, cfg)
	http.HandleFunc("/", s.HandRequest)

	logger.InfoLogger.Println("Running contractServer...", cfg.metemaskCfg.ListenPort)
	err := http.ListenAndServe(cfg.metemaskCfg.ListenPort, nil)
	if err != nil {
		log.Println("start fasthttp fail:", err.Error())
		os.Exit(1)
	}
}
