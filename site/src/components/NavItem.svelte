<script>
    import { getContext } from 'svelte';
    import { locale } from 'svelte-i18n';
    export let segment = null;
    export let active = false;
	export let external = null;
    const current = getContext('nav');
</script>

{#if external}
	<li><a href={external}><slot></slot></a></li>
{:else}
	<li class:active="{`${$current}` === segment || active }"><a rel="prefetch" href={segment}><i class="bars las la-bars"></i><slot></slot></a></li>
{/if}