package model

import (
	"github.com/ethereum/go-ethereum/common"
)

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

type Log struct {
	TransactionID common.Hash `json:"-"`
	Data          []byte      `json:"data"`
	Index         uint        `json:"index"`
}

func (dao *DAO) GetTransactionByTxHash(txHash []byte) (*Transaction, error) {
	tx := &Transaction{}
	if err := dao.DB.Preload("Logs").Take(tx, "tx_hash = ?", common.BytesToHash(txHash)).Error; err != nil {
		return nil, err
	}

	return tx, nil
}
