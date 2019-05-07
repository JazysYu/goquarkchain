package qkcapi

import (
	"github.com/QuarkChain/goquarkchain/account"
	"github.com/QuarkChain/goquarkchain/cluster/config"
	qkcRPC "github.com/QuarkChain/goquarkchain/cluster/rpc"
	"github.com/QuarkChain/goquarkchain/consensus"
	"github.com/QuarkChain/goquarkchain/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
)

type Backend interface {
	AddTransaction(tx *types.Transaction) error
	//AddRawMinorBlock(branch account.Branch, blockData []byte) error
	AddRootBlockFromMine(block *types.RootBlock) error
	ExecuteTransaction(tx *types.Transaction, address account.Address, height *uint64) ([]byte, error)
	GetMinorBlockByHash(blockHash common.Hash, branch account.Branch) (*types.MinorBlock, error)
	GetMinorBlockByHeight(height *uint64, branch account.Branch) (*types.MinorBlock, error)
	GetTransactionByHash(txHash common.Hash, branch account.Branch) (*types.MinorBlock, uint32, error)
	GetTransactionReceipt(txHash common.Hash, branch account.Branch) (*types.MinorBlock, uint32, *types.Receipt, error)
	GetTransactionsByAddress(address account.Address, start []byte, limit uint32) ([]*qkcRPC.TransactionDetail, []byte, error)
	GetLogs(branch account.Branch, address []account.Address, topics []*qkcRPC.Topic, startBlock, endBlock rpc.BlockNumber) ([]*types.Log, error)
	EstimateGas(tx *types.Transaction, address account.Address) (uint32, error)
	GetStorageAt(address account.Address, key common.Hash, height *uint64) (common.Hash, error)
	GetCode(address account.Address, height *uint64) ([]byte, error)
	GasPrice(branch account.Branch) (uint64, error)
	GetWork(branch *account.Branch) consensus.MiningWork
	SubmitWork(branch *account.Branch, headerHash common.Hash, nonce uint64, mixHash common.Hash) bool

	GetRootBlockByNumber(blockNr *uint64) (*types.RootBlock, error)
	GetRootBlockByHash(hash common.Hash) (*types.RootBlock, error)
	NetWorkInfo() map[string]interface{}
	GetPrimaryAccountData(address account.Address, blockHeight *uint64) (*qkcRPC.AccountBranchData, error)
	CurrentBlock() *types.RootBlock
	GetAccountData(address account.Address) (map[account.Branch]*qkcRPC.AccountBranchData, error)
	GetClusterConfig() *config.ClusterConfig
	GetPeers() []qkcRPC.PeerInfoForDisPlay
	GetStats() map[string]interface{}
	GetBlockCount() map[string]interface{}
	SetTargetBlockTime(rootBlockTime *uint32, minorBlockTime *uint32) error
	SetMining(mining bool) error
}

func GetAPIs(apiBackend Backend) []rpc.API {
	return []rpc.API{
		{
			Namespace: "qkc",
			Version:   "1.0",
			Service:   NewPublicBlockChainAPI(apiBackend),
			Public:    true,
		},
		{
			Namespace: "qkc",
			Version:   "1.0",
			Service:   NewPrivateBlockChainAPI(apiBackend),
			Public:    false,
		},
	}
}
