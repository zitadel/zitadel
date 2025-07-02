import { NextRequest, NextResponse } from "next/server";

export async function GET(request: NextRequest) {
  const searchParams = request.nextUrl.searchParams;
  const url = searchParams.get("url");
  const relayState = searchParams.get("RelayState");
  const samlResponse = searchParams.get("SAMLResponse");

  if (!url || !relayState || !samlResponse) {
    return new NextResponse("Missing required parameters", { status: 400 });
  }

  // Respond with an HTML form that auto-submits via POST
  const html = `
    <html>
      <body onload="document.forms[0].submit()">
        <form action="${url}" method="post">
          <input type="hidden" name="RelayState" value="${relayState}" />
          <input type="hidden" name="SAMLResponse" value="${samlResponse}" />
          <noscript>
            <button type="submit">Continue</button>
          </noscript>
        </form>
      </body>
    </html>
  `;
  return new NextResponse(html, {
    headers: { "Content-Type": "text/html" },
  });
}
