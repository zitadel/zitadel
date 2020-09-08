<script context="module">
    export async function preload(page, session) {
        const {params} = page;
    }
</script>

<script>
    import Split from "../components/Split.svelte";
    import Section from '../components/Section.svelte';
    import { _ } from 'svelte-i18n';
</script>

<style>
    h2 {
        margin-bottom: 2rem;
    }
    .caos-back {
        position: absolute;
        top: 80px;
        height: 60vh;
        right: 0;
    }
    .section {
        width: 100%;
        margin-top: 70vh;
        min-height: 80vh;
    }
    @media screen and (min-width: 768px) {
        .caos-back {
            right: 0;
            height: 70vh;
        }
        .section {
            width: 50%;
            margin-top: 50px;
        }
    }
</style>

<svelte:head>
  <title>
    ZITADEL • Documentation
  </title>
</svelte:head>

    <img class="caos-back" src="logos/zitadel-logo-solo-darkdesign.svg" alt="caos logo">

    <Section>
        <div class="section">
            <h2>{$_('title')}</h2>
            <p>{$_('description')}</p>
            <p>{$_('description2')}</p>
            <p>{$_('description3')}</p>
            
            <ul style="list-style-type: none;">
            <li><a href="/start" >{$_('startlink')}</a></li>
            <li><a href="/integrate" >{$_('integratelink')}</a></li>
            <li><a href="/administrate" >{$_('administratelink')}</a></li>
            <li><a href="/develop" >{$_('developlink')}</a></li>
            <li><a href="/documentation" >{$_('docslink')}</a></li>
            </ul>
            
        </div>
    </Section>
