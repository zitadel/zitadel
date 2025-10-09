import { Gaxios, GaxiosResponse } from "gaxios";

// With exponential backoff and default retry settings, 6 attempts should take about 6-7 minutes
const NOTIFICATION_RETRY_COUNT = 6;

const awaitNotification = new Gaxios({
  method: "GET",
  retryConfig: {
    httpMethodsToRetry: ["GET"],
    statusCodesToRetry: [[404, 404]],
    retry: NOTIFICATION_RETRY_COUNT,
    onRetryAttempt: (error) => {
      console.warn(`Retrying request to sink notification service: ${error.message}`);
    },
  },
});

export async function eventualSMSOTP(number: string): Promise<any> {
  return awaitNotification.request({ url: `${process.env.MOCK_URL}/notifications/sms/${number}` }).then((response) => {
    expectSuccess(response);
    const otp = response?.data?.args?.oTP;
    if (!otp) {
      throw new Error(`Response does not contain an otp property: ${JSON.stringify(response.data, null, 2)}`);
    }
    return otp;
  });
}

export type OTPResponseProperty = "code" | "oTP";

export async function eventualEmailOTP(recipient: string, property: OTPResponseProperty = "code"): Promise<any> {
  return awaitNotification.request({ url: `${process.env.MOCK_URL}/notifications/email/${recipient}` }).then((response) => {
    expectSuccess(response);
    const code = response?.data?.args[property]
    if (!code) {
      throw new Error(`Response does not contain a code property: ${JSON.stringify(response.data, null, 2)}`);
    }
    return code;
  });
}

function expectSuccess(response: GaxiosResponse): void {
  if (response.status !== 200) {
    throw new Error(`Expected HTTP status 200, but got: ${response.status} - ${response.statusText}`);
  }
}
