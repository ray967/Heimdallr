package conf

import "time"

// RPCEndpoint represents a BSC RPC Endpoints: Testnet
const RPCEndpoint string = "https://data-seed-prebsc-2-s3.binance.org:8545"

const MysqlAddress string = "root:0000@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True"

const IndexerTimeInterval time.Duration = time.Second * 1
