<script lang="ts">
	import { Router, Link, Route } from 'svelte-navigator'
	import { nodeStore } from '../stores/Node'
</script>

<div class="{$$props.class} drawer-side">
	<div class="flex self-center text-center p-2">[LOGO]</div>
	<ul class="menu p-4 overflow-y-auto bg-base-100 text-base-content">
		<li>
			<Link to="/">Dashboard</Link>
		</li>
		{#if $nodeStore !== undefined}
			<li>
				<Link to="node/{$nodeStore.id}">{$nodeStore.name}</Link>
			</li>
			<div class="ml-4">
				<li>
					<Link to="node/{$nodeStore.id}/accounts">Accounts (coming soon)</Link>
				</li>
				<li>
					<Link to="node/{$nodeStore.id}/contracts">Contracts</Link>
				</li>
				<li>
					<Link to="node/{$nodeStore.id}/console">Console (coming soon)</Link>
				</li>
			</div>
		{/if}
		<li>
			<Link to="/settings">Settings</Link>
		</li>
	</ul>
</div>

<style>
	.drawer-side {
		display: grid;
		grid-template-rows: 1fr 11fr;
	}
</style>
