package service

import (
	"context"
	"fmt"
	"heimdallr/model"
	"heimdallr/model/repository"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	preHead *big.Int = nil
)

// GetNewBlocks insert pending blocks into channel
func GetNewBlocks(client *ethclient.Client, repo repository.RepositoryService, ch chan *types.Block) {

	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// check preHead has val
	if preHead == nil {
		preHead = header.Number
		preHead = preHead.Add(preHead, big.NewInt(-1))
	}

	// while has new blocks
	for preHead.Cmp(header.Number) == -1 {
		preHead.Add(preHead, big.NewInt(1))
		block, err := client.BlockByNumber(context.Background(), header.Number)
		if err != nil {
			panic(fmt.Sprintf("BlockByNumber failed:%v\n", err))
		}
		// save to db
		repo.DAO.CreateBlock(&model.Block{
			BlockNum:   block.Header().Number.Int64(),
			BlockHash:  block.Hash(),
			BlockTime:  block.Header().Time,
			ParentHash: block.ParentHash(),
			IsPending:  true,
		})
		// add to channel
		ch <- block
	}
}
