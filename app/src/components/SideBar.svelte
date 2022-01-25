<script lang="ts">
	import { Router, Link, Route } from 'svelte-navigator'
	import { nodeStore } from '../stores/Node'
</script>

<div class="{$$props.class} drawer-side">
	<div class="flex self-center text-center p-2">
		<div>[LOGO]</div>
		<button class="p-2">
			<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512" class="w-5 h-5 fill-current text-coolGray-800">
				<rect width="352" height="32" x="80" y="96" />
				<rect width="352" height="32" x="80" y="240" />
				<rect width="352" height="32" x="80" y="384" />
			</svg>
		</button>
	</div>
	<ul class="menu p-4 overflow-y-auto bg-base-100 text-base-content">
		<li>
			<Link to="/" class="flex gap-2">
				<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512" class="w-5 h-5 fill-current text-coolGray-600">
					<path
						d="M469.666,216.45,271.078,33.749a34,34,0,0,0-47.062.98L41.373,217.373,32,226.745V496H208V328h96V496H480V225.958ZM248.038,56.771c.282,0,.108.061-.013.18C247.9,56.832,247.756,56.771,248.038,56.771ZM448,464H336V328a32,32,0,0,0-32-32H208a32,32,0,0,0-32,32V464H64V240L248.038,57.356c.013-.012.014-.023.024-.035L448,240Z"
					/>
				</svg>
				<span>Dashboard</span>
			</Link>
		</li>
		{#if $nodeStore !== undefined}
			<li>
				<Link to="node/{$nodeStore.id}" class="flex gap-2">{$nodeStore.name}</Link>
			</li>
			<div class="ml-4">
				<li>
					<Link to="node/{$nodeStore.id}/accounts" class="flex gap-2">Accounts (coming soon)</Link>
				</li>
				<li>
					<Link to="node/{$nodeStore.id}/contracts" class="flex gap-2">Contracts</Link>
				</li>
				<li>
					<Link to="node/{$nodeStore.id}/console" class="flex gap-2">Console (coming soon)</Link>
				</li>
			</div>
		{/if}
		<li>
			<Link to="/settings" class="flex gap-2">
				<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512" class="w-5 h-5 fill-current text-coolGray-600">
					<path
						d="M245.151,168a88,88,0,1,0,88,88A88.1,88.1,0,0,0,245.151,168Zm0,144a56,56,0,1,1,56-56A56.063,56.063,0,0,1,245.151,312Z"
					/>
					<path
						d="M464.7,322.319l-31.77-26.153a193.081,193.081,0,0,0,0-80.332l31.77-26.153a19.941,19.941,0,0,0,4.606-25.439l-32.612-56.483a19.936,19.936,0,0,0-24.337-8.73l-38.561,14.447a192.038,192.038,0,0,0-69.54-40.192L297.49,32.713A19.936,19.936,0,0,0,277.762,16H212.54a19.937,19.937,0,0,0-19.728,16.712L186.05,73.284a192.03,192.03,0,0,0-69.54,40.192L77.945,99.027a19.937,19.937,0,0,0-24.334,8.731L21,164.245a19.94,19.94,0,0,0,4.61,25.438l31.767,26.151a193.081,193.081,0,0,0,0,80.332l-31.77,26.153A19.942,19.942,0,0,0,21,347.758l32.612,56.483a19.937,19.937,0,0,0,24.337,8.73l38.562-14.447a192.03,192.03,0,0,0,69.54,40.192l6.762,40.571A19.937,19.937,0,0,0,212.54,496h65.222a19.936,19.936,0,0,0,19.728-16.712l6.763-40.572a192.038,192.038,0,0,0,69.54-40.192l38.564,14.449a19.938,19.938,0,0,0,24.334-8.731L469.3,347.755A19.939,19.939,0,0,0,464.7,322.319Zm-50.636,57.12-48.109-18.024-7.285,7.334a159.955,159.955,0,0,1-72.625,41.973l-10,2.636L267.6,464h-44.89l-8.442-50.642-10-2.636a159.955,159.955,0,0,1-72.625-41.973l-7.285-7.334L76.241,379.439,53.8,340.562l39.629-32.624-2.7-9.973a160.9,160.9,0,0,1,0-83.93l2.7-9.972L53.8,171.439l22.446-38.878,48.109,18.024,7.285-7.334a159.955,159.955,0,0,1,72.625-41.973l10-2.636L222.706,48H267.6l8.442,50.642,10,2.636a159.955,159.955,0,0,1,72.625,41.973l7.285,7.334,48.109-18.024,22.447,38.877-39.629,32.625,2.7,9.972a160.9,160.9,0,0,1,0,83.93l-2.7,9.973,39.629,32.623Z"
					/>
				</svg>
				<span>Settings</span>
			</Link>
		</li>
		<li>
			<Link to="/settings" class="flex gap-2">
				<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512" class="w-5 h-5 fill-current text-coolGray-600">
					<path
						d="M440,424V88H352V13.005L88,58.522V424H16v32h86.9L352,490.358V120h56V456h88V424ZM320,453.642,120,426.056V85.478L320,51Z"
					/>
					<rect width="32" height="64" x="256" y="232" />
				</svg>
				<span>Logout</span>
			</Link>
		</li>
	</ul>

	<div class="text-center my-2">Zeth <sup>v1.0.0</sup></div>
	<!-- TODO: github icon -->
	<!-- TODO: gitbook icon -->
</div>

<style>
	.drawer-side {
		display: grid;
		grid-template-rows: 1fr 11fr;
	}
</style>
