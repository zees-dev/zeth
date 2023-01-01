
# list of api requests

## settings [GET]
curl \
	-H "Content-Type: application/json" \
	http://localhost:7000/api/v1/settings

## settings update default node [PUT]
curl -X PUT \
	-H "Content-Type: application/json" \
	-d '{"uuid": "d0a7c824-67e7-4bab-905d-077220260698"}' \
	http://localhost:7000/api/v1/settings/node

---

## nodes all [GET]
curl --silent http://localhost:3000/api/v1/endpoints | jq .

## node [GET]
curl --silent http://localhost:3000/api/v1/endpoints/endpoints:x02kzu2ecx6o48yhbe8b | jq .

## node [DELETE]
curl -X DELETE http://localhost:3000/api/v1/endpoints/endpoints:x02kzu2ecx6o48yhbe8b

## node (register remote) [POST]
curl --silent \
	-X POST \
	-H "Content-Type: application/json" \
	-d '{"name": "eth-node", "is_dev": false, "enabled": true, "date_added": "2022-12-26T23:15:51.581789Z", "rpc_http": "https://rpc.flashbots.net", "rpc_ws": "wss://mainnet.infura.io/ws/v3/8792dc3bbc3743f6b884807fb6a22525" }' \
	http://localhost:3000/api/v1/endpoints

## node (update remote) [PUT]
curl --silent \
	-X PUT \
	-H "Content-Type: application/json" \
	-d '{"name": "test-node", "is_dev": false, "enabled": false, "date_added": "2022-12-26T23:15:51.581789Z", "rpc_http": "https://data-seed-prebsc-1-s1.binance.org:8545" }' \
	http://localhost:3000/api/v1/endpoints/endpoints:8od88gdzdfxpq4tyun6k

## node RPC (http)
curl -X POST -v \
	http://localhost:3000/api/v1/endpoints/endpoints:wmu2af50bcbz0a5q93ic/rpc \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_accounts","params":[],"id":1}'

## node RPC (ws)
websocat ws://localhost:3000/api/v1/endpoints/endpoints:wmu2af50bcbz0a5q93ic/rpc
wscat -c ws://localhost:3000/api/v1/endpoints/endpoints:wmu2af50bcbz0a5q93ic/rpc
{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}


## node RPC sse
curl -v http://localhost:3000/api/v1/endpoints/endpoints:wmu2af50bcbz0a5q93ic/rpc/events

---

curl https://rpc.flashbots.net \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"eth_chainId"}'

curl http://localhost:8545 \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"eth_chainId"}'

---

## start local geth node
docker run --rm -it --name geth -p 8545:8545 -p 8546:8546 ethereum/client-go --dev --http --http.addr=0.0.0.0
docker run --rm -it --name geth -p 8545:8545 -p 8546:8546 ethereum/client-go --dev --ws --ws.addr=0.0.0.0

## use ws client to send payloads
wscat -c ws://127.0.0.1:8546
websocat ws://127.0.0.1:8546
{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}

---

```sh
curl -k -L -s --compressed POST \
--header "Accept: application/json" \
--header "NS: endpoints" \
--header "DB: zeth-db" \
--user "root:root" \
--data 'CREATE endpoints:test SET name="test";' \
http://localhost:8000/sql
```

```sh
curl -k -L -s --compressed POST \
--header "Accept: application/json" \
--header "NS: endpoints" \
--header "DB: zeth-db" \
--user "root:root" \
--data "SELECT * FROM endpoints WHERE id == endpoints:test;" \
http://localhost:8000/sql
```


---


--http2-prior-knowledge

curl -v https://bsc-dataseed.binance.org \
  -X POST \
  -H "X-Forwarded-For: localhost" \
  -H "Content-Type:" \
  -H "User-Agent:" \
  -H "Accept:" \
  -d '{"jsonrpc":"2.0","id":2,"method":"web3_clientVersion","params":[]}'
# > POST / HTTP/2
# > Host: bsc-dataseed.binance.org
# > user-agent: curl/7.77.0
# > accept: */*
# > content-type: application/json
# > content-length: 66

---

curl -v http://localhost:7000 \
  -X POST \
  -H "Host: bsc-dataseed.binance.org" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":2,"method":"web3_clientVersion","params":[]}'


# > Host: bsc-dataseed.binance.org
# > user-agent: curl/7.77.0
# > accept: */*
# > content-type: application/json
# > content-length: 66

TODO generate cert and serve https endpoint
