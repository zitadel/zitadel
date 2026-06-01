import { Page } from "@playwright/test";
import { registerWithPasskey } from "./register";
import { activateOTP, addTOTP, addUser, eventualNewUser, getUserByUsername, removeUser } from "./zitadel";

export interface userProps {
  email: string;
  isEmailVerified?: boolean;
  firstName: string;
  lastName: string;
  organization: string;
  password: string;
  passwordChangeRequired?: boolean;
  phone: string;
  isPhoneVerified?: boolean;
}

class User {
  private readonly props: userProps;
  private user: string;

  constructor(userProps: userProps) {
    this.props = userProps;
  }

  async ensure(page: Page) {
    const response = await addUser(this.props);

    this.setUserId(response.userId);
  }

  async cleanup() {
    await removeUser(this.getUserId());
  }

  public setUserId(userId: string) {
    this.user = userId;
  }

  public getUserId() {
    return this.user;
  }

  public getUsername() {
    return this.props.email;
  }

  public getPassword() {
    return this.props.password;
  }

  public getFirstname() {
    return this.props.firstName;
  }

  public getLastname() {
    return this.props.lastName;
  }

  public getPhone() {
    return this.props.phone;
  }

  public getFullName() {
    return `${this.props.firstName} ${this.props.lastName}`;
  }
}

export class PasswordUser extends User {
  async ensure(page: Page) {
    await super.ensure(page);
    await eventualNewUser(this.getUserId());
  }
}

export enum OtpType {
  sms = "sms",
  email = "email",
}

export interface otpUserProps {
  email: string;
  isEmailVerified?: boolean;
  firstName: string;
  lastName: string;
  organization: string;
  password: string;
  passwordChangeRequired?: boolean;
  phone: string;
  isPhoneVerified?: boolean;
  type: OtpType;
}

export class PasswordUserWithOTP extends User {
  private type: OtpType;

  constructor(props: otpUserProps) {
    super({
      email: props.email,
      firstName: props.firstName,
      lastName: props.lastName,
      organization: props.organization,
      password: props.password,
      phone: props.phone,
      isEmailVerified: props.isEmailVerified,
      isPhoneVerified: props.isPhoneVerified,
      passwordChangeRequired: props.passwordChangeRequired,
    });
    this.type = props.type;
  }

  async ensure(page: Page) {
    await super.ensure(page);
    await activateOTP(this.getUserId(), this.type);
    await eventualNewUser(this.getUserId());
  }
}

export class PasswordUserWithTOTP extends User {
  private secret: string;

  async ensure(page: Page) {
    await super.ensure(page);
    this.secret = await addTOTP(this.getUserId());
    await eventualNewUser(this.getUserId());
  }

  public getSecret(): string {
    return this.secret;
  }
}

export interface passkeyUserProps {
  email: string;
  firstName: string;
  lastName: string;
  organization: string;
  phone: string;
  isEmailVerified?: boolean;
  isPhoneVerified?: boolean;
}

export class PasskeyUser extends User {
  private authenticatorId: string;

  constructor(props: passkeyUserProps) {
    super({
      email: props.email,
      firstName: props.firstName,
      lastName: props.lastName,
      organization: props.organization,
      password: "",
      phone: props.phone,
      isEmailVerified: props.isEmailVerified,
      isPhoneVerified: props.isPhoneVerified,
    });
  }

  public async ensure(page: Page) {
    const authId = await registerWithPasskey(page, this.getFirstname(), this.getLastname(), this.getUsername());
    this.authenticatorId = authId;

    // wait for projection of user
    await page.waitForTimeout(10000);
  }

  async cleanup() {
    const resp: any = await getUserByUsername(this.getUsername());
    if (!resp || !resp.result || !resp.result[0]) {
      return;
    }
    await removeUser(resp.result[0].userId);
  }

  public getAuthenticatorId(): string {
    return this.authenticatorId;
  }
}
