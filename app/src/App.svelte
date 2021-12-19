<script lang="ts">
	import { Router, Route } from 'svelte-routing'
	import NavBar from './components/NavBar.svelte'
	import Menu from './components/Menu.svelte'

	// content views
	import Dashboard from './routes/Dashboard.svelte'
	import Node from './routes/Node.svelte'

	export let url = ''
</script>

<Router {url}>
	<main class="main">
		<NavBar />
		<Menu />
		<section class="overflow-x-auto">
			<Route path="node/:id" let:params>
				<Node id={params.id} />
			</Route>
			<Route path="/"><Dashboard /></Route>
		</section>
	</main>
</Router>

<style global lang="postcss">
	@tailwind base;
	@tailwind components;
	@tailwind utilities;

	.main {
		color: var(--text);
		background: var(--bg);
		transition: background 500ms ease-in-out, color 1000ms ease-in-out;

		margin: 0 auto;
		padding: 0;
		width: 100%;
		height: 100%;
		display: grid;
		grid-template: 1fr 11fr/ 15rem auto;
	}
</style>
