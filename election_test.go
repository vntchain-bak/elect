package elect

import (
	"reflect"
	"testing"

	"github.com/vntchain/go-vnt/common"
)

func TestElectionLoadCfg(t *testing.T) {
	e := &Election{}
	if err := e.loadCfg("./tests/config.json"); err != nil {
		t.Errorf("want no error, got: %s", err)
	}

	cfg := &Config{
		Sender:      common.HexToAddress("0x122369f04f32269598789998de33e3d56e2c507a"),
		Password:    "",
		KeystoreDir: "./tests/keystore",
		RpcUrl:      "http://localhost:8880",
		ChainID:     1234,
	}

	if !reflect.DeepEqual(e.cfg, cfg) {
		t.Errorf("config want: %v, got: %v", cfg, e.cfg)
	}
}
