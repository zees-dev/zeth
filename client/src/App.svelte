<script lang="ts">
  import './app.css'
  
  import Surreal, { type Auth } from 'surrealdb.js';
  import Home from './components/Home.svelte';
  import Login from './components/login/Login.svelte';
  import Signup from './components/login/Signup.svelte';
  import { loginStore } from './stores/login';
  import { dbStore } from './stores/db';

  let signup = false;
  const toggleSignUp = () => signup = !signup;
  (window as any).db = $dbStore.db; // for browser debugging

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

  $: {
    const loggedIn = $loginStore.loggedIn && $loginStore.token;
    if (!loggedIn) {
      loginStore.logout();
    };
  }
</script>

<svelte:head>
  <title>Zeth</title>
</svelte:head>

{#await dbStore.connect()}
  <p>Connecting to db...</p>
{:then}
  {#if $loginStore.loggedIn && $loginStore.token}
    {#await $dbStore.db.authenticate($loginStore.token)}
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
