<script context="module">
    export async function preload({params}) {
        const {language, slug} = params;
        const {docs, seo} = await this.fetch(`${language}/${slug}.json`).then(r => r.json());
        return { sections: docs, seo, slug, language };
    }
</script>

<script>
    import manifest from '../../../static/manifest.json';
    import Docs from "../../components/Docs.svelte";
    export let slug;
    export let sections;
    export let seo;
    export let language;
    import { onMount } from 'svelte';
    import { initPhotoSwipeFromDOM } from '../../utils/photoswipe.js';

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

<Docs {sections} {language} dir="{slug}"/>