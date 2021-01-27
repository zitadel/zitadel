<script>
	import { onMount, setContext } from 'svelte';
	import { writable } from 'svelte/store';
	import Icon from './Icon.svelte';
	export let segment;
    export let page;
    export let title;
	export let logo;
    // export let home = 'Home';
    import { _ } from 'svelte-i18n';

	const current = writable(null);
	setContext('nav', current);
	let open = false;
	let visible = true;
    // hide nav whenever we navigate
    // if (page) {
        page.subscribe(() => {
            open = false;
        });
    // }
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
		padding: 0 1rem;
		margin: 0 auto;
		box-shadow: 0 -0.4rem 0.9rem 0.2rem rgba(0,0,0,.5);
		z-index: 100;
		user-select: none;
		transform: translate(0,calc(-100% - 1rem));
        transition: transform 0.2s;
        backdrop-filter: saturate(100%) blur(10px);
        -webkit-backdrop-filter: saturate(100%) blur(10px);
        background-color: #1b153030;
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
		padding: 0 1rem 0 1rem;
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
        display: flex;
        align-items: center;
        line-height: 22px;
        font-size: 22px;
    }

    .home .docs {
        color: var(--second);
        margin-left: 3px;
    }

    a img {
        width: 160px;
        max-height: 45px;
        padding: 0;
    }

	.primary {
		list-style: none;
		margin: 0;
		line-height: 1;
	}

	ul :global(li),
    ul :global(.login-button) {
		display: block;
        display: none;
        font-size: 14px;
        line-height: 14px;
	}

	ul :global(li).active {
		display: block;
	}

	ul {
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
		padding: 0 0 1em 0;
        background-color: #20244e;
		align-self: start;
        border-bottom-left-radius: .5rem;
        border-bottom-right-radius: .5rem;
        box-shadow: 0 5px 10px rgba(0, 0, 0, .12);
	}

	ul.open :global(li),
    ul.open :global(.login-button){
		display: block;
		text-align: right;
	}

	ul.open::after {
		display: none;
    }

	ul :global(li) :global(a) {
        font-weight: 500;
        padding: 1rem 1rem;
		border: none;
        color: inherit;
        text-decoration: none;
	}
	ul.open :global(li) :global(a) {
        padding: 1.5rem 1rem;
		display: block;
    }

    ul :global(li).active :global(a){
        display: flex;
        align-items: center;
        justify-content: flex-end;
    }

    ul :global(li).active :global(.bars) {
        display: inline-block;
    }

    ul.open :global(li) :global(.bars) {
        display: none;
    }
    
	.primary :global(svg) {
		width: 2rem;
		height: 2rem;
    }
    
	.home {
        width: 160px;
    }

    .home img {
        display: block;
    }
    
     /* TODO remove global color */
	ul :global(li).active :global(a) {
        color: var(--second); 
        white-space: nowrap;
    }

    ul :global(i) {
        color: var(--second);
    }
    
	.modal-background {
		position: fixed;
		width: 100%;
		height: 100%;
		left: 0;
		top: 0;
		/* background-color: rgba(255, 255, 255, 0.9); */
    }
    
	a {
		color: inherit;
		border-bottom: none;
        transition: none;
    }

    .show-if-desktop {
        display: none;
    }

	@media (min-width: 1040px) {
		ul {
			padding: 0;
            background: none;
        }

        ul :global(.login-button) {
            display: inline;
        }
        
		ul.open {
			padding: 0;
            background-color: transparent;
			border: none;
			align-self: initial;
        }
        
		ul.open :global(li),
        ul.open :global(.login-button){
			display: inline;
			text-align: left;
        }
        
		ul.open :global(li) :global(a) {
            display: inline;
            white-space: nowrap;
        }

        ul :global(li).active :global(.bars) {
            display: none;
        }

        ul :global(li).active :global(a){
            display: inline;
        }
        
		ul::after {
			display: none;
        }
        
		ul :global(li) {
            display: inline !important;
            white-space: nowrap;
            vertical-align: middle;
        }

        ul :global(li) :global(.bars) {
            display: none;
        }
        
		.hide-if-desktop {
			display: none !important;
        }

        .show-if-desktop {
            display: block;
        }
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



</style>

<svelte:window on:hashchange={handle_hashchange} on:scroll={handle_scroll} />

<header class:visible="{visible || open}">
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
            <span class="docs">DOCS</span>
        </a>

		{#if open}
			<div class="modal-background hide-if-desktop" on:click="{() => open = false}"></div>
		{/if}

        <span class="fill-space"></span>

		<ul
			class="primary"
			class:open
			on:touchstart|capture={intercept_touchstart}
			on:mouseenter="{() => open = true}"
			on:mouseleave="{() => open = false}">
			<!-- <li class:active="{!segment}"><a rel="prefetch" href="."><i class="bars las la-bars"></i>{home}</a></li> -->
			<slot></slot>

            <!-- <a class="show-if-desktop" href='https://console.zitadel.ch'>
                <button>
                    <span>{$_('login')}</span>
                </button>
            </a> -->
		</ul>

        <span class="fill-space show-if-desktop"></span>

        <a class="show-if-desktop" href='https://console.zitadel.ch'>
            <button>
                <span>{$_('login')}</span>
            </button>
        </a>

        <!-- <a href='https://accounts.zitadel.ch/register'>
            <button style="border-color: var(--second); margin-left: 0;">
                <span>{$_('register')}</span>
            </button>
        </a> -->
   
        <!-- <a class="login-button" rel=prefetch href='https://console.zitadel.ch' target="_blank"><button>{$_('nav_login')}</button></a> -->

	</nav>
</header>