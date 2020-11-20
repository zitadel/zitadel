<script>
  import { stores } from "@sapper/app";
  import Icon from "../components/Icon.svelte";
  import Nav from "../components/Nav.svelte";
  import NavItem from "../components/NavItem.svelte";
  import manifest from '../../static/manifest.json';
  import PhotoSwipe from '../components/PhotoSwipe.svelte';
  const { page } = stores();
  export let segment;
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

<main>
  <slot />
  <PhotoSwipe/>
</main>
