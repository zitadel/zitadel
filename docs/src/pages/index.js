import React from 'react';
import clsx from 'clsx';
import Layout from '@theme/Layout';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import useBaseUrl from '@docusaurus/useBaseUrl';
import styles from './styles.module.css';

const features = [
  {
    title: 'Manuals',
    imageUrl: 'img/undraw_docusaurus_mountain.svg',
    link: 'docs/manuals',
    description: (
      <>
        Follow this guide to get started with ZITADEL as a user.
      </>
    ),
  },
  {
    title: 'Quickstarts',
    imageUrl: 'img/undraw_docusaurus_tree.svg',
    link: 'docs/quickstarts/quickstarts',
    description: (
      <>
        Learn how to integrate your applications and build secure workflows and APIs with ZITADEL
      </>
    ),
  },
  {
    title: 'Guides',
    imageUrl: 'img/undraw_docusaurus_react.svg',
    link: 'docs',
    description: (
      <>
        Read our guides on how to manage your data and role associations in ZITADEL and on what we recommend.
      </>
    ),
  },
  {
    title: 'APIs',
    imageUrl: 'img/undraw_docusaurus_react.svg',
    link: '/docs/apis/apis',
    description: (
      <>
        Learn more about our APIs and how to integrate them in your apps.
      </>
    ),
  },
  {
    title: 'Architecture',
    imageUrl: 'img/undraw_docusaurus_react.svg',
    link: 'docs/architecture',
    description: (
      <>
        Learn more about engineering and design principles, ZITADELs architecture and used technologies.
      </>
    ),
  },
];

function Feature({imageUrl, title, description, link}) {
  const imgUrl = useBaseUrl(imageUrl);
  return (
        <div className={clsx('col col--4', styles.feature)}>
          <Link to={useBaseUrl(link)}>
          {imgUrl && (
              <div className="text--center">
                <img className={styles.featureImage} src={imgUrl} alt={title} />
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
                'button button--outline button--secondary button--lg',
                styles.getStarted,
              )}
              to={useBaseUrl('docs/quickstarts/quickstarts')}>
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
