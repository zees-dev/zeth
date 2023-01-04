<script lang="ts">
  import './app.css'
  import Home from './components/Home.svelte';
  import Login from './components/login/Login.svelte';

  // TODO: get app state (if app ever initialized)
  // TODO: if not, show signup page
  // TODO: get user state (if user logged in - if not, redirect to login)
  let loggedin = true;
</script>

<svelte:head>
  <title>Zeth</title>
</svelte:head>

{#if loggedin}
  <Home />
{:else}
  <Login />
{/if}

<style>
</style>