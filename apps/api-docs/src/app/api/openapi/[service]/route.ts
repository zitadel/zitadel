import { NextRequest, NextResponse } from "next/server";
import { readFile } from "fs/promises";
import { join } from "path";

export async function GET(
  request: NextRequest,
  { params }: { params: { service: string } }
) {
  try {
    const { service } = params;
    const artifactsPath = join(
      process.cwd(),
      ".artifacts",
      "openapi3",
      "zitadel"
    );

    // Convert service name back to file path
    // e.g., "user/v2/user_service" -> "user/v2/user_service.openapi.yaml"
    const filePath = join(artifactsPath, `${service}.openapi.yaml`);

    try {
      const content = await readFile(filePath, "utf-8");
      return new NextResponse(content, {
        headers: {
          "Content-Type": "application/yaml",
          "Access-Control-Allow-Origin": "*",
        },
      });
    } catch (error) {
      console.error(
        `Error reading OpenAPI spec for service ${service}:`,
        error
      );
      return NextResponse.json(
        { error: `OpenAPI specification not found for service: ${service}` },
        { status: 404 }
      );
    }
  } catch (error) {
    console.error(
      `Error processing request for service ${params.service}:`,
      error
    );
    return NextResponse.json(
      {
        error: `OpenAPI specification not found for service: ${params.service}`,
      },
      { status: 404 }
    );
  }
}
