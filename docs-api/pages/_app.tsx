import React from "react";
import Head from "next/head";
import Link from "next/link";

import "prismjs";
// Import other Prism themes here
import "prismjs/components/prism-bash.min";
import "prismjs/themes/prism.css";

import "../styles/globals.css";

import type { AppProps } from "next/app";
import type { MarkdocNextJsPageProps } from "@markdoc/next.js";
import { TopNav } from "../components/TopNav";
import { TableOfContents } from "../components/TableOfContents";
import { ThemeProvider } from "next-themes";

const TITLE = "Markdoc";
const DESCRIPTION = "A powerful, flexible, Markdown-based authoring framework";

function collectHeadings(node, sections = []) {
  if (node) {
    if (node.name === "Heading") {
      const title = node.children[0];

      if (typeof title === "string") {
        sections.push({
          ...node.attributes,
          title,
        });
      }
    }

    if (node.children) {
      for (const child of node.children) {
        collectHeadings(child, sections);
      }
    }
  }

  return sections;
}

export type MyAppProps = MarkdocNextJsPageProps;

export default function MyApp({ Component, pageProps }: AppProps<MyAppProps>) {
  const { markdoc } = pageProps;

  let title = TITLE;
  let description = DESCRIPTION;
  if (markdoc) {
    if (markdoc.frontmatter.title) {
      title = markdoc.frontmatter.title;
    }
    if (markdoc.frontmatter.description) {
      description = markdoc.frontmatter.description;
    }
  }

  const toc = pageProps.markdoc?.content
    ? collectHeadings(pageProps.markdoc.content)
    : [];

  console.log(toc);

  return (
    <>
      <Head>
        <title>{title}</title>
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <meta name="referrer" content="strict-origin" />
        <meta name="title" content={title} />
        <meta name="description" content={description} />
        <link rel="shortcut icon" href="/favicon.ico" />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <ThemeProvider
        attribute="class"
        defaultTheme="system"
        storageKey="cp-theme"
      >
        <div className="h-screen flex w-screen flex-grow bg-white dark:bg-background-dark-500 text-gray-800 dark:text-gray-100">
          <TableOfContents toc={toc} />
          <div className="overflow-auto flex-grow pt-0 pl-4 pr-8 pb-8">
            <TopNav></TopNav>
            <main className="">
              <Component {...pageProps} />
            </main>
          </div>
        </div>
      </ThemeProvider>
    </>
  );
}
