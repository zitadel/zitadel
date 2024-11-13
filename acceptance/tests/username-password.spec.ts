import {test as base} from "@playwright/test";
import {PasswordUser} from './user';
import path from 'path';
import dotenv from 'dotenv';
import {loginScreenExpect, loginWithPassword, startLogin} from "./login";
import {loginnameScreenExpect} from "./loginname-screen";
import {passwordScreenExpect} from "./password-screen";
import {loginname} from "./loginname";
import {password} from "./password";

// Read from ".env" file.
dotenv.config({path: path.resolve(__dirname, '.env.local')});

const test = base.extend<{ user: PasswordUser }>({
    user: async ({page}, use) => {
        const user = new PasswordUser({
            email: "password@example.com",
            firstName: "first",
            lastName: "last",
            password: "Password1!",
            organization: "",
        });
        await user.ensure(page);
        await use(user);
    },
});

test("username and password login", async ({user, page}) => {
    await loginWithPassword(page, user.getUsername(), user.getPassword())
    await loginScreenExpect(page, user.getFullName());
});

test("username and password login, unknown username", async ({page}) => {
    const username = "unknown"
    await startLogin(page);
    await loginname(page, username)
    await loginnameScreenExpect(page, username)
});

test("username and password login, wrong password", async ({user, page}) => {
    await startLogin(page);
    await loginname(page, user.getUsername())
    await password(page, "wrong")
    await passwordScreenExpect(page, "wrong")
});
