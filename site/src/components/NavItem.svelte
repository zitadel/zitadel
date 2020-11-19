<style>
    a {
        display: flex;
        align-items: center;
        color: white;
        text-decoration: none;
        border: none;
        padding: 0;
        font-size: 14px;
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
    export let title = '';
    const current = getContext('nav');
</script>

{#if external}
	<a href={external}><slot></slot></a>
{:else}
	<a rel="prefetch" alt="{title}" href={segment}><slot></slot></a>
{/if}