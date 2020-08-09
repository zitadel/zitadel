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
    export let tags;
</script>

<svelte:head>
  <title>{manifest.name} • {slug}</title>    

    {#each tags as { name, content }, i}
     <meta name={name} content={content} />
	{/each}
</svelte:head>

<Docs {sections} project="site" dir="{slug}"/>