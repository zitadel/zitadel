<script context="module">
    import { setCookie } from '../modules/cookie.js';
    import { goto } from '@sapper/app';
    import { docLanguages } from '../modules/language-store.js'
    import {LANGUAGES} from '../../config.js';
</script>

<script>
    import { locale } from 'svelte-i18n';
    import { startClient } from '../i18n.js';

    let group= $locale;

    function setLocale(language) {
        if (typeof window !== 'undefined') {
            locale.set(language);
            // return goto('', {hl: language});
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
        width: var(--sidebar-w);
        justify-content: center;
    }

    a {
        height: var(--height);
        margin: .5rem 1rem;
        font-size: 12px;
        color: white;
        display: flex;
        align-items: center;
        cursor: pointer;
        justify-content: center;
        border: none;
    }

    a.current {
        color: var(--prime);
    }
</style>

<div class="language-switcher">
	{#each LANGUAGES as lang}
        <a href="&hl={lang}" class="{lang == group ? 'current': ''}">{lang == 'de'? 'Deutsch' : 'English'}</a>
	{/each}
</div>
