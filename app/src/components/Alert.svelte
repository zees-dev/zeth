<script lang="ts">
	import { onMount } from 'svelte'
	import { fade, fly } from 'svelte/transition'

	type AlertType = 'success' | 'warning' | 'error' | 'info'

	export let type: AlertType = 'success'
	export let message = ''
	export let displayTime = 2000
	let displayed = true

	onMount(() => {
		const timeout = setTimeout(() => (displayed = false), displayTime)
		return () => clearTimeout(timeout)
	})
</script>

<!-- success -->
{#if displayed}
	<div
		class="flex w-full max-w-sm mx-auto overflow-hidden bg-white rounded-lg shadow-md dark:bg-gray-800 right-side"
		in:fly={{ y: 200, duration: 500 }}
		out:fade
	>
		{#if type === 'success'}
			<div class="flex items-center justify-center w-12 bg-emerald-500">
				<svg class="w-6 h-6 text-white fill-current" viewBox="0 0 40 40" xmlns="http://www.w3.org/2000/svg">
					<path
						d="M20 3.33331C10.8 3.33331 3.33337 10.8 3.33337 20C3.33337 29.2 10.8 36.6666 20 36.6666C29.2 36.6666 36.6667 29.2 36.6667 20C36.6667 10.8 29.2 3.33331 20 3.33331ZM16.6667 28.3333L8.33337 20L10.6834 17.65L16.6667 23.6166L29.3167 10.9666L31.6667 13.3333L16.6667 28.3333Z"
					/>
				</svg>
			</div>
			<div class="px-4 py-2 -mx-3">
				<div class="mx-3">
					<span class="font-semibold text-emerald-500 dark:text-emerald-400">Success</span>
					<p class="text-sm text-gray-600 dark:text-gray-200">{message}</p>
				</div>
			</div>
		{:else if type === 'info'}
			<div class="flex items-center justify-center w-12 bg-blue-500">
				<svg class="w-6 h-6 text-white fill-current" viewBox="0 0 40 40" xmlns="http://www.w3.org/2000/svg">
					<path
						d="M20 3.33331C10.8 3.33331 3.33337 10.8 3.33337 20C3.33337 29.2 10.8 36.6666 20 36.6666C29.2 36.6666 36.6667 29.2 36.6667 20C36.6667 10.8 29.2 3.33331 20 3.33331ZM21.6667 28.3333H18.3334V25H21.6667V28.3333ZM21.6667 21.6666H18.3334V11.6666H21.6667V21.6666Z"
					/>
				</svg>
			</div>

			<div class="px-4 py-2 -mx-3">
				<div class="mx-3">
					<span class="font-semibold text-blue-500 dark:text-blue-400">Info</span>
					<p class="text-sm text-gray-600 dark:text-gray-200">{message}</p>
				</div>
			</div>
		{:else if type === 'warning'}
			<div class="flex items-center justify-center w-12 bg-yellow-400">
				<svg class="w-6 h-6 text-white fill-current" viewBox="0 0 40 40" xmlns="http://www.w3.org/2000/svg">
					<path
						d="M20 3.33331C10.8 3.33331 3.33337 10.8 3.33337 20C3.33337 29.2 10.8 36.6666 20 36.6666C29.2 36.6666 36.6667 29.2 36.6667 20C36.6667 10.8 29.2 3.33331 20 3.33331ZM21.6667 28.3333H18.3334V25H21.6667V28.3333ZM21.6667 21.6666H18.3334V11.6666H21.6667V21.6666Z"
					/>
				</svg>
			</div>

			<div class="px-4 py-2 -mx-3">
				<div class="mx-3">
					<span class="font-semibold text-yellow-400 dark:text-yellow-300">Warning</span>
					<p class="text-sm text-gray-600 dark:text-gray-200">{message}</p>
				</div>
			</div>
		{:else if type === 'error'}
			<div class="flex items-center justify-center w-12 bg-red-500">
				<svg class="w-6 h-6 text-white fill-current" viewBox="0 0 40 40" xmlns="http://www.w3.org/2000/svg">
					<path
						d="M20 3.36667C10.8167 3.36667 3.3667 10.8167 3.3667 20C3.3667 29.1833 10.8167 36.6333 20 36.6333C29.1834 36.6333 36.6334 29.1833 36.6334 20C36.6334 10.8167 29.1834 3.36667 20 3.36667ZM19.1334 33.3333V22.9H13.3334L21.6667 6.66667V17.1H27.25L19.1334 33.3333Z"
					/>
				</svg>
			</div>

			<div class="px-4 py-2 -mx-3">
				<div class="mx-3">
					<span class="font-semibold text-red-500 dark:text-red-400">Error</span>
					<p class="text-sm text-gray-600 dark:text-gray-200">{message}</p>
				</div>
			</div>
		{:else}
			<div class="w-2 bg-gray-800 dark:bg-gray-900" />

			<div class="flex items-center px-2 py-3">
				<!-- pass component prop down here -->
				<img
					class="object-cover w-10 h-10 rounded-full"
					alt="User avatar"
					src="https://images.unsplash.com/photo-1477118476589-bff2c5c4cfbb?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=200&q=200"
				/>

				<div class="mx-3">
					<p class="text-gray-600 dark:text-gray-200">
						Sara has replied on the
						<a class="text-blue-500 dark:text-blue-300 hover:text-blue-400 hover:underline">uploaded image</a>.
					</p>
				</div>
			</div>
		{/if}
	</div>
{/if}

<style>
	.right-side {
		position: absolute;
		/* top: 0; */
		right: 0;
		margin-top: 1rem;
		margin-block-end: 1rem;
	}
</style>
