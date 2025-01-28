import { NextRequest, NextResponse } from "next/server";

export const dynamic = "force-dynamic";
export const revalidate = false;
export const fetchCache = "default-no-store";

export async function GET(request: NextRequest) {
  const headers = request.headers;

  const headersObject: Record<string, string> = {};
  headers.forEach((value, key) => {
    headersObject[key] = value;
  });

  return NextResponse.json(headersObject);
}
