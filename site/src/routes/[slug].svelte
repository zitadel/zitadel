<script context="module">
    export async function preload({params}) {
        const {lang, slug} = params;
        const {docs, seo} = await this.fetch(`${slug}.json`).then(r => r.json());
        return { sections: docs, seo, slug };
    }
</script>

<script>
    import DocsHeader from '../components/DocsHeader.svelte';
    import manifest from '../../static/manifest.json';
    import Docs from "../components/Docs.svelte";
    export let slug;
    export let sections;
    export let seo;
    import { onMount } from 'svelte';
    import { initPhotoSwipeFromDOM } from '../utils/photoswipe.js';

    onMount(() => {
        initPhotoSwipeFromDOM('.zitadel-gallery');
    });
</script>

<style>
    @media (min-width: 832px) {
        :global(main) {
            padding: 0 !important;
        }
    }
</style>

<svelte:head>
  <title>{manifest.name} • {slug}</title>    

    {#if seo}
    { @html seo}
   {/if}
</svelte:head>

<DocsHeader {slug}></DocsHeader>
<Docs {sections} dir="{slug}"/>