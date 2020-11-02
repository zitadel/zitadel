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
        top: 0;
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
		padding: 0 16px 0 16px;
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
        line-height: 22px;
        font-size: 22px;
        display: flex;
        align-items: center;
    }

    .home:hover {
        color: inherit;
        text-decoration:none;
        border: none;
    }
    
	a {
		color: inherit;
		border-bottom: none;
        transition: none;
        padding: 0;
    }

    .home span {
        color: var(--second);
        margin-left: 3px;
    }

    a img {
        width: 160px;
        max-height: 45px;
        padding: 0;
    }

    button {
        display: flex;
        align-items: center;
        border-radius: 8px;
        border: 1px solid hsla(0,0%,100%,.12);
        box-shadow: 0 0 0 0 rgba(0,0,0,.2), 0 0 0 0 rgba(0,0,0,.14), 0 0 0 0 rgba(0,0,0,.12);
        padding: 0 15px;
        height: 36px;
        color: white;
        transition: background-color .2 ease;
        margin: 0 1rem;
        min-width: 120px;
        /* background: #2a2f45; */
    }

    button:hover {
        background-color: var(--back-hover);
    }
    button:active {
        background-color: var(--back-hover);
    }

    button span {
        font-size: 14px;
        line-height: 14px;
        text-align: center;
        margin: auto;
    }

    .show-on-desktop {
        display: none;
    }

    @media (min-width: 832px) {
        .show-on-desktop {
            display: flex;
        }
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
            <span>DOCS</span>
        </a>

        <span class="fill-space"></span>
        <div class="show-on-desktop">
            <NavItem external="https://zitadel.ch" title="GitHub Repo">
                {$_('moreabout')}
            </NavItem>

            <a href='https://console.zitadel.ch'>
                <button>
                    <span>{$_('login')}</span>
                </button>
            </a>

            <a href='https://accounts.zitadel.ch/register'>
                <button style="border-color: var(--second); margin-left: 0;">
                    <span>{$_('register')}</span>
                </button>
            </a>
        </div>
	</nav>
</header>