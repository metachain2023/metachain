// Package blockchain define the interface of blockchain and implement its object
package blockchain

import (
	"math/big"

	"metechain/pkg/block"
	"metechain/pkg/transaction"

	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/ethereum/go-ethereum/core/types"
)

//Blockchains interface specification of blockchain
type Blockchains interface {
	// NewBlock create a new block for the blockchain
	NewBlock([]*transaction.SignedTransaction, *common.Address) (*block.Block, error)
	// AddBlock add blocks to blockchain
	AddUncleBlock(*block.Block) error
	// AddBlock add blocks to blockchain
	AddBlock(*block.Block) error
	// DeleteBlock delete some blocks from the blockchain
	DeleteBlock(uint64) error
	// DeleteUncleBlock delete some blocks from the blockchain
	DeleteUncleBlock(*block.Block) error
	// GetBalance retrieves the balance from the given address or 0 if object not found
	GetBalance(*common.Address) (*big.Int, error)
	// GetAvailableBalance get available balance of address
	GetAvailableBalance(*common.Address) (*big.Int, error)
	// GetNonce get the nonce of the address
	GetNonce(*common.Address) (uint64, error)
	// GetHash get the hash corresponding to the block height
	GetHash(uint64) ([]byte, error)
	// GetMaxBlockHeight get maximum block height
	GetMaxBlockHeight() (uint64, error)
	// GetBlockByHeight get the block corresponding to the block height
	GetBlockByHeight(uint64) (*block.Block, error)
	// GetBlockByHash get block data through hash
	GetBlockByHash([]byte) (*block.Block, error)

	ReorganizeChain([][]byte, uint64) error
	Tip() (*block.Block, error)

	//get binding mete address by eth address
	GetBindingmeteAddress(ethAddr string) (*common.Address, error)
	//get binding eth address by mete address
	GetBindingEthAddress(meteAddr *common.Address) (string, error)
	//call contract
	CallSmartContract(contractAddr, origin, callInput, value string) (string, string, error)
	//get code
	GetCode(contractAddr string) []byte
	//set code
	SetCode(contractAddr common.Address, code []byte)
	//get storage by hash
	GetStorageAt(addr, hash string) common.Hash
	// GetTransactionByHash get the transaction corresponding to the transaction hash
	GetTransactionByHash([]byte) (*transaction.FinishedTransaction, error)
	//Get logs
	GetLogs() []*evmtypes.Log
	GetLogByHeight(height uint64) []*evmtypes.Log

	DifficultDetection(b *block.Block) error
	CheckBlockRegular(b *block.Block) error
}
