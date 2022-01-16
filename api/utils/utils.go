package utils

import (
	"heimdallr/model"
	"heimdallr/model/repository"
	"strconv"
)

// LatestBlocks get the latest limit blocks
func LatestBlocks(repo repository.Service, limit string) (*[]model.Block, error) {
	n, err := strconv.Atoi(limit)
	if err != nil {
		return nil, err
	}

	blocks, err := repo.DAO.GetBlocks(n)
	if err != nil {
		return nil, err
	}

	return blocks, nil
}

// BlockByBlockID get the block by block id
func BlockByBlockID(repo repository.Service, blockID string) (*model.Block, error) {
	id, err := strconv.ParseInt(blockID, 10, 64)
	if err != nil {
		return nil, err
	}

	block, err := repo.DAO.GetBlocksByID(id)
	if err != nil {
		return nil, err
	}

	return block, nil
}

// TransactionByTxHash get the transaction by hash
func TransactionByTxHash(repo repository.Service, hash string) (*model.Transaction, error) {
	txHash := []byte(hash)

	tx, err := repo.DAO.GetTransactionByTxHash(txHash)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
