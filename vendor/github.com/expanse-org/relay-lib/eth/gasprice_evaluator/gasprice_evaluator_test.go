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

package gasprice_evaluator_test

import (
	"encoding/json"
	"github.com/expanse-org/relay-lib/cache"
	"github.com/expanse-org/relay-lib/cache/redis"
	"github.com/expanse-org/relay-lib/eth/accessor"
	"github.com/expanse-org/relay-lib/eth/gasprice_evaluator"
	"github.com/expanse-org/relay-lib/log"
	"github.com/expanse-org/relay-lib/zklock"
	"go.uber.org/zap"
	"math/big"
	"testing"
	"time"
)

func init() {
	logConfig := `{
	  "level": "debug",
	  "development": false,
	  "encoding": "json",
	  "outputPaths": ["stdout"],
	  "errorOutputPaths": ["stderr"],
	  "encoderConfig": {
	    "messageKey": "message",
	    "levelKey": "level",
	    "levelEncoder": "lowercase"
	  }
	}`
	rawJSON := []byte(logConfig)

	var (
		cfg zap.Config
		err error
	)
	if err = json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	log.Initialize(cfg)
	cache.NewCache(redis.RedisOptions{Host: "127.0.0.1", Port: "6379"})

	accessor.Initialize(accessor.AccessorOptions{RawUrls: []string{"http://13.230.23.98:8545"}})

	zkconfig := zklock.ZkLockConfig{}
	zkconfig.ZkServers = "127.0.0.1:2181"
	zkconfig.ConnectTimeOut = 10000
	zklock.Initialize(zkconfig)
}

func TestInitGasPriceEvaluator(t *testing.T) {
	gasprice_evaluator.InitGasPriceEvaluator()
	time.Sleep(5 * time.Minute)
}

func TestEstimateGasPrice(t *testing.T) {
	gasPrice := gasprice_evaluator.EstimateGasPrice(big.NewInt(int64(100000000)), big.NewInt(int64(1000000000000000)))
	t.Log(gasPrice.String())
}
