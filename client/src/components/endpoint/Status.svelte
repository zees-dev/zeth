<script lang="ts">
  import { onMount } from "svelte";
    import { dbStore } from "../../stores/db";
    import Spinner from "../Spinner.svelte";
  import { getEndpoint } from "./queries";
    import type { Endpoint } from "./types";

  export let params: { endpointId: string };

  let loading = true;
  let endpoint: Endpoint | undefined = undefined;
  let err: string | undefined | unknown = undefined;

  onMount(() => getEndpointInComponent());

  async function getEndpointInComponent() {
    const { loading: l, result, error } = await getEndpoint($dbStore.db, params.endpointId);
    loading = l;
    endpoint = result;
    err = error;
  }
</script>

<div>
  Status {params.endpointId}!
  {#if loading}
    <Spinner size="lg" />
  {:else if endpoint}
    <div>
      {endpoint.name}
    </div>
  {:else if err}
    <div>
      <p class="bg-error">Failed to fetch endpoint {params.endpointId}</p>
    </div>
  {/if}
</div>