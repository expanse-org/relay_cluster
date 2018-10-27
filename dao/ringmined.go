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

package dao

import (
	"github.com/expanse-org/relay-lib/types"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

// todo(fuk): delete field isRinghashReserved
type RingMinedEvent struct {
	ID                 int    `gorm:"column:id;primary_key" json:"id"`
	Protocol           string `gorm:"column:contract_address;type:varchar(42)" json:"protocol"`
	DelegateAddress    string `gorm:"column:delegate_address;type:varchar(42)" json:"delegateAddress"`
	RingIndex          string `gorm:"column:ring_index;type:varchar(40)" json:"ringIndex"`
	RingHash           string `gorm:"column:ring_hash;type:varchar(82)" json:"ringHash"`
	TxHash             string `gorm:"column:tx_hash;type:varchar(82)" json:"txHash"`
	OrderHashList      string `gorm:"column:order_hash_list;type:text"`
	Miner              string `gorm:"column:miner;type:varchar(42);" json:"miner"`
	FeeRecipient       string `gorm:"column:fee_recipient;type:varchar(42)" json:"feeRecipient"`
	IsRinghashReserved bool   `gorm:"column:is_ring_hash_reserved;" json:"isRinghashReserved"`
	BlockNumber        int64  `gorm:"column:block_number;type:bigint" json:"blockNumber"`
	TotalPexFee        string `gorm:"column:total_pex_fee;type:varchar(40)" json:"totalPexFee"`
	TradeAmount        int    `gorm:"column:trade_amount" json:"tradeAmount"`
	Time               int64  `gorm:"column:time;type:bigint" json:"timestamp"`
	Status             uint8  `gorm:"column:status;type:tinyint(4)"`
	Fork               bool   `gorm:"column:fork"`
	GasLimit           string `gorm:"column:gas_limit;type:varchar(50)"`
	GasUsed            string `gorm:"column:gas_used;type:varchar(50)"`
	GasPrice           string `gorm:"column:gas_price;type:varchar(50)"`
	Err                string `gorm:"column:err;type:text" json:"err"`
}

func (r *RingMinedEvent) ConvertDown(event *types.RingMinedEvent) error {
	r.RingIndex = event.RingIndex.String()
	r.TotalPexFee = event.TotalPexFee.String()
	r.Protocol = event.Protocol.Hex()
	r.DelegateAddress = event.DelegateAddress.Hex()
	r.Miner = event.Miner.Hex()
	r.FeeRecipient = event.FeeRecipient.Hex()
	r.RingHash = event.Ringhash.Hex()
	r.TxHash = event.TxHash.Hex()
	r.BlockNumber = event.BlockNumber.Int64()
	r.Time = event.BlockTime
	r.TradeAmount = event.TradeAmount
	r.Status = uint8(event.Status)
	r.GasLimit = event.GasLimit.String()
	r.GasUsed = event.GasUsed.String()
	r.GasPrice = event.GasPrice.String()
	r.Err = ""
	r.Fork = false

	return nil
}

func (r *RingMinedEvent) ConvertUp(event *types.RingMinedEvent) error {
	event.RingIndex, _ = new(big.Int).SetString(r.RingIndex, 0)
	event.TotalPexFee, _ = new(big.Int).SetString(r.TotalPexFee, 0)
	event.Ringhash = common.HexToHash(r.RingHash)
	event.TxHash = common.HexToHash(r.TxHash)
	event.Miner = common.HexToAddress(r.Miner)
	event.FeeRecipient = common.HexToAddress(r.FeeRecipient)
	event.BlockNumber = big.NewInt(r.BlockNumber)
	event.BlockTime = r.Time
	event.TradeAmount = r.TradeAmount
	event.Protocol = common.HexToAddress(r.Protocol)
	event.DelegateAddress = common.HexToAddress(r.DelegateAddress)
	event.Status = types.TxStatus(r.Status)
	event.GasLimit, _ = new(big.Int).SetString(r.GasLimit, 0)
	event.GasUsed, _ = new(big.Int).SetString(r.GasUsed, 0)
	event.GasPrice, _ = new(big.Int).SetString(r.GasPrice, 0)
	event.Err = r.Err

	return nil
}

func (r *RingMinedEvent) FromSubmitRingMethod(event *types.SubmitRingMethodEvent) error {
	r.Protocol = event.Protocol.Hex()
	r.DelegateAddress = event.DelegateAddress.Hex()
	r.TxHash = event.TxHash.Hex()
	r.BlockNumber = event.BlockNumber.Int64()
	r.Status = uint8(event.Status)
	r.GasLimit = event.GasLimit.String()
	r.GasUsed = event.GasUsed.String()
	r.GasPrice = event.GasPrice.String()
	r.Err = event.Err

	var list []common.Hash
	for _, v := range event.OrderList {
		list = append(list, v.Hash)
	}
	r.OrderHashList = MarshalHashListToStr(list)

	return nil
}

func (r *RingMinedEvent) GetOrderHashList() []common.Hash {
	return UnmarshalStrToHashList(r.OrderHashList)
}

func (s *RdsService) FindRingMined(txhash string) (*RingMinedEvent, error) {
	var (
		model RingMinedEvent
		err   error
	)

	err = s.Db.Where("tx_hash=?", txhash).Where("fork = ?", false).First(&model).Error

	return &model, err
}

func (s *RdsService) RollBackRingMined(from, to int64) error {
	return s.Db.Model(&RingMinedEvent{}).Where("block_number > ? and block_number <= ?", from, to).Update("fork", true).Error
}

func (s *RdsService) RingMinedPageQuery(query map[string]interface{}, pageIndex, pageSize int) (res PageResult, err error) {
	ringMined := make([]RingMinedEvent, 0)
	res = PageResult{PageIndex: pageIndex, PageSize: pageSize, Data: make([]interface{}, 0)}

	err = s.Db.Where(query).Where("fork = ?", false).Order("time desc").Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&ringMined).Error

	if err != nil {
		return res, err
	}
	err = s.Db.Model(&RingMinedEvent{}).Where(query).Where("fork = ?", false).Count(&res.Total).Error
	if err != nil {
		return res, err
	}

	for _, rm := range ringMined {
		res.Data = append(res.Data, rm)
	}
	return
}

func (s *RdsService) GetRingminedMethods(lastId int, limit int) ([]RingMinedEvent, error) {
	var (
		list []RingMinedEvent
		err  error
	)

	err = s.Db.Where("id > ?", lastId).
		Find(&list).
		Limit(limit).
		Error

	return list, err
}

func (s *RdsService) IsMiner(miner common.Address) bool {
	var data RingMinedEvent
	err := s.Db.Where("miner=?", miner.Hex()).Order("id DESC").First(&data).Error
	if err != nil {
		return false
	}
	return true
}
