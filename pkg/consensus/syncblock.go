//Sync  Blockchain
package consensus

import (
	"encoding/hex"
	"math/big"
	"time"

	"metechain/pkg/block"
	"metechain/pkg/logger"

	"go.uber.org/zap"
)

const (
	MaxEqealBlockWeight = 10
	MaxExpiration       = time.Hour
	MaxOrphanBlocks     = 200
)

//ProcessBlock is management block function
func (b *BlockChain) ProcessBlock(newblock *block.Block, globalDifficulty *big.Int) bool {

	//newblcok hash is exist
	defer logger.Info(" ProcessBlock  end ", zap.Uint64("height", newblock.Height), zap.String("hash", hex.EncodeToString(newblock.Hash)))

	if b.BlockExists(newblock.Hash) {
		logger.SugarLogger.Info("Block is exist", hex.EncodeToString(newblock.Hash))
		return false
	}

	hash := BytesToHash(newblock.Hash)
	if _, exist := b.Oranphs[hash]; exist {
		logger.Info("orphan is exist", zap.String("hash", hex.EncodeToString(newblock.Hash)))
		return false
	}

	if !b.BlockExists(newblock.PrevHash) {
		logger.Info("prevhash not exist")
		b.AddOrphanBlock(newblock)
		return false
	}

	//maybeAcceptBlock return longest chain flag
	succ, mainChain := b.maybeAcceptBlock(newblock)
	if !succ {
		return false
	}
	ok := b.ProcessOrphan(newblock)
	if ok {
		mainChain = ok
	}

	return mainChain
}

func bigToCompact(n *big.Int) uint32 {
	if n.Sign() == 0 {
		return 0
	}
	var mantissa uint32
	exponent := uint(len(n.Bytes()))
	if exponent <= 3 {
		mantissa = uint32(n.Bits()[0])
		mantissa <<= 8 * (3 - exponent)
	} else {
		tn := new(big.Int).Set(n)
		mantissa = uint32(tn.Rsh(tn, 8*(exponent-3)).Bits()[0])
	}
	if mantissa&0x00800000 != 0 {
		mantissa >>= 8
		exponent++
	}
	compact := uint32(exponent<<24) | mantissa
	if n.Sign() < 0 {
		compact |= 0x00800000
	}
	return compact
}
