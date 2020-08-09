<script>
  import { stores } from "@sapper/app";
  import Icon from "../components/Icon.svelte";
  import Nav from "../components/Nav.svelte";
  import NavItem from "../components/NavItem.svelte";
  import manifest from '../../static/manifest.json';

  export let segment;
  const { page } = stores();
</script>

<style>
  main {
    position: relative;
    margin: 0 auto;
    padding: var(--nav-h) 0 0 0;
    overflow-x: hidden;
  }
</style>

<Nav {segment} {page} title="{manifest.name}">
    <NavItem external="https://github.com/caos" title="GitHub Repo">
        <Icon name="lab la-github" size="24px"></Icon>
    </NavItem>
</Nav>

<main>
  <slot />
</main>
