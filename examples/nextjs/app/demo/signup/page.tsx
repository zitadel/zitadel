"use client";

import { type FormEvent, useState } from "react";

type SignupPayload = {
  email: string;
  password: string;
  givenName: string;
  familyName: string;
  username?: string;
};

type SignupApiResponse = {
  ok: boolean;
  request?: unknown;
  result?: unknown;
  error?: {
    message?: string;
    code?: string;
    rawMessage?: string;
  };
};

const initialFormState = {
  email: "",
  password: "",
  givenName: "",
  familyName: "",
  username: "",
};

export default function SignupDemoPage() {
  const [formState, setFormState] = useState(initialFormState);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [lastStatus, setLastStatus] = useState<number | null>(null);
  const [lastRequest, setLastRequest] = useState<unknown>(null);
  const [lastResult, setLastResult] = useState<unknown>(null);
  const [lastError, setLastError] = useState<string | null>(null);

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    const payload: SignupPayload = {
      email: formState.email.trim(),
      password: formState.password,
      givenName: formState.givenName.trim(),
      familyName: formState.familyName.trim(),
      ...(formState.username.trim() ? { username: formState.username.trim() } : {}),
    };

    setIsSubmitting(true);
    setLastStatus(null);
    setLastError(null);
    setLastResult(null);
    setLastRequest({
      ...payload,
      password: "********",
      hasPassword: payload.password.length > 0,
    });

    try {
      const response = await fetch("/api/demo/signup/password", {
        method: "POST",
        headers: {
          "content-type": "application/json",
        },
        body: JSON.stringify(payload),
      });

      const data = (await response.json().catch(() => null)) as SignupApiResponse | null;
      setLastStatus(response.status);
      setLastResult(data);

      if (!response.ok || !data?.ok) {
        setLastError(data?.error?.message ?? `Signup failed with HTTP ${response.status}.`);
      }
    } catch (error) {
      setLastError(
        error instanceof Error
          ? error.message
          : "Unexpected error while calling the signup demo API.",
      );
    } finally {
      setIsSubmitting(false);
    }
  }

  const canSubmit = Boolean(
    formState.email.trim() &&
      formState.password &&
      formState.givenName.trim() &&
      formState.familyName.trim() &&
      !isSubmitting,
  );

  return (
    <section className="ztdl-lane ztdl-noise">
      <h2>Signup demo lane</h2>
      <p>
        Live demo for password signup via <code>POST /api/demo/signup/password</code>.
      </p>
      <form onSubmit={onSubmit}>
        <p>
          <label>
            Email
            <br />
            <input
              type="email"
              name="email"
              value={formState.email}
              onChange={(event) => setFormState((current) => ({ ...current, email: event.target.value }))}
              autoComplete="email"
              required
            />
          </label>
        </p>
        <p>
          <label>
            Password
            <br />
            <input
              type="password"
              name="password"
              value={formState.password}
              onChange={(event) => setFormState((current) => ({ ...current, password: event.target.value }))}
              autoComplete="new-password"
              required
            />
          </label>
        </p>
        <p>
          <label>
            Given name
            <br />
            <input
              type="text"
              name="givenName"
              value={formState.givenName}
              onChange={(event) => setFormState((current) => ({ ...current, givenName: event.target.value }))}
              autoComplete="given-name"
              required
            />
          </label>
        </p>
        <p>
          <label>
            Family name
            <br />
            <input
              type="text"
              name="familyName"
              value={formState.familyName}
              onChange={(event) => setFormState((current) => ({ ...current, familyName: event.target.value }))}
              autoComplete="family-name"
              required
            />
          </label>
        </p>
        <p>
          <label>
            Username (optional)
            <br />
            <input
              type="text"
              name="username"
              value={formState.username}
              onChange={(event) => setFormState((current) => ({ ...current, username: event.target.value }))}
              autoComplete="username"
            />
          </label>
        </p>
        <p>
          <button type="submit" disabled={!canSubmit}>
            {isSubmitting ? "Submitting..." : "Create user with password"}
          </button>
        </p>
      </form>
      {lastStatus ? <p>Last HTTP status: {lastStatus}</p> : null}
      {lastError ? <p className="ztdl-status-error">Error: {lastError}</p> : null}
      <h3>Last request</h3>
      <pre>{JSON.stringify(lastRequest ?? { message: "No request submitted yet." }, null, 2)}</pre>
      <h3>Last result</h3>
      <pre>{JSON.stringify(lastResult ?? { message: "No result yet." }, null, 2)}</pre>
    </section>
  );
}
