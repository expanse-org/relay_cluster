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

package gateway

import (
	"fmt"
	"github.com/expanse-org/relay_cluster/dao"
	"github.com/expanse-org/relay_cluster/test"
	"github.com/expanse-org/relay-lib/eth/accessor"
	"github.com/expanse-org/relay-lib/eth/contract"
	"github.com/expanse-org/relay-lib/eth/loopringaccessor"
	ethtyp "github.com/expanse-org/relay-lib/eth/types"
	"github.com/expanse-org/relay-lib/kafka"
	util "github.com/expanse-org/relay-lib/marketutil"
	"github.com/expanse-org/relay-lib/types"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
	"testing"
	"time"
)

var (
	miner            = test.Entity().Creator
	account1         = test.Entity().Accounts[0].Address
	account2         = test.Entity().Accounts[1].Address
	lrcTokenAddress  = util.AllTokens["LRC"].Protocol
	wexpTokenAddress = util.AllTokens["WEXP"].Protocol
	delegateAddress  = test.Delegate()
	gas              = big.NewInt(200000)
	gasPrice         = big.NewInt(21000000000)
)

func TestEthNodeAccessor_SetBalance(t *testing.T) {
	test.SetTokenBalances()
}

func TestEthNodeAccessor_WexpDeposit(t *testing.T) {
	account := account1
	wexpAddr := wexpTokenAddress
	amount := new(big.Int).Mul(big.NewInt(1e18), big.NewInt(2))
	callMethod := accessor.ContractSendTransactionMethod("latest", test.WexpAbi(), wexpAddr)

	hash, _, err := callMethod(account, contract.METHOD_WEXP_DEPOSIT, gas, gasPrice, amount)
	if nil != err {
		t.Fatalf("call method wexp-deposit error:%s", err.Error())
	}

	if err := sendPendingTx(hash); err != nil {
		t.Fatalf("send tx err:%s", err.Error())
	}
	t.Logf("wexp-deposit result:%s", hash)
}

func TestEthNodeAccessor_WexpWithdrawal(t *testing.T) {
	account := account1
	wexpAddr := wexpTokenAddress
	amount := new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1))
	callMethod := accessor.ContractSendTransactionMethod("latest", test.WexpAbi(), wexpAddr)

	hash, _, err := callMethod(account, "withdraw", gas, gasPrice, nil, amount)
	if nil != err {
		t.Fatalf("call method wexp-withdraw error:%s", err.Error())
	}

	if err := sendPendingTx(hash); err != nil {
		t.Fatalf("send tx err:%s", err.Error())
	}

	t.Logf("wexp-withdraw result:%s", hash)
}

func TestEthNodeAccessor_WexpTransfer(t *testing.T) {
	from := account1
	to := account2
	wexpAddr := wexpTokenAddress
	amount := new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1))

	callMethod := accessor.ContractSendTransactionMethod("latest", test.WexpAbi(), wexpAddr)
	hash, _, err := callMethod(from, "transfer", gas, gasPrice, nil, to, amount)
	if nil != err {
		t.Fatalf("call method wexp-transfer error:%s", err.Error())
	}
	if err := sendPendingTx(hash); err != nil {
		t.Fatalf("send tx err:%s", err.Error())
	}
	t.Logf("wexp-transfer result:%s", hash)
}

func TestEthNodeAccessor_EthTransfer(t *testing.T) {
	sender := account1
	receiver := account2
	amount := new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1))

	hash, _, err := accessor.SignAndSendTransaction(sender, receiver, gas, gasPrice, amount, []byte("test"), false)
	if err != nil {
		t.Errorf(err.Error())
	}

	if err := sendPendingTx(hash); err != nil {
		t.Fatalf("send tx err:%s", err.Error())
	}
	t.Logf("eth transfer txhash:%s", hash)
}

func TestEthNodeAccessor_EthBalance(t *testing.T) {
	account := account1

	var balance types.Big
	if err := accessor.GetBalance(&balance, common.HexToAddress("0x8311804426A24495bD4306DAf5f595A443a52E32"), "0x53ca90"); err != nil {
		t.Fatalf(err.Error())
	} else {
		amount := new(big.Rat).SetFrac(balance.BigInt(), big.NewInt(1e18)).FloatString(2)
		t.Logf("eth account:%s amount:%s", account.Hex(), amount)
	}

	//time.Sleep(5 * time.Second)
}

func TestEthNodeAccessor_ERC20Transfer(t *testing.T) {
	from := account1
	to := account2
	amount := new(big.Int).Mul(big.NewInt(1e18), big.NewInt(150))

	erc20abi := loopringaccessor.Erc20Abi()
	tokenAddress := lrcTokenAddress
	callMethod := accessor.ContractSendTransactionMethod("latest", erc20abi, tokenAddress)

	hash, _, err := callMethod(from, "transfer", gas, gasPrice, nil, to, amount)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if err := sendPendingTx(hash); err != nil {
		t.Fatalf("send tx err:%s", err.Error())
	}
	t.Logf("txhash:%s", hash)
}

func TestEthNodeAccessor_ERC20Balance(t *testing.T) {
	accounts := []common.Address{common.HexToAddress("0x750ad4351bb728cec7d639a9511f9d6488f1e259"), common.HexToAddress("0x251f3bd45b06a8b29cb6d171131e192c1254fec1")}
	tokens := []common.Address{common.HexToAddress("0x639687b7f8501f174356d3acb1972f749021ccd0"), common.HexToAddress("0xe1C541BA900cbf212Bc830a5aaF88aB499931751")}

	for _, tokenAddress := range tokens {
		for _, account := range accounts {
			balance, err := loopringaccessor.Erc20Balance(tokenAddress, account, "latest")
			if err != nil {
				t.Fatalf("accessor get erc20 balance error:%s", err.Error())
			}
			//amount := new(big.Rat).SetFrac(balance, big.NewInt(1e18)).FloatString(2)
			amount := balance.String()
			symbol, _ := util.GetSymbolWithAddress(tokenAddress)
			t.Logf("token:%s account:%s amount:%s", symbol, account.Hex(), amount)
		}
	}
}

func TestEthNodeAccessor_Approval(t *testing.T) {
	accounts := []common.Address{account1, account2}
	spender := delegateAddress
	tokenAddress := lrcTokenAddress
	amount := new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1000000))
	//amount,_ := big.NewInt(0).SetString("9223372036854775806000000000000000000", 0)

	for _, account := range accounts {
		callMethod := accessor.ContractSendTransactionMethod("latest", loopringaccessor.Erc20Abi(), tokenAddress)
		if hash, _, err := callMethod(account, "approve", gas, gasPrice, nil, spender, amount); nil != err {
			t.Fatalf("call method approve error:%s", err.Error())
		} else {
			sendPendingTx(hash)
			t.Logf("approve result:%s", hash)
		}
	}
}

func TestEthNodeAccessor_Allowance(t *testing.T) {
	/*accounts := []common.Address{account1, account2}
	tokens := []common.Address{lrcTokenAddress, wexpTokenAddress}*/
	accounts := []common.Address{common.HexToAddress("0x750ad4351bb728cec7d639a9511f9d6488f1e259"), common.HexToAddress("0x251f3bd45b06a8b29cb6d171131e192c1254fec1")}
	tokens := []common.Address{common.HexToAddress("0x639687b7f8501f174356d3acb1972f749021ccd0"), common.HexToAddress("0xe1C541BA900cbf212Bc830a5aaF88aB499931751")}

	spender := delegateAddress

	for _, tokenAddress := range tokens {
		for _, account := range accounts {
			if allowance, err := loopringaccessor.Erc20Allowance(tokenAddress, account, spender, "latest"); err != nil {
				t.Fatalf("accessor get erc20 approval error:%s", err.Error())
			} else {
				amount := new(big.Rat).SetFrac(allowance, big.NewInt(1e18)).FloatString(2)
				symbol, _ := util.GetSymbolWithAddress(tokenAddress)
				t.Logf("token:%s, account:%s, amount:%s", symbol, account.Hex(), amount)
			}
		}
	}
}

func TestEthNodeAccessor_CancelOrder(t *testing.T) {
	var (
		model        *dao.Order
		state        types.OrderState
		err          error
		result       string
		orderhash    = common.HexToHash("0x2fd51638ad98d79aef0e1aecd5d85c5a471914792b9aaf74afdb1fafd5c25ff2")
		cancelAmount = new(big.Int).Mul(big.NewInt(1e18), big.NewInt(2))
	)

	// get order
	rds := test.Rds()
	if model, err = rds.GetOrderByHash(orderhash); err != nil {
		t.Fatalf(err.Error())
	}
	if err := model.ConvertUp(&state); err != nil {
		t.Fatalf(err.Error())
	}

	account := accounts.Account{Address: state.RawOrder.Owner}

	// create cancel order contract function parameters
	addresses := [5]common.Address{state.RawOrder.Owner, state.RawOrder.TokenS, state.RawOrder.TokenB, state.RawOrder.WalletAddress, state.RawOrder.AuthAddr}
	values := [6]*big.Int{state.RawOrder.AmountS, state.RawOrder.AmountB, state.RawOrder.ValidSince, state.RawOrder.ValidUntil, state.RawOrder.LrcFee, cancelAmount}
	buyNoMoreThanB := state.RawOrder.BuyNoMoreThanAmountB
	marginSplitPercentage := state.RawOrder.MarginSplitPercentage
	v := state.RawOrder.V
	s := state.RawOrder.S
	r := state.RawOrder.R

	// call cancel order
	protocol := test.Protocol()
	implAddress := loopringaccessor.ProtocolAddresses()[protocol].ContractAddress
	callMethod := accessor.ContractSendTransactionMethod("latest", loopringaccessor.ProtocolImplAbi(), implAddress)
	if result, _, err = callMethod(account.Address, "cancelOrder", gas, gasPrice, nil, addresses, values, buyNoMoreThanB, marginSplitPercentage, v, r, s); nil != err {
		t.Fatalf("call method cancelOrder error:%s", err.Error())
	}

	if err := sendPendingTx(result); err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf("cancelOrder result:%s", result)
}

func TestEthNodeAccessor_GetCancelledOrFilled(t *testing.T) {
	orderhash := common.HexToHash("0x77aecf96a71d260074ab6ad9352365f0e83cd87d4d3a424071e76cafc393f549")

	protocol := test.Protocol()
	implAddress := loopringaccessor.ProtocolAddresses()[protocol].ContractAddress
	if amount, err := loopringaccessor.GetCancelledOrFilled(implAddress, orderhash, "latest"); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("cancelOrFilled amount:%s", amount.String())
	}
}

// cutoff的值必须在两个块的timestamp之间
func TestEthNodeAccessor_CutoffAll(t *testing.T) {
	account := common.HexToAddress("0x1B978a1D302335a6F2Ebe4B8823B5E17c3C84135")
	cutoff := big.NewInt(1531808145)

	callMethod := accessor.ContractSendTransactionMethod("latest", test.LprAbi(), test.Protocol())
	result, _, err := callMethod(account, contract.METHOD_CUTOFF_ALL, gas, gasPrice, nil, cutoff)
	if nil != err {
		t.Fatalf("call method cancelAllOrders error:%s", err.Error())
	}
	if err := sendPendingTx(result); err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf("cutoff result:%s", result)
}

func TestEthNodeAccessor_GetCutoffAll(t *testing.T) {
	owner := account1

	var res types.Big
	if err := loopringaccessor.GetCutoff(&res, test.Protocol(), owner, "latest"); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("cutoff timestamp:%s", res.BigInt().String())
	}
}

func TestEthNodeAccessor_CutoffPair(t *testing.T) {
	account := common.HexToAddress("0xb1018949b241D76A1AB2094f473E9bEfeAbB5Ead")
	cutoff := big.NewInt(1531107175)
	token1 := lrcTokenAddress
	token2 := wexpTokenAddress

	callMethod := accessor.ContractSendTransactionMethod("latest", test.LprAbi(), test.Protocol())
	result, _, err := callMethod(account, contract.METHOD_CUTOFF_PAIR, gas, gasPrice, nil, token1, token2, cutoff)
	if nil != err {
		t.Fatalf("call method cancelAllOrdersByTradingPair error:%s", err.Error())
	}
	if err := sendPendingTx(result); err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf("cutoff result:%s", result)
}

func TestEthNodeAccessor_GetCutoffPair(t *testing.T) {
	owner := accounts.Account{Address: account2}
	token1 := lrcTokenAddress
	token2 := wexpTokenAddress

	var res types.Big
	if err := loopringaccessor.GetCutoffPair(&res, test.Protocol(), owner.Address, token1, token2, "latest"); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("cutoffpair timestamp:%s", res.BigInt().String())
	}
}

func TestEthNodeAccessor_TokenRegister(t *testing.T) {
	account := accounts.Account{Address: test.Entity().Creator.Address}

	address := lrcTokenAddress
	symbol := "LRC"

	callMethod := accessor.ContractSendTransactionMethod("latest", test.TokenRegisterAbi(), test.TokenRegisterAddress())
	if result, _, err := callMethod(account.Address, contract.METHOD_TOKEN_REGISTRY, gas, gasPrice, nil, address, symbol); nil != err {
		t.Fatalf("call method registerToken error:%s", err.Error())
	} else {
		t.Logf("registerToken result:%s", result)
	}
}

func TestEthNodeAccessor_TokenUnRegister(t *testing.T) {
	account := accounts.Account{Address: test.Entity().Creator.Address}

	address := lrcTokenAddress
	symbol := "LRC"

	callMethod := accessor.ContractSendTransactionMethod("latest", test.TokenRegisterAbi(), test.TokenRegisterAddress())
	if result, _, err := callMethod(account.Address, contract.METHOD_TOKEN_UNREGISTRY, gas, gasPrice, nil, address, symbol); nil != err {
		t.Fatalf("call method unregisterToken error:%s", err.Error())
	} else {
		t.Logf("unregisterToken result:%s", result)
	}
}

func TestEthNodeAccessor_IsTokenRegistried(t *testing.T) {
	var result string

	symbol := "LRC"

	callMethod := accessor.ContractCallMethod(test.TokenRegisterAbi(), test.TokenRegisterAddress())
	if err := callMethod(&result, "isTokenRegisteredBySymbol", "latest", symbol); nil != err {
		t.Fatalf("call method isTokenRegisteredBySymbol error:%s", err.Error())
	} else {
		t.Logf("isTokenRegisteredBySymbol result:%s", result)
	}
}

func TestEthNodeAccessor_GetAddressBySymbol(t *testing.T) {
	var result string
	symbol := "LRC"

	callMethod := accessor.ContractCallMethod(test.TokenRegisterAbi(), test.TokenRegisterAddress())
	if err := callMethod(&result, "getAddressBySymbol", "latest", symbol); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("symbol map:%s->%s", symbol, common.HexToAddress(result).Hex())
	}
}

// 注册合约
func TestEthNodeAccessor_AuthorizedAddress(t *testing.T) {
	account := accounts.Account{Address: test.Entity().Creator.Address}

	callMethod := accessor.ContractSendTransactionMethod("latest", test.DelegateAbi(), test.DelegateAddress())
	if result, _, err := callMethod(account.Address, contract.METHOD_ADDRESS_AUTHORIZED, gas, gasPrice, nil, test.Protocol()); nil != err {
		t.Fatalf("call method authorizeAddress error:%s", err.Error())
	} else {
		t.Logf("authorizeAddress result:%s", result)
	}
}

func TestEthNodeAccessor_DeAuthorizedAddress(t *testing.T) {
	account := accounts.Account{Address: test.Entity().Creator.Address}

	protocol := test.Protocol()
	delegateAddress := loopringaccessor.ProtocolAddresses()[protocol].DelegateAddress
	delegateAbi := loopringaccessor.DelegateAbi()
	callMethod := accessor.ContractSendTransactionMethod("latest", delegateAbi, delegateAddress)
	if result, _, err := callMethod(account.Address, contract.METHOD_ADDRESS_DEAUTHORIZED, gas, gasPrice, nil, protocol); nil != err {
		t.Fatalf("call method deauthorizeAddress error:%s", err.Error())
	} else {
		t.Logf("deauthorizeAddress result:%s", result)
	}
}

func TestEthNodeAccessor_IsAddressAuthorized(t *testing.T) {
	var result string

	callMethod := accessor.ContractCallMethod(test.DelegateAbi(), test.DelegateAddress())
	if err := callMethod(&result, "isAddressAuthorized", "latest", test.DelegateAddress()); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("isAddressAuthorized result:%s", result)
	}
}

func TestEthNodeAccessor_TokenAddress(t *testing.T) {
	symbol := "LRC"
	protocol := test.Protocol()
	tokenRegistryAddress := loopringaccessor.ProtocolAddresses()[protocol].TokenRegistryAddress
	tokenRegisterAbi := loopringaccessor.TokenRegistryAbi()
	callMethod := accessor.ContractCallMethod(tokenRegisterAbi, tokenRegistryAddress)
	var result string
	if err := callMethod(&result, "getAddressBySymbol", "latest", symbol); nil != err {
		t.Fatalf("call method tokenAddress error:%s", err.Error())
	} else {
		t.Logf("symbol:%s-> address:%s", symbol, result)
	}
}

func TestEthNodeAccessor_BlockTransactionStatus(t *testing.T) {
	const (
		startBlock = 5444365
		endBlock   = startBlock + 2
	)

	for i := startBlock; i < endBlock; i++ {
		blockNumber := big.NewInt(int64(i))

	blockMark:
		var blockWithTxHash ethtyp.BlockWithTxHash
		if err := accessor.GetBlockByNumber(&blockWithTxHash, blockNumber, false); err != nil {
			time.Sleep(1 * time.Second)
			fmt.Printf("........err:%s nil\n", err.Error())
			goto blockMark
		} else if len(blockWithTxHash.Transactions) == 0 {
			time.Sleep(1 * time.Second)
			fmt.Printf("........tx 0\n")
			goto blockMark
		}

		blockWithTxAndReceipt := &ethtyp.BlockWithTxAndReceipt{}
		blockWithTxAndReceipt.Block = blockWithTxHash.Block
		blockWithTxAndReceipt.Transactions = []ethtyp.Transaction{}
		blockWithTxAndReceipt.Receipts = []ethtyp.TransactionReceipt{}

		txno := len(blockWithTxHash.Transactions)
		var rcReqs = make([]*accessor.BatchTransactionRecipientReq, txno)
		for idx, txstr := range blockWithTxHash.Transactions {
			var (
				rcreq accessor.BatchTransactionRecipientReq
				rc    ethtyp.TransactionReceipt
				rcerr error
			)

			rcreq.TxHash = txstr
			rcreq.TxContent = rc
			rcreq.Err = rcerr

			rcReqs[idx] = &rcreq
		}

		if err := accessor.BatchTransactionRecipients(rcReqs, blockWithTxAndReceipt.Number.BigInt().String()); err != nil {
			t.Fatalf(err.Error())
		}

		for idx, _ := range rcReqs {
			blockWithTxAndReceipt.Receipts = append(blockWithTxAndReceipt.Receipts, rcReqs[idx].TxContent)
		}

		var (
			success = 0
			failed  = 0
			invalid = 0
		)
		for _, v := range blockWithTxAndReceipt.Receipts {
			if v.StatusInvalid() {
				invalid++
				fmt.Printf("tx:%s status is nil\n", v.TransactionHash)
			} else if v.Status.BigInt().Cmp(big.NewInt(1)) < 0 {
				failed++
				fmt.Printf("tx:%s status:%s\n", v.TransactionHash, v.Status.BigInt().String())
			} else {
				success++
			}
		}
		fmt.Printf("blockNumber:%s, blockHash:%s, txNumber:%d, successTx:%d failed:%d nil:%d \n",
			blockNumber.String(), blockWithTxHash.Hash.Hex(), txno, success, failed, invalid)
	}
}

func TestEthNodeAccessor_GetTransaction(t *testing.T) {
	tx := &ethtyp.Transaction{}
	if err := accessor.GetTransactionByHash(tx, "0x26383249d29e13c4c5f73505775813829875d0b0bf496f2af2867548e2bf8108", "pending"); err == nil {
		t.Logf("tx blockNumber:%s, from:%s, to:%s, gas:%s value:%s", tx.BlockNumber.BigInt().String(), tx.From, tx.To, tx.Gas.BigInt().String(), tx.Value.BigInt().String())
		t.Logf("tx input:%s", tx.Input)
	} else {
		t.Fatalf(err.Error())
	}
}

func TestEthNodeAccessor_GetTransactionReceipt(t *testing.T) {
	var tx ethtyp.TransactionReceipt
	if err := accessor.GetTransactionReceipt(&tx, "0x26383249d29e13c4c5f73505775813829875d0b0bf496f2af2867548e2bf8108", "latest"); err == nil {
		t.Logf("tx blockNumber:%s gasUsed:%s status:%s logs:%d", tx.BlockNumber.BigInt().String(), tx.GasUsed.BigInt().String(), tx.Status.BigInt().String(), len(tx.Logs))
		idx := len(tx.Logs) - 1
		t.Logf("tx event:%d data:%s", idx, tx.Logs[idx].Data)
		for _, v := range tx.Logs[idx].Topics {
			t.Logf("topic:%s", v)
		}
	} else {
		t.Fatalf(err.Error())
	}
}

func TestEthNodeAccessor_GetBlock(t *testing.T) {
	hash := "0x25d526f4d913a563783fd09a1e5472c505d644fc2f3ac17eae8f2704943dd033"
	var block ethtyp.Block
	if err := accessor.GetBlockByHash(&block, hash, false); err != nil {
		t.Fatalf(err.Error())
	} else {
		t.Logf("number:%s, hash:%s, time:%s", block.Number.BigInt().String(), block.Hash.Hex(), block.Timestamp.BigInt().String())
	}
}

func TestEthNodeAccessor_GetTransactionCount(t *testing.T) {
	var count types.Big
	user := common.HexToAddress("0x71c079107b5af8619d54537a93dbf16e5aab4900")
	if err := accessor.GetTransactionCount(&count, user, "latest"); err != nil {
		t.Fatalf(err.Error())
	} else {
		t.Logf("transaction count:%d", count.Int64())
	}
}

func TestEthNodeAccessor_GetFullBlock(t *testing.T) {
	blockNumber := big.NewInt(5514801)
	withObject := true
	ret, err := accessor.GetFullBlock(blockNumber, withObject)
	if err != nil {
		t.Fatalf(err.Error())
	}
	block := ret.(*ethtyp.BlockWithTxAndReceipt)
	for _, v := range block.Transactions {
		t.Logf("hash:%s", v.Hash)
	}
	t.Logf("length of block:%s is %d", blockNumber.String(), len(block.Transactions))
}

// 使用rpc.client调用eth call时应该使用 arg参数应该指针 保证unmarshal的正确性
func TestEthNodeAccessor_Call(t *testing.T) {
	var (
		arg1 ethtyp.CallArg
		res1 string
	)

	arg1.To = common.HexToAddress("0x45245bc59219eeaAF6cD3f382e078A461FF9De7B")
	arg1.Data = "0x95d89b41"
	if err := accessor.Call(&res1, &arg1, "latest"); err != nil {
		t.Fatal(err)
	}

	t.Log(res1)

	type CallArg struct {
		From     common.Address `json:"from"`
		To       common.Address `json:"to"`
		Gas      string         `json:"gas"`
		GasPrice string         `json:"gasPrice"`
		Value    string         `json:"value"`
		Data     string         `json:"data"`
		Nonce    string         `json:"nonce"`
	}

	var (
		client *rpc.Client
		err    error
		arg2   CallArg
		res2   string
	)

	url := "http://ec2-13-115-183-194.ap-northeast-1.compute.amazonaws.com:8545"
	if client, err = rpc.Dial(url); nil != err {
		t.Fatalf("rpc.Dail err : %s, url:%s", err.Error(), url)
	}

	arg2.To = common.HexToAddress("0x45245bc59219eeaAF6cD3f382e078A461FF9De7B")
	arg2.Data = "0x95d89b41"
	if err := client.Call(&res2, "eth_call", arg2, "latest"); err != nil {
		t.Fatal(err)
	}

	t.Log(res2)
}

func TestAccessor_MutilClient(t *testing.T) {
	test.Delegate()

}

func sendPendingTx(txhash string) error {
	var tx ethtyp.Transaction
	if err := accessor.GetTransactionByHash(&tx, txhash, "latest"); err != nil {
		return err
	}
	test.Producer().SendMessage(kafka.Kafka_Topic_Extractor_PendingTransaction, &tx, "extractor")
	return nil
}
