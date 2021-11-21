<script lang="ts">
	import { utils } from 'ethers'
	import { ethersProvider } from '../stores/Settings'
	import { EthNetworks } from '../lib/const'

	const username = 'john'
	const theme = 'dark'

	let useEIP1559 = false
	let blockNumber = 0
	let networkId = 0
	let gasFees = {
		gasPrice: '0',
		maxFeePerGas: '0',
		maxPriorityFeePerGas: '0',
		baseFeePerGas: '0'
	}

	async function getFees() {
		const [feeData] = await Promise.allSettled([$ethersProvider.getFeeData()])

		if (feeData.status === 'fulfilled') {
			gasFees.gasPrice = utils.formatUnits(feeData.value.gasPrice ?? 0, 'gwei')
			gasFees.maxFeePerGas = utils.formatUnits(feeData.value.maxFeePerGas ?? 0, 'gwei')
			gasFees.maxPriorityFeePerGas = utils.formatUnits(feeData.value.maxPriorityFeePerGas ?? 0, 'gwei')

			// basefee = (maxFeePerGas - maxPriorityFeePerGas) / 2
			const baseFeeWei = feeData.value.maxFeePerGas?.sub(feeData.value.maxPriorityFeePerGas ?? 0).div(2)
			gasFees.baseFeePerGas = utils.formatUnits(baseFeeWei ?? 0, 'gwei')
		}
	}

	$: {
		if ($ethersProvider) {
			// initial call on-page-load
			$ethersProvider.getBlockNumber().then((bn) => {
				blockNumber = bn
				networkId = $ethersProvider.network?.chainId
				getFees()
			})

			// update on block-number change
			$ethersProvider.on('block', (bn) => {
				blockNumber = bn
				getFees()
			})
		}
	}
</script>

<header>
	<div class="flex place-items-center py-1 px-2">
		<svg width="32" height="32" viewBox="0 0 32 32" xmlns="http://www.w3.org/2000/svg">
			<g fill="none" fillRule="evenodd">
				<path d="M10 0h12a10 10 0 0110 10v12a10 10 0 01-10 10H10A10 10 0 010 22V10A10 10 0 0110 0z" fill="#FFF" />
				<path d="M5.3 10.6l10.4 6v11.1l-10.4-6v-11zm11.4-6.2l9.7 5.5-9.7 5.6V4.4z" fill="#555AB9" />
				<path d="M27.2 10.6v11.2l-10.5 6V16.5l10.5-6zM15.7 4.4v11L6 10l9.7-5.5z" fill="#91BAF8" />
			</g>
		</svg>
		<div class="ml-2">
			<div class="mt-6">
				{EthNetworks[networkId] ?? 'Unknown'}<sup title="Network ID: {networkId}">{networkId}</sup>
			</div>
			<div title="block no. {blockNumber}" class="ml-6 w-16">{blockNumber}</div>
		</div>
	</div>

	<div class="flex place-items-center gap-2">
		<!-- TODO: USE SWITCH HERE -->
		<input type="checkbox" bind:checked={useEIP1559} class="toggle toggle-sm" />
		<div>EIP1559</div>
		{#if useEIP1559}
			<div>
				<div>base fee: {Math.round(+gasFees.baseFeePerGas) || '?'}</div>
				<div>priority fee: {Math.trunc(+gasFees.maxPriorityFeePerGas) || '?'}</div>
			</div>
		{:else}
			<div>gas price (legacy): {Math.round(+gasFees.gasPrice) || '?'}</div>
		{/if}
	</div>

	<button class="btn btn-sm btn-outline">
		<svg
			xmlns="http://www.w3.org/2000/svg"
			fill="none"
			viewBox="0 0 24 24"
			class="inline-block w-6 h-6 mr-2 stroke-current"
		>
			<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
		</svg>
		Connect Wallet
	</button>

	<div class="user">
		<div>[N]</div>
		<!-- TODO: USE SWITCH HERE -->
		<div>[Theme]</div>
		<div class="px-3">
			<p>{username}</p>
			<p class="flex"><span class="mr-2">[X]</span> logout</p>
		</div>
	</div>
</header>

<style>
	header {
		display: grid;
		align-items: center;
		grid-template-columns: min-content 1fr auto min-content;
		grid-gap: 0.5rem;
		border: 1px solid #555ab9;
	}

	.user {
		justify-self: end;
		display: grid;
		column-gap: 1rem;
		grid-template-columns: min-content min-content 1fr;
		place-items: center;
	}
</style>
