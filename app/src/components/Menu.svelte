<script lang="ts">
	import { Router, Link, Route } from 'svelte-routing'

	// import { ethersProvider } from '../stores'

	// let blockNumber = -1
	// $: {
	// 	if ($ethersProvider) {
	// 		$ethersProvider.getBlockNumber().then((bn) => {
	// 			blockNumber = bn
	// 		})
	// 		$ethersProvider.on('block', (bn) => {
	// 			blockNumber = bn
	// 		})
	// 	}
	// }
</script>

<nav>
	<section>menu</section>
	<Link to="/">Dashboard</Link>
	<Link to="node/1">Node</Link>
	<!-- <section>{blockNumber}</section> -->
</nav>

<style>
	nav {
		display: grid;
		grid-row: 1 / -1;
		background: var(--bg-nav);
		color: var(--text);
	}
</style>
