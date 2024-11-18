import {test as base} from "@playwright/test";
import path from 'path';
import dotenv from 'dotenv';
import {PasskeyUser} from "./user";
import {loginScreenExpect, loginWithPasskey} from "./login";

// Read from ".env" file.
dotenv.config({path: path.resolve(__dirname, '.env.local')});

const test = base.extend<{ user: PasskeyUser }>({
    user: async ({page}, use) => {
        const user = new PasskeyUser({
            email: "passkey@example.com",
            firstName: "first",
            lastName: "last",
            organization: "",
        });
        await user.ensure(page);
        await use(user);
    },
});

test("username and passkey login", async ({user, page}) => {
    await loginWithPasskey(page, user.getAuthenticatorId(), user.getUsername())
    await loginScreenExpect(page, user.getFullName());
});

test("username and passkey login, if passkey enabled", async ({user, page}) => {
    // Given passkey is enabled on the organization of the user
    // Given the user has only passkey enabled as authentication

    // enter username
    // passkey popup is directly shown
    // user verifies passkey
    // user is redirected to app
});

test("username and passkey login, multiple auth methods", async ({user, page}) => {
    // Given passkey and password is enabled on the organization of the user
    // Given the user has password and passkey registered

    // enter username
    // passkey popup is directly shown
    // user aborts passkey authentication
    // user switches to password authentication
    // user enters password
    // user is redirected to app
});
