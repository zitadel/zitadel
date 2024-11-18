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

test("username, password and otp login", async ({user, page}) => {
    const server = startSink()
    await loginWithPassword(page, user.getUsername(), user.getPassword())


    await loginScreenExpect(page, user.getFullName());
    server.close()
});


test("username, password and sms otp login", async ({user, page}) => {
    // Given sms otp is enabled on the organizaiton of the user
    // Given the user has only sms otp configured as second factor

    // User enters username
    // User enters password
    // User receives an sms with a verification code
    // User enters the code into the ui
    // User is redirected to the app (default redirect url)
});


test("username, password and sms otp login, resend code", async ({user, page}) => {
    // Given sms otp is enabled on the organizaiton of the user
    // Given the user has only sms otp configured as second factor

    // User enters username
    // User enters password
    // User receives an sms with a verification code
    // User clicks resend code
    // User receives a new sms with a verification code
    // User is redirected to the app (default redirect url)
});


test("username, password and sms otp login, wrong code", async ({user, page}) => {
    // Given sms otp is enabled on the organizaiton of the user
    // Given the user has only sms otp configured as second factor

    // User enters username
    // User enters password
    // User receives an sms with a verification code
    // User enters a wrond code
    // Error message - "Invalid code" is shown
});
