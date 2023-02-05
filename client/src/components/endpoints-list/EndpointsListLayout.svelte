<script lang="ts">
  import { onMount } from "svelte";
  import { dbStore } from "../../stores/db";
  import { web3ProviderStore } from "../../stores/web3provider";
  import AddEndpoint from "../endpoint/AddEndpoint.svelte";
  import SidePanel from "../SidePanel.svelte";

  let showAddEndpointPanel = false;
  // TODO: metamask connection button UI

  let loading = true;
  let endpoints: any[] | undefined = undefined;
  let err: string | undefined = undefined;

  onMount(async () => getEndpoints());

  async function getEndpoints() {
    loading = true;
    try {
      const response = await $dbStore.db.query("SELECT * FROM endpoint;");
      endpoints = (response as any[])[0].result;
    } catch(err) {
      endpoints = undefined;
      err = err;
    } finally {
      loading = false;
    }
  }
</script>

<div class="content">
  <nav class="sub-nav">
    <div class="search-bar">search....</div>
    {#if $web3ProviderStore.provider}
      <a
        href="#/"
        class="menu-item"
        aria-label="Metamask"
        on:click={web3ProviderStore.disconnect}
      >
        Disconnect
        <i class="ri-database-2-line" />
      </a>
    {:else}
      <a
        href="#/"
        class="menu-item"
        aria-label="Metamask"
        on:click={web3ProviderStore.connect}
      >
        Metamask
        <i class="ri-database-2-line" />
      </a>
    {/if}
    <a
      href={`#/endpoints/1`}
      class="menu-item"
      aria-label="Node1"
    >
      Node1
      <i class="ri-database-2-line" />
    </a>
  </nav>

  <section>
    <h1>Endpoints</h1>
    <button on:click={() => showAddEndpointPanel = true}>Add endpoint</button>
    <!-- CARD (node stats) - left -->
    <!-- Button add node (right) -->
    <SidePanel bind:show={showAddEndpointPanel} component={AddEndpoint} onSuccessfulSubmission={getEndpoints} />

    <!-- Connect to metamask -->
    <div class="flex flex-row">
      {#if $web3ProviderStore.provider}
        <p>Metamask node details</p>
      {:else}
        <button on:click={web3ProviderStore.connect}>Connect metamask</button>
      {/if}
    </div>

    <!-- node list to metamask -->
    {#if loading}
      <p>fetching endpoints...</p>
    {:else if endpoints}
      {#each endpoints as ep (ep.id)}
        <div class="flex flex-row">
          <h2>{ep.name}</h2>
          <p>-{ep.rpc_url}</p>
        </div>
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
  }
</style>