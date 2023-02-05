<script lang="ts">
  import Router, { replace, link, location } from "svelte-spa-router";
  import active from "svelte-spa-router/active";
  // import {wrap} from 'svelte-spa-router/wrap'
  import routes from '../lib/routes';
  import { loginStore } from "../stores/login";

  $: $location === "/" ? replace("/endpoints") : null; // redirect '/' to /endpoints
  $: routeSegments = $location.split('/').slice(1); // /endpoints/1 -> ['endpoints', '1']
</script>

<main>
  <a href="/" class="logo logo-sm" use:link>
    <img
      src="{import.meta.env.BASE_URL}images/logo.svg"
      alt="Zeth logo"
      width="40"
      height="40"
    />
  </a>

  <div class="heading">
    <span class="heading-route">
      {#each routeSegments as segment, i}
        <a
          style="text-decoration: underline;"
          href={'/' + routeSegments.slice(0, i+1).join('/')}
          use:link
        >
          {segment}
        </a>
        {#if i < routeSegments.length - 1}
          {" > "}
        {/if}
      {/each}
    </span>

    <img
      class="m-auto"
      width="20"
      src={`https://api.dicebear.com/5.x/bottts-neutral/svg?seed=${$loginStore.userId}&radius=25`}
      alt="avatar"
    />

    <div>
      notifications
    </div>

    <button on:click={loginStore.logout}>logout</button>
  </div>

  <aside class="sidebar">
    <a
      href="/endpoints"
      class="menu-item"
      aria-label="Endpoints"
      use:link
      use:active={{ path: "/endpoints/?.*", className: "current-route" }}
    >
      Endpoints
      <i class="ri-database-2-line" />
    </a>
    <a
      href="/accounts"
      class="menu-item"
      aria-label="Accounts"
      use:link
      use:active={{ path: "/accounts/?.*", className: "current-route" }}
    >
      Accounts
      <i class="ri-user-2-line" />
    </a>
    <a
      href="/contracts"
      class="menu-item"
      aria-label="Contracts"
      use:link
      use:active={{ path: "/contracts/?.*", className: "current-route" }}
    >
      Contracts
      <i class="ri-file-list-2-line" />
    </a>
    <a
      href="/settings"
      class="menu-item"
      aria-label="Settings"
      use:link
      use:active={{ path: "/settings/?.*", className: "current-route" }}
    >
      Settings
      <i class="ri-settings-2-line" />
    </a>
  </aside>

  <Router {routes} />
</main>

<!-- <Router {url}>
	<main class="main">
		<SideBar class="sidebar" />
		<NavBar />
		<section class="content overflow-x-auto">
			<Route path="/"><Dashboard /></Route>
			<Route path="node/:id/*" let:params>
				<NodeRouter id={params.id} />
			</Route>
			<Route path="settings"><Settings /></Route>
		</section>
	</main>
</Router> -->

<style>
  main {
		color: var(--text);
		background: var(--bg);
		transition: background 500ms ease-in-out, color 1000ms ease-in-out;

		margin: 0 auto;
		padding: 0;
		width: 100vw;
		height: 100vh;

		display: grid;
		grid-template: 1fr 15fr/ 5rem auto;
	}

  .logo {
    margin: auto;
  }

  .heading {
    display: grid;
    grid-auto-flow: column;
    align-content: center;
    grid-template-columns: 1fr auto;
    padding: 0rem 1rem;
    grid-gap: 0.5rem;
  }
</style>