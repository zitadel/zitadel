import { NextRequest, NextResponse } from "next/server";
import { readdir, readFile } from "fs/promises";
import { join } from "path";
import { existsSync } from "fs";

interface VersionMetadata {
  version: string;
  generatedAt?: string;
  gitCommit?: string;
  gitBranch?: string;
}

export async function GET(request: NextRequest) {
  try {
    const versions: Array<{
      id: string;
      name: string;
      isDefault: boolean;
      available: boolean;
      metadata?: VersionMetadata;
    }> = [];

    // Add current/latest version
    const currentArtifactsPath = join(
      process.cwd(),
      ".artifacts",
      "openapi3",
      "zitadel"
    );
    versions.push({
      id: "latest",
      name: "Latest (Current Branch)",
      isDefault: true,
      available: existsSync(currentArtifactsPath),
    });

    // Check for organized version folders (.artifacts/versions/)
    const versionsDir = join(process.cwd(), ".artifacts", "versions");

    if (existsSync(versionsDir)) {
      try {
        const versionDirs = await readdir(versionsDir);

        for (const versionDir of versionDirs) {
          const versionPath = join(
            versionsDir,
            versionDir,
            "openapi3",
            "zitadel"
          );
          const metadataPath = join(versionsDir, versionDir, "metadata.json");

          // Read metadata if available
          let metadata: VersionMetadata | undefined;
          try {
            if (existsSync(metadataPath)) {
              const metadataContent = await readFile(metadataPath, "utf-8");
              metadata = JSON.parse(metadataContent);
            }
          } catch (error) {
            console.warn(
              `Failed to read metadata for version ${versionDir}:`,
              error
            );
          }

          const displayName =
            metadata?.gitBranch === "main"
              ? `${versionDir} (Main Branch)`
              : versionDir;

          versions.push({
            id: versionDir,
            name: displayName,
            isDefault: false,
            available: existsSync(versionPath),
            metadata,
          });
        }
      } catch (error) {
        console.error("Error reading version directories:", error);
      }
    }

    // Also check legacy versioned artifacts for backward compatibility
    const legacyVersionedArtifactsDir = join(
      process.cwd(),
      ".artifacts-versioned"
    );

    if (existsSync(legacyVersionedArtifactsDir)) {
      try {
        const versionDirs = await readdir(legacyVersionedArtifactsDir);

        for (const versionDir of versionDirs) {
          const versionPath = join(
            legacyVersionedArtifactsDir,
            versionDir,
            "openapi3",
            "zitadel"
          );

          // Skip if already added from organized versions
          if (!versions.find((v) => v.id === versionDir)) {
            versions.push({
              id: versionDir,
              name: versionDir === "main" ? "Main Branch" : versionDir,
              isDefault: false,
              available: existsSync(versionPath),
            });
          }
        }
      } catch (error) {
        console.error("Error reading legacy versioned artifacts:", error);
      }
    }

    // Sort versions (latest first, then by semantic version)
    const sortedVersions = versions.sort((a, b) => {
      if (a.isDefault) return -1;
      if (b.isDefault) return 1;

      // Simple version sorting - you might want to use a proper semver library
      if (a.id === "main") return -1;
      if (b.id === "main") return 1;

      return b.id.localeCompare(a.id, undefined, { numeric: true });
    });

    return NextResponse.json({
      versions: sortedVersions.filter((v) => v.available),
      default: sortedVersions.find((v) => v.isDefault)?.id || "latest",
    });
  } catch (error) {
    console.error("Error fetching versions:", error);
    return NextResponse.json(
      { error: "Failed to fetch available versions" },
      { status: 500 }
    );
  }
}
