import { NextRequest, NextResponse } from "next/server";

export async function GET(request: NextRequest) {
  return NextResponse.json({ found: true });
}

export async function POST(request: NextRequest) {
  const body = await request.json();
  if (body) {
    return NextResponse.json({ found: true });
  } else {
    return NextResponse.error();
  }
}
