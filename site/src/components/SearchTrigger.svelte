<script>
    import {onMount} from 'svelte';
    let buttonname = '';

    onMount(() => {
        buttonname = navigator.platform.indexOf('Mac') > -1 ? 'Cmd': 'Ctrl';
	});
</script>

<style>
    button {
        height: 20px;
        border-radius: 8px;
        display: flex;
        align-items: center;
        color: var(--grey-text);
        margin: 1rem 2rem;
        background-color: #2a2f45;
        border: 1px solid #ffffff20;
        height: 35px;
        padding: 0 1rem;
        width: calc(100% - 4rem);
        box-sizing: border-box;
        position: fixed;
        top: 80px;
        left: 0;
    }

    button:hover {
        background-color: #373d58;
    }

    button i {
        margin-right: 1rem;
    }

    button span {
        font-size: 1.3rem;
        margin-right: 1rem;
    }

    .fill-space {
        flex: 1;
    }

    .strg {
        font-size: 1.3rem;
        margin-right: 1rem;
    }

    @media (min-width: 832px) {
        button {
            width: calc(var(--sidebar-w) - 4rem);
        }
    }
</style>

<button on:click class="search-trigger">
    <i class="las la-search"></i>
    <span>Search this site</span>
    <span class="fill-space"></span>
    <span class="strg">{buttonname || 'Ctrl'} F</span>
</button>