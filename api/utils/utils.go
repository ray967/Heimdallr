package utils

import (
	"heimdallr/model"
	"heimdallr/model/repository"
	"strconv"
)

func LatestBlocks(repo repository.RepositoryService, limit string) (*[]model.Block, error) {
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

func BlockByBlockID(repo repository.RepositoryService, blockID string) (*model.Block, error) {
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

func TransactionByTxHash(repo repository.RepositoryService, hash string) (*model.Transaction, error) {
	txHash := []byte(hash)

	tx, err := repo.DAO.GetTransactionByTxHash(txHash)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
