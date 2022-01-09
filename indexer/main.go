package main

import (
	"fmt"
	"heimdallr/conf"
	"heimdallr/indexer/service"
	"heimdallr/model/repository"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	ch := make(chan *types.Block, 10000)
	client, err := ethclient.Dial(conf.RPCEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	repo := repository.RepositoryService{}
	err = repo.Init()
	if err != nil {
		panic(fmt.Sprintf("Init repo failed:%v\n", err))
	}
	for {
		go service.GetNewBlocks(client, repo, ch)
		go service.GetTransactions(client, repo, ch)
		time.Sleep(conf.IndexerTimeInterval)
	}
}
