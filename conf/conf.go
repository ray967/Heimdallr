package conf

import "time"

// RPCEndpoint represents a BSC RPC Endpoints: Testnet
const RPCEndpoint string = "https://data-seed-prebsc-2-s3.binance.org:8545"

// MysqlAddress ...
const MysqlAddress string = "root:0000@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True"

// IndexerTimeInterval represents the interval between indexer services
const IndexerTimeInterval time.Duration = time.Second * 1

// TransactionPendingTimeInterval represents the interval between indexer services check the transactions
const TransactionPendingTimeInterval time.Duration = time.Millisecond
