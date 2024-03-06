import { headers } from "next/headers";

export default function Page() {
  const headersList = headers();
  const hds = [
    "x-zitadel-login-client",
    "forwarded",
    "x-zitadel-forwarded",
    "host",
    "referer",
  ];
  return (
    <div className="space-y-8">
      <h1 className="text-xl font-medium">Headers</h1>
      {hds.map((h) => (
        <p>
          {h}:{headersList.get(h)}
        </p>
      ))}
    </div>
  );
}
