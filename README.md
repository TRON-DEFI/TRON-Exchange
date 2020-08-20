## 1. Query KChart
- url:/api/v1/market/kline/query
- method:get
```
input:param
*exchange ID  exchangeId
*start time   startTime
*end time   endTime
*granu      granu (1min 5min 15min 30min 1h 4h 1d 5d 1w 1m)
eg:
http://127.0.0.1:21110/api/v1/market/kline/query?exchangeId=2&startTime=0&endTime=67686767676&granu=1d
```
output:json
```json
{
  "code": 0,
  "msg": "OK",
  "data": {
    "exchangeId": "3",
    "startTime": "0",
    "endTime": "67686710400",
    "gran": "1d",
    "data": [
      {
        "time": "1553817600",
        "open": "0.100000",
        "high": "0.100000",
        "low": "0.000000",
        "close": "0.100000",
        "volume": "3000.000000"
      }
    ]
  }
}
```

## 2.Query Exchange Pair Information
- url:/api/v1/market/pair/query
- method:get
```
input:param
start index:  start
size per page:  limit
eg:
http://127.0.0.1:21110/api/v1/market/pair/query?start=0&limit=20
```
output:json
```
volume 24h volume (first token)
gain 24h gain percent
highestPrice24h 24h hight price
lowestPrice24h 24h lowest price
volume24h 24h volume(second token)
pairType exchange pair type 1:trc10 2:trc20
{
  "code": 0,
  "msg": "OK",
  "data": {
    "total": 3,
    "rows": [
      {
        "id": 3,
        "volume": 8000,
        "gain": "2.333333",
        "price": 0.1,
        "precision1": 6,
        "precision2": 6,
        "tokenName1": "DDCToken",
        "tokenName2": "TRX",
        "shortName1": "DDC",
        "shortName2": "TRX",
        "tokenAddr1": "TTiRL5d49CMXEvFrm6UWpmbFLaW7Aa27EN",
        "tokenAddr2": "T9yD14Nj9j7xAB4dbGeiX9h8unkKHxuWwb",
        "highestPrice24h": 0.1,
        "lowestPrice24h": 0.1,
        "volume24h": 800,
        "unit": "TRX",
        "pairType": 2,
        "logoUrl": ""
      }
    ]
  }
}
```

## 3.Query Transactions for User
- url:/api/v1/market/user/order
- method:get
```
input:param
start index:  start
size per page:  limit
*user address: userAddr
first token address: tokenAddr1
second token address: tokenAddr2
order status: status(0:pending 1:done 2:cancel pending 3:cancel done)
eg:
http://127.0.0.1:21110/api/v1/market/user/order?start=0&limit=20&userAddr=TYvRQn5ycznkN3Mv6SBpFQyxVqa3xbbe2e
```
output:json
```
volume: user ordered volume
price  price
schedule transaction schedule
curTurnover business volume
{
  "code": 0,
  "msg": "OK",
  "data": {
    "total": 40,
    "rows": [
      {
        "id": 29,
        "shortName1": "John01",
        "shortName2": "TRX",
        "volume": 100,
        "price": 3.5,
        "orderType": 1,
        "orderTime": "1553827240815",
        "orderID": 2,
        "schedule": "0.0000",
        "curTurnover": 0,
        "orderStatus": 0,
        "pairType": 2
      }     
    ]
  }
}
```

## 4.Query Latest Transactions
- url:/api/v1/market/common/order/latest
- method:get
```
input:param
start index:  start
size per page:  limit
*Exchange ID pairID
eg:
http://127.0.0.1:21110/api/v1/market/common/order/latest?pairID=2
```
output:json
```
volume volume
price  price
orderType order type 0:buy order 1:sell order
{
  "code": 0,
  "msg": "OK",
  "data": {
    "total": 6,
    "rows": [
      {
        "blockID": "0",
        "buyAddr": "TGvPbgh63Hub5LR783hf4wJRyC8kn4FkhT",
        "sellAddr": "TGvPbgh63Hub5LR783hf4wJRyC8kn4FkhT",
        "volume": 1000,
        "price": 0.1,
        "orderTime": "1553843182270",
        "unit": "",
        "orderType": 1,
        "pairID": 3
      }
    ]
  }
}
```

## 5.Query All Orders By ExchangeID
- url:/api/v1/market/common/order/list/:pairID
- method:get
```
input:param
*ExchangeID pairID
eg:
http://127.0.0.1:21110/api/v1/market/common/order/list/2
```
output:json
```
orderID order ID
exchangeID exchange ID
bsFlag order type 0:buy order 1:sell order
amount user ordered amount
price price
isCancel status false:pending true:canceled
curTurnover business volume
{
  "code": 0,
  "msg": "OK",
  "data": {
    "buy": [
      {
        "orderID": 0,
        "blockID": 0,
        "exchangeID": 1,
        "ownerAddress": "TGvPbgh63Hub5LR783hf4wJRyC8kn4FkhT",
        "bsFlag": 0,
        "amount": 100,
        "buyTokenAddress": "TXs3V3yi4eMeg7CS4XA7thbDsXaoysv57b",
        "buyTokenAmount": 100,
        "sellTokenAddress": "T9yD14Nj9j7xAB4dbGeiX9h8unkKHxuWwb",
        "sellTokenAmount": 300,
        "price": 3,
        "trxHash": "",
        "orderTime": "",
        "isCancel": false,
        "status": 0,
        "curTurnover": 0
      }
    ],
    "sell": [
      {
        "orderID": 0,
        "blockID": 0,
        "exchangeID": 1,
        "ownerAddress": "TGvPbgh63Hub5LR783hf4wJRyC8kn4FkhT",
        "bsFlag": 1,
        "amount": 100,
        "buyTokenAddress": "TXs3V3yi4eMeg7CS4XA7thbDsXaoysv57b",
        "buyTokenAmount": 100,
        "sellTokenAddress": "T9yD14Nj9j7xAB4dbGeiX9h8unkKHxuWwb",
        "sellTokenAmount": 350,
        "price": 3.5,
        "trxHash": "",
        "orderTime": "",
        "isCancel": false,
        "status": 0,
        "curTurnover": 0
      }
    ],
    "price": "10000"
  }
}
```


## 6.Query Highest&Lowest Price By Exchange ID
- url:/market/price/:exchangeID
- method:get
```
input:param
*exchange ID exchangeID
eg:
http://127.0.0.1:21110/api/v1/market/price/2
```
output:json
```
exchangeID exchange ID
buyHighPrice highest buy order price
sellHighPrice highest sell order price
price  current price
{
  "code": 0,
  "msg": "OK",
  "data": {
    "exchangeID": 2,
    "buyHighPrice": 1.973,
    "sellLowPrice": 1.974002,
    "price": 1.963935
  }
}
```

## 7.Query Token Information
- url:/market/price/:exchangeID
- method:get
```
input:param
*token address address
eg:
http://127.0.0.1:21110/api/v1/market/common/tokenInfo/query?address=TTKAhucKwU3zpjbXkodSofVwMGSrNxA75U
```
output:json
```
address token address
fullName full name
shortName short name
circulation total circulation
precision token precision
description token description
websiteUrl web url
logoUrl logo url
{
  "code": 0,
  "msg": "OK",
  "data": {
    "id": 3,
    "address": "TTKAhucKwU3zpjbXkodSofVwMGSrNxA75U",
    "fullName": "AAAToken1",
    "shortName": "AAA1",
    "circulation": 10000000,
    "precision": 2,
    "description": "AAAToken AAA1",
    "websiteUrl": "https://me",
    "logoUrl": "http://127.0.0.1/images/token/TEKg9MaPXUsHdEyci2GMgKsytywkFCMxk1-111.png"
  }
}
```
