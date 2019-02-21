package elect

import (
	"encoding/json"
	"fmt"
	"os"

	"context"

	"math/big"

	"github.com/vntchain/go-vnt/accounts"
	"github.com/vntchain/go-vnt/accounts/keystore"
	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/core/types"
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

// Stake return a tx hash if send transaction success or an error when f the tx may be failed in execution.
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

// signAndSendTx return tx hash if sign and send transaction success.
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
