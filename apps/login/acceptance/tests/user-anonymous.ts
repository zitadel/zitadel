import { Page } from "@playwright/test";
import { emailVerify } from "./email-verify.js";
import { passkeyRegister } from "./passkey.js";
import { eventualEmailOTP } from "./mock.js";
import { test as base, UserService } from "./api.js";
import { faker } from "@faker-js/faker";

export class AnonymousUser {

  private readonly passwordField = "password-text-input";
  private readonly passwordConfirmField = "password-confirm-text-input";
  public username: string | null = null;

  defaultFirstName = faker.person.firstName();
  defaultLastName = faker.person.lastName()
  defaultEmail = faker.internet.email();
  

  constructor(private page: Page, private svc: UserService) { }

  public async registerWithPassword(
    firstname: string = this.defaultFirstName,
    lastname: string = this.defaultLastName,
    email: string = this.defaultEmail,
    password1: string,
    password2: string,
  ) {
    await this.page.goto("/ui/v2/login/register");
    await this.registerUserScreenPassword(firstname, lastname, email);
    await this.page.getByTestId("submit-button").click();
    this.username = email;
    await this.registerPasswordScreen(password1, password2);
    await this.page.getByTestId("submit-button").click();
    await this.verifyEmail(email);
  }

  public async registerWithPasskey(
    firstname: string = this.defaultFirstName,
    lastname: string = this.defaultLastName,
     email: string = this.defaultEmail
  ): Promise<string> {
    await this.page.goto("/ui/v2/login/register");
    await this.registerUserScreenPasskey(firstname, lastname, email);
    await this.page.getByTestId("submit-button").click();
    this.username = email;

    // wait for projection of user
    await this.page.waitForTimeout(10000);
    const authId = await passkeyRegister(this.page);

    await this.verifyEmail(email);
    return authId;
  }

  private async verifyEmail(email: string) {
    const c = await eventualEmailOTP(email);
    await emailVerify(this.page, c);
  }

  async registerUserScreenPassword(firstname: string, lastname: string, email: string) {
    await this.registerUserScreen(firstname, lastname, email);
    await this.page.getByTestId("password-radio").click();
  }

  async registerUserScreenPasskey(firstname: string, lastname: string, email: string) {
    await this.registerUserScreen(firstname, lastname, email);
    await this.page.getByTestId("passkey-radio").click();
  }

  async registerPasswordScreen(password1: string, password2: string) {
    await this.page.getByTestId(this.passwordField).pressSequentially(password1);
    await this.page.getByTestId(this.passwordConfirmField).pressSequentially(password2);
  }

  async registerUserScreen(firstname: string, lastname: string, email: string) {
    await this.page.getByTestId("firstname-text-input").pressSequentially(firstname);
    await this.page.getByTestId("lastname-text-input").pressSequentially(lastname);
    await this.page.getByTestId("email-text-input").pressSequentially(email);
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

export const test = base.extend<{ anonymousUser: AnonymousUser }>({
  anonymousUser: async ({ page, userService }, use) => {
    console.log("Setting up user");
    const user = new AnonymousUser(page, userService);
    await use(user);
    await user.cleanup();
  }
});

export { expect } from '@playwright/test';