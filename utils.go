package rpc

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/unichainplatform/unichain/common"
	"github.com/unichainplatform/unichain/crypto"
	"github.com/unichainplatform/unichain/params"
	"github.com/unichainplatform/unichain/types"
	"github.com/unichainplatform/unichain/utils/rlp"
	"github.com/unichainplatform/unichain/rpc"
)

var (
	once           sync.Once
	clientInstance *rpc.Client
	clientInstanceMap map[string]*rpc.Client
	hostIp         = "127.0.0.1"
	port           = 8545 //
	gasPrice       = new(big.Int).Mul(big.NewInt(100000), big.NewInt(100000))
	SentTxNum 	   = int32(0)
	GasLimit       = uint64(1000000)
)

func GenerateStatInfo() {
	atomic.AddInt32(&SentTxNum, 1)
}

type GenAction struct {
	*types.Action
	PrivateKey *ecdsa.PrivateKey
}

// DefultURL default rpc url
func DefultURL() string {
	return fmt.Sprintf("http://%s:%d", hostIp, port)
}

func SetUrl(_hostIp string, _port int) {
	hostIp = _hostIp
	port = _port
}

func setGasPrice(newGasPrice *big.Int) {
	gasPrice = newGasPrice
}

func setGasLimit(newGasLimit uint64) {
	GasLimit = newGasLimit
}

func GeneratePubKey() (common.PubKey, *ecdsa.PrivateKey) {
	prikey, _ := crypto.GenerateKey()
	return common.BytesToPubKey(crypto.FromECDSAPub(&prikey.PublicKey)), prikey
}

func NewGeAction(at types.ActionType, from, to common.Name, nonce uint64, assetid uint64, gaslimit uint64, amount *big.Int, payload []byte, remark []byte, prikey *ecdsa.PrivateKey) *GenAction {
	action := types.NewAction(at, from, to, nonce, assetid, gaslimit, amount, payload, remark)
	return &GenAction{
		Action:     action,
		PrivateKey: prikey,
	}
}
func SendTxTest(gcs []*GenAction) (common.Hash, error) {
	//nonce := GetNonce(sendaddr, "latest")
	signer := types.NewSigner(params.DefaultChainconfig.ChainID)
	var actions []*types.Action
	for _, v := range gcs {
		actions = append(actions, v.Action)
	}
	tx := types.NewTransaction(uint64(1), gasPrice, actions...)
	for _, v := range gcs {
		keys := make([]*types.KeyPair, 0)
		keys = append(keys, &types.KeyPair{v.PrivateKey, []uint64{uint64(0)}})
		err := types.SignActionWithMultiKey(v.Action, tx, signer, 0, keys)
		if err != nil {
			return common.Hash{}, err
		}
	}
	rawtx, _ := rlp.EncodeToBytes(tx)
	hash, err := SendRawTx(rawtx)
	GenerateStatInfo()
	return hash, err
}

//SendRawTx send raw transaction
func SendRawTx(rawTx []byte) (common.Hash, error) {
	hash := new(common.Hash)
	err := ClientCall("uni_sendRawTransaction", hash, hexutil.Bytes(rawTx))
	return *hash, err
}

// MustRPCClient Wraper rpc's client
func MustRPCClient() (*rpc.Client, error) {
	once.Do(func() {
		client, err := rpc.DialHTTP(DefultURL())
		if err != nil {
			return
		}
		clientInstance = client
	})

	return clientInstance, nil
}

func MustRPCClientWithAddr(nodeIp string, nodePort int64) (*rpc.Client, error) {
	endPoint := fmt.Sprintf("http://%s:%d", nodeIp, nodePort)
	if clientInstanceMap[endPoint] != nil {
		return clientInstanceMap[endPoint], nil
	}
	client, err := rpc.DialHTTP(endPoint)
	if err != nil {
		return nil, err
	}
	clientInstanceMap[endPoint] = client
	return client, nil
}

// ClientCall Wrapper rpc call api.
func ClientCall(method string, result interface{}, args ...interface{}) error {
	client, err := MustRPCClient()
	if err != nil {
		return err
	}
	err = client.CallContext(context.Background(), result, method, args...)
	if err != nil {
		return err
	}
	return nil
}

// ClientCall Wrapper rpc call api.
func ClientCallWithAddr(nodeIp string, nodePort int64, method string, result interface{}, args ...interface{}) error {
	client, err := MustRPCClientWithAddr(nodeIp, nodePort)
	if err != nil {
		return err
	}
	err = client.CallContext(context.Background(), result, method, args...)
	if err != nil {
		return err
	}
	return nil
}

// GasPrice suggest gas price
func GasPrice() (*big.Int, error) {
	gp := big.NewInt(0)
	err := ClientCall("uni_gasPrice", gp)
	return gp, err
}

// GetNonce get nonce by address and block number.
func GetNonce(accountName common.Name) (uint64, error) {
	nonce := uint64(0)
	err := ClientCall("account_getNonce", &nonce, accountName.String())
	if err != nil {
		return 0, err
	}
	return nonce, nil
}

