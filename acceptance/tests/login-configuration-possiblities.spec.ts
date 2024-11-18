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

test("login with mfa setup, mfa setup prompt", async ({user, page}) => {
    // Given the organization has set "multifactor init check time" to 40
    // Given the organization has enabled all possible mfa types 
    // Given the user has a password but no mfa registered and never authenticated

    // enter login name
    // enter password
    // User is prompted to setup a mfa, all possible mfa providers are listed, the user can choose the provider
});

test("login with mfa setup, mfa setup prompt", async ({user, page}) => {
    // Given the organization has set "multifactor init check time" to 0
    // Given the organization has enabled all possible mfa types 
    // Given the user has a password but no mfa registered and never authenticated

    // enter login name
    // enter password
    // user is redirected to app
});

test("login with mfa setup, force mfa for local authenticated users", async ({user, page}) => {
    // Given the organization has enabled force mfa for local authentiacted users
    // Given the organization has enabled all possible mfa types 
    // Given the user has a password but no mfa registered

    // enter login name
    // enter password
    // User is prompted to setup a mfa, all possible mfa providers are listed, the user can choose the provider
});


test("login with mfa setup, force mfa - local user", async ({user, page}) => {
    // Given the organization has enabled force mfa for local authentiacted users
    // Given the organization has enabled all possible mfa types 
    // Given the user has a password but no mfa registered

    // enter login name
    // enter password
    // User is prompted to setup a mfa, all possible mfa providers are listed, the user can choose the provider
});


test("login with mfa setup, force mfa - external user", async ({user, page}) => {
    // Given the organization has enabled force mfa
    // Given the organization has enabled all possible mfa types 
    // Given the user has an idp but no mfa registered

    // enter login name
    // redirect to configured external idp
    // User is prompted to setup a mfa, all possible mfa providers are listed, the user can choose the provider
});


test("login with mfa setup, force mfa - external user", async ({user, page}) => {
    // Given the organization has a password lockout policy set to 1 on the max password attempts
    // Given the user has only a password as auth methos

    // enter login name
    // enter wrong password
    // User will get an error "Wrong password"
    // enter password
    // User will get an error "Max password attempts reached - user is locked. Please reach out to your administrator"
});

