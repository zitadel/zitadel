"use client";

import { useEffect, useState, useRef } from "react";
import { createApiReference } from "@scalar/api-reference";
import yaml from "js-yaml";

interface OpenApiSpec {
  name: string;
  fileName: string;
  content: string;
}

interface ApiResponse {
  specs: OpenApiSpec[];
}

export function ApiReferenceComponent() {
  const [specs, setSpecs] = useState<OpenApiSpec[]>([]);
  const [selectedSpec, setSelectedSpec] = useState<string>("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    async function loadSpecs() {
      try {
        const response = await fetch("/api/openapi");
        if (!response.ok) {
          throw new Error("Failed to load API specifications");
        }
        const data: ApiResponse = await response.json();

        // Filter out specs with no endpoints (only schema definitions)
        const specsWithEndpoints = data.specs.filter((spec) => {
          try {
            const parsed = yaml.load(spec.content) as any;
            return parsed.paths && Object.keys(parsed.paths).length > 0;
          } catch {
            return false;
          }
        });

        setSpecs(specsWithEndpoints);

        if (specsWithEndpoints.length > 0) {
          // Default to a v2 user service if available, otherwise management service
          const userV2Service = specsWithEndpoints.find((spec) =>
            spec.name.includes("user/v2/user_service")
          );
          const managementService = specsWithEndpoints.find(
            (spec) => spec.name === "management"
          );

          if (userV2Service) {
            setSelectedSpec(userV2Service.name);
          } else if (managementService) {
            setSelectedSpec("management");
          } else {
            setSelectedSpec(specsWithEndpoints[0].name);
          }
        }
      } catch (err) {
        setError(err instanceof Error ? err.message : "Unknown error");
      } finally {
        setLoading(false);
      }
    }

    loadSpecs();
  }, []);

  useEffect(() => {
    if (selectedSpec && containerRef.current) {
      const selectedSpecData = specs.find((spec) => spec.name === selectedSpec);
      if (selectedSpecData) {
        try {
          const parsedSpec = yaml.load(selectedSpecData.content) as any;

          // Debug: Log the parsed spec
          console.log("Selected spec:", selectedSpec);
          console.log("Parsed spec:", parsedSpec);
          console.log("Parsed spec paths:", parsedSpec.paths);
          console.log(
            "Number of paths:",
            Object.keys(parsedSpec.paths || {}).length
          );

          // Clear the container
          containerRef.current.innerHTML = "";

          // Create a div for Scalar
          const scalarDiv = document.createElement("div");
          scalarDiv.id = `api-reference-${selectedSpec}`;
          containerRef.current.appendChild(scalarDiv);

          // Create the API reference with correct configuration
          createApiReference(scalarDiv, {
            spec: {
              content: parsedSpec,
            },
            theme: "default",
          } as any);
        } catch (err) {
          console.error("Error parsing YAML or creating API reference:", err);
          if (containerRef.current) {
            containerRef.current.innerHTML = `<div style="padding: 20px; color: red;">Error loading API documentation: ${err}</div>`;
          }
        }
      }
    }
  }, [selectedSpec, specs]);

  if (loading) {
    return (
      <div
        style={{
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          height: "100vh",
          fontSize: "18px",
        }}
      >
        Loading API documentation...
      </div>
    );
  }

  if (error) {
    return (
      <div
        style={{
          display: "flex",
          flexDirection: "column",
          justifyContent: "center",
          alignItems: "center",
          height: "100vh",
          fontSize: "18px",
          color: "#e74c3c",
        }}
      >
        <h2>Error loading API documentation</h2>
        <p>{error}</p>
        <p style={{ marginTop: "20px", fontSize: "14px", color: "#666" }}>
          Make sure to run <code>pnpm run generate</code> to generate the
          OpenAPI specifications.
        </p>
      </div>
    );
  }

  if (specs.length === 0) {
    return (
      <div
        style={{
          display: "flex",
          flexDirection: "column",
          justifyContent: "center",
          alignItems: "center",
          height: "100vh",
          fontSize: "18px",
        }}
      >
        <h2>No API specifications found</h2>
        <p style={{ marginTop: "20px", fontSize: "14px", color: "#666" }}>
          Run <code>pnpm run generate</code> to generate the OpenAPI
          specifications from proto files.
        </p>
      </div>
    );
  }

  const getServiceDisplayName = (serviceName: string): string => {
    // Handle service paths like "user/v2/user_service"
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

    // Handle v1 services
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
  };

  const sortedSpecs = [...specs].sort((a, b) => {
    // Sort by service name first, then by version (v2 before v1)
    const aService = a.name.split("/")[0];
    const bService = b.name.split("/")[0];

    if (aService !== bService) {
      return aService.localeCompare(bService);
    }

    // Same service, sort by version (v2+ first)
    const aIsV1 = !a.name.includes("/");
    const bIsV1 = !b.name.includes("/");

    if (aIsV1 && !bIsV1) return 1; // v1 comes after v2+
    if (!aIsV1 && bIsV1) return -1; // v2+ comes before v1

    return a.name.localeCompare(b.name);
  });

  return (
    <div style={{ display: "flex", height: "100vh" }}>
      {/* Sidebar for service selection */}
      <div
        style={{
          width: "250px",
          borderRight: "1px solid #e1e4e8",
          padding: "20px",
          overflowY: "auto",
          backgroundColor: "#f6f8fa",
        }}
      >
        <h2
          style={{ marginBottom: "20px", fontSize: "18px", fontWeight: "bold" }}
        >
          ZITADEL API Services
        </h2>
        <nav>
          {sortedSpecs.map((spec) => (
            <button
              key={spec.name}
              onClick={() => setSelectedSpec(spec.name)}
              style={{
                display: "block",
                width: "100%",
                padding: "10px",
                marginBottom: "8px",
                border: "none",
                borderRadius: "6px",
                backgroundColor:
                  selectedSpec === spec.name ? "#0969da" : "transparent",
                color: selectedSpec === spec.name ? "white" : "#24292f",
                textAlign: "left",
                cursor: "pointer",
                fontSize: "14px",
                transition: "background-color 0.2s",
              }}
              onMouseEnter={(e) => {
                if (selectedSpec !== spec.name) {
                  e.currentTarget.style.backgroundColor = "#f3f4f6";
                }
              }}
              onMouseLeave={(e) => {
                if (selectedSpec !== spec.name) {
                  e.currentTarget.style.backgroundColor = "transparent";
                }
              }}
            >
              {getServiceDisplayName(spec.name)}
            </button>
          ))}
        </nav>
      </div>

      {/* Main content area */}
      <div
        ref={containerRef}
        style={{ flex: 1, height: "100vh", overflow: "auto" }}
      />
    </div>
  );
}
