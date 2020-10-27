<script context="module">
    export async function preload({params}) {
        const {lang, slug} = params;
        const sections = await this.fetch(`${slug}.json`).then(r => r.json());
        const tags = [];
        return { sections, slug, tags };
    }
</script>

<script>
    import manifest from '../../static/manifest.json';
    import Docs from "../components/Docs.svelte";
    export let slug;
    export let sections;
    import { onMount } from 'svelte';
    export let tags;
    import { initPhotoSwipeFromDOM } from '../utils/photoswipe.js';
    import SearchSelector from '../components/SearchSelector.svelte';

    onMount(() => {
        initPhotoSwipeFromDOM('.zitadel-gallery');
    });
</script>

<svelte:head>
  <title>{manifest.name} • {slug}</title>    

    {#each tags as { name, content }, i}
     <meta name={name} content={content} />
	{/each}
</svelte:head>

<Docs {sections} project="zitadel/site" dir="{slug}"/>

<SearchSelector></SearchSelector>