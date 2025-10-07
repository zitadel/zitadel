import { test as base } from "@playwright/test";
import { readFileSync } from "fs";
import { createServerTransport } from "@zitadel/client/node";
import { Transport, Client } from "@connectrpc/connect";
import { createClientFor, fromJson } from "@zitadel/client";
import { ListUsersRequestSchema, UserService as NativeUserService } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { Authenticator } from "@otplib/core";
import { createDigest, createRandomBytes } from "@otplib/plugin-crypto";
import { keyDecoder, keyEncoder } from "@otplib/plugin-thirty-two";

export const test = base.extend<{}, { transport: Transport, userService: UserService, orgId: string }>({
    transport: [async ({ }, use) => {
        console.log("Setting up transport");
        const adminToken = readFileSync(process.env.ZITADEL_ADMIN_TOKEN_FILE!).toString().trim()
        const transport = createServerTransport(adminToken, { baseUrl: process.env.ZITADEL_API_URL! });
        await use(transport);
    }, { scope: 'worker', auto: true }],
    userService: [async ({ transport }, use) => {
        console.log("Setting up user service");
        const nativeUserService = createClientFor(NativeUserService)(transport);
        const svc = new UserService(nativeUserService);
        await use(svc);
    }, { scope: 'worker', auto: true }]
});

export class UserService {

    private orgIdCache?: string;

    constructor(public readonly native: Client<typeof NativeUserService>) {}

    async getByUsername(username: string) {
        const res = await this.native.listUsers(fromJson(ListUsersRequestSchema, {
            query: {
                limit: 1,
            },
            queries: [{
                userNameQuery: {
                    userName: username,
                }
            }]
        }));
        if (res.result?.length !== 1) {
            throw new Error(`User with username ${username} not found`);
        }
        return res.result[0];
    }

    public async orgId(): Promise<string> {
        if (this.orgIdCache) {
            return this.orgIdCache;
        }
        const adminUser = await this.getByUsername(process.env.ZITADEL_ADMIN_USER!);
        this.orgIdCache = adminUser.details?.resourceOwner;
        if (!this.orgIdCache) {
            throw new Error(`Admin user ${process.env.ZITADEL_ADMIN_USER} has no orgId`);
        }
        return this.orgIdCache;
    }

    public generateTOTPToken(secret: string) {
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

