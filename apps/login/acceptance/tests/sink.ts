import { Gaxios, GaxiosResponse } from "gaxios";
import { Config } from "./config";

const awaitNotification = (cfg: Config, since: Date) => new Gaxios({
  url: `${cfg.sinkNotificationUrl}?since=${since.toISOString()}`,
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

export async function getOtpFromSink(cfg: Config, recipient: string, since: Date): Promise<any> {
  console.log(`Awaiting notification from url ${cfg.sinkNotificationUrl} for recipient ${recipient}`);
  return awaitNotification(cfg, since).request({ data: { recipient } }).then((response) => {
    expectSuccess(response);
    const otp = response?.data?.args?.oTP;
    if (!otp) {
      throw new Error(`Response does not contain an otp property: ${JSON.stringify(response.data, null, 2)}`);
    }
    return otp;
  });
}

export async function getCodeFromSink(cfg: Config, recipient: string, since: Date): Promise<any> {
  console.log(`Awaiting notification from url ${cfg.sinkNotificationUrl} for recipient ${recipient}`);
  return awaitNotification(cfg, since).request({ data: { recipient } }).then((response) => {
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
