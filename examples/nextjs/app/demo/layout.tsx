import type { ReactNode } from "react";

const demoLinks = [
  { href: "/demo/oidc", label: "OIDC" },
  { href: "/demo/username-password", label: "Username/password" },
  { href: "/demo/signup", label: "Signup" },
  { href: "/demo/org-registration", label: "Org registration" },
];

export default function DemoLayout({ children }: Readonly<{ children: ReactNode }>) {
  return (
    <main className="ztdl-page ztdl-noise">
      <p className="ztdl-eyebrow">ZITADEL SDK demo</p>
      <h1 className="ztdl-title">Demo lanes</h1>
      <p className="ztdl-subtitle">
        <a href="/">Back to landing page</a>
      </p>
      <nav aria-label="Demo lanes" className="ztdl-nav">
        <ul className="ztdl-nav-list">
          {demoLinks.map((link) => (
            <li key={link.href}>
              <a className="ztdl-nav-link" href={link.href}>
                {link.label}
              </a>
            </li>
          ))}
        </ul>
      </nav>
      {children}
    </main>
  );
}
