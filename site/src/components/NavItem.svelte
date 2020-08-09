<script>
	import { getContext } from 'svelte';
	export let segment = null;
    export let external = null;
    export let lang = false;
    const current = getContext('nav');
</script>

{#if external}
	<li><a href={external}><slot></slot></a></li>
{:else}
	<li class:active="{$current === segment}" class:lang="{lang}"><a rel="prefetch" href={segment}><slot></slot></a></li>
{/if}