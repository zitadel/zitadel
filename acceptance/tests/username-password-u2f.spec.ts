import {test as base} from "@playwright/test";
import {OtpType, PasswordUserWithOTP} from './user';
import path from 'path';
import dotenv from 'dotenv';
import {loginScreenExpect, loginWithPassword} from "./login";
import {startSink} from "./otp";

// Read from ".env" file.
dotenv.config({path: path.resolve(__dirname, '.env.local')});

const test = base.extend<{ user: PasswordUserWithOTP }>({
    user: async ({page}, use) => {
        const user = new PasswordUserWithOTP({
            email: "otp_sms@example.com",
            firstName: "first",
            lastName: "last",
            password: "Password1!",
            organization: "",
            type: OtpType.sms,
        });

        await user.ensure(page);
        await use(user);
    },
});


test("username, password and u2f login", async ({user, page}) => {
    // Given u2f is enabled on the organizaiton of the user
    // Given the user has only u2f configured as second factor

    // User enters username
    // User enters password
    // Popup for u2f is directly opened
    // User verifies u2f
    // User is redirected to the app (default redirect url)
});


test("username, password and u2f login, multiple mfa options", async ({user, page}) => {
    // Given u2f and semailms otp is enabled on the organizaiton of the user
    // Given the user has u2f and email otp configured as second factor

    // User enters username
    // User enters password
    // Popup for u2f is directly opened
    // User aborts u2f verification
    // User clicks button to use email otp as second factor
    // User receives an email with a verification code
    // User enters code in ui
    // User is redirected to the app (default redirect url)
});
