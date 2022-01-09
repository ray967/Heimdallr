package model

import (
	"github.com/ethereum/go-ethereum/common"
)

type Block struct {
	BlockNum     int64         `json:"block_num"`
	BlockHash    common.Hash   `json:"block_hash"`
	BlockTime    uint64        `json:"block_time"`
	ParentHash   common.Hash   `json:"parent_hash"`
	IsPending    bool          `json:"is_pending"`
	Transactions []Transaction `gorm:"foreignKey:BlockID;references:BlockNum" json:"-"`
}

func (dao *DAO) CreateBlock(block *Block) (*Block, error) {
	// TODO validation
	if err := dao.DB.Create(block).Error; err != nil {
		return nil, err
	}
	return block, nil
}

func (dao *DAO) GetBlocks(n int) (*[]Block, error) {
	blocks := &[]Block{}
	if err := dao.DB.Order("block_time desc").Limit(n).Find(blocks).Error; err != nil {
		return nil, err
	}
	return blocks, nil
}

func (dao *DAO) GetBlocksByID(id int64) (*Block, error) {
	block := &Block{}
	if err := dao.DB.Preload("Transactions").Take(block, "block_num = ?", id).Error; err != nil {
		return nil, err
	}

	return block, nil
}

func (dao *DAO) UpdateBlockDone(blockID int64) error {
	// TODO validation
	if err := dao.DB.Model(&Block{}).Where("block_num = ?", blockID).Update("is_pending", false).Error; err != nil {
		return err
	}
	return nil
}
