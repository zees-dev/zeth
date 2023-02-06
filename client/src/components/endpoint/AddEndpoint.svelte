<script lang="ts">
  import { z } from "zod";
  import { dbStore } from "../../stores/db";
  import { loginStore } from "../../stores/login";

  export let onSuccessfulSubmission: () => void = () => {};

  const formSchema = z.object({
    endpointName: z.string()
      .min(3)
      .max(64)
      .trim(),
    rpcUrl: z.string().regex(/^(http|https|ws|wss):\/\/.+/, "Invalid RPC URL, must be of scheme http, https, ws or wss"),
    testConnection: z.boolean(),
    symbol: z.string().min(3).max(10).trim(),
  });

  let submissionError: string | undefined;
  let loading = false;
  let endpointName = '';
  let rpcUrl = '';
  let testConnection = true;
  let symbol = 'ETH';
  let blockExplorerUrl = '';

  async function testRPCConnection(type: "http" | "ws") {
    if (type === 'http') {
      const res = await fetch(rpcUrl, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          "jsonrpc": "2.0",
          "method": "eth_blockNumber",
          "params": [],
          "id": 1
        }),
      });
      if (!res.ok) {
        throw new Error(`Failed to connect to ${rpcUrl}`);
      }
    } else {
      const socket = new WebSocket(rpcUrl);
      socket.onopen = () => {
        socket.send(JSON.stringify({
          "jsonrpc": "2.0",
          "method": "eth_blockNumber",
          "params": [],
          "id": 1
        }));
      }
      try {
        await new Promise((resolve, reject) => {
          socket.onmessage = (msg) => {
            const res = JSON.parse(msg.data);
            if (res.error) {
              reject(res.error);
              return;
            }
            resolve(res);
          }
          socket.onerror = (err) => reject(err);
        })
      } catch (err) {
        throw new Error(`Failed to connect to ${rpcUrl}`);
      } finally {
        socket!.close();
      }
    }
  }

  async function handleSubmit() {
    if (validationErrors.length) return;
    loading = true;
    try {
      const type = rpcUrl.startsWith('http') ? 'http' : 'ws'
      // perform external API request if test connection is checked
      if (testConnection) {
        await testRPCConnection(type);
      }

      const created = await $dbStore.db.create("endpoint", {
        user: $loginStore.userId, // TODO: this should be internally set
        name: endpointName,
        enabled: true,
        date_added: (new Date()).toISOString(),
        rpc_url: rpcUrl,
        type,
      });
      console.info("created", created);
      onSuccessfulSubmission();
    } catch(err) {
      submissionError = err as string;
      setTimeout(() => submissionError = undefined, 5000);
    } finally {
      loading = false;
    }
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

<form style="width: 100%;" on:submit|preventDefault={handleSubmit}>
  <div class="flex flex-col mb-4">
    <label class="text-sm font-medium mb-2">Endpoint name</label>
    <input type="text" class="border p-2 rounded" bind:value={endpointName} required />
    {#if validationErrors.find(i => i.path.includes('endpointName'))}
      <div class="text-sm text-red-500 self-end mr-2">
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
      pattern="^(http|https|ws|wss):\/\/.+"
      placeholder="wss://mainnet.infura.io/ws/v3/abcdef, https://mainnet.infura.io/v3/abcdef" />
    <div class="mt-2 self-end mr-2" on:click={() => testConnection = !testConnection}>
      <input type="checkbox" checked={testConnection} />
      <label class="text-sm font-medium ml-2">Test connection</label>
    </div>
    {#if validationErrors.find(i => i.path.includes('rpcUrl'))}
      <div class="text-sm text-red-500 self-end mr-2">
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
      <div class="text-sm text-red-500 self-end mr-2">
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
  
  {#if submissionError}
    <div class="bg-red-100 border-l-4 border-red-500 text-red-700 p-4 mb-4" role="alert">
      <p>{submissionError}</p>
    </div>
  {/if}

  <button
    type="submit"
    class={(validationErrors.length > 0 ? "disabled bg-gray-300" : "bg-blue-500 hover:bg-blue-600") + " text-white p-2 rounded flex m-auto"}
    disabled={validationErrors.length > 0}
  >
    Create
  </button>
</form>
