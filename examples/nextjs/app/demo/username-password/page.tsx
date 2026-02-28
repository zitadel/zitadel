"use client";

import { UsernamePasswordForm } from "@zitadel/react";
import { useState, type FormEvent } from "react";

interface SessionResult {
  sessionId: string;
  sessionToken: string;
  details?: unknown;
  challenges?: unknown;
}

interface CallbackResult {
  callbackUrl: string;
}

interface ApiError {
  error?: string;
  details?: unknown;
}

async function requestJson<T>(input: string, init?: RequestInit): Promise<T> {
  const response = await fetch(input, init);
  const payload = (await response.json().catch(() => ({}))) as T & ApiError;

  if (!response.ok) {
    const detailText = payload.details
      ? ` (${JSON.stringify(payload.details)})`
      : "";
    throw new Error(
      (payload.error ?? `Request failed with status ${response.status}`) + detailText,
    );
  }

  return payload;
}

export default function UsernamePasswordDemoPage() {
  const [loginName, setLoginName] = useState("");
  const [password, setPassword] = useState("");
  const [authRequestId, setAuthRequestId] = useState("");
  const [sessionId, setSessionId] = useState("");
  const [sessionToken, setSessionToken] = useState("");

  const [sessionError, setSessionError] = useState<string | null>(null);
  const [callbackError, setCallbackError] = useState<string | null>(null);
  const [inspectError, setInspectError] = useState<string | null>(null);

  const [sessionResult, setSessionResult] = useState<unknown>(null);
  const [callbackResult, setCallbackResult] = useState<unknown>(null);
  const [inspectResult, setInspectResult] = useState<unknown>(null);

  const [isCreatingSession, setIsCreatingSession] = useState(false);
  const [isCreatingCallback, setIsCreatingCallback] = useState(false);
  const [isInspectingSession, setIsInspectingSession] = useState(false);

  async function onCreateSession(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setSessionError(null);
    setSessionResult(null);
    setInspectResult(null);

    setIsCreatingSession(true);
    try {
      const result = await requestJson<SessionResult>(
        "/api/demo/username-password/session",
        {
          method: "POST",
          headers: {
            "content-type": "application/json",
          },
          body: JSON.stringify({ loginName, password }),
        },
      );

      setSessionId(result.sessionId);
      setSessionToken(result.sessionToken);
      setSessionResult(result);
    } catch (error) {
      setSessionError(error instanceof Error ? error.message : "Failed to create session.");
    } finally {
      setIsCreatingSession(false);
    }
  }

  async function onCreateCallback(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setCallbackError(null);
    setCallbackResult(null);

    setIsCreatingCallback(true);
    try {
      const result = await requestJson<CallbackResult>(
        "/api/demo/username-password/callback",
        {
          method: "POST",
          headers: {
            "content-type": "application/json",
          },
          body: JSON.stringify({ authRequestId, sessionId, sessionToken }),
        },
      );

      setCallbackResult(result);
    } catch (error) {
      setCallbackError(error instanceof Error ? error.message : "Failed to create callback.");
    } finally {
      setIsCreatingCallback(false);
    }
  }

  async function onInspectSession() {
    setInspectError(null);
    setInspectResult(null);

    const params = new URLSearchParams({ sessionId });
    if (sessionToken) {
      params.set("sessionToken", sessionToken);
    }

    setIsInspectingSession(true);
    try {
      const result = await requestJson<unknown>(
        `/api/demo/username-password/session?${params.toString()}`,
      );
      setInspectResult(result);
    } catch (error) {
      setInspectError(error instanceof Error ? error.message : "Failed to inspect session.");
    } finally {
      setIsInspectingSession(false);
    }
  }

  return (
    <section className="ztdl-lane">
      <h2>Username/password demo lane</h2>
      <p>Creates a verified session from username/password and optionally builds an OIDC callback URL.</p>

      <UsernamePasswordForm
        error={sessionError}
        isLoading={isCreatingSession}
        loginName={loginName}
        onLoginNameChange={(event) => setLoginName(event.currentTarget.value)}
        onPasswordChange={(event) => setPassword(event.currentTarget.value)}
        onSubmit={onCreateSession}
        password={password}
      />

      <h4>Session result</h4>
      <pre>{JSON.stringify(sessionResult ?? { status: "No session created yet." }, null, 2)}</pre>

      <p>
        <button
          disabled={!sessionId || isInspectingSession}
          onClick={onInspectSession}
          type="button"
        >
          {isInspectingSession ? "Loading session..." : "Inspect current session"}
        </button>
      </p>

      {inspectError ? (
        <p className="ztdl-status-error" role="alert">
          <strong>Inspect error:</strong> {inspectError}
        </p>
      ) : null}

      <h4>Session inspection result</h4>
      <pre>{JSON.stringify(inspectResult ?? { status: "Not fetched yet." }, null, 2)}</pre>

      <form onSubmit={onCreateCallback}>
        <h3>2) Create callback URL</h3>
        <p>
          <label>
            OIDC auth request ID
            <br />
            <input
              name="authRequestId"
              onChange={(event) => setAuthRequestId(event.currentTarget.value)}
              placeholder="oidc_..."
              required
              value={authRequestId}
            />
          </label>
        </p>
        <p>
          <label>
            Session ID
            <br />
            <input
              name="sessionId"
              onChange={(event) => setSessionId(event.currentTarget.value)}
              required
              value={sessionId}
            />
          </label>
        </p>
        <p>
          <label>
            Session token
            <br />
            <input
              name="sessionToken"
              onChange={(event) => setSessionToken(event.currentTarget.value)}
              required
              value={sessionToken}
            />
          </label>
        </p>
        <p>
          <button disabled={isCreatingCallback} type="submit">
            {isCreatingCallback ? "Creating callback..." : "Create callback URL"}
          </button>
        </p>
      </form>

      {callbackError ? (
        <p className="ztdl-status-error" role="alert">
          <strong>Callback error:</strong> {callbackError}
        </p>
      ) : null}

      <h4>Callback result</h4>
      <pre>{JSON.stringify(callbackResult ?? { status: "No callback created yet." }, null, 2)}</pre>
      {callbackResult && typeof callbackResult === "object" && "callbackUrl" in callbackResult ? (
        <p>
          <a href={(callbackResult as CallbackResult).callbackUrl}>Open callback URL</a>
        </p>
      ) : null}
    </section>
  );
}
