import {expect, Page} from "@playwright/test";
import {CDPSession} from "playwright-core";

const BASE64_ENCODED_PK =
    "MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQDbBOu5Lhs4vpowbCnmCyLUpIE7JM9sm9QXzye2G+jr+Kr" +
    "MsinWohEce47BFPJlTaDzHSvOW2eeunBO89ZcvvVc8RLz4qyQ8rO98xS1jtgqi1NcBPETDrtzthODu/gd0sjB2Tk3TLuBGV" +
    "oPXt54a+Oo4JbBJ6h3s0+5eAfGplCbSNq6hN3Jh9YOTw5ZA6GCEy5l8zBaOgjXytd2v2OdSVoEDNiNQRkjJd2rmS2oi9AyQ" +
    "FR3B7BrPSiDlCcITZFOWgLF5C31Wp/PSHwQhlnh7/6YhnE2y9tzsUvzx0wJXrBADW13+oMxrneDK3WGbxTNYgIi1PvSqXlq" +
    "GjHtCK+R2QkXAgMBAAECggEAVc6bu7VAnP6v0gDOeX4razv4FX/adCao9ZsHZ+WPX8PQxtmWYqykH5CY4TSfsuizAgyPuQ0" +
    "+j4Vjssr9VODLqFoanspT6YXsvaKanncUYbasNgUJnfnLnw3an2XpU2XdmXTNYckCPRX9nsAAURWT3/n9ljc/XYY22ecYxM" +
    "8sDWnHu2uKZ1B7M3X60bQYL5T/lVXkKdD6xgSNLeP4AkRx0H4egaop68hoW8FIwmDPVWYVAvo8etzWCtibRXz5FcNld9MgD" +
    "/Ai7ycKy4Q1KhX5GBFI79MVVaHkSQfxPHpr7/XcmpQOEAr+BMPon4s4vnKqAGdGB3j/E3d/+4F2swykoQKBgQD8hCsp6FIQ" +
    "5umJlk9/j/nGsMl85LgLaNVYpWlPRKPc54YNumtvj5vx1BG+zMbT7qIE3nmUPTCHP7qb5ERZG4CdMCS6S64/qzZEqijLCqe" +
    "pwj6j4fV5SyPWEcpxf6ehNdmcfgzVB3Wolfwh1ydhx/96L1jHJcTKchdJJzlfTvq8wwKBgQDeCnKws1t5GapfE1rmC/h4ol" +
    "L2qZTth9oQmbrXYohVnoqNFslDa43ePZwL9Jmd9kYb0axOTNMmyrP0NTj41uCfgDS0cJnNTc63ojKjegxHIyYDKRZNVUR/d" +
    "xAYB/vPfBYZUS7M89pO6LLsHhzS3qpu3/hppo/Uc/AM/r8PSflNHQKBgDnWgBh6OQncChPUlOLv9FMZPR1ZOfqLCYrjYEqi" +
    "uzGm6iKM13zXFO4AGAxu1P/IAd5BovFcTpg79Z8tWqZaUUwvscnl+cRlj+mMXAmdqCeO8VASOmqM1ml667axeZDIR867ZG8" +
    "K5V029Wg+4qtX5uFypNAAi6GfHkxIKrD04yOHAoGACdh4wXESi0oiDdkz3KOHPwIjn6BhZC7z8mx+pnJODU3cYukxv3WTct" +
    "lUhAsyjJiQ/0bK1yX87ulqFVgO0Knmh+wNajrb9wiONAJTMICG7tiWJOm7fW5cfTJwWkBwYADmkfTRmHDvqzQSSvoC2S7aa" +
    "9QulbC3C/qgGFNrcWgcT9kCgYAZTa1P9bFCDU7hJc2mHwJwAW7/FQKEJg8SL33KINpLwcR8fqaYOdAHWWz636osVEqosRrH" +
    "zJOGpf9x2RSWzQJ+dq8+6fACgfFZOVpN644+sAHfNPAI/gnNKU5OfUv+eav8fBnzlf1A3y3GIkyMyzFN3DE7e0n/lyqxE4H" +
    "BYGpI8g==";

interface session {
    client: CDPSession
    authenticatorId: string
}

async function client(page: Page): Promise<session> {
    const cdpSession = await page.context().newCDPSession(page);
    await cdpSession.send('WebAuthn.enable', {enableUI: false});
    const result = await cdpSession.send('WebAuthn.addVirtualAuthenticator', {
        options: {
            protocol: 'ctap2',
            transport: 'internal',
            hasResidentKey: true,
            hasUserVerification: true,
            isUserVerified: true,
            automaticPresenceSimulation: true,
        },
    });
    return {client: cdpSession, authenticatorId: result.authenticatorId};
}

export async function passkeyRegister(page: Page): Promise<string> {
    const session = await client(page)

    await passkeyNotExisting(session.client, session.authenticatorId);
    await simulateSuccessfulPasskeyRegister(
        session.client,
        session.authenticatorId,
        () =>
            page.getByTestId("submit-button").click()
    );
    await passkeyRegistered(session.client, session.authenticatorId);

    return session.authenticatorId
}

export async function passkey(page: Page, authenticatorId: string) {
    const cdpSession = await page.context().newCDPSession(page);
    await cdpSession.send('WebAuthn.enable', {enableUI: false});

    const signCount = await passkeyExisting(cdpSession, authenticatorId);

    await simulateSuccessfulPasskeyInput(
        cdpSession,
        authenticatorId,
        () =>
            page.getByTestId("submit-button").click()
    );

    await passkeyUsed(cdpSession, authenticatorId, signCount);
}

async function passkeyNotExisting(client: CDPSession, authenticatorId: string) {
    const result = await client.send('WebAuthn.getCredentials', {authenticatorId});
    expect(result.credentials).toHaveLength(0);
}

async function passkeyRegistered(client: CDPSession, authenticatorId: string) {
    const result = await client.send('WebAuthn.getCredentials', {authenticatorId});
    expect(result.credentials).toHaveLength(1);
    await passkeyUsed(client, authenticatorId, 0);
}

async function passkeyExisting(client: CDPSession, authenticatorId: string): Promise<number> {
    const result = await client.send('WebAuthn.getCredentials', {authenticatorId});
    expect(result.credentials).toHaveLength(1);
    return result.credentials[0].signCount
}

async function passkeyUsed(client: CDPSession, authenticatorId: string, signCount: number) {
    const result = await client.send('WebAuthn.getCredentials', {authenticatorId});
    expect(result.credentials).toHaveLength(1);
    expect(result.credentials[0].signCount).toBeGreaterThan(signCount);
}

async function simulateSuccessfulPasskeyRegister(client: CDPSession, authenticatorId: string, operationTrigger: () => Promise<void>) {
    // initialize event listeners to wait for a successful passkey input event
    const operationCompleted = new Promise<void>(resolve => {
        client.on('WebAuthn.credentialAdded', () => {
            console.log('Credential Added!');
            resolve()
        });
    });

    // perform a user action that triggers passkey prompt
    await operationTrigger();

    // wait to receive the event that the passkey was successfully registered or verified
    await operationCompleted;
}

async function simulateSuccessfulPasskeyInput(client: CDPSession, authenticatorId: string, operationTrigger: () => Promise<void>) {
    // initialize event listeners to wait for a successful passkey input event
    const operationCompleted = new Promise<void>(resolve => {
        client.on('WebAuthn.credentialAsserted', () => {
            console.log('Credential Asserted!');
            resolve()
        });
    });

    // perform a user action that triggers passkey prompt
    await operationTrigger();

    // wait to receive the event that the passkey was successfully registered or verified
    await operationCompleted;
}

async function simulateFailedPasskeyInput(client: CDPSession, authenticatorId: string, operationTrigger: () => Promise<void>, postOperationCheck: () => Promise<void>) {
    // set isUserVerified option to false
    // (so that subsequent passkey operations will fail)
    await client.send('WebAuthn.setUserVerified', {
        authenticatorId: authenticatorId,
        isUserVerified: false,
    });

    // set automaticPresenceSimulation option to true
    // (so that the virtual authenticator will respond to the next passkey prompt)
    await client.send('WebAuthn.setAutomaticPresenceSimulation', {
        authenticatorId: authenticatorId,
        enabled: true,
    });

    // perform a user action that triggers passkey prompt
    await operationTrigger();

    // wait for an expected UI change that indicates the passkey operation has completed
    await postOperationCheck();

    // set automaticPresenceSimulation option back to false
    await client.send('WebAuthn.setAutomaticPresenceSimulation', {
        authenticatorId,
        enabled: false,
    });
}

