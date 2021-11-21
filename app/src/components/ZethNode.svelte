<script lang="ts">
	import { ethers } from 'ethers'
	import { navigate } from 'svelte-routing'
	import { settingsStore } from '../stores'
	import { httpNodeRPCURL, NodeType, EthNetworks } from '../lib/const'
	import type { Node, SyncStatus } from '../lib/types'
	import SyncIndicator from './SyncIndicator.svelte'
	import BlockSyncBar from './BlockSyncBar.svelte'

	export let node: Node

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

	async function getNodeDataFromRPC(node: ProviderNode): Promise<ProviderNode> {
		const provider = new ethers.providers.JsonRpcProvider(httpNodeRPCURL(node.id))

		const [networkP, blockP, versionP, syncingP, peersP, modulesP, miningP] = await Promise.allSettled([
			provider.getNetwork(),
			provider.getBlockNumber(),
			provider.send('web3_clientVersion', []),
			provider.send('eth_syncing', []),
			provider.send('net_peerCount', []),
			provider.send('rpc_modules', []),
			provider.send('eth_mining', [])
		])

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
	}

	function getNetworkName(chainId: number) {
		// TODO: account for other non-eth chains
		return EthNetworks[chainId]
	}

	function getNodeSyncStatus(node: ProviderNode) {
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

	let pNode = node as ProviderNode
	if (node.enabled) {
		getNodeDataFromRPC(pNode).then((responseNode) => {
			pNode = responseNode
		})
	}
</script>

<div
	class="card card-side bordered bg-neutral hover:bg-accent m-2 cursor-pointer"
	on:click={() => navigate(`/node/${pNode.id}`)}
>
	<figure class="!w-24 m-auto w-">
		<SyncIndicator class="absolute m-2" status={getNodeSyncStatus(pNode)} />
		<picture>
			<img src="images/geth-mascot.png" width="72" alt="" class="m-auto" />
		</picture>
	</figure>
	<div class="card-body p-0">
		<h2 class="card-title pt-2 pl-2 mb-0 card-title-grid">
			{pNode.name}
			<div class="badge badge-sm mx-2">{NodeType[pNode.nodeType]}</div>
			{#if pNode.isDefault}
				<div class="badge badge-md badge-info mx-2 justify-self-end">Default</div>
			{/if}
		</h2>

		<div class="pl-6 pr-2 pb card-content-grid">
			<div class="text-sm flex place-items-center span-column">
				<p>Network ID: {pNode.network?.chainId ?? '-'}</p>
				{#if getNetworkName(pNode.network?.chainId)}
					<p class="badge badge-sm badge-info mx-2">{getNetworkName(pNode.network.chainId)}</p>
				{/if}
			</div>
			<div data-tip={pNode.version} class="tooltip flex span-column">
				<p class="text-sm">
					Version: {(pNode.version && getVersion(pNode.version)) || '-'}
				</p>
			</div>
			<div class="text-sm flex span-column">
				<p class="whitespace-nowrap">RPC modules:</p>
				<div class="flex flex-wrap gap-2">
					{#if pNode.modules}
						{#each getSortedModules(pNode.modules) as { module, version }}
							<p class="badge first:ml-1 h-6 bg-yellow-400">{module}<sup>{version}</sup></p>
						{/each}
					{:else}
						<p class="badge ml-1 h-6 badge-warning">unknown</p>
					{/if}
				</div>
			</div>
			<p class="text-sm">Peers: {pNode.peers ?? '-'}</p>
			<p class="text-sm justify-self-end">Date added: {dateWithoutTZ(new Date(pNode.dateAdded)) ?? '-'}</p>
		</div>
		<BlockSyncBar blockNumber={pNode.block} syncing={pNode.syncing} />
	</div>
</div>

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
