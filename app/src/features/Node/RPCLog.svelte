<script lang="ts">
	import { httpNodeRPCURL } from '../../lib/const'
	import { nodeStore } from '../../stores/Node'

	const node = $nodeStore

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
</script>

<div class="text-sm span-column w-11/12 {$$props.class}">
	<h2 class="whitespace-nowrap text-center" title="{responseCount} sent, {requestCount} recv">
		RPC Log ({responseCount} ⬆︎/{requestCount} ⬇︎)
	</h2>
	<table class="table w-full table-compact">
		<thead>
			<tr>
				<th>ID</th>
				<th>Request</th>
				<th>Response</th>
				<th>Duration (ms)</th>
				<th>(replay)</th>
				<th>(pin)</th>
				<th>(copy)</th>
			</tr>
		</thead>
		<tbody>
			{#each rpcEventData as event}
				<!-- <p>{JSON.stringify(event)}</p> -->
				<tr>
					<td><code>{event.request.body.id}</code></td>
					<td>
						<code>{event.request.body.method}</code>
						{#if event.request.body.params.length > 0}
							<p>params:</p>
							{#each event.request.body.params as param}
								<code>{JSON.stringify(param)}</code>
							{/each}
						{/if}
					</td>
					<td>
						{#if event.response.body}
							{#if event.response.body.error}
								<code>error: {event.response.body.error.message}</code>
							{/if}
							{#if event.response.body.result}
								<code class="tooltip tooltip-info" data-tip={JSON.stringify(event.response.body.result)}>
									{event.response.body.result}
								</code>
							{/if}
						{:else}
							<p>pending...</p>
						{/if}
					</td>
					<td><code>{event.duration}</code></td>
					<td>todo.</td>
					<td>todo.</td>
					<td>todo.</td>
				</tr>
			{/each}
		</tbody>
	</table>
</div>
