<script lang="ts">
	import { ethersProvider } from '../stores'

	let blockNumber = -1
	$: {
		if ($ethersProvider) {
			$ethersProvider.getBlockNumber().then((bn) => {
				blockNumber = bn
			})
			$ethersProvider.on('block', (bn) => {
				blockNumber = bn
			})
		}
	}
</script>

<nav>
	<section>menu</section>
	<section>{blockNumber}</section>
</nav>

<style>
	nav {
		display: grid;
		grid-row: 1 / -1;
		background: var(--bg-nav);
		color: var(--text);
	}
</style>
