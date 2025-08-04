import { expect, Page } from "@playwright/test";
import { CDPSession } from "playwright-core";

interface session {
  client: CDPSession;
  authenticatorId: string;
}

async function client(page: Page): Promise<session> {
  const cdpSession = await page.context().newCDPSession(page);
  await cdpSession.send("WebAuthn.enable", { enableUI: false });
  const result = await cdpSession.send("WebAuthn.addVirtualAuthenticator", {
    options: {
      protocol: "ctap2",
      transport: "internal",
      hasResidentKey: true,
      hasUserVerification: true,
      isUserVerified: true,
      automaticPresenceSimulation: true,
    },
  });
  return { client: cdpSession, authenticatorId: result.authenticatorId };
}

export async function passkeyRegister(page: Page): Promise<string> {
  const session = await client(page);

  await passkeyNotExisting(session.client, session.authenticatorId);
  await simulateSuccessfulPasskeyRegister(session.client, session.authenticatorId, () =>
    page.getByTestId("submit-button").click(),
  );
  await passkeyRegistered(session.client, session.authenticatorId);

  return session.authenticatorId;
}

export async function passkey(page: Page, authenticatorId: string) {
  const cdpSession = await page.context().newCDPSession(page);
  await cdpSession.send("WebAuthn.enable", { enableUI: false });

  const signCount = await passkeyExisting(cdpSession, authenticatorId);

  await simulateSuccessfulPasskeyInput(cdpSession, authenticatorId, () => page.getByTestId("submit-button").click());

  await passkeyUsed(cdpSession, authenticatorId, signCount);
}

async function passkeyNotExisting(client: CDPSession, authenticatorId: string) {
  const result = await client.send("WebAuthn.getCredentials", { authenticatorId });
  expect(result.credentials).toHaveLength(0);
}

async function passkeyRegistered(client: CDPSession, authenticatorId: string) {
  const result = await client.send("WebAuthn.getCredentials", { authenticatorId });
  expect(result.credentials).toHaveLength(1);
  await passkeyUsed(client, authenticatorId, 0);
}

async function passkeyExisting(client: CDPSession, authenticatorId: string): Promise<number> {
  const result = await client.send("WebAuthn.getCredentials", { authenticatorId });
  expect(result.credentials).toHaveLength(1);
  return result.credentials[0].signCount;
}

async function passkeyUsed(client: CDPSession, authenticatorId: string, signCount: number) {
  const result = await client.send("WebAuthn.getCredentials", { authenticatorId });
  expect(result.credentials).toHaveLength(1);
  expect(result.credentials[0].signCount).toBeGreaterThan(signCount);
}

async function simulateSuccessfulPasskeyRegister(
  client: CDPSession,
  authenticatorId: string,
  operationTrigger: () => Promise<void>,
) {
  // initialize event listeners to wait for a successful passkey input event
  const operationCompleted = new Promise<void>((resolve) => {
    client.on("WebAuthn.credentialAdded", () => {
      console.log("Credential Added!");
      resolve();
    });
  });

  // perform a user action that triggers passkey prompt
  await operationTrigger();

  // wait to receive the event that the passkey was successfully registered or verified
  await operationCompleted;
}

async function simulateSuccessfulPasskeyInput(
  client: CDPSession,
  authenticatorId: string,
  operationTrigger: () => Promise<void>,
) {
  // initialize event listeners to wait for a successful passkey input event
  const operationCompleted = new Promise<void>((resolve) => {
    client.on("WebAuthn.credentialAsserted", () => {
      console.log("Credential Asserted!");
      resolve();
    });
  });

  // perform a user action that triggers passkey prompt
  await operationTrigger();

  // wait to receive the event that the passkey was successfully registered or verified
  await operationCompleted;
}
