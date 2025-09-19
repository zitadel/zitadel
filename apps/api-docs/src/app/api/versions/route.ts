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

interface VersionConfig {
  versions: Array<{
    id: string;
    name: string;
    gitRef: string;
    enabled: boolean;
    isStable: boolean;
  }>;
  settings: {
    defaultVersion: string;
    autoGenerate: boolean;
    maxVersions: number;
    includePrerelease: boolean;
  };
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

    // Load version config for fallback
    let versionConfig: VersionConfig | null = null;
    try {
      const configPath = join(process.cwd(), "versions.config.json");
      if (existsSync(configPath)) {
        const configContent = await readFile(configPath, "utf-8");
        versionConfig = JSON.parse(configContent);
      }
    } catch (error) {
      console.warn("Failed to load versions config:", error);
    }

    // Add current/latest version
    const currentArtifactsPath = join(
      process.cwd(),
      ".artifacts",
      "versions",
      "main",
      "zitadel"
    );

    const isDefaultLatest =
      !versionConfig?.settings?.defaultVersion ||
      versionConfig.settings.defaultVersion === "latest";

    versions.push({
      id: "latest",
      name: "Latest (Current Branch)",
      isDefault: isDefaultLatest,
      available: existsSync(currentArtifactsPath),
    });

    // Check for organized version folders (.artifacts/versions/)
    const versionsDir = join(process.cwd(), ".artifacts", "versions");
    const foundVersions = new Set<string>();

    if (existsSync(versionsDir)) {
      try {
        const versionDirs = await readdir(versionsDir);

        for (const versionDir of versionDirs) {
          foundVersions.add(versionDir);
          const versionPath = join(versionsDir, versionDir, "zitadel");
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

          const isConfigDefault =
            versionConfig?.settings?.defaultVersion === versionDir;

          versions.push({
            id: versionDir,
            name: displayName,
            isDefault: isConfigDefault,
            available: existsSync(versionPath),
            metadata,
          });
        }
      } catch (error) {
        console.error("Error reading version directories:", error);
      }
    }

    // Add versions from config that weren't found in filesystem (fallback)
    if (versionConfig) {
      for (const configVersion of versionConfig.versions) {
        if (
          configVersion.enabled &&
          !foundVersions.has(configVersion.id) &&
          configVersion.id !== "latest"
        ) {
          const isConfigDefault =
            versionConfig.settings.defaultVersion === configVersion.id;

          versions.push({
            id: configVersion.id,
            name: configVersion.name,
            isDefault: isConfigDefault,
            available: false, // Not available since artifacts don't exist
            metadata: {
              version: configVersion.id,
              gitBranch: configVersion.gitRef,
            },
          });
        }
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
