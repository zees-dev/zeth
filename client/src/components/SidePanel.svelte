
<script lang="ts">
  export let component: any = null;
  export let show = false;
  export let onSuccessfulSubmission: () => void;

  function hidePanelAndCallback() {
    show = false;
    onSuccessfulSubmission();
  }
</script>

{#if show}
  <div class="overlay" on:click={() => { show = false }} />
  <aside class="panel">
    <svelte:component this={component} onSuccessfulSubmission={hidePanelAndCallback} />
  </aside>
{/if}

<style>
  .overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: rgba(0, 0, 0, 0.5);
    z-index: 1;
    backdrop-filter: blur(2px);
  }

  .panel {
    position: fixed;
    top: 0;
    bottom: 0;
    right: 0;
    background-color: white;
    z-index: 2;
    transition: transform 0.3s ease-out;

    width: max(25rem, 35%);
    transform: translateX(0);
    display: flex;
    padding: 1rem;
  }
</style>