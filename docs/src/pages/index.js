import React from 'react';
import clsx from 'clsx';
import Layout from '@theme/Layout';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import useBaseUrl from '@docusaurus/useBaseUrl';
import styles from './styles.module.css';
import ThemedImage from '@theme/ThemedImage';

const features = [
  {
    title: 'Manuals',
    darkImageUrl: 'img/index/Manual-dark.svg',
    lightImageUrl: 'img/index/Manual-light.svg',
    link: 'docs/manuals/introduction',
    description: (
      <>
        Follow this guide to get started with ZITADEL as a user.
      </>
    ),
  },
  {
    title: 'Quickstarts',
    darkImageUrl: 'img/index/Quickstarts-dark.svg',
    lightImageUrl: 'img/index/Quickstarts-light.svg',
    link: 'docs/quickstarts/introduction',
    description: (
      <>
        Learn how to integrate your applications and build secure workflows and APIs with ZITADEL
      </>
    ),
  },
  {
    title: 'Guides',
    darkImageUrl: 'img/index/Guides-dark.svg',
    lightImageUrl: 'img/index/Guides-light.svg',
    link: 'docs/guides/introduction',
    description: (
      <>
        Read our guides on how to manage your data and role associations in ZITADEL and on what we recommend.
      </>
    ),
  },
  {
    title: 'APIs',
    darkImageUrl: 'img/index/APIs-dark.svg',
    lightImageUrl: 'img/index/APIs-light.svg',
    link: '/docs/apis/introduction',
    description: (
      <>
        Learn more about our APIs and how to integrate them in your apps.
      </>
    ),
  },
  {
    title: 'Concepts',
    darkImageUrl: 'img/index/Concepts-dark.svg',
    lightImageUrl: 'img/index/Concepts-light.svg',
    link: 'docs/concepts/introduction',
    description: (
      <>
        Learn more about engineering and design principles, ZITADELs architecture and used technologies.
      </>
    ),
  },
];

function Feature({darkImageUrl, lightImageUrl, title, description, link}) {
  const darkImgUrl = useBaseUrl(darkImageUrl);
  const lightImgUrl = useBaseUrl(lightImageUrl);
  return (
        <div className={clsx('col col--4 docs-link', styles.feature)}>
          <Link to={useBaseUrl(link)}>
          {darkImgUrl && lightImgUrl && (
              <div className="text--center">
                <ThemedImage
                    className={styles.featureImage}
                    alt={title}
                    sources={{
                      light: lightImgUrl,
                      dark: darkImgUrl,
                    }}
                />
              </div>
          )}
          <h3>{title}</h3>
          <p>{description}</p>
          </Link>
        </div>
  );
}

export default function Home() {
  const context = useDocusaurusContext();
  const {siteConfig = {}} = context;
  return (
    <Layout
      title={`${siteConfig.title}`}
      description="This site bundles ZITADELs Documentations">
      <header className={clsx('hero hero--primary', styles.heroBanner)}>
        <div className="container">
          <h1 className="hero__title">{siteConfig.title}</h1>
          <p className="hero__subtitle">{siteConfig.tagline}</p>
          <div className={styles.buttons}>
            <Link
              className={clsx(
                'button button--outline button--lg get-started',
                styles.getStarted,
              )}
              to={useBaseUrl('docs/guides/get-started')}>
              Get Started
            </Link>
          </div>
        </div>
      </header>
      <main>
        {features && features.length > 0 && (
          <section className={styles.features}>
            <div className="container">
              <div className="row">
                {features.map((props, idx) => (
                  <Feature key={idx} {...props} />
                ))}
              </div>
            </div>
          </section>
        )}
      </main>
    </Layout>
  );
}
