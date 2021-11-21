<script lang="ts">
	import { fetchData } from '../stores/http'
	import { nodesURL } from '../lib/const'
	import type { NodesListResponse, NodeResponse } from '../types'
	import ZethNode from '../features/ZethNode/ZethNode.svelte'

	const [nodes] = fetchData<NodesListResponse>(nodesURL)

	let nodeList: NodeResponse[] = []
	$: {
		if ($nodes.data.response?.nodes) {
			nodeList = $nodes.data.response.nodes.sort((a, b) => +new Date(a.dateAdded) - +new Date(b.dateAdded))
		}
	}
</script>

<section>
	<h1>Dashboard</h1>

	{#each nodeList as nodeResponse}
		<ZethNode bind:nodeResponse />
	{/each}
</section>
