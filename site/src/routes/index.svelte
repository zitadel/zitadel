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
    Caos • Documentation
  </title>
</svelte:head>

    <img class="caos-back" src="logos/zitadel-logo-solo-darkdesign.svg" alt="caos logo">

    <Section>
        <div class="section">
            <h2>{$_('title')}</h2>
            <p>{$_('description')}</p>
            <p>{$_('description2')}</p>
            <p>{$_('description3')}</p>

            <a href="/get_started" >{$_('button')}</a>
        </div>
    </Section>
