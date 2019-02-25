package elect

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"time"

	"unicode"

	"github.com/vntchain/go-vnt/accounts"
	"github.com/vntchain/go-vnt/accounts/keystore"
	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/core/types"
	vntelection "github.com/vntchain/go-vnt/core/vm/election"
	"github.com/vntchain/go-vnt/vntclient"
)

var (
	emptyHash = common.Hash{}
)

type Election struct {
	cfg    *Config
	vc     *vntclient.Client
	wallet accounts.Wallet // 用于签名的钱包

	ctx     context.Context
	account accounts.Account
}

func NewElection() (*Election, error) {
	e := &Election{
		ctx: context.Background(),
	}
	if err := e.init(); err != nil {
		return nil, err
	}
	return e, nil
}

// init load config and set vntclient
func (e *Election) init() error {
	if err := e.loadCfg("./config.json"); err != nil {
		return err
	}

	if err := e.newClient(); err != nil {
		return err
	}

	e.account = accounts.Account{Address: e.cfg.Sender}
	e.wallet = loadKSWallet(e.cfg.KeystoreDir, e.account)
	if e.wallet == nil {
		return fmt.Errorf("Not find keystore file of account: %s, in directory: %s\n", e.cfg.Sender, e.cfg.KeystoreDir)
	}

	return nil
}

func (e *Election) loadCfg(cfgPath string) error {
	f, err := os.Open(cfgPath)
	if err != nil {
		return fmt.Errorf("Open config file error: %s\n", err)
	}
	defer f.Close()

	config := Config{}
	e.cfg = &config
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&config); err != nil {
		return fmt.Errorf("Decode config file error: %s\n", err)
	}

	return nil
}

func (e *Election) newClient() error {
	var err error
	e.vc, err = vntclient.Dial(e.cfg.RpcUrl)
	if err != nil {
		return fmt.Errorf("Connect to ethereum RPC server failed. url: %s, err: %v\n", e.cfg.RpcUrl, err)
	}
	return err
}

func loadKSWallet(ksDir string, account accounts.Account) accounts.Wallet {
	n, p := keystore.StandardScryptN, keystore.StandardScryptP
	ks := keystore.NewKeyStore(ksDir, n, p)

	for _, wa := range ks.Wallets() {
		if wa.Contains(account) {
			return wa
		}
	}

	return nil
}

// Stake returns a tx hash of staking VNT if passed condition check and tx has been send, or an error if failed.
func (e *Election) Stake(stakeCnt int) (common.Hash, error) {
	// 至少1个VNT
	if stakeCnt <= 1 {
		return emptyHash, fmt.Errorf("stakeCnt = %d is less than 1 VNT", stakeCnt)
	}

	// 抵押数不得多于自己的VNT数量
	b, err := e.vc.BalanceAt(e.ctx, e.cfg.Sender, nil)
	if err != nil {
		return emptyHash, fmt.Errorf("Query balance of account:%s failed, err: %s\n", e.cfg.Sender.String(), err)
	}
	stake := big.NewInt(int64(stakeCnt))
	stakeWei := big.NewInt(0).Mul(stake, big.NewInt(1e+18))
	if stakeWei.Cmp(b) > 0 {
		return emptyHash, fmt.Errorf("stake more than your balance. stake = %s wei, balance = %s wei", stakeWei.String(), b.String())
	}

	unSignTx, err := e.vc.NewElectionTx(e.ctx, e.cfg.Sender, 30000, big.NewInt(18000000000), "stake", stake)
	if err != nil {
		return emptyHash, err
	}

	return e.signAndSendTx(unSignTx)
}

// Unstake returns a tx hash of staking VNT if passed condition check and tx has been send, or an error if failed.
func (e *Election) Unstake() (common.Hash, error) {
	// 用户当前有抵押的VNT
	stake, err := e.vc.StakeAt(e.ctx, e.cfg.Sender)
	if err != nil {
		return emptyHash, err
	}

	if stake != nil {
		// 距离上次抵押超过24小时
		unstakeTime := big.NewInt(0).Add(stake.LastStakeTimeStamp, big.NewInt(vntelection.OneDay))
		now := big.NewInt(time.Now().Unix())
		if now.Cmp(unstakeTime) < 0 {
			return emptyHash, fmt.Errorf("cannot unstake in 24 hours")
		}
	}

	unSignTx, err := e.vc.NewElectionTx(e.ctx, e.cfg.Sender, 30000, big.NewInt(18000000000), "unStake")
	if err != nil {
		return emptyHash, err
	}

	return e.signAndSendTx(unSignTx)
}

func (e *Election) RegisterWitness(nodeName, nodeUrl, website string) (common.Hash, error) {
	// 节点名称必填，由数字和小写字母组成，长度在[3,20]区间。
	// 网址必填，长度在[3,60]区间。
	if err := checkCandi(nodeName, website); err != nil {
		return emptyHash, err
	}

	// 名称和网址不得与其他候选人有重复，不可重复注册
	candidates, err := e.vc.WitnessCandidates(e.ctx)
	if err != nil && err.Error() != errNotFound {
		return emptyHash, err
	}
	if candidates != nil {
		for _, c := range candidates {
			if c.Owner != e.cfg.Sender.String() {
				if c.Url == nodeUrl || c.Website == website {
					return emptyHash, fmt.Errorf("candidate's name or website url is duplicated with a candidate")
				}
			} else if c.Active {
				return emptyHash, fmt.Errorf("candidate is already registered")
			}
		}
	}

	unSignTx, err := e.vc.NewElectionTx(e.ctx, e.cfg.Sender, 30000, big.NewInt(18000000000),
		"registerWitness", []byte(nodeUrl), []byte(website), []byte(nodeName))
	if err != nil {
		return emptyHash, err
	}

	return e.signAndSendTx(unSignTx)
}

func (e *Election) UnregisterWitness() (common.Hash, error) {
	// 账号已注册为见证人
	candidates, err := e.vc.WitnessCandidates(e.ctx)
	if err != nil && err.Error() != errNotFound {
		return emptyHash, err
	}
	if candidates != nil {
		find := false
		for _, c := range candidates {
			if c.Owner == e.cfg.Sender.String() && c.Active {
				find = true
			}
		}
		if !find {
			return emptyHash, fmt.Errorf("account: %s is not registered", e.cfg.Sender)
		}
	}

	unSignTx, err := e.vc.NewElectionTx(e.ctx, e.cfg.Sender, 30000,
		big.NewInt(18000000000), "unregisterWitness")
	if err != nil {
		return emptyHash, err
	}

	return e.signAndSendTx(unSignTx)
}

func checkCandi(name string, website string) error {
	// length check
	if len(name) < 3 || len(name) > 20 {
		return fmt.Errorf("the length of candidate's name should between [3, 20]")
	}
	if len(website) < 3 || len(website) > 60 {
		return fmt.Errorf("the length of candidate's website url should between [3, 60]")
	}

	digitalAndLower := func(s string) bool {
		for _, ru := range s {
			if !unicode.IsDigit(ru) && !unicode.IsLower(ru) {
				return false
			}
		}
		return true
	}
	if !digitalAndLower(name) {
		return fmt.Errorf("andidate's name should consist of digits and lowercase letters")
	}

	return nil
}

// signAndSendTx returns tx hash if sign and send transaction success.
func (e *Election) signAndSendTx(unSignTx *types.Transaction) (common.Hash, error) {
	tx, err := e.wallet.SignTxWithPassphrase(e.account, e.cfg.Password, unSignTx, big.NewInt(int64(e.cfg.ChainID)))
	if err != nil {
		return emptyHash, err
	}
	if err := e.vc.SendTransaction(e.ctx, tx); err != nil {
		return emptyHash, fmt.Errorf("send transaction occur error: %s", err)
	}
	return tx.Hash(), nil
}
