<script lang="ts">
  import { z } from "zod";

  let submissionDisabled = true;

  const formSchema = z.object({
    endpointName: z.string()
      .min(3)
      .max(64)
      .trim(),
    rpcUrl: z.string().regex(/^(http|https|ws|wss):\/\/.+/, "Invalid RPC URL, must be of scheme http, https, ws or wss"),
    testConnection: z.boolean(),
    symbol: z.string().min(3).max(10).trim(),
  });

  let endpointName = '';
  let rpcUrl = '';
  let testConnection = true;
  let symbol = 'ETH';
  let blockExplorerUrl = '';

  function handleSubmit() {
    if (validationErrors.length) return;

    // Perform external API request if test connection is checked
    if (testConnection) {
      // Perform API request here
    }

    // TODO: ensure endpoint name is unique
    // TODO: ensure endpoint rpc url is unique

    // Submit form data
    console.log({ endpointName, rpcUrl, symbol, blockExplorerUrl });
  }

  interface ValidationError {
    code: string;
    message: string;
    path: string[];
  }
  let validationErrors: ValidationError[] = [];
  $: {
    try {
      formSchema.parse({ endpointName, rpcUrl, testConnection, symbol, blockExplorerUrl})
      validationErrors = []
    } catch (err) {
      validationErrors = (err as any).issues;
    }
  }
</script>

<form on:submit|preventDefault={handleSubmit}>
  <div class="flex flex-col mb-4">
    <label class="text-sm font-medium mb-2">Endpoint name</label>
    <input type="text" class="border p-2 rounded" bind:value={endpointName} required />
    {#if validationErrors.find(i => i.path.includes('endpointName'))}
      <div class="text-red-500 self-end mr-2">
        ❌ {validationErrors.find(i => i.path.includes('endpointName'))?.message}
      </div>
    {/if}
  </div>
  <div class="flex flex-col mb-4">
    <label class="text-sm font-medium mb-2">RPC URL</label>
    <input type="text" 
      class="border p-2 rounded"
      bind:value={rpcUrl} 
      required
      pattern="^(wss|ws|https|http)://[a-zA-Z0-9._-]+"
      placeholder="wss://mainnet.infura.io/ws/v3/abcdef, https://mainnet.infura.io/v3/abcdef" />
    <div class="mt-2 self-end mr-2" on:click={() => testConnection = !testConnection}>
      <input type="checkbox" checked={testConnection} />
      <label class="text-sm font-medium ml-2">Test connection</label>
    </div>
    {#if validationErrors.find(i => i.path.includes('rpcUrl'))}
      <div class="text-red-500 self-end mr-2">
        ❌ {validationErrors.find(i => i.path.includes('rpcUrl'))?.message}
      </div>
    {/if}
  </div>
  <div class="flex flex-col mb-4">
    <label class="text-sm font-medium mb-2">Symbol</label>
    <input type="text"
      class="border p-2 rounded"
      bind:value={symbol}
      required
      placeholder="ETH" />
    {#if validationErrors.find(i => i.path.includes('symbol'))}
      <div class="text-red-500 self-end mr-2">
        ❌ {validationErrors.find(i => i.path.includes('symbol'))?.message}
      </div>
    {/if}
  </div>
  <div class="flex flex-col mb-4">
    <label class="text-sm font-medium mb-2">Block explorer URL</label>
    <input type="text"
      class="border p-2 rounded"
      bind:value={blockExplorerUrl}
      placeholder="https://etherscan.io" />
  </div>

  <button
    type="submit"
    class={(validationErrors.length > 0 ? "disabled" : "bg-blue-500 hover:bg-blue-600") + "text-white p-2 rounded flex m-auto"}
    disabled={validationErrors.length > 0}
  >
    Create
  </button>
</form>
