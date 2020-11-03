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
    import Button from '../components/Button.svelte';

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

    .subsection .getstarted-btn {
        border-radius: 8px;
        color: white;
        padding: 0;
        border: none;
    }
    
    .doc-container {
        display: flex;
        flex-wrap: wrap;
        margin: 0 -1rem;
    }

    .doc-container .doc {
        flex-basis: 45%;
        margin: 1rem;
        padding: 1.5rem;
        background: #2a2f45;
        border-radius: 10px;
        display: flex;
        box-sizing: border-box;
    }

    .doc-container .doc img{
        width: 180px;
        margin-left: 1rem;
        object-fit: cover;
        margin-right: -1.5rem;
        margin-top: -1.5rem;
        margin-bottom: -1.5rem;
        border-top-right-radius: 10px;
        border-bottom-right-radius: 10px;
    }

    .doc-container .doc .text{
        flex: 1;
        min-width: 250px;
    }

    .doc-container .doc .text p{
        font-size: 15px;
    }

    .doc-container .doc .text h4 {
        font-size: 20px;
        margin-bottom: 1.5rem;
    }

    .doc-container .doc .text .section{
        font-weight: 700;
        text-transform: uppercase;
        font-size: 12px;
        margin: .5rem 0;
    }

    .doc-container .doc a {
        display: block;
        font-size: 1.5rem;
        border: none;
        margin-bottom: 1rem;
        padding: 0;
        font-weight: 700;
        margin-bottom: 2rem;
    }

    .doc-container .doc a i {
        color: var(--prime);
        margin-left: .5rem;
    }

    .doc-container .doc a:hover {
        padding: 0;
    }

    blockquote {
        padding: 0;
        color: var(--text-color);
    }

    blockquote p {
        font-size: 14px;
    }

    .sublinks {
        display: block;
        margin: 1rem 0;
    }

    .sublinks .sublink {
        font-size: 15px;
        margin-right: 2rem;
        border: none;
    }

    .sectionlinks {
        display: block;
        margin: 1rem 0;
    }

    .sectionlinks .link {
        font-size: 15px;
        font-size: 13px !important;
        margin-bottom: 4px !important;
        border: none;
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
                        <a href="/" on:click="{() => reload(lang)}" class="{lang == $locale ? 'current': ''}">{lang == 'de'? 'Deutsch' : 'English'}</a>
                    {/if}
                {/each}
                </p>
            </blockquote>
            <h2>{$_('title')}</h2>
            <p>{$_('description')}</p>
            <p>{$_('description2')}</p>
            <p>{$_('description3')}</p>

        </div>

        <div class="subsection">
            <!-- <h3>{$_('subheader_title')}</h3>
            <p>{$_('subheader_description')}</p> -->

            <h3>{$_('startlink')}</h3>
            <p>{$_('startlink_desc')}</p>

            <a href="/start" class="getstarted-btn"> 
                <Button selected style="margin: 2rem 0;">
                    {$_('startlink')}
                    <i class="las la-arrow-right"></i>
                </Button>
            </a>

            <p class="section">{$_('inthissection')}</p>
            <div class="sublinks">
                <a class="sublink" href="start#Use_ORBOS_to_install_ZITADEL">{$_('startlink_useorbos')}</a>
                <a class="sublink" href="start#Use_ORBOS_to_install_ZITADEL">{$_('startlink_setupapp')}</a>
            </div>
        </div>

        <div class="doc-container">
            <div class="doc">
                <div class="text">
                    <h4>{$_('integratelink')}</h4>
                    <p>{$_('integratelink_desc')}</p>
                    <a href="/integrate">{$_('learnmore')}<i class="las la-arrow-right"></i></a>

                    <p class="section">{$_('inthissection')}</p>
                    <div class="sectionlinks">
                        <a class="link" href="integrate#Single_Page_Application">{$_('integratelink_spa')}</a>
                        <a class="link" href="integrate#Server_Side_Application">{$_('integratelink_ssr')}</a>
                        <a class="link" href="integrate#Mobile_App_Native_App">{$_('integratelink_nativeapp')}</a>
                    </div>
                </div>
                <img src="img/developcropped.png" alt="Develop" />
            </div>

            <div class="doc">
                <div class="text">
                    <h4>{$_('administratelink')}</h4>
                    <p>{$_('administratelink_desc')}</p>
                    <a href="/administrate">{$_('learnmore')}<i class="las la-arrow-right"></i></a>

                    <p class="section">{$_('inthissection')}</p>
                    <div class="sectionlinks">
                        <a class="link" href="administrate#Organisations">{$_('administratelink_orgs')}</a>
                        <a class="link" href="administrate#Projects">{$_('administratelink_projects')}</a>
                    </div>
                </div>
                <img src="img/projects.png" alt="Develop" />
            </div>

            <div class="doc">
                <div class="text">
                    <h4>{$_('developlink')}</h4>
                    <p>{$_('developlink_desc')}</p>
                    <a href="/develop">{$_('learnmore')}<i class="las la-arrow-right"></i></a>

                    <p class="section">{$_('inthissection')}</p>
                    <div class="sectionlinks">
                        <a class="link" href="develop#Authentication_API">{$_('developlink_authapi')}</a>
                        <a class="link" href="develop#Management_API">{$_('developlink_mgmtapi')}</a>
                        <a class="link" href="develop#Admin_API">{$_('developlink_adminapi')}</a>
                    </div>
                </div>
            </div>

            <div class="doc">
                <div class="text">
                    <h4>{$_('docslink')}</h4>
                    <p>{$_('docslink_desc')}</p>
                    <a href="/documentation" >{$_('learnmore')}<i class="las la-arrow-right"></i></a>

                    <p class="section">{$_('inthissection')}</p>
                    <div class="sectionlinks">
                        <a class="link" href="documentation#Principles">{$_('docslink_principles')}</a>
                        <a class="link" href="documentation#Architecture">{$_('docslink_architecture')}</a>
                        <a class="link" href="documentation#OpenID_Connect_1_0_and_OAuth_2_0">{$_('docslink_oidc')}</a>
                    </div>
                </div>
            </div>

            <div class="doc">
                <div class="text">
                    <h4>{$_('uselink')}</h4>
                    <p>{$_('uselink_desc')}</p>
                    <a href="/use" >{$_('learnmore')}<i class="las la-arrow-right"></i></a>
                </div>
                <img src="img/usermanual.png" alt="Develop" />
            </div>
        </div>
    </Section>
