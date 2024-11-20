import { Page } from "@playwright/test";
import axios from "axios";
import { registerWithPasskey } from "./register";
import { getUserByUsername, removeUser } from "./zitadel";

export interface userProps {
  email: string;
  firstName: string;
  lastName: string;
  organization: string;
  password: string;
  phone: string;
}

class User {
  private readonly props: userProps;
  private user: string;

  constructor(userProps: userProps) {
    this.props = userProps;
  }

  async ensure(page: Page) {
    await this.remove();

    const body = {
      username: this.props.email,
      organization: {
        orgId: this.props.organization,
      },
      profile: {
        givenName: this.props.firstName,
        familyName: this.props.lastName,
      },
      email: {
        email: this.props.email,
        isVerified: true,
      },
      phone: {
        phone: this.props.phone!,
        isVerified: true,
      },
      password: {
        password: this.props.password!,
      },
    };

    try {
      const response = await axios.post(`${process.env.ZITADEL_API_URL}/v2/users/human`, body, {
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${process.env.ZITADEL_SERVICE_USER_TOKEN}`,
        },
      });

      if (response.status >= 400 && response.status !== 409) {
        const error = `HTTP Error: ${response.status} - ${response.statusText}`;
        console.error(error);
        throw new Error(error);
      }
      this.setUserId(response.data.userId);
    } catch (error) {
      console.error("Error making request:", error);
      throw error;
    }

    // wait for projection of user
    await page.waitForTimeout(3000);
  }

  async remove() {
    const resp: any = await getUserByUsername(this.getUsername());
    if (!resp || !resp.result || !resp.result[0]) {
      return;
    }
    await removeUser(resp.result[0].userId);
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

export class PasswordUser extends User {}

export enum OtpType {
  sms = "sms",
  email = "email",
}

export interface otpUserProps {
  email: string;
  firstName: string;
  lastName: string;
  organization: string;
  password: string;
  phone: string;
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
    });
    this.type = props.type;
  }

  async ensure(page: Page) {
    await super.ensure(page);

    let url = "otp_";
    switch (this.type) {
      case OtpType.sms:
        url = url + "sms";
        break;
      case OtpType.email:
        url = url + "email";
        break;
    }

    try {
      const response = await axios.post(
        `${process.env.ZITADEL_API_URL}/v2/users/${this.getUserId()}/${url}`,
        {},
        {
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${process.env.ZITADEL_SERVICE_USER_TOKEN}`,
          },
        },
      );

      if (response.status >= 400 && response.status !== 409) {
        const error = `HTTP Error: ${response.status} - ${response.statusText}`;
        console.error(error);
        throw new Error(error);
      }
    } catch (error) {
      console.error("Error making request:", error);
      throw error;
    }

    // wait for projection of user
    await page.waitForTimeout(2000);
  }
}

export interface passkeyUserProps {
  email: string;
  firstName: string;
  lastName: string;
  organization: string;
  phone: string;
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
    });
  }

  public async ensure(page: Page) {
    await this.remove();
    const authId = await registerWithPasskey(page, this.getFirstname(), this.getLastname(), this.getUsername());
    this.authenticatorId = authId;

    // wait for projection of user
    await page.waitForTimeout(2000);
  }

  public async remove() {
    await super.remove();
  }

  public getAuthenticatorId(): string {
    return this.authenticatorId;
  }
}
