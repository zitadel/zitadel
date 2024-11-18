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
            type: OtpType.email,
        });

        await user.ensure(page);
        await use(user);
    },
});


test("username, password and email otp login, enter code manually", async ({user, page}) => {
    // Given email otp is enabled on the organizaiton of the user
    // Given the user has only email otp configured as second factor

    // User enters username
    // User enters password
    // User receives an email with a verification code
    // User enters the code into the ui
    // User is redirected to the app
});


test("username, password and email otp login, click link in email", async ({user, page}) => {
    // Given email otp is enabled on the organizaiton of the user
    // Given the user has only email otp configured as second factor

    // User enters username
    // User enters password
    // User receives an email with a verification code
    // User clicks link in the email
    // User is redirected to the app
});

test("username, password and email otp login, resend code", async ({user, page}) => {
    // Given email otp is enabled on the organizaiton of the user
    // Given the user has only email otp configured as second factor

    // User enters username
    // User enters password
    // User receives an email with a verification code
    // User clicks resend code
    // User receives a new email with a verification code
    // User enters the new code in the ui
    // User is redirected to the app
});

test("username, password and email otp login, wrong code", async ({user, page}) => {
    // Given email otp is enabled on the organizaiton of the user
    // Given the user has only email otp configured as second factor

    // User enters username
    // User enters password
    // User receives an email with a verification code
    // User enters a wrond code
    // Error message - "Invalid code" is shown
});

test("username, password and email otp login, multiple mfa options", async ({user, page}) => {
    // Given email otp and sms otp is enabled on the organizaiton of the user
    // Given the user has email and sms otp configured as second factor

    // User enters username
    // User enters password
    // User receives an email with a verification code
    // User clicks button to use sms otp as second factor
    // User receives an sms with a verification code
    // User enters code in ui
    // User is redirected to the app
});
