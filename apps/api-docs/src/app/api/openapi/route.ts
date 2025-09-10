import { NextRequest, NextResponse } from "next/server";
import { readdir, readFile, stat } from "fs/promises";
import { join } from "path";

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
        // Recursively search subdirectories
        const subFiles = await getAllOpenApiFiles(fullPath, entryRelativePath);
        files.push(...subFiles);
      } else if (entry.endsWith("_service.openapi.yaml")) {
        // Only include service files that contain actual API endpoints
        files.push({ path: fullPath, relativePath: entryRelativePath });
      } else if (entry.endsWith(".openapi.yaml") && relativePath === "") {
        // Include top-level v1 API files (management, admin, auth, system)
        files.push({ path: fullPath, relativePath: entryRelativePath });
      }
    }
  } catch (error) {
    console.error(`Error reading directory ${dir}:`, error);
  }

  return files;
}

export async function GET(request: NextRequest) {
  try {
    const artifactsPath = join(
      process.cwd(),
      ".artifacts",
      "openapi3",
      "zitadel"
    );

    // Get all OpenAPI spec files recursively
    const allFiles = await getAllOpenApiFiles(artifactsPath);

    const specs = await Promise.all(
      allFiles.map(async (file: { path: string; relativePath: string }) => {
        try {
          const content = await readFile(file.path, "utf-8");
          const serviceName = file.relativePath.replace(/\.openapi\.yaml$/, "");

          return {
            name: serviceName,
            fileName: file.relativePath,
            content: content,
          };
        } catch (error) {
          console.error(`Error reading file ${file.path}:`, error);
          return null;
        }
      })
    );

    // Filter out null entries and return valid specs
    const validSpecs = specs.filter((spec) => spec !== null);

    return NextResponse.json({ specs: validSpecs });
  } catch (error) {
    console.error("Error reading OpenAPI specs:", error);
    return NextResponse.json(
      { error: "Failed to load OpenAPI specifications" },
      { status: 500 }
    );
  }
}
