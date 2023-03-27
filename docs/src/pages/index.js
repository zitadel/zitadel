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
    title: "Documentation", // TODO: Plausible
    darkImageUrl: "img/index/Guides-dark.svg",
    lightImageUrl: "img/index/Guides-light.svg",
    link: "guides/overview",
    description: (
      <>
        Read our documentation and learn how you can setup, customize, and
        integrate authentication and authorization to your project.
      </>
    ),
    content: (
      <ListWrapper>
        <Column>
          <div>
            <ListElement
              link="/docs/guides/start/quickstart"
              type={ICONTYPE.START}
              title="Get started"
              description=""
            />
            <ListElement
              link="/docs/examples/sdks"
              type={ICONTYPE.APIS}
              title="SDKs"
              description=""
            />
            <ListElement
              link="/docs/examples/introduction"
              type={ICONTYPE.APIS}
              title="Example Apps"
              description=""
            />
            <ListElement
              link="/docs/guides/manage/console/overview"
              type={ICONTYPE.LOGIN}
              title="Manage"
              description="All about Console"
            />
            <ListElement
              link="/docs/guides/integrate"
              type={ICONTYPE.LOGIN}
              title="Integrate"
              description="Access our APIs and configure services and tools"
            />
          </div>
          <div>
            <ListElement
              link="/docs/guides/solution-scenarios/introduction"
              iconClasses="las la-paragraph"
              roundClasses="custom-rounded custom-rounded-split"
              label="B2C"
              title="Solution Scenarios"
              description=""
            />
            <ListElement
              link="/docs/concepts/introduction"
              type={ICONTYPE.TASKS}
              title="Concepts"
              description=""
            />
            <ListElement
              link="/docs/concepts/architecture/software"
              type={ICONTYPE.ARCHITECTURE}
              title="Architecture"
              description=""
            />
            <ListElement
              link="/docs/guides/manage/customize/branding"
              type={ICONTYPE.PRIVATELABELING}
              title="Customization"
              description=""
            />
            <ListElement
              link="/docs/support/troubleshooting"
              type={ICONTYPE.HELP}
              title="Support"
              description=""
            />
          </div>
        </Column>
      </ListWrapper>
    ),
  },
  {
    title: "Get Started",
    darkImageUrl: "/docs/img/index/Quickstarts-dark.svg",
    lightImageUrl: "img/index/Quickstarts-light.svg",
    link: "examples/introduction",
    description: (
      <>
        Learn how to integrate your applications and build secure workflows and
        APIs with ZITADEL.
      </>
    ),
    content: (
      <div className={styles.apilinks}>
        <ListWrapper>
          <ListElement
            link="/docs/guides/start/quickstart"
            type={ICONTYPE.START}
            title="Quick Start Guide"
            description="The ultimate guide to get started with ZITADEL."
          />
          <ListElement
            link="/docs/examples/login/angular"
            type={ICONTYPE.APIS}
            title="Frontend Quickstart Guides"
            description=""
          />
          <ListElement
            link="/docs/examples/secure-api/go"
            type={ICONTYPE.APIS}
            title="Backend Quickstart Guides"
            description=""
          />
          <ListElement
            link="/docs/examples/introduction"
            type={ICONTYPE.APIS}
            title="Examples"
            description="Clone an existing example application."
          />
        </ListWrapper>
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
            link="/docs/apis/auth/authentication-api-aka-auth"
            type={ICONTYPE.APIS}
            title="Authenticated User"
            description="All operations on the currently authenticated user."
          />
          <ListElement
            link="/docs/apis/mgmt/management-api"
            type={ICONTYPE.APIS}
            title="Organization Objects"
            description="Mutate IAM objects like organizations, projects, clients, users etc."
          />
          <ListElement
            link="/docs/apis/admin/administration-api-aka-admin"
            type={ICONTYPE.APIS}
            title="Instance Objects"
            description="Configure and manage the IAM instance."
          />
          <ListElement
            link="/docs/apis/openidoauth/endpoints"
            type={ICONTYPE.APIS}
            title="OIDC Endpoints"
            description=""
          />
          <ListElement
            link="/docs/apis/saml/endpoints"
            type={ICONTYPE.APIS}
            title="SAML Endpoints"
            description=""
          />
          <ListElement
            link="/docs/apis/actions/introduction"
            type={ICONTYPE.APIS}
            title="Actions"
            description="Customize and integrate ZITADEL into your landscape"
          />
        </ListWrapper>
      </div>
    ),
  },
  {
    title: "Self-hosting",
    darkImageUrl: "img/index/Concepts-dark.svg",
    lightImageUrl: "img/index/Concepts-light.svg",
    link: "/docs/self-hosting/deploy/overview",
    description: <>Everything you need to know about self-hosting ZITADEL.</>,
    content: (
      <ListWrapper>
        <ListElement
          link="/docs/self-hosting/deploy/overview"
          type={ICONTYPE.SYSTEM}
          title="Deploy"
          description=""
        />
        <ListElement
          link="/docs/self-hosting/manage/production"
          type={ICONTYPE.TASKS}
          title="Production Setup"
          description=""
        />
        <ListElement
          link="/docs/self-hosting/manage/configure"
          type={ICONTYPE.APIS}
          title="Configuration"
          description=""
        />
        <ListElement
          link="/docs/self-hosting/manage/updating_scaling"
          type={ICONTYPE.APIS}
          title="Update and Scaling"
          description=""
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
        <div className={styles.gigibanner}>Try ZITADEL Cloud for FREE 🚀</div>
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
