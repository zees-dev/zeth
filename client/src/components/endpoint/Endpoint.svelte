<script lang="ts">
  // route params
  export let params: { endpointId: string };

  // TODO query endpoint by id

  console.log(params.endpointId);
</script>

<h1>
  Selected Endpoint {params.endpointId}!
</h1>