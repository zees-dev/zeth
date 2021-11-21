<script lang="ts">
	import { ethers } from 'ethers'
	import { navigate } from 'svelte-routing'
	import { settingsStore } from '../../stores/Settings'
	import { nodeStore } from '../../stores/Node'
	import { httpNodeRPCURL, NodeType, EthNetworks } from '../../lib/const'
	import type { NodeResponse, RPCModules } from '../../types'
	import { Node } from '../../lib/Models/Node'
	import SyncIndicator from '../../components/SyncIndicator.svelte'
	import BlockSyncBar from '../../components/BlockSyncBar.svelte'

	export let nodeResponse: NodeResponse

	// handleClick globally sets node and navigates to node view
	function handleClick(node: Node) {
		nodeStore.set(node)
		navigate(`/node/${node.id}`)
	}

	function getNetworkName(chainId: number) {
		// TODO: account for other non-eth chains
		return EthNetworks[chainId]
	}

	function getNodeSyncStatus(node: Node) {
		if (!node.connected) {
			return undefined
		}
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

	function getSortedModules(modules: RPCModules) {
		return Object.keys(modules)
			.sort()
			.map((key: string) => ({ module: key, version: modules[key as any] }))
	}

	let node = new Node(nodeResponse)
	node.setHTTPProvider(new ethers.providers.JsonRpcProvider(httpNodeRPCURL(node.id)))
	node.isDefault = $settingsStore.nodeSettings.defaultNodeID === node.id

	if (node.enabled && node.nodeType == NodeType.RemoteNode && !node.connected) {
		node
			.getRPCData()
			.then((n) => (node = n))
			.catch((err) => {
				// TODO: snackbar/notification
				console.log(err)
			})
	}
</script>

<div class="card card-side bordered bg-neutral hover:bg-accent m-2 cursor-pointer" on:click={() => handleClick(node)}>
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
				{#if node.network}
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
