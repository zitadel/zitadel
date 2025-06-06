"use client";

import { useSearchParams } from "next/navigation";
import { useEffect } from "react";

export default function SamlPost() {
  const searchParams = useSearchParams();

  const url = searchParams.get("url");
  const relayState = searchParams.get("RelayState");
  const samlResponse = searchParams.get("SAMLResponse");

  useEffect(() => {
    // Automatically submit the form after rendering
    const form = document.getElementById("samlForm") as HTMLFormElement;
    if (form) {
      form.submit();
    }
  }, []);

  if (!url || !relayState || !samlResponse) {
    return (
      <p className="text-center">Missing required parameters for SAML POST.</p>
    );
  }

  return (
    <html lang="en">
      <head>
        <meta charSet="UTF-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <title>Redirecting...</title>
      </head>
      <body>
        <form id="samlForm" action={url} method="POST">
          <input type="hidden" name="RelayState" value={relayState} />
          <input type="hidden" name="SAMLResponse" value={samlResponse} />
        </form>
        <p>Redirecting...</p>
      </body>
    </html>
  );
}
