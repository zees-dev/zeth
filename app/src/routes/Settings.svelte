<script lang="ts">
	let nodeName = ''
	let nodeHTTPRPC = ''
	let nodeWSRPC = ''
	let nodeTestHTTPRPC = true
	let nodeExplorerURL = ''

	let registeringNode = false
	let registerNodeError = ''
	async function registerNode() {
		if (registeringNode || disabled) return
		registeringNode = true
		registerNodeError = ''
		try {
			const res = await fetch('/api/v1/nodes', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					name: nodeName,
					rpc: {
						http: nodeHTTPRPC,
						ws: nodeWSRPC,
					},
					test: nodeTestHTTPRPC,
					explorer: nodeExplorerURL,
				}),
			})
			if (!res.ok) throw new Error(await res.text())

			const data = await res.json()
			nodeName = data.name
			nodeHTTPRPC = data.rpc.http
			nodeWSRPC = data.rpc.ws
			nodeExplorerURL = data.explorerUrl
		} catch (e) {
			console.error(e)
			registerNodeError = e as string
		} finally {
			registeringNode = false
		}
	}

	$: disabled = !nodeName || !nodeHTTPRPC
</script>

<form class="flex gap-2 flex-col w-96 p-4" on:submit|preventDefault={registerNode}>
	<div class="container">
		<label>Name<sup class="text-red-500">*</sup></label>
		<input type="text" placeholder="Mainnet" bind:value={nodeName} />
	</div>

	<label>RPC</label>

	<div class="container">
		<label>HTTP<sup class="text-red-500">*</sup></label>
		<input type="text" placeholder="https://mainnet.infura.io/v3/abcdef" bind:value={nodeHTTPRPC} />

		<div class="test flex justify-self-end gap-2 place-items-center">
			<label>test connection</label>
			<input type="checkbox" bind:checked={nodeTestHTTPRPC} />
		</div>
	</div>

	<div class="container">
		<label>WS</label>
		<input type="text" placeholder="wss://mainnet.infura.io/ws/v3/abcdef" bind:value={nodeWSRPC} />
	</div>

	<div class="container">
		<label>Explorer</label>
		<input type="text" placeholder="https://etherscan.io" bind:value={nodeExplorerURL} />
	</div>

	<button class="col-span-2 btn btn-secondary {registeringNode ? 'loading' : ''} disabled" type="submit" {disabled}>
		Register
	</button>

	{#if registerNodeError}
		<div class="text-red-500">{registerNodeError}</div>
	{/if}
</form>

<style>
	.container {
		display: grid;
		grid-column: 1/-1;
		grid-template-columns: 1fr 3fr;
	}

	.test {
		grid-column: 2;
	}
</style>
