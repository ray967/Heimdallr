package model

import (
	"gorm.io/gorm"
)

// DAO represents the data access object
type DAO struct {
	DB *gorm.DB
	// TODO cache
}

// NewDAO ...
var NewDAO = func(db *gorm.DB) DAOAbstracter {
	dao := &DAO{
		DB: db,
	}
	return dao
}

// DAOAbstracter represents the DAO interface
type DAOAbstracter interface {
	GetBlocks(n int) (*[]Block, error)
	GetBlocksByID(id int64) (*Block, error)
	GetTransactionByTxHash(txHash []byte) (*Transaction, error)
	CreateBlock(block *Block) (*Block, error)
	CreateTransaction(tx *Transaction) (*Transaction, error)
	UpdateBlockDone(blockID int64) error
}
