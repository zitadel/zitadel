import { test as base } from "@playwright/test";
import { Transport } from "@connectrpc/connect";
import { UserService } from "./api.js";
import { CreateUserRequest,  CreateUserRequestSchema, CreateUserResponse, UserService as NativeUserService } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import minimalRequest from './user-registered-request.json' with { type: "json" };
import { create } from "@zitadel/client";
import { faker } from "@faker-js/faker";

export class CreateUserRequestBuilder {

    public req = minimalRequest
    constructor() { }

    build(): CreateUserRequest {
        return create(CreateUserRequestSchema, {
            ...minimalRequest,
            ... {
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
        })
    }

    withPasswordChangeRequired(): CreateUserRequestBuilder {
        this.req.human.password.changeRequired = true;
        return this;
   }
}


export class RegisteredUser {
    constructor(private svc: UserService) { }

    public res: CreateUserResponse | null = null;
    public builder: CreateUserRequestBuilder = new CreateUserRequestBuilder();

    async create(req?: CreateUserRequestBuilder) {
        if (req) {
            this.builder = req;
        }
        this.res = await this.svc.native.createUser(this.builder.build());
        console.log("Created user", this.res);
        return this.res;
    }

    async cleanup() {
        if (this.res) {
            await this.svc.native.deleteUser({ userId: this.res.id });
        }
    }

    get username(): string {
        return this.builder.req.human.email.email
    }

    get password(): string {
        return this.builder.req.human.password.password;
    }

    get phone(): string {
        return this.builder.req.human.phone.phone;
    }

    get fullName(): string {
        return `${this.builder.req.human.profile.givenName} ${this.builder.req.human.profile.familyName}`;
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

