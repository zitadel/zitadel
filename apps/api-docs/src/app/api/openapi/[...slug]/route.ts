import { NextRequest, NextResponse } from "next/server";
import { readFile } from "fs/promises";
import { join } from "path";
import { existsSync } from "fs";

export async function GET(
  request: NextRequest,
  { params }: { params: Promise<{ slug: string[] }> }
) {
  let version: string;
  let service: string;

  try {
    const resolvedParams = await params;
    const slugs = resolvedParams.slug;

    // Parse URL structure:
    // /api/openapi/service -> ["service"] (latest version)
    // /api/openapi/service/path -> ["service", "path"] (latest version)
    // /api/openapi/service/path/more -> ["service", "path", "more"] (latest version)
    // /api/openapi/version/service -> ["version", "service"] (specific version if version matches pattern)
    // /api/openapi/version/service/path -> ["version", "service", "path"] (specific version)
    // /api/openapi/version/service/path/more -> ["version", "service", "path", "more"] (specific version)

    console.log("URL slugs:", slugs);

    // Check if first slug looks like a version (starts with 'v' followed by numbers/dots)
    const versionPattern = /^v\d+(\.\d+)*(\.\d+)*(-.*)?$/;
    const firstSlugIsVersion = versionPattern.test(slugs[0]);

    if (firstSlugIsVersion && slugs.length >= 2) {
      // First slug is a version, rest is service path
      version = slugs[0];
      service = slugs.slice(1).join("/");
    } else {
      // All slugs form the service path, use latest version
      version = "latest";
      service = slugs.join("/");
    }

    console.log("Parsed version:", version, "service:", service);
  } catch (error) {
    console.error("Error resolving params:", error);
    return NextResponse.json(
      { error: "Invalid request parameters" },
      { status: 400 }
    );
  }

  try {
    // Determine artifacts path based on version
    let artifactsPath: string;

    if (version === "latest") {
      // Use current artifacts
      artifactsPath = join(process.cwd(), ".artifacts", "openapi3", "zitadel");
    } else {
      // Try organized version folders first
      const organizedPath = join(
        process.cwd(),
        ".artifacts",
        "versions",
        version,
        "openapi3",
        "zitadel"
      );

      if (existsSync(organizedPath)) {
        artifactsPath = organizedPath;
      } else {
        // Fallback to legacy versioned artifacts
        const legacyPath = join(
          process.cwd(),
          ".artifacts-versioned",
          version,
          "openapi3",
          "zitadel"
        );
        if (existsSync(legacyPath)) {
          artifactsPath = legacyPath;
        } else {
          return NextResponse.json(
            { error: `Version ${version} not found` },
            { status: 404 }
          );
        }
      }
    }

    // Convert service name back to file path
    const filePath = join(artifactsPath, `${service}.openapi.yaml`);

    try {
      const content = await readFile(filePath, "utf-8");

      // Add version metadata to the response
      return new NextResponse(content, {
        headers: {
          "Content-Type": "application/yaml",
          "Access-Control-Allow-Origin": "*",
          "X-API-Version": version,
        },
      });
    } catch (error) {
      console.error(
        `Error reading OpenAPI spec for service ${service} version ${version}:`,
        error
      );
      return NextResponse.json(
        {
          error: `OpenAPI specification not found for service: ${service} in version: ${version}`,
        },
        { status: 404 }
      );
    }
  } catch (error) {
    console.error(
      `Error processing request for service ${service} version ${version}:`,
      error
    );
    return NextResponse.json(
      {
        error: `OpenAPI specification not found for service: ${service} in version: ${version}`,
      },
      { status: 404 }
    );
  }
}
