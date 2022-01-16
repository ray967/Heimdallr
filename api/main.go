package main

import (
	"fmt"
	"heimdallr/api/utils"
	"heimdallr/model/repository"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

var repo repository.Service

func latestBlocks(c *gin.Context) {
	blocks, err := utils.LatestBlocks(repo, c.Query("limit"))
	if err != nil {
		fmt.Print(err.Error())
	}

	c.JSON(200, blocks)
}

func blockByBlockID(c *gin.Context) {
	block, err := utils.BlockByBlockID(repo, c.Param("id"))
	if err != nil {
		fmt.Print(err.Error())
	}

	txs := []common.Hash{}
	for _, v := range block.Transactions {
		txs = append(txs, v.TxHash)
	}

	c.JSON(200, struct {
		BlockNum     int64         `json:"block_num"`
		BlockHash    common.Hash   `json:"block_hash"`
		BlockTime    uint64        `json:"block_time"`
		ParentHash   common.Hash   `json:"parent_hash"`
		Transactions []common.Hash `json:"transactions"`
	}{
		BlockNum:     block.BlockNum,
		BlockHash:    block.BlockHash,
		BlockTime:    block.BlockTime,
		ParentHash:   block.ParentHash,
		Transactions: txs,
	})
}

func transactionByTxHash(c *gin.Context) {
	tx, err := utils.TransactionByTxHash(repo, c.Param("txHash"))
	if err != nil {
		fmt.Print(err.Error())
	}
	c.JSON(200, tx)
}

func main() {

	err := repo.Init()
	if err != nil {
		fmt.Print(err.Error())
	}

	server := gin.Default()
	server.GET("/blocks", latestBlocks)
	server.GET("/blocks/:id", blockByBlockID)
	server.GET("/transaction/:txHash", transactionByTxHash)
	server.Run(":8888")

}
