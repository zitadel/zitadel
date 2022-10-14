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
    link: "docs/guides/overview",
    description: (
      <>
        Read our guides on how to manage your data and role associations in
        ZITADEL and on what we recommend.
      </>
    ),
  },
  {
    title: "Integrate",
    darkImageUrl: "img/index/Quickstarts-dark.svg",
    lightImageUrl: "img/index/Quickstarts-light.svg",
    link: "docs/examples/introduction",
    description: (
      <>
        Learn how to integrate your applications and build secure workflows and
        APIs with ZITADEL
      </>
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
          <img height="151px" width="256px" src="/img/gigi.svg" />
          <div className={styles.gigibanner}>ZITADEL Cloud OUT NOW! 🚀</div>
        </div>
      </div>
    
  );
};

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
              to={useBaseUrl("docs/guides/start/quickstart")}
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
