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
