<script>
    import Section from '../components/Section.svelte';
    import Button from '../components/Button.svelte';
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

    .doc-container .doc .text .logo-cloud {
        display: flex;
        align-items: center;
        margin-bottom: 1rem;
    }

     .doc-container .doc .text .logo-cloud img {
         max-width: 60px;
         max-height: 60px;
         margin: .5rem;
         object-fit: contain;
         border-radius: 4px;
     }

    .doc-container .doc a {
        display: block;
        font-size: 1.7rem;
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

    @media screen and (min-width: 899px) {
        .section {
            display: flex;
            flex-direction: row;
            flex-wrap: wrap;
        }

        .section .left {
            flex: 2;
        }

        .section .caos-back {
            flex: 1;
        }

        .section .caos-back {
            display: block;
            margin: 50px 0 20px 20px;
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
            flex: 1 0 auto;
            max-height: 350px;
            transition: box-shadow .2 ease;
        }

        .doc-container .doc:hover {
            box-shadow: 0 1px 8px rgba(0,0,0,0.2);
        }
    }
</style>

<svelte:head>
  <title>
    ZITADEL • Documentation
  </title>
	<title>ZITADEL • Documentation</title>
    <meta property="description" content="Learn more about the ZITADEL, the cloud native IAM from switzerland." />

    <meta property="og:url" content="https://docs.zitadel.ch" />
    <meta property="og:type" content="website" />
    <meta property="og:title" content="ZITADEL • Documentation" />
    <meta property="og:description" content="Learn more about the ZITADEL, the cloud native IAM from switzerland." />
    <meta property="og:image"
        content="https://www.zitadel.ch/zitadel-social-preview25.png" />

    <meta name="twitter:card" content="summary">
    <meta name="twitter:site" content="@caos_ch">
    <meta name="twitter:title" content="ZITADEL • Documentation" />
    <meta name="twitter:description" content="Learn more about the ZITADEL, the cloud native IAM from switzerland." />
    <meta name="twitter:image" content="https://www.zitadel.ch/zitadel-social-preview25.png">
</svelte:head>
    <Section>
        <div class="section">
            <div class="left">
                <h2>ZITADEL Identity and Access Management</h2>
                <p>Most applications need to know the identity of a user allowing to securely store user data in the cloud and provide the same personalised experience across all of the user's devices.</p>
                <p>ZITADEL's authentication provides backend services, easy-to-use SDKs, and ready-made UI libraries to authenticate users in your application. It supports authentication using passwords and applies additional security with the help of a second factor, for example OTP, to ensure a safe and secure access.</p>
                <p>ZITADEL's authentication leverages industry standards like OAuth 2.0 and OpenID Connect so that it can be easily integrated in your custom backend.</p>
            </div>
            <img class="caos-back" src="logos/zitadel-logo-solo-darkdesign.svg" alt="caos logo">
        </div>

        <div class="subsection">
            <p>Follow this guide to get started with ZITADEL in general.</p>

            <a href="start" class="getstarted-btn"> 
                <Button selected style="margin: 2rem 0; padding: 16px 25px; height: auto;">
                    Start using ZITADEL
                    <i class="las la-arrow-right"></i>
                </Button>
            </a>
            <div class="sublinks">
                <a class="sublink" href="start#ORBOS">Use ORBOS to install ZITADEL</a>
                <a class="sublink" href="start#Integrating_an_application">Create application</a>
            </div>
        </div>

        <div class="doc-container">
            <div class="doc">
                <div class="text">
                    <h4>Quickstarts</h4>
                    <p>Learn how to integrate your applications and build secure workflows and APIs with ZITADEL</p>

                    <div class="logo-cloud">
                        <img src="tech/react.png" style="margin: -.5rem;" alt="react" />
                        <img src="tech/angular.svg" alt="angular" />
                        <img style="padding: 1rem;" src="tech/flutter.svg" alt="flutter" />
                        <img src="tech/golang.svg" alt="golang" />
                        <img style="padding: .5rem" src="tech/net.svg" alt="aspnet" />
                    </div>
                    <a href="quickstarts">Learn more<i class="las la-arrow-right"></i></a>

                    <p class="section">In this section</p>
                    <div class="sectionlinks">
                        <a class="link" href="angular">Single Page Applications (Angular)</a>
                        <a class="link" href="go">Serverside rendered Applications</a>
                        <a class="link" href="dart">Native Applications (Flutter)</a>
                    </div>
                </div>
            </div>

            <div class="doc">
                <div class="text">
                    <h4>Administration</h4>
                    <p>Learn how to manage your data and role associations in ZITADEL</p>
                    <a href="administrate">Learn more<i class="las la-arrow-right"></i></a>

                    <p class="section">In this section</p>
                    
                    <div class="sectionlinks">
                        <a class="link" href="administrate#Organizations">Manage Organisations</a>
                        <a class="link" href="administrate#Projects">Manage Projects</a>
                    </div>
                </div>
            </div>

            <div class="doc">
                <div class="text">
                    <h4>APIs</h4>
                    <i class="las la-code" style="font-size: 100px; height: 100px; color:#6c8eef" ></i>
                    <p>Learn more about our APIs and how to integrate them in your apps.</p>
                    <a href="apis">Learn more<i class="las la-arrow-right"></i></a>

                    <p class="section">In this section</p>
                    
                    <div class="sectionlinks">
                        <a class="link" href="apis#Authentication_API">Authentication API</a>
                        <a class="link" href="apis#Management_API">Management API</a>
                        <a class="link" href="apis#Admin_API">Admin API</a>
                    </div>
                
                </div>
            </div>

            <div class="doc">
                <div class="text">
                    <h4>Architecture and Technologies</h4>
                    <i class="las la-project-diagram" style="font-size: 100px; height: 100px; color:#6c8eef"></i>
                    <p>Learn more about engineering and design principles, ZITADELs architecture and used technologies.</p>
                    <a href="architecture" sapper:prefetch>Learn more<i class="las la-arrow-right"></i></a>

                    <p class="section">In this section</p>
                    <div class="sectionlinks">
                        <a class="link" href="architecture#Principles">Principles</a>
                        <a class="link" href="architecture#Architecture">Architecture</a>
                        <a class="link" href="architecture#OpenID_Connect_1_0_and_OAuth_2_0">OIDC & OAuth 2.0</a>
                    </div>
                </div>
            </div>

            <div class="doc">
                <div class="text">
                    <h4>User Manuals</h4>
                    <i class="las la-book-reader" style="font-size: 100px; height: 100px; color:#6c8eef"></i>

                    <p>Follow this guide to get started with ZITADEL as a user.</p>
                    <a href="use" sapper:prefetch>Learn more<i class="las la-arrow-right"></i></a>
                </div>
            </div>
        </div>
    </Section>
