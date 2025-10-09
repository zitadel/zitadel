import { test as base } from "@playwright/test";
import { Transport } from "@connectrpc/connect";
import { UserService } from "./api.js";
import { CreateUserRequestSchema, CreateUserResponseJson, CreateUserResponseSchema } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import minimalRequest from './user-registered-request.json' with { type: "json" };
import { fromJson, toJson } from "@zitadel/client";
import { faker } from "@faker-js/faker";

export class UserCreator {

    public req = {
        ...minimalRequest,
        human: {
            ...minimalRequest.human,
            email: {
                ...minimalRequest.human.email,
                email: faker.internet.email(),
            },
            profile: {
                ...minimalRequest.human.profile,
                givenName: faker.person.firstName(),
                familyName: faker.person.lastName(),
            },
            phone: {
                ...minimalRequest.human.phone,
                phone: faker.phone.number({ style: "international" }),
            },
        }
    }
    private res?: CreateUserResponseJson;

    constructor(private svc: UserService) { }

    async create() {
        const req = { ...this.req, organizationId: await this.svc.orgId() };
        console.log("Creating user", req, null, 2);
        const res = await this.svc.native.createUser(fromJson(CreateUserRequestSchema, req));
        this.res = toJson(CreateUserResponseSchema, res);
        console.log("Created user", JSON.stringify(this.res, null, 2));
        return this.res;
    }

    withPasswordChangeRequired(): UserCreator {
        this.req.human.password.changeRequired = true;
        return this;
    }

    withEmailUnverified(): UserCreator {
        this.req.human.email.isVerified = false;
        return this;
    }

    async addTOTPFactor(): Promise<string> {
        if (!this.res) {
            throw new Error("User must be created before adding a TOTP factor");
        }
        const response = await this.svc.native.registerTOTP({ userId: this.res.id });
        const code = this.svc.generateTOTPToken(response.secret);
        await this.svc.native.verifyTOTPRegistration({ userId: this.res.id, code });
        return response.secret;
    }

    async addEmailOTPFactor(): Promise<void> {
        if (!this.res) {
            throw new Error("User must be created before adding an email OTP factor");
        }
        await this.svc.native.addOTPEmail({ userId: this.res.id });
    }

    async addSMSOTPFactor(): Promise<void> {
        if (!this.res) {
            throw new Error("User must be created before adding an SMS OTP factor");
        }
        await this.svc.native.addOTPSMS({ userId: this.res.id });
    }

    async cleanup() {
        if (this.res) {
            await this.svc.native.deleteUser({ userId: this.res.id });
        }
    }

    get username(): string {
        return this.req.human.email.email
    }

    get password(): string {
        return this.req.human.password.password;
    }

    get phone(): string {
        return this.req.human.phone.phone;
    }

    get fullName(): string {
        return `${this.req.human.profile.givenName} ${this.req.human.profile.familyName}`;
    }    
}

export const test = base.extend<{ transport: Transport, userService: UserService, userCreator: UserCreator }>({
    userCreator: async ({ userService }, use) => {
        const user = new UserCreator(userService);
        await use(user);
        await user.cleanup();
    }
});

export { expect } from '@playwright/test';

