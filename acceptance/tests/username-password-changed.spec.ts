import {test as base} from "@playwright/test";
import {PasswordUser} from './user';
import path from 'path';
import dotenv from 'dotenv';

// Read from ".env" file.
dotenv.config({path: path.resolve(__dirname, '.env.local')});

const test = base.extend<{ user: PasswordUser }>({
    user: async ({page}, use) => {
        const user = new PasswordUser({
            email: "password-changed@example.com",
            firstName: "first",
            lastName: "last",
            password: "Password1!",
            organization: "",
        });
        await user.ensure();
        await use(user);
    },
});

test("username and password changed login", async ({user, page}) => {
    await user.changePassword(page, "ChangedPw1!")
    await user.login(page)
});
