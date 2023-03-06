<script lang="ts">
  import { onMount } from "svelte";
  import { dbStore } from "../../stores/db";
  import { web3ProviderStore } from "../../stores/web3provider";
  import AddEndpoint from "../endpoint/AddEndpoint.svelte";
  import SidePanel from "../SidePanel.svelte";
  import Spinner from "../Spinner.svelte";
  import type { Endpoint } from "../endpoint/types";
  import EndpointPanel from "./EndpointPanel.svelte";
  import ProtocolBadge from "./ProtocolBadge.svelte";
    import { endpointType } from "../../lib/utils";

  let showAddEndpointPanel = false;
  // TODO: metamask connection button UI

  let loading = true;
  let endpoints: Endpoint[] | undefined = undefined;
  let err: string | undefined = undefined;

  onMount(async () => getEndpoints());

  async function getEndpoints() {
    loading = true;
    try {
      const response = await $dbStore.db.query("SELECT * FROM endpoint;");
      endpoints = (response as any[])[0].result as Endpoint[];
      console.log(response);
    } catch(err) {
      endpoints = undefined;
      err = err;
    } finally {
      loading = false;
    }
  }
</script>

<div class="content">
  <nav class="sub-nav overflow-y-auto overflow-x-hidden">
    <div class="search-bar">search....</div>
    {#if $web3ProviderStore.provider}
      <button
        class="btn btn-sm gap-2 text-sm menu-item"
        on:click={web3ProviderStore.disconnect}
      >
        Disconnect
      </button>
    {:else}
      <button
        class="btn btn-sm gap-2 text-sm menu-item"
        on:click={web3ProviderStore.connect}
      >
        Metamask
      </button>
    {/if}
    {#if loading}
      <Spinner size="lg" />
    {:else if endpoints}
      {#each endpoints as ep (ep.id)}
        <div class="indicator w-[90%]">
          <div class="indicator-item">
            <div class="tooltip" data-tip={ep.rpc_url}>
              <ProtocolBadge type={endpointType(ep.rpc_url)} />
            </div>
          </div>
          <a
            href={`#/endpoints/${ep.id}`}
            class="btn btn-sm gap-2 text-sm w-full menu-item"
            aria-label={ep.name}
          >
            {ep.name}
          </a>
        </div>
      {/each}
    {:else if err}
      <p class="bg-error">failed to fetch endpoints: {err}</p>
    {/if}
  </nav>

  <section class="mr-2">
    <h1 class="text-2xl font-bold">Endpoints</h1>

    {#if loading}
      <p>fetching endpoints...</p>
    {:else if endpoints}

      <!-- TODO: migrate this to a standalone component which performs requests against DB to get stats -->
      <div class="stats flex">
        <!-- TODO CARD (node stats) - left -->

        <div class="stat place-items-center">
          <div class="stat-title">Total endpoints</div>
          <div class="stat-value">{endpoints.length}</div>
          <div class="stat-desc gap-2">
            <div class="badge badge-secondary">{endpoints.filter(e => endpointType(e.rpc_url) === 'http').length} HTTP</div>
            <div class="badge badge-accent">{endpoints.filter(e => endpointType(e.rpc_url) === 'ws').length} WS</div>
          </div>
        </div>
        
        <div class="stat place-items-center">
          <div class="stat-title">Total RPC requests</div>
          <div class="stat-value text-secondary">TODO</div>
          <!-- <div class="stat-desc text-secondary">↗︎ 40 (2%)</div> -->
        </div>
        
        <div class="stat place-items-center">
          <div class="stat-title">Total users</div>
          <div class="stat-value">TODO</div>
          <!-- <div class="stat-desc">↘︎ 90 (14%)</div> -->
        </div>

        <div class="stat place-items-center">
          <button class="stat-title btn btn-primary btn-sm mt-3" on:click={() => showAddEndpointPanel = true}>Add endpoint</button>
          <SidePanel bind:show={showAddEndpointPanel} component={AddEndpoint} onSuccessfulSubmission={getEndpoints} />
          <div class="stat-value">TODO</div>
        </div>
      </div>

      <!-- Connect to metamask -->
      <div class="flex flex-row">
        {#if $web3ProviderStore.provider}
          <button class="btn btn-secondary w-full" on:click={web3ProviderStore.disconnect}>Metamask node details</button>
        {:else}
          <button class="btn btn-primary w-full" on:click={web3ProviderStore.connect}>Connect metamask</button>
        {/if}
      </div>

      {#each endpoints as ep (ep.id)}
        <EndpointPanel endpoint={ep} class="my-2"/>
      {/each}
    {:else if err}
      <p class="bg-error">failed to fetch endpoints: {err}</p>
    {/if}

    <!-- {#await getEndpoints()}
      <p>fetching endpoints...</p>
    {:then endpoints }
      {#each endpoints as ep (ep.id)}
        <div class="flex flex-row">
          <h2>{ep.name}</h2>
          <p>-{ep.rpc_url}</p>
        </div>
      {/each}
    {:catch error}
      <p class="bg-error">failed to fetch endpoints: {error}</p>
    {/await} -->
  </section>
</div>

<style>
  .content {
    width: 100%;
    height: 100%;

    display: grid;
    grid-template-columns: 12rem auto;
  }

  .sub-nav {
    width: 100%;
    height: 100%;

    display: flex;
    flex-flow: column;
    gap: 0.5rem;
    padding: 0rem 1rem 0rem 1rem;
  }
</style>