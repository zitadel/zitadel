<script context="module">
    import { setCookie } from '../modules/cookie.js';
    import { docLanguages } from '../modules/language-store.js'
    import {LANGUAGES} from '../../config.js';
</script>

<script>
    import { locale } from 'svelte-i18n';
    import { startClient } from '../i18n.js';

    let group= $locale;

    $:setLocale(group);
    function setLocale(language) {
        if (typeof window !== 'undefined') {
            setCookie('locale', language);
            startClient();
        }
    }
</script>

<style>
    :root {
        --deep-blue: #1e3470;
        --speed3: cubic-bezier(0.26, 0.48, 0.08, 0.9);
        --height: 30px;
    }

    .language-switcher {
        display: flex;
        align-items: center;
        position: relative;
    }
    
    .language-switcher input {
        appearance: none;
        display: none;
    }

    .language-switcher .select {
        height: var(--height);
        width: var(--height);
        border-radius: 50vw;
        font-size: calc(var(--height) / 2.5);
        color: #fff;
        mix-blend-mode: difference;    
        display: flex;
        align-items: center;
        cursor: pointer;
        justify-content: center;
    }

    .language-switcher .current {
        background-color: white;
        color: black;
    }
</style>

<div class="language-switcher">
	{#each LANGUAGES as lang}
		<label class="select {lang == group ? 'current' : 'notcurrent'}">
            <input type=radio bind:group value={lang}>
            {lang.toUpperCase()}
        </label>
	{/each}
</div>
