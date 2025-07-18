import { getSAMLFormCookie } from "@/lib/saml";
import { NextRequest, NextResponse } from "next/server";

export async function GET(request: NextRequest) {
  const searchParams = request.nextUrl.searchParams;
  const url = searchParams.get("url");
  const id = searchParams.get("id");

  if (!url) {
    return new NextResponse("Missing url parameter", { status: 400 });
  }

  if (!id) {
    return new NextResponse("Missing id parameter", { status: 400 });
  }

  const formData = await getSAMLFormCookie(id);

  const formDataParsed = formData ? JSON.parse(formData) : null;

  if (!formDataParsed) {
    return new NextResponse("SAML form data not found", { status: 404 });
  }

  // Generate hidden input fields for all key-value pairs in formDataParsed
  const hiddenInputs = Object.entries(formDataParsed)
    .map(
      ([key, value]) =>
        `<input type="hidden" name="${key}" value="${value}" />`,
    )
    .join("\n    ");

  // Respond with an HTML form that auto-submits via POST
  const html = `
    <html>
      <body onload="document.forms[0].submit()">
        <form action="${url}" method="post">
          ${hiddenInputs}
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
