import { test as base } from "@playwright/test";
import { readFileSync } from "fs";
import { createServerTransport } from "@zitadel/client/node";
import { Transport, Client } from "@connectrpc/connect";
import { createClientFor } from "@zitadel/client";
import { UserService as NativeUserService } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { Authenticator } from "@otplib/core";
import { createDigest, createRandomBytes } from "@otplib/plugin-crypto";
import { keyDecoder, keyEncoder } from "@otplib/plugin-thirty-two"; // use your chosen base32 plugin

export const test = base.extend<{ transport: Transport, userService: UserService }>({
    transport: async ({ }, use) => {
        console.log("Setting up transport");
        const adminToken = readFileSync(process.env.ZITADEL_ADMIN_TOKEN_FILE!).toString().trim()
        const transport = createServerTransport(adminToken, { baseUrl: process.env.ZITADEL_API_URL! });
        await use(transport);
    },
    userService: async ({ transport }, use) => {
        console.log("Setting up user service");
        const nativeUserService = createClientFor(NativeUserService)(transport);
        const svc = new UserService(nativeUserService);
        await use(svc);
    }
});

export class UserService {
    constructor(public readonly native: Client<typeof NativeUserService>) { }

    async getByUsername(username: string) {
        const res = await this.native.listUsers({
            query: {
                limit: 1,
            },
            queries: [{
                query: {
                    case: "userNameQuery",
                    value: {
                        userName: username,
                    }
                }
            }]
        })
        if (res.result?.length !== 1) {
            throw new Error(`User with username ${username} not found`);
        }
        return res.result[0];
    }

    async addTOTP(userId: string): Promise<string> {
        const response = await this.native.registerTOTP({ userId });
        const code = this.totp(response.secret);
        await this.native.verifyTOTPRegistration({ userId, code });
        return response.secret;
    }

    public totp(secret: string) {
        const authenticator = new Authenticator({
            createDigest,
            createRandomBytes,
            keyDecoder,
            keyEncoder,
        });
        // google authenticator usage
        const token = authenticator.generate(secret);

        // check if token can be used
        if (!authenticator.verify({ token: token, secret: secret })) {
            const error = `Generated token could not be verified`;
            console.error(error);
            throw new Error(error);
        }

        return token;
    }
}

export { expect } from '@playwright/test';

