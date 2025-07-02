import { Gaxios, GaxiosResponse } from "gaxios";

const awaitNotification = new Gaxios({
  url: process.env.SINK_NOTIFICATION_URL,
  method: "POST",
  retryConfig: {
    httpMethodsToRetry: ["POST"],
    statusCodesToRetry: [[404, 404]],
    retry: Number.MAX_SAFE_INTEGER, // totalTimeout limits the number of retries
    totalTimeout: 10000, // 10 seconds
    onRetryAttempt: (error) => {
      console.warn(`Retrying request to sink notification service: ${error.message}`);
    },
  },
});

export async function getOtpFromSink(recipient: string): Promise<any> {
  return awaitNotification.request({ data: { recipient } }).then((response) => {
    expectSuccess(response);
    const otp = response?.data?.args?.otp;
    if (!otp) {
      throw new Error(`Response does not contain an otp property: ${JSON.stringify(response.data, null, 2)}`);
    }
    return otp;
  });
}

export async function getCodeFromSink(recipient: string): Promise<any> {
  return awaitNotification.request({ data: { recipient } }).then((response) => {
    expectSuccess(response);
    const code = response?.data?.args?.code;
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
