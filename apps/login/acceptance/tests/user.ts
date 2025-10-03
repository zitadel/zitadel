import { CreateUserRequest, CreateUserResponse, UserService as NativeUserService } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { faker } from "@faker-js/faker";
import { test as base } from "@playwright/test";
import { readFileSync } from "fs";
import { createServerTransport } from "@zitadel/client/node";
import { Transport, Client } from "@connectrpc/connect";
import { createClientFor } from "@zitadel/client";
import { Authenticator } from "@otplib/core";
import { createDigest, createRandomBytes } from "@otplib/plugin-crypto";
import { keyDecoder, keyEncoder } from "@otplib/plugin-thirty-two"; // use your chosen base32 plugin

class UserService {
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

    private totp(secret: string) {
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

class User {
    constructor(private svc: UserService) { }

    public readonly default: CreateUserRequest = {
        $typeName: "zitadel.user.v2.CreateUserRequest",
        organizationId: "340565276842066283",
        userType: {
            case: "human",
            value: {
                $typeName: "zitadel.user.v2.CreateUserRequest.Human",
                metadata: [],
                idpLinks: [],
                email: {
                    $typeName: "zitadel.user.v2.SetHumanEmail",
                    email: faker.internet.email(),
                    verification: {
                        case: "isVerified",
                        value: true
                    }
                },
                profile: {
                    $typeName: "zitadel.user.v2.SetHumanProfile",
                    givenName: faker.person.firstName(),
                    familyName: faker.person.lastName(),
                },
                phone: {
                    $typeName: "zitadel.user.v2.SetHumanPhone",
                    phone: faker.phone.number(),
                    verification: {
                        case: "isVerified",
                        value: true
                    }
                },
                passwordType: {
                    case: "password",
                    value: {
                        $typeName: "zitadel.user.v2.Password",
                        password: "Password1!",
                        changeRequired: false,
                    }
                }
            },
        }
    };
    public res: CreateUserResponse | null = null;
    public req: CreateUserRequest = { ...this.default };

    async create(req: CreateUserRequest = this.default) {
        this.req = req;
        this.res = await this.svc.native.createUser(req);
    }

    async cleanup() {
        if (this.res) {
            await this.svc.native.deleteUser({ userId: this.res.id });
        }
    }

    get username(): string {
        return this.req.username!;
    }

    get password(): string {
        if (this.req.userType?.case !== "human" || this.req.userType.value.passwordType.case !== "password") {
            throw new Error("User has no password in the request.");
        }
        return this.req.userType.value.passwordType.value.password;
    }

    get fullName(): string {
        if (this.req.userType?.case !== "human" || !this.req.userType.value.profile) {
            throw new Error("User has no profile in the request.");
        }
        return `${this.req.userType.value.profile.givenName} ${this.req.userType.value.profile.familyName}`;
    }
}

export const test = base.extend<{ user: User; transport: Transport, userService: UserService }>({
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
    },
    user: async ({ userService }, use) => {
        console.log("Setting up user");
        const user = new User(userService);
        await use(user);
        await user.cleanup();
    }
});

export { expect } from '@playwright/test';

