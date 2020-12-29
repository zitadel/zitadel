<script context="module">
    import { goto } from '@sapper/app';
    import { docLanguages } from '../modules/language-store.js'
    import {LANGUAGES} from '../../config.js';
</script>

<script>
    import { locale } from 'svelte-i18n';
    import { startClient } from '../i18n.js';

    let group= $locale;

    function reload(language) {
        if (typeof window !== 'undefined') {
            locale.set(language);
            location.reload();
        }
    }
</script>

<style>
    :root {
        --speed3: cubic-bezier(0.26, 0.48, 0.08, 0.9);
        --height: 30px;
    }

    .language-switcher {
        position: fixed;
        left: 0;
        bottom: 0;
        display: flex;
        align-items: center;
        z-index: 1;
        justify-content: center;
    }

    button {
        height: var(--height);
        margin: .5rem 1rem;
        font-size: 12px;
        display: flex;
        align-items: center;
        cursor: pointer;
        justify-content: center;
        border: none;
    }

    button.current {
        color: var(--grey-text);
    }
</style>

<div class="language-switcher">
	{#each LANGUAGES as lang}
        <button on:click="{() => reload(lang)}" disabled="{lang == group}" class="{lang == group ? 'current': ''}">{lang == 'de'? 'Deutsch' : 'English'}</button>
	{/each}
</div>
