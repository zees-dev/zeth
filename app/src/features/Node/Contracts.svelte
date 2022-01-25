<script lang="ts">
	import Alert from '../../components/Alert.svelte'
	import { nodeStore } from '../../stores/Node'

	let viewType = 0
	let pairAddress = '',
		token1Address = '',
		token2Address = ''

	let token1 = '',
		token2 = ''

	let swapEnabled = false
	let isSwapping = false
	let swapError = ''
</script>

<div>WELCOME {$nodeStore.name}</div>

<div class="max-w-2xl px-8 py-4 m-2 bg-white rounded-lg shadow-md dark:bg-gray-800">
	<div class="flex items-center justify-between">
		<div class="text-sm font-light text-gray-600 dark:text-gray-400">
			<select class="select select-bordered w-full max-w-xs">
				<option disabled selected>Automated Market Maker</option>
				<option>TraderJoe</option>
				<option>Pangolin</option>
			</select>
		</div>
		<div class="tabs tabs-boxed">
			<a class="tab {viewType === 0 ? 'tab-active' : ''}" on:click={() => (viewType = 0)}>Use Pair</a>
			<a class="tab {viewType === 1 ? 'tab-active' : ''}" on:click={() => (viewType = 1)}>Use Tokens</a>
		</div>
	</div>

	<form class="flex gap-2 flex-col mt-2 p-4" on:submit|preventDefault={() => {}}>
		{#if viewType === 0}
			<div class="container">
				<label>Pair address</label>
				<input type="text" placeholder="0x0000000000000000000000000000000000000000" bind:value={pairAddress} />
			</div>
		{/if}

		{#if viewType === 1}
			<label for="token-1" class="btn btn-primary modal-button w-32">{token1 || 'select token 1'}</label>
			<input type="checkbox" id="token-1" class="modal-toggle" />
			<div class="modal">
				<div class="modal-box">
					<div class="container">
						<label>Token 0 address</label>
						<input type="text" placeholder="0x0000000000000000000000000000000000000000" bind:value={token1Address} />
					</div>
					<div class="modal-action">
						<label for="token-1" class="btn btn-primary">Use</label>
						<label for="token-1" class="btn">Close</label>
					</div>
				</div>
			</div>

			<label for="token-2" class="btn btn-primary modal-button w-32">{token2 || 'select token 2'}</label>
			<input type="checkbox" id="token-2" class="modal-toggle" />
			<div class="modal">
				<div class="modal-box">
					<div class="container">
						<label>Token 1 address</label>
						<input type="text" placeholder="0x0000000000000000000000000000000000000000" bind:value={token2Address} />
					</div>
					<div class="modal-action">
						<label for="token-2" class="btn btn-primary">Use</label>
						<label for="token-2" class="btn">Close</label>
					</div>
				</div>
			</div>
		{/if}

		{#if token1 && token2}
			<div class="container">
				<label>{token1} amount</label>
				<input type="text" placeholder="0.0" bind:value={pairAddress} />
			</div>

			->

			<div class="container">
				<label>{token1} amount</label>
				<input type="text" placeholder="0.0" bind:value={pairAddress} />
			</div>

			<button
				class="col-span-2 btn btn-secondary {isSwapping ? 'loading' : ''} disabled"
				type="submit"
				disabled={swapEnabled}
			>
				SWAP
			</button>
		{/if}

		{#if swapError}
			<Alert type="error" message={swapError} />
		{/if}
	</form>
</div>

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
