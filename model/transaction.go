package model

import (
	"github.com/ethereum/go-ethereum/common"
)

// Transaction is an Ethereum transaction.
type Transaction struct {
	BlockID int64       `json:"-"`
	TxHash  common.Hash `json:"tx_hash"`
	From    common.Hash `json:"from"`
	To      common.Hash `json:"to"`
	Nonce   uint64      `json:"nonce"`
	Data    string      `json:"data"`
	Value   int64       `json:"value"`
	Logs    []Log       `gorm:"foreignKey:TransactionID;references:TxHash" json:"logs"`
}

// Log represents a contract log event.
type Log struct {
	TransactionID common.Hash `json:"-"`
	Data          []byte      `json:"data"`
	Index         uint        `json:"index"`
}

// GetTransactionByTxHash gets the transaction by transaction hash
func (dao *DAO) GetTransactionByTxHash(txHash []byte) (*Transaction, error) {
	tx := &Transaction{}
	if err := dao.DB.Preload("Logs").Take(tx, "tx_hash = ?", common.BytesToHash(txHash)).Error; err != nil {
		return nil, err
	}

	return tx, nil
}

// CreateTransaction create a transaction
func (dao *DAO) CreateTransaction(tx *Transaction) (*Transaction, error) {
	// TODO validation
	if err := dao.DB.Create(tx).Error; err != nil {
		return nil, err
	}
	return tx, nil
}
