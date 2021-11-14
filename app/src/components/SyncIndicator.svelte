<script lang="ts">
	type SyncIndicator = 'syncing' | 'synced' | undefined
	export let status: SyncIndicator
</script>

<span class="{$$props.class} dot {status ?? 'offline'}" title={status ?? 'offline'} />

<style>
	.dot {
		display: inline-block;
		width: 0.5rem;
		height: 0.5rem;
		border-radius: 50%;
	}

	.offline {
		background-color: #9ca3af3d;
		box-shadow: 0px 0px 4px 1px #9ca3af;
	}

	.synced {
		background-color: rgb(3 255 46);
		box-shadow: 0px 0px 4px 1px rgb(3 255 46);
	}

	.syncing {
		background-color: #859aff;
		box-shadow: 0px 0px 4px 1px #859aff;
		/* TODO: pulse animation */
		animation: pulse 2s infinite;
	}

	@keyframes pulse {
		0% {
			transform: scale(0.7);
			box-shadow: 0 0 4px 2px #859aff;
		}

		70% {
			transform: scale(1);
			box-shadow: 0 0 4px 5px #859aff00;
		}

		100% {
			transform: scale(0.7);
			box-shadow: 0 0 4px 2px #859affff;
		}
	}
</style>
