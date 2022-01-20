<script lang="ts">
	import { navigate } from 'svelte-navigator'
	import { settingsStore } from '../../stores/Settings'
	import { nodeStore } from '../../stores/Node'
	import type { NodeResponse } from '../../types'
	import {
		Node,
		getNodeSyncStatus,
		getNetworkName,
		getVersion,
		getSortedModules,
		dateWithoutTZ,
	} from '../../lib/Models/Node'
	import SyncIndicator from '../../components/SyncIndicator.svelte'
	import BlockSyncBar from '../../components/BlockSyncBar.svelte'

	export let nodeResponse: NodeResponse

	// handleClick globally sets node and navigates to node view
	function handleClick(node: Node) {
		nodeStore.set(node)
		navigate(`/node/${node.id}`)
	}

	let node = new Node(nodeResponse)
	node.isDefault = $settingsStore.nodeSettings.defaultNodeID === node.id

	if (node.enabled && !node.connected) {
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
