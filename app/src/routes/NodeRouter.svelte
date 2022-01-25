<script lang="ts">
	import { nodeStore } from '../stores/Node'
	import { settingsStore } from '../stores/Settings'
	import { Node } from '../lib/Models/Node'
	import { fetchStates, fetchData } from '../stores/http'
	import { nodesURL } from '../lib/const'
	import type { NodeResponse } from '../types'
	import { onDestroy } from 'svelte'
	import { Route } from 'svelte-navigator'
	import Contracts from '../features/Node/Contracts.svelte'
	import NodeComponent from '../features/Node/index.svelte'

	export let id: string
	const [nodeData] = fetchData<NodeResponse>(nodesURL + '/' + id)

	let node: Node
	let status: string
	const unsubscribe = nodeData.subscribe(async (res) => {
		status = res.status
		if (res.data.response) {
			const nodeResponse = res.data.response
			if (nodeResponse) {
				node = new Node(nodeResponse)

				if (!node.connected) {
					node = await node.getRPCData()
				}

				nodeStore.set(node)
			}
		}
	})

	onDestroy(() => {
		unsubscribe
	})
</script>

{#if status === fetchStates.LOADING}
	<h1>loading...</h1>
{:else if status === fetchStates.SUCCESS && $nodeStore?.id}
	<Route path="/" let:params>
		<NodeComponent />
	</Route>
	<Route path="/contracts" let:params>
		<Contracts />
	</Route>
{:else if status === fetchStates.ERROR}
	<h1>Node not found</h1>
{:else}
	<h1>Invalid state: {status}</h1>
{/if}
