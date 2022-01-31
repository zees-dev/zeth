<script lang="ts">
	import { onMount } from 'svelte'

	import Alert from '../../components/Alert.svelte'
	import { nodeStore } from '../../stores/Node'
	import type { AMM } from '../../types'

	let viewType = 0
	let pairAddress = '',
		token1Address = '',
		token2Address = ''

	let token1 = '',
		token2 = ''

	let swapEnabled = false
	let isSwapping = false
	let swapError = ''
	let selectedIndex: number = -1
	let ammList: AMM[] = []

	let isRetrievingAMMs = false
	async function getAMMs() {
		if (isRetrievingAMMs) return
		isRetrievingAMMs = true
		try {
			const ammsData = await fetch(`/api/v1/defi/${$nodeStore.network?.chainId}/amm`)
			ammList = await ammsData.json()
		} catch (e) {
			console.log(e)
		} finally {
			isRetrievingAMMs = false
		}
	}

	onMount(() => {
		getAMMs()
	})
</script>

<div>WELCOME {$nodeStore.name}</div>

<div class="max-w-2xl px-8 py-4 m-2 bg-white rounded-lg shadow-md dark:bg-gray-800">
	<div class="flex items-center justify-between">
		<div class="text-sm font-light text-gray-600 dark:text-gray-400">
			<select
				class="select select-bordered w-full max-w-xs"
				bind:value={selectedIndex}
				on:change={() => console.log(selectedIndex)}
			>
				<option disabled selected value={-1}>Automated Market Maker {isRetrievingAMMs ? 'loading...' : ''}</option>
				{#each ammList as amm, index}
					<option value={index}>{amm.name}</option>
				{/each}
			</select>
		</div>
		<div class="tabs tabs-boxed">
			<a class="tab {viewType === 0 ? 'tab-active' : ''}" on:click={() => (viewType = 0)}>Pair address</a>
			<a class="tab {viewType === 1 ? 'tab-active' : ''}" on:click={() => (viewType = 1)}>Token addresses</a>
		</div>
	</div>

	<form class="flex gap-2 flex-col mt-2 p-4" on:submit|preventDefault={() => {}}>
		{#if selectedIndex >= 0}
			{#if ammList[selectedIndex].url}
				<a href={ammList[selectedIndex].url} target="_blank" class="text-blue-500 m-auto">
					{ammList[selectedIndex].name}
				</a>
			{:else}
				<div class="font-light text-gray-600 dark:text-gray-400 m-auto">
					{ammList[selectedIndex].name}
				</div>
			{/if}
			<p>
				Router address:
				{#if $nodeStore.explorerUrl}
					<a
						href={`${$nodeStore.explorerUrl}/address/${ammList[selectedIndex].routerAddress}`}
						target="_blank"
						class="text-blue-500"
					>
						{ammList[selectedIndex].routerAddress}
					</a>
				{:else}
					{ammList[selectedIndex].routerAddress}
				{/if}
			</p>
			<p>
				Factory address:
				{#if $nodeStore.explorerUrl}
					<a
						href={`${$nodeStore.explorerUrl}/address/${ammList[selectedIndex].factoryAddress}`}
						target="_blank"
						class="text-blue-500"
					>
						{ammList[selectedIndex].factoryAddress}
					</a>
				{:else}
					{ammList[selectedIndex].factoryAddress}
				{/if}
			</p>
		{/if}

		{#if viewType === 0}
			<div class="container">
				<label>Pair address</label>
				<input type="text" placeholder="0x0000000000000000000000000000000000000000" bind:value={pairAddress} />
			</div>
		{/if}

		{#if viewType === 1}
			<div class="container">
				<label>Token 1</label>
				<input type="text" placeholder="0x0000000000000000000000000000000000000000" bind:value={token1Address} />
			</div>

			<div class="container">
				<label>Token 2</label>
				<input type="text" placeholder="0x0000000000000000000000000000000000000000" bind:value={token2Address} />
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
				class="col-span-2 btn btn-secondary {isSwapping ? 'loading' : ''} {!swapEnabled && !isSwapping
					? 'disabled'
					: ''}"
				type="submit"
				disabled={!swapEnabled && !isSwapping}
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
</style>
