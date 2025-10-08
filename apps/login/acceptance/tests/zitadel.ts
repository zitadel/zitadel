import { Authenticator } from "@otplib/core";
import { createDigest, createRandomBytes } from "@otplib/plugin-crypto";
import { keyDecoder, keyEncoder } from "@otplib/plugin-thirty-two"; // use your chosen base32 plugin
import axios from "axios";
import dotenv from "dotenv";
import { request } from "gaxios";
import path from "path";
import { OtpType, userProps } from "./user";

dotenv.config({ path: path.resolve(__dirname, "../../login/.env.test.local") });

export async function addUser(props: userProps) {
  const body = {
    username: props.email,
    organization: {
      orgId: props.organization,
    },
    profile: {
      givenName: props.firstName,
      familyName: props.lastName,
    },
    email: {
      email: props.email,
      isVerified: true,
    },
    phone: {
      phone: props.phone,
      isVerified: true,
    },
    password: {
      password: props.password,
      changeRequired: props.passwordChangeRequired ?? false,
    },
  };
  if (!props.isEmailVerified) {
    delete body.email.isVerified;
  }
  if (!props.isPhoneVerified) {
    delete body.phone.isVerified;
  }

  return await listCall(`${process.env.ZITADEL_API_URL}/v2/users/human`, body);
}

export async function removeUserByUsername(username: string) {
  const resp = await getUserByUsername(username);
  if (!resp || !resp.result || !resp.result[0]) {
    return;
  }
  await removeUser(resp.result[0].userId);
}

export async function removeUser(id: string) {
  await deleteCall(`${process.env.ZITADEL_API_URL}/v2/users/${id}`);
}

async function deleteCall(url: string) {
  try {
    const response = await axios.delete(url, {
      headers: {
        Authorization: `Bearer ${process.env.ZITADEL_ADMIN_TOKEN}`,
      },
    });

    if (response.status >= 400 && response.status !== 404) {
      const error = `HTTP Error: ${response.status} - ${response.statusText}`;
      console.error(error);
      throw new Error(error);
    }
  } catch (error) {
    console.error("Error making request:", error);
    throw error;
  }
}

export async function getUserByUsername(username: string): Promise<any> {
  const listUsersBody = {
    queries: [
      {
        userNameQuery: {
          userName: username,
        },
      },
    ],
  };

  return await listCall(`${process.env.ZITADEL_API_URL}/v2/users`, listUsersBody);
}

async function listCall(url: string, data: any): Promise<any> {
  try {
    const response = await axios.post(url, data, {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${process.env.ZITADEL_ADMIN_TOKEN}`,
      },
    });

    if (response.status >= 400) {
      const error = `HTTP Error: ${response.status} - ${response.statusText}`;
      console.error(error);
      throw new Error(error);
    }

    return response.data;
  } catch (error) {
    console.error("Error making request:", error);
    throw error;
  }
}

export async function activateOTP(userId: string, type: OtpType) {
  let url = "otp_";
  switch (type) {
    case OtpType.sms:
      url = url + "sms";
      break;
    case OtpType.email:
      url = url + "email";
      break;
  }

  await pushCall(`${process.env.ZITADEL_API_URL}/v2/users/${userId}/${url}`, {});
}

async function pushCall(url: string, data: any) {
  try {
    const response = await axios.post(url, data, {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${process.env.ZITADEL_ADMIN_TOKEN}`,
      },
    });

    if (response.status >= 400) {
      const error = `HTTP Error: ${response.status} - ${response.statusText}`;
      console.error(error);
      throw new Error(error);
    }
  } catch (error) {
    console.error("Error making request:", error);
    throw error;
  }
}

export async function addTOTP(userId: string): Promise<string> {
  const response = await listCall(`${process.env.ZITADEL_API_URL}/v2/users/${userId}/totp`, {});
  const code = totp(response.secret);
  await pushCall(`${process.env.ZITADEL_API_URL}/v2/users/${userId}/totp/verify`, { code: code });
  return response.secret;
}

export function totp(secret: string) {
  const authenticator = new Authenticator({
    createDigest,
    createRandomBytes,
    keyDecoder,
    keyEncoder,
  });
  // google authenticator usage
  const token = authenticator.generate(secret);

  // check if token can be used
  if (!authenticator.verify({ token: token, secret: secret })) {
    const error = `Generated token could not be verified`;
    console.error(error);
    throw new Error(error);
  }

  return token;
}

export async function eventualNewUser(id: string) {
  return request({
    url: `${process.env.ZITADEL_API_URL}/v2/users/${id}`,
    method: "GET",
    headers: {
      Authorization: `Bearer ${process.env.ZITADEL_ADMIN_TOKEN}`,
      "Content-Type": "application/json",
    },
    retryConfig: {
      statusCodesToRetry: [[404, 404]],
      retry: Number.MAX_SAFE_INTEGER, // totalTimeout limits the number of retries
      totalTimeout: 10000, // 10 seconds
      onRetryAttempt: (error) => {
        console.warn(`Retrying to query new user ${id}: ${error.message}`);
      },
    },
  });
}
