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

package node

import (
	"errors"
	"os"
	"reflect"

	"github.com/expanse-org/relay_cluster/accountmanager"
	"github.com/expanse-org/relay_cluster/gateway"
	"github.com/expanse-org/relay_cluster/market"
	ordermanager "github.com/expanse-org/relay_cluster/ordermanager/common"
	"github.com/expanse-org/relay_cluster/usermanager"
	"github.com/expanse-org/relay-lib/cache/redis"
	"github.com/expanse-org/relay-lib/cloudwatch"
	"github.com/expanse-org/relay-lib/dao"
	"github.com/expanse-org/relay-lib/eth/accessor"
	"github.com/expanse-org/relay-lib/eth/loopringaccessor"
	"github.com/expanse-org/relay-lib/kafka"
	"github.com/expanse-org/relay-lib/marketcap"
	util "github.com/expanse-org/relay-lib/marketutil"
	"github.com/expanse-org/relay-lib/motan"
	"github.com/expanse-org/relay-lib/sns"
	"github.com/expanse-org/relay-lib/zklock"
	"github.com/naoina/toml"
	"go.uber.org/zap"
)

func LoadConfig(file string) *GlobalConfig {
	if "" == file {
		dir, _ := os.Getwd()
		file = dir + "/config/relay.toml"
	}

	io, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer io.Close()

	c := &GlobalConfig{}
	if err := toml.NewDecoder(io).Decode(c); err != nil {
		panic(err)
	}
	return c
}

type GlobalConfig struct {
	Title            string `required:"true"`
	Log              zap.Config
	Mysql            dao.MysqlOptions
	Redis            redis.RedisOptions
	OrderManager     ordermanager.OrderManagerOptions
	Gateway          gateway.GateWayOptions
	Accessor         accessor.AccessorOptions
	LoopringProtocol loopringaccessor.LoopringProtocolOptions
	Market           util.MarketOptions
	MarketCap        marketcap.MarketCapOptions
	GatewayFilters   gateway.GatewayFiltersOptions
	UserManager      usermanager.UserManagerOptions
	ZkLock           zklock.ZkLockConfig
	Sns              sns.SnsConfig
	Kafka            kafka.KafkaOptions
	MotanServer      motan.MotanServerOptions
	Jsonrpc          gateway.JsonrpcOptions
	Websocket        gateway.WebsocketOptions
	AccountManager   accountmanager.AccountManagerOptions
	MyToken          market.MyTokenConfig
	CloudWatch       cloudwatch.CloudWatchConfig
}

func Validator(cv reflect.Value) (bool, error) {
	for i := 0; i < cv.NumField(); i++ {
		cvt := cv.Type().Field(i)

		if cv.Field(i).Type().Kind() == reflect.Struct {
			if res, err := Validator(cv.Field(i)); nil != err {
				return res, err
			}
		} else {
			if "true" == cvt.Tag.Get("required") {
				if !isSet(cv.Field(i)) {
					return false, errors.New("The field " + cvt.Name + " in config must be setted")
				}
			}
		}
	}

	return true, nil
}

func isSet(v reflect.Value) bool {
	switch v.Type().Kind() {
	case reflect.Invalid:
		return false
	case reflect.String:
		return v.String() != ""
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() != 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() != 0
	case reflect.Map:
		return len(v.MapKeys()) != 0
	}
	return true
}
