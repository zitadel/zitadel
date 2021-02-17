<script>
    import { getContext } from 'svelte';
    export let segment = null;
    export let active = false;
    export let external = null;
    export let prefetch = false;
    const current = getContext('nav');
</script>

{#if external}
	<li><a href={external}><slot></slot></a></li>
{:else}
    {#if prefetch}
	    <li class:active="{`${$current}` === segment || active }"><a rel="prefetch" sapper:prefetch href={segment}><i class="bars las la-bars"></i><slot></slot></a></li>
    {:else}
    	<li class:active="{`${$current}` === segment || active }"><a rel="prefetch"  href={segment}><i class="bars las la-bars"></i><slot></slot></a></li>
    {/if}
{/if}