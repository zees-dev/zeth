<script lang="ts">
	import Header from './components/Header.svelte'
	import Nav from './components/Nav.svelte'
	import Content from './components/Content.svelte'
</script>

<main class="dark main">
	<Header />
	<Nav />
	<Content />
</main>

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
