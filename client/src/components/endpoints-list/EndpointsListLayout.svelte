<script lang="ts">
</script>

<div class="content">
  <nav class="sub-nav">
    <div class="search-bar">search....</div>
    <a
      href="#/"
      class="menu-item"
      aria-label="Metamask"
    >
      Metamask
      <i class="ri-database-2-line" />
    </a>
    <a
      href={`#/endpoints/1`}
      class="menu-item"
      aria-label="Node1"
    >
      Node1
      <i class="ri-database-2-line" />
    </a>
  </nav>

  <section>
    <h1>Endpoints</h1>
  </section>
</div>

<style>
  .content {
    width: 100%;
    height: 100%;

    display: grid;
    grid-template-columns: 12rem auto;
  }

  .sub-nav {
    width: 100%;
    height: 100%;

    display: flex;
    flex-flow: column;
  }
</style>