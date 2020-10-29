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
    header {
        box-sizing: border-box;
		position: fixed;
		display: flex;
		align-items: center;
		justify-content: space-between;
		width: 100vw;
		height: var(--nav-h);
		padding: 0 16px;
		margin: 0 auto;
		box-shadow: 0 -0.4rem 0.9rem 0.2rem rgba(0,0,0,.5);
		z-index: 100;
		user-select: none;
		transform: translate(0,calc(-100% - 1rem));
        transition: transform 0.2s;
        backdrop-filter: saturate(100%) blur(10px);
    }
    
    header img {
        margin: auto;
        display: block;
    }
</style>

<div>
    <header>
        <img src={logo} alt={title} />
    </header>
</div>