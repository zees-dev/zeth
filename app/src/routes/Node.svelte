<script lang="ts">
	import { ethers } from 'ethers'
	import { nodeStore } from '../stores/Node'
	import { settingsStore } from '../stores/Settings'
	import {
		Node,
		getNodeSyncStatus,
		getNetworkName,
		getVersion,
		getSortedModules,
		dateWithoutTZ
	} from '../lib/Models/Node'
	import { NodeType, httpNodeRPCURL, wsNodeRPCURL } from '../lib/const'
	import { fetchData } from '../stores/http'
	import { nodesURL } from '../lib/const'
	import type { NodeResponse } from '../types'
	import { onDestroy } from 'svelte'
	import SyncIndicator from '../components/SyncIndicator.svelte'
	import BlockSyncBar from '../components/BlockSyncBar.svelte'

	export let id: string
	const [nodeData] = fetchData<NodeResponse>(nodesURL + '/' + id)

	// TODO: make this a store
	let rpcEvents: EventSource
	let requestCount = 0
	let responseCount = 0
	let rpcEventData: any[] = []

	rpcEvents = new EventSource(httpNodeRPCURL(id) + '/sse')
	rpcEvents.onmessage = (event) => {
		const data = JSON.parse(event.data)

		// create new temp variable (for list re-assignment - to update list in svelte)
		let updatedData = rpcEventData
		if (data.response.statusCode) {
			// find matching request, remove it from the list, add object with response in its place
			const index = updatedData.findIndex((rpcEvent) => rpcEvent.id === data.id)
			updatedData[index] = data
			responseCount++
		} else {
			updatedData.push(data)
			requestCount++
		}
		// re-assign var to update list in svelte
		rpcEventData = updatedData

		console.log(data)
	}

	// TODO: use global nodeStore if available

	let node: Node
	let coinbase: string | undefined
	const unsubscribe = nodeData.subscribe(async (res) => {
		if (res.data.response) {
			const nodeResponse = res.data.response
			if (nodeResponse) {
				node = new Node(nodeResponse)
				node.setHTTPProvider(new ethers.providers.JsonRpcProvider(httpNodeRPCURL(node.id)))
				node.isDefault = $settingsStore.nodeSettings.defaultNodeID === node.id

				if (NodeType.RemoteNode && !node.connected) {
					node = await node.getRPCData()
				}

				node
					.getCoinbase()
					.then((cb) => {
						coinbase = cb
					})
					.catch(() => {})
			}
		}
	})

	onDestroy(() => {
		unsubscribe
		rpcEvents.close()
	})
</script>

{#if node}
	<figure class="!w-24 m-auto w-">
		<SyncIndicator class="m-2" status={getNodeSyncStatus(node)} />
		<!-- <picture>
			<img src="images/geth-mascot.png" width="72" alt="" class="m-auto" />
		</picture> -->
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
					Version: {node.version || '-'}
				</p>
			</div>
			<div class="text-sm flex span-column my-1">
				<p>Explorer URL: {node.explorerUrl || '-'}</p>
			</div>
			<div class="text-sm span-column mb-1">
				<div class="whitespace-nowrap">Proxy URLs:</div>
				<div class="text-sm">
					{#if node.rpc.http}
						<p class="badge first:ml-1 h-6 bg-red-400"><code>{httpNodeRPCURL(node.id)}</code></p>
					{/if}
					{#if node.rpc.ws}
						<p class="badge first:ml-1 h-6 bg-red-400">{wsNodeRPCURL(node.id)}</p>
					{/if}
				</div>
			</div>
			<div class="text-sm span-column mb-1">
				<div class="whitespace-nowrap ">RPC URLs:</div>
				<div class="text-sm">
					{#if node.rpc.http}
						<p class="badge first:ml-1 h-6 bg-red-400"><code>{node.rpc.http}</code></p>
					{/if}
					{#if node.rpc.ws}
						<p class="badge first:ml-1 h-6 bg-red-400">{node.rpc.ws}</p>
					{/if}
				</div>
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
		<div class="text-sm span-column">
			<p class="whitespace-nowrap">Accounts:</p>
			<div class="flex flex-wrap gap-2">
				{#await node.httpProvider?.listAccounts()}
					<p>...waiting</p>
				{:then accounts}
					{#if accounts?.length}
						{#each accounts as account}
							<p class="badge first:ml-1 h-6 bg-yellow-400">
								{account}<sup>{coinbase?.toLowerCase() === account.toLowerCase() ? ' (coinbase)' : ''}</sup>
							</p>
						{/each}
					{:else}
						<p class="badge ml-1 h-6 badge-warning">none</p>
					{/if}
				{:catch error}
					<p style="color: red">{error.message}</p>
				{/await}
			</div>
		</div>
	</div>
	<div class="text-sm span-column">
		<h2 class="whitespace-nowrap">Stats:</h2>
		<div class="flex flex-wrap gap-2">
			<div class="text-sm span-column">
				<p class="whitespace-nowrap">Logs:</p>
				<div class="flex flex-wrap gap-2">
					<p>TODO</p>
				</div>
			</div>
		</div>
		<div class="text-sm span-column">
			<p class="whitespace-nowrap">Metrics:</p>
			<div class="flex flex-wrap gap-2">
				<p>TODO</p>
			</div>
		</div>
	</div>
	<div class="text-sm span-column">
		<h2 class="whitespace-nowrap">RPC Log ({responseCount}/{requestCount})</h2>
		<div class="flex flex-wrap gap-2 overflow-x-scroll">
			{#each rpcEventData as event}
				<p>{JSON.stringify(event)}</p>
			{/each}
		</div>
	</div>
{:else}
	<h1>Node not found</h1>
{/if}

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
