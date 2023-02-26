<script lang="ts">
  export let type: 'http' | 'ws'
  export let text: string | undefined = type
</script>

{#if type === 'http'}
  <div class="{$$props.class} indicator-item badge badge-secondary backdrop-brightness-50 text-xs">{text}</div>
{:else if type === 'ws'}
  <div class="{$$props.class} indicator-item badge badge-accent text-xs">{text}</div>
{/if}
