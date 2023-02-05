<script lang="ts">
  import { web3ProviderStore } from "../../stores/web3provider";
  import AddEndpoint from "../endpoint/AddEndpoint.svelte";
  import SidePanel from "../SidePanel.svelte";

  let showAddEndpointPanel = false;
  // TODO: metamask connection button UI
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
    <SidePanel bind:show={showAddEndpointPanel} component={AddEndpoint} />
    <!-- CARD (node stats) - left -->
    <!-- Button add node (right) -->

    <!-- Connect to metamask -->
    <!-- node list to metamask -->
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