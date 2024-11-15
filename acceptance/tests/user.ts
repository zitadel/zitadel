import fetch from "node-fetch";
import {Page} from "@playwright/test";
import {registerWithPasskey} from "./register";
import {getUserByUsername, removeUser} from './zitadel';

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
        const resp = await getUserByUsername(this.getUsername())
        if (!resp || !resp.result || !resp.result[0]) {
            return
        }
        await removeUser(resp.result[0].userId)
        return
    }

    public setUserId(userId: string) {
        this.user = userId
    }

    public getUserId() {
        return this.user;
    }

    public getUsername() {
        return this.props.email;
    }

    public getPassword() {
        return this.props.password;
    }

    public getFirstname() {
        return this.props.firstName
    }

    public getLastname() {
        return this.props.lastName
    }

    public getFullName() {
        return this.props.firstName + " " + this.props.lastName
    }
}

export class PasswordUser extends User {
}

export enum OtpType {
    sms = "sms",
    email = "email",
}

export interface otpUserProps {
    email: string;
    firstName: string;
    lastName: string;
    organization: string;
    password: string,
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
            password: props.password,
        })
        this.type = props.type
    }

    async ensure(page: Page) {
        await super.ensure(page)

        let url = "otp_"
        switch (this.type) {
            case OtpType.sms:
                url = url + "sms"
            case OtpType.email:
                url = url + "email"
        }

        const response = await fetch(process.env.ZITADEL_API_URL! + "/v2/users/" + this.getUserId() + "/" + url, {
            method: 'POST',
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

        // TODO: get code from SMS or Email provider
        this.code = ""
        return
    }

    public getCode() {
        return this.code
    }
}

export interface passkeyUserProps {
    email: string;
    firstName: string;
    lastName: string;
    organization: string;
}

export class PasskeyUser extends User {
    private authenticatorId: string

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
        const authId = await registerWithPasskey(page, this.getFirstname(), this.getLastname(), this.getUsername())
        this.authenticatorId = authId
    }

    public async remove() {
        await super.remove()
    }

    public getAuthenticatorId(): string {
        return this.authenticatorId
    }
}
