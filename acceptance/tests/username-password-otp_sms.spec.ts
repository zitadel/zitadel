import { test as base } from "@playwright/test";
import dotenv from "dotenv";
import path from "path";
import { loginScreenExpect, loginWithPassword } from "./login";
import { OtpType, PasswordUserWithOTP } from "./user";

// Read from ".env" file.
dotenv.config({ path: path.resolve(__dirname, ".env.local") });

const test = base.extend<{ user: PasswordUserWithOTP }>({
  user: async ({ page }, use) => {
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

/*
test("username, password and otp login", async ({ user, page }) => {
  //const server = startSink()
  await loginWithPassword(page, user.getUsername(), user.getPassword());

  await loginScreenExpect(page, user.getFullName());
  //server.close()
});
*/