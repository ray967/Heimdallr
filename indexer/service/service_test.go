package service

import (
	"heimdallr/model"
	"heimdallr/model/repository"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type mockDAO struct {
}

func (dao *mockDAO) GetBlocks(n int) (*[]model.Block, error) {

	return nil, nil
}
func (dao *mockDAO) GetBlocksByID(id int64) (*model.Block, error) {

	return nil, nil
}
func (dao *mockDAO) GetTransactionByTxHash(txHash []byte) (*model.Transaction, error) {

	return nil, nil
}
func (dao *mockDAO) CreateBlock(block *model.Block) (*model.Block, error) {

	return nil, nil
}
func (dao *mockDAO) CreateTransaction(tx *model.Transaction) (*model.Transaction, error) {

	return nil, nil
}
func (dao *mockDAO) UpdateBlockDone(blockID int64) error {

	return nil
}

type mockTrieHasher struct {
}

func (*mockTrieHasher) Reset()                {}
func (*mockTrieHasher) Update([]byte, []byte) {}
func (*mockTrieHasher) Hash() common.Hash     { return common.Hash{} }

func Test_GetNewBlocks(t *testing.T) {
	repo := repository.RepositoryService{DAO: &mockDAO{}}

	testCases := []struct {
		Name           string
		PreHead        *big.Int
		HeaderByNumber func(client *ethclient.Client) (*types.Header, error)
		BlockByNumber  func(client *ethclient.Client, headerNumber *big.Int) (*types.Block, error)
		Channel        chan *types.Block
	}{
		{
			Name:    "empty",
			PreHead: nil,
			HeaderByNumber: func(client *ethclient.Client) (*types.Header, error) {
				return &types.Header{
					Number: big.NewInt(1),
				}, nil
			},
			BlockByNumber: func(client *ethclient.Client, headerNumber *big.Int) (*types.Block, error) {
				return &types.Block{}, nil
			},
			Channel: make(chan *types.Block, 10000),
		},
		{
			Name:    "valid",
			PreHead: big.NewInt(0),
			HeaderByNumber: func(client *ethclient.Client) (*types.Header, error) {
				return &types.Header{
					Number: big.NewInt(1),
				}, nil
			},
			BlockByNumber: func(client *ethclient.Client, headerNumber *big.Int) (*types.Block, error) {
				return types.NewBlock(&types.Header{
					Number: big.NewInt(1),
				}, nil, nil, nil, &mockTrieHasher{}), nil
			},
			Channel: make(chan *types.Block, 10000),
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			preHead = test.PreHead
			headerByNumber = test.HeaderByNumber
			blockByNumber = test.BlockByNumber
			GetNewBlocks(nil, repo, test.Channel)
		})
	}
}

func mockChannel(blocks []*types.Block) chan *types.Block {
	ch := make(chan *types.Block, 10)
	for _, v := range blocks {
		ch <- v
	}
	return ch
}
func Test_GetTransactions(t *testing.T) {
	repo := repository.RepositoryService{DAO: &mockDAO{}}

	testCases := []struct {
		Name               string
		TransactionByHash  func(client *ethclient.Client, hash common.Hash) (tx *types.Transaction, isPending bool, err error)
		TransactionReceipt func(client *ethclient.Client, hash common.Hash) (*types.Receipt, error)
		Channel            chan *types.Block
	}{
		{
			Name: "empty",
			TransactionByHash: func(client *ethclient.Client, hash common.Hash) (tx *types.Transaction, isPending bool, err error) {
				return &types.Transaction{}, false, nil
			},
			TransactionReceipt: func(client *ethclient.Client, hash common.Hash) (*types.Receipt, error) {
				return &types.Receipt{}, nil
			},
			Channel: make(chan *types.Block, 10000),
		},
		{
			Name: "valid",
			TransactionByHash: func(client *ethclient.Client, hash common.Hash) (tx *types.Transaction, isPending bool, err error) {
				return types.NewTransaction(1, common.Address{}, big.NewInt(1), 1, big.NewInt(1), nil), false, nil
			},
			TransactionReceipt: func(client *ethclient.Client, hash common.Hash) (*types.Receipt, error) {
				return &types.Receipt{BlockNumber: big.NewInt(1), Logs: []*types.Log{{}}}, nil
			},
			Channel: mockChannel([]*types.Block{types.NewBlock(&types.Header{
				Number: big.NewInt(1),
			}, []*types.Transaction{types.NewTransaction(1, common.Address{}, big.NewInt(1), 1, big.NewInt(1), nil)}, nil, nil, &mockTrieHasher{})}),
		},
	}
	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			transactionByHash = test.TransactionByHash
			transactionReceipt = test.TransactionReceipt
			GetTransactions(nil, repo, test.Channel)
		})
	}

}
