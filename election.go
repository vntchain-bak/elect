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
	emptyAddr = common.Address{}
)

// Election creates transactions and sends to hubble network nodes using RPC.
// It checks conditions before creating transaction, signs the transaction
// with your account and password.
type Election struct {
	cfgPath string // config.json的路径
	cfg     *Config
	wallet  accounts.Wallet  // 用于签名的钱包
	account accounts.Account // config中配置的账号

	vc  *vntclient.Client
	ctx context.Context
}

// NewElection returns a Election, or an error if initializing Election failed.
func NewElection(cp string) (*Election, error) {
	e := &Election{
		cfgPath: cp,
		ctx:     context.Background(),
	}
	if err := e.init(); err != nil {
		return nil, err
	}
	return e, nil
}

// init load config and set vntclient
func (e *Election) init() error {
	if err := e.loadCfg(e.cfgPath); err != nil {
		return err
	}

	if err := e.newClient(); err != nil {
		return err
	}

	e.account = accounts.Account{Address: e.cfg.Sender}
	e.wallet = loadKSWallet(e.cfg.KeystoreDir, e.account)
	if e.wallet == nil {
		return fmt.Errorf("Not find keystore file of account: %s, in directory: %s\n", e.cfg.Sender.String(), e.cfg.KeystoreDir)
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
		return emptyHash, fmt.Errorf("you have no stake")
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

// RegisterWitness returns tx hash of registering witness if passed condition check and tx has been send, or an error if failed.
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

// UnregisterWitness returns tx hash of unregistering witness if passed condition check and tx has been send, or an error if failed.
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
			return emptyHash, fmt.Errorf("account: %s is not registered", e.cfg.Sender.String())
		}
	}

	unSignTx, err := e.vc.NewElectionTx(e.ctx, e.cfg.Sender, 30000,
		big.NewInt(18000000000), "unregisterWitness")
	if err != nil {
		return emptyHash, err
	}

	return e.signAndSendTx(unSignTx)
}

// VoteWitness returns tx hash of voting for witness if passed condition check and tx has been send, or an error if failed.
func (e *Election) VoteWitness(wits []string) (common.Hash, error) {
	// 所投候选人不得超过30人
	if len(wits) > 30 {
		return emptyHash, fmt.Errorf("vote too may witness, at most 30")
	}

	// 有抵押的VNT代币
	_, err := e.vc.StakeAt(e.ctx, e.cfg.Sender)
	if err != nil {
		if err.Error() == errNotFound {
			err = fmt.Errorf("please stake before vote")
		}
		return emptyHash, err
	}

	// 距离上次投票或设置代理超过24小时
	vote, err := e.vc.VoteAt(e.ctx, e.cfg.Sender)
	if err != nil && err.Error() != errNotFound {
		return emptyHash, err
	}
	if vote != nil {
		nextVoteTime := big.NewInt(0).Add(vote.LastVoteTimeStamp, big.NewInt(vntelection.OneDay))
		now := big.NewInt(time.Now().Unix())
		if now.Cmp(nextVoteTime) < 0 {
			return emptyHash, fmt.Errorf("cannot vote or set proxy twice within 24 hours")
		}
	}

	// 需要转换为地址
	witnesses := make([]common.Address, len(wits))
	for i, w := range wits {
		if len(w) != len(emptyAddr.String()) {
			return emptyHash, fmt.Errorf("invalid witness address: %s", w)
		}
		witnesses[i] = common.HexToAddress(w)
	}

	unSignTx, err := e.vc.NewElectionTx(e.ctx, e.cfg.Sender, 30000,
		big.NewInt(18000000000), "voteWitnesses", witnesses)
	if err != nil {
		return emptyHash, err
	}

	return e.signAndSendTx(unSignTx)
}

// CancelVote returns tx hash of cancellation vote for witness if passed condition check and tx has been send, or an error if failed.
func (e *Election) CancelVote() (common.Hash, error) {
	// 未设置代理、被投票的见证人列表为空
	vote, err := e.vc.VoteAt(e.ctx, e.cfg.Sender)
	if err != nil {
		if err.Error() == errNotFound {
			return emptyHash, fmt.Errorf("not vote before")
		}
		return emptyHash, err
	}
	if vote != nil {
		if vote.Proxy != emptyAddr {
			return emptyHash, fmt.Errorf("please use cancelProxy to unset your proxy")
		} else if len(vote.VoteCandidates) == 0 {
			return emptyHash, fmt.Errorf("you didn't vote for any witness")
		}
	}

	unSignTx, err := e.vc.NewElectionTx(e.ctx, e.cfg.Sender, 30000,
		big.NewInt(18000000000), "cancelVote")
	if err != nil {
		return emptyHash, err
	}

	return e.signAndSendTx(unSignTx)
}

// StartProxy returns tx hash of becoming a vote proxy if passed condition check and tx has been send, or an error if failed.
func (e *Election) StartProxy() (common.Hash, error) {
	// 已经开启了代理功能，不可重复开启。
	// 已经设置了代理人，不可开启代理功能。
	voter, err := e.vc.VoteAt(e.ctx, e.cfg.Sender)
	if err != nil && err.Error() != errNotFound {
		return emptyHash, err
	}
	if voter != nil {
		if voter.IsProxy {
			return emptyHash, fmt.Errorf("you are vote proxy, no need start proxy again")
		} else if voter.Proxy != emptyAddr {
			return emptyHash, fmt.Errorf("can not become a vote proxy, when you have a vote proxy")
		}
	}

	unSignTx, err := e.vc.NewElectionTx(e.ctx, e.cfg.Sender, 30000,
		big.NewInt(18000000000), "startProxy")
	if err != nil {
		return emptyHash, err
	}

	return e.signAndSendTx(unSignTx)
}

// StopProxy returns tx hash of back to a normal voter if passed condition check and tx has been send, or an error if failed.
func (e *Election) StopProxy() (common.Hash, error) {
	// 是代理人
	voter, err := e.vc.VoteAt(e.ctx, e.cfg.Sender)
	if err != nil {
		if err.Error() == errNotFound {
			return emptyHash, fmt.Errorf("you are not a vote proxy, no need stop proxy")
		}
		return emptyHash, err
	}
	if voter != nil {
		if !voter.IsProxy {
			return emptyHash, fmt.Errorf("you are not a vote proxy, no need stop proxy")
		}
	}

	unSignTx, err := e.vc.NewElectionTx(e.ctx, e.cfg.Sender, 30000,
		big.NewInt(18000000000), "stopProxy")
	if err != nil {
		return emptyHash, err
	}

	return e.signAndSendTx(unSignTx)
}

// SetProxy returns tx hash of setting vote proxy if passed condition check and tx has been send, or an error if failed.
func (e *Election) SetProxy(p string) (common.Hash, error) {
	proxyAddr := common.HexToAddress(p)
	// 不可将自己设置为自己的代理人
	if e.cfg.Sender.String() == proxyAddr.String() {
		return emptyHash, fmt.Errorf("can not set self as your proxy")
	}
	// 有抵押的VNT代币
	_, err := e.vc.StakeAt(e.ctx, e.cfg.Sender)
	if err != nil {
		if err.Error() == errNotFound {
			err = fmt.Errorf("please stake before vote")
		}
		return emptyHash, err
	}

	vote, err := e.vc.VoteAt(e.ctx, e.cfg.Sender)
	if err != nil && err.Error() != errNotFound {
		return emptyHash, err
	}
	if vote != nil {
		// 自己是代理人不可设置他人为代理
		if vote.IsProxy {
			return emptyHash, fmt.Errorf("can not set proxy when you are a proxy")
		}

		// 距离上次投票或设置代理超过24小时
		nextVoteTime := big.NewInt(0).Add(vote.LastVoteTimeStamp, big.NewInt(vntelection.OneDay))
		now := big.NewInt(time.Now().Unix())
		if now.Cmp(nextVoteTime) < 0 {
			return emptyHash, fmt.Errorf("cannot vote or set proxy twice within 24 hours")
		}
	}

	// 要设置的代理人必须是代理
	proxy, err := e.vc.VoteAt(e.ctx, proxyAddr)
	if err != nil && err.Error() != errNotFound {
		return emptyHash, fmt.Errorf("%s is not a proxy", p)
	}
	if proxy != nil && !proxy.IsProxy {
		return emptyHash, fmt.Errorf("%s is not a proxy", p)
	}

	// 需要转换为地址
	unSignTx, err := e.vc.NewElectionTx(e.ctx, e.cfg.Sender, 30000,
		big.NewInt(18000000000), "setProxy", proxyAddr)
	if err != nil {
		return emptyHash, err
	}

	return e.signAndSendTx(unSignTx)
}

// CancelProxy returns tx hash of cancel setting vote proxy if passed condition check and tx has been send, or an error if failed.
func (e *Election) CancelProxy() (common.Hash, error) {
	// 设置过代理人
	voter, err := e.vc.VoteAt(e.ctx, e.cfg.Sender)
	if err != nil {
		if err.Error() == errNotFound {
			return emptyHash, fmt.Errorf("you have no proxy, no need cancel proxy")
		}
		return emptyHash, err
	}
	if voter != nil && voter.Proxy == emptyAddr {
		return emptyHash, fmt.Errorf("you have no proxy, no need cancel proxy")
	}

	unSignTx, err := e.vc.NewElectionTx(e.ctx, e.cfg.Sender, 30000,
		big.NewInt(18000000000), "cancelProxy")
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
