<script lang="ts">
  import './app.css'
  
  import Surreal, { type Auth } from 'surrealdb.js';
  import Home from './components/Home.svelte';
  import Login from './components/login/Login.svelte';
  import Signup from './components/login/Signup.svelte';
  import { loginStore } from './stores/login';

  // docker run --rm -it --name surrealdb -p 8000:8000 surrealdb/surrealdb:latest start --log trace --user admin --pass admin memory
  const db = Surreal.Instance;
  const surealdbUrl = 'http://127.0.0.1:8000/rpc';

  let signup = false;
  const toggleSignUp = () => signup = !signup;
  (window as any).db = Surreal.Instance; // for browser debugging

  // async function admin() {
  //   try {
  //     await db.connect(surealdbUrl);
  //     await db.signin({
  //       // NS: 'test',
  //       // DB: 'test',
  //       // SC: 'user',
  //       user: 'admin',
  //       pass: 'admin',
  //     });
  //     await db.use('test', 'test');

  //     console.log(db.token);

  //     // db.authenticate(token);
  //   } catch (err) {
  //     console.log('error connecting to db', err);
  //   }
  // }
  // admin();
</script>

<svelte:head>
  <title>Zeth</title>
</svelte:head>

{#await db.connect(surealdbUrl)}
  <p>Connecting to db...</p>
{:then}
  {#if $loginStore.loggedIn && $loginStore.token}
    {#await db.authenticate($loginStore.token)}
      <p>Signing in to db...</p>
    {:then}
      <Home />
    {:catch error}
      <p class="bg-error">{error}</p>
      {#if signup}
        <Signup {toggleSignUp} />
      {:else}
        <Login {toggleSignUp} />
      {/if}
    {/await}
  {:else}
    {#if signup}
      <Signup {toggleSignUp} />
    {:else}
      <Login {toggleSignUp} />
    {/if}
  {/if}
{:catch}
  <p class="bg-error">Could not connect</p>
{/await}

<style>
</style>