<script lang="ts">
  import { onMount } from 'svelte';

  let username = '';
  let password = '';
  let passwordError = '';

  function checkPassword(password: string) {
    // perform basic checks on the password
    if (password.length < 8) {
      passwordError = 'Password must be at least 8 characters long';
      return false;
    }
    if (!/\d/.test(password)) {
      passwordError = 'Password must contain at least one number';
      return false;
    }
    if (!/[a-z]/.test(password)) {
      passwordError = 'Password must contain at least one lowercase letter';
      return false;
    }
    if (!/[A-Z]/.test(password)) {
      passwordError = 'Password must contain at least one uppercase letter';
      return false;
    }

    passwordError = '';
    return true;
  }

  function handleSignup(event: Event) {
    event.preventDefault();

    if (!checkPassword(password)) {
      return;
    }

    // make an arbitrary request to an endpoint with the username and password
  }

  function handleLogin(event) {
    event.preventDefault();

    // make an arbitrary request to an endpoint with the username and password
  }
</script>

<form class="flex flex-col items-center h-screen">
  <div class="w-full max-w-xs">
    <label for="username" class="block font-bold mb-2 text-xl text-center">Username</label>
    <input
      id="username"
      type="text"
      bind:value={username}
      class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
    />
  </div>

  <div class="w-full max-w-xs mt-6">
    <label for="password" class="block font-bold mb-2 text-xl text-center">Password</label>
    <input
      id="password"
      type="password"
      bind:value={password}
      class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
    />
    {#if passwordError}
      <p class="text-red-500 text-center mt-2">{passwordError}</p>
    {/if}
  </div>
  <div class="w-full max-w-xs mt-6">
    <button
      type="submit"
      on:click={handleSignup}
      class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
    >
      Sign Up
    </button>
  </div>

  <div class="w-full max-w-xs mt-6">
    <button
      type="submit"
      on:click={handleLogin}
      class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
    >
      Login
    </button>
  </div>
</form>