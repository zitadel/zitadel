<script>
    import { stores } from "@sapper/app";
    import { _ } from 'svelte-i18n';
    const { page } = stores();
    let menuOpen = false;
    export let slug;

    page.subscribe(() => {
		menuOpen = false;
    });

</script>

<style>
    .content {
        width: 100%;
        margin: 0;
        padding: 0 var(--side-nav) 0 0;
        position: fixed;
        top: 0;
        left: 0;
        right: 0;
        display: flex;
        align-items: center;
        justify-content: flex-end;
        height: var(--nav-h);
        box-shadow: 0 -0.4rem 0.9rem 0.2rem rgba(0,0,0,.5);
        backdrop-filter: saturate(100%) blur(10px);
        z-index: 1;
    }

    .content .home {
        width: var(--nav-h);
        height: var(--nav-h);
        border: none;
        object-fit: contain;
        padding: 1rem;
    }

    .content .fill-space {
        flex: 1;
    }

    .content .list-item {
        margin: 0 15px;
        border: none;
        padding: 0;
        color: var(--prime);
        font-weight: 500;
    }

    button {
        display: flex;
        align-items: center;
        margin-top: 20px !important;
        margin-bottom: 22px !important;
    }

    button.list-item i{
        color: var(--prime);
        margin-top: 2px;
        margin-right: 3px;
    }

    .content :last-child {
        margin-right: 0;
    }

    .content .btn-wrapper {
        position: relative;
    }

    .content .wrapper #menu {
        position: absolute;
        top: calc(var(--nav-h) - 5px);
        left: 50%;
        z-index: 2;
    }

    /* menu appearance*/
    #menu {
        position: relative;
        color: #999;
        width: 200px;
        padding: 10px;
        margin: auto;
        border-radius: 8px;
        background: #2a2f45;
        box-shadow: 0 1px 8px rgba(0,0,0,0.05);
        transform: translateX(-50%);
    }
    #menu:after {
        position: absolute;
        top: -10px;
        left: 85px;
        content: "";
        display: block;
        border-left: 15px solid transparent;
        border-right: 15px solid transparent;
        border-bottom: 20px solid #2a2f45;
    }

    ul, li, li a {
        list-style: none;
        display: block;
        margin: 0;
        padding: 4px 0;
        font-size: 1.5rem;
    }

    li a {
        text-decoration: none;
        border: none;
    }

    li.active a {
        color: var(--second);
    }

    .modal-background {
		position: fixed;
		width: 100%;
		height: 100%;
		left: 0;
		top: 0;
	}

    @media (min-width: 832px) {
        .content {
            position: relative;
            padding-left: calc(var(--sidebar-w) + var(--side-nav));
            box-shadow: none;
            background: none;
            backdrop-filter: none;
            z-index: auto;
        }

        .content .home {
            display: none;
        }
    }
</style>

<div class="content">
    <a class="home" href="." ><img src="logos/zitadel-logo-solo-darkdesign.svg" alt="zitadel logo" /></a>
    <span class="fill-space"></span>

    <div class="btn-wrapper">
        <button class="list-item" on:click="{() => menuOpen = !menuOpen}"><i class="las la-bars"></i><span>{$_('references')}</span></button>
        <div class="wrapper">
            {#if menuOpen}
                <div id="menu">
                    <ul>
                        <li class="{slug == 'start' ? 'active' : ''}"><a href="/start" >{$_('startlink')}</a></li>
                        <li class="{slug == 'integrate' ? 'active' : ''}"><a href="/integrate">{$_('integratelink')}</a></li>
                        <li class="{slug == 'administrate' ? 'active' : ''}"><a href="/administrate" >{$_('administratelink')}</a></li>
                        <li class="{slug == 'develop' ? 'active' : ''}"><a href="/develop" >{$_('developlink')}</a></li>
                        <li class="{slug == 'documentation' ? 'active' : ''}"><a href="/documentation" >{$_('docslink')}</a></li>
                        <li class="{slug == 'use' ? 'active' : ''}"><a href="/use" >{$_('uselink')}</a></li>
                    </ul>
                </div>
            {/if}
        </div>
    </div>

    <a class="list-item" href="https://zitadel.ch" target="_blank">More about ZITADEL</a>
    <!-- <a class="list-item" href="https://console.zitadel.ch" target="_blank">Sign in</a> -->
</div>

{#if menuOpen}
    <div class="modal-background" on:click="{() => menuOpen = false}"></div>
{/if}