package rpc

import (
	"crypto/ecdsa"
	"github.com/unichainplatform/unichain/common"
	"github.com/unichainplatform/unichain/consensus/dpos"
	"github.com/unichainplatform/unichain/types"
	"github.com/unichainplatform/unichain/utils/rlp"
	"math/big"
)

func RegisterCandidate(fromAccount string, fromPriKey *ecdsa.PrivateKey, url string, stake *big.Int) (common.Hash, error) {
	rp := dpos.RegisterCandidate{
		Info:   url,
	}

	rawdata, err := rlp.EncodeToBytes(rp)
	if err != nil {
		return common.Hash{}, err
	}

	accountName := common.Name(fromAccount)
	nonce, err := GetNonce(accountName)
	if err != nil {
		return common.Hash{}, err
	}
	gc := NewGeAction(types.RegCandidate, accountName, accountName, nonce, 1, GasLimit, stake, rawdata, nil, fromPriKey)
	var gcs []*GenAction
	gcs = append(gcs, gc)
	txHash, err := SendTxTest(gcs)
	return txHash, err
}

func UpdateCandidate(fromAccount string, fromPriKey *ecdsa.PrivateKey, url string, stake *big.Int) (common.Hash, error) {
	rp := dpos.UpdateCandidate{
		Info:   url,
	}

	rawdata, err := rlp.EncodeToBytes(rp)
	if err != nil {
		return common.Hash{}, err
	}

	accountName := common.Name(fromAccount)
	nonce, err := GetNonce(accountName)
	if err != nil {
		return common.Hash{}, err
	}
	gc := NewGeAction(types.UpdateCandidate, accountName, "", nonce, 1, GasLimit, stake, rawdata, nil, fromPriKey)
	var gcs []*GenAction
	gcs = append(gcs, gc)
	txHash, err := SendTxTest(gcs)
	return txHash, err
}

func UnregCandidate(fromAccount string, fromPriKey *ecdsa.PrivateKey) (common.Hash, error) {
	accountName := common.Name(fromAccount)
	nonce, err := GetNonce(accountName)
	if err != nil {
		return common.Hash{}, err
	}
	gc := NewGeAction(types.UnregCandidate, accountName, "", nonce, 1, GasLimit, nil, nil, nil,  fromPriKey)
	var gcs []*GenAction
	gcs = append(gcs, gc)
	txHash, err := SendTxTest(gcs)
	return txHash, err
}

func RefundCandidate(fromAccount string, fromPriKey *ecdsa.PrivateKey) (common.Hash, error) {
	accountName := common.Name(fromAccount)
	nonce, err := GetNonce(accountName)
	if err != nil {
		return common.Hash{}, err
	}
	gc := NewGeAction(types.RefundCandidate, accountName, "", nonce, 1, GasLimit, nil, nil, nil,  fromPriKey)
	var gcs []*GenAction
	gcs = append(gcs, gc)
	txHash, err := SendTxTest(gcs)
	return txHash, err
}

func UpdateCandidatePubKey(fromAccount string, fromPriKey *ecdsa.PrivateKey, publicKey common.PubKey) (common.Hash, error) {
	rp := dpos.UpdateCandidatePubKey{
		PubKey:   publicKey,
	}

	rawdata, err := rlp.EncodeToBytes(rp)
	if err != nil {
		return common.Hash{}, err
	}

	accountName := common.Name(fromAccount)
	nonce, err := GetNonce(accountName)
	if err != nil {
		return common.Hash{}, err
	}
	gc := NewGeAction(types.UpdateCandidatePubKey, accountName, "", nonce, 1, GasLimit, nil, rawdata, nil, fromPriKey)
	var gcs []*GenAction
	gcs = append(gcs, gc)
	txHash, err := SendTxTest(gcs)
	return txHash, err
}

func VoteCandidate(fromAccount string, fromPriKey *ecdsa.PrivateKey, candidate string, stake *big.Int) (common.Hash, error) {
	arg := &dpos.VoteCandidate{
		Candidate: candidate,
		Stake:    stake,
	}
	payload, err := rlp.EncodeToBytes(arg)
	if err != nil {
		panic(err)
	}
	accountName := common.Name(fromAccount)
	nonce, err := GetNonce(accountName)
	if err != nil {
		return common.Hash{}, err
	}
	gc := NewGeAction(types.VoteCandidate, accountName, "", nonce, 1, GasLimit, nil, payload, nil, fromPriKey)
	var gcs []*GenAction
	gcs = append(gcs, gc)
	txHash, err := SendTxTest(gcs)
	return txHash, err
}

func GetCandidate(epoch int, accountName common.Name) (*dpos.CandidateInfo, error) {
	fields := &dpos.CandidateInfo{}
	err := ClientCall("dpos_candidate", fields, epoch, accountName.String())
	return fields, err
}

func GetCandidateSize(epoch int) (uint64, error) {
	fields := uint64(0)
	err := ClientCall("dpos_candidatesSize", &fields, epoch)
	return fields, err
}

func GetCandidates(epoch int, bDetailInfo bool) ([]*dpos.CandidateInfo, error) {
	fields := make([]*dpos.CandidateInfo, 0)
	err := ClientCall("dpos_candidates", fields, epoch, bDetailInfo)
	return fields, err
}

func GetActivedCandidateNumber(epoch int) (uint64, error) {
	fields := uint64(0)
	err := ClientCall("dpos_getActivedCandidateSize", &fields, epoch)
	return fields, err
}

func GetActivedCandidate(epoch int, index int) (map[string]interface{}, error) {
	fields := map[string]interface{}{}
	err := ClientCall("dpos_getActivedCandidate", &fields, epoch, index)
	return fields, err
}

func GetDposInfo() (*dpos.Config, error) {
	fields := &dpos.Config{}
	err := ClientCall("dpos_info", &fields)
	return fields, err
}

func GetDposIrreversibleInfo() (map[string]interface{}, error) {
	fields := map[string]interface{}{}
	err := ClientCall("dpos_irreversible", &fields)
	return fields, err
}

func GetValidCandidates(epoch int) (*dpos.GlobalState, error) {
	fields := &dpos.GlobalState{}
	err := ClientCall("dpos_validCandidates", fields, epoch)
	return fields, err
}

func GetAvailableStake(epoch int, accountName string) (*big.Int, error) {
	fields := big.NewInt(0)
	err := ClientCall("dpos_availableStake", fields, epoch, accountName)
	return fields, err
}

func GetVotersByCandidate(epoch int, candidateName string, bDetailInfo bool) ([]string, error) {
	fields := make([]string, 0)
	err := ClientCall("dpos_votersByCandidate", &fields, epoch, candidateName, bDetailInfo)
	return fields, err
}

func GetVotersByVoter(epoch int, voterName string, bDetailInfo bool) ([]string, error) {
	fields := make([]string, 0)
	err := ClientCall("dpos_votersByVoter", &fields, epoch, voterName, bDetailInfo)
	return fields, err
}

func GetNextValidCandidates() (*dpos.GlobalState, error) {
	fields := &dpos.GlobalState{}
	err := ClientCall("dpos_nextValidCandidates", fields)
	return fields, err
}

func GetSnapShotTime(epoch int) (map[string]interface{}, error) {
	fields := map[string]interface{}{}
	err := ClientCall("dpos_snapShotTime", &fields, epoch)
	return fields, err
}

func GetEpochByHeight(blockHeight int) (uint64, error) {
	fields := uint64(0)
	err := ClientCall("dpos_epoch", &fields, blockHeight)
	return fields, err
}

func GetPreEpoch(epoch int) (uint64, error) {
	fields := uint64(0)
	err := ClientCall("dpos_prevEpoch", &fields, epoch)
	return fields, err
}

func GetNextEpoch(epoch int) (uint64, error) {
	fields := uint64(0)
	err := ClientCall("dpos_nextEpoch", &fields, epoch)
	return fields, err
}