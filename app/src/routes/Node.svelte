<script lang="ts">
	import { ethers } from 'ethers'
	import { nodeStore } from '../stores/Node'
	import { settingsStore } from '../stores/Settings'
	import { Node } from '../lib/Models/Node'
	import { NodeType, httpNodeRPCURL } from '../lib/const'
	import { fetchData } from '../stores/http'
	import { nodesURL } from '../lib/const'
	import type { NodeResponse } from '../types'
	import { onDestroy } from 'svelte'

	export let id: string
	const [nodeData] = fetchData<NodeResponse>(nodesURL + '/' + id)

	// TODO: use global nodeStore if available

	let node: Node
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
			}
		}
	})

	onDestroy(unsubscribe)
</script>

<section>
	<h1>Node: {id}</h1>
</section>
