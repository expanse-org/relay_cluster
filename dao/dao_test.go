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

package dao_test

import (
	"github.com/expanse-org/relay_cluster/dao"
	"github.com/expanse-org/relay_cluster/test"
	"github.com/ethereum/go-ethereum/common"
	"testing"
)

// gorm没有合适的insert返回id的方法，但是create的data本身的ID会被赋值
func TestCreate(t *testing.T) {
	rds := test.Rds()
	data := &dao.RingMinedEvent{TxHash: "0x123"}
	if err := rds.Db.Create(data).Error; err != nil {
		t.Fatalf(err.Error())
	} else {
		t.Log(data.ID)
	}
}

// 找不到数据会报error:"record not found"
func TestFindError(t *testing.T) {
	// 数据不存在, err不为空
	rds := test.Rds()
	txhash := "0x125"
	if model, err := rds.FindRingMined(txhash); err != nil {
		t.Log(err.Error())
	} else {
		t.Log(model.ID)
	}
}

func TestRdsService_GetMaxNonce(t *testing.T) {
	rds := test.Rds()
	owner := common.HexToAddress("0xb1018949b241D76A1AB2094f473E9bEfeAbB5Ea")
	nonce, err := rds.GetMaxNonce(owner)
	if err != nil {
		t.Logf(err.Error())
	} else {
		t.Logf("nonce is:%d", nonce.Int64())
	}
}

func TestRdsService_GetMaxSuccessNonce(t *testing.T) {
	rds := test.Rds()
	owner := common.HexToAddress("0xb1018949b241D76A1AB2094f473E9bEfeAbB5Ead")
	nonce, err := rds.GetMaxSuccessNonce(owner)
	if err != nil {
		t.Logf(err.Error())
	} else {
		t.Logf("nonce is:%d", nonce.Int64())
	}
}
