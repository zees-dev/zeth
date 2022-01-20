<script lang="ts">
	import { Router, Link, Route } from 'svelte-navigator'
	import NavBar from './components/NavBar.svelte'

	// content views
	import Dashboard from './routes/Dashboard.svelte'
	import SideBar from './components/SideBar.svelte'
	import Contracts from './features/Node/Contracts.svelte'
	import NodeRouter from './routes/NodeRouter.svelte'
	import Settings from './routes/Settings.svelte'

	export let url = ''
</script>

<Router {url}>
	<main class="main">
		<SideBar class="sidebar" />
		<NavBar />
		<section class="content overflow-x-auto">
			<Route path="/"><Dashboard /></Route>
			<Route path="node/:id/*" let:params>
				<NodeRouter id={params.id} />
			</Route>
			<Route path="settings"><Settings /></Route>
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

	.content {
		grid-column-start: 2;
	}

	.sidebar {
		grid-row: 1/-1;
	}
</style>
