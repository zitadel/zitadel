import { NextRequest, NextResponse } from "next/server";
import { readdir, readFile, stat } from "fs/promises";
import { join } from "path";

interface SearchResult {
  serviceName: string;
  serviceDisplayName: string;
  path: string;
  method: string;
  operationId?: string;
  summary?: string;
  description?: string;
  tags?: string[];
}

async function getAllOpenApiFiles(
  dir: string,
  relativePath = ""
): Promise<Array<{ path: string; relativePath: string }>> {
  const files: Array<{ path: string; relativePath: string }> = [];

  try {
    const entries = await readdir(dir);

    for (const entry of entries) {
      const fullPath = join(dir, entry);
      const entryRelativePath = relativePath
        ? join(relativePath, entry)
        : entry;
      const stats = await stat(fullPath);

      if (stats.isDirectory()) {
        const subFiles = await getAllOpenApiFiles(fullPath, entryRelativePath);
        files.push(...subFiles);
      } else if (entry.endsWith("_service.swagger.json")) {
        files.push({ path: fullPath, relativePath: entryRelativePath });
      } else if (entry.endsWith(".swagger.json") && relativePath === "") {
        files.push({ path: fullPath, relativePath: entryRelativePath });
      }
    }
  } catch (error) {
    console.error(`Error reading directory ${dir}:`, error);
  }

  return files;
}

function getServiceDisplayName(serviceName: string): string {
  if (serviceName.includes("/")) {
    const parts = serviceName.split("/");
    const service = parts[0];
    const version = parts[1];
    const file = parts[2] || "";

    if (file.includes("_service")) {
      return `${
        service.charAt(0).toUpperCase() + service.slice(1)
      } API ${version.toUpperCase()}`;
    } else {
      return `${
        service.charAt(0).toUpperCase() + service.slice(1)
      } ${version.toUpperCase()}`;
    }
  }

  const nameMap: { [key: string]: string } = {
    management: "Management API (v1)",
    admin: "Admin API (v1)",
    auth: "Authentication API (v1)",
    system: "System API (v1)",
  };

  return (
    nameMap[serviceName] ||
    serviceName.charAt(0).toUpperCase() + serviceName.slice(1)
  );
}

export async function GET(request: NextRequest) {
  try {
    const { searchParams } = new URL(request.url);
    const query = searchParams.get("q");

    if (!query || query.trim().length < 2) {
      return NextResponse.json({ results: [] });
    }

    const searchTerm = query.toLowerCase().trim();
    const artifactsPath = join(
      process.cwd(),
      ".artifacts",
      "openapi",
      "zitadel"
    );
    const allFiles = await getAllOpenApiFiles(artifactsPath);

    const results: SearchResult[] = [];

    for (const file of allFiles) {
      try {
        const content = await readFile(file.path, "utf-8");
        const serviceName = file.relativePath.replace(/\.swagger\.json$/, "");
        const serviceDisplayName = getServiceDisplayName(serviceName);

        const parsed = JSON.parse(content) as any;

        if (!parsed.paths || Object.keys(parsed.paths).length === 0) {
          continue;
        }

        // Search through all paths and operations
        for (const [path, pathData] of Object.entries(parsed.paths)) {
          const pathStr = path as string;

          for (const [method, operation] of Object.entries(pathData as any)) {
            const methodStr = method.toLowerCase();

            if (methodStr === "parameters" || methodStr === "$ref") continue;

            const op = operation as any;
            const operationId = op.operationId || "";
            const summary = op.summary || "";
            const description = op.description || "";
            const tags = op.tags || [];

            // Search in various fields
            const searchableText = [
              pathStr,
              operationId,
              summary,
              description,
              ...tags,
              serviceName,
              serviceDisplayName,
            ]
              .join(" ")
              .toLowerCase();

            if (searchableText.includes(searchTerm)) {
              results.push({
                serviceName,
                serviceDisplayName,
                path: pathStr,
                method: methodStr.toUpperCase(),
                operationId,
                summary,
                description,
                tags,
              });
            }
          }
        }
      } catch (error) {
        console.error(`Error processing file ${file.path}:`, error);
      }
    }

    // Sort results by relevance (exact matches first, then partial matches)
    results.sort((a, b) => {
      const aExact =
        a.operationId?.toLowerCase().includes(searchTerm) ||
        a.summary?.toLowerCase().includes(searchTerm) ||
        a.path.toLowerCase().includes(searchTerm);
      const bExact =
        b.operationId?.toLowerCase().includes(searchTerm) ||
        b.summary?.toLowerCase().includes(searchTerm) ||
        b.path.toLowerCase().includes(searchTerm);

      if (aExact && !bExact) return -1;
      if (!aExact && bExact) return 1;

      return a.serviceName.localeCompare(b.serviceName);
    });

    return NextResponse.json({
      results: results.slice(0, 50), // Limit to 50 results
      total: results.length,
    });
  } catch (error) {
    console.error("Error searching APIs:", error);
    return NextResponse.json(
      { error: "Failed to search APIs" },
      { status: 500 }
    );
  }
}
