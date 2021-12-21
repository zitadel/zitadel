import Link from '@docusaurus/Link';
import useBaseUrl from '@docusaurus/useBaseUrl';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import ThemedImage from '@theme/ThemedImage';
import clsx from 'clsx';
import React from 'react';

import Column from '../components/column';
import { HomeListWrapper, ICONTYPE, ListElement, ListWrapper } from '../components/list';
import styles from './styles.module.css';

const features = [
  {
    title: "Guides",
    darkImageUrl: "img/index/Guides-dark.svg",
    lightImageUrl: "img/index/Guides-light.svg",
    link: "docs/guides/overview",
    description: (
      <>
        Read our guides on how to manage your data and role associations in
        ZITADEL and on what we recommend.
      </>
    ),
    content: (
      <ListWrapper>
        <Column>
          <div>
            <ListElement
              link="docs/guides/basics/get-started"
              type={ICONTYPE.START}
              title="Get started"
              description=""
            />
            <ListElement
              link="docs/guides/authentication/login-users"
              type={ICONTYPE.LOGIN}
              title="Authentication"
              description=""
            />
            <ListElement
              link="docs/guides/authorization/oauth-recommended-flows"
              type={ICONTYPE.LOGIN}
              title="Authorzation"
              description=""
            />
            <ListElement
              link="docs/guides/api/access-zitadel-apis"
              type={ICONTYPE.APIS}
              title="Access APIs"
              description=""
            />
          </div>
          <div>
            {/* <ListElement link="docs/guides/architecture-scenarios/introduction" iconClasses="las la-paragraph" roundClasses="rounded rounded-split" label="B2C" title="Architecture Scenarios" description="" /> */}
            <ListElement
              link="docs/guides/customization/branding"
              type={ICONTYPE.PRIVATELABELING}
              title="Customization"
              description=""
            />
            <ListElement
              link="docs/guides/installation/shared-cloud"
              type={ICONTYPE.STORAGE}
              title="Installation"
              description=""
            />
            <ListElement
              link="docs/guides/trainings/introduction"
              type={ICONTYPE.STORAGE}
              title="Trainings"
              description=""
            />
          </div>
        </Column>
      </ListWrapper>
    ),
  },
  {
    title: "Quickstarts",
    darkImageUrl: "img/index/Quickstarts-dark.svg",
    lightImageUrl: "img/index/Quickstarts-light.svg",
    link: "docs/quickstarts/introduction",
    description: (
      <>
        Learn how to integrate your applications and build secure workflows and
        APIs with ZITADEL
      </>
    ),
    content: (
      <div className={styles.quickstartcontainer}>
        <QuickstartLink
          link="/docs/quickstarts/login/angular"
          imageSource="/img/tech/angular.svg"
          title="Angular"
          description="Add the user login to your application and query some data from the userinfo endpoint"
        />
        <QuickstartLink
          link="/docs/quickstarts/login/react"
          imageSource="/img/tech/react.png"
          title="React"
          description="Logs into your application and queries some data from the userinfo endpoint"
        />
        <QuickstartLink
          link="/docs/quickstarts/login/flutter"
          imageSource="/img/tech/flutter.svg"
          title="Flutter"
          description="Mobile Application working for iOS and Android that authenticates your user."
        />
        <QuickstartLink
          link="/docs/quickstarts/login/nextjs"
          imageSource="/img/tech/nextjslight.svg"
          lightImageSource="/img/tech/nextjs.svg"
          title="NextJS"
          description="A simple application to log into your user account and query some data from User endpoint."
        />
      </div>
    ),
  },
  {
    title: "APIs",
    darkImageUrl: "img/index/APIs-dark.svg",
    lightImageUrl: "img/index/APIs-light.svg",
    link: "/docs/apis/introduction",
    description: (
      <>Learn more about our APIs and how to integrate them in your apps.</>
    ),
    content: (
      <div className={styles.apilinks}>
        <ListWrapper>
          <ListElement
            link="./docs/apis/proto/auth"
            type={ICONTYPE.APIS}
            title="Proto Definitions"
            description=""
          />
          <ListElement
            link="./docs/apis/openidoauth/endpoints"
            type={ICONTYPE.APIS}
            title="OpenID Connect and OAuth"
            description="Scopes, Claims, Authentication Methods, Grant Types"
          />
        </ListWrapper>
      </div>
    ),
  },
  {
    title: "Concepts",
    darkImageUrl: "img/index/Concepts-dark.svg",
    lightImageUrl: "img/index/Concepts-light.svg",
    link: "docs/concepts/introduction",
    description: (
      <>
        Learn more about engineering and design principles, ZITADELs
        architecture and used technologies.
      </>
    ),
    content: (
      <ListWrapper>
        <ListElement
          link="./docs/concepts/principles"
          type={ICONTYPE.TASKS}
          title="Principles"
          description="Design and engineering principles"
        />
        <ListElement
          link="./docs/concepts/eventstore"
          type={ICONTYPE.STORAGE}
          title="Eventstore"
          description="Learn how ZITADEL stores data"
        />
        <ListElement
          link="./docs/concepts/architecture"
          type={ICONTYPE.ARCHITECTURE}
          title="Architecture"
          description="Sotware-, Cluster- and Multi Cluster Architecture"
        />
        <ListElement
          link="./docs/concepts/structure/overview"
          type={ICONTYPE.ARCHITECTURE}
          title="Structure"
          description="Object structure of ZITADEL"
        />
      </ListWrapper>
    ),
  },
];

function QuickstartLink({ link, title, imageSource, lightImageSource }) {
  return (
    <Link href={link} className={clsx("", styles.quickstart)}>
      {/* <img className={styles.quickstartlinkimg} src={imageSource} alt={`${title}`}/> */}
      <ThemedImage
        className={styles.quickstartlinkimg}
        alt={title}
        sources={{
          light: lightImageSource ? lightImageSource : imageSource,
          dark: imageSource,
        }}
      />
      <p>{title}</p>
    </Link>
  );
}

function Feature({
  darkImageUrl,
  lightImageUrl,
  title,
  description,
  link,
  content,
}) {
  const darkImgUrl = useBaseUrl(darkImageUrl);
  const lightImgUrl = useBaseUrl(lightImageUrl);

  const themedImage = (
    <ThemedImage
      className={styles.featureImage}
      alt={title}
      sources={{
        light: lightImgUrl,
        dark: darkImgUrl,
      }}
    />
  );
  return (
    <div className={clsx("col col--6 docs-link", styles.feature)}>
      {darkImgUrl && lightImgUrl && (
        <div className="">
          <HomeListWrapper image={themedImage}>
            <Link to={useBaseUrl(link)}>
              <h3 className={styles.homelink}>
                {title}
                <i
                  className={clsx("las la-angle-right", styles.homelinkicon)}
                ></i>
              </h3>
            </Link>
            <p className="">{description}</p>

            {content}
          </HomeListWrapper>
        </div>
      )}
    </div>
  );
}

export default function Home() {
  const context = useDocusaurusContext();
  const { siteConfig = {} } = context;
  return (
    <Layout
      title={`${siteConfig.title}`}
      description="This site bundles ZITADELs Documentations"
    >
      <header className={clsx("hero", styles.heroBanner)}>
        <div className="container">
          <h1 className="hero__title">{siteConfig.title}</h1>
          <p className="hero__subtitle">{siteConfig.tagline}</p>
          <div className={styles.buttons}>
            <Link
              className={clsx(
                "button button--outline button--lg get-started",
                styles.getStarted
              )}
              to={useBaseUrl("docs/guides/basics/get-started")}
            >
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
