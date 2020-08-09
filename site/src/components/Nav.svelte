<script>
    import LanguageSwitcher from './LanguageSwitcher.svelte'
	import { onMount, setContext } from 'svelte';
	import { writable } from 'svelte/store';
	import Icon from './Icon.svelte';
	export let segment;
	export let page;
    export let logo;
    export let title;
	export let home_title = 'Homepage';
	const current = writable(null);
	setContext('nav', current);
	let open = false;
	let visible = true;
	// hide nav whenever we navigate
	page.subscribe(() => {
		open = false;
	});
	function intercept_touchstart(event) {
		if (!open) {
			event.preventDefault();
			event.stopPropagation();
			open = true;
		}
	}
	// Prevents navbar to show/hide when clicking in docs sidebar
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

	.primary {
		list-style: none;
		margin: 0;
		line-height: 1;
	}
	ul :global(li) {
		display: block;
        display: none;
	}
	ul :global(li).active {
		display: block;
    }

    ul :global(li).lang :global(a) {
        font-size: 16px;
    }
    
	ul {
        /* display: flex;
        align-items: center; */
        position: relative;
        text-align: center;
	}
	ul::after {
		/* prevent clicks from registering if nav is closed */
		position: absolute;
		content: '';
		width: 100%;
		height: 100%;
		left: 0;
		top: 0;
	}
	ul.open {
		padding: 3rem 1rem;
        background-color: #212224;
        align-self: start;
        border: 1px solid #ffffff;
	}
	ul.open :global(li) {
		display: block;
		text-align: right
	}
	ul.open::after {
		display: none;
	}
	ul :global(li) :global(a) {
        font-weight: 500;
        padding: 1rem .5rem;
		border: none;
        color: inherit;
        text-decoration: none;
	}
	ul.open :global(li) :global(a) {
        display: block;
        padding: .5rem .5rem;
        height: 30px;
    }
    
	.primary :global(svg) {
		width: 2rem;
		height: 2rem;
    }
    
	.home {
        width: 200px;
        line-height: 22px;
        font-size: 22px;
    }

    .home:hover {
        color: inherit;
        text-decoration:none;
        border: none;
    }

    .home img {
        display: block;
    }
    
	ul :global(li).active :global(a) {
		color: rgb(187,89,131);
	}

	.modal-background {
		position: fixed;
		width: 100%;
		height: 100%;
		left: 0;
		top: 0;
    }
    
	a {
		color: inherit;
		border-bottom: none;
        transition: none;
    }
    
	@media (min-width: 840px) {
		ul {
			padding: 0;
            background: none;
        }
        
		ul.open {
			padding: 0;
            background-color: transparent;
			border: none;
			align-self: initial;
        }
        
		ul.open :global(li) {
			display: inline;
			text-align: left;
        }
        
		ul.open :global(li) :global(a) {
			display: inline;
        }
        
		ul::after {
			display: none;
        }
        
		ul :global(li) {
            display: inline !important;
        }
        
		.hide-if-desktop {
			display: none !important;
        }
    }
    
    .hide {
        display: none;
    }

    .switcher-wrapper {
        padding: 0 1rem;
    }
</style>

<svelte:window on:hashchange={handle_hashchange} on:scroll={handle_scroll} />

<header class:visible="{visible || open}">
	<nav>
		<a
			rel="prefetch"
			href="."
			class="home"
			title="{home_title}"
		>
            {#if logo}
                <img src={logo} alt={home_title} />
            {/if}
            {#if title}
                {title}
            {/if}
        </a>

		{#if open}
			<div class="modal-background hide-if-desktop" on:click="{() => open = false}"></div>
		{/if}

        <span class="fill-space"></span>

        <div class="switcher-wrapper">
            <LanguageSwitcher></LanguageSwitcher>
        </div>

        <ul
			class="primary"
			class:open
			on:touchstart|capture={intercept_touchstart}
			on:mouseenter="{() => open = true}"
			on:mouseleave="{() => open = false}"
		>
			<slot></slot>
            <i class="hide-if-desktop las la-chevron-down" class:hide={open}></i>
		</ul>
	</nav>
</header>