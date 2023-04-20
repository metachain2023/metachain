package grpcserver

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"

	"metechain/pkg/blockchain"
	"metechain/pkg/config"
	_ "metechain/pkg/crypto/sigs/ed25519"
	_ "metechain/pkg/crypto/sigs/secp"
	"metechain/pkg/logger"
	"metechain/pkg/miner"
	"metechain/pkg/p2p"
	"metechain/pkg/server"
	"metechain/pkg/server/grpcserver/message"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/handlers"
	"github.com/pkg/errors"
	"golang.org/x/crypto/sha3"

	http "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"

	"metechain/pkg/transaction"
	"metechain/pkg/txpool"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"google.golang.org/protobuf/types/known/timestamppb"
)

const DECIWEI = 1000000000

var Deciwei = new(big.Int).Mul(big.NewInt(DECIWEI), big.NewInt(10000000))

type Greeter struct {
	Bc       blockchain.Blockchains
	Tp       *txpool.Pool
	Cfg      *config.CfgInfo
	Node     *p2p.Node
	Miner    *miner.Miner
	NodeName string
}

var _ message.GreeterServer = (*Greeter)(nil)
var _ message.GreeterHTTPServer = (*Greeter)(nil)

func NewGreeter(bc blockchain.Blockchains, tp *txpool.Pool, cfg *config.CfgInfo) *Greeter {
	return &Greeter{Bc: bc, Tp: tp, Cfg: cfg}
}

func (g *Greeter) RunGrpc() {
	lis, err := net.Listen("tcp", g.Cfg.SververCfg.GRpcAddress)
	if err != nil {
		logger.Error("net.Listen", zap.Error(err))
		os.Exit(-1)
	}

	rpvServ := grpc.NewServer(grpc.UnaryInterceptor(server.IpInterceptor))
	message.RegisterGreeterServer(rpvServ, g)

	// http server
	{
		opts := []http.ServerOption{
			http.Address(g.Cfg.SververCfg.WebAddress),
			http.Filter(handlers.CORS(
				handlers.AllowedOrigins([]string{"*"}),
				handlers.AllowedMethods([]string{"GET", "POST"}),
				handlers.AllowedHeaders([]string{"Content-Type", "X-Requested-With"}),
			)),
		}
		httpServ := http.NewServer(opts...)
		openAPIhandler := openapiv2.NewHandler()
		httpServ.HandlePrefix("/q/", openAPIhandler)
		message.RegisterGreeterHTTPServer(httpServ, g)

		webAddr, err := net.Listen("tcp", g.Cfg.SververCfg.WebAddress)
		if err != nil {
			panic(err)
		}
		go httpServ.Serve(webAddr)

	}

	if err := rpvServ.Serve(lis); err != nil {
		panic(err)
	}
}

func (g *Greeter) GetBalance(ctx context.Context, in *message.ReqBalance) (*message.ResBalance, error) {

	addr := common.HexToAddress(in.Address)
	balance, err := g.Bc.GetBalance(&addr)
	if err != nil {
		logger.Error("g.Bc.GetBalance", zap.Error(err), zap.String("address", in.Address))
		return nil, err
	}
	return &message.ResBalance{Balance: balance.String()}, nil
}

func (g *Greeter) SendTransaction(ctx context.Context, in *message.SendTransactionRequest) (*message.SendTransactionResponse, error) {

	from := common.HexToAddress(in.From)
	to := common.HexToAddress(in.To)

	balance, err := g.Bc.GetAvailableBalance(&from)
	if err != nil {
		return nil, err
	}

	amount, ok := new(big.Int).SetString(in.Amount, 10)
	if !ok {
		return nil, fmt.Errorf("invalid parameter:%s", fmt.Sprintf("amount:%s", in.Amount))
	}

	gasPrice, ok := new(big.Int).SetString(in.GasPrice, 10)
	if !ok {
		return nil, fmt.Errorf("invalid parameter:%s", fmt.Sprintf("gas price:%s", in.GasPrice))
	}

	gaslimit, ok := new(big.Int).SetString(in.GasLimit, 10)
	if !ok {
		return nil, fmt.Errorf("invalid parameter:%s", fmt.Sprintf("gas limit:%s", in.GasLimit))
	}

	gasFeeCap, ok := new(big.Int).SetString(in.GasFeeCap, 10)
	if !ok {
		return nil, fmt.Errorf("invalid parameter:%s", fmt.Sprintf(" gas fee cap:%s", in.GasFeeCap))
	}

	gas := new(big.Int).Mul(gasPrice, gaslimit)
	if gas.Cmp(big.NewInt(0)) == 0 {
		return nil, fmt.Errorf("error: one of gasprice[%v] and gaslimit[%v] is 0", in.GasLimit, in.GasPrice)
	}

	if balance.Cmp(new(big.Int).Add(amount, gas)) == -1 {
		return nil, fmt.Errorf("from(%v) balance(%v) is not enough or out of gas.", from, balance)
	}

	signature, err := hex.DecodeString(in.Signature)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode siganature")
	}

	tx := &transaction.SignedTransaction{
		Transaction: transaction.Transaction{
			From:      &from,
			To:        &to,
			Amount:    amount,
			Nonce:     in.Nonce,
			GasLimit:  gaslimit,
			GasPrice:  gasPrice,
			GasFeeCap: gasFeeCap,
			Type:      transaction.TransferTransaction,
		},
		Signature: signature,
	}

	err = g.Tp.Add(tx)
	if err != nil {
		return nil, err
	}

	data, err := tx.Serialize()
	if err != nil {
		return nil, err
	}

	if g.Node != nil {
		g.Node.SendMessage(p2p.PayloadMessageType, append([]byte{0}, data...))
	}

	hash := hex.EncodeToString(tx.Hash())

	return &message.SendTransactionResponse{Hash: hash}, nil
}

func (g *Greeter) GetBlockByNum(ctx context.Context, in *message.ReqBlockByNumber) (*message.RespBlock, error) {
	b, err := g.Bc.GetBlockByHeight(in.Height)
	if err != nil {
		logger.Error("GetBlockByHeight", zap.Error(fmt.Errorf("error:%v,height %v", err.Error(), in.Height)))
		return &message.RespBlock{Data: []byte{}, Code: -1, Message: err.Error()}, nil
	}
	blockbyte, err := b.Serialize()
	if err != nil {
		logger.Error("Serialize", zap.String("error", err.Error()))
		return &message.RespBlock{Data: []byte{}, Code: -1, Message: err.Error()}, nil
	}
	return &message.RespBlock{Data: blockbyte, Code: 0}, nil

}

func (g *Greeter) GetBlockByHash(ctx context.Context, in *message.ReqBlockByHash) (*message.RespBlockData, error) {
	hash, err := transaction.StringToHash(blockchain.Check0x(in.Hash))
	if err != nil {
		logger.Error("StringToHash", zap.String("error", err.Error()), zap.String("hash", in.Hash))
		return &message.RespBlockData{Data: []byte{}, Code: -1, Message: err.Error()}, nil
	}
	block, err := g.Bc.GetBlockByHash(hash)
	if err != nil {
		//logger.Error("GetBlockByHash", zap.Error(fmt.Errorf("error:%v,hash: %v", err.Error(), in.Hash)))
		return &message.RespBlockData{Data: []byte{}, Code: -1, Message: err.Error()}, nil
	}
	blockbyte, err := block.Serialize()
	if err != nil {
		logger.Error("Serialize", zap.String("error", err.Error()))
		return &message.RespBlockData{Data: []byte{}, Code: -1, Message: err.Error()}, nil
	}
	return &message.RespBlockData{Data: blockbyte, Code: 0}, nil
}

func (g *Greeter) GetTxByHash(ctx context.Context, in *message.ReqTxByHash) (*message.RespTxByHash, error) {
	hash, _ := transaction.StringToHash(blockchain.Check0x(in.Hash))
	tx, err := g.Bc.GetTransactionByHash(hash)
	if err != nil {
		//logger.Error("GetTransactionByHash", zap.Error(fmt.Errorf("error:%v,hash:%v", err.Error(), in.Hash)))
		st, err := g.Tp.GetTxByHash(in.Hash)
		if err != nil {
			return &message.RespTxByHash{Data: []byte{}, Code: -1, Message: err.Error()}, nil
		}
		tx = transaction.NewFinishedTransaction(st, nil, 0)
	}

	txbytes, err := tx.Serialize()
	if err != nil {
		logger.Error("tx Serialize", zap.String("error", err.Error()))
		return &message.RespTxByHash{Data: []byte{}, Code: -1, Message: err.Error()}, nil
	}

	return &message.RespTxByHash{Data: txbytes, Code: 0}, nil
}

func (g *Greeter) GetAddressNonceAt(ctx context.Context, in *message.ReqNonce) (*message.ResposeNonce, error) {

	addr := common.HexToAddress(in.Address)
	resp, err := g.Bc.GetNonce(&addr)
	if err != nil {
		return nil, err
	}

	return &message.ResposeNonce{Nonce: resp}, nil
}

//get max block height
func (g *Greeter) GetMaxBlockHeight(ctx context.Context, in *message.ReqMaxBlockHeight) (*message.ResMaxBlockHeight, error) {
	maxH, err := g.Bc.GetMaxBlockHeight()
	if err != nil {
		logger.Error("rpc GetMaxBlockHeight", zap.String("error", err.Error()))
		return nil, err
	}
	return &message.ResMaxBlockHeight{MaxHeight: maxH}, nil
}

func (g *Greeter) GetBlockDetails(ctx context.Context, req *message.GetBlockDetailsRequest) (replay *message.GetBlockDetailsResponse, err error) {
	block, err := g.Bc.GetBlockByHeight(req.Height)
	if err != nil {
		return nil, err
	}

	txData, _ := block.Transactions[0].Serialize()
	blockData, _ := block.Serialize()
	blockHash := sha3.Sum256(blockData)

	proposer := new(big.Int).SetBytes(blockHash[len(blockHash)-2:])
	replay = &message.GetBlockDetailsResponse{
		Epoch:        block.Height / 32,
		Slot:         block.Height,
		Time:         timestamppb.New(time.Unix(int64(block.Timestamp), 0)),
		Status:       "Proposed",
		Proposer:     new(big.Int).Abs(proposer).Int64(),
		BlockRoot:    hex.EncodeToString(block.Hash),
		ParentRoot:   hex.EncodeToString(block.PrevHash),
		StateRoot:    hex.EncodeToString(block.Root),
		Signature:    hex.EncodeToString(txData),
		RandaoReceal: hex.EncodeToString(blockHash[:]),
		Graffiti:     "",
	}

	for _, tx := range block.Transactions {
		rtx := message.FinalTransaction{
			Stx: &message.SignedTransaction{
				Utx: &message.UnsignedTransaction{
					From:      tx.From.Hex(),
					To:        tx.To.Hex(),
					Amount:    tx.Amount.String(),
					Nonce:     tx.Nonce,
					GasLimit:  tx.GasLimit.String(),
					GasPrice:  tx.GasPrice.String(),
					GasFeeCap: tx.GasFeeCap.String(),
				},
				Signature: hex.EncodeToString(tx.Signature),
			},
			GasUsed:  tx.GasUsed.String(),
			BlockNum: tx.BlockNum,
		}

		replay.Ftxs = append(replay.Ftxs, &rtx)
	}

	return replay, nil
}

func (g *Greeter) GetTransactionDetails(ctx context.Context, in *message.GetTransactionDetailsRequest) (*message.GetTransactionDetailsResponse, error) {
	hash, _ := transaction.StringToHash(blockchain.Check0x(in.Hash))
	tx, err := g.Bc.GetTransactionByHash(hash)
	if err != nil {
		st, err := g.Tp.GetTxByHash(in.Hash)
		if err != nil {
			return nil, err
		}
		tx = transaction.NewFinishedTransaction(st, nil, 0)
	}

	logger.Info("GetTransactionDetails", zap.Any("tx", *tx))

	replay := &message.GetTransactionDetailsResponse{
		Version:   int64(tx.Version),
		Type:      int32(tx.Type),
		From:      tx.From.Hex(),
		To:        tx.To.Hex(),
		Nonce:     tx.Nonce,
		Amount:    tx.Amount.String(),
		GasLimit:  tx.GasLimit.String(),
		GasPrice:  tx.GasPrice.String(),
		GasFeeCap: tx.GasFeeCap.String(),
		Input:     tx.Input,
		Signature: tx.Signature,
		GasUsed:   tx.GasUsed.String(),
		BlockNum:  tx.BlockNum,
	}

	return replay, nil
}

func (g *Greeter) Sign(ctx context.Context, in *message.SginRequest) (*message.SginResponse, error) {
	privData, err := hex.DecodeString(in.Priv)
	if err != nil {
		return nil, errors.Wrap(err, "priv key error")
	}

	priv, err := crypto.ToECDSA(privData)
	if err != nil {
		return nil, errors.Wrap(err, "priv key error")
	}

	from, to := common.HexToAddress(in.Utx.From), common.HexToAddress(in.Utx.To)
	amount, ok := new(big.Int).SetString(in.Utx.Amount, 10)
	if !ok {
		return nil, errors.Wrap(err, "failed to parse amount")
	}

	limit, ok := new(big.Int).SetString(in.Utx.GasLimit, 10)
	if !ok {
		return nil, errors.Wrap(err, "failed to pasr gas limit")
	}
	feeGap, ok := new(big.Int).SetString(in.Utx.GasFeeCap, 10)
	if !ok {
		return nil, errors.Wrap(err, "failed to parse gas fee cap")
	}

	price, ok := new(big.Int).SetString(in.Utx.GasPrice, 10)
	if !ok {
		return nil, errors.Wrap(err, "failed to parse fas price")
	}

	tx := &transaction.Transaction{
		From:   &from,
		To:     &to,
		Amount: amount,
		Nonce:  in.Utx.Nonce,
		Type:   transaction.TransferTransaction,

		GasLimit:  limit,
		GasFeeCap: feeGap,
		GasPrice:  price,
	}

	signature, err := crypto.Sign(tx.SignHash(), priv)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign")
	}

	return &message.SginResponse{Signature: hex.EncodeToString(signature)}, nil
}
