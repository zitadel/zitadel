<script>
	import { getContext } from 'svelte';
	export let segment = null;
	export let external = null;
	const current = getContext('nav');
</script>

{#if external}
	<li><a href={external}><slot></slot></a></li>
{:else}
	<li class:active="{$current === segment}"><a rel="prefetch" href={segment}><i class="bars las la-bars"></i><slot></slot></a></li>
{/if}