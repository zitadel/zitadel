import { Authenticator } from "@otplib/core";
import { createDigest, createRandomBytes } from "@otplib/plugin-crypto";
import { keyDecoder, keyEncoder } from "@otplib/plugin-thirty-two"; // use your chosen base32 plugin
import axios from "axios";
import { request } from "gaxios";
import { OtpType, userProps } from "./user";
import { Config } from "./config";

export async function addUser(props: userProps, cfg: Config) {
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
      isVerified: props.isEmailVerified || undefined,
    },
    phone: {
      phone: props.phone,
      isVerified: props.isPhoneVerified || undefined,
    },
    password: {
      password: props.password,
      changeRequired: props.passwordChangeRequired ?? false,
    },
  };
  return await listCall(`${cfg.zitadelApiUrl}/v2/users/human`, body, cfg);
}

export async function removeUserByUsername(username: string, cfg: Config) {
  const resp = await getUserByUsername(username, cfg);
  if (!resp || !resp.result || !resp.result[0]) {
    return;
  }
  await removeUser(resp.result[0].userId, cfg);
}

export async function removeUser(id: string, cfg: Config) {
  await deleteCall(`${cfg.zitadelApiUrl}/v2/users/${id}`, cfg);
}

async function deleteCall(url: string, cfg: Config) {
  try {
    const response = await axios.delete(url, {
      headers: {
        Authorization: `Bearer ${cfg.adminToken}`,
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

export async function getUserByUsername(username: string, cfg: Config): Promise<any> {
  const listUsersBody = {
    queries: [
      {
        userNameQuery: {
          userName: username,
        },
      },
    ],
  };

  return await listCall(`${cfg.zitadelApiUrl}/v2/users`, listUsersBody, cfg);
}

async function listCall(url: string, data: any, cfg: Config): Promise<any> {
  try {
    const response = await axios.post(url, data, {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${cfg.adminToken}`,
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

export async function activateOTP(userId: string, type: OtpType, cfg: Config) {
  let url = "otp_";
  switch (type) {
    case OtpType.sms:
      url = url + "sms";
      break;
    case OtpType.email:
      url = url + "email";
      break;
  }

  await pushCall(`${cfg.zitadelApiUrl}/v2/users/${userId}/${url}`, {}, cfg);
}

async function pushCall(url: string, data: any, cfg: Config) {
  try {
    const response = await axios.post(url, data, {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${cfg.adminToken}`,
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

export async function addTOTP(userId: string, cfg: Config): Promise<string> {
  const response = await listCall(`${cfg.zitadelApiUrl}/v2/users/${userId}/totp`, {}, cfg);
  const code = totp(response.secret);
  await pushCall(`${cfg.zitadelApiUrl}/v2/users/${userId}/totp/verify`, { code: code }, cfg);
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

export async function eventualNewUser(id: string, cfg: Config) {
  return request({
    url: `${cfg.zitadelApiUrl}/v2/users/${id}`,
    method: "GET",
    headers: {
      Authorization: `Bearer ${cfg.adminToken}`,
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
