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

package accountmanager

import (
	"errors"
	rcache "github.com/expanse-org/relay-lib/cache"
	"github.com/expanse-org/relay-lib/eth/loopringaccessor"
	"github.com/expanse-org/relay-lib/eventemitter"
	"github.com/expanse-org/relay-lib/log"

	"fmt"
	"github.com/Hamzaahmed742/relay_cluster/util"
	"github.com/expanse-org/relay-lib/kafka"
	"github.com/expanse-org/relay-lib/marketutil"
	"github.com/expanse-org/relay-lib/types"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type AccountManager struct {
	tokenCacheDuration int64
	ethCacheDuration   int64
	//maxBlockLength uint64
	cachedBlockCount *big.Int
	//block           *ChangedOfBlock
	producerWrapped *kafka.MessageProducer
}

func isPackegeReady() error {
	if !log.IsInit() {
		return fmt.Errorf("log must be init first")
	}
	if !rcache.IsInit() || !loopringaccessor.IsInit() || !marketutil.IsInit() {
		return fmt.Errorf("cache、loopringaccessor、 marketutil must be init first")
	}
	return nil
}

func Initialize(options *AccountManagerOptions, brokers []string) AccountManager {
	if nil != accManager {
		log.Fatalf("AccountManager has been init")
	}
	if err := isPackegeReady(); nil != err {
		log.Fatalf(err.Error())
	}

	accountManager := AccountManager{}
	if options.CacheDuration > 0 {
		accountManager.tokenCacheDuration = options.CacheDuration
	} else {
		accountManager.tokenCacheDuration = 3600 * 24 * 100
	}
	if options.EthCacheDuration > 0 {
		accountManager.ethCacheDuration = options.EthCacheDuration
	} else {
		accountManager.ethCacheDuration = 1800
	}
	//accountManager.maxBlockLength = 3000
	accountManager.cachedBlockCount = big.NewInt(int64(500))

	//b := &ChangedOfBlock{}
	//b.cachedDuration = big.NewInt(int64(500))
	//accountManager.block = b

	if len(brokers) > 0 {
		accountManager.producerWrapped = &kafka.MessageProducer{}
		if err := accountManager.producerWrapped.Initialize(brokers); nil != err {
			log.Fatalf("Failed init producerWrapped %s", err.Error())
		}
	} else {
		log.Errorf("There is not brokers of kafka to send msg.")
	}
	accManager = &accountManager
	return accountManager
}

func sendBlockEndKafkaMsg(msg interface{}) error {
	topic, key := kafka.Kafka_Topic_RelayCluster_BlockEnd, "relaycluster_blockend"
	_, _, err := accManager.producerWrapped.SendMessage(topic, msg, key)
	return err
}

func (accountManager *AccountManager) Start() {
	transferWatcher := &eventemitter.Watcher{Concurrent: false, Handle: accountManager.handleTokenTransfer}
	approveWatcher := &eventemitter.Watcher{Concurrent: false, Handle: accountManager.handleApprove}
	wethDepositWatcher := &eventemitter.Watcher{Concurrent: false, Handle: accountManager.handleWethDeposit}
	wethWithdrawalWatcher := &eventemitter.Watcher{Concurrent: false, Handle: accountManager.handleWethWithdrawal}
	blockForkWatcher := &eventemitter.Watcher{Concurrent: false, Handle: accountManager.handleBlockFork}
	blockEndWatcher := &eventemitter.Watcher{Concurrent: false, Handle: accountManager.handleBlockEnd}
	blockNewWatcher := &eventemitter.Watcher{Concurrent: false, Handle: accountManager.handleBlockNew}
	ethTransferWatcher := &eventemitter.Watcher{Concurrent: false, Handle: accountManager.handleEthTransfer}
	cancelOrderWather := &eventemitter.Watcher{Concurrent: false, Handle: accountManager.handleCancelOrder}
	cutoffAllWatcher := &eventemitter.Watcher{Concurrent: false, Handle: accountManager.handleCutOff}
	cutoffPairAllWatcher := &eventemitter.Watcher{Concurrent: false, Handle: accountManager.handleCutOffPair}
	unsupportedContractAllWatcher := &eventemitter.Watcher{Concurrent: false, Handle: accountManager.handleUnsupportedContract}

	eventemitter.On(eventemitter.Transfer, transferWatcher)
	eventemitter.On(eventemitter.Approve, approveWatcher)
	eventemitter.On(eventemitter.EthTransfer, ethTransferWatcher)
	eventemitter.On(eventemitter.Block_End, blockEndWatcher)
	eventemitter.On(eventemitter.Block_New, blockNewWatcher)
	eventemitter.On(eventemitter.WethDeposit, wethDepositWatcher)
	eventemitter.On(eventemitter.WethWithdrawal, wethWithdrawalWatcher)
	eventemitter.On(eventemitter.ChainForkDetected, blockForkWatcher)

	eventemitter.On(eventemitter.CancelOrder, cancelOrderWather)
	eventemitter.On(eventemitter.CutoffAll, cutoffAllWatcher)
	eventemitter.On(eventemitter.CutoffPair, cutoffPairAllWatcher)
	eventemitter.On(eventemitter.UnsupportedContract, unsupportedContractAllWatcher)
}

func (a *AccountManager) handleTokenTransfer(input eventemitter.EventData) (err error) {
	event := input.(*types.TransferEvent)
	log.Debugf("transfer,owner:%s,to:%s, token:%s", event.Receiver.Hex(), event.To.Hex(), event.Protocol.Hex())
	if event == nil || event.Status != types.TX_STATUS_SUCCESS {
		log.Info("received wrong status event, drop it")
		return nil
	}
	//log.Debugf("received transfer event, from:%s, sender:%s, receiver:%s, protocol:%s", event.From.Hex(), event.Sender.Hex(), event.Receiver.Hex(), event.Protocol.Hex())

	block := &ChangedOfBlock{}
	block.cachedDuration = new(big.Int).Set(a.cachedBlockCount)
	block.currentBlockNumber = new(big.Int).Set(event.BlockNumber)
	//balance
	block.saveBalanceKey(event.Sender, event.Protocol)
	block.saveBalanceKey(event.From, types.NilAddress)
	block.saveBalanceKey(event.Receiver, event.Protocol)

	//allowance
	if spender, err := loopringaccessor.GetSpenderAddress(event.To); nil == err {
		log.Debugf("handleTokenTransfer allowance owner:%s", event.Sender.Hex(), event.Protocol.Hex(), spender.Hex())
		block.saveAllowanceKey(event.Sender, event.Protocol, spender)
	}

	return nil
}

func (a *AccountManager) handleApprove(input eventemitter.EventData) error {
	event := input.(*types.ApprovalEvent)
	log.Debugf("approve,owner:%s, token:%s", event.Owner.Hex(), event.Protocol.Hex())
	if event == nil || event.Status != types.TX_STATUS_SUCCESS {
		log.Info("received wrong status event, drop it")
		return nil
	}

	block := &ChangedOfBlock{}
	block.cachedDuration = new(big.Int).Set(a.cachedBlockCount)
	block.currentBlockNumber = new(big.Int).Set(event.BlockNumber)

	block.saveAllowanceKey(event.Owner, event.Protocol, event.Spender)

	block.saveBalanceKey(event.Owner, types.NilAddress)

	return nil
}

func (a *AccountManager) handleWethDeposit(input eventemitter.EventData) (err error) {
	event := input.(*types.WethDepositEvent)
	log.Debugf("wethDeposit, owner:%s", event.Dst.Hex())

	if event == nil || event.Status != types.TX_STATUS_SUCCESS {
		log.Info("received wrong status event, drop it")
		return nil
	}
	block := &ChangedOfBlock{}
	block.cachedDuration = new(big.Int).Set(a.cachedBlockCount)
	block.currentBlockNumber = new(big.Int).Set(event.BlockNumber)
	block.saveBalanceKey(event.Dst, event.Protocol)
	block.saveBalanceKey(event.From, types.NilAddress)
	return
}

func (a *AccountManager) handleWethWithdrawal(input eventemitter.EventData) (err error) {
	event := input.(*types.WethWithdrawalEvent)
	log.Debugf("wethWithdrawal owner:%s", event.Src.Hex())
	if event == nil || event.Status != types.TX_STATUS_SUCCESS {
		log.Info("received wrong status event, drop it")
		return nil
	}

	block := &ChangedOfBlock{}
	block.cachedDuration = new(big.Int).Set(a.cachedBlockCount)
	block.currentBlockNumber = new(big.Int).Set(event.BlockNumber)

	block.saveBalanceKey(event.Src, event.Protocol)
	block.saveBalanceKey(event.From, types.NilAddress)

	return
}

func (a *AccountManager) handleBlockEnd(input eventemitter.EventData) error {
	event := input.(*types.BlockEvent)
	log.Debugf("handleBlockEndhandleBlockEndhandleBlockEnd:%s", event.BlockNumber.String())

	block := &ChangedOfBlock{}
	block.cachedDuration = new(big.Int).Set(a.cachedBlockCount)
	block.currentBlockNumber = new(big.Int).Set(event.BlockNumber)
	changedAddrs, _ := block.syncAndSaveBalances(a.tokenCacheDuration, a.ethCacheDuration)
	changedAllowanceAddrs, _ := block.syncAndSaveAllowances()

	removeExpiredBlock(block.currentBlockNumber, block.cachedDuration)

	for addr, _ := range changedAllowanceAddrs {
		changedAddrs[addr] = true
	}
	for addr, _ := range changedAddrs {
		event := &types.BalanceUpdateEvent{}
		event.Owner = addr.Hex()
		util.NotifyAccountBalanceUpdate(event)
	}

	// send blockEnd, miner use only
	if err := sendBlockEndKafkaMsg(event); nil != err {
		log.Errorf("err:%s", err.Error())
	}

	return nil
}

func (a *AccountManager) handleBlockNew(input eventemitter.EventData) error {
	event := input.(*types.BlockEvent)
	log.Debugf("handleBlockNewhandleBlockNewhandleBlockNewhandleBlockNew:%s", event.BlockNumber.String())
	block := &ChangedOfBlock{}
	block.cachedDuration = new(big.Int).Set(a.cachedBlockCount)
	block.currentBlockNumber = new(big.Int).Set(event.BlockNumber)
	block.currentBlockNumber = new(big.Int).Set(event.BlockNumber)
	return nil
}

func (a *AccountManager) handleEthTransfer(input eventemitter.EventData) error {
	event := input.(*types.EthTransferEvent)
	if event == nil || event.Status != types.TX_STATUS_SUCCESS {
		log.Info("received wrong status event, drop it")
		return nil
	}
	log.Debugf("transfer owner:%s, to:%s", event.From.Hex(), event.To.Hex())
	block := &ChangedOfBlock{}
	block.cachedDuration = new(big.Int).Set(a.cachedBlockCount)
	block.currentBlockNumber = new(big.Int).Set(event.BlockNumber)
	block.saveBalanceKey(event.From, types.NilAddress)
	block.saveBalanceKey(event.To, types.NilAddress)
	return nil
}

func (a *AccountManager) handleCancelOrder(input eventemitter.EventData) error {
	event := input.(*types.OrderCancelledEvent)
	if event == nil || event.Status != types.TX_STATUS_SUCCESS {
		log.Info("received wrong status event, drop it")
		return nil
	}
	block := &ChangedOfBlock{}
	block.cachedDuration = new(big.Int).Set(a.cachedBlockCount)
	block.currentBlockNumber = new(big.Int).Set(event.BlockNumber)
	block.saveBalanceKey(event.From, types.NilAddress)
	return nil
}

func (a *AccountManager) handleCutOff(input eventemitter.EventData) error {
	event := input.(*types.CutoffEvent)
	if event == nil || event.Status != types.TX_STATUS_SUCCESS {
		log.Info("received wrong status event, drop it")
		return nil
	}
	block := &ChangedOfBlock{}
	block.cachedDuration = new(big.Int).Set(a.cachedBlockCount)
	block.currentBlockNumber = new(big.Int).Set(event.BlockNumber)
	block.saveBalanceKey(event.From, types.NilAddress)
	return nil
}

func (a *AccountManager) handleCutOffPair(input eventemitter.EventData) error {
	event := input.(*types.CutoffPairEvent)
	if event == nil || event.Status != types.TX_STATUS_SUCCESS {
		log.Info("received wrong status event, drop it")
		return nil
	}
	block := &ChangedOfBlock{}
	block.cachedDuration = new(big.Int).Set(a.cachedBlockCount)
	block.currentBlockNumber = new(big.Int).Set(event.BlockNumber)
	block.saveBalanceKey(event.From, types.NilAddress)
	return nil
}

func (a *AccountManager) handleUnsupportedContract(input eventemitter.EventData) error {
	event := input.(*types.UnsupportedContractEvent)
	if event == nil || event.Status != types.TX_STATUS_SUCCESS {
		log.Info("received wrong status event, drop it")
		return nil
	}
	block := &ChangedOfBlock{}
	block.cachedDuration = new(big.Int).Set(a.cachedBlockCount)
	block.currentBlockNumber = new(big.Int).Set(event.BlockNumber)
	block.saveBalanceKey(event.From, types.NilAddress)
	return nil
}

func (a *AccountManager) UnlockedWallet(owner string) (err error) {
	if !common.IsHexAddress(owner) {
		return errors.New("owner isn't a valid hex-address")
	}

	//accountBalances := AccountBalances{}
	//accountBalances.Owner = common.HexToAddress(owner)
	//accountBalances.Balances = make(map[common.Address]Balance)
	//err = accountBalances.getOrSave(a.cacheDuration)
	rcache.Set(unlockCacheKey(common.HexToAddress(owner)), []byte("true"), -1)
	return
}

func (a *AccountManager) handleBlockFork(input eventemitter.EventData) (err error) {
	event := input.(*types.ForkedEvent)
	log.Infof("the eth network may be forked. flush all cache, detectedBlock:%s", event.DetectedBlock.String())

	i := new(big.Int).Set(event.DetectedBlock)
	changedAddrs := make(map[common.Address]bool)
	for i.Cmp(event.ForkBlock) >= 0 {
		changedOfBlock := &ChangedOfBlock{}
		changedOfBlock.currentBlockNumber = i
		changedBalanceAddrs, _ := changedOfBlock.syncAndSaveBalances(a.tokenCacheDuration, a.ethCacheDuration)
		changedAllowanceAddrs, _ := changedOfBlock.syncAndSaveAllowances()
		for addr, _ := range changedBalanceAddrs {
			changedAddrs[addr] = true
		}
		for addr, _ := range changedAllowanceAddrs {
			changedAddrs[addr] = true
		}
		i.Sub(i, big.NewInt(int64(1)))
	}

	for addr, _ := range changedAddrs {
		event := &types.BalanceUpdateEvent{}
		event.Owner = addr.Hex()
		util.NotifyAccountBalanceUpdate(event)
	}
	return nil
}
