<script lang="ts">
	import { fetchData } from '../stores/http'
	import { nodesURL } from '../lib/const'
	import type { NodesResponse, Node } from '../types'
	import ZethNode from '../features/ZethNode.svelte'

	const [nodes] = fetchData<NodesResponse>(nodesURL)

	let nodeList: Node[] = []
	$: {
		if ($nodes.data.response?.nodes) {
			nodeList = $nodes.data.response.nodes.sort((a, b) => +new Date(a.dateAdded) - +new Date(b.dateAdded))
		}
	}
</script>

<section>
	<h1>Dashboard</h1>

	{#each nodeList as node}
		<ZethNode bind:node />
	{/each}
</section>
