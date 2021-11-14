<script lang="ts">
	import type { SyncStatus } from '../lib/types'

	export let syncing: SyncStatus
	export let blockNumber: number

	let syncPercentage = '0'
	$: {
		if (syncing) {
			syncPercentage = ((parseInt(syncing.currentBlock, 16) / parseInt(syncing.highestBlock, 16)) * 100).toFixed(2)
		}
	}
</script>

<div class="mb-1 px-2 flex flex-col">
	{#if syncing}
		<div class="mb-1 flex self-end">
			<span class="text-sm font-medium text-blue-700">
				{syncPercentage}%
			</span>
		</div>
		<div class="w-full bg-gray-200 rounded-full">
			<div class="absolute w-full p-0.5 text-center text-blue-500 text-xs font-medium">
				{parseInt(syncing.currentBlock, 16)} / {parseInt(syncing.highestBlock, 16)}
			</div>
			<div
				class="bg-blue-600 text-xs font-medium text-blue-500 text-center p-0.5 leading-none rounded-full"
				style="width: {syncPercentage}%"
			>
				{'syncing...'}
			</div>
		</div>
	{:else if blockNumber}
		<div class="flex self-end">
			<span class="text-sm font-medium text-green-700">100%</span>
		</div>
		<div class="w-full bg-gray-200 rounded-full">
			<div class="bg-green-500 w-full text-xs font-medium text-blue-100 text-center p-0.5 leading-none rounded-full">
				{blockNumber}
			</div>
		</div>
	{:else}
		<div class="flex self-end">
			<span class="text-sm font-medium text-gray-500">-%</span>
		</div>
		<div class="bg-gray-200 rounded-full">
			<div class="bg-gray-400 w-full text-xs font-medium text-blue-100 text-center p-0.5 leading-none rounded-full">
				-
			</div>
		</div>
	{/if}
</div>
