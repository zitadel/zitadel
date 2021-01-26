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

    locale.subscribe(l => {
        console.log(l);
    })
</script>

<style>
    h2 {
        margin-bottom: 2rem;
    }
    
    .section {
        width: 100%;
        margin-top: 100px;
    }

    .section .caos-back {
        display: none;
    }

    .subsection {
        display: block;
        margin-bottom: 50px;
    }

    .subsection .getstarted-btn {
        border-radius: 8px;
        color: white;
        padding: 0;
        border: none;
    }

    .doc-container .doc {
        padding: 1.5rem;
        background: #2a2f45;
        border-radius: 10px;
        box-sizing: border-box;
        display: block;
        margin-bottom: 2rem;
    }

    .doc-container .doc img{
        display: none;
    }

    .doc-container .doc .text{
        top: 180px;
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
        white-space: nowrap;
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
        white-space: nowrap;
    }

    @media screen and (min-width: 768px) {
        .section {
            display: flex;
            flex-direction: row;
            flex-wrap: wrap;
        }

        .section .left,
        .section .caos-back {
            flex: 1;
        }

        .section .caos-back {
            display: block;
            margin: 50px;
            max-width: 400px;
        }

        .doc-container {
            display: flex;
            flex-wrap: wrap;
            margin: 0 -1rem;
        }

        .doc-container .doc {
            display: flex;
            margin: 1rem;
            max-width: 500px;
            flex: 1 0 auto;
            max-height: 350px;
            transition: box-shadow .2 ease;
        }

        .doc-container .doc:hover {
            box-shadow: 0 1px 8px rgba(0,0,0,0.2);
        }

        .doc-container .doc img{
            display: block;
            width: 180px;
            height: auto;
            margin-left: 1rem;
            object-fit: cover;
            margin-right: -1.5rem;
            margin-top: -1.5rem;
            margin-bottom: -1.5rem;
            border-top-right-radius: 10px;
            border-bottom-right-radius: 10px;
            object-position: 0 0;
        }
    }
</style>

<svelte:head>
  <title>
    ZITADEL • Documentation
  </title>
	<title>{$_('title')}</title>
    <meta property="description" content="{$_('home_seo.description')}" />

    <meta property="og:url" content="https://docs.zitadel.ch" />
    <meta property="og:type" content="website" />
    <meta property="og:title" content="{$_('home_seo.title')}" />
    <meta property="og:description" content="{$_('home_seo.description')}" />
    <meta property="og:image"
        content="https://www.zitadel.ch/zitadel-social-preview25.png" />

    <meta name="twitter:card" content="summary">
    <meta name="twitter:site" content="@caos_ch">
    <meta name="twitter:title" content="{$_('home_seo.title')}" />
    <meta name="twitter:description" content="{$_('home_seo.description')}" />
    <meta name="twitter:image" content="https://www.zitadel.ch/zitadel-social-preview25.png">
</svelte:head>
    <Section>
        <div class="section">
            <div class="left">
                <blockquote>
                
                    <p>{$_('languagealsoavailable')}:
                    {#each LANGUAGES as lang}
                        {#if lang != $locale}
                            <a href="/" on:click="{() => reload(lang)}" class="{lang == $locale ? 'current': ''}">{lang == 'de'? 'Deutsch (WIP)' : 'English'}</a>
                        {/if}
                    {/each}
                    </p>
                </blockquote>
                <h2>{$_('title')}</h2>
                <p>{$_('description')}</p>
                <p>{$_('description2')}</p>
                <p>{$_('description3')}</p>
            </div>
            <img class="caos-back" src="logos/zitadel-logo-solo-darkdesign.svg" alt="caos logo">
        </div>

        <div class="subsection">
            <h3>{$_('startlink')}</h3>
            <p>{$_('startlink_desc')}</p>

            <a href="/start" class="getstarted-btn"> 
                <Button selected style="margin: 2rem 0;">
                    {$_('startlink')}
                    <i class="las la-arrow-right"></i>
                </Button>
            </a>
            {#if $locale == 'en'}
                <div class="sublinks">
                    <a class="sublink" href="start#Use_ORBOS_to_install_ZITADEL">{$_('startlink_useorbos')}</a>
                    <a class="sublink" href="start#Use_ORBOS_to_install_ZITADEL">{$_('startlink_setupapp')}</a>
                </div>
            {:else if $locale == 'de'}
                <div class="sublinks">
                    <a class="sublink" href="start#Use_ORBOS_to_install_ZITADEL">{$_('startlink_useorbos')}</a>
                    <a class="sublink" href="start#Use_ORBOS_to_install_ZITADEL">{$_('startlink_setupapp')}</a>
                </div>
            {/if}
        </div>

        <div class="doc-container">
            <div class="doc">
                <div class="text">
                    <h4>{$_('integratelink')}</h4>
                    <p>{$_('integratelink_desc')}</p>
                    <a href="/integrate">{$_('learnmore')}<i class="las la-arrow-right"></i></a>

                    {#if $locale == 'en'}
                        <p class="section">{$_('inthissection')}</p>
                        <div class="sectionlinks">
                            <a class="link" href="integrate#Single_Page_Application">{$_('integratelink_spa')}</a>
                            <a class="link" href="integrate#Server_Side_Application">{$_('integratelink_ssr')}</a>
                            <a class="link" href="integrate#Mobile_App_Native_App">{$_('integratelink_nativeapp')}</a>
                        </div>
                    {/if}
                </div>
                <img src="img/develop2.png" alt="Develop" />
            </div>

            <div class="doc">
                <div class="text">
                    <h4>{$_('administratelink')}</h4>
                    <p>{$_('administratelink_desc')}</p>
                    <a href="/administrate">{$_('learnmore')}<i class="las la-arrow-right"></i></a>

                    <p class="section">{$_('inthissection')}</p>
                    
                    {#if $locale == 'en'}
                        <div class="sectionlinks">
                            <a class="link" href="administrate#Organisations">{$_('administratelink_orgs')}</a>
                            <a class="link" href="administrate#Projects">{$_('administratelink_projects')}</a>
                        </div>
                    {:else if $locale == 'de'}
                        <div class="sectionlinks">
                            <a class="link" href="administrate#Organisationen">{$_('administratelink_orgs')}</a>
                            <a class="link" href="administrate#Projekte">{$_('administratelink_projects')}</a>
                        </div>
                    {/if}
                </div>
                <img src="img/projects2.png" alt="Develop" />
            </div>

            <div class="doc">
                <div class="text">
                    <h4>{$_('developlink')}</h4>
                    <p>{$_('developlink_desc')}</p>
                    <a href="/develop">{$_('learnmore')}<i class="las la-arrow-right"></i></a>

                    <p class="section">{$_('inthissection')}</p>

                    {#if $locale == 'en'}
                       <div class="sectionlinks">
                            <a class="link" href="develop#Authentication_API">{$_('developlink_authapi')}</a>
                            <a class="link" href="develop#Management_API">{$_('developlink_mgmtapi')}</a>
                            <a class="link" href="develop#Admin_API">{$_('developlink_adminapi')}</a>
                        </div>
                    {:else if $locale == 'de'}
                        <div class="sectionlinks">
                            <a class="link" href="develop#Authentication_API">{$_('developlink_authapi')}</a>
                            <a class="link" href="develop#Management_API">{$_('developlink_mgmtapi')}</a>
                            <a class="link" href="develop#Admin_API">{$_('developlink_adminapi')}</a>
                        </div>
                    {/if}
                </div>
            </div>

            <div class="doc">
                <div class="text">
                    <h4>{$_('docslink')}</h4>
                    <p>{$_('docslink_desc')}</p>
                    <a href="/documentation" >{$_('learnmore')}<i class="las la-arrow-right"></i></a>

                    <p class="section">{$_('inthissection')}</p>
                    {#if $locale == 'en'}
                       <div class="sectionlinks">
                            <a class="link" href="documentation#Principles">{$_('docslink_principles')}</a>
                            <a class="link" href="documentation#Architecture">{$_('docslink_architecture')}</a>
                            <a class="link" href="documentation#OpenID_Connect_1_0_and_OAuth_2_0">{$_('docslink_oidc')}</a>
                        </div>
                    {:else if $locale == 'de'}
                        <div class="sectionlinks">
                            <a class="link" href="documentation#Prinzipien">{$_('docslink_principles')}</a>
                            <a class="link" href="documentation#Architektur">{$_('docslink_architecture')}</a>
                            <a class="link" href="documentation#OpenID_Connect_1_0_and_OAuth_2_0">{$_('docslink_oidc')}</a>
                        </div>
                    {/if}
                </div>
            </div>

            <div class="doc">
                <div class="text">
                    <h4>{$_('uselink')}</h4>
                    <p>{$_('uselink_desc')}</p>
                    <a href="/use" >{$_('learnmore')}<i class="las la-arrow-right"></i></a>
                </div>

                <img src="img/personal2.png" alt="Develop" />
            </div>
        </div>
    </Section>
