import Link from "@docusaurus/Link";
import useBaseUrl from "@docusaurus/useBaseUrl";
import useDocusaurusContext from "@docusaurus/useDocusaurusContext";
import Layout from "@theme/Layout";
import ThemedImage from "@theme/ThemedImage";
import clsx from "clsx";
import React from "react";

import Column from "../components/column";
import {
  HomeListWrapper,
  ICONTYPE,
  ListElement,
  ListWrapper,
} from "../components/list";
import styles from "./styles.module.css";

const features = [
  {
    title: "Guides",
    darkImageUrl: "img/index/Guides-dark.svg",
    lightImageUrl: "img/index/Guides-light.svg",
    link: "guides/overview",
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
              link="guides/start/quickstart"
              type={ICONTYPE.START}
              title="Get started"
              description=""
            />
            <ListElement
              link="guides/manage/cloud/overview"
              type={ICONTYPE.LOGIN}
              title="ZITADEL Cloud"
              description=""
            />
            <ListElement
              link="guides/integrate/login-users"
              type={ICONTYPE.LOGIN}
              title="Login Users"
              description=""
            />
            <ListElement
              link="guides/integrate/access-zitadel-apis"
              type={ICONTYPE.APIS}
              title="Access APIs"
              description=""
            />
          </div>
          <div>
            <ListElement
              link="guides/solution-scenarios/introduction"
              iconClasses="las la-paragraph"
              roundClasses="custom-rounded custom-rounded-split"
              label="B2C"
              title="Solution Scenarios"
              description=""
            />
            <ListElement
              link="guides/manage/customize/branding"
              type={ICONTYPE.PRIVATELABELING}
              title="Customization"
              description=""
            />
            <ListElement
              link="guides/deploy/overview"
              type={ICONTYPE.SYSTEM}
              title="Deploy"
              description=""
            />
            <ListElement
              link="guides/trainings/introduction"
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
    darkImageUrl: "/docs/img/index/Quickstarts-dark.svg",
    lightImageUrl: "img/index/Quickstarts-light.svg",
    link: "examples/introduction",
    description: (
      <>
        Learn how to integrate your applications and build secure workflows and
        APIs with ZITADEL
      </>
    ),
    content: (
      <div className={styles.quickstartcontainer}>
        <QuickstartLink
          link="/examples/login/angular"
          imageSource="/docs/img/tech/angular.svg"
          title="Angular"
          description="Add the user login to your application and query some data from the userinfo endpoint"
        />
        <QuickstartLink
          link="/examples/login/react"
          imageSource="/docs/img/tech/react.png"
          title="React"
          description="Logs into your application and queries some data from the userinfo endpoint"
        />
        <QuickstartLink
          link="/examples/login/flutter"
          imageSource="/docs/img/tech/flutter.svg"
          title="Flutter"
          description="Mobile Application working for iOS and Android that authenticates your user."
        />
        <QuickstartLink
          link="/examples/login/nextjs"
          imageSource="/docs/img/tech/nextjslight.svg"
          lightImageSource="/docs/img/tech/nextjs.svg"
          title="NextJS"
          description="A simple application to log into your user account and query some data from User endpoint."
        />
      </div>
    ),
  },
  {
    title: "APIs",
    darkImageUrl: "/docs/img/index/APIs-dark.svg",
    lightImageUrl: "/docs/img/index/APIs-light.svg",
    link: "/apis/introduction",
    description: (
      <>Learn more about our APIs and how to integrate them in your apps.</>
    ),
    content: (
      <div className={styles.apilinks}>
        <ListWrapper>
          <ListElement
            link="./apis/proto/auth"
            type={ICONTYPE.APIS}
            title="Proto Definitions"
            description=""
          />
          <ListElement
            link="./apis/openidoauth/endpoints"
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
    link: "concepts/introduction",
    description: (
      <>
        Learn more about engineering and design principles, ZITADELs
        architecture and used technologies.
      </>
    ),
    content: (
      <ListWrapper>
        <ListElement
          link="./concepts/principles"
          type={ICONTYPE.TASKS}
          title="Principles"
          description="Design and engineering principles"
        />
        <ListElement
          link="./concepts/architecture/software"
          type={ICONTYPE.ARCHITECTURE}
          title="Architecture"
          description="Sotware-, Cluster- and Multi Cluster Architecture"
        />
        <ListElement
          link="./concepts/structure/overview"
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

const Gigi = () => {
  return (
    <div className={styles.gigiwrapper}>
      <div className={styles.gigiwrapperrelative}>
        <img height="151px" width="256px" src="/docs/img/gigi.svg" />
        <div className={styles.gigibanner}>ZITADEL Cloud OUT NOW! 🚀</div>
      </div>
    </div>
  );
};

export default function Home() {
  const context = useDocusaurusContext();
  const { siteConfig = {} } = context;
  return (
    <Layout description={`${siteConfig.customFields.description}`}>
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
              to={useBaseUrl("guides/start/quickstart")}
            >
              Get Started
            </Link>
          </div>
        </div>
        <Link to="https://zitadel.com">
          <Gigi />
        </Link>
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
