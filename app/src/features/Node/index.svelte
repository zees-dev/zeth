<script lang="ts">
	import { nodeStore } from '../../stores/Node'
	import { Node, getNodeSyncStatus, getNetworkName, getSortedModules, dateWithoutTZ } from '../../lib/Models/Node'
	import { httpNodeRPCURL, wsNodeRPCURL } from '../../lib/const'
	import { onDestroy } from 'svelte'
	import SyncIndicator from '../../components/SyncIndicator.svelte'
	import GasPrice from '../../components/GasPrice.svelte'
	import BlockSyncBar from '../../components/BlockSyncBar.svelte'
	import { navigate } from 'svelte-navigator'
	import RPCLog from './RPCLog.svelte'

	interface RPCEvent {
		id: string
		request: {
			headers: { [key: string]: string | number | boolean | Array<any> }
			body: {
				// [key: string]: any,
				id: number
				jsonrpc: string
				method: string
				params: Array<any>
			}
		}
		response: {
			headers?: { [key: string]: string | number | boolean | Array<any> }
			body?: {
				// [key: string]: any
				id: number
				jsonrpc: string
				result?: any
				error?: { code: number; message: string }
			}
			statusCode: number
		}
		duration: number
	}

	const node: Node = $nodeStore

	// TODO: make this a store
	let rpcEvents: EventSource
	let requestCount = 0
	let responseCount = 0
	let rpcEventData: RPCEvent[] = []

	rpcEvents = new EventSource(httpNodeRPCURL(node.id) + '/sse')
	rpcEvents.onmessage = (event) => {
		const data = JSON.parse(event.data)

		// JSONify request headers and body
		data.request.headers = JSON.parse(data.request.headers)
		data.request.body = JSON.parse(data.request.body)

		// create new temp variable (for list re-assignment - to update list in svelte)
		let updatedData = rpcEventData
		if (data.response.statusCode) {
			// JSONify response headers and body
			if (data.response.headers) {
				data.response.headers = JSON.parse(data.response.headers)
			}
			try {
				if (data.response.body) {
					data.response.body = JSON.parse(data.response.body)
				}
			} catch (e) {
				// console.log(data)
				console.log('bad...')
				data.response.body = { result: 'error...' }
			}
			console.log(data)

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
	}

	// TODO: use global nodeStore if available

	let coinbase: string | undefined
	node
		.getCoinbase()
		.then((cb) => {
			coinbase = cb
		})
		.catch(() => {})

	onDestroy(() => {
		rpcEvents.close()
	})

	let isRemoving = false
	async function removeNode() {
		if (isRemoving) return
		isRemoving = true
		try {
			await fetch(`/api/v1/nodes/${node.id}`, {
				method: 'DELETE',
			})
			nodeStore.remove()
			navigate('/')
		} catch (e) {
			console.log(e)
		} finally {
			isRemoving = false
		}
	}
</script>

<figure class="!w-24 m-auto w-">
	<SyncIndicator class="m-2" status={getNodeSyncStatus(node)} />
	<!-- <picture>
		<img src="images/geth-mascot.png" width="72" alt="" class="m-auto" />
	</picture> -->
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

<GasPrice {node} />

<div class="w-96">
	<label for="my-modal-2" class="btn btn-error modal-button w-full">Remove</label>
	<input type="checkbox" id="my-modal-2" class="modal-toggle" />
	<div class="modal">
		<div class="modal-box">
			<p>
				Enim dolorem dolorum omnis atque necessitatibus. Consequatur aut adipisci qui iusto illo eaque. Consequatur
				repudiandae et. Nulla ea quasi eligendi. Saepe velit autem minima.
			</p>
			<div class="modal-action">
				<label for="my-modal-2" class="btn btn-error btn-primary" on:click|preventDefault={removeNode}>Yes</label>
				<label for="my-modal-2" class="btn btn-info">Cancel</label>
			</div>
		</div>
	</div>
</div>

<RPCLog class="my-8" />

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
