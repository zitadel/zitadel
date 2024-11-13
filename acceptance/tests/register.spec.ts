import {test} from "@playwright/test";
import {registerWithPasskey, registerWithPassword} from './register';
import {loginScreenExpect} from "./login";
import {removeUserByUsername} from './zitadel';
import path from 'path';
import dotenv from 'dotenv';

// Read from ".env" file.
dotenv.config({path: path.resolve(__dirname, '.env.local')});

test("register with password", async ({page}) => {
    const username = "register-password@example.com"
    const password = "Password1!"
    const firstname = "firstname"
    const lastname = "lastname"

    await removeUserByUsername(username)
    await registerWithPassword(page, firstname, lastname, username, password, password)
    await loginScreenExpect(page, firstname + " " + lastname);
});

test("register with passkey", async ({page}) => {
    const username = "register-passkey@example.com"
    const firstname = "firstname"
    const lastname = "lastname"

    await removeUserByUsername(username)
    await registerWithPasskey(page, firstname, lastname, username)
    await loginScreenExpect(page, firstname + " " + lastname);
});
