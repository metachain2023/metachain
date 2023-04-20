package transaction

import (
	"bytes"
	"errors"
	"fmt"

	_ "metechain/pkg/crypto/sigs/ed25519"
	_ "metechain/pkg/crypto/sigs/secp"

	"github.com/ethereum/go-ethereum/crypto"
)

func (st *SignedTransaction) VerifySign() error {
	switch st.Type {
	case TransferTransaction, WithdrawToEthTransaction:

		sigPub, err := crypto.SigToPub(st.SignHash(), st.Signature)
		if err != nil {
			return err
		}

		sigAdde := crypto.PubkeyToAddress(*sigPub)
		if !bytes.Equal(sigAdde.Bytes(), st.From.Bytes()) {
			return fmt.Errorf("signature verification failed")
		}
	case EvmContractTransaction, EvmmeteTransaction:
		evm, err := DecodeEvmData(st.Input)
		if err != nil {
			return err
		}
		if !VerifyEthSign(evm.EthData) {
			return errors.New("verify eth transaction failed!")
		}
	}

	return nil
}
