import { test as base } from "@playwright/test";
import { Transport } from "@connectrpc/connect";
import { UserService } from "./api.js";
import { CreateUserRequest, CreateUserResponse, UserService as NativeUserService } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { faker } from "@faker-js/faker";

export class RegisteredUser {
    constructor(private svc: UserService) { }

    public readonly minimal: CreateUserRequest = {
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
                    },
                },
            },
        }
    };
    public res: CreateUserResponse | null = null;
    public req: CreateUserRequest = { ...this.minimal };

    async create(req: CreateUserRequest = this.minimal) {
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

export const test = base.extend<{ transport: Transport, userService: UserService, registeredUser: RegisteredUser }>({
    registeredUser: async ({ userService }, use) => {
        console.log("Setting up user");
        const user = new RegisteredUser(userService);
        await use(user);
        await user.cleanup();
    }
});

export { expect } from '@playwright/test';

