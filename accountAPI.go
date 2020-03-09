package rpc

import (
	"crypto/ecdsa"
	"github.com/unichainplatform/unichain/accountmanager"
	"github.com/unichainplatform/unichain/common"
	"github.com/unichainplatform/unichain/types"
	"github.com/unichainplatform/unichain/asset"
	"math/big"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/unichainplatform/unichain/utils/rlp"
)

func CreateAccount(fromAccount string, fromPriKey *ecdsa.PrivateKey, newAccountName string, pubKey common.PubKey, amount *big.Int) (common.Hash, common.Name, error) {
	accountName := common.Name(fromAccount)
	nonce, err := GetNonce(accountName)
	if err != nil {
		return common.Hash{}, "", err
	}
	createdAccountName := common.Name(newAccountName) //common.Name(newAccountName + strconv.FormatInt(int64(nonce),10))
	gc := NewGeAction(types.CreateAccount, accountName, createdAccountName, nonce, 1, GasLimit, amount, pubKey[:], nil, fromPriKey)
	var gcs []*GenAction
	gcs = append(gcs, gc)
	txHash, err := SendTxTest(gcs)
	return txHash, createdAccountName, err
}

func Transfer(from, to string, assetId uint64, amount *big.Int, prikey *ecdsa.PrivateKey) (common.Hash, error) {
	nonce, err := GetNonce(common.Name(from))
	if err != nil {
		return common.Hash{}, err
	}
	gc := NewGeAction(types.Transfer, common.Name(from), common.Name(to), nonce, assetId, GasLimit, amount, nil, nil, prikey)
	var gcs []*GenAction
	gcs = append(gcs, gc)
	return SendTxTest(gcs)
}

func IsAccountExist(accountName string) (bool, error) {
	isExist := new(bool)
	err := ClientCall("account_accountIsExist", isExist, accountName)
	if err != nil {
		return false, err
	}
	return *isExist, nil
}

func GetAccountByName(accountName string) (*accountmanager.Account, error) {
	account := &accountmanager.Account{}
	err := ClientCall("account_getAccountByName", account, accountName)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func GetAccountExByName(accountName string) (*accountmanager.Account, error) {
	account := &accountmanager.Account{}
	err := ClientCall("account_getAccountExByName", account, accountName)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func GetAccountById(accountId uint64) (*accountmanager.Account, error) {
	account := &accountmanager.Account{}
	err := ClientCall("account_getAccountById", account, accountId)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func GetAccountExById(accountId uint64) (*accountmanager.Account, error) {
	account := &accountmanager.Account{}
	err := ClientCall("account_getAccountExById", account, accountId)
	if err != nil {
		return nil, err
	}
	return account, nil
}

// GetAccountBalanceByID get balance by address ,assetID and number.
func GetAccountBalanceByID(accountName string, assetID uint64) (*big.Int, error) {
	balance := big.NewInt(0)
	err := ClientCall("account_getAccountBalanceByID", balance, common.Name(accountName), assetID)
	return balance, err
}

// GetAccountBalanceByID get balance by address ,assetID and number.
func GetCode(accountName string) (hexutil.Bytes, error) {
	balance := hexutil.Bytes{}
	err := ClientCall("account_getAccountBalanceByID", balance, common.Name(accountName))
	return balance, err
}

func GetAssetInfoByName(assetName string) (*asset.AssetObject, error) {
	assetInfo := &asset.AssetObject{}
	err := ClientCall("account_getAssetInfoByName", assetInfo, assetName)
	return assetInfo, err
}


func GetAssetInfoById(assetID uint64) (*asset.AssetObject, error) {
	assetInfo := &asset.AssetObject{}
	err := ClientCall("account_getAssetInfoByID", assetInfo, assetID)
	return assetInfo, err
}

func IssueAsset(from, owner common.Name, amount *big.Int, assetName string, symbol string, nonce uint64, decimals uint64, prikey *ecdsa.PrivateKey) (common.Hash, error) {
	assetObj := &asset.AssetObject{
		AssetName: assetName,
		Symbol:    symbol,
		Amount:    amount,
		Decimals:  decimals,
		Owner:     owner,
	}
	payload, err := rlp.EncodeToBytes(assetObj)
	if err != nil {
		return common.Hash{}, err
	}
	gc := NewGeAction(types.IssueAsset, from, "", nonce, 1, GasLimit, nil, payload, nil, prikey)
	var gcs []*GenAction
	gcs = append(gcs, gc)
	return SendTxTest(gcs)
}

func IncreaseAsset(from common.Name, to common.Name, assetId uint64, increasedAmount *big.Int, nonce uint64, prikey *ecdsa.PrivateKey) (common.Hash, error) {
	ast := &asset.AssetObject{
		AssetID:   assetId,
		AssetName: "",
		Symbol:    "",
		Amount:    increasedAmount,
	}
	payload, err := rlp.EncodeToBytes(ast)
	if err != nil {
		return common.Hash{}, err
	}
	gc := NewGeAction(types.IncreaseAsset, from, to, nonce, assetId, GasLimit, nil, payload, nil, prikey)
	var gcs []*GenAction
	gcs = append(gcs, gc)
	return SendTxTest(gcs)
}


func SetAssetOwner(from, newOwner common.Name, assetId uint64, prikey *ecdsa.PrivateKey) (common.Hash, error) {
	ast := &asset.AssetObject{
		AssetID: assetId,
		Owner:   newOwner,
	}
	payload, err := rlp.EncodeToBytes(ast)
	if err != nil {
		return common.Hash{}, err
	}
	nonce, err := GetNonce(common.Name(from))
	if err != nil {
		return common.Hash{}, err
	}

	gc := NewGeAction(types.SetAssetOwner, from, "", nonce, assetId, GasLimit, nil, payload, nil, prikey)
	var gcs []*GenAction
	gcs = append(gcs, gc)
	return SendTxTest(gcs)
}