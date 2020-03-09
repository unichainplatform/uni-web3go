package rpc

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/unichainplatform/unichain/common"
	"github.com/unichainplatform/unichain/types"
)

func GetTransactionByTxHash(hash common.Hash) (*types.RPCTransaction, error) {
	txInfo := &types.RPCTransaction{}
	err := ClientCall("uni_getTransactionByHash", txInfo, hash.Hex())
	if err != nil {
		fmt.Println("uni_getTransactionByHash:" + hash.Hex() + ", err=" + err.Error())
	}
	return txInfo, err
}

func GetReceiptByTxHash(hash common.Hash) (*types.RPCReceipt, error) {
	receipt := &types.RPCReceipt{}
	err := ClientCall("uni_getTransactionReceipt", receipt, hash.Hex())
	if err != nil {
		fmt.Println("uni_getTransactionReceipt:" + hash.Hex() + ", err=" + err.Error())
	}
	return receipt, err
}

func GetAllTxInPendingAndQueue() (map[string]map[string]map[string]string, error) {
	allNotExecutedTxs := map[string]map[string]map[string]string{}
	err := ClientCall("txpool_inspect", &allNotExecutedTxs)

	return allNotExecutedTxs, err
}

func GetTxpoolStatus() (int, int) {
	result := map[string]int{}
	ClientCall("txpool_status", &result)
	return result["pending"], result["queue"]
}

func CheckTxInPendingOrQueue(accountName string, txHash common.Hash) (bool, bool, error) {
	allNotExecutedTxs, err := GetAllTxInPendingAndQueue()
	if err != nil {
		return false, false, errors.New("can't get txs in pending and queue")
	}
	pendingTxs, _ := allNotExecutedTxs["pending"]
	queuedTxs, _ := allNotExecutedTxs["queued"]

	bInPending := false
	bInQueued := false
	accountTxs, ok := pendingTxs[accountName]
	if ok {
		if _, ok = accountTxs[txHash.Hex()]; ok {
			bInPending = true
		}
	}

	accountTxs, ok = queuedTxs[accountName]
	if ok {
		if _, ok = accountTxs[txHash.Hex()]; ok {
			bInQueued = true
		}
	}
	return bInPending, bInQueued, nil
}