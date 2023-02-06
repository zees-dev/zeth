<script lang="ts">
  import Router, { link } from "svelte-spa-router";
  import wrap from "svelte-spa-router/wrap";
  import active from "svelte-spa-router/active";
  import { dbStore } from "../../stores/db";
  import type { Endpoint } from "./types";
  import Spinner from "../Spinner.svelte";

  const routes = {
    "/endpoints/:endpointId": wrap({ asyncComponent: () => import('./Status.svelte')}),
    "/endpoints/:endpointId/rpc-log": wrap({ asyncComponent: () => import('./RPCLog.svelte')}),
    "*": wrap({ asyncComponent: () => import('../NotFound.svelte')}),
  }

  // route params
  export let params: { endpointId: string };

  async function getEndpoint(endpointId: string): Promise<Endpoint[]> {
    const response = await $dbStore.db.query(`SELECT * FROM endpoint WHERE id=${endpointId}`);
    return (response as any[])[0].result as Endpoint[];
  }
</script>

<div class="content">
  <nav class="sub-nav overflow-y-auto">
    {#await getEndpoint(params.endpointId)}
      <Spinner size="lg" />
    {:then endpoints}
      {#each endpoints as endpoint}
        <h1 class="search-bar">&gt {endpoint.name} &lt</h1>
        <a
          href={`/endpoints/${params.endpointId}`}
          class="menu-item"
          aria-label="Status"
          use:link
          use:active={{ path: `/endpoints/${params.endpointId}`, className: "current-route" }}
        >
          Status
          <i class="ri-database-2-line" />
        </a>
        <a
          href={`/endpoints/${params.endpointId}/rpc-log`}
          class="menu-item"
          aria-label="RPC Log"
          use:link
          use:active={{ path: `/endpoints/${params.endpointId}/rpc-log`, className: "current-route" }}
        >
          RPC Log
          <i class="ri-database-2-line" />
        </a>
        <!-- TODO -->
        <!--
        <ul>
          <li>Accounts</li>
          <li>Contracts</li>
          <li>Metrics</li>
          <li>Console</li>
        </ul> -->
      {/each}
    {:catch}
      <p class="bg-error">Failed to fetch endpoint {params.endpointId}</p>
    {/await}
  </nav>

  <section>
    <Router {routes} />
  </section>
</div>

<style>
  .sub-nav {
    width: 100%;
    height: 100%;
    
    display: flex;
    flex-flow: column;
  }

  .content {
    width: 100%;
    height: 100%;

    display: grid;
    grid-template-columns: 12rem auto;
  }
</style>