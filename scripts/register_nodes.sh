#!/bin/sh

curl -X POST \
	-H "Content-Type: application/json" \
	-d '{"name": "mainnet", "httpRPCURL": "https://mainnet.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161", "test": true }' \
	http://localhost:7000/api/v1/nodes/remote

curl -X POST \
	-H "Content-Type: application/json" \
	-d '{"name": "local", "httpRPCURL": "http://localhost:8545", "test": true }' \
	http://localhost:7000/api/v1/nodes/remote

curl -X POST \
	-H "Content-Type: application/json" \
	-d '{"name": "bsc testnet", "httpRPCURL": "https://data-seed-prebsc-1-s1.binance.org:8545", "test": true }' \
	http://localhost:7000/api/v1/nodes/remote


curl -X POST \
	-H "Content-Type: application/json" \
	-d '{"name": "bsc mainnet", "httpRPCURL": "https://bsc-dataseed.binance.org/", "test": true }' \
	http://localhost:7000/api/v1/nodes/remote

curl -X POST \
	-H "Content-Type: application/json" \
	-d '{"name": "local-geth", "httpRPCURL": "http://localhost:8545", "websocketRPCURL": "ws://localhost:8546", "test": true }' \
	http://localhost:7000/api/v1/nodes/remote

curl -X POST \
	-H "Content-Type: application/json" \
	-d '{"name": "avalanche mainnet", "httpRPCURL": "https://api.avax.network/ext/bc/C/rpc", "test": true }' \
	http://localhost:7000/api/v1/nodes/remote

curl -X POST \
	-H "Content-Type: application/json" \
	-d '{"name": "fantom mainnet", "httpRPCURL": "https://rpc.ftm.tools/", "test": true }' \
	http://localhost:7000/api/v1/nodes/remote

curl -X POST \
	-H "Content-Type: application/json" \
	-d '{"name": "polygon mainnet", "httpRPCURL": "https://polygon-rpc.com/", "test": true }' \
	http://localhost:7000/api/v1/nodes/remote

curl -X POST \
	-H "Content-Type: application/json" \
	-d '{"name": "harmony mainnet", "httpRPCURL": "https://api.harmony.one", "test": true }' \
	http://localhost:7000/api/v1/nodes/remote

