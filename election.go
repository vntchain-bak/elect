package elect

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/vntchain/go-vnt/accounts"
	"github.com/vntchain/go-vnt/accounts/keystore"
	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/vntclient"
)

type Election struct {
	cfg    *Config
	vc     *vntclient.Client
	wallet accounts.Wallet // 用于签名的钱包
}

func NewElection() (*Election, error) {
	e := &Election{}
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

	account := accounts.Account{Address: common.HexToAddress(e.cfg.Sender)}
	e.wallet = loadKSWallet(e.cfg.KeystoreDir, account)
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
	fmt.Printf("Connect to ethereum RPC server success. url: %s", e.cfg.RpcUrl)
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
