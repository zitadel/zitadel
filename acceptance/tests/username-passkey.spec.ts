import {test as base} from "@playwright/test";
import path from 'path';
import dotenv from 'dotenv';
import {PasskeyUser} from "./user";

// Read from ".env" file.
dotenv.config({path: path.resolve(__dirname, '.env.local')});

/*
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

const test = base.extend<{ user: PasskeyUser }>({
    user: async ({page}, use) => {

        // Initialize a CDP session for the current page
        const client = await page.context().newCDPSession(page);
        // Enable WebAuthn environment in this session
        await client.send('WebAuthn.enable', {enableUI: true});

        // Attach a virtual authenticator with specific options
        const result = await client.send('WebAuthn.addVirtualAuthenticator', {
            options: {
                protocol: 'ctap2',
                transport: 'usb',
                hasResidentKey: true,
                hasUserVerification: true,
                isUserVerified: true,
            },
        });
        const authenticatorId = result.authenticatorId;

        const url = new URL(process.env.ZITADEL_API_URL!)
        await client.send('WebAuthn.addCredential', {
            credential: {
                credentialId: "",
                rpId: url.hostname,
                privateKey: BASE64_ENCODED_PK,
                isResidentCredential: false,
                signCount: 0,
            },
            authenticatorId: authenticatorId
        });

        await client.send('WebAuthn.setUserVerified', {
            authenticatorId: authenticatorId,
            isUserVerified: true,
        });
        await client.send('WebAuthn.setAutomaticPresenceSimulation', {
            authenticatorId: authenticatorId,
            enabled: true,
        });

        const user = new PasskeyUser({
            email: "password@example.com",
            firstName: "first",
            lastName: "last",
            organization: "",
        });
        await user.ensure();
        const respJson = await user.ensurePasskeyRegister();

        const credential = await navigator.credentials.create({
            publicKey: respJson.publicKeyCredentialCreationOptions
        });

        await user.ensurePasskeyVerify(respJson.passkeyId, respJson.publicKeyCredentialCreationOptions)
        use(user);
        await client.send('WebAuthn.setAutomaticPresenceSimulation', {
            authenticatorId,
            enabled: false,
        });
    },
});*/

const test = base.extend<{ user: PasskeyUser }>({
    user: async ({page}, use) => {
        const user = new PasskeyUser({
            email: "passkey@example.com",
            firstName: "first",
            lastName: "last",
            organization: "",
        });
        await user.ensure(page);
        await use(user)
    },
});

test("username and passkey login", async ({user, page}) => {
    await user.login(page)
    await page.getByRole("heading", {name: "Welcome " + user.fullName() + "!"}).click();
});
