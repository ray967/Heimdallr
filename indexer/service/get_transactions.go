package service

import (
	"context"
	"fmt"
	"heimdallr/conf"
	"heimdallr/model"
	"heimdallr/model/repository"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var transactionByHash = func(client *ethclient.Client, hash common.Hash) (tx *types.Transaction, isPending bool, err error) {
	return client.TransactionByHash(context.Background(), hash)
}

var transactionReceipt = func(client *ethclient.Client, hash common.Hash) (*types.Receipt, error) {
	return client.TransactionReceipt(context.Background(), hash)
}

// GetTransactions gets the transactions when it's done
func GetTransactions(client *ethclient.Client, repo repository.Service, ch chan *types.Block) {
	for len(ch) > 0 {
		select {
		case block, ok := <-ch:
			if ok {
				wg := sync.WaitGroup{}
				txs := block.Transactions()
				for i := 0; i < txs.Len(); i++ {
					wg.Add(1)
					go getTransaction(client, repo, txs[i].Hash(), &wg)
				}
				wg.Wait()
				updateBlockDone(repo, block.Header().Number.Int64())
			} else {
				fmt.Println("Channel closed!")
			}
		default:
			fmt.Println("No value ready, moving on.")
		}
	}
}

func getTransaction(client *ethclient.Client, repo repository.Service, hash common.Hash, wg *sync.WaitGroup) {
	defer wg.Done()

	var (
		tx        *types.Transaction
		isPending bool = true
		err       error
	)

	for isPending {
		tx, isPending, err = transactionByHash(client, hash)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		var res *model.Transaction
		if isPending {
			// if it's still pending then retry
			// TODO hadle never be done
			time.Sleep(conf.TransactionPendingTimeInterval)
			continue
		} else {
			to := common.Hash{}
			if tx.To() != nil {
				// null when its a contract creation transaction.
				to = tx.To().Hash()
			}
			from, err := types.Sender(types.NewEIP155Signer(tx.ChainId()), tx)
			if err != nil {
				from, err = types.Sender(types.HomesteadSigner{}, tx)
			}
			res = &model.Transaction{
				TxHash: tx.Hash(),
				From:   from.Hash(),
				To:     to,
				Nonce:  tx.Nonce(),
				Value:  tx.Value().Int64(),
			}
			recp, err := transactionReceipt(client, hash)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			res.BlockID = recp.BlockNumber.Int64()
			res.Logs = []model.Log{}

			for _, v := range recp.Logs {
				res.Logs = append(res.Logs, model.Log{
					Data:  v.Data,
					Index: v.Index,
				})
			}
		}
		// insert transaction into db
		repo.DAO.CreateTransaction(res)
	}
}

func updateBlockDone(repo repository.Service, blockID int64) {
	err := repo.DAO.UpdateBlockDone(blockID)
	if err != nil {
		panic(fmt.Sprintf("UpdateBlockDone failed, block_num: %d", blockID))
	}
}
