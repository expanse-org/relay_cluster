// Code generated by github.com/fjl/gencodec. DO NOT EDIT.

package types

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/expanse-org/relay-lib/crypto"
	"github.com/ethereum/go-ethereum/common"
)

var _ = (*orderJsonRequestMarshaling)(nil)

// MarshalJSON marshals as JSON.
func (o OrderJsonRequest) MarshalJSON() ([]byte, error) {
	type OrderJsonRequest struct {
		Protocol              common.Address             `json:"protocol" gencodec:"required"`
		DelegateAddress       common.Address             `json:"delegateAddress" gencodec:"required"`
		TokenS                common.Address             `json:"tokenS" gencodec:"required"`
		TokenB                common.Address             `json:"tokenB" gencodec:"required"`
		AuthAddr              common.Address             `json:"authAddr" gencodec:"required"`
		AuthPrivateKey        crypto.EthPrivateKeyCrypto `json:"authPrivateKey"`
		WalletAddress         common.Address             `json:"walletAddress" gencodec:"required"`
		AmountS               *Big                       `json:"amountS" gencodec:"required"`
		AmountB               *Big                       `json:"amountB" gencodec:"required"`
		ValidSince            *Big                       `json:"validSince" gencodec:"required"`
		ValidUntil            *Big                       `json:"validUntil" gencodec:"required"`
		LrcFee                *Big                       `json:"lrcFee" `
		BuyNoMoreThanAmountB  bool                       `json:"buyNoMoreThanAmountB" gencodec:"required"`
		MarginSplitPercentage uint8                      `json:"marginSplitPercentage" gencodec:"required"`
		V                     uint8                      `json:"v" gencodec:"required"`
		R                     Bytes32                    `json:"r" gencodec:"required"`
		S                     Bytes32                    `json:"s" gencodec:"required"`
		Price                 *big.Rat                   `json:"price"`
		Owner                 common.Address             `json:"owner"`
		Hash                  common.Hash                `json:"hash"`
		CreateTime            int64                      `json:"createTime"`
		PowNonce              uint64                     `json:"powNonce"`
		Side                  string                     `json:"side"`
		OrderType             string                     `json:"orderType"`
		P2PSide               string                     `json:"p2pSide"`
		SourceId              string                     `json:"sourceId"`
	}
	var enc OrderJsonRequest
	enc.Protocol = o.Protocol
	enc.DelegateAddress = o.DelegateAddress
	enc.TokenS = o.TokenS
	enc.TokenB = o.TokenB
	enc.AuthAddr = o.AuthAddr
	enc.AuthPrivateKey = o.AuthPrivateKey
	enc.WalletAddress = o.WalletAddress
	enc.AmountS = (*Big)(o.AmountS)
	enc.AmountB = (*Big)(o.AmountB)
	enc.ValidSince = (*Big)(o.ValidSince)
	enc.ValidUntil = (*Big)(o.ValidUntil)
	enc.LrcFee = (*Big)(o.LrcFee)
	enc.BuyNoMoreThanAmountB = o.BuyNoMoreThanAmountB
	enc.MarginSplitPercentage = o.MarginSplitPercentage
	enc.V = o.V
	enc.R = o.R
	enc.S = o.S
	enc.Price = o.Price
	enc.Owner = o.Owner
	enc.Hash = o.Hash
	enc.CreateTime = o.CreateTime
	enc.PowNonce = o.PowNonce
	enc.Side = o.Side
	enc.OrderType = o.OrderType
	enc.P2PSide = o.P2PSide
	enc.SourceId = o.SourceId
	return json.Marshal(&enc)
}

// UnmarshalJSON unmarshals from JSON.
func (o *OrderJsonRequest) UnmarshalJSON(input []byte) error {
	type OrderJsonRequest struct {
		Protocol              *common.Address             `json:"protocol" gencodec:"required"`
		DelegateAddress       *common.Address             `json:"delegateAddress" gencodec:"required"`
		TokenS                *common.Address             `json:"tokenS" gencodec:"required"`
		TokenB                *common.Address             `json:"tokenB" gencodec:"required"`
		AuthAddr              *common.Address             `json:"authAddr" gencodec:"required"`
		AuthPrivateKey        *crypto.EthPrivateKeyCrypto `json:"authPrivateKey"`
		WalletAddress         *common.Address             `json:"walletAddress" gencodec:"required"`
		AmountS               *Big                        `json:"amountS" gencodec:"required"`
		AmountB               *Big                        `json:"amountB" gencodec:"required"`
		ValidSince            *Big                        `json:"validSince" gencodec:"required"`
		ValidUntil            *Big                        `json:"validUntil" gencodec:"required"`
		LrcFee                *Big                        `json:"lrcFee" `
		BuyNoMoreThanAmountB  *bool                       `json:"buyNoMoreThanAmountB" gencodec:"required"`
		MarginSplitPercentage *uint8                      `json:"marginSplitPercentage" gencodec:"required"`
		V                     *uint8                      `json:"v" gencodec:"required"`
		R                     *Bytes32                    `json:"r" gencodec:"required"`
		S                     *Bytes32                    `json:"s" gencodec:"required"`
		Price                 *big.Rat                    `json:"price"`
		Owner                 *common.Address             `json:"owner"`
		Hash                  *common.Hash                `json:"hash"`
		CreateTime            *int64                      `json:"createTime"`
		PowNonce              *uint64                     `json:"powNonce"`
		Side                  *string                     `json:"side"`
		OrderType             *string                     `json:"orderType"`
		P2PSide               *string                     `json:"p2pSide"`
		SourceId              *string                     `json:"sourceId"`
	}
	var dec OrderJsonRequest
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.Protocol == nil {
		return errors.New("missing required field 'protocol' for OrderJsonRequest")
	}
	o.Protocol = *dec.Protocol
	if dec.DelegateAddress == nil {
		return errors.New("missing required field 'delegateAddress' for OrderJsonRequest")
	}
	o.DelegateAddress = *dec.DelegateAddress
	if dec.TokenS == nil {
		return errors.New("missing required field 'tokenS' for OrderJsonRequest")
	}
	o.TokenS = *dec.TokenS
	if dec.TokenB == nil {
		return errors.New("missing required field 'tokenB' for OrderJsonRequest")
	}
	o.TokenB = *dec.TokenB
	if dec.AuthAddr == nil {
		return errors.New("missing required field 'authAddr' for OrderJsonRequest")
	}
	o.AuthAddr = *dec.AuthAddr
	if dec.AuthPrivateKey != nil {
		o.AuthPrivateKey = *dec.AuthPrivateKey
	}
	if dec.WalletAddress == nil {
		return errors.New("missing required field 'walletAddress' for OrderJsonRequest")
	}
	o.WalletAddress = *dec.WalletAddress
	if dec.AmountS == nil {
		return errors.New("missing required field 'amountS' for OrderJsonRequest")
	}
	o.AmountS = (*big.Int)(dec.AmountS)
	if dec.AmountB == nil {
		return errors.New("missing required field 'amountB' for OrderJsonRequest")
	}
	o.AmountB = (*big.Int)(dec.AmountB)
	if dec.ValidSince == nil {
		return errors.New("missing required field 'validSince' for OrderJsonRequest")
	}
	o.ValidSince = (*big.Int)(dec.ValidSince)
	if dec.ValidUntil == nil {
		return errors.New("missing required field 'validUntil' for OrderJsonRequest")
	}
	o.ValidUntil = (*big.Int)(dec.ValidUntil)
	if dec.LrcFee != nil {
		o.LrcFee = (*big.Int)(dec.LrcFee)
	}
	if dec.BuyNoMoreThanAmountB == nil {
		return errors.New("missing required field 'buyNoMoreThanAmountB' for OrderJsonRequest")
	}
	o.BuyNoMoreThanAmountB = *dec.BuyNoMoreThanAmountB
	if dec.MarginSplitPercentage == nil {
		return errors.New("missing required field 'marginSplitPercentage' for OrderJsonRequest")
	}
	o.MarginSplitPercentage = *dec.MarginSplitPercentage
	if dec.V == nil {
		return errors.New("missing required field 'v' for OrderJsonRequest")
	}
	o.V = *dec.V
	if dec.R == nil {
		return errors.New("missing required field 'r' for OrderJsonRequest")
	}
	o.R = *dec.R
	if dec.S == nil {
		return errors.New("missing required field 's' for OrderJsonRequest")
	}
	o.S = *dec.S
	if dec.Price != nil {
		o.Price = dec.Price
	}
	if dec.Owner != nil {
		o.Owner = *dec.Owner
	}
	if dec.Hash != nil {
		o.Hash = *dec.Hash
	}
	if dec.CreateTime != nil {
		o.CreateTime = *dec.CreateTime
	}
	if dec.PowNonce != nil {
		o.PowNonce = *dec.PowNonce
	}
	if dec.Side != nil {
		o.Side = *dec.Side
	}
	if dec.OrderType != nil {
		o.OrderType = *dec.OrderType
	}
	if dec.P2PSide != nil {
		o.P2PSide = *dec.P2PSide
	}
	if dec.SourceId != nil {
		o.SourceId = *dec.SourceId
	}
	return nil
}
