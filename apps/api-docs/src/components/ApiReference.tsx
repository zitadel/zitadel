"use client";

import { useEffect, useState, useRef } from "react";
import { createApiReference } from "@scalar/api-reference";
import yaml from "js-yaml";

// Add CSS to handle scroll offset for fixed header
if (typeof document !== "undefined") {
  const style = document.createElement("style");
  style.textContent = `
    /* Fix scroll anchorin  return (
    <div style={{ 
      height: "100vh", 
      width: "100vw",
      display: "flex", 
      flexDirection: "column",
      overflow: "hidden"
    }}>
      {/* Navigation bar - fixed height */}
      <div
        style={{
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
          padding: "8px 24px",
          backgroundColor: "var(--scalar-background-1, #ffffff)",
          borderBottom: "1px solid var(--scalar-border-color, #e1e4e8)",
          boxShadow: "0 2px 4px rgba(0, 0, 0, 0.1)",
          zIndex: 1000,
          minHeight: "70px",
          flexShrink: 0, // Prevent shrinking
        }}
      > fixed header */
    html {
      scroll-padding-top: 80px;
    }
    
    /* Target Scalar's specific elements that might be used as scroll targets */
    [id]:target,
    .scalar-api-reference [id]:target,
    .scalar-api-reference h1[id],
    .scalar-api-reference h2[id],
    .scalar-api-reference h3[id],
    .scalar-api-reference h4[id],
    .scalar-api-reference h5[id],
    .scalar-api-reference h6[id] {
      scroll-margin-top: 80px;
    }
    
    /* Also fix any operation/endpoint scrolling */
    .scalar-api-reference [data-operation-id],
    .scalar-api-reference [data-section-id] {
      scroll-margin-top: 80px;
    }
  `;

  if (!document.head.querySelector("#scroll-offset-fix")) {
    style.id = "scroll-offset-fix";
    document.head.appendChild(style);
  }
}

interface OpenApiSpec {
  name: string;
  fileName: string;
  content: string;
}

interface ApiResponse {
  specs: OpenApiSpec[];
}

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

interface SearchResponse {
  results: SearchResult[];
  total: number;
}

export function ApiReferenceComponent() {
  const [specs, setSpecs] = useState<OpenApiSpec[]>([]);
  const [selectedSpec, setSelectedSpec] = useState<string>("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState<string>("");
  const [searchResults, setSearchResults] = useState<SearchResult[]>([]);
  const [isSearching, setIsSearching] = useState(false);
  const [showSearchResults, setShowSearchResults] = useState(false);
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
            const parsed = JSON.parse(spec.content);
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

  // Search function with debouncing
  useEffect(() => {
    const timeoutId = setTimeout(async () => {
      if (searchQuery.trim().length >= 2) {
        setIsSearching(true);
        try {
          const response = await fetch(
            `/api/search?q=${encodeURIComponent(searchQuery)}`
          );
          if (response.ok) {
            const data: SearchResponse = await response.json();
            setSearchResults(data.results);
            setShowSearchResults(true);
          }
        } catch (error) {
          console.error("Search error:", error);
          setSearchResults([]);
        } finally {
          setIsSearching(false);
        }
      } else {
        setSearchResults([]);
        setShowSearchResults(false);
      }
    }, 300);

    return () => clearTimeout(timeoutId);
  }, [searchQuery]);

  const handleSearchResultClick = (result: SearchResult) => {
    setSelectedSpec(result.serviceName);
    setShowSearchResults(false);
    setSearchQuery("");

    // Wait for the spec to load, then navigate using Scalar's navigation patterns
    setTimeout(() => {
      if (result.operationId || result.path) {
        const method = result.method.toLowerCase();
        const operationId = result.operationId;
        const tag = result.tags && result.tags.length > 0 ? result.tags[0] : "";

        // Based on Scalar's code, they use patterns like:
        // - tag/{tag}/{method}{path} for operations under tags
        // - operation/{operationId} for direct operation access
        let targetHash = "";

        if (tag && result.path) {
          // This is the most common pattern in Scalar: tag/TagName/method/path
          targetHash = `tag/${tag}/${method}${result.path}`;
        } else if (operationId) {
          // Fallback to operation ID
          targetHash = `operation/${operationId}`;
        } else {
          // Last resort: method + path
          targetHash = `${method}${result.path}`;
        }

        console.log("Navigation attempt:", {
          result,
          targetHash,
          currentLocation: window.location.href,
        });

        // Set the hash and let Scalar handle it
        window.location.hash = targetHash;

        // Fallback: try to find and scroll to elements after Scalar processes the hash
        setTimeout(() => {
          // Try to find any element with the operation ID or path
          const possibleSelectors = [
            `[id="${targetHash}"]`,
            `[id="${operationId}"]`,
            `[data-operation-id="${operationId}"]`,
            `[id*="${operationId}"]`,
            // Look for the operation in the content area
            `.scalar-api-reference [id*="${method}"][id*="${result.path.replace(
              /\//g,
              ""
            )}"]`,
            // Try to find by text content (operation summary)
            ...Array.from(document.querySelectorAll("h1, h2, h3, h4, h5, h6"))
              .filter(
                (el) =>
                  el.textContent?.includes(result.summary || "") &&
                  result.summary
              )
              .map((el) => `#${el.id}`)
              .filter((id) => id !== "#"),
            // Look for method badges
            `[data-method="${method}"]`,
          ];

          console.log(
            "Searching for elements with selectors:",
            possibleSelectors.slice(0, 5)
          );

          for (const selector of possibleSelectors) {
            if (!selector || selector === "#") continue;

            try {
              const element = document.querySelector(selector);
              if (element && (element as HTMLElement).offsetParent !== null) {
                // Check if element is visible
                console.log(
                  "Found and scrolling to element:",
                  selector,
                  element
                );
                element.scrollIntoView({
                  behavior: "smooth",
                  block: "start",
                });

                // Highlight the element briefly to confirm navigation
                const htmlElement = element as HTMLElement;
                const originalStyle = htmlElement.style.border;
                htmlElement.style.border = "2px solid #007bff";
                setTimeout(() => {
                  htmlElement.style.border = originalStyle;
                }, 2000);

                return; // Exit once we find and scroll to an element
              }
            } catch (e) {
              console.warn("Selector failed:", selector, e);
            }
          }

          // If nothing worked, log what's available for debugging
          console.log(
            "No elements found. Available IDs in document:",
            Array.from(document.querySelectorAll("[id]"))
              .map((el) => el.id)
              .filter((id) => id)
              .slice(0, 20)
          );
        }, 1000);
      }
    }, 1500);
  };

  useEffect(() => {
    if (selectedSpec && containerRef.current) {
      const selectedSpecData = specs.find((spec) => spec.name === selectedSpec);
      if (selectedSpecData) {
        try {
          const parsedSpec = JSON.parse(selectedSpecData.content);

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
            theme: "github",
            layout: "modern",
            customCss: `
              /* Fix Scalar's sidebar height to work with our navigation bar */
              .scalar-api-reference {
                height: 100% !important;
              }
              
              /* Adjust sidebar wrapper to account for our navigation */
              .scalar-api-reference .sidebar,
              .scalar-api-reference .scalar-api-reference__sidebar,
              .scalar-api-reference [data-sidebar] {
                max-height: calc(100vh - 70px) !important;
                top: 0 !important;
              }
              
              /* Ensure sidebar content scrolls properly */
              .scalar-api-reference .sidebar-content,
              .scalar-api-reference .scalar-api-reference__sidebar-content,
              .scalar-api-reference .sidebar .scalar-api-reference__navigation {
                max-height: calc(100vh - 70px) !important;
                overflow-y: auto !important;
              }
              
              /* Let Scalar handle its own positioning but fix the viewport calculation */
              .scalar-api-reference .scalar-api-reference__container {
                height: 100% !important;
              }
            `,
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
    <div style={{ height: "100vh", display: "flex", flexDirection: "column" }}>
      {/* Top navigation bar */}
      <div
        style={{
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
          padding: "6px 24px",
          backgroundColor: "var(--scalar-background-1, #ffffff)",
          borderBottom: "1px solid var(--scalar-border-color, #e1e4e8)",
          boxShadow: "0 2px 4px rgba(0, 0, 0, 0.1)",
          zIndex: 1000,
          minHeight: "70px",
        }}
      >
        {/* Left side: Logo, Title and Service selector */}
        <div style={{ display: "flex", alignItems: "center", gap: "16px" }}>
          {/* ZITADEL Logo */}
          <a
            href="/"
            style={{
              display: "flex",
              alignItems: "center",
              textDecoration: "none",
              cursor: "pointer",
            }}
            title="Go to homepage"
          >
            <img
              src="/zitadel-logo-light@2x.png"
              alt="ZITADEL"
              width="160"
              height="48"
              style={{
                marginRight: "8px",
                objectFit: "contain",
              }}
            />
          </a>

          <div
            style={{
              backgroundColor: "var(--scalar-background-2, #f6f8fa)",
              border: "1px solid var(--scalar-border-color, #e1e4e8)",
              borderRadius: "6px",
              padding: "8px 12px",
            }}
          >
            <select
              value={selectedSpec}
              onChange={(e) => setSelectedSpec(e.target.value)}
              style={{
                backgroundColor: "transparent",
                border: "none",
                color: "var(--scalar-color-1, #24292f)",
                fontSize: "14px",
                fontWeight: "500",
                cursor: "pointer",
                outline: "none",
                minWidth: "200px",
                padding: "4px 8px",
              }}
            >
              {sortedSpecs.map((spec) => (
                <option key={spec.name} value={spec.name}>
                  {getServiceDisplayName(spec.name)}
                </option>
              ))}
            </select>
          </div>
        </div>

        {/* Right side: Global search */}
        <div style={{ position: "relative" }}>
          <div
            style={{
              backgroundColor: "var(--scalar-background-2, #f6f8fa)",
              border: "1px solid var(--scalar-border-color, #e1e4e8)",
              borderRadius: "6px",
              padding: "6px 12px",
              minWidth: "300px",
            }}
          >
            <div style={{ position: "relative" }}>
              <input
                type="text"
                placeholder="Search across all APIs..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                style={{
                  width: "100%",
                  padding: "8px 12px",
                  border: "1px solid var(--scalar-border-color, #e1e4e8)",
                  borderRadius: "4px",
                  fontSize: "14px",
                  backgroundColor: "var(--scalar-background-1, #ffffff)",
                  color: "var(--scalar-color-1, #24292f)",
                  outline: "none",
                }}
              />
              {isSearching && (
                <div
                  style={{
                    position: "absolute",
                    right: "8px",
                    top: "50%",
                    transform: "translateY(-50%)",
                    fontSize: "12px",
                    color: "var(--scalar-color-2, #666)",
                  }}
                >
                  Searching...
                </div>
              )}
            </div>

            {/* Search results dropdown */}
            {showSearchResults && searchResults.length > 0 && (
              <div
                style={{
                  position: "absolute",
                  top: "100%",
                  left: "0",
                  right: "0",
                  backgroundColor: "var(--scalar-background-1, #ffffff)",
                  border: "1px solid var(--scalar-border-color, #e1e4e8)",
                  borderRadius: "6px",
                  boxShadow: "0 4px 12px rgba(0, 0, 0, 0.15)",
                  maxHeight: "400px",
                  overflowY: "auto",
                  marginTop: "4px",
                  zIndex: 1001,
                }}
              >
                {searchResults.map((result, index) => (
                  <div
                    key={`${result.serviceName}-${result.path}-${result.method}-${index}`}
                    onClick={() => handleSearchResultClick(result)}
                    style={{
                      padding: "12px",
                      borderBottom:
                        index < searchResults.length - 1
                          ? "1px solid var(--scalar-border-color, #e1e4e8)"
                          : "none",
                      cursor: "pointer",
                      transition: "background-color 0.2s",
                    }}
                    onMouseEnter={(e) => {
                      e.currentTarget.style.backgroundColor =
                        "var(--scalar-background-2, #f6f8fa)";
                    }}
                    onMouseLeave={(e) => {
                      e.currentTarget.style.backgroundColor = "transparent";
                    }}
                  >
                    <div
                      style={{
                        display: "flex",
                        alignItems: "center",
                        gap: "8px",
                        marginBottom: "4px",
                      }}
                    >
                      <span
                        style={{
                          backgroundColor:
                            result.method === "GET"
                              ? "#28a745"
                              : result.method === "POST"
                              ? "#007bff"
                              : result.method === "PUT"
                              ? "#ffc107"
                              : result.method === "DELETE"
                              ? "#dc3545"
                              : "#6c757d",
                          color: "white",
                          fontSize: "10px",
                          fontWeight: "bold",
                          padding: "2px 6px",
                          borderRadius: "3px",
                          minWidth: "45px",
                          textAlign: "center",
                        }}
                      >
                        {result.method}
                      </span>
                      <span
                        style={{
                          fontFamily: "monospace",
                          fontSize: "13px",
                          color: "var(--scalar-color-1, #24292f)",
                          fontWeight: "500",
                        }}
                      >
                        {result.path}
                      </span>
                    </div>
                    <div
                      style={{
                        fontSize: "12px",
                        color: "var(--scalar-color-2, #666)",
                        marginBottom: "2px",
                      }}
                    >
                      {result.serviceDisplayName}
                    </div>
                    {result.summary && (
                      <div
                        style={{
                          fontSize: "12px",
                          color: "var(--scalar-color-2, #666)",
                          fontStyle: "italic",
                        }}
                      >
                        {result.summary}
                      </div>
                    )}
                  </div>
                ))}
              </div>
            )}

            {showSearchResults &&
              searchResults.length === 0 &&
              searchQuery.trim().length >= 2 &&
              !isSearching && (
                <div
                  style={{
                    position: "absolute",
                    top: "100%",
                    left: "0",
                    right: "0",
                    backgroundColor: "var(--scalar-background-1, #ffffff)",
                    border: "1px solid var(--scalar-border-color, #e1e4e8)",
                    borderRadius: "6px",
                    boxShadow: "0 4px 12px rgba(0, 0, 0, 0.15)",
                    padding: "12px",
                    marginTop: "4px",
                    fontSize: "12px",
                    color: "var(--scalar-color-2, #666)",
                    textAlign: "center",
                  }}
                >
                  No results found for "{searchQuery}"
                </div>
              )}
          </div>
        </div>
      </div>

      {/* API Documentation content - fills remaining space */}
      <div
        ref={containerRef}
        style={{
          flex: 1,
          overflow: "auto", // Allow scrolling in the content area
          minHeight: 0, // Important: allows flex item to shrink below content size
        }}
        onClick={() => setShowSearchResults(false)}
      />
    </div>
  );
}
