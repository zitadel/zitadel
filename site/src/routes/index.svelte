<script context="module">
    export async function preload(page, session) {
        const {params} = page;
    }
</script>

<script>
    import Split from "../components/Split.svelte";
    import Section from '../components/Section.svelte';
    import { _ , locale} from 'svelte-i18n';
    import LanguageSwitcher from '../components/LanguageSwitcher.svelte';
    import Nav from '../components/Nav.svelte';
    export let segment;
    import {LANGUAGES} from '../../config.js';

    function reload(language) {
        if (typeof window !== 'undefined') {
            locale.set(language);
        }
    }
</script>

<style>
    h2 {
        margin-bottom: 2rem;
    }
    .caos-back {
        position: absolute;
        top: calc(2rem + var(--nav-h));
        height: 30vh;
        max-height: 500px;
        right: 2rem;
    }
    .section {
        width: 100%;
        margin-top: 40vh;
    }
    @media screen and (min-width: 768px) {
        .caos-back {
            right: 50px;
            height: 60vh;
        }
        .section {
            width: 50%;
            margin: 50px 0;
        }
    }

    .doc-container {
        display: flex;
        flex-wrap: wrap;
        margin: 0 -1rem;
    }

    .doc-container .doc {
        flex: 1;
        min-width: 270px;
        margin: 1rem;
        padding: 2rem;
        background: #2a2f45;
        border-radius: 10px;
    }

    .doc-container .doc a {
        display: block;
        font-size: 2rem;
        border: none;
        margin-bottom: 1rem;
        padding: 0;
    }

    .doc-container .doc a:hover {
        padding: 0;
    }

    blockquote {
        color: white;
        background: #2a2f45; 
    }

    blockquote button {
        color: var(--prime);
    }
</style>

<svelte:head>
  <title>
    ZITADEL • Documentation
  </title>
</svelte:head>
    <Nav {segment} title="Zitadel docs" logo="logos/zitadel-logo-light.svg"></Nav>
    <img class="caos-back" src="logos/zitadel-logo-solo-darkdesign.svg" alt="caos logo">

    <Section>
        <div class="section">
        <blockquote>
            <p>This site is also available in: 
            {#each LANGUAGES as lang}
                {#if lang != $locale}
                    <button on:click="{() => reload(lang)}" class="{lang == $locale ? 'current': ''}">{lang == 'de'? 'Deutsch' : 'English'}</button>
                {/if}
            {/each}
            </p>
        </blockquote>
            <h2>{$_('title')}</h2>
            <p>{$_('description')}</p>
            <p>{$_('description2')}</p>
            <p>{$_('description3')}</p>
        </div>

        <div class="doc-container">
            <div class="doc">
                <a href="/start" >{$_('startlink')}</a>
                <p>{$_('startlink_desc')}</p>
            </div>

            <div class="doc">
                <a href="/integrate">{$_('integratelink')}</a>
                <p>{$_('integratelink_desc')}</p>
            </div>

            <div class="doc">
                <a href="/administrate" >{$_('administratelink')}</a>
                <p>{$_('administratelink_desc')}</p>
            </div>

            <div class="doc">
                <a href="/develop" >{$_('developlink')}</a>
                <p>{$_('developlink_desc')}</p>
            </div>

            <div class="doc">
                <a href="/documentation" >{$_('docslink')}</a>
                <p>{$_('docslink_desc')}</p>
            </div>

            <div class="doc">
                <a href="/use" >{$_('uselink')}</a>
                <p>{$_('uselink_desc')}</p>
            </div>
        </div>
    </Section>
