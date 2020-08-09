<style>
    a {
        display: flex;
        align-items: center;
        color: var(--text);
        text-decoration: none;
        border: none;
        padding: 0;
    }

    a:hover {
        text-decoration: none;
        border: none;
        padding: 0;
    }
</style>

<script>
	import { getContext } from 'svelte';
	export let segment = null;
    export let external = null;
    const current = getContext('nav');
</script>

{#if external}
	<a href={external}><slot></slot></a>
{:else}
	<a rel="prefetch" href={segment}><slot></slot></a>
{/if}