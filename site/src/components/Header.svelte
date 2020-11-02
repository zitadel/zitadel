<script>
    import LanguageSwitcher from './LanguageSwitcher.svelte'
    import NavItem from './NavItem.svelte'
	import { setContext } from 'svelte';
	import { writable } from 'svelte/store';
	import Icon from './Icon.svelte';
	export let segment;
    export let logo;
    export let title;
    import { _ } from 'svelte-i18n';
	const current = writable(null);
	setContext('nav', current);
	$: $current = segment;
</script>

<style>
    .header {
        position: fixed;
        box-sizing: border-box;
		display: flex;
        align-items: center;
        justify-content: flex-end;
        width: 100%;
		height: var(--nav-h);
        padding: 0 16px;
        top: 0;
		margin: 0 auto;
    }
    
    .header img {
        margin: auto;
        display: block;
    }
</style>

<div class="header">
    <img src={logo} alt={title} />
</div>
