<script>
    import LanguageSwitcher from './LanguageSwitcher.svelte'
    import NavItem from './NavItem.svelte'
	import { onMount, setContext } from 'svelte';
	import { writable } from 'svelte/store';
	import Icon from './Icon.svelte';
	export let segment;
    export let logo;
    export let title;
    import { _ } from 'svelte-i18n';
	const current = writable(null);
	setContext('nav', current);
	let visible = true;
	
	let hash_changed = false;
	function handle_hashchange() {
		hash_changed = true;
	}
	let last_scroll = 0;
	function handle_scroll() {
		const scroll = window.pageYOffset;
		if (!hash_changed) {
			visible = (scroll < 50 || scroll < last_scroll);
		}
		last_scroll = scroll;
		hash_changed = false;
	}
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
		padding: 0 3rem;
		margin: 0 auto;
		box-shadow: 0 -0.4rem 0.9rem 0.2rem rgba(0,0,0,.5);
		z-index: 100;
		user-select: none;
		transform: translate(0,calc(-100% - 1rem));
        transition: transform 0.2s;
        backdrop-filter: saturate(100%) blur(10px);
	}
	header.visible {
		transform: none;
	}
	nav {
        box-sizing: border-box;
		position: fixed;
		top: 0;
		left: 0;
		width: 100vw;
		height: var(--nav-h);
		padding: 0 3rem 0 3rem;
		display: flex;
		align-items: center;
		background-color: transparent;
		transform: none;
		transition: none;
		box-shadow: none;
    }
    
    .fill-space {
        flex: 1;
    }
    
	.home {
        width: 200px;
        line-height: 22px;
        font-size: 22px;
        display: none;
    }

    .home:hover {
        color: inherit;
        text-decoration:none;
        border: none;
    }

    .home img {
        display: block;
    }
    
	a {
		color: inherit;
		border-bottom: none;
        transition: none;
    }

    @media (min-width: 400px) {
        .home {
            display: inline-block;
        }
    }

    a img {
        max-height: 40px;
    }

    .switcher-wrapper {
        padding: 0 1rem;
    }

    button {
        display: flex;
        align-items: center;
        border-radius: 8px;
        border: 1px solid hsla(0,0%,100%,.12);
        box-shadow: 0 0 0 0 rgba(0,0,0,.2), 0 0 0 0 rgba(0,0,0,.14), 0 0 0 0 rgba(0,0,0,.12);
        padding: 0 15px;
        height: 36px;
        color: var(--prime);
        transition: background-color .2 ease;
        margin: 0 1rem;
        min-width: 120px;
    }

    button:hover {
        background-color: #5282c110;
    }
    button:active {
        background-color: #5282c120;
    }

    button span {
        font-size: 14px;
        line-height: 14px;
    }
</style>

<svelte:window on:hashchange={handle_hashchange} on:scroll={handle_scroll} />

<header class:visible="{visible}">
	<nav>
		<a
			rel="prefetch"
			href="."
			class="home"
			title="{title}"
		>
            {#if logo}
                <img src={logo} alt={title} />
            {:else if title}
                {title}
            {/if}
        </a>

        <span class="fill-space"></span>

        <a href='https://console.zitadel.ch'><button>
            <span>{$_('toconsole')}</span>
        </button>
        </a>

        <NavItem external="https://github.com/caos" title="GitHub Repo">
            <Icon name="lab la-github" size="24px"></Icon>
        </NavItem>

        <div class="switcher-wrapper">
            <LanguageSwitcher></LanguageSwitcher>
        </div>
	</nav>
</header>