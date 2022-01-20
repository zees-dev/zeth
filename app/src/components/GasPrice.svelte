<script lang="ts">
	import { BigNumber, constants, utils } from 'ethers'
	import { EthNetworks } from '../lib/const'
	import { Node } from '../lib/Models/Node'

	export let node: Node
	export let realtime: boolean = false

	let useEIP1559 = false
	let blockNumber = 0
	let networkId = 0
	let gasFees = {
		gasPrice: '0',
		maxFeePerGas: '0',
		maxPriorityFeePerGas: '0',
		baseFeePerGas: '0',
	}

	async function getFees() {
		const feeData = await node.httpProvider?.getFeeData()

		const gasPrice = utils.formatUnits(feeData?.gasPrice ?? constants.Zero, 'gwei')
		const maxFeePerGas = utils.formatUnits(feeData?.maxFeePerGas ?? constants.Zero, 'gwei')
		const maxPriorityFeePerGas = utils.formatUnits(feeData?.maxPriorityFeePerGas ?? constants.Zero, 'gwei')

		// basefee = (maxFeePerGas - maxPriorityFeePerGas) / 2
		const baseFeeWei = feeData?.maxFeePerGas?.sub(feeData?.maxPriorityFeePerGas ?? constants.Zero).div(2)
		const baseFeePerGas = utils.formatUnits(baseFeeWei ?? constants.Zero, 'gwei')

		return {
			gasPrice,
			maxFeePerGas,
			maxPriorityFeePerGas,
			baseFeePerGas,
		}
	}
</script>

<h2>Gas Prices</h2>
{#await getFees()}
	<p>...retrieving gas fee data</p>
{:then fee}
	<p>Gas Price: {parseFloat(fee.gasPrice).toFixed(2)} Gwei</p>
	<p>MaxFeePerGas: {fee.maxFeePerGas}</p>
	<p>MaxPriorityFeePerGas: {fee.maxPriorityFeePerGas}</p>
	<p>Base Fee: {fee.baseFeePerGas}</p>
{:catch error}
	<p style="color: red">{error.message}</p>
{/await}
