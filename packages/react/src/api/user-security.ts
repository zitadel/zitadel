"use server";

import { create, toJson } from "@zitadel/client";
import {
  ListAuthenticationMethodTypesRequestSchema,
  ListAuthenticationMethodTypesResponseSchema,
  ListAuthenticationFactorsRequestSchema,
  ListAuthenticationFactorsResponseSchema,
  ListPasskeysRequestSchema,
  ListPasskeysResponseSchema,
  RemovePasskeyRequestSchema,
  RemoveTOTPRequestSchema,
  RemoveOTPSMSRequestSchema,
  RemoveOTPEmailRequestSchema,
  RemoveU2FRequestSchema,
  UserService,
} from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { createClient } from "@connectrpc/connect";
import { getTransport } from "./transport";

/**
 * List configured authentication method types for a user.
 * Returns an array of enum strings like "AUTHENTICATION_METHOD_TYPE_PASSWORD", etc.
 */
export async function listAuthMethodTypes(userId: string): Promise<string[]> {
  const transport = getTransport();
  const client = createClient(UserService, transport);
  const request = create(ListAuthenticationMethodTypesRequestSchema, { userId });
  const response = await client.listAuthenticationMethodTypes(request);
  const json = toJson(ListAuthenticationMethodTypesResponseSchema, response) as any;
  return json.authMethodTypes ?? [];
}

/**
 * List authentication factors (MFA) for a user with state and type details.
 * Returns factors like OTP (TOTP), U2F, OTP SMS, OTP Email with their state.
 */
export async function listAuthFactors(userId: string): Promise<any[]> {
  const transport = getTransport();
  const client = createClient(UserService, transport);
  const request = create(ListAuthenticationFactorsRequestSchema, { userId });
  const response = await client.listAuthenticationFactors(request);
  const json = toJson(ListAuthenticationFactorsResponseSchema, response) as any;
  return json.result ?? [];
}

/**
 * List passkeys registered for a user.
 */
export async function listPasskeys(userId: string): Promise<any[]> {
  const transport = getTransport();
  const client = createClient(UserService, transport);
  const request = create(ListPasskeysRequestSchema, { userId });
  const response = await client.listPasskeys(request);
  const json = toJson(ListPasskeysResponseSchema, response) as any;
  return json.result ?? [];
}

/**
 * Remove a passkey from a user.
 */
export async function removePasskey(userId: string, passkeyId: string) {
  const transport = getTransport();
  const client = createClient(UserService, transport);
  const request = create(RemovePasskeyRequestSchema, { userId, passkeyId });
  await client.removePasskey(request);
}

/**
 * Remove TOTP factor from a user.
 */
export async function removeTOTP(userId: string) {
  const transport = getTransport();
  const client = createClient(UserService, transport);
  const request = create(RemoveTOTPRequestSchema, { userId });
  await client.removeTOTP(request);
}

/**
 * Remove OTP SMS factor from a user.
 */
export async function removeOTPSMS(userId: string) {
  const transport = getTransport();
  const client = createClient(UserService, transport);
  const request = create(RemoveOTPSMSRequestSchema, { userId });
  await client.removeOTPSMS(request);
}

/**
 * Remove OTP Email factor from a user.
 */
export async function removeOTPEmail(userId: string) {
  const transport = getTransport();
  const client = createClient(UserService, transport);
  const request = create(RemoveOTPEmailRequestSchema, { userId });
  await client.removeOTPEmail(request);
}

/**
 * Remove U2F factor from a user.
 */
export async function removeU2F(userId: string, u2fId: string) {
  const transport = getTransport();
  const client = createClient(UserService, transport);
  const request = create(RemoveU2FRequestSchema, { userId, u2fId });
  await client.removeU2F(request);
}
