<script context="module">
    export async function preload(page, session) {
        const {params} = page;
    }
</script>

<script>
    import Split from "../../components/Split.svelte";
    import Section from '../../components/Section.svelte';
    import LanguageSwitcher from '../../components/LanguageSwitcher.svelte';
    import Nav from '../../components/Nav.svelte';
    import {LANGUAGES} from '../../../config.js';
    import Button from '../../components/Button.svelte';
    import {locale} from 'svelte-i18n';

    locale.set('de');
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
	<title>ZITADEL • Dokumentation</title>
    <meta property="description" content="Erfahren Sie mehr über ZITADEL, das Cloud-Native IAM aus der Schweiz" />

    <meta property="og:url" content="https://docs.zitadel.ch" />
    <meta property="og:type" content="website" />
    <meta property="og:title" content="ZITADEL • Dokumentation" />
    <meta property="og:description" content="Erfahren Sie mehr über ZITADEL, das Cloud-Native IAM aus der Schweiz" />
    <meta property="og:image"
        content="https://www.zitadel.ch/zitadel-social-preview25.png" />

    <meta name="twitter:card" content="summary">
    <meta name="twitter:site" content="@caos_ch">
    <meta name="twitter:title" content="ZITADEL • Dokumentation" />
    <meta name="twitter:description" content="Erfahren Sie mehr über ZITADEL, das Cloud-Native IAM aus der Schweiz" />
    <meta name="twitter:image" content="https://www.zitadel.ch/zitadel-social-preview25.png">
</svelte:head>
    <Section>
        <div class="section">
            <div class="left">
                <blockquote>
                    <p>Diese Website ist auch verfügbar in:
                        <a href="/">English'</a>
                    </p>
                </blockquote>
                <h2>ZITADEL Identity and Access Management</h2>
                <p>Die meisten Applikationen müssen die Identität eines Benutzers kennen, um Benutzerdaten sicher zu speichern und auf allen Geräten des Benutzers dieselbe personalisierte Erfahrung zu bieten.</p>
                <p>Die Authentifizierung von ZITADEL bietet Backend-Dienste, benutzerfreundliche SDKs und vorgefertigte UI-Bibliotheken zur Authentifizierung von Benutzern in Ihrer Applikation. Es unterstützt die Authentifizierung mithilfe von Passwörtern und bietet weitere Sicherheitsmechanismen, wie einen zweiten Faktor beispielsweise über OTP, um einen sicheren Zugang zu gewährleisten.</p>
                <p>ZITADEL nutzt Industriestandards wie OAuth 2.0 und OpenID Connect, die leicht in Ihr benutzerdefiniertes Backend integriert werden können.</p>
            </div>
            <img class="caos-back" src="logos/zitadel-logo-solo-darkdesign.svg" alt="caos logo">
        </div>

        <div class="subsection">
            <h3>Erfahren Sie hier, wie Sie starten</h3>
            <p>Befolgen Sie diese Anleitung, um mit ZITADEL im Allgemeinen zu beginnen.</p>

            <a href="/de/start" class="getstarted-btn"> 
                <Button selected style="margin: 2rem 0;">
                    Erfahren Sie hier, wie Sie starten
                    <i class="las la-arrow-right"></i>
                </Button>
            </a>

            <div class="sublinks">
                <a class="sublink" href="de/start#Use_ORBOS_to_install_ZITADEL">Orbos verwenden um ZITADEL zu installieren</a>
                <a class="sublink" href="de/start#Use_ORBOS_to_install_ZITADEL">Applikation erstellen</a>
            </div>
        </div>

        <div class="doc-container">
            <div class="doc">
                <div class="text">
                    <h4>Quickstarts</h4>
                    <p>Erfahren Sie, wie Sie Ihre Anwendungen integrieren und sichere Applikationen und Schnittstellen erstellen.</p>
                    <a href="/de/quickstarts">>Mehr erfahren><i class="las la-arrow-right"></i></a>
                </div>
                <img src="img/develop2.png" alt="Develop" />
            </div>

            <div class="doc">
                <div class="text">
                    <h4>Administrieren</h4>
                    <p>Erfahren Sie, wie Sie Ihre Daten und Rollenzuordnungen in ZITADEL verwalten</p>
                    <a href="/de/administrate">Mehr erfahren<i class="las la-arrow-right"></i></a>

                    <p class="section">In dieser Sektion</p>
                    
                    <div class="sectionlinks">
                        <a class="link" href="de/administrate#Organisationen">Organisationen bearbeiten</a>
                        <a class="link" href="de/administrate#Projekte">Projekte bearbeiten</a>
                    </div>
                </div>
                <img src="img/projects2.png" alt="Develop" />
            </div>

            <div class="doc">
                <div class="text">
                    <h4>ZITADEL APIs</h4>
                    <p>Erfahren Sie mehr über die ZITADEL-APIs und wie Sie mit ihnen entwickeln können.</p>
                    <a href="de/apis">Mehr erfahren<i class="las la-arrow-right"></i></a>

                    <p class="section">In dieser Sektion</p>

                    <div class="sectionlinks">
                        <a class="link" href="de/apis#Authentication_API">Authentication API</a>
                        <a class="link" href="de/apis#Management_API">Management API</a>
                        <a class="link" href="de/apis#Admin_API">Admin API</a>
                    </div>
                </div>
            </div>

            <div class="doc">
                <div class="text">
                    <h4>Architektur und Technologien</h4>
                    <p>Erfahren Sie mehr über Konstruktion- und Designprinzipien und eingesetzte Technologien</p>
                    <a href="de/architecture" >Mehr erfahren<i class="las la-arrow-right"></i></a>

                    <p class="section">In dieser Sektion</p>
                  
                    <div class="sectionlinks">
                        <a class="link" href="de/architecture#Prinzipien">Prinzipien</a>
                        <a class="link" href="de/architecture#Architektur">Architektur</a>
                        <a class="link" href="de/architecture#OpenID_Connect_1_0_and_OAuth_2_0">OIDC & OAuth 2.0</a>
                    </div>
                </div>
            </div>

            <div class="doc">
                <div class="text">
                    <h4>Benutzer Handbuch</h4>
                    <p>Erfahren Sie hier, wie Sie als Benutzer ZITADEL verwenden können.</p>
                    <a href="de/use" >Mehr erfahren<i class="las la-arrow-right"></i></a>
                </div>

                <img src="img/personal2.png" alt="Develop" />
            </div>
        </div>
    </Section>
