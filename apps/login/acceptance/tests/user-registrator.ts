import { Page } from "@playwright/test";
import { emailVerify } from "./email-verify.js";
import { passkeyRegister } from "./passkey.js";
import { eventualEmailOTP } from "./mock.js";
import { test as base, UserService } from "./api.js";
import { faker } from "@faker-js/faker";

export class UserRegistrator {

  private readonly passwordField = "password-text-input";
  private readonly passwordConfirmField = "password-confirm-text-input";

  public firstName = faker.person.firstName();
  public lastName = faker.person.lastName();
  public username = faker.internet.email();

  constructor(private page: Page, private svc: UserService) { }

  get fullName() {
    return this.firstName + " " + this.lastName;
  }

  public async registerWithPassword() {
    await this.page.goto("/ui/v2/login/register");
    await this.registerUserScreenPassword();
    await this.page.getByTestId("submit-button").click();
    await this.registerPasswordScreen();
    await this.page.getByTestId("submit-button").click();
    await this.verifyEmail(this.username);
  }

  public async registerWithPasskey(): Promise<string> {
    await this.page.goto("/ui/v2/login/register");
    await this.registerUserScreenPasskey();
    await this.page.getByTestId("submit-button").click();

    // wait for projection of user
    await this.page.waitForTimeout(10000);
    const authId = await passkeyRegister(this.page);

    await this.verifyEmail(this.username);
    return authId;
  }

  private async verifyEmail(email: string) {
    const c = await eventualEmailOTP(email);
    await emailVerify(this.page, c);
  }

  private async registerUserScreenPassword() {
    await this.registerUserScreen();
    await this.page.getByTestId("password-radio").click();
  }

  private async registerUserScreenPasskey() {
    await this.registerUserScreen();
    await this.page.getByTestId("passkey-radio").click();
  }

  private async registerPasswordScreen() {
    await this.page.getByTestId(this.passwordField).pressSequentially("Password2!");
    await this.page.getByTestId(this.passwordConfirmField).pressSequentially("Password2!");
  }

  private async registerUserScreen() {
    await this.page.getByTestId("firstname-text-input").pressSequentially(this.firstName);
    await this.page.getByTestId("lastname-text-input").pressSequentially(this.lastName);
    await this.page.getByTestId("email-text-input").pressSequentially(this.username);
    await this.page.getByTestId("privacy-policy-checkbox").check();
    await this.page.getByTestId("tos-checkbox").check();
  }

  public async cleanup() {
    if (!this.username) {
      console.log("No user to clean up");
      return;
    }
    // TODO: retry until user is found
    await this.page.waitForTimeout(10000);
    const user = await this.svc.getByUsername(this.username);
    await this.svc.native.deleteUser({ userId: user.userId });
  }
}

export const test = base.extend<{ userRegistrator: UserRegistrator }>({
  userRegistrator: async ({ page, userService }, use) => {
    console.log("Setting up user");
    const user = new UserRegistrator(page, userService);
    await use(user);
    await user.cleanup();
  }
});

export { expect } from '@playwright/test';