<script lang="ts">
	import { fetchData, fetchStates } from '../stores/http'
	import { nodesURL, httpNodeRPCURL } from '../lib/const'
	import type { NodesResponse, Node } from '../lib/types'
	import { ethers } from 'ethers'

	interface ProviderNode extends Node {
		network: ethers.providers.Network
		block: number
		version: string
		syncing: boolean
		peers: number
		modules: {
			admin: string
			debug: string
			eth: string
			net: string
			personal: string
			rpc: string
			web3: string
		}
		mining: boolean
	}

	const [nodes] = fetchData<NodesResponse>(nodesURL)

	async function getNodeDataFromRPC(nodes: Node[]) {
		const nodeList = await Promise.all(
			nodes.map(async (n) => {
				const provider = new ethers.providers.JsonRpcProvider(httpNodeRPCURL(n.id))
				const [networkP, blockP, versionP, syncingP, peersP, modulesP, miningP] = await Promise.allSettled([
					provider.getNetwork(),
					provider.getBlockNumber(),
					provider.send('web3_clientVersion', []),
					provider.send('eth_syncing', []),
					provider.send('net_peerCount', []),
					provider.send('rpc_modules', []),
					provider.send('eth_mining', [])
				])

				const pNode = n as ProviderNode
				if (networkP.status === 'fulfilled') {
					pNode.network = networkP.value
				}
				if (blockP.status === 'fulfilled') {
					pNode.block = blockP.value
				}
				if (versionP.status === 'fulfilled') {
					pNode.version = versionP.value
				}
				if (syncingP.status === 'fulfilled') {
					pNode.syncing = syncingP.value
				}
				if (peersP.status === 'fulfilled') {
					pNode.peers = parseInt(peersP.value, 16)
				}
				if (modulesP.status === 'fulfilled') {
					pNode.modules = modulesP.value
				}
				if (miningP.status === 'fulfilled') {
					pNode.mining = miningP.value
				}

				return pNode
			})
		)
		return nodeList
	}

	// $: {
	// 	if ($nodes.data.response?.nodes) {
	// 		getNodeDataFromRPC($nodes.data.response?.nodes)
	// 	}
	// }

	// - node status entry
	// - type: geth
	// - version (geth version)
	// - name
	// - date added
	// - status (on/off)
	// - chain (mainnet, bsc, poly)
	// - testnet?
	// - blocknumber, is syncing?
	// - IPC, HTTP, Websocket enabled
	// - RPC modules
	// - is miner/mining
	// - is dev mode
	// - connected peers

	// // {"jsonrpc":"2.0","id":2,"method":"web3_clientVersion","params":[]}
	// export const type: string = 'geth'
	// export const version: string = 'v1.0.0'
	// export const name: string = 'Ethereum'

	// // {"jsonrpc":"2.0","id":1,"method":"net_version"}
	// export const networkID: number = 1
	// export const running: boolean = true
	// export const dateAdded: string = '2018-01-01'

	// // {"jsonrpc":"2.0","id":3,"method":"eth_blockNumber","params":[]}
	// // {"jsonrpc":"2.0","method":"eth_syncing","params":[],"id":1}
	// export const isSyncing: boolean = false // isValidating

	// // {"jsonrpc":"2.0","method":"net_peerCount","params":[],"id":1}
	// export const peers: number = 0

	// // {"jsonrpc":"2.0","id":1,"method":"rpc_modules"}
	// export const enabledModules: string[] = []

	// // {"jsonrpc":"2.0","id":1,"method":"eth_mining"}
	// export const isMining: boolean = false

	// export const isDev: boolean = false
</script>

{#if $nodes.data.response?.nodes}
	{#await getNodeDataFromRPC($nodes.data.response?.nodes)}
		<p>...waiting</p>
	{:then nodes}
		{#each nodes as node}
			<p>node block {node.block}</p>
		{/each}
	{:catch error}
		<p style="color: red">{error.message}</p>
	{/await}
{/if}

<style>
</style>
