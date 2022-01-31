<script lang="ts">
	import { onDestroy, onMount } from 'svelte'

	import Alert from '../../components/Alert.svelte'
	import { nodeStore } from '../../stores/Node'
	import type { AMM } from '../../types'
	import { ERC20 } from '../../../typechain/abi/ERC20'
	import ERC20Abi from '../../../abi/ERC20.json'
	import * as UniswapV2PairABI from '../../../abi/uniswap/UniswapV2Pair.json'
	import { debounce, getContract } from '../../utils'
	import { ethers, logger } from 'ethers'
	import { UniswapV2Pair } from '../../../typechain/abi'
	import { floatToBigNumber } from '../../utils/etherutils'

	interface Token {
		address: string
		symbol: string
		name: string
		decimals: number
		reserve: ethers.BigNumber
	}

	onMount(() => {
		getAMMs()
	})

	let viewType = 0
	let pairAddress = '',
		token0Address = '',
		token1Address = ''

	let token0: Token | undefined, token1: Token | undefined
	let token0Amount: number = 0.0,
		token1Amount: number = 0.0

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

	let pairError = ''
	let pairLoading = false
	const getPairData = debounce(async (pairAddress: string) => {
		if (!ethers.utils.isAddress(pairAddress)) {
			// TODO: show error
			pairError = 'value must be a valid address'
			return
		}
		pairError = ''
		pairLoading = true

		const uniswapv2PairC: UniswapV2Pair = getContract<UniswapV2Pair>(
			pairAddress,
			(UniswapV2PairABI as any).default,
			$nodeStore.httpProvider as ethers.providers.JsonRpcProvider
		)

		try {
			const [token0Address, token1Address, { _reserve0, _reserve1 }] = await Promise.all([
				uniswapv2PairC.token0(),
				uniswapv2PairC.token1(),
				uniswapv2PairC.getReserves(),
			])

			const token0C = getContract<ERC20>(
				token0Address,
				ERC20Abi,
				$nodeStore.httpProvider as ethers.providers.JsonRpcProvider
			)
			const token1C = getContract<ERC20>(
				token1Address,
				ERC20Abi,
				$nodeStore.httpProvider as ethers.providers.JsonRpcProvider
			)

			const [token0Symbol, token1Symbol, token0Name, token1Name, token0Decimals, token1Decimals] = await Promise.all([
				token0C.symbol(),
				token1C.symbol(),
				token0C.name(),
				token1C.name(),
				token0C.decimals(),
				token1C.decimals(),
			])

			token0 = {
				address: token0Address,
				symbol: token0Symbol,
				name: token0Name,
				decimals: token0Decimals,
				reserve: _reserve0,
			}

			token1 = {
				address: token1Address,
				symbol: token1Symbol,
				name: token1Name,
				decimals: token1Decimals,
				reserve: _reserve1,
			}

			// TODO: use reserves and inputs to calculate amounts
		} catch (err) {
			console.log(err)
			pairError = 'Unable to fetch pair with specified address'
		} finally {
			pairLoading = false
		}

		// const pair = await ammFactory.get
		// const s = await getContract<ERC20>('a', ERC20Abi, $nodeStore.httpProvider as ethers.providers.JsonRpcProvider)
	})

	$: {
		getPairData(pairAddress)
	}

	$: {
		if (token0 && token1 && token0Amount) {
			// TODO: float no. to big number
			// TODO: calculate amountOut and set to token1Amount
			if (token0Amount > 0) {
				token1Amount = floatToBigNumber(token0Amount).mul(token0.reserve.div(token1.reserve)).toNumber()
				console.log(floatToBigNumber(token0Amount).mul(token0.reserve.div(token1.reserve)))
			} else if (token1Amount > 0) {
				token0Amount = floatToBigNumber(token1Amount).mul(token1.reserve.div(token0.reserve)).toNumber()
				console.log(floatToBigNumber(token1Amount).mul(token1.reserve.div(token0.reserve)))
			}
		}
	}
	// onDestroy(() => {
	// 	clearTimeout(timer)
	// })
</script>

<div>WELCOME {$nodeStore.name}</div>

<div class="max-w-2xl px-8 py-4 m-2 bg-white rounded-lg shadow-md dark:bg-gray-800">
	<div class="flex items-center justify-between">
		<div class="text-sm font-light text-gray-600 dark:text-gray-400">
			<select class="select select-bordered w-full max-w-xs" bind:value={selectedIndex}>
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
				<span>
					<input
						class="input w-full {pairError ? 'input-error' : ''}"
						type="text"
						name="pairAddress"
						id="pairAddress"
						placeholder="0x0000000000000000000000000000000000000000"
						disabled={selectedIndex < 0}
						bind:value={pairAddress}
					/>
					{#if pairLoading}
						<div
							style="border-top-color:green-500"
							class="w-8 h-8 border-4 border-green-500 border-dashed rounded-full animate-spin"
						/>
					{/if}
					{#if selectedIndex >= 0 && pairError}
						<label class="label">
							<span class="label-text-alt text-red-500">{pairError}</span>
						</label>
					{/if}
				</span>
			</div>
		{/if}

		{#if viewType === 1}
			<div class="container">
				<label>Token 1</label>
				<input type="text" placeholder="0x0000000000000000000000000000000000000000" bind:value={token0Address} />
			</div>

			<div class="container">
				<label>Token 2</label>
				<input type="text" placeholder="0x0000000000000000000000000000000000000000" bind:value={token1Address} />
			</div>
		{/if}

		{#if token0 && token1}
			<div class="container">
				<label title={token0.name}>{token0.symbol} amount</label>
				<input type="number" placeholder="0.0" bind:value={token0Amount} />
			</div>

			->

			<div class="container">
				<label title={token1.name}>{token1.symbol} amount</label>
				<input type="number" placeholder="0.0" bind:value={token1Amount} />
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
