

Heimdallr is an Ethereum blockchain service based on PostgreSQL and Binance Chain RPC web3 API.

## Services
1. API service
2. Ethereum block indexer service

### API service
- [GET] /blocks?limit=n 
  - 回傳最新的 n 個 blocks (不包含所有 tx hash) 
```json
{ 
    “blocks”: [ 
        { 
            “block_num”: 1, 
            “block_hash”: “”, 
            “block_time”: 123456789, 
            “parent_hash”: “”, 
        } 
    ] 
}
```

- [GET] /blocks/:id 
  - 回傳單一 block by block id (包含所有 tx hash) 
```json
{ 
    “block_num”: 1, 
    “block_hash”: “”, 
    “block_time”: 123456789, 
    “parent_hash”: “”, 
    “transactions”: [ 
        “0x12345678”, 
        “0x87654321” 
    ] 
}
```

- [GET] /transaction/:txHash 
  - 回傳 transaction data with event logs 
```json
{ 
    “tx_hash”: “0x6666”, 
    “from”: “0x4321”, 
    “to”: “0x1234”, 
    “nonce”: 1, 
    “data”: “0xeb12”, 
    “value”: “12345678” 
    “logs”: [ 
        { 
            “index”: 0, 
            “data”: “0x12345678”,
        } 
    ] 
}
```

### Ethereum block indexer service

- 根據 web3 API 透過 RPC 將區塊內的資料掃進 db
  - 需以平行的方式從 block n 開始進行掃描,直到掃到最新區塊後繼續執行

## TODO
- [x] Database
- API service
  - [x] /blocks?limit=n
  - [x] /blocks/:id
  - [x] /transaction/:txHash
  - [ ] Test
- [x] Ethereum block indexer service
  - [ ] Test
- Bonus
  - [ ] 掃的速度以及量會直接影響到記憶體使用,如何在有限度的記憶體下,做可量化的設計
  - [ ] 可使用你所知的額外服務加強存取效能
  - [ ] 區塊鏈有常見的分叉現象,意味著在當下最新區塊的前 n 塊都處於不穩定狀態,假設 n = 20,如何做到可以掃到最新塊又可以隨後將不穩定區塊替換成穩替區塊,且 API 回傳會註明是否是穩定區塊
  - [ ] 如何最小限度的使用 RPC,達到需求
- [ ] Logs

## Refs
- lib & framework
  - gorm 
  - gin 
  - go-ethereum, ethclient 
- web3 API
  - https://eth.wiki/json-rpc/API 
  - https://eth.wiki/json-rpc/API#eth_getblockbynumber 
  - https://eth.wiki/json-rpc/API#eth_gettransactionbyhash 
  - https://eth.wiki/json-rpc/API#eth_gettransactionreceipt 
- RPC endpoint 
  - https://data-seed-prebsc-2-s3.binance.org:8545/
