import React from "react";
import { useRouter } from "next/router";
import Link from "next/link";

const items = [
  {
    title: "Get started",
    links: [
      { href: "/guides", children: "Guides" },
      { href: "/examples", children: "Examples" },
      { href: "/apis", children: "APIs" },
      { href: "/concepts", children: "Concepts" },
      { href: "/help", children: "Help" },
      { href: "/legal", children: "Legal" },
    ],
  },
];

export function SideNav() {
  const router = useRouter();

  return (
    <nav className="sticky top-0 h-screen bottom-0 bg-zinc-500 py-4 border-r border-border-light dark:border-border-dark">
      {items.map((item) => (
        <div key={item.title}>
          <span>{item.title}</span>
          <ul className="flex column">
            {item.links.map((link) => {
              const active = router.pathname === link.href;
              return (
                <li
                  key={link.href}
                  className={active ? "active" : "text-red-500"}
                >
                  <Link {...link} />
                </li>
              );
            })}
          </ul>
        </div>
      ))}
      {/* <style jsx>
        {`
          nav {
            position: sticky;
            height: calc(100vh - var(--top-nav-height));
            flex: 0 0 auto;
            overflow-y: auto;
            padding: 2.5rem 2rem 2rem;
            border-right: 1px solid var(--border-color);
          }
          span {
            font-size: larger;
            font-weight: 500;
            padding: 0.5rem 0 0.5rem;
          }
          ul {
            padding: 0;
          }
          li {
            list-style: none;
            margin: 0;
          }
          li :global(a) {
            text-decoration: none;
          }
          li :global(a:hover),
          li.active :global(a) {
            text-decoration: underline;
          }
        `}
      </style> */}
    </nav>
  );
}
