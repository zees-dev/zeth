#!/bin/sh

curl -X POST \
	-H "Content-Type: application/json" \
	-d '{"name": "mainnet", "rpc": { "http": "https://mainnet.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161" }, "test": true }' \
	http://localhost:7000/api/v1/nodes/remote

curl -X POST \
	-H "Content-Type: application/json" \
	-d '{"name": "local", "rpc": { "http": "http://localhost:8545" }, "test": true }' \
	http://localhost:7000/api/v1/nodes/remote

curl -X POST \
	-H "Content-Type: application/json" \
	-d '{"name": "bsc testnet", "rpc": { "http": "https://data-seed-prebsc-1-s1.binance.org:8545" }, "test": true }' \
	http://localhost:7000/api/v1/nodes/remote

curl -X POST \
	-H "Content-Type: application/json" \
	-d '{"name": "bsc mainnet", "rpc": { "http": "https://bsc-dataseed.binance.org/" }, "test": true }' \
	http://localhost:7000/api/v1/nodes/remote

curl -X POST \
	-H "Content-Type: application/json" \
	-d '{"name": "local-geth", "rpc": { "http": "http://localhost:8545", "ws": "ws://localhost:8546" }, "test": true }' \
	http://localhost:7000/api/v1/nodes/remote

curl -X POST \
	-H "Content-Type: application/json" \
	-d '{"name": "avalanche mainnet", "rpc": { "http": "https://api.avax.network/ext/bc/C/rpc", "ws": "wss://api.avax.network/ext/bc/C/ws" }, "test": true }' \
	http://localhost:7000/api/v1/nodes/remote

curl -X POST \
	-H "Content-Type: application/json" \
	-d '{"name": "fantom mainnet", "rpc": { "http": "https://rpc.ftm.tools/" }, "test": true }' \
	http://localhost:7000/api/v1/nodes/remote

curl -X POST \
	-H "Content-Type: application/json" \
	-d '{"name": "polygon mainnet", "rpc": { "http": "https://polygon-rpc.com/" }, "test": true }' \
	http://localhost:7000/api/v1/nodes/remote

curl -X POST \
	-H "Content-Type: application/json" \
	-d '{"name": "harmony mainnet", "rpc": { "http": "https://api.harmony.one" }, "test": true }' \
	http://localhost:7000/api/v1/nodes/remote

curl -X POST \
	-H "Content-Type: application/json" \
	-d '{"name": "Arbitrum One", "rpc": { "http": "https://arb1.arbitrum.io/rpc" }, "test": true }' \
	http://localhost:7000/api/v1/nodes/remote

curl -X POST \
	-H "Content-Type: application/json" \
	-d '{"name": "aurora mainnet", "rpc": { "http": "https://mainnet.aurora.dev" }, "test": true }' \
	http://localhost:7000/api/v1/nodes/remote
