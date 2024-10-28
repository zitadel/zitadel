import fetch from 'node-fetch';
import {Page} from "@playwright/test";
import {registerWithPasskey} from "./register";
import {loginWithPasskey, loginWithPassword} from "./login";
import {changePassword} from "./password";

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

    async ensure() {
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
        const response = await fetch(process.env.ZITADEL_API_URL! + "/v2/users/" + this.userId(), {
            method: 'DELETE',
            headers: {
                'Authorization': "Bearer " + process.env.ZITADEL_SERVICE_USER_TOKEN!
            }
        });
        if (response.statusCode >= 400 && response.statusCode != 404) {
            const error = 'HTTP Error: ' + response.statusCode + ' - ' + response.statusMessage;
            console.error(error);
            throw new Error(error);
        }
        return
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

export interface passkeyUserProps {
    email: string;
    firstName: string;
    lastName: string;
    organization: string;
}

export class PasskeyUser {
    private props: passkeyUserProps

    constructor(props: passkeyUserProps) {
        this.props = props
    }

    async ensurePasskey(page: Page) {
        await registerWithPasskey(page, this.props.firstName, this.props.lastName, this.props.email)
    }

    public async login(page: Page) {
        await loginWithPasskey(page, this.props.email)
    }

    public fullName() {
        return this.props.firstName + " " + this.props.lastName
    }

    async ensurePasskeyRegister() {
        const url = new URL(process.env.ZITADEL_API_URL!)
        const registerBody = {
            domain: url.hostname,
        }
        const userId = ""
        const registerResponse = await fetch(process.env.ZITADEL_API_URL! + "/v2/users/" + userId + "/passkeys", {
            method: 'POST',
            body: JSON.stringify(registerBody),
            headers: {
                'Content-Type': 'application/json',
                'Authorization': "Bearer " + process.env.ZITADEL_SERVICE_USER_TOKEN!
            }
        });
        if (registerResponse.statusCode >= 400 && registerResponse.statusCode != 409) {
            const error = 'HTTP Error: ' + registerResponse.statusCode + ' - ' + registerResponse.statusMessage;
            console.error(error);
            throw new Error(error);
        }
        const respJson = await registerResponse.json()
        return respJson
    }

    async ensurePasskeyVerify(passkeyId: string, credential: Credential) {
        const verifyBody = {
            publicKeyCredential: credential,
            passkeyName: "passkey",
        }
        const userId = ""
        const verifyResponse = await fetch(process.env.ZITADEL_API_URL! + "/v2/users/" + userId + "/passkeys/" + passkeyId, {
            method: 'POST',
            body: JSON.stringify(verifyBody),
            headers: {
                'Content-Type': 'application/json',
                'Authorization': "Bearer " + process.env.ZITADEL_SERVICE_USER_TOKEN!
            }
        });
        if (verifyResponse.statusCode >= 400 && verifyResponse.statusCode != 409) {
            const error = 'HTTP Error: ' + verifyResponse.statusCode + ' - ' + verifyResponse.statusMessage;
            console.error(error);
            throw new Error(error);
        }
        return
    }
}