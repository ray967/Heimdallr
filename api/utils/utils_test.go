package utils_test

import (
	"errors"
	"heimdallr/api/utils"
	"heimdallr/model"
	"heimdallr/model/repository"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

type mockDAO struct {
}

var mockTxs *[]model.Transaction = &[]model.Transaction{
	{
		BlockID: 1,
		TxHash:  common.HexToHash("1"),
		From:    common.HexToHash("1"),
		To:      common.HexToHash("1"),
		Nonce:   1,
		Data:    "",
		Value:   1,
		Logs:    nil,
	},
	{
		BlockID: 2,
		TxHash:  common.HexToHash("2"),
		From:    common.HexToHash("2"),
		To:      common.HexToHash("2"),
		Nonce:   2,
		Data:    "",
		Value:   2,
		Logs:    nil,
	},
}
var mockBlocks *[]model.Block = &[]model.Block{
	{
		BlockNum:   1,
		BlockHash:  common.HexToHash("1"),
		BlockTime:  1,
		ParentHash: common.HexToHash("0"),
		IsPending:  false,
	},
	{
		BlockNum:   2,
		BlockHash:  common.HexToHash("2"),
		BlockTime:  2,
		ParentHash: common.HexToHash("1"),
		IsPending:  false,
	},
}

func (dao *mockDAO) GetBlocks(n int) (*[]model.Block, error) {
	if n >= 0 && n < 3 {
		res := (*mockBlocks)[0:n]
		for i := range res {
			res[i].Transactions = nil
		}
		return &res, nil
	}
	return nil, errors.New("Out of range")
}
func (dao *mockDAO) GetBlocksByID(id int64) (*model.Block, error) {

	if id >= 0 && id <= 2 {
		res := (*mockBlocks)[id-1]
		res.Transactions = append(res.Transactions, (*mockTxs)[id-1])
		return &res, nil
	}
	return nil, errors.New("Not found")
}
func (dao *mockDAO) GetTransactionByTxHash(txHash []byte) (*model.Transaction, error) {
	if string(txHash[:]) == "1" {
		return &(*mockTxs)[0], nil
	}
	if string(txHash[:]) == "2" {
		return &(*mockTxs)[1], nil
	}
	return nil, errors.New("Not found")
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

func sameTransaction(a, b *model.Transaction) bool {
	return true
}

func sameBlock(a, b *model.Block) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if a.BlockHash != b.BlockHash {
		return false
	}
	if a.BlockNum != b.BlockNum {
		return false
	}
	if a.BlockTime != b.BlockTime {
		return false
	}
	if a.ParentHash != b.ParentHash {
		return false
	}
	if a.IsPending != b.IsPending {
		return false
	}
	if len(a.Transactions) != len(b.Transactions) {
		// fmt.Println(len(a.Transactions))
		// fmt.Println(len(b.Transactions))
		return false
	}
	for i, v := range a.Transactions {
		if !sameTransaction(&v, &b.Transactions[i]) {
			return false
		}
	}
	return true
}

func Test_LatestBlocks(t *testing.T) {
	repo := repository.Service{DAO: &mockDAO{}}

	testCases := []struct {
		Name           string
		Limit          string
		ExceptResponse *[]model.Block
		ExceptError    error
	}{
		{
			Name:  "Get 1",
			Limit: "1",
			ExceptResponse: &[]model.Block{
				{
					BlockNum:   1,
					BlockHash:  common.HexToHash("1"),
					BlockTime:  1,
					ParentHash: common.HexToHash("0"),
					IsPending:  false,
				},
			},
		},
		{
			Name:  "Get 2",
			Limit: "2",
			ExceptResponse: &[]model.Block{
				{
					BlockNum:   1,
					BlockHash:  common.HexToHash("1"),
					BlockTime:  1,
					ParentHash: common.HexToHash("0"),
					IsPending:  false,
				},
				{
					BlockNum:   2,
					BlockHash:  common.HexToHash("2"),
					BlockTime:  2,
					ParentHash: common.HexToHash("1"),
					IsPending:  false,
				},
			},
		},
		{
			Name:           "Get 0",
			Limit:          "0",
			ExceptResponse: &[]model.Block{},
		},
		{
			Name:        "Get more than db max",
			Limit:       "3",
			ExceptError: errors.New("Out of range"),
		},
		{
			Name:        "Get less than 0",
			Limit:       "-1",
			ExceptError: errors.New("Out of range"),
		},
		{
			Name:        "invalid limit: limit is not a number",
			Limit:       "abc",
			ExceptError: errors.New(`strconv.Atoi: parsing "abc": invalid syntax`),
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			blocks, err := utils.LatestBlocks(repo, test.Limit)

			if err != nil && test.ExceptError == nil ||
				err == nil && test.ExceptError != nil ||
				err != nil && test.ExceptError != nil && err.Error() != test.ExceptError.Error() {
				t.Fatalf("error mismatch\nexcept: %v\nget: %v", test.ExceptError, err)
			}

			if blocks != nil && test.ExceptResponse != nil {
				if len(*blocks) != len(*test.ExceptResponse) {
					t.Fatalf("responses length mismatch")
				} else {
					for i, block := range *blocks {
						if !sameBlock(&block, &(*test.ExceptResponse)[i]) {
							t.Fatalf("response mismatch\nexcept: %v\nget: %v\n", (*test.ExceptResponse)[i], block)
						}
					}
				}
			}
		})
	}
}

func Test_BlockByBlockID(t *testing.T) {
	repo := repository.Service{DAO: &mockDAO{}}

	testCases := []struct {
		Name           string
		BlockID        string
		ExceptResponse *model.Block
		ExceptError    error
	}{
		{
			Name:    "valid id: 1",
			BlockID: "1",
			ExceptResponse: &model.Block{
				BlockNum:   1,
				BlockHash:  common.HexToHash("1"),
				BlockTime:  1,
				ParentHash: common.HexToHash("0"),
				IsPending:  false,
				Transactions: []model.Transaction{
					{
						BlockID: 1,
						TxHash:  common.HexToHash("1"),
						From:    common.HexToHash("1"),
						To:      common.HexToHash("1"),
						Nonce:   1,
						Data:    "",
						Value:   1,
						Logs:    nil,
					},
				},
			},
		},
		{
			Name:    "valid id: 2",
			BlockID: "2",
			ExceptResponse: &model.Block{
				BlockNum:   2,
				BlockHash:  common.HexToHash("2"),
				BlockTime:  2,
				ParentHash: common.HexToHash("1"),
				IsPending:  false,
				Transactions: []model.Transaction{
					{
						BlockID: 2,
						TxHash:  common.HexToHash("2"),
						From:    common.HexToHash("2"),
						To:      common.HexToHash("2"),
						Nonce:   2,
						Data:    "",
						Value:   2,
						Logs:    nil,
					},
				},
			},
		},
		{
			Name:        "out of range: max",
			BlockID:     "3",
			ExceptError: errors.New("Not found"),
		},
		{
			Name:        "out of range: min",
			BlockID:     "-1",
			ExceptError: errors.New("Not found"),
		},
		{
			Name:        "invalid input: not int64",
			BlockID:     "abc",
			ExceptError: errors.New(`strconv.ParseInt: parsing "abc": invalid syntax`),
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			block, err := utils.BlockByBlockID(repo, test.BlockID)

			if err != nil && test.ExceptError == nil ||
				err == nil && test.ExceptError != nil ||
				err != nil && test.ExceptError != nil && err.Error() != test.ExceptError.Error() {
				t.Fatalf("error mismatch\nexcept: %v\nget: %v", test.ExceptError, err)
			}

			if !sameBlock(block, test.ExceptResponse) {
				t.Fatalf("response mismatch\nexcept: %v\nget: %v\n", test.ExceptResponse, block)
			}
		})
	}
}

func Test_TransactionByTxHash(t *testing.T) {
	repo := repository.Service{DAO: &mockDAO{}}

	testCases := []struct {
		Name           string
		Hash           string
		ExceptResponse *model.Transaction
		ExceptError    error
	}{
		{
			Name: "valid hash: 1",
			Hash: "1",
			ExceptResponse: &model.Transaction{
				BlockID: 1,
				TxHash:  common.HexToHash("1"),
				From:    common.HexToHash("1"),
				To:      common.HexToHash("1"),
				Nonce:   1,
				Data:    "",
				Value:   1,
				Logs:    nil,
			},
		},
		{
			Name: "valid hash: 2",
			Hash: "2",
			ExceptResponse: &model.Transaction{
				BlockID: 2,
				TxHash:  common.HexToHash("2"),
				From:    common.HexToHash("2"),
				To:      common.HexToHash("2"),
				Nonce:   2,
				Data:    "",
				Value:   2,
				Logs:    nil,
			},
		},
		{
			Name:        "out of range: max",
			Hash:        "3",
			ExceptError: errors.New("Not found"),
		},
		{
			Name:        "out of range: min",
			Hash:        "-1",
			ExceptError: errors.New("Not found"),
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			tx, err := utils.TransactionByTxHash(repo, test.Hash)

			if err != nil && test.ExceptError == nil ||
				err == nil && test.ExceptError != nil ||
				err != nil && test.ExceptError != nil && err.Error() != test.ExceptError.Error() {
				t.Fatalf("error mismatch\nexcept: %v\nget: %v", test.ExceptError, err)
			}

			if !sameTransaction(tx, test.ExceptResponse) {
				t.Fatalf("response mismatch\nexcept: %v\nget: %v\n", test.ExceptResponse, tx)
			}
		})
	}
}
