package elect

import "github.com/vntchain/go-vnt/common"

// Config contains accounts information and RPC information of a vnt node.
type Config struct {
	// Account information
	Sender      common.Address `json:sender`
	Password    string         `json:password`
	KeystoreDir string         `json:keystoreDir`

	// Network information
	RpcUrl  string `json:rpcUrl` // ip:port, example: localhost:8080
	ChainID int    `json:chainID`
}
