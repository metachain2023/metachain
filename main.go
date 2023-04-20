package main

import (
	"flag"
	"fmt"
	"net"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"os"
	"runtime"

	"metechain/pkg/blockchain"
	"metechain/pkg/config"
	"metechain/pkg/consensus"
	"metechain/pkg/controller"
	"metechain/pkg/logger"
	"metechain/pkg/p2p"
	"metechain/pkg/storage/store/pb"
	"metechain/pkg/txpool"
	"metechain/pkg/util/ntp"

	"github.com/cockroachdb/pebble"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gofrs/uuid"

	"net/http"
	_ "net/http/pprof"

	"metechain/pkg/miner"
	"metechain/pkg/server/grpcserver"
	"metechain/pkg/server/rpcserver"

	"go.uber.org/zap"
)

var RunMode *int
var NoMining *bool
var MinerAddr *string
var GreamHost *string

const Version = "version: matechain v0.0.0"

func init() {
	if err := ntp.UpdateTimeFromNtp(); err != nil {
		panic(fmt.Errorf("failed to set time:%s", err))
	}

	if err := logger.InitLogger(logger.DefaultConfig()); err != nil {
		panic(err)
	}

	// if err := logger.RewriteStderrFile("runtime_err"); err != nil {
	// 	panic(err)
	// }

	avlNum := flag.Int("n", 0, "要使用的cpu核心数量")
	cycle := flag.Int("c", 10, "每次cycle个块调整下难度")
	RunMode = flag.Int("m", 0, "运行模式")
	NoMining = flag.Bool("nomining", false, "不挖矿节点启用此命令")
	MinerAddr = flag.String("mineraddr", "", "矿工地址")
	GreamHost = flag.String("greamhost", "", "种子节点")

	flag.Parse()

	if err := miner.SetConf(*avlNum, *cycle); err != nil {
		fmt.Println("-b 0x12345678", err)
		os.Exit(1)
	}

	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(1)
	runtime.GOMAXPROCS(4)
	go http.ListenAndServe(":8090", nil)
}

func main() {
	var nodeName string
	id, _ := uuid.NewV4()
	nodeName = id.String()

	plog := logger.Logger.Named("pebble").Sugar()
	db, err := pebble.Open("pebble.db", &pebble.Options{Logger: plog})
	if err != nil {
		panic(err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	if cfg.ChainCfg == nil || cfg.SververCfg == nil {
		panic("load nil config")
	}
	//
	cfg.MinerConfig.GenesisHash = Version

	// address.SetNetWork(cfg.NetWorkType)

	// mineraddr, _ := address.StringToAddress(cfg.MinerConfig.MiningAddr)
	mineraddr := common.HexToAddress(cfg.MinerConfig.MiningAddr)
	cfg.ChainCfg.Miner = &mineraddr
	b, err := blockchain.New(pb.New(db), cfg.ChainCfg)
	if err != nil {
		panic(err)
	}

	pool, err := txpool.NewPool(txpool.Config{b, logger.Logger})

	//contract server
	/* 	logger.Info("metemashk", zap.Int64("chain ID", cfg.ChainCfg.ChainId))
	   	go contractServer.RunmetemaskServer(b, pool, cfg) */
	//chain server
	// csv := chainserver.NewServer(b, cfg.SververCfg.ChainServerPort, cfg.SververCfg.GRpcAddress)
	// go csv.RunServer()

	var node *p2p.Node
	if *RunMode <= 0 {
		p2pConf := p2p.DefaultConfig()
		p2pConf.NodeName = nodeName
		p2pConf.Logger = zap.NewStdLog(logger.Logger)

		p2pConf.MemberlistConfig.BindPort = cfg.P2PConfig.Port
		p2pConf.MemberlistConfig.AdvertisePort = cfg.P2PConfig.Port
		p2pConf.MemberlistConfig.AdvertiseAddr = cfg.P2PConfig.AdvertiseAddr
		p2pConf.MemberlistConfig.TCPTimeout = 20 * time.Second

		node, err = p2p.Create(p2pConf)
		if err != nil {
			panic(err)
		}

		node.Join(cfg.P2PConfig.JionMembers)

	}

	cbc := consensus.New(b)

	cfg.MinerConfig.NoMining = *NoMining
	if len(*MinerAddr) > 0 {
		cfg.MinerConfig.MiningAddr = *MinerAddr
	}
	m, err := miner.New(cfg.MinerConfig, b, pool, node, cbc, *RunMode, cfg.P2PConfig.AdvertiseAddr)
	if err != nil {
		logger.Error("miner.New ", zap.Error(err))
		panic(err)
	}

	coll, err := controller.New(pool, b, m, logger.Logger, cfg.P2PConfig.AdvertiseAddr, m.InitBits)
	if err != nil {
		panic(err)
	}

	//初始化块数据
	var greamhost string
	if len(cfg.SververCfg.GreamHost) == 0 && len(*GreamHost) == 0 {
		logger.ErrorLogger.Println("GreamHost error")
		return
	}

	if len(*GreamHost) > 0 {
		greamhost = *GreamHost
	} else {
		greamhost = cfg.SververCfg.GreamHost
	}

	gm := net.ParseIP(greamhost)
	if gm == nil || gm.String() == "127.0.0.1" {
		logger.ErrorLogger.Println("greamhost error:", greamhost)
		return
	}

	//go test(&blockchain.BlockHeight)

	logger.InfoLogger.Println("greamhost:", greamhost)
	if greamhost != "0.0.0.0" {
		err = miner.InitBlockChain(greamhost, b, pool, m)
		if err != nil {
			logger.Error("InitBlockChain ", zap.Error(err))
			panic(err)
		}
	}

	if *RunMode <= 0 {
		node.RegisterHandleFunc(coll.RegisterHandleFunc())
	}

	if *RunMode > 0 {
		miner.Start(greamhost, b, pool, m)
	}

	m.Start()

	inside := rpcserver.NewInsideGreeter(b, pool, node, cfg, m)
	go inside.RunInsideGrpc()

	greet := grpcserver.Greeter{Bc: b, Tp: pool, Cfg: cfg, Node: node, Miner: m, NodeName: nodeName}
	go greet.RunGrpc()
	go func() {
		for {
			time.Sleep(20 * time.Second)
			runtime.GC()
			debug.FreeOSMemory()
		}
	}()

	{
		ctrlC := make(chan os.Signal)
		signal.Notify(ctrlC, os.Interrupt, syscall.SIGTERM)
		<-ctrlC
		b.Close()
	}
}
