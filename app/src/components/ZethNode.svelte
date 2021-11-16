<script lang="ts">
	import { ethers } from 'ethers'
	import { navigate } from 'svelte-routing'
	import { settingsStore } from '../stores'
	import { fetchData } from '../stores/http'
	import { nodesURL, httpNodeRPCURL, NodeType, EthNetworks, BSCNetworks, AvaNetworks, FTMNetworks } from '../lib/const'
	import type { NodesResponse, Node, SyncStatus } from '../lib/types'
	import SyncIndicator from './SyncIndicator.svelte'
	import BlockSyncBar from './BlockSyncBar.svelte'

	interface ProviderNode extends Node {
		connected: boolean
		network: ethers.providers.Network
		block: number
		version: string
		syncing: SyncStatus
		peers: number
		modules: [string, string][]
		mining: boolean
		isDefault: boolean
	}

	const [nodes] = fetchData<NodesResponse>(nodesURL)

	async function getNodeDataFromRPC(nodes: Node[]): Promise<ProviderNode[]> {
		const enabledNodeList = await Promise.all(
			nodes
				.filter((n) => n.enabled)
				.map(async (n) => {
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
						pNode.connected = true
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

					pNode.isDefault = $settingsStore.nodeSettings.defaultNodeID === pNode.id
					return pNode
				})
		)
		const disabledNodeList = nodes.filter((n) => !n.enabled) as ProviderNode[]
		return [...enabledNodeList, ...disabledNodeList]
	}

	function getNetworkName(chainId: number) {
		// TODO: account for other non-eth chains
		return EthNetworks[chainId] ?? BSCNetworks[chainId] ?? AvaNetworks[chainId] ?? FTMNetworks[chainId]
	}

	function getNodeSyncStatus(node: ProviderNode) {
		debugger
		if (node.syncing) {
			return 'syncing'
		}
		if (node.syncing === undefined) {
			return undefined
		}
		return 'synced'
	}

	/**
	 * getVersion returns the geth version from the version string.
	 * It gets the string between first '/' and second '/' characters.
	 * @param gethVersion example: "Geth/v1.10.9-omnibus-e03773e6/linux-amd64/go1.17.2"
	 * @returns example: "v1.10.9-omnibus-e03773e6"
	 */
	function getVersion(gethVersion: string) {
		return gethVersion.split('/')[1] ?? gethVersion
	}

	function dateWithoutTZ(date: Date) {
		return date.toString().substring(0, date.toString().indexOf('GMT') - 1)
	}

	function getSortedModules(modules: [string, string][]) {
		return Object.keys(modules)
			.sort()
			.map((key: string) => ({ module: key, version: modules[key as any] }))
	}

	let nodeList: ProviderNode[] = []
	$: {
		if ($nodes.data.response?.nodes) {
			nodeList = $nodes.data.response.nodes.sort(
				(a, b) => +new Date(a.dateAdded) - +new Date(b.dateAdded)
			) as ProviderNode[]
			getNodeDataFromRPC($nodes.data.response?.nodes ?? []).then((nodes) => {
				nodeList = nodes.sort((a, b) => +new Date(a.dateAdded) - +new Date(b.dateAdded))
			})
		}
	}
</script>

<!-- {#if $nodes.data.response?.nodes} -->

{#each nodeList as node}
	<div
		class="card card-side bordered bg-neutral hover:bg-accent m-2 cursor-pointer"
		on:click={() => navigate(`/node/${node.id}`)}
	>
		<figure class="!w-24 m-auto w-">
			<SyncIndicator class="absolute m-2" status={getNodeSyncStatus(node)} />
			<picture>
				<img src="images/geth-mascot.png" width="72" alt="" class="m-auto" />
			</picture>
		</figure>
		<div class="card-body p-0">
			<h2 class="card-title pt-2 pl-2 mb-0 card-title-grid">
				{node.name}
				<div class="badge badge-sm mx-2">{NodeType[node.nodeType]}</div>
				{#if node.isDefault}
					<div class="badge badge-md badge-info mx-2 justify-self-end">Default</div>
				{/if}
			</h2>

			<div class="pl-6 pr-2 pb card-content-grid">
				<div class="text-sm flex place-items-center span-column">
					<p>Network ID: {node.network?.chainId ?? '-'}</p>
					{#if getNetworkName(node.network?.chainId)}
						<p class="badge badge-sm badge-info mx-2">{getNetworkName(node.network.chainId)}</p>
					{/if}
				</div>
				<div data-tip={node.version} class="tooltip flex span-column">
					<p class="text-sm">
						Version: {(node.version && getVersion(node.version)) || '-'}
					</p>
				</div>
				<div class="text-sm flex span-column">
					<p class="whitespace-nowrap">RPC modules:</p>
					<div class="flex flex-wrap gap-2">
						{#if node.modules}
							{#each getSortedModules(node.modules) as { module, version }}
								<p class="badge first:ml-1 h-6 bg-yellow-400">{module}<sup>{version}</sup></p>
							{/each}
						{:else}
							<p class="badge ml-1 h-6 badge-warning">unknown</p>
						{/if}
					</div>
				</div>
				<p class="text-sm">Peers: {node.peers ?? '-'}</p>
				<p class="text-sm justify-self-end">Date added: {dateWithoutTZ(new Date(node.dateAdded)) ?? '-'}</p>
			</div>
			<BlockSyncBar blockNumber={node.block} syncing={node.syncing} />
		</div>
	</div>
{/each}

<!-- {/if} -->
<style>
	.span-column {
		grid-column: 1/ -1;
	}
	.card-title-grid {
		display: grid;
		grid-template-columns: auto auto 1fr;
	}

	.card-content-grid {
		display: grid;
		grid-template: 1fr 1fr / auto auto;
	}
</style>
