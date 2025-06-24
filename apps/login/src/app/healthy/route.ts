import { NextResponse } from "next/server";

export async function GET() {
  return NextResponse.json({token: process.env.ZITADEL_SERVICE_USER_TOKEN}, { status: 200 });
}
