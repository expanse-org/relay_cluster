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
	"encoding/json"
	"fmt"
	util "github.com/expanse-org/relay-lib/marketutil"
	"github.com/expanse-org/relay-lib/types"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type TransactionEntity struct {
	From        common.Address `json:"from"`
	To          common.Address `json:"to"`
	Protocol    common.Address `json:"protocol"`
	BlockNumber int64          `json:"block_number"`
	Hash        common.Hash    `json:"hash"`
	LogIndex    int64          `json:"log_index"`
	Value       *big.Int       `json:"value"`
	Content     string         `json:"content"`
	Status      types.TxStatus `json:"status"`
	GasLimit    *big.Int       `json:"gas_limit"`
	GasUsed     *big.Int       `json:"gas_used"`
	GasPrice    *big.Int       `json:"gas_price"`
	Nonce       *big.Int       `json:"nonce"`
	BlockTime   int64          `json:"block_time"`
}

type ApproveContent struct {
	Owner   string `json:"owner"`
	Spender string `json:"spender"`
	Amount  string `json:"amount"`
}

type CancelContent struct {
	OrderHash string `json:"order_hash"`
	Amount    string `json:"amount"`
}

type CutoffContent struct {
	Owner           string `json:"owner"`
	CutoffTimeStamp int64  `json:"cutoff"`
}

type CutoffPairContent struct {
	Owner           string `json:"owner"`
	Token1          string `json:"token1"`
	Token2          string `json:"token2"`
	CutoffTimeStamp int64  `json:"cutoff"`
}

type WexpWithdrawalContent struct {
	Src    string `json:"src"`
	Amount string `json:"amount"`
}

type WexpDepositContent struct {
	Dst    string `json:"dst"`
	Amount string `json:"amount"`
}

type TransferContent struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Amount   string `json:"amount"`
}

type OrderFilledContent struct {
	RingHash  string `json:"ring_hash"`
	OrderHash string `json:"order_hash"`
	Owner     string `json:"owner"`
	TokenS    string `json:"token_s"`
	TokenB    string `json:"token_b"`
	RingIndex string `json:"ring_index"`
	AmountS   string `json:"amount_s"`
	AmountB   string `json:"amount_b"`
	PexReward string `json:"pex_reward"`
	PexFee    string `json:"pex_fee"`
	SplitS    string `json:"split_s"`
	SplitB    string `json:"split_b"`
	Market    string `json:"market"`
	FillIndex string `json:"fill_index"`
}

func (tx *TransactionEntity) FromApproveEvent(src *types.ApprovalEvent) error {
	if err := tx.fullFilled(src.TxInfo); err != nil {
		return err
	}

	var content ApproveContent
	content.Owner = src.Owner.Hex()
	content.Spender = src.Spender.Hex()
	content.Amount = src.Amount.String()

	bs, err := json.Marshal(&content)
	if err != nil {
		return err
	}

	tx.Content = string(bs)
	return nil
}

func (tx *TransactionEntity) FromCancelEvent(src *types.OrderCancelledEvent) error {
	if err := tx.fullFilled(src.TxInfo); err != nil {
		return err
	}

	var content CancelContent
	content.OrderHash = src.OrderHash.Hex()
	content.Amount = src.AmountCancelled.String()

	bs, err := json.Marshal(&content)
	if err != nil {
		return err
	}

	tx.Content = string(bs)
	return nil
}

func (tx *TransactionEntity) FromCutoffEvent(src *types.CutoffEvent) error {
	if err := tx.fullFilled(src.TxInfo); err != nil {
		return err
	}

	var content CutoffContent
	content.Owner = src.Owner.Hex()
	content.CutoffTimeStamp = src.Cutoff.Int64()

	bs, err := json.Marshal(&content)
	if err != nil {
		return err
	}

	tx.Content = string(bs)
	return nil
}

func (tx *TransactionEntity) FromCutoffPairEvent(src *types.CutoffPairEvent) error {
	if err := tx.fullFilled(src.TxInfo); err != nil {
		return err
	}

	var content CutoffPairContent
	content.Owner = src.Owner.Hex()
	content.Token1 = src.Token1.Hex()
	content.Token2 = src.Token2.Hex()
	content.CutoffTimeStamp = src.Cutoff.Int64()

	bs, err := json.Marshal(&content)
	if err != nil {
		return err
	}

	tx.Content = string(bs)
	return nil
}

// 充值和提现from和to都是用户钱包自己的地址，因为合约限制了发送方msg.sender
func (tx *TransactionEntity) FromWexpDepositEvent(src *types.WexpDepositEvent) error {
	if err := tx.fullFilled(src.TxInfo); err != nil {
		return err
	}

	var content WexpDepositContent
	content.Dst = src.Dst.Hex()
	content.Amount = src.Amount.String()

	bs, err := json.Marshal(&content)
	if err != nil {
		return err
	}

	tx.Content = string(bs)
	return nil
}

func (tx *TransactionEntity) FromWexpWithdrawalEvent(src *types.WexpWithdrawalEvent) error {
	if err := tx.fullFilled(src.TxInfo); err != nil {
		return err
	}

	var content WexpWithdrawalContent
	content.Src = src.Src.Hex()
	content.Amount = src.Amount.String()

	bs, err := json.Marshal(&content)
	if err != nil {
		return err
	}

	tx.Content = string(bs)
	return nil
}

func (tx *TransactionEntity) FromTransferEvent(src *types.TransferEvent) error {
	if err := tx.fullFilled(src.TxInfo); err != nil {
		return err
	}

	var content TransferContent
	content.Sender = src.Sender.Hex()
	content.Receiver = src.Receiver.Hex()
	content.Amount = src.Amount.String()

	bs, err := json.Marshal(&content)
	if err != nil {
		return err
	}

	tx.Content = string(bs)
	return nil
}

func (tx *TransactionEntity) FromEthTransferEvent(src *types.EthTransferEvent) error {
	if err := tx.fullFilled(src.TxInfo); err != nil {
		return err
	}

	tx.Content = ""
	return nil
}

func (tx *TransactionEntity) FromUnsupportedContractEvent(src *types.UnsupportedContractEvent) error {
	if err := tx.fullFilled(src.TxInfo); err != nil {
		return err
	}

	tx.Content = ""
	return nil
}

func (tx *TransactionEntity) FromOrderFilledEvent(src *types.OrderFilledEvent) error {
	if err := tx.fullFilled(src.TxInfo); err != nil {
		return err
	}

	var content OrderFilledContent
	content.RingHash = src.Ringhash.Hex()
	content.OrderHash = src.OrderHash.Hex()
	content.TokenS = src.TokenS.Hex()
	content.TokenB = src.TokenB.Hex()
	content.RingIndex = src.RingIndex.String()
	content.AmountS = src.AmountS.String()
	content.AmountB = src.AmountB.String()
	content.PexReward = src.PexReward.String()
	content.PexFee = src.PexFee.String()
	content.SplitS = src.SplitS.String()
	content.SplitB = src.SplitB.String()
	content.Market, _ = util.WrapMarketByAddress(src.TokenB.Hex(), src.TokenS.Hex())
	content.FillIndex = src.FillIndex.String()

	bs, err := json.Marshal(&content)
	if err != nil {
		return err
	}

	tx.Content = string(bs)
	return nil
}

func (tx *TransactionEntity) fullFilled(src types.TxInfo) error {
	if src.Nonce == nil || src.GasPrice == nil || src.GasLimit == nil || src.Value == nil {
		return fmt.Errorf("transaction manager, full fill tx entity error: nonce/gas/gasPrice/value cann't be nill")
	}

	tx.Hash = src.TxHash
	tx.Protocol = src.Protocol
	tx.From = src.From
	tx.To = src.To
	tx.LogIndex = src.TxLogIndex
	tx.Value = src.Value
	tx.Status = src.Status
	tx.GasLimit = src.GasLimit
	tx.GasUsed = src.GasUsed
	tx.GasPrice = src.GasPrice
	tx.Nonce = src.Nonce
	tx.BlockTime = src.BlockTime

	if src.BlockNumber != nil {
		tx.BlockNumber = src.BlockNumber.Int64()
	} else {
		tx.BlockNumber = 0
	}
	if src.GasUsed != nil {
		tx.GasUsed = src.GasUsed
	} else {
		tx.GasUsed = big.NewInt(0)
	}

	return nil
}

// Compare return true: is the same
func (tx *TransactionEntity) Compare(src *TransactionEntity) bool {
	if tx.Hash != src.Hash {
		return false
	}
	if tx.LogIndex != src.LogIndex {
		return false
	}
	if tx.Nonce != src.Nonce {
		return false
	}
	if tx.Status != src.Status {
		return false
	}
	return true
}
