import { Gaxios, GaxiosResponse } from "gaxios";

const awaitNotification = new Gaxios({
  method: "GET",
  retryConfig: {
    httpMethodsToRetry: ["GET"],
    statusCodesToRetry: [[404, 404]],
    retry: 6,
    onRetryAttempt: (error) => {
      console.warn(`Retrying request to sink notification service: ${error.message}`);
    },
  },
});

export async function eventualSMSOTP(recipient: string): Promise<any> {
  return awaitNotification.request({ url: `${process.env.SINK_NOTIFICATION_URL}/${recipient}` }).then((response) => {
    expectSuccess(response);
    const otp = response?.data?.args?.otp;
    if (!otp) {
      throw new Error(`Response does not contain an otp property: ${JSON.stringify(response.data, null, 2)}`);
    }
    return otp;
  });
}

export async function eventualEmailOTP(recipient: string): Promise<any> {
  return awaitNotification.request({ url: `${process.env.SINK_NOTIFICATION_URL}/${recipient}` }).then((response) => {
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
