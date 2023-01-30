<!-- Component imports and setup -->
<script lang="ts">
  import Surreal from 'surrealdb.js';
  import { loginStore } from '../../stores/login';
  import { z } from "zod";

  export let toggleSignUp: () => void;

  const formSchema = z.object({
    email: z.string().email().min(3),
    password: z.string().min(5),
  });

  let loginDisabled = true;
  let email = '';
  let password = '';
  let error = '';

  async function handleLogin() {
    if (loginDisabled) return;
    try {
      const token = await Surreal.Instance.signin({
        NS: 'test',
        DB: 'test',
        SC: 'allusers',
        user: email.toLowerCase(),
        pass: password,
      });
      console.info('signed in to db: ', token);
      loginStore.login(token);
    } catch (err) {
      error = err as string;
      setTimeout(() => error = '', 3000);
    }
  };

  $: {
    loginDisabled = !formSchema.safeParse({ email, password }).success
  }
</script>

<!-- Component template -->
<div class="flex justify-center items-center h-screen">
  <form class="bg-white p-6 rounded-lg flex flex-col" on:submit|preventDefault={handleLogin}>
    <h2 class="text-lg font-medium mb-4 text-center">Log In</h2>
    <div class="mb-4">
      <label class="block text-gray-700 font-medium mb-2" for="email">
        Email
      </label>
      <input
        class="border border-gray-400 p-2 rounded-lg w-full"
        type="email"
        id="email"
        bind:value={email}
        required
      >
    </div>
  
    <div class="mb-4">
      <label class="block text-gray-700 font-medium mb-2" for="password">
        Password
      </label>
      <input
        class="border border-gray-400 p-2 rounded-lg w-full"
        type="password"
        id="password"
        bind:value={password}
        required
      >
    </div>
    {#if error}
      <div class="text-red-500 text-sm mb-4">{error}</div>
    {/if}
  
    <button
      class={(loginDisabled ? "disabled" : "bg-blue-500 hover:bg-blue-600") + "text-white py-2 px-4 rounded-lg w-full mb-4"}
      disabled={loginDisabled}
    >
      Log In
    </button>
  
    <button
      class="bg-white text-gray-600 py-2 px-4 rounded-lg hover:bg-gray-200 w-full"
      on:click={toggleSignUp}
    >
      Sign Up
    </button>
  </form>
</div>

<style>
  form {
    width: 20rem;
  }
</style>
