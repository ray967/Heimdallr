package model

import (
	"gorm.io/gorm"
)

type DAO struct {
	DB *gorm.DB
	// TODO cache
}

var NewDAO = func(db *gorm.DB) DAOAbstracter {
	dao := &DAO{
		DB: db,
	}
	return dao
}

type DAOAbstracter interface {
	GetBlocks(n int) (*[]Block, error)
	GetBlocksByID(id int64) (*Block, error)
	GetTransactionByTxHash(txHash []byte) (*Transaction, error)
	CreateBlock(block *Block) (*Block, error)
}
