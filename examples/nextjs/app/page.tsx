const lanes = [
  {
    href: "/demo/oidc",
    label: "OIDC",
    description: "Validate redirect sign-in, callback processing, and sign-out behavior.",
  },
  {
    href: "/demo/username-password",
    label: "Username/password",
    description: "Create and inspect Session API-backed sessions from credentials.",
  },
  {
    href: "/demo/signup",
    label: "Signup",
    description: "Create users with password and inspect request/response details.",
  },
  {
    href: "/demo/org-registration",
    label: "Org registration",
    description: "Create organizations and inspect API permissions/error behavior.",
  },
];

export default function HomePage() {
  return (
    <main className="ztdl-page ztdl-noise">
      <p className="ztdl-eyebrow">ZITADEL SDK demo</p>
      <h1 className="ztdl-title">Next.js reference lanes</h1>
      <p className="ztdl-subtitle">
        Explore the core integration paths with website-aligned styling while developing against workspace-linked SDK packages.
      </p>
      <ul className="ztdl-card-grid">
        {lanes.map((lane) => (
          <li key={lane.href} className="ztdl-card ztdl-noise">
            <h3>
              <a href={lane.href}>{lane.label}</a>
            </h3>
            <p>{lane.description}</p>
          </li>
        ))}
      </ul>
    </main>
  );
}
