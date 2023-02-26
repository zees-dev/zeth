<script lang="ts">
    import { onMount } from "svelte";
    import { formatDate, testRPCConnection } from "../../lib/utils";
  import type { ConnectionStatus, Endpoint } from "../endpoint/types";
    import ConnectionIndicator from "./ConnectionIndicator.svelte";
    import ProtocolBadge from "./ProtocolBadge.svelte";

  export let endpoint: Endpoint;

  let connectionState: ConnectionStatus = 'connecting'
  onMount(async () => {
    try {
      await testRPCConnection(endpoint.rpc_url);
      connectionState = 'connected';
    } catch(err) {
      connectionState = undefined;
    }
  });
</script>

<!-- <div class="card card-side bg-base-100 shadow-xl w-[100%]">
  <figure>
    <div class="avatar online h-[unset]">
      <div class="rounded">
        <img src="https://picsum.photos/100" alt="Endpoint" />
      </div>
    </div>
  </figure>

  <div class="card-body">
    <h2 class="card-title">{endpoint.name}</h2>
    <div class="text-sm text-gray-500 truncate">{endpoint.symbol}</div>
    <div class="text-sm text-gray-500 truncate">{endpoint.rpc_url}</div>
    <div class="text-sm text-gray-500 truncate">{formatDate(endpoint.date_added)}</div>
    {#if endpoint.type === 'http'}
      <div class="text-sm text-gray-500 truncate">HTTP</div>
    {:else if endpoint.type === 'ws'}
      <div class="text-sm text-gray-500 truncate">WS</div>
    {/if}

    <div class="text-sm text-gray-500 truncate">{endpoint.block_explorer_url}</div>
  </div>
</div> -->

<a
  class="card card-side bordered bg-neutral hover:bg-base-300 w-[100%] py-2 {$$props.class}"
  href={`#/endpoints/${endpoint.id}`}
  aria-label={endpoint.name}
>
  <figure class="!w-24 m-auto relative">
    <ConnectionIndicator class="absolute top-0 right-[0.5rem]" status={connectionState} />
    <picture>
      <img src="https://picsum.photos/100" alt="Endpoint" width="72" class="m-auto"/>
    </picture>
  </figure>
  <div class="card-body p-0">
    <h2 class="card-title card-title-grid">
      {endpoint.name}
    </h2>

    <div class="pr-2 pb card-content-grid gap-2">
      <ProtocolBadge class="text-sm truncate" type={endpoint.type} text={endpoint.rpc_url} />
      <div class="text-sm text-gray-500 truncate">{endpoint.symbol}</div>
      <div class="text-sm text-gray-500 truncate">{formatDate(endpoint.date_added)}</div>
      <div class="text-sm text-gray-500 truncate">{endpoint.block_explorer_url}</div>

      <!-- <div class="text-sm flex place-items-center span-column">
        <p>Network ID: {node.network?.chainId ?? '-'}</p>
        {#if node.network}
          <p class="badge badge-sm badge-info mx-2">{getNetworkName(node.network.chainId)}</p>
        {/if}
      </div>
      <div data-tip={node.version} class="tooltip flex span-column">
        <p class="text-sm">
          Version: {(node.version && getVersion(node.version)) || '-'}
        </p>
      </div>
      <div class="text-sm flex span-column">
        <p class="whitespace-nowrap">RPC modules:</p>
        <div class="flex flex-wrap gap-2">
          {#if node.modules}
            {#each getSortedModules(node.modules) as { module, version }}
              <p class="badge first:ml-1 h-6 bg-yellow-400">{module}<sup>{version}</sup></p>
            {/each}
          {:else}
            <p class="badge ml-1 h-6 badge-warning">unknown</p>
          {/if}
        </div>
      </div>
      <p class="text-sm">Peers: {node.peers ?? '-'}</p>
      <p class="text-sm justify-self-end">Date added: {dateWithoutTZ(new Date(node.dateAdded)) ?? '-'}</p> -->
    </div>
  </div>
</a>

<style>
	.card-title-grid {
		display: grid;
		grid-template-columns: auto auto 1fr;
	}
	.card-content-grid {
		display: grid;
		grid-template: 1fr 1fr / auto auto;
    justify-content: space-between;
	}
</style>

