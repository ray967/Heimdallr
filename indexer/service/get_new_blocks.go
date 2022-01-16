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

var headerByNumber = func(client *ethclient.Client) (*types.Header, error) {
	return client.HeaderByNumber(context.Background(), nil)
}

var blockByNumber = func(client *ethclient.Client, headerNumber *big.Int) (*types.Block, error) {
	return client.BlockByNumber(context.Background(), headerNumber)
}

// GetNewBlocks insert pending blocks into channel
func GetNewBlocks(client *ethclient.Client, repo repository.RepositoryService, ch chan *types.Block) {

	header, err := headerByNumber(client)
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
		block, err := blockByNumber(client, header.Number)
		if err != nil {
			fmt.Printf("ERROR: BlockByNumber failed:%v\n", err)
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
