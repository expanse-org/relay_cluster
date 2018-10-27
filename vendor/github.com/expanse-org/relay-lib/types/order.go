/*

  Copyright 2017 Loopring Project Ltd (Loopring Foundation).

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.

*/

package types

import (
	"github.com/expanse-org/relay-lib/crypto"
	"github.com/expanse-org/relay-lib/log"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"time"
)

type OrderStatus uint8

const (
	ORDER_UNKNOWN     OrderStatus = 0
	ORDER_NEW         OrderStatus = 1
	ORDER_PARTIAL     OrderStatus = 2
	ORDER_FINISHED    OrderStatus = 3
	ORDER_CANCEL      OrderStatus = 4
	ORDER_CUTOFF      OrderStatus = 5
	ORDER_EXPIRE      OrderStatus = 6
	ORDER_PENDING     OrderStatus = 7
	ORDER_CANCELLING  OrderStatus = 8
	ORDER_CUTOFFING   OrderStatus = 9
	ORDER_FLEX_CANCEL OrderStatus = 10
	//ORDER_BALANCE_INSUFFICIENT   OrderStatus = 9
	//ORDER_ALLOWANCE_INSUFFICIENT OrderStatus = 10
)

const (
	ORDER_TYPE_MARKET = "market_order"
	ORDER_TYPE_P2P    = "p2p_order"
)

//go:generate gencodec -type Order -field-override orderMarshaling -out gen_order_json.go
type Order struct {
	Protocol              common.Address             `json:"protocol" gencodec:"required"`        // 智能合约地址
	DelegateAddress       common.Address             `json:"delegateAddress" gencodec:"required"` // 智能合约地址
	AuthAddr              common.Address             `json:"authAddr" gencodec:"required"`        //
	AuthPrivateKey        crypto.EthPrivateKeyCrypto `json:"authPrivateKey" gencodec:"required"`  //
	WalletAddress         common.Address             `json:"walletAddress" gencodec:"required"`
	TokenS                common.Address             `json:"tokenS" gencodec:"required"`     // 卖出erc20代币智能合约地址
	TokenB                common.Address             `json:"tokenB" gencodec:"required"`     // 买入erc20代币智能合约地址
	AmountS               *big.Int                   `json:"amountS" gencodec:"required"`    // 卖出erc20代币数量上限
	AmountB               *big.Int                   `json:"amountB" gencodec:"required"`    // 买入erc20代币数量上限
	ValidSince            *big.Int                   `json:"validSince" gencodec:"required"` //
	ValidUntil            *big.Int                   `json:"validUntil" gencodec:"required"` // 订单过期时间
	PexFee                *big.Int                   `json:"pexFee" `                        // 交易总费用,部分成交的费用按该次撮合实际卖出代币额与比例计算
	BuyNoMoreThanAmountB  bool                       `json:"buyNoMoreThanAmountB" gencodec:"required"`
	MarginSplitPercentage uint8                      `json:"marginSplitPercentage" gencodec:"required"` // 不为0时支付给交易所的分润比例，否则视为100%
	V                     uint8                      `json:"v" gencodec:"required"`
	R                     Bytes32                    `json:"r" gencodec:"required"`
	S                     Bytes32                    `json:"s" gencodec:"required"`
	Price                 *big.Rat                   `json:"price"`
	Owner                 common.Address             `json:"owner"`
	Hash                  common.Hash                `json:"hash"`
	Market                string                     `json:"market"`
	CreateTime            int64                      `json:"createTime"`
	PowNonce              uint64                     `json:"powNonce"`
	Side                  string                     `json:"side"`
	OrderType             string                     `json:"orderType"`
}

type orderMarshaling struct {
	AmountS    *Big
	AmountB    *Big
	ValidSince *Big
	ValidUntil *Big
	PexFee     *Big
}

//go:generate gencodec -type OrderJsonRequest -field-override orderJsonRequestMarshaling -out gen_order_request_json.go
type OrderJsonRequest struct {
	Protocol        common.Address             `json:"protocol" gencodec:"required"`        // 智能合约地址
	DelegateAddress common.Address             `json:"delegateAddress" gencodec:"required"` // 智能合约地址
	TokenS          common.Address             `json:"tokenS" gencodec:"required"`          // 卖出erc20代币智能合约地址
	TokenB          common.Address             `json:"tokenB" gencodec:"required"`          // 买入erc20代币智能合约地址
	AuthAddr        common.Address             `json:"authAddr" gencodec:"required"`        //
	AuthPrivateKey  crypto.EthPrivateKeyCrypto `json:"authPrivateKey"`                      //
	WalletAddress   common.Address             `json:"walletAddress" gencodec:"required"`
	AmountS         *big.Int                   `json:"amountS" gencodec:"required"`    // 卖出erc20代币数量上限
	AmountB         *big.Int                   `json:"amountB" gencodec:"required"`    // 买入erc20代币数量上限
	ValidSince      *big.Int                   `json:"validSince" gencodec:"required"` //
	ValidUntil      *big.Int                   `json:"validUntil" gencodec:"required"` // 订单过期时间
	// Salt                  int64          `json:"salt" gencodec:"required"`
	PexFee                *big.Int       `json:"pexFee" ` // 交易总费用,部分成交的费用按该次撮合实际卖出代币额与比例计算
	BuyNoMoreThanAmountB  bool           `json:"buyNoMoreThanAmountB" gencodec:"required"`
	MarginSplitPercentage uint8          `json:"marginSplitPercentage" gencodec:"required"` // 不为0时支付给交易所的分润比例，否则视为100%
	V                     uint8          `json:"v" gencodec:"required"`
	R                     Bytes32        `json:"r" gencodec:"required"`
	S                     Bytes32        `json:"s" gencodec:"required"`
	Price                 *big.Rat       `json:"price"`
	Owner                 common.Address `json:"owner"`
	Hash                  common.Hash    `json:"hash"`
	CreateTime            int64          `json:"createTime"`
	PowNonce              uint64         `json:"powNonce"`
	Side                  string         `json:"side"`
	OrderType             string         `json:"orderType"`
}

type orderJsonRequestMarshaling struct {
	AmountS    *Big
	AmountB    *Big
	ValidSince *Big
	ValidUntil *Big
	PexFee     *Big
}

func (o *Order) GenerateHash() common.Hash {
	h := &common.Hash{}

	buyNoMoreThanAmountB := byte(0)
	if o.BuyNoMoreThanAmountB {
		buyNoMoreThanAmountB = byte(1)
	}

	hashBytes := crypto.GenerateHash(
		o.DelegateAddress.Bytes(),
		o.Owner.Bytes(),
		o.TokenS.Bytes(),
		o.TokenB.Bytes(),
		o.WalletAddress.Bytes(),
		o.AuthAddr.Bytes(),
		common.LeftPadBytes(o.AmountS.Bytes(), 32),
		common.LeftPadBytes(o.AmountB.Bytes(), 32),
		common.LeftPadBytes(o.ValidSince.Bytes(), 32),
		common.LeftPadBytes(o.ValidUntil.Bytes(), 32),
		common.LeftPadBytes(o.PexFee.Bytes(), 32),
		[]byte{buyNoMoreThanAmountB},
		[]byte{byte(o.MarginSplitPercentage)},
	)

	h.SetBytes(hashBytes)
	return *h
}

func (o *Order) GenerateAndSetSignature(singerAddr common.Address) error {
	if IsZeroHash(o.Hash) {
		o.Hash = o.GenerateHash()
	}

	if sig, err := crypto.Sign(o.Hash.Bytes(), singerAddr); nil != err {
		return err
	} else {
		v, r, s := crypto.SigToVRS(sig)
		o.V = uint8(v)
		o.R = BytesToBytes32(r)
		o.S = BytesToBytes32(s)
		return nil
	}
}

func (o *Order) ValidateSignatureValues() bool {
	return crypto.ValidateSignatureValues(byte(o.V), o.R.Bytes(), o.S.Bytes())
}

func (o *Order) SignerAddress() (common.Address, error) {
	address := &common.Address{}
	if IsZeroHash(o.Hash) {
		o.Hash = o.GenerateHash()
	}

	sig, _ := crypto.VRSToSig(o.V, o.R.Bytes(), o.S.Bytes())

	if addressBytes, err := crypto.SigToAddress(o.Hash.Bytes(), sig); nil != err {
		log.Errorf("type,order signer address error:%s", err.Error())
		return *address, err
	} else {
		address.SetBytes(addressBytes)
		return *address, nil
	}
}

func (o *Order) GeneratePrice() {
	o.Price = new(big.Rat).SetFrac(o.AmountS, o.AmountB)
}

//RateAmountS、FeeSelection 需要提交到contract
type FilledOrder struct {
	OrderState       OrderState `json:"orderState" gencodec:"required"`
	FeeSelection     uint8      `json:"feeSelection"`     //0 -> pex
	RateAmountS      *big.Rat   `json:"rateAmountS"`      //提交需要
	AvailableAmountS *big.Rat   `json:"availableAmountS"` //需要，也是用于计算fee
	AvailableAmountB *big.Rat   //需要，也是用于计算fee
	FillAmountS      *big.Rat   `json:"fillAmountS"`
	FillAmountB      *big.Rat   `json:"fillAmountB"` //计算需要
	PexReward        *big.Rat   `json:"pexReward"`
	PexFee           *big.Rat   `json:"pexFee"`
	LegalPexFee      *big.Rat   `json:"legalPexFee"`
	FeeS             *big.Rat   `json:"feeS"`
	LegalFeeS        *big.Rat   `json:"legalFeeS"`
	LegalFee         *big.Rat   `json:"legalFee"` //法币计算的fee

	SPrice *big.Rat `json:"SPrice"`
	BPrice *big.Rat `json:"BPrice"`

	AvailablePexBalance    *big.Rat
	AvailableTokenSBalance *big.Rat
}

func ConvertOrderStateToFilledOrder(orderState OrderState, pexBalance, tokenSBalance *big.Rat, pexAddress common.Address) *FilledOrder {
	filledOrder := &FilledOrder{}
	filledOrder.OrderState = orderState
	filledOrder.AvailablePexBalance = new(big.Rat).Set(pexBalance)
	filledOrder.AvailableTokenSBalance = new(big.Rat).Set(tokenSBalance)

	filledOrder.AvailableAmountS, filledOrder.AvailableAmountB = filledOrder.OrderState.RemainedAmount()
	sellPrice := new(big.Rat).SetFrac(filledOrder.OrderState.RawOrder.AmountS, filledOrder.OrderState.RawOrder.AmountB)

	availableBalance := new(big.Rat).Set(filledOrder.AvailableTokenSBalance)
	if availableBalance.Cmp(filledOrder.AvailableAmountS) < 0 {
		filledOrder.AvailableAmountS = availableBalance
		filledOrder.AvailableAmountB.Mul(filledOrder.AvailableAmountS, new(big.Rat).Inv(sellPrice))
	}
	if filledOrder.OrderState.RawOrder.BuyNoMoreThanAmountB {
		filledOrder.AvailableAmountS.Mul(filledOrder.AvailableAmountB, sellPrice)
	} else {
		filledOrder.AvailableAmountB.Mul(filledOrder.AvailableAmountS, new(big.Rat).Inv(sellPrice))
	}

	if orderState.RawOrder.TokenB == pexAddress && pexBalance.Cmp(filledOrder.AvailableAmountB) < 0 {
		filledOrder.AvailablePexBalance.Set(filledOrder.AvailableAmountB)
	}
	return filledOrder
}

// 从[]byte解析时使用json.Unmarshal
type OrderState struct {
	RawOrder         Order       `json:"rawOrder"`
	UpdatedBlock     *big.Int    `json:"updatedBlock"`
	DealtAmountS     *big.Int    `json:"dealtAmountS"`
	DealtAmountB     *big.Int    `json:"dealtAmountB"`
	SplitAmountS     *big.Int    `json:"splitAmountS"`
	SplitAmountB     *big.Int    `json:"splitAmountB"`
	CancelledAmountS *big.Int    `json:"cancelledAmountS"`
	CancelledAmountB *big.Int    `json:"cancelledAmountB"`
	Status           OrderStatus `json:"status"`
	BroadcastTime    int         `json:"broadcastTime"`
}

func (state *OrderState) IsOrderFullFinished(tokenSPrice, tokenBPrice *big.Rat, dustValue *big.Rat) bool {
	remainedAmountS, _ := state.RemainedAmount()
	remainedValue := new(big.Rat)
	remainedValue.Mul(remainedAmountS, tokenSPrice)
	return isValueDusted(remainedValue, dustValue)
}

func isValueDusted(remainedValue, dustValue *big.Rat) bool {
	return remainedValue.Cmp(dustValue) <= 0
}

type OrderDelayList struct {
	OrderHash    []common.Hash
	DelayedCount int64
}

func (ord *OrderState) IsExpired() bool {
	if (ord.Status == ORDER_NEW || ord.Status == ORDER_PARTIAL) && ord.RawOrder.ValidUntil.Int64() < time.Now().Unix() {
		return true
	}
	return false
}

func (ord *OrderState) IsEffective() bool {
	if (ord.Status == ORDER_NEW || ord.Status == ORDER_PARTIAL) &&
		ord.RawOrder.ValidSince.Int64() <= time.Now().Unix() &&
		ord.RawOrder.ValidUntil.Int64() > time.Now().Unix() {
		return true
	}
	return false
}

func (orderState *OrderState) RemainedAmount() (remainedAmountS *big.Rat, remainedAmountB *big.Rat) {
	remainedAmountS = new(big.Rat)
	remainedAmountB = new(big.Rat)
	if orderState.RawOrder.BuyNoMoreThanAmountB {
		reducedAmountB := new(big.Rat)
		reducedAmountB.Add(reducedAmountB, new(big.Rat).SetInt(orderState.DealtAmountB)).
			Add(reducedAmountB, new(big.Rat).SetInt(orderState.CancelledAmountB)).
			Add(reducedAmountB, new(big.Rat).SetInt(orderState.SplitAmountB))
		sellPrice := new(big.Rat).SetFrac(orderState.RawOrder.AmountS, orderState.RawOrder.AmountB)
		remainedAmountB.Sub(new(big.Rat).SetInt(orderState.RawOrder.AmountB), reducedAmountB)
		remainedAmountS.Mul(remainedAmountB, sellPrice)
	} else {
		reducedAmountS := new(big.Rat)
		reducedAmountS.Add(reducedAmountS, new(big.Rat).SetInt(orderState.DealtAmountS)).
			Add(reducedAmountS, new(big.Rat).SetInt(orderState.CancelledAmountS)).
			Add(reducedAmountS, new(big.Rat).SetInt(orderState.SplitAmountS))
		buyPrice := new(big.Rat).SetFrac(orderState.RawOrder.AmountB, orderState.RawOrder.AmountS)
		remainedAmountS.Sub(new(big.Rat).SetInt(orderState.RawOrder.AmountS), reducedAmountS)
		remainedAmountB.Mul(remainedAmountS, buyPrice)
	}

	return remainedAmountS, remainedAmountB
}

func (state *OrderState) DealtAndSplitAmount() (totalAmountS *big.Rat, totalAmountB *big.Rat) {
	totalAmountS = new(big.Rat)
	totalAmountB = new(big.Rat)

	if state.RawOrder.BuyNoMoreThanAmountB {
		totalAmountB = totalAmountB.SetInt(new(big.Int).Add(state.DealtAmountB, state.SplitAmountB))
		sellPrice := new(big.Rat).SetFrac(state.RawOrder.AmountS, state.RawOrder.AmountB)
		totalAmountS = totalAmountS.Mul(totalAmountB, sellPrice)
	} else {
		totalAmountS = totalAmountS.SetInt(new(big.Int).Add(state.DealtAmountS, state.SplitAmountS))
		buyPrice := new(big.Rat).SetFrac(state.RawOrder.AmountB, state.RawOrder.AmountS)
		totalAmountB = totalAmountB.Mul(totalAmountS, buyPrice)
	}

	return totalAmountS, totalAmountB
}

func ToOrder(request *OrderJsonRequest) *Order {
	order := &Order{}
	order.Protocol = request.Protocol
	order.DelegateAddress = request.DelegateAddress
	order.TokenS = request.TokenS
	order.TokenB = request.TokenB
	order.AmountS = request.AmountS
	order.AmountB = request.AmountB
	order.ValidSince = request.ValidSince
	order.ValidUntil = request.ValidUntil
	order.AuthAddr = request.AuthAddr
	order.AuthPrivateKey = request.AuthPrivateKey
	order.PexFee = request.PexFee
	order.BuyNoMoreThanAmountB = request.BuyNoMoreThanAmountB
	order.MarginSplitPercentage = request.MarginSplitPercentage
	order.V = request.V
	order.R = request.R
	order.S = request.S
	order.Owner = request.Owner
	order.WalletAddress = request.WalletAddress
	order.PowNonce = request.PowNonce
	order.OrderType = request.OrderType
	return order
}
