<script>
  import { stores } from "@sapper/app";
  import Icon from "../components/Icon.svelte";
  import Nav from "../components/Nav.svelte";
  import NavItem from "../components/NavItem.svelte";
  import manifest from '../../static/manifest.json';
  import PhotoSwipe from '../components/PhotoSwipe.svelte';
  const { page } = stores();
  export let segment;
  import { _ } from 'svelte-i18n';

</script>

<script context="module">
  import { waitLocale } from 'svelte-i18n';
  export async function preload(page) {
    return waitLocale();
  }
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
    <NavItem segment="start">{$_('nav_start')}</NavItem>
    <NavItem segment="quickstarts">{$_('nav_quickstarts')}</NavItem>
    <NavItem segment="administrate">{$_('nav_administrate')}</NavItem>
    <NavItem segment="apis">{$_('nav_apis')}</NavItem>
    <NavItem segment="architecture">{$_('nav_architecture')}</NavItem>
    <NavItem segment="use">{$_('nav_use')}</NavItem>
</Nav>

<main>
  <slot />
  <PhotoSwipe/>
</main>
