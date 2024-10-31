import fetch from "node-fetch";
import {Page} from "@playwright/test";
import {registerWithPasskey} from "./register";
import {loginWithPasskey, loginWithPassword} from "./login";
import {changePassword} from "./password";
import {removeUser, getUserByUsername} from './zitadel';

export interface userProps {
    email: string;
    firstName: string;
    lastName: string;
    organization: string;
    password: string;
}

class User {
    private readonly props: userProps;
    private user: string;

    constructor(userProps: userProps) {
        this.props = userProps;
    }

    async ensure(page: Page) {
        await this.remove()

        const body = {
            username: this.props.email,
            organization: {
                orgId: this.props.organization
            },
            profile: {
                givenName: this.props.firstName,
                familyName: this.props.lastName,
            },
            email: {
                email: this.props.email,
                isVerified: true,
            },
            password: {
                password: this.props.password!,
            }
        }

        const response = await fetch(process.env.ZITADEL_API_URL! + "/v2/users/human", {
            method: 'POST',
            body: JSON.stringify(body),
            headers: {
                'Content-Type': 'application/json',
                'Authorization': "Bearer " + process.env.ZITADEL_SERVICE_USER_TOKEN!
            }
        });
        if (response.statusCode >= 400 && response.statusCode != 409) {
            const error = 'HTTP Error: ' + response.statusCode + ' - ' + response.statusMessage;
            console.error(error);
            throw new Error(error);
        }
        return
    }

    async remove() {
        await removeUser(this.userId())
        return
    }

    public setUserId(userId: string) {
        this.user = userId
    }

    public userId() {
        return this.user;
    }

    public username() {
        return this.props.email;
    }

    public password() {
        return this.props.password;
    }

    public firstname() {
        return this.props.firstName
    }

    public lastname() {
        return this.props.lastName
    }

    public fullName() {
        return this.props.firstName + " " + this.props.lastName
    }

    public async login(page: Page) {
        await loginWithPassword(page, this.username(), this.password())
    }

    public async changePassword(page: Page, password: string) {
        await loginWithPassword(page, this.username(), this.password())
        await changePassword(page, this.username(), password)
        this.props.password = password
    }
}

export class PasswordUser extends User {
}

enum OtpType {
    time = "time-based",
    sms = "sms",
    email = "email",
}

export interface otpUserProps {
    email: string;
    firstName: string;
    lastName: string;
    organization: string;
    type: OtpType,
}

export class PasswordUserWithOTP extends User {
    private type: OtpType
    private code: string

    constructor(props: otpUserProps) {
        super({
            email: props.email,
            firstName: props.firstName,
            lastName: props.lastName,
            organization: props.organization,
            password: ""
        })
        this.type = props.type
    }

    async ensure(page: Page) {
        await super.ensure(page)

        const body = {
            username: this.props.email,
            organization: {
                orgId: this.props.organization
            },
            profile: {
                givenName: this.props.firstName,
                familyName: this.props.lastName,
            },
            email: {
                email: this.props.email,
                isVerified: true,
            },
            password: {
                password: this.props.password!,
            }
        }

        const response = await fetch(process.env.ZITADEL_API_URL! + "/v2/users/human", {
            method: 'POST',
            body: JSON.stringify(body),
            headers: {
                'Content-Type': 'application/json',
                'Authorization': "Bearer " + process.env.ZITADEL_SERVICE_USER_TOKEN!
            }
        });
        if (response.statusCode >= 400 && response.statusCode != 409) {
            const error = 'HTTP Error: ' + response.statusCode + ' - ' + response.statusMessage;
            console.error(error);
            throw new Error(error);
        }
        return
    }
}

export interface passkeyUserProps {
    email: string;
    firstName: string;
    lastName: string;
    organization: string;
}

export class PasskeyUser extends User {
    constructor(props: passkeyUserProps) {
        super({
            email: props.email,
            firstName: props.firstName,
            lastName: props.lastName,
            organization: props.organization,
            password: ""
        })
    }

    public async ensure(page: Page) {
        await this.remove()
        await registerWithPasskey(page, this.firstname(), this.lastname(), this.username())
    }

    public async login(page: Page) {
        await loginWithPasskey(page, this.username())
    }

    public async remove() {
        const resp = await getUserByUsername(this.username())
        if (!resp || !resp.result || !resp.result[0]) {
            return
        }
        this.setUserId(resp.result[0].userId)
        await super.remove()
    }
}
