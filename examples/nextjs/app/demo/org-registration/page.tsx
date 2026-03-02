"use client";

import type { FormEvent } from "react";
import { useState } from "react";

interface RegisterOrganizationResponse {
  organizationId: string;
  createdAdmins: Array<{
    userId: string;
    emailCode?: string | null;
    phoneCode?: string | null;
  }>;
}

async function readPayload(response: Response): Promise<unknown> {
  const bodyText = await response.text();
  if (!bodyText) {
    return null;
  }

  try {
    return JSON.parse(bodyText) as unknown;
  } catch {
    return { raw: bodyText };
  }
}

export default function OrgRegistrationDemoPage() {
  const [name, setName] = useState("");
  const [organizationId, setOrganizationId] = useState("");
  const [result, setResult] = useState<RegisterOrganizationResponse | null>(null);
  const [apiError, setApiError] = useState<unknown>(null);
  const [statusCode, setStatusCode] = useState<number | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setIsSubmitting(true);
    setApiError(null);
    setResult(null);
    setStatusCode(null);

    const payload: { name: string; organizationId?: string } = {
      name: name.trim(),
    };
    const normalizedOrganizationId = organizationId.trim();
    if (normalizedOrganizationId) {
      payload.organizationId = normalizedOrganizationId;
    }

    try {
      const response = await fetch("/api/demo/org-registration/register", {
        method: "POST",
        headers: {
          "content-type": "application/json",
        },
        body: JSON.stringify(payload),
      });
      const body = await readPayload(response);
      setStatusCode(response.status);

      if (!response.ok) {
        setApiError(body);
        return;
      }

      setResult(body as RegisterOrganizationResponse);
    } catch (error) {
      setStatusCode(0);
      setApiError({
        error: "Request failed before reaching the API route",
        message: error instanceof Error ? error.message : "Unknown network error",
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <section className="ztdl-lane ztdl-noise">
      <h2>Organization registration demo lane</h2>
      <p>
        This calls the v2 Organization API (<code>organizationService.addOrganization</code>) via{" "}
        <code>/api/demo/org-registration/register</code>.
      </p>
      <p>
        Sign in first (<a href="/api/auth/signin">/api/auth/signin</a>) or provide service-user credentials for{" "}
        <code>createZitadelApiClient</code>.
      </p>
      <h3>Permission guidance</h3>
      <ul>
        <li>
          Required permission for this demo: <code>org.create</code>
        </li>
        <li>
          <code>403</code> usually means the token/session is missing <code>org.create</code>.
        </li>
        <li>
          <code>409</code> usually means the organization name or ID is already in use.
        </li>
      </ul>
      <form onSubmit={handleSubmit}>
        <p>
          <label>
            Organization name
            <br />
            <input
              name="name"
              value={name}
              onChange={(event) => setName(event.target.value)}
              required
              autoComplete="organization"
            />
          </label>
        </p>
        <p>
          <label>
            Organization ID (optional)
            <br />
            <input
              name="organizationId"
              value={organizationId}
              onChange={(event) => setOrganizationId(event.target.value)}
              autoComplete="off"
            />
          </label>
        </p>
        <button type="submit" disabled={isSubmitting}>
          {isSubmitting ? "Creating organization..." : "Create organization"}
        </button>
      </form>
      {result && (
        <>
          <h3>Registration result{statusCode ? ` (HTTP ${statusCode})` : ""}</h3>
          <pre>{JSON.stringify(result, null, 2)}</pre>
        </>
      )}
      {apiError !== null && (
        <>
          <h3>API error{statusCode !== null ? ` (HTTP ${statusCode})` : ""}</h3>
          <pre>{JSON.stringify(apiError, null, 2)}</pre>
        </>
      )}
    </section>
  );
}
