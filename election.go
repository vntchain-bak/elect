package elect

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/vntchain/go-vnt/vntclient"
)

type Election struct {
	cfg *Config
	vc  *vntclient.Client
}

func NewElection() *Election {
	e := &Election{}
	e.init()
	return e
}

// init load config and set vntclient
func (e *Election) init() {
	e.loadCfg("./config.json")

	// TODO vc
	// TODO wallet
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
