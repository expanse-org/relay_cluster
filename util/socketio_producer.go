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

package util

import (
	"github.com/expanse-org/relay-lib/kafka"
	"log"
)

var socketIOProducer *kafka.MessageProducer

func Initialize(brokers []string) {
	if socketIOProducer == nil {
		socketIOProducer = &kafka.MessageProducer{}
		if err := socketIOProducer.Initialize(brokers); nil != err {
			log.Fatalf("Failed init producerWrapped %s", err.Error())
		}
	}
}

func ProducerSocketIOMessage(eventKey string, data interface{}) error {
	_, _, err := socketIOProducer.SendMessage(eventKey, data, "1")
	return err
}

func ProducerNormalMessage(topic string, data interface{}) error {
	_, _, err := socketIOProducer.SendMessage(topic, data, "1")
	return err
}
