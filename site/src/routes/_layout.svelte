<script>
    import { stores } from "@sapper/app";
    import Nav from "../components/Nav.svelte";
    import NavItem from "../components/NavItem.svelte";
    import PhotoSwipe from '../components/PhotoSwipe.svelte';
    const { page } = stores();
    export let segment;
</script>

<style>
  main {
    position: relative;
    margin: 0 auto;
    padding: var(--nav-h) 0 0 0;
    overflow-x: hidden;
  }
</style>

<Nav title="Zitadel docs" {segment} {page} logo="logos/zitadel-logo-light.svg">
    <NavItem segment="start">Get started</NavItem>
    <NavItem segment="quickstarts">Quickstarts</NavItem>
    <NavItem segment="administrate">Administrate</NavItem>
    <NavItem segment="apis">APIs</NavItem>
    <NavItem segment="architecture">Architecture</NavItem>
    <NavItem segment="use">Manual</NavItem>
</Nav>

<main>
  <slot />
</main>

<PhotoSwipe/>